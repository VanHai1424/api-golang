package coursetrpt

import (
	coursebiz "crawdata/module/course/business"
	coursemodel "crawdata/module/course/model"
	coursestorage "crawdata/module/course/storage"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// @Summary Cập nhật khóa học
// @Description Cập nhật thông tin khóa học
// @tag courses
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "ID khóa học"
// @Param title form string false "Tiêu đề"
// @Param version form string false "Phiên bản"
// @Param titleCate form string false "Loại" Enums(Phổ biến, Miễn Phí, Mới cập nhật) default(Phổ biến)
// @Param category form string false "Danh mục"
// @Param developer form string false "Nhà phát triển"
// @Param desc form string false "Mô tả"
// @Param playId form string false "ID phát hành"
// @Param img file ignored false "Ảnh đại diện"
// @Success 201 {object} coursemodel.Course
// @Failure 400 {string} string "Lỗi khi phân tích dữ liệu"
// @Failure 500 {string} string "Lỗi khi thêm khóa học"
// @Router /courses/{id} [put]
func HandleUpdateCourse(db *sql.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		if id == "" {
			return ctx.Status(400).SendString("Thiếu tham số id")
		}

		course := new(coursemodel.Course)
		if err := ctx.BodyParser(course); err != nil {
			return ctx.Status(400).SendString("Lỗi khi phân tích dữ liệu")
		}

		now := time.Now().Format("2 thg 1, 2006")

		if course.Update == "" {
			course.Update = now
		}

		if course.Downloads == "" {
			course.Downloads = "1.000.000+"
		}

		if course.TitleCate == "" {
			course.TitleCate = course.Type
		}

		file, err := ctx.FormFile("img")
		if err == nil {
			if err := os.MkdirAll("./public/images", os.ModePerm); err != nil {
				return ctx.Status(500).SendString("Lỗi khi tạo thư mục lưu ảnh")
			}

			filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
			savePath := fmt.Sprintf("./public/images/%s", filename)

			if err := ctx.SaveFile(file, savePath); err != nil {
				return ctx.Status(500).SendString("Lỗi khi lưu ảnh")
			}

			course.Img = filename
		} else {
			course.Img = ""
		}

		storage := coursestorage.NewMySQLStorage(db)
		biz := coursebiz.NewUpdateCourseBiz(storage)

		if err := biz.UpdateCourse(id, course); err != nil {
			return ctx.Status(500).SendString("Lỗi khi cập nhật khóa học")
		}

		return ctx.JSON(course)
	}
}
