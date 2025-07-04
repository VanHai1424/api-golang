package coursetrpt

import (
	coursebiz "crawdata/module/course/business"
	coursemodel "crawdata/module/course/model"
	coursestorage "crawdata/module/course/storage"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// @Summary Tạo khóa học mới
// @Description Tạo khóa học mới với thông tin từ form data
// @tag courses
// @Accept multipart/form-data
// @Produce json
// @Param title form string false "Tiêu đề"
// @Param version form string false "Phiên bản"
// @Param titleCate form string false "Loại khóa học" enums(Phổ biến,Miễn phí,Mới cập nhật)
// @Param category form string false "Danh mục"
// @Param developer form string false "Nhà phát triển"
// @Param desc form string false "Mô tả"
// @Param playId form string false "ID phát hành"
// @Param img file ignored false "Ảnh đại diện"
// @Success 201 {object} coursemodel.Course
// @Failure 400 {string} string "Lỗi khi phân tích dữ liệu"
// @Failure 500 {string} string "Lỗi khi thêm khóa học"
// @Router /courses [post]
func HandleCreateCourse(db *gorm.DB) fiber.Handler {
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

		if course.TitleCate == "" && course.Type != "" {
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
