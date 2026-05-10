package academic_provider

import (
	"context"
	"fmt"
)

type StudentSnapshot struct {
	Name    string
	Major   string
	College string
	Grade   string
}

type Provider interface {
	GenerateCaptcha(ctx context.Context, studentID string) (string, error)
	LoadStudentSnapshot(ctx context.Context, studentID, password string) (StudentSnapshot, error)
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
