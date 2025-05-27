package coursetrpt

import (
	coursebiz "crawdata/module/course/business"
	coursestorage "crawdata/module/course/storage"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// @Summary Lấy tất cả khóa học
// @Description Lấy danh sách tất cả khóa học
// @tag courses
// @Accept json
// @Produce json
// @Param name query string false "Tên khóa học"
// @Param startDate query string false "Ngày bắt đầu khóa học" format(date)
// @Param endDate query string false "Ngày kết thúc khóa học" format(date)
// @Success 200 {array} coursemodel.Course
// @Failure 500 {string} string "Lỗi khi lấy danh sách khóa học"
// @Router /courses [get]
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
