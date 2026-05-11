package academic_provider

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type StudentSnapshot struct {
	Name    string
	Major   string
	College string
	Grade   string
}

type SemesterSnapshot struct {
	ID        string
	Name      string
	StartAt   time.Time
	EndAt     time.Time
	IsCurrent bool
}

type CourseScheduleSnapshot struct {
	ID           string
	SemesterID   string
	CourseName   string
	Teacher      string
	Location     string
	Weekday      int
	StartSection int
	EndSection   int
	Weeks        []int
	Credit       float64
}

type ExamSnapshot struct {
	ID         string
	SemesterID string
	CourseName string
	Location   string
	SeatNo     string
	Status     string
	StartsAt   time.Time
	EndsAt     time.Time
}

type GradeSnapshot struct {
	ID         string
	SemesterID string
	CourseName string
	CourseType string
	Credit     float64
	Score      float64
	GradePoint float64
	Status     string
}

type Provider interface {
	GenerateCaptcha(ctx context.Context, studentID string) (string, error)
	LoadStudentSnapshot(ctx context.Context, studentID, password string) (StudentSnapshot, error)
	ListSemesters(ctx context.Context, studentID string) ([]SemesterSnapshot, error)
	ListSchedule(ctx context.Context, studentID, semesterID string) ([]CourseScheduleSnapshot, error)
	ListExams(ctx context.Context, studentID, semesterID string) ([]ExamSnapshot, error)
	ListGrades(ctx context.Context, studentID, semesterID string) ([]GradeSnapshot, error)
}

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) GenerateCaptcha(_ context.Context, _ string) (string, error) {
	return "123456", nil
}

func (p *MockProvider) LoadStudentSnapshot(_ context.Context, studentID, password string) (StudentSnapshot, error) {
	if password == "" {
		return StudentSnapshot{}, fmt.Errorf("password is required")
	}

	suffix := studentID
	if len(suffix) > 4 {
		suffix = suffix[len(suffix)-4:]
	}

	colleges := []string{"海洋科学学院", "信息工程学院", "管理学院", "食品科学学院"}
	majors := []string{"软件工程", "信息管理", "海洋技术", "食品质量与安全"}
	grades := []string{"2022级", "2023级", "2024级"}
	index := len(studentID) % len(colleges)

	return StudentSnapshot{
		Name:    "同学" + suffix,
		Major:   majors[index%len(majors)],
		College: colleges[index],
		Grade:   grades[index%len(grades)],
	}, nil
}

func (p *MockProvider) ListSemesters(_ context.Context, studentID string) ([]SemesterSnapshot, error) {
	if strings.TrimSpace(studentID) == "" {
		return nil, fmt.Errorf("student id is required")
	}

	loc := time.FixedZone("CST", 8*60*60)
	return []SemesterSnapshot{
		{
			ID:        "2025-2026-2",
			Name:      "2025-2026 学年第二学期",
			StartAt:   time.Date(2026, time.February, 24, 0, 0, 0, 0, loc),
			EndAt:     time.Date(2026, time.July, 5, 23, 59, 59, 0, loc),
			IsCurrent: true,
		},
		{
			ID:        "2025-2026-1",
			Name:      "2025-2026 学年第一学期",
			StartAt:   time.Date(2025, time.September, 1, 0, 0, 0, 0, loc),
			EndAt:     time.Date(2026, time.January, 18, 23, 59, 59, 0, loc),
			IsCurrent: false,
		},
		{
			ID:        "2024-2025-2",
			Name:      "2024-2025 学年第二学期",
			StartAt:   time.Date(2025, time.February, 24, 0, 0, 0, 0, loc),
			EndAt:     time.Date(2025, time.July, 6, 23, 59, 59, 0, loc),
			IsCurrent: false,
		},
	}, nil
}

