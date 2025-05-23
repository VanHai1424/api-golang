package coursebiz

import (
	coursemodel "crawdata/module/course/model"
)

type UpdateCourseStorage interface {
	Update(id string, course *coursemodel.Course) error
}

type updateCourseBiz struct {
	store UpdateCourseStorage
}

func NewUpdateCourseBiz(store UpdateCourseStorage) *updateCourseBiz {
	return &updateCourseBiz{store: store}
}

func (biz *updateCourseBiz) UpdateCourse(id string, course *coursemodel.Course) error {
	err := biz.store.Update(id, course)
	if err != nil {
		return err
	}
	return nil
}
