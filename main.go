package main

import (
	coursetrpt "crawdata/module/course/transport"
	"crawdata/pkg/db"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	sqlConfig := db.Sql{
		Host:     "localhost",
		User:     "root",
		Password: "",
		Port:     "3306",
		Database: "test",
	}

	// Khởi tạo kết nối cơ sở dữ liệu
	err := sqlConfig.InitDB()
	if err != nil {
		log.Fatal("Kết nối không thành công: ", err)
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Cấu hình thư mục public thành static folder, truy cập bằng đường dẫn /public/...
	app.Static("/public", "./public")

	courseGroup := app.Group("/api/courses")

	courseGroup.Get("/", coursetrpt.HandleListCourse(sqlConfig.Db))
	courseGroup.Get("/:id", coursetrpt.HandleFindCourse(sqlConfig.Db))
	courseGroup.Post("/", coursetrpt.HandleCreateCourse(sqlConfig.Db))
	courseGroup.Put("/:id", coursetrpt.HandleUpdateCourse(sqlConfig.Db))
	courseGroup.Delete("/:id", coursetrpt.HandleDeleteCourse(sqlConfig.Db))

	fmt.Println("Server đang chạy tại http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
