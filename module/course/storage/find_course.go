package coursestorage

import coursemodel "crawdata/module/course/model"

func (s *MySQLStorage) GetCourseById(id string) (coursemodel.Course, error) {
	var course coursemodel.Course
	err := s.db.First(&course, "id = ?", id).Error

	if err != nil {
		return coursemodel.Course{}, err
	}
	return course, nil
}
