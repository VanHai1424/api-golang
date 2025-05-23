package coursestorage

import coursemodel "crawdata/module/course/model"

func (s *MySQLStorage) GetAllCourses(name string) ([]coursemodel.Course, error) {
	query := "SELECT * FROM courses WHERE 1"
	args := []interface{}{}

	if name != "" {
		query += " AND title LIKE ?"
		args = append(args, "%"+name+"%")
	}

	query += " ORDER BY id DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []coursemodel.Course
	for rows.Next() {
		var course coursemodel.Course
		err := rows.Scan(&course.Id, &course.Title, &course.Img, &course.Desc,
			&course.Version, &course.Update, &course.Developer,
			&course.Category, &course.PlayId, &course.Downloads,
			&course.Link, &course.TitleCate, &course.CreatedAt, &course.UpdatedAt)

		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}
