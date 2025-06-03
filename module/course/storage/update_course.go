package coursestorage

import coursemodel "crawdata/module/course/model"

func (s *MySQLStorage) Update(id string, course *coursemodel.Course) error {
	// Tạo map chứa các trường cần update
	updates := map[string]interface{}{
		"title":     course.Title,
		"desc":      course.Desc,
		"version":   course.Version,
		"update":    course.Update,
		"developer": course.Developer,
		"category":  course.Category,
		"playId":    course.PlayId,
		"downloads": course.Downloads,
		"link":      course.Link,
		"titleCate": course.TitleCate,
	}

	// Nếu có ảnh mới thì thêm vào map
	if course.Img != "" {
		updates["img"] = course.Img
	}

	err := s.db.Model(&coursemodel.Course{}).Where("id = ?", id).Updates(updates).Error
	return err
}
