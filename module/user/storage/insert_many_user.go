package userstorage

import usermodel "crawdata/module/user/model"

func (s *MySQLStorage) InsertMany(users []usermodel.User) error {
	if err := s.db.Create(&users).Error; err != nil {
		return err
	}
	return nil
}
