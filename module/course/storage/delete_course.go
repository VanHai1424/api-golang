package coursestorage

import coursemodel "crawdata/module/course/model"

func (s *MySQLStorage) Delete(id string) error {
	return s.db.Delete(&coursemodel.Course{}, id).Error
}
