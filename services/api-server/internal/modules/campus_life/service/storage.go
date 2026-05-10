package service

import (
	"context"
	"strings"

	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/storage_provider"
)

func resolveManagedURL(ctx context.Context, provider storage_provider.Provider, value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return value
	}
	if provider == nil {
		return value
	}

	url, err := provider.PresignGetURL(ctx, value)
	if err != nil {
		return ""
	}

	return url
}

func resolveManagedURLs(ctx context.Context, provider storage_provider.Provider, values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if resolved := resolveManagedURL(ctx, provider, value); resolved != "" {
			result = append(result, resolved)
		}
	}
	return result
}

func resolveResourceFiles(ctx context.Context, provider storage_provider.Provider, files []cltypes.ResourceFile) []cltypes.ResourceFile {
	if len(files) == 0 {
		return nil
	}

	result := make([]cltypes.ResourceFile, 0, len(files))
	for _, file := range files {
		next := file
		if next.Path == "" && next.URL != "" && !strings.HasPrefix(next.URL, "http://") && !strings.HasPrefix(next.URL, "https://") {
			next.Path = next.URL
			next.URL = ""
		}
		if next.URL == "" {
			next.URL = resolveManagedURL(ctx, provider, next.Path)
		} else {
			next.URL = resolveManagedURL(ctx, provider, next.URL)
		}
		result = append(result, next)
	}

	return result
}

func firstResourceURL(files []cltypes.ResourceFile, fallback string) string {
	for _, file := range files {
		if strings.TrimSpace(file.URL) != "" {
			return file.URL
		}
	}

	return fallback
}
