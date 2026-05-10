package types

import "time"

type HealthStatus struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
}

type DependencyStatus struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Required bool   `json:"required"`
	Detail   string `json:"detail,omitempty"`
}

type ReadinessStatus struct {
	Status       string             `json:"status"`
	Dependencies []DependencyStatus `json:"dependencies"`
	Timestamp    time.Time          `json:"timestamp"`
}

func (r ReadinessStatus) IsReady() bool {
	return r.Status == "ready"
}

type Profile struct {
	Service ProfileService `json:"service"`
	Auth    ProfileAuth    `json:"auth"`
}

type ProfileService struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type ProfileAuth struct {
	Authenticated bool     `json:"authenticated"`
	UserID        string   `json:"user_id,omitempty"`
	Roles         []string `json:"roles,omitempty"`
	Permissions   []string `json:"permissions,omitempty"`
	AcademicBound bool     `json:"academic_bound"`
}
