package coursestorage

import coursemodel "crawdata/module/course/model"

func (s *MySQLStorage) Update(id string, course *coursemodel.Course) error {
	query := "UPDATE courses SET title = ?, `desc` = ?, version = ?, `update` = ?, developer = ?, category = ?, playId = ?, downloads = ?, link = ?, titleCate = ?"
	args := []interface{}{course.Title, course.Desc, course.Version, course.Update, course.Developer, course.Category, course.PlayId, course.Downloads, course.Link, course.TitleCate}

	if course.Img != "" {
		query += ", img = ?"
		args = append(args, course.Img)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	_, err := s.db.Exec(query, args...)
	return err
}
