package coursestorage

import coursemodel "crawdata/module/course/model"

func (s *MySQLStorage) GetCourseById(id string) (coursemodel.Course, error) {
	query := "SELECT * FROM courses WHERE id = ?"
	row := s.db.QueryRow(query, id)

	var course coursemodel.Course
	err := row.Scan(&course.Id, &course.Title, &course.Img, &course.Desc,
		&course.Version, &course.Update, &course.Developer,
		&course.Category, &course.PlayId, &course.Downloads,
		&course.Link, &course.TitleCate, &course.CreatedAt, &course.UpdatedAt)

	if err != nil {
		return coursemodel.Course{}, err
	}
	return course, nil
}
