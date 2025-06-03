package main

import (
	coursetrpt "crawdata/module/course/transport"
	"crawdata/pkg/db"
	"crawdata/pkg/migration"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Định nghĩa flag migration
	migrationFlag := flag.Bool("migration", false, "Chạy migration và insert data mẫu")
	flag.Parse()

	sqlConfig := db.Sql{
		Host:     "localhost",
		User:     "root",
		Password: "",
		Port:     "3306",
		Database: "test",
	}

	// Khởi tạo kết nối DB
	err := sqlConfig.InitDB()
	if err != nil {
		log.Fatal("Kết nối DB không thành công: ", err)
	}

	if *migrationFlag {
		fmt.Println("Chạy migration và insert data mẫu...")

		// Chạy migrate
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

	// Nếu không có flag migration thì chạy server bình thường
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Static("/public", "./public")

	app.Static("/docs", "./docs")
	app.Static("/docs/oas.json", "./docs/oas.json")

	courseGroup := app.Group("/api/courses")

	courseGroup.Get("/", coursetrpt.HandleListCourse(sqlConfig.Db))
	courseGroup.Get("/:id", coursetrpt.HandleFindCourse(sqlConfig.Db))
	courseGroup.Post("/", coursetrpt.HandleCreateCourse(sqlConfig.Db))
	courseGroup.Put("/:id", coursetrpt.HandleUpdateCourse(sqlConfig.Db))
	courseGroup.Delete("/:id", coursetrpt.HandleDeleteCourse(sqlConfig.Db))

	fmt.Println("Server đang chạy tại http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
