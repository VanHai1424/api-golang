package usertrpt

import (
	userbiz "crawdata/module/user/business"
	usermodel "crawdata/module/user/model"
	userstorage "crawdata/module/user/storage"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func HandleCreateUsers(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		users := make([]usermodel.User, 0, 1000000)

		for i := 0; i < 1000000; i++ {
			user := usermodel.User{
				Username: fmt.Sprintf("Username %d", i+1),
				Password: "123456",
			}
			users = append(users, user)
		}

		storage := userstorage.NewMySQLStorage(db)
		biz := userbiz.NewCreateUsersBiz(storage)

		if err := biz.CreateNewUsers(users); err != nil {
			return ctx.Status(500).SendString("Lỗi khi thêm 1tr người dùng")
		}

		return ctx.JSON(fiber.Map{
			"message": "Tạo 1 triệu user fake thành công",
		})
	}
}
