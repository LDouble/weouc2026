package service

import (
	"context"
	"strings"

	fcconfig "github.com/liangluo/weouc2026/services/api-server/internal/modules/file_center/config"
	fctypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/file_center/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/storage_provider"
)

type Service struct {
	config          fcconfig.ModuleConfig
	storageProvider storage_provider.Provider
}

func New(config fcconfig.ModuleConfig, storageProvider storage_provider.Provider) *Service {
	return &Service{
		config:          config,
		storageProvider: storageProvider,
	}
}

func (s *Service) IssueUploadCredentials(
	ctx context.Context,
	principal auth.Principal,
	scene string,
) (fctypes.UploadSTSCredentials, error) {
	if !principal.Authenticated || principal.UserID == "" {
		return fctypes.UploadSTSCredentials{}, httpx.Unauthorized("")
	}
	if s.storageProvider == nil {
		return fctypes.UploadSTSCredentials{}, httpx.Unavailable("对象存储未启用")
	}

	credentials, err := s.storageProvider.IssueUploadCredentials(ctx, storage_provider.IssueUploadCredentialsInput{
		UserID: principal.UserID,
		Scene:  s.config.NormalizeScene(scene),
	})
	if err != nil {
		return fctypes.UploadSTSCredentials{}, httpx.Internal("签发对象存储临时凭证失败", err)
	}

	return fctypes.UploadSTSCredentials{
		TmpSecretID:  credentials.TmpSecretID,
		TmpSecretKey: credentials.TmpSecretKey,
		SessionToken: credentials.SessionToken,
		StartTime:    credentials.StartTime,
		ExpiredTime:  credentials.ExpiredTime,
		Bucket:       credentials.Bucket,
		Region:       credentials.Region,
		PathPrefix:   credentials.PathPrefix,
	}, nil
}

func (s *Service) PresignGet(ctx context.Context, request fctypes.PresignedGetRequest) (fctypes.PresignedGetResponse, error) {
	objectPath := strings.TrimSpace(request.Path)
	if objectPath == "" {
		return fctypes.PresignedGetResponse{}, httpx.BadRequest("path 不能为空", nil)
	}
	if s.storageProvider == nil {
		return fctypes.PresignedGetResponse{}, httpx.Unavailable("对象存储未启用")
	}

	url, err := s.storageProvider.PresignGetURL(ctx, objectPath)
	if err != nil {
		return fctypes.PresignedGetResponse{}, httpx.Internal("签发对象下载地址失败", err)
	}

	return fctypes.PresignedGetResponse{
		Path: objectPath,
		URL:  url,
	}, nil
}
