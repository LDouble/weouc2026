package types

import "time"

type Semester struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	StartAt   time.Time `json:"start_at"`
	EndAt     time.Time `json:"end_at"`
	IsCurrent bool      `json:"is_current"`
}

type ScheduleCourse struct {
	ID           string  `json:"id"`
	SemesterID   string  `json:"semester_id"`
	CourseName   string  `json:"course_name"`
	Teacher      string  `json:"teacher"`
	Location     string  `json:"location"`
	Weekday      int     `json:"weekday"`
	StartSection int     `json:"start_section"`
	EndSection   int     `json:"end_section"`
	Weeks        []int   `json:"weeks"`
	Credit       float64 `json:"credit"`
}

type ExamItem struct {
	ID         string    `json:"id"`
	SemesterID string    `json:"semester_id"`
	CourseName string    `json:"course_name"`
	Location   string    `json:"location"`
	SeatNo     string    `json:"seat_no"`
	Status     string    `json:"status"`
	StartsAt   time.Time `json:"starts_at"`
	EndsAt     time.Time `json:"ends_at"`
}

type GradeItem struct {
	ID         string  `json:"id"`
	SemesterID string  `json:"semester_id"`
	CourseName string  `json:"course_name"`
	CourseType string  `json:"course_type"`
	Credit     float64 `json:"credit"`
	Score      float64 `json:"score"`
	GradePoint float64 `json:"grade_point"`
	Status     string  `json:"status"`
}
