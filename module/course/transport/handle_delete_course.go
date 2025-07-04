package coursetrpt

import (
	coursebiz "crawdata/module/course/business"
	coursestorage "crawdata/module/course/storage"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// @Summary Xóa khóa học
// @Description Xóa khóa học theo ID
// @tag courses
// @Accept json
// @Produce json
// @Param id path int true "ID khóa học"
// @Success 204 {string} string "No Content"
// @Router /courses/{id} [delete]
func HandleDeleteCourse(db *gorm.DB) fiber.Handler {
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
