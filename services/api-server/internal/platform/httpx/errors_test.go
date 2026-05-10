package httpx

import "testing"

func TestErrorFactoriesUseStableCodes(t *testing.T) {
	cases := []*AppError{
		BadRequest("参数错误", nil),
		Unauthorized("未登录"),
		Forbidden("无权限", map[string]any{"permission": "portal:publish"}),
		NotFound("不存在", nil),
		Unavailable("服务暂不可用"),
		Internal("内部错误", nil),
	}

	seen := map[string]struct{}{}
	for _, appErr := range cases {
		if appErr.Code == "" {
			t.Fatal("expected non-empty error code")
		}
		if _, exists := seen[appErr.Code]; exists {
			t.Fatalf("duplicate error code detected: %s", appErr.Code)
		}
		seen[appErr.Code] = struct{}{}
		if appErr.HTTPStatus < 400 {
			t.Fatalf("expected error http status >= 400, got %d", appErr.HTTPStatus)
		}
	}
}
