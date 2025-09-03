package main

import (
	"context"
	"crawdata/pkg/config"
	"crawdata/pkg/db"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

// ==== Struct Models ====
type (
	Category struct {
		ID   int    `gorm:"primaryKey;autoIncrement"`
		Name string `json:"name"`
	}

	Quiz struct {
		ID            int    `gorm:"primaryKey;autoIncrement"`
		Name          string `json:"name"`
		Desc          string `json:"desc"`
		Thumbnail     string `json:"thumbnail"`
		TotalQuestion int    `json:"total_question"`
		CategoryId    int    `json:"category_id"`
	}

	Question struct {
		ID      int    `gorm:"primaryKey;autoIncrement"`
		Content string `json:"content"`
		QuizId  int    `json:"quiz_id"`
	}

	Answer struct {
		ID         int    `gorm:"primaryKey;autoIncrement"`
		Content    string `json:"content"`
		IsCorrect  bool   `json:"is_correct"`
		QuestionId int    `json:"question_id"`
	}
)

var sqlConfig db.Sql

// ==== Helper Functions ====
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getMap(m map[string]interface{}, key string) map[string]interface{} {
	if v, ok := m[key].(map[string]interface{}); ok {
		return v
	}
	return nil
}

func getSlice(m map[string]interface{}, key string) []interface{} {
	if v, ok := m[key].([]interface{}); ok {
		return v
	}
	return nil
}

// ==== Crawl Quiz ====
func CrawlQuiz(url string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var raw string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(5*time.Second),
		chromedp.Evaluate(`JSON.stringify(window.Mendel)`, &raw),
	); err != nil {
		return fmt.Errorf("lỗi chromedp khi truy cập %s: %w", url, err)
	}

	if raw == "" || raw == "undefined" {
		return fmt.Errorf("không tìm thấy dữ liệu window.Mendel trên trang %s", url)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return fmt.Errorf("lỗi parse JSON từ %s: %w", url, err)
	}

	config := getMap(data, "config")
	if config == nil {
		return fmt.Errorf("không tìm thấy config trong dữ liệu từ %s", url)
	}

	quizCfg := getMap(config, "quiz")
	if quizCfg == nil {
		return fmt.Errorf("không tìm thấy quiz config trong dữ liệu từ %s", url)
	}

	// Xử lý và lưu vào database
	return sqlConfig.Db.Transaction(func(tx *gorm.DB) error {
		// Process Category
		if err := processCategoryAndQuiz(tx, quizCfg, config, url); err != nil {
			return err
		}

		log.Printf("✅ Đã lưu thành công quiz từ: %s", url)
		return nil
	})
}

func processCategoryAndQuiz(tx *gorm.DB, quizCfg, config map[string]interface{}, url string) error {
	// Category
	cats := getSlice(quizCfg, "categories")
	if len(cats) == 0 {
		return fmt.Errorf("không tìm thấy categories trong quiz từ %s", url)
	}

	firstCat, ok := cats[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("format category không đúng trong quiz từ %s", url)
	}

	label := getString(firstCat, "label")
	if label == "" {
		return fmt.Errorf("category label rỗng trong quiz từ %s", url)
	}

	var category Category
	if err := tx.Where("name = ?", label).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			category = Category{Name: label}
			if err := tx.Create(&category).Error; err != nil {
				return fmt.Errorf("lỗi tạo category '%s': %w", label, err)
			}
			log.Printf("🆕 Đã tạo category mới: %s", label)
		} else {
			return fmt.Errorf("lỗi truy vấn category: %w", err)
		}
	}

	// Quiz
	quizTitle := getString(quizCfg, "title")
	if quizTitle == "" {
		return fmt.Errorf("quiz title rỗng từ %s", url)
	}

	questions := getSlice(quizCfg, "questions")
	if len(questions) == 0 {
		return fmt.Errorf("không có questions trong quiz từ %s", url)
	}

	quiz := Quiz{
		Name:          quizTitle,
		Desc:          getString(config, "quizDescription"),
		Thumbnail:     getString(getMap(quizCfg, "image"), "fullUrl"),
		TotalQuestion: len(questions),
		CategoryId:    category.ID,
	}
	if err := tx.Create(&quiz).Error; err != nil {
		return fmt.Errorf("lỗi tạo quiz '%s': %w", quizTitle, err)
	}

	log.Printf("📝 Đã tạo quiz: %s (%d questions)", quizTitle, len(questions))

	// Process Questions and Answers
	return processQuestionsAndAnswers(tx, questions, quiz.ID, url)
}

