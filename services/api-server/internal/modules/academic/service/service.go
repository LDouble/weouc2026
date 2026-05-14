package service

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	academicrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/academic/repo"
	academictypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/academic/types"
	iamrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/repo"
	iamtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/audit"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

type Service struct {
	repository academicrepo.Repository
	users      iamrepo.UserRepository
	recorder   audit.Recorder
}

func New(repository academicrepo.Repository, users iamrepo.UserRepository, recorder audit.Recorder) *Service {
	if repository == nil {
		repository = academicrepo.NewProviderRepository(nil)
	}

	return &Service{
		repository: repository,
		users:      users,
		recorder:   recorder,
	}
}

func (s *Service) ListSemesters(ctx context.Context, principal auth.Principal) (map[string]any, error) {
	profile, err := s.loadBoundProfile(ctx, principal)
	if err != nil {
		return nil, err
	}

	semesters, err := s.repository.ListSemesters(ctx, profile.StudentID)
	if err != nil {
		return nil, httpx.Internal("读取学期列表失败", err)
	}

	current, err := resolveSemester(semesters, "")
	if err != nil {
		return nil, err
	}

	audit.RecordBestEffort(ctx, s.recorder, audit.Entry{
		ActorID:      principal.UserID,
		ActorName:    principal.DisplayName,
		Action:       "academic.semesters.view",
		ResourceType: "academic_semester",
		ResourceID:   current.ID,
		Message:      "查看学期列表",
		Details: map[string]any{
			"student_id":          profile.StudentID,
			"current_semester_id": current.ID,
			"semester_count":      len(semesters),
		},
	})

	return map[string]any{
		"current_semester_id": current.ID,
		"list":                semesters,
	}, nil
}

func (s *Service) GetSchedule(
	ctx context.Context,
	principal auth.Principal,
	semesterID string,
) (map[string]any, error) {
	profile, err := s.loadBoundProfile(ctx, principal)
	if err != nil {
		return nil, err
	}

	semester, err := s.resolveSemester(ctx, profile.StudentID, semesterID)
	if err != nil {
		return nil, err
	}

	items, err := s.repository.ListSchedule(ctx, profile.StudentID, semester.ID)
	if err != nil {
		return nil, httpx.Internal("读取课程表失败", err)
	}

	teachingDays := make(map[int]struct{})
	for _, item := range items {
		teachingDays[item.Weekday] = struct{}{}
	}

	audit.RecordBestEffort(ctx, s.recorder, audit.Entry{
		ActorID:      principal.UserID,
		ActorName:    principal.DisplayName,
		Action:       "academic.schedule.view",
		ResourceType: "academic_schedule",
		ResourceID:   semester.ID,
		Message:      "查看课程表",
		Details: map[string]any{
			"student_id":   profile.StudentID,
			"semester_id":  semester.ID,
			"course_count": len(items),
		},
	})

	return map[string]any{
		"semester": semester,
		"list":     items,
		"summary": map[string]any{
			"course_count":  len(items),
			"teaching_days": len(teachingDays),
		},
	}, nil
}

func (s *Service) ListExams(
	ctx context.Context,
	principal auth.Principal,
	semesterID string,
) (map[string]any, error) {
	profile, err := s.loadBoundProfile(ctx, principal)
	if err != nil {
		return nil, err
	}

	semester, err := s.resolveSemester(ctx, profile.StudentID, semesterID)
	if err != nil {
		return nil, err
	}

	items, err := s.repository.ListExams(ctx, profile.StudentID, semester.ID)
	if err != nil {
		return nil, httpx.Internal("读取考试安排失败", err)
	}

	upcomingCount := 0
	now := time.Now().UTC()
	for _, item := range items {
		if item.StartsAt.UTC().After(now) {
			upcomingCount++
		}
	}

	audit.RecordBestEffort(ctx, s.recorder, audit.Entry{
		ActorID:      principal.UserID,
		ActorName:    principal.DisplayName,
		Action:       "academic.exams.view",
		ResourceType: "academic_exam",
		ResourceID:   semester.ID,
		Message:      "查看考试安排",
		Details: map[string]any{
			"student_id":     profile.StudentID,
			"semester_id":    semester.ID,
			"exam_count":     len(items),
			"upcoming_count": upcomingCount,
		},
	})

	return map[string]any{
		"semester": semester,
		"list":     items,
		"summary": map[string]any{
			"count":          len(items),
			"upcoming_count": upcomingCount,
		},
	}, nil
}

