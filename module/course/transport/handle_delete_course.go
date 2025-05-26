package coursetrpt

import (
	coursebiz "crawdata/module/course/business"
	coursestorage "crawdata/module/course/storage"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// @Summary Xóa khóa học
// @Description Xóa khóa học theo ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path int true "ID khóa học"
// @Success 204 {string} string "No Content"
// @Router /courses/{id} [delete]
func HandleDeleteCourse(db *sql.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		if id == "" {
			return ctx.Status(400).SendString("Thiếu tham số id")
		}

		storage := coursestorage.NewMySQLStorage(db)
		biz := coursebiz.NewDeleteCourseBiz(storage)

		if err := biz.DeleteCourse(id); err != nil {
			return ctx.Status(500).SendString("Lỗi khi xóa khóa học")
		}

		return ctx.SendStatus(204)
	}
}
