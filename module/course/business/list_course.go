package coursebiz

import (
	coursemodel "crawdata/module/course/model"
	"log"
	"time"
)

type ListCourseStorage interface {
	GetAllCourses(name string) ([]coursemodel.Course, error)
}

type listCourseBiz struct {
	store ListCourseStorage
}

func NewListCourseBiz(store ListCourseStorage) *listCourseBiz {
	return &listCourseBiz{store: store}
}

func (biz *listCourseBiz) ListCourse(name, startDate, endDate string) ([]coursemodel.Course, error) {
	// Gọi repo để lấy tất cả courses theo name
	courses, err := biz.store.GetAllCourses(name)
	if err != nil {
		return nil, err
	}

	// Xử lý filter theo thời gian
	layoutVN := "2 thg 1, 2006"
	layoutISO := "2006-01-02"
	var filtered []coursemodel.Course
	var startTime, endTime time.Time

	if startDate != "" {
		startTime, err = time.Parse(layoutISO, startDate)
		if err != nil {
			log.Println("Lỗi parse startDate:", err)
			startDate = ""
		}
	}

	if endDate != "" {
		endTime, err = time.Parse(layoutISO, endDate)
		if err != nil {
			log.Println("Lỗi parse endDate:", err)
			endDate = ""
		}
	}

	for _, course := range courses {
		updateTime, err := time.Parse(layoutVN, course.Update)
		if err != nil {
			log.Printf("Lỗi parse update [%s] cho ID %d: %v\n", course.Update, course.Id, err)
			continue
		}

		if (startDate == "" || !updateTime.Before(startTime)) &&
			(endDate == "" || !updateTime.After(endTime)) {
			filtered = append(filtered, course)
		}
	}

	return filtered, nil
}