func processQuestionsAndAnswers(tx *gorm.DB, questions []interface{}, quizID int, url string) error {
	for qIndex, q := range questions {
		qMap, ok := q.(map[string]interface{})
		if !ok {
			return fmt.Errorf("format question %d không đúng", qIndex+1)
		}

		// Question content
		content := getString(qMap, "text")
		if content == "" {
			content = getString(getMap(qMap, "image"), "fullUrl")
		}
		if content == "" {
			return fmt.Errorf("question %d không có content", qIndex+1)
		}

		question := Question{
			Content: content,
			QuizId:  quizID,
		}
		if err := tx.Create(&question).Error; err != nil {
			return fmt.Errorf("lỗi tạo question %d: %w", qIndex+1, err)
		}

		// Process Answers
		if err := processAnswers(tx, qMap, question.ID, qIndex+1); err != nil {
			return err
		}
	}
	return nil
}

func processAnswers(tx *gorm.DB, qMap map[string]interface{}, questionID, qIndex int) error {
	correctAnswerIndex, ok := qMap["correctAnswerIndex"].(float64)
	if !ok {
		return fmt.Errorf("correctAnswerIndex không đúng format trong question %d", qIndex)
	}

	correctIdx := int(correctAnswerIndex)
	answers := getSlice(qMap, "answers")
	if len(answers) == 0 {
		return fmt.Errorf("question %d không có answers", qIndex)
	}

	for ansIndex, ans := range answers {
		answerText, ok := ans.(string)
		if !ok {
			return fmt.Errorf("answer %d trong question %d không phải string", ansIndex+1, qIndex)
		}

		answer := Answer{
			Content:    answerText,
			IsCorrect:  ansIndex == correctIdx,
			QuestionId: questionID,
		}
		if err := tx.Create(&answer).Error; err != nil {
			return fmt.Errorf("lỗi tạo answer %d cho question %d: %w", ansIndex+1, qIndex, err)
		}
	}
	return nil
}

// ==== Get Categories ====
func GetCategories(url string) ([]string, error) {
	log.Printf("🔍 Đang lấy danh sách categories từ: %s", url)

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var html string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second),
		chromedp.OuterHTML(`.quiz-category-select`, &html, chromedp.ByQuery),
	); err != nil {
		return nil, fmt.Errorf("lỗi chromedp khi lấy categories: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("lỗi parse HTML categories: %w", err)
	}

	var links []string
	doc.Find("option").Each(func(_ int, s *goquery.Selection) {
		val, _ := s.Attr("value")
		if val != "-1" && strings.TrimSpace(val) != "" {
			links = append(links, "https://www.britannica.com"+val)
		}
	})

	if len(links) == 0 {
		return nil, fmt.Errorf("không tìm thấy category links nào")
	}

	return links, nil
}

func GetQuizzesByCategory(url string) ([]string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(4*time.Second),
		chromedp.OuterHTML(`#all-quizzes`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		return nil, fmt.Errorf("lỗi chromedp khi truy cập category: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("lỗi parse HTML: %w", err)
	}

	var quizLinks []string
	seen := make(map[string]bool) // Để tránh duplicate

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.HasPrefix(href, "/quiz/") {
			fullURL := fmt.Sprintf("https://www.britannica.com%s", href)
			if !seen[fullURL] {
				seen[fullURL] = true
				quizLinks = append(quizLinks, fullURL)
			}
		}
	})

	if len(quizLinks) == 0 {
		return nil, fmt.Errorf("không tìm thấy quiz links nào")
	}

	log.Printf("📊 Tìm thấy %d quiz links", len(quizLinks))
	return quizLinks, nil
}

