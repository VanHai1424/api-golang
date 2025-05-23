package coursetrpt

import (
	coursebiz "crawdata/module/course/business"
	coursestorage "crawdata/module/course/storage"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func HandleListCourse(db *sql.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		name := ctx.Query("name")
		startDate := ctx.Query("startDate")
		endDate := ctx.Query("endDate")

		storage := coursestorage.NewMySQLStorage(db)
		biz := coursebiz.NewListCourseBiz(storage)

		courses, err := biz.ListCourse(name, startDate, endDate)
		if err != nil {
			return ctx.Status(500).SendString("Lỗi khi lấy danh sách khóa học")
		}

		return ctx.JSON(courses)
	}
}
