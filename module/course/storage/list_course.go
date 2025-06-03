package coursestorage

import coursemodel "crawdata/module/course/model"

func (s *MySQLStorage) GetAllCourses(name string) ([]coursemodel.Course, error) {
	var courses []coursemodel.Course
	db := s.db.Order("id DESC")

	if name != "" {
		db = db.Where("title LIKE ?", "%"+name+"%")
	}

	err := db.Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}
