// @Version 1.0.0
// @Title Course API
// @Description API Course.
// @Server http://localhost:8080/api

package main

import (
	"context"
	coursetrpt "crawdata/module/course/transport"
	"crawdata/pkg/db"
	"crawdata/pkg/migration"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"crawdata/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

// Google Performance Example
type GooglePerformance struct {
	Id            int    `gorm:"primaryKey" json:"id"`
	Mcc           string `json:"mcc"`
	Status        string `json:"status"`
	Placement     string `json:"placement "`
	PlacementType string `json:"placement_type"`
	Impression    int    `json:"impression"`
	Date          string `json:"date"`
	Flag          int    `json:"flag"`
	Ads           string `json:"ads"`
}

func (GooglePerformance) TableName() string {
	return "google_performance"
}

type Task struct {
	Domain string
}

// Result struct để lưu kết quả xử lý
type ProcessResult struct {
	Domain  string
	Content string
	Error   error
}

// BatchResult để lưu kết quả của 1 batch
type BatchResult struct {
	Domain  string
	Content string
}

// Tối ưu HTTP client với connection pooling
var optimizedClient = &http.Client{
	Timeout: 10 * time.Second, // Tăng timeout
	Transport: &http.Transport{
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 50,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
		MaxConnsPerHost:     100,
	},
}

func fetchAdsTxtOptimized(domain string) (string, error) {
	url := "https://" + domain + "/ads.txt"

	// Tạo request với context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	// Set headers để giảm response size nếu có thể
	req.Header.Set("User-Agent", "AdsTxt-Crawler/1.0")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	resp, err := optimizedClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Giới hạn kích thước đọc để tránh memory leak
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024)) // 1MB max
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Worker xử lý - đã bỏ retry logic
func optimizedWorker(id int, tasks <-chan string, results chan<- ProcessResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for domain := range tasks {
		content, err := fetchAdsTxtOptimized(domain)

		results <- ProcessResult{
			Domain:  domain,
			Content: content,
			Error:   err,
		}
	}
}

// Batch update DB để giảm số lần truy vấn
func batchUpdateDB(db *gorm.DB, results []BatchResult) error {
	if len(results) == 0 {
		return nil
	}

	// Sử dụng transaction để đảm bảo consistency
	return db.Transaction(func(tx *gorm.DB) error {
		for _, result := range results {
			err := tx.Model(&GooglePerformance{}).
				Where("placement = ?", result.Domain).
				Update("ads", result.Content).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func handleOptimizedBatchUpdate(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		const (
			batchSize  = 50000 // Tăng batch size
			numWorkers = 100   // Giảm số workers
			maxResults = 10000 // Buffer cho results channel
		)

		// Đếm tổng số domain unique
		var total int64
		if err := db.Model(&GooglePerformance{}).
			Select("COUNT(DISTINCT placement)").
			Where("placement != '' AND placement IS NOT NULL").
			Count(&total).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Không thể đếm domain: " + err.Error(),
			})
		}

		log.Printf("🚀 Bắt đầu xử lý %d domain unique", total)

		numBatches := int((total + int64(batchSize) - 1) / int64(batchSize))
		successCount := 0
		errorCount := 0

		for batchIndex := 0; batchIndex < numBatches; batchIndex++ {
			batchStartTime := time.Now()
			log.Printf("📦 Bắt đầu xử lý batch %d/%d", batchIndex+1, numBatches)

			// Lấy domains theo batch
			var domains []string
			offset := batchIndex * batchSize

			if err := db.Model(&GooglePerformance{}).
				Select("DISTINCT placement").
				Where("placement != '' AND placement IS NOT NULL").
				Limit(1000).
				Offset(offset).
				Limit(1000).
				Pluck("placement", &domains).Error; err != nil {
				log.Printf("❌ Lỗi lấy domain batch %d: %v", batchIndex+1, err)
				continue
			}

			if len(domains) == 0 {
				continue
			}

			// Tạo channels
			taskChan := make(chan string, len(domains))
			resultChan := make(chan ProcessResult, maxResults)

			// Khởi động workers
			var workerWg sync.WaitGroup
			actualWorkers := numWorkers
			if len(domains) < numWorkers {
				actualWorkers = len(domains)
			}

			for w := 0; w < actualWorkers; w++ {
				workerWg.Add(1)
				go optimizedWorker(w, taskChan, resultChan, &workerWg)
			}

			// Gửi tasks
			for _, domain := range domains {
				taskChan <- domain
			}
			close(taskChan)

			// Thu thập kết quả trong goroutine riêng
			var resultWg sync.WaitGroup
			var batchResults []BatchResult
			var mu sync.Mutex

			resultWg.Add(1)
			go func() {
				defer resultWg.Done()
				processedCount := 0

				for result := range resultChan {
					processedCount++

					if result.Error != nil {
						errorCount++
						log.Printf("❌ [%s] %v", result.Domain, result.Error)
					} else {
						successCount++
						mu.Lock()
						batchResults = append(batchResults, BatchResult{
							Domain:  result.Domain,
							Content: result.Content,
						})
						mu.Unlock()
					}

					// Log progress mỗi 100 domain để theo dõi chi tiết hơn
					if processedCount%100 == 0 {
						log.Printf("⏳ Batch %d/%d: Đã xử lý %d/%d domain (Success: %d, Error: %d)",
							batchIndex+1, numBatches, processedCount, len(domains), successCount, errorCount)
					}

					if processedCount >= len(domains) {
						break
					}
				}

				// Log kết quả cuối batch
				log.Printf("🎯 Hoàn thành batch %d/%d: %d/%d domain (Success: %d, Error: %d)",
					batchIndex+1, numBatches, processedCount, len(domains), successCount, errorCount)
			}()

			// Đợi workers hoàn thành
			workerWg.Wait()
			close(resultChan)

			// Đợi result collector hoàn thành
			resultWg.Wait()

			// Batch update database
			if len(batchResults) > 0 {
				updateStart := time.Now()

				// Chia nhỏ batch update để tránh lock DB quá lâu
				chunkSize := 1000
				for i := 0; i < len(batchResults); i += chunkSize {
					end := i + chunkSize
					if end > len(batchResults) {
						end = len(batchResults)
					}

					chunk := batchResults[i:end]
					if err := batchUpdateDB(db, chunk); err != nil {
						log.Printf("❌ Lỗi update DB chunk %d-%d: %v", i, end, err)
					}
				}

				log.Printf("💾 Cập nhật %d records vào DB trong %v",
					len(batchResults), time.Since(updateStart))
			}

			elapsed := time.Since(startTime)
			batchElapsed := time.Since(batchStartTime)
			log.Printf("✅ Hoàn thành batch %d/%d trong %v - Tổng Success: %d, Tổng Error: %d, Tổng thời gian: %v",
				batchIndex+1, numBatches, batchElapsed, successCount, errorCount, elapsed)
		}

		totalTime := time.Since(startTime)
		avgDomainsPerSecond := float64(successCount+errorCount) / totalTime.Seconds()

		log.Printf("🏁 KẾT THÚC XỬ LÝ: Tổng %d domain, Success: %d, Error: %d, Thời gian: %v, Tốc độ: %.2f domain/giây",
			total, successCount, errorCount, totalTime, avgDomainsPerSecond)

		return c.JSON(fiber.Map{
			"message":            "Hoàn thành xử lý",
			"total_domains":      total,
			"success_count":      successCount,
			"error_count":        errorCount,
			"total_time":         totalTime.String(),
			"domains_per_second": avgDomainsPerSecond,
		})
	}
}

