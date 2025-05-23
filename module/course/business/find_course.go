package coursebiz

import (
	coursemodel "crawdata/module/course/model"
)

type FindCourseStorage interface {
	GetCourseById(id string) (coursemodel.Course, error)
}

type findCourseBiz struct {
	store FindCourseStorage
}

func NewFindCourseBiz(store FindCourseStorage) *findCourseBiz {
	return &findCourseBiz{store: store}
}

func (biz *findCourseBiz) FindCourse(id string) (coursemodel.Course, error) {
	course, err := biz.store.GetCourseById(id)
	if err != nil {
		return coursemodel.Course{}, err
	}
	return course, nil
}
