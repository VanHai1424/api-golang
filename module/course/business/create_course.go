package coursebiz

import (
	coursemodel "crawdata/module/course/model"
)

type CreateCourseStorage interface {
	Insert(course *coursemodel.Course) error
}

type createCourseBiz struct {
	store CreateCourseStorage
}

func NewCreateCourseBiz(store CreateCourseStorage) *createCourseBiz {
	return &createCourseBiz{store: store}
}

func (biz *createCourseBiz) CreateNewCourse(course *coursemodel.Course) error {
	err := biz.store.Insert(course)
	if err != nil {
		return err
	}
	return nil
}