func main() {
	// Định nghĩa flag migration
	migrationFlag := flag.Bool("migration", false, "Chạy migration và insert data mẫu")
	flag.Parse()

	// Load cấu hình
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Load config lỗi: ", err)
	}

	// Khởi tạo kết nối DB dựa trên config
	sqlConfig := db.Sql{
		Host:     cfg.DB.Host,
		User:     cfg.DB.User,
		Password: cfg.DB.Pass,
		Port:     cfg.DB.Port,
		Database: cfg.DB.Name,
	}

	err = sqlConfig.InitDB()
	if err != nil {
		log.Fatal("Kết nối DB không thành công: ", err)
	}

	// Nếu có flag migration thì chạy migration + insert data mẫu
	if *migrationFlag {
		fmt.Println("Chạy migration và insert data mẫu...")

		if err = migration.Migrate(sqlConfig.Db); err != nil {
			log.Fatal("Migration lỗi: ", err)
		}

		if err = migration.InsertSampleData(sqlConfig.Db); err != nil {
			log.Fatal("Insert data mẫu lỗi: ", err)
		}

		fmt.Println("Hoàn thành migration và insert data mẫu.")
	}

	// Khởi tạo Fiber app và middleware
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Static("/public", "./public")
	app.Static("/docs", "./docs")

	// Định nghĩa route
	courseGroup := app.Group("/api/courses")
	courseGroup.Get("/", coursetrpt.HandleListCourse(sqlConfig.Db))
	courseGroup.Get("/:id", coursetrpt.HandleFindCourse(sqlConfig.Db))
	courseGroup.Post("/", coursetrpt.HandleCreateCourse(sqlConfig.Db))
	courseGroup.Put("/:id", coursetrpt.HandleUpdateCourse(sqlConfig.Db))
	courseGroup.Delete("/:id", coursetrpt.HandleDeleteCourse(sqlConfig.Db))

	// Route test cập nhật ads.txt
	app.Get("/api/test", handleOptimizedBatchUpdate(sqlConfig.Db))

	// Chạy server
	fmt.Println("Server đang chạy tại http://localhost:8070")
	log.Fatal(app.Listen(":8070"))
}
