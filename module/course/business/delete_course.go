package coursebiz

type DeleteCourseStorage interface {
	Delete(id string) error
}

type deleteCourseBiz struct {
	store DeleteCourseStorage
}

func NewDeleteCourseBiz(store DeleteCourseStorage) *deleteCourseBiz {
	return &deleteCourseBiz{store: store}
}

func (biz *deleteCourseBiz) DeleteCourse(id string) error {
	err := biz.store.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
