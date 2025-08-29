package userstorage

import usermodel "crawdata/module/user/model"

func (s *MySQLStorage) GetAllUser(page, pageSize int) ([]usermodel.User, error) {
	var users []usermodel.User
	// Tính offset
	offset := (page - 1) * pageSize

	if err := s.db.Order("id DESC").Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
