package coursestorage

import (
	coursemodel "crawdata/module/course/model"
	"time"
)

func (s *MySQLStorage) Insert(course *coursemodel.Course) error {
	now := time.Now().In(vnLoc)

	query := "INSERT INTO courses(title, img, `desc`, version, `update`, developer, category, playId, downloads, link, titleCate, created_at, updated_at) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	_, err := s.db.Exec(query,
		course.Title, course.Img, course.Desc, course.Version, course.Update,
		course.Developer, course.Category, course.PlayId, course.Downloads,
		course.Link, course.TitleCate, now, now)

	return err
}
