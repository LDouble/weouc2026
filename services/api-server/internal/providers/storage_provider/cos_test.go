package storage_provider

import (
	"testing"
	"time"
)

func TestBuildObjectPrefixIncludesSceneUserAndDate(t *testing.T) {
	prefix := buildObjectPrefix("miniapp", "market", "wx-user_1", time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC))
	if prefix != "miniapp/market/wx-user_1/20260510/" {
		t.Fatalf("unexpected prefix %q", prefix)
	}
}

func TestSanitizePathSegmentFallsBackWhenEmpty(t *testing.T) {
	if got := sanitizePathSegment("   ", "general"); got != "general" {
		t.Fatalf("expected fallback, got %q", got)
	}
}

func TestNormalizePolicyPrefixAlwaysEndsWithSlash(t *testing.T) {
	if got := normalizePolicyPrefix("/miniapp/market/u1/"); got != "miniapp/market/u1/" {
		t.Fatalf("unexpected normalized prefix %q", got)
	}
}
