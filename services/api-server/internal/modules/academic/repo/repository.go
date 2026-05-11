package repo

import (
	"context"

	academictypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/academic/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/academic_provider"
)

type Repository interface {
	ListSemesters(ctx context.Context, studentID string) ([]academictypes.Semester, error)
	ListSchedule(ctx context.Context, studentID, semesterID string) ([]academictypes.ScheduleCourse, error)
	ListExams(ctx context.Context, studentID, semesterID string) ([]academictypes.ExamItem, error)
	ListGrades(ctx context.Context, studentID, semesterID string) ([]academictypes.GradeItem, error)
}

type ProviderRepository struct {
	provider academic_provider.Provider
}

func NewProviderRepository(provider academic_provider.Provider) *ProviderRepository {
	if provider == nil {
		provider = academic_provider.NewMockProvider()
	}

	return &ProviderRepository{provider: provider}
}

func (r *ProviderRepository) ListSemesters(ctx context.Context, studentID string) ([]academictypes.Semester, error) {
	items, err := r.provider.ListSemesters(ctx, studentID)
	if err != nil {
		return nil, err
	}

	result := make([]academictypes.Semester, 0, len(items))
	for _, item := range items {
		result = append(result, academictypes.Semester{
			ID:        item.ID,
			Name:      item.Name,
			StartAt:   item.StartAt,
			EndAt:     item.EndAt,
			IsCurrent: item.IsCurrent,
		})
	}
	return result, nil
}

func (r *ProviderRepository) ListSchedule(
	ctx context.Context,
	studentID, semesterID string,
) ([]academictypes.ScheduleCourse, error) {
	items, err := r.provider.ListSchedule(ctx, studentID, semesterID)
	if err != nil {
		return nil, err
	}

	result := make([]academictypes.ScheduleCourse, 0, len(items))
	for _, item := range items {
		result = append(result, academictypes.ScheduleCourse{
			ID:           item.ID,
			SemesterID:   item.SemesterID,
			CourseName:   item.CourseName,
			Teacher:      item.Teacher,
			Location:     item.Location,
			Weekday:      item.Weekday,
			StartSection: item.StartSection,
			EndSection:   item.EndSection,
			Weeks:        append([]int(nil), item.Weeks...),
			Credit:       item.Credit,
		})
	}
	return result, nil
}

func (r *ProviderRepository) ListExams(
	ctx context.Context,
	studentID, semesterID string,
) ([]academictypes.ExamItem, error) {
	items, err := r.provider.ListExams(ctx, studentID, semesterID)
	if err != nil {
		return nil, err
	}

	result := make([]academictypes.ExamItem, 0, len(items))
	for _, item := range items {
		result = append(result, academictypes.ExamItem{
			ID:         item.ID,
			SemesterID: item.SemesterID,
			CourseName: item.CourseName,
			Location:   item.Location,
			SeatNo:     item.SeatNo,
			Status:     item.Status,
			StartsAt:   item.StartsAt,
			EndsAt:     item.EndsAt,
		})
	}
	return result, nil
}

func (r *ProviderRepository) ListGrades(
	ctx context.Context,
	studentID, semesterID string,
) ([]academictypes.GradeItem, error) {
	items, err := r.provider.ListGrades(ctx, studentID, semesterID)
	if err != nil {
		return nil, err
	}

	result := make([]academictypes.GradeItem, 0, len(items))
	for _, item := range items {
		result = append(result, academictypes.GradeItem{
			ID:         item.ID,
			SemesterID: item.SemesterID,
			CourseName: item.CourseName,
			CourseType: item.CourseType,
			Credit:     item.Credit,
			Score:      item.Score,
			GradePoint: item.GradePoint,
			Status:     item.Status,
		})
	}
	return result, nil
}
