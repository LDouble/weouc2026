package openapi_test

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCanonicalOpenAPIContainsCoreRoutes(t *testing.T) {
	content, err := os.ReadFile("openapi.yaml")
	if err != nil {
		t.Fatalf("read openapi.yaml failed: %v", err)
	}

	var document struct {
		OpenAPI string                 `yaml:"openapi"`
		Paths   map[string]interface{} `yaml:"paths"`
	}
	if err := yaml.Unmarshal(content, &document); err != nil {
		t.Fatalf("unmarshal openapi document failed: %v", err)
	}

	if document.OpenAPI != "3.1.0" {
		t.Fatalf("expected openapi version 3.1.0, got %q", document.OpenAPI)
	}

	expectedPaths := []string{
		"/healthz",
		"/readyz",
		"/api/v1/system/profile",
		"/api/auth/wechat/login",
		"/api/student",
		"/api/edu/send-captcha",
		"/api/feed/list",
		"/api/market/detail/{id}",
		"/api/errand/accept",
		"/api/resource/list",
		"/api/lostFound/detail/{id}",
		"/api/carpool/list",
		"/api/carpool/publish",
		"/api/meetup/list",
		"/api/meetup/publish",
		"/api/admin/campus-life/review/list",
		"/api/admin/campus-life/review/update",
	}
	for _, path := range expectedPaths {
		if _, exists := document.Paths[path]; !exists {
			t.Fatalf("expected path %s in openapi document", path)
		}
	}
}
