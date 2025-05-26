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

// @Summary Tạo khóa học mới
// @Description Tạo khóa học mới với thông tin từ form data
// @Tags courses
// @Accept multipart/form-data
// @Produce json
// @Param title formData string false "Tiêu đề"
// @Param version formData string false "Phiên bản"
// @Param titleCate formData string false "Loại" Enums(Phổ biến, Miễn Phí, Mới cập nhật) default(Phổ biến)
// @Param category formData string false "Danh mục"
// @Param developer formData string false "Nhà phát triển"
// @Param desc formData string false "Mô tả"
// @Param playId formData string false "ID phát hành"
// @Param img formData file false "Ảnh đại diện"
// @Success 201 {object} coursemodel.Course
// @Failure 400 {string} string "Lỗi khi phân tích dữ liệu"
// @Failure 500 {string} string "Lỗi khi thêm khóa học"
// @Router /courses [post]
func HandleCreateCourse(db *sql.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		course := new(coursemodel.Course)

		if err := ctx.BodyParser(course); err != nil {
			return ctx.Status(400).SendString("Lỗi khi phân tích dữ liệu")
		}
		now := time.Now().Format("2 thg 1, 2006") // định dạng ngày tháng theo yêu cầu

		if course.Update == "" {
			course.Update = now
		}

		if course.Downloads == "" {
			course.Downloads = "1.000.000+"
		}

		if course.TitleCate == "" {
			course.TitleCate = course.Type
		}

		// 2. Lấy file ảnh từ form data (key 'img')
		file, err := ctx.FormFile("img")
		if err == nil {
			// Tạo thư mục lưu file nếu chưa có
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
		biz := coursebiz.NewCreateCourseBiz(storage)

		if err := biz.CreateNewCourse(course); err != nil {
			return ctx.Status(500).SendString("Lỗi khi thêm khóa học")
		}

		return ctx.Status(201).JSON(course)
	}
}
