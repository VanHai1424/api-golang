// @Version 1.0.0
// @Title Course API
// @Description API Course.
// @Server http://localhost:8080/api

package main

import (
	coursetrpt "crawdata/module/course/transport"
	"crawdata/pkg/db"
	"crawdata/pkg/migration"
	"flag"
	"fmt"
	"log"

	"crawdata/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

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

		err = migration.Migrate(sqlConfig.Db)
		if err != nil {
			log.Fatal("Migration lỗi: ", err)
		}

		err = migration.InsertSampleData(sqlConfig.Db)
		if err != nil {
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
	app.Static("/docs/oas.json", "./docs/oas.json")

	// Định nghĩa route
	courseGroup := app.Group("/api/courses")
	courseGroup.Get("/", coursetrpt.HandleListCourse(sqlConfig.Db))
	courseGroup.Get("/:id", coursetrpt.HandleFindCourse(sqlConfig.Db))
	courseGroup.Post("/", coursetrpt.HandleCreateCourse(sqlConfig.Db))
	courseGroup.Put("/:id", coursetrpt.HandleUpdateCourse(sqlConfig.Db))
	courseGroup.Delete("/:id", coursetrpt.HandleDeleteCourse(sqlConfig.Db))

	fmt.Println("Server đang chạy tại http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