func (p *MockProvider) ListSchedule(_ context.Context, studentID, semesterID string) ([]CourseScheduleSnapshot, error) {
	if strings.TrimSpace(studentID) == "" || strings.TrimSpace(semesterID) == "" {
		return nil, fmt.Errorf("student id and semester id are required")
	}

	suffix := shortSuffix(studentID)
	switch semesterID {
	case "2025-2026-2":
		return []CourseScheduleSnapshot{
			{
				ID:           "course-math-" + suffix,
				SemesterID:   semesterID,
				CourseName:   "高等数学 B(2)",
				Teacher:      "王老师",
				Location:     "理工楼 A201",
				Weekday:      1,
				StartSection: 1,
				EndSection:   2,
				Weeks:        weekRange(1, 16),
				Credit:       4,
			},
			{
				ID:           "course-struct-" + suffix,
				SemesterID:   semesterID,
				CourseName:   "数据结构",
				Teacher:      "陈老师",
				Location:     "信息楼 302",
				Weekday:      3,
				StartSection: 3,
				EndSection:   4,
				Weeks:        weekRange(1, 16),
				Credit:       3.5,
			},
			{
				ID:           "course-english-" + suffix,
				SemesterID:   semesterID,
				CourseName:   "大学英语 IV",
				Teacher:      "Lily Zhang",
				Location:     "外语楼 105",
				Weekday:      4,
				StartSection: 1,
				EndSection:   2,
				Weeks:        weekRange(1, 14),
				Credit:       2,
			},
			{
				ID:           "course-policy-" + suffix,
				SemesterID:   semesterID,
				CourseName:   "形势与政策",
				Teacher:      "刘老师",
				Location:     "文科楼 407",
				Weekday:      5,
				StartSection: 7,
				EndSection:   8,
				Weeks:        weekRange(2, 10),
				Credit:       1,
			},
		}, nil
	case "2025-2026-1":
		return []CourseScheduleSnapshot{
			{
				ID:           "course-os-" + suffix,
				SemesterID:   semesterID,
				CourseName:   "操作系统",
				Teacher:      "李老师",
				Location:     "信息楼 405",
				Weekday:      2,
				StartSection: 1,
				EndSection:   2,
				Weeks:        weekRange(1, 16),
				Credit:       3.5,
			},
			{
				ID:           "course-db-" + suffix,
				SemesterID:   semesterID,
				CourseName:   "数据库原理",
				Teacher:      "孙老师",
				Location:     "理工楼 B301",
				Weekday:      4,
				StartSection: 3,
				EndSection:   4,
				Weeks:        weekRange(1, 16),
				Credit:       3,
			},
		}, nil
	case "2024-2025-2":
		return []CourseScheduleSnapshot{
			{
				ID:           "course-c-" + suffix,
				SemesterID:   semesterID,
				CourseName:   "程序设计基础",
				Teacher:      "周老师",
				Location:     "实验楼 201",
				Weekday:      1,
				StartSection: 3,
				EndSection:   4,
				Weeks:        weekRange(1, 16),
				Credit:       4,
			},
			{
				ID:           "course-linear-" + suffix,
				SemesterID:   semesterID,
				CourseName:   "线性代数",
				Teacher:      "张老师",
				Location:     "理工楼 C102",
				Weekday:      5,
				StartSection: 1,
				EndSection:   2,
				Weeks:        weekRange(1, 16),
				Credit:       3,
			},
		}, nil
	default:
		return nil, fmt.Errorf("semester %s not found", semesterID)
	}
}

