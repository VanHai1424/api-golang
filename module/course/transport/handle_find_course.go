package coursetrpt

import (
	coursebiz "crawdata/module/course/business"
	coursestorage "crawdata/module/course/storage"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// @Summary Lấy thông tin khóa học
// @Description Lấy thông tin khóa học
// @tag courses
// @Accept json
// @Produce json
// @Param id path int true "ID khóa học"
// @Success 200 {object} coursemodel.Course
// @Failure 500 {string} string "Lỗi khi lấy thông tin khóa học"
// @Router /courses/{id} [get]
func HandleFindCourse(db *sql.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")

		storage := coursestorage.NewMySQLStorage(db)
		biz := coursebiz.NewFindCourseBiz(storage)

		course, err := biz.FindCourse(id)
		if err != nil {
			return ctx.Status(500).SendString("Lỗi khi lấy thông tin khóa học")
		}

		return ctx.JSON(course)
	}
}
