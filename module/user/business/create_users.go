package userbiz

import (
	usermodel "crawdata/module/user/model"
	"sync"
)

type CreateUsersStorage interface {
	InsertMany(users []usermodel.User) error
}

type createUsersBiz struct {
	store CreateUsersStorage
}

func NewCreateUsersBiz(storage CreateUsersStorage) *createUsersBiz {
	return &createUsersBiz{
		store: storage,
	}
}

func (biz *createUsersBiz) CreateNewUsers(users []usermodel.User) error {
	if len(users) == 0 {
		return nil
	}

	batches := splitIntoBatches(users, 1000)
	batchChan := make(chan []usermodel.User, len(batches))
	errChan := make(chan error, 100)
	var wg sync.WaitGroup

	// Gửi tất cả batch vào channel
	for _, batch := range batches {
		batchChan <- batch
	}
	close(batchChan)

	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range batchChan {
				if err := biz.store.InsertMany(batch); err != nil {
					errChan <- err
					return
				}
			}
		}()
	}

	wg.Wait()
	close(errChan)

	// Trả về lỗi đầu tiên nếu có
	for err := range errChan {
		return err
	}

	return nil
}

// Tách danh sách user thành các batch
func splitIntoBatches(users []usermodel.User, batchSize int) [][]usermodel.User {
	var batches [][]usermodel.User

	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}
		batches = append(batches, users[i:end])
	}
	return batches
}