func (p *MockProvider) ListExams(_ context.Context, studentID, semesterID string) ([]ExamSnapshot, error) {
	if strings.TrimSpace(studentID) == "" || strings.TrimSpace(semesterID) == "" {
		return nil, fmt.Errorf("student id and semester id are required")
	}

	loc := time.FixedZone("CST", 8*60*60)
	suffix := shortSuffix(studentID)
	switch semesterID {
	case "2025-2026-2":
		return []ExamSnapshot{
			{
				ID:         "exam-math-" + suffix,
				SemesterID: semesterID,
				CourseName: "高等数学 B(2)",
				Location:   "理工楼 A201",
				SeatNo:     "A-" + suffix,
				Status:     "scheduled",
				StartsAt:   time.Date(2026, time.June, 18, 9, 0, 0, 0, loc),
				EndsAt:     time.Date(2026, time.June, 18, 11, 0, 0, 0, loc),
			},
			{
				ID:         "exam-struct-" + suffix,
				SemesterID: semesterID,
				CourseName: "数据结构",
				Location:   "信息楼 302",
				SeatNo:     "C-" + suffix,
				Status:     "scheduled",
				StartsAt:   time.Date(2026, time.June, 20, 14, 0, 0, 0, loc),
				EndsAt:     time.Date(2026, time.June, 20, 16, 0, 0, 0, loc),
			},
		}, nil
	case "2025-2026-1":
		return []ExamSnapshot{
			{
				ID:         "exam-os-" + suffix,
				SemesterID: semesterID,
				CourseName: "操作系统",
				Location:   "信息楼 405",
				SeatNo:     "B-" + suffix,
				Status:     "completed",
				StartsAt:   time.Date(2026, time.January, 9, 9, 0, 0, 0, loc),
				EndsAt:     time.Date(2026, time.January, 9, 11, 0, 0, 0, loc),
			},
		}, nil
	case "2024-2025-2":
		return []ExamSnapshot{
			{
				ID:         "exam-c-" + suffix,
				SemesterID: semesterID,
				CourseName: "程序设计基础",
				Location:   "实验楼 201",
				SeatNo:     "D-" + suffix,
				Status:     "completed",
				StartsAt:   time.Date(2025, time.June, 28, 14, 0, 0, 0, loc),
				EndsAt:     time.Date(2025, time.June, 28, 16, 0, 0, 0, loc),
			},
		}, nil
	default:
		return nil, fmt.Errorf("semester %s not found", semesterID)
	}
}

func (p *MockProvider) ListGrades(_ context.Context, studentID, semesterID string) ([]GradeSnapshot, error) {
	if strings.TrimSpace(studentID) == "" || strings.TrimSpace(semesterID) == "" {
		return nil, fmt.Errorf("student id and semester id are required")
	}

	suffix := shortSuffix(studentID)
	switch semesterID {
	case "2025-2026-2":
		return []GradeSnapshot{
			{
				ID:         "grade-math-" + suffix,
				SemesterID: semesterID,
				CourseName: "高等数学 B(2)",
				CourseType: "required",
				Credit:     4,
				Score:      88,
				GradePoint: 3.8,
				Status:     "passed",
			},
			{
				ID:         "grade-struct-" + suffix,
				SemesterID: semesterID,
				CourseName: "数据结构",
				CourseType: "required",
				Credit:     3.5,
				Score:      92,
				GradePoint: 4.0,
				Status:     "passed",
			},
			{
				ID:         "grade-english-" + suffix,
				SemesterID: semesterID,
				CourseName: "大学英语 IV",
				CourseType: "required",
				Credit:     2,
				Score:      85,
				GradePoint: 3.7,
				Status:     "passed",
			},
		}, nil
	case "2025-2026-1":
		return []GradeSnapshot{
			{
				ID:         "grade-os-" + suffix,
				SemesterID: semesterID,
				CourseName: "操作系统",
				CourseType: "required",
				Credit:     3.5,
				Score:      90,
				GradePoint: 4.0,
				Status:     "passed",
			},
			{
				ID:         "grade-db-" + suffix,
				SemesterID: semesterID,
				CourseName: "数据库原理",
				CourseType: "required",
				Credit:     3,
				Score:      87,
				GradePoint: 3.8,
				Status:     "passed",
			},
		}, nil
	case "2024-2025-2":
		return []GradeSnapshot{
			{
				ID:         "grade-c-" + suffix,
				SemesterID: semesterID,
				CourseName: "程序设计基础",
				CourseType: "required",
				Credit:     4,
				Score:      91,
				GradePoint: 4.0,
				Status:     "passed",
			},
			{
				ID:         "grade-linear-" + suffix,
				SemesterID: semesterID,
				CourseName: "线性代数",
				CourseType: "required",
				Credit:     3,
				Score:      84,
				GradePoint: 3.5,
				Status:     "passed",
			},
		}, nil
	default:
		return nil, fmt.Errorf("semester %s not found", semesterID)
	}
}

func weekRange(start, end int) []int {
	if end < start {
		return nil
	}

	result := make([]int, 0, end-start+1)
	for value := start; value <= end; value++ {
		result = append(result, value)
	}
	return result
}

func shortSuffix(studentID string) string {
	trimmed := strings.TrimSpace(studentID)
	if len(trimmed) <= 2 {
		return trimmed
	}
	return trimmed[len(trimmed)-2:]
}
