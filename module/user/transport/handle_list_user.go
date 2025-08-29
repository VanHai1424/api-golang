package usertrpt

import (
	userbiz "crawdata/module/user/business"
	userstorage "crawdata/module/user/storage"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func HandleListUser(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		pageStr := ctx.Query("page", "1")
		perPageStr := ctx.Query("page_size", "100")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		perPage, err := strconv.Atoi(perPageStr)
		if err != nil || perPage < 1 {
			perPage = 100
		}

		storage := userstorage.NewMySQLStorage(db)
		biz := userbiz.NewListUserBiz(storage)

		users, err := biz.ListUser(page, perPage)
		if err != nil {
			return ctx.Status(500).SendString("Lỗi khi lấy danh sách người dùng")
		}
		return ctx.JSON(users)
	}
}