func (s *Service) ListGrades(
	ctx context.Context,
	principal auth.Principal,
	semesterID string,
) (map[string]any, error) {
	profile, err := s.loadBoundProfile(ctx, principal)
	if err != nil {
		return nil, err
	}

	semester, err := s.resolveSemester(ctx, profile.StudentID, semesterID)
	if err != nil {
		return nil, err
	}

	items, err := s.repository.ListGrades(ctx, profile.StudentID, semester.ID)
	if err != nil {
		return nil, httpx.Internal("读取成绩单失败", err)
	}

	var totalScore float64
	var totalGradePoint float64
	passedCount := 0
	for _, item := range items {
		totalScore += item.Score
		totalGradePoint += item.GradePoint
		if strings.EqualFold(item.Status, "passed") {
			passedCount++
		}
	}

	averageScore := 0.0
	averageGradePoint := 0.0
	if len(items) > 0 {
		averageScore = roundTo(totalScore/float64(len(items)), 2)
		averageGradePoint = roundTo(totalGradePoint/float64(len(items)), 2)
	}

	audit.RecordBestEffort(ctx, s.recorder, audit.Entry{
		ActorID:      principal.UserID,
		ActorName:    principal.DisplayName,
		Action:       "academic.grades.view",
		ResourceType: "academic_grade",
		ResourceID:   semester.ID,
		Message:      "查看成绩单",
		Details: map[string]any{
			"student_id":          profile.StudentID,
			"semester_id":         semester.ID,
			"course_count":        len(items),
			"average_score":       averageScore,
			"average_grade_point": averageGradePoint,
		},
	})

	return map[string]any{
		"semester": semester,
		"list":     items,
		"summary": map[string]any{
			"course_count":        len(items),
			"passed_count":        passedCount,
			"average_score":       averageScore,
			"average_grade_point": averageGradePoint,
		},
	}, nil
}

func (s *Service) resolveSemester(ctx context.Context, studentID, semesterID string) (academictypes.Semester, error) {
	semesters, err := s.repository.ListSemesters(ctx, studentID)
	if err != nil {
		return academictypes.Semester{}, httpx.Internal("读取学期列表失败", err)
	}
	return resolveSemester(semesters, semesterID)
}

func (s *Service) loadBoundProfile(ctx context.Context, principal auth.Principal) (iamtypes.StudentProfile, error) {
	user, err := s.users.FindByID(ctx, principal.UserID)
	if errors.Is(err, iamrepo.ErrUserNotFound) {
		return iamtypes.StudentProfile{}, httpx.Unauthorized("需要登录后访问")
	}
	if err != nil {
		return iamtypes.StudentProfile{}, httpx.Internal("读取当前用户资料失败", err)
	}
	if user.StudentProfile == nil || !user.StudentProfile.IsBound || strings.TrimSpace(user.StudentProfile.StudentID) == "" {
		return iamtypes.StudentProfile{}, httpx.Forbidden("当前账号尚未完成教务绑定", map[string]any{
			"required_state": "academic_bound",
		})
	}

	return *user.StudentProfile, nil
}

func resolveSemester(semesters []academictypes.Semester, requestedID string) (academictypes.Semester, error) {
	if len(semesters) == 0 {
		return academictypes.Semester{}, httpx.NotFound("当前账号暂无可用学期数据", nil)
	}

	requestedID = strings.TrimSpace(requestedID)
	if requestedID != "" {
		for _, semester := range semesters {
			if semester.ID == requestedID {
				return semester, nil
			}
		}
		return academictypes.Semester{}, httpx.BadRequest("semester_id 不存在", map[string]any{
			"semester_id": requestedID,
		})
	}

	for _, semester := range semesters {
		if semester.IsCurrent {
			return semester, nil
		}
	}
	return semesters[0], nil
}

func roundTo(value float64, places int) float64 {
	pow := math.Pow10(places)
	return math.Round(value*pow) / pow
}
