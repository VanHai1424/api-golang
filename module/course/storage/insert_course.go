package coursestorage

import (
	coursemodel "crawdata/module/course/model"
)

func (s *MySQLStorage) Insert(course *coursemodel.Course) error {
	return s.db.Create(course).Error
}
