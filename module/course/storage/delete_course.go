package coursestorage

func (s *MySQLStorage) Delete(id string) error {
	query := "DELETE FROM courses WHERE id = ?"
	_, err := s.db.Exec(query, id)
	return err
}
