package userbiz

import usermodel "crawdata/module/user/model"

type ListUserStorage interface {
	GetAllUser(page, pageSize int) ([]usermodel.User, error)
}

type listUserBiz struct {
	storage ListUserStorage
}

func NewListUserBiz(storage ListUserStorage) *listUserBiz {
	return &listUserBiz{
		storage: storage,
	}
}

func (biz *listUserBiz) ListUser(page, pageSize int) ([]usermodel.User, error) {
	users, err := biz.storage.GetAllUser(page, pageSize)
	if err != nil {
		return nil, err
	}
	return users, nil
}
