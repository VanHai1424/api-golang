package migration

import (
	coursemodel "crawdata/module/course/model"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// Seed data mẫu
func SeedData(db *gorm.DB) error {
	// var count int64
	// db.Model(&coursemodel.Course{}).Count(&count)
	// if count > 0 {
	// 	log.Println("Data already seeded. Skip seeding.")
	// 	return nil
	// }

	imgs := []string{
		"1747967776142227000_Bài 1 - ảnh 1.jpg",
		"1747967825642925200_Add-a-heading.png",
		"1747967834354149600_e-learning-business-continuity-basics.jpg",
		"1747967839457947700_khoa-hoc-javascript.jpg",
		"https://play-lh.googleusercontent.com/iypHl2e5gjAqE0oKbNl2D4MaNMkjKjNuDXusAS1vZvw5i2z3-ezZl3nr6lllMlQov-Q=s75-rw",
		"https://play-lh.googleusercontent.com/mllJKslgVX5eZrhQJbgQRp53EfPmKvop6QOVLt1IcTb9ANyn6zC60aBNZB-zkc0dVZ4=s75-rw",
		"https://play-lh.googleusercontent.com/dySbFLHEUG2ZnKNaLHtBRcz1v-dnRNx-fjFrwK8401Y4ZuOkTZw73N01b7Zqgsl-NA=s75-rw",
		"https://play-lh.googleusercontent.com/eUaosZZxd0viXFgdoCfFsoYleTsLZfIc9il32dFM-hOsF30gfDAC5dcFFJfXQDiRfXw=s75-rw",
		"https://play-lh.googleusercontent.com/SzAOFY3bOy-DlcMRYrlNzD5aTiN4Qv_kknQ3iyEzb5q6KzjekjhDxTWVeD2x4LkUi5E=s75-rw",
		"https://play-lh.googleusercontent.com/mllJKslgVX5eZrhQJbgQRp53EfPmKvop6QOVLt1IcTb9ANyn6zC60aBNZB-zkc0dVZ4=s75-rw",
	}

	// Tạo vài data mẫu
	courses := make([]coursemodel.Course, 0, 10)

	for i := 1; i <= 10; i++ {
		c := coursemodel.Course{
			Title:     fmt.Sprintf("Khóa học Golang cơ bản %d", i),
			Img:       imgs[i-1],
			Desc:      "Học lập trình Golang từ đầu",
			Version:   "1.1.0",
			Update:    time.Now().Format("2 thg 1, 2006"),
			Developer: "Hai",
			Category:  "Programming",
			PlayId:    fmt.Sprintf("go%03d", i),
			Downloads: "500.000+",
			Link:      "https://example.com/course/golang",
			TitleCate: "Programming",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		courses = append(courses, c)
	}

	if err := db.Create(&courses).Error; err != nil {
		return err
	}
	log.Println("Seed data completed.")
	return nil
}