// Crawl categories
func FetchAllQuizzes(categories []string) ([]string, error) {
	var allQuizzes []string
	var mu sync.Mutex

	categoryCh := make(chan string)
	wg := sync.WaitGroup{}

	for w := 0; w < 2; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for catURL := range categoryCh {
				log.Printf("📂 [Worker-%d] Đang xử lý category: %s", workerID, catURL)

				quizzes, err := GetQuizzesByCategory(catURL)
				if err != nil {
					log.Printf("❌ Worker-%d lỗi lấy quiz từ category %s: %v", workerID, catURL, err)
					continue
				}

				mu.Lock()
				allQuizzes = append(allQuizzes, quizzes...)
				mu.Unlock()

				log.Printf("✅ Worker-%d lấy được %d quizzes. Tổng hiện tại: %d", workerID, len(quizzes), len(allQuizzes))

				// Delay tránh bị rate limit
				time.Sleep(2 * time.Second)
			}
		}(w + 1)
	}

	// Gửi category vào channel
	go func() {
		for _, cat := range categories {
			categoryCh <- cat
		}
		close(categoryCh)
	}()

	wg.Wait()
	return allQuizzes, nil
}

// Crawl quizzes
func CrawlAllQuizzes(quizLinks []string) int {
	total := len(quizLinks)
	successCount := 0
	var mu sync.Mutex

	quizCh := make(chan string)
	wg := sync.WaitGroup{}

	for w := 0; w < 2; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for quizURL := range quizCh {
				log.Printf("🔄 [Worker-%d] Đang crawl quiz: %s", workerID, quizURL)

				if err := CrawlQuiz(quizURL); err != nil {
					log.Printf("❌ Worker-%d lỗi crawl quiz %s: %v", workerID, quizURL, err)
					continue
				}

				mu.Lock()
				successCount++
				done := successCount
				mu.Unlock()

				log.Printf("✅ [Worker-%d] Đã crawl thành công quiz (%d/%d)", workerID, done, total)

				if done%10 == 0 {
					log.Printf("📊 Tiến độ: %d/%d quizzes (%.1f%%)", done, total, float64(done)/float64(total)*100)
					runtime.GC()
				}

				// Delay để tránh rate limit
				time.Sleep(3 * time.Second)
			}
		}(w + 1)
	}

	go func() {
		for _, quiz := range quizLinks {
			quizCh <- quiz
		}
		close(quizCh)
	}()

	wg.Wait()
	runtime.GC()
	return successCount
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Load config error: ", err)
	}

	sqlConfig = db.Sql{
		Host:     cfg.DB.Host,
		User:     cfg.DB.User,
		Password: cfg.DB.Pass,
		Port:     cfg.DB.Port,
		Database: cfg.DB.Name,
	}

	if err := sqlConfig.InitDB(); err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	// Fiber
	app := fiber.New()
	app.Use(cors.New())

	app.Get("/api/crawl", func(c *fiber.Ctx) error {
		log.Println("🚀 Bắt đầu quá trình crawl...")

		// Lấy categories
		categories, err := GetCategories("https://www.britannica.com/quiz/browse")
		if err != nil {
			log.Fatalf("❌ FATAL: Không thể lấy categories: %v", err)
		}
		log.Printf("✅ Lấy được %d categories", len(categories))

		// Lấy quiz links
		allQuizzes, err := FetchAllQuizzes(categories)
		if err != nil {
			log.Fatalf("❌ FATAL: Lỗi khi lấy quizzes: %v", err)
		}
		log.Printf("📊 Tổng cộng có %d quizzes cần crawl", len(allQuizzes))

		// Crawl quizzes
		successCount := CrawlAllQuizzes(allQuizzes)

		log.Printf("🎉 HOÀN THÀNH! Đã crawl thành công %d/%d quizzes", successCount, len(allQuizzes))

		return c.JSON(fiber.Map{
			"success": true,
			"total":   len(allQuizzes),
			"crawled": successCount,
			"message": "Crawl completed successfully",
		})
	})

	fmt.Println("Server running at http://localhost:8070")
	log.Fatal(app.Listen(":8070"))
}
