package service

import (
	"context"
	"errors"
	"path"
	"strings"
	"time"

	clrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/repo"
	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

func (s *Service) ListResources(ctx context.Context, principal auth.Principal, query cltypes.ResourceQuery) (map[string]any, error) {
	currentUserID, visibleStatuses, includeAllStatus := buildVisibilityFilter(principal)

	items, total, err := s.repository.ListByType(ctx, cltypes.ContentTypeResource, cltypes.ContentFilter{
		Pagination:       query.Pagination,
		Category:         query.Category,
		Keyword:          query.Keyword,
		CurrentUserID:    currentUserID,
		VisibleStatuses:  visibleStatuses,
		IncludeAllStatus: includeAllStatus,
	})
	if err != nil {
		return nil, httpx.Internal("读取资料列表失败", err)
	}

	list := make([]map[string]any, 0, len(items))
	for _, item := range items {
		rp, _ := unmarshalPayload[cltypes.ResourcePayload](item.TypePayload)
		resolvedFiles := resolveResourceFiles(ctx, s.storageProvider, rp.Files)
		list = append(list, map[string]any{
			"id":                item.ID.Hex(),
			"title":             item.Title,
			"desc":              item.Desc,
			"publisher":         "校园用户",
			"publisher_initial": "校",
			"created_at":        item.CreatedAt.Format(time.RFC3339),
			"status":            item.Status,
			"extra": map[string]any{
				"category":     rp.Category,
				"course_name":  rp.CourseName,
				"contact":      item.Contact,
				"files":        resolvedFiles,
				"file_size":    rp.FileSize,
				"file_type":    rp.FileType,
				"download_url": firstResourceURL(resolvedFiles, resolveManagedURL(ctx, s.storageProvider, rp.DownloadURL)),
				"likes":        rp.Likes,
				"views":        rp.Views,
			},
		})
	}
	return listEnvelope(list, int(total), query.Pagination), nil
}

func (s *Service) GetResourceDetail(ctx context.Context, principal auth.Principal, id string) (map[string]any, error) {
	item, err := s.repository.GetByID(ctx, id)
	if errors.Is(err, clrepo.ErrNotFound) {
		return nil, httpx.NotFound("资料不存在", nil)
	}
	if err != nil {
		return nil, httpx.Internal("读取资料详情失败", err)
	}
	if item.ContentType != cltypes.ContentTypeResource {
		return nil, httpx.NotFound("资料不存在", nil)
	}
	if err := ensureContentVisible(principal, item, "资料不存在"); err != nil {
		return nil, err
	}
	canView := canViewContact(principal, item.PublisherUserID)
	role := simpleUserRole(item.PublisherUserID, principal)
	isOwner := role == "publisher"
	rp, _ := unmarshalPayload[cltypes.ResourcePayload](item.TypePayload)
	resolvedFiles := resolveResourceFiles(ctx, s.storageProvider, rp.Files)
	return map[string]any{
		"id":                item.ID.Hex(),
		"title":             item.Title,
		"desc":              item.Desc,
		"publisher":         publisherName(item, principal),
		"publisher_initial": initialOf(publisherName(item, principal)),
		"created_at":        item.CreatedAt.Format(time.RFC3339),
		"status":            item.Status,
		"user_role":         role,
		"is_owner":          isOwner,
		"can_view_contact":  canView,
		"can_edit":          canEditContent(isOwner, item.Status),
		"can_delete":        canDeleteContent(isOwner, item.Status),
		"can_download":      principal.Authenticated && canView,
		"extra": map[string]any{
			"category":     rp.Category,
			"course_name":  rp.CourseName,
			"contact":      visibleValue(canView, item.Contact),
			"files":        resolvedFiles,
			"file_size":    rp.FileSize,
			"file_type":    rp.FileType,
			"download_url": firstResourceURL(resolvedFiles, resolveManagedURL(ctx, s.storageProvider, rp.DownloadURL)),
			"likes":        rp.Likes,
			"views":        rp.Views,
		},
	}, nil
}

func (s *Service) PublishResource(ctx context.Context, principal auth.Principal, request cltypes.ResourcePublishRequest) (map[string]any, error) {
	if strings.TrimSpace(request.Title) == "" || len(request.FilePaths) == 0 {
		return nil, httpx.BadRequest("资料标题和文件路径不能为空", nil)
	}
	files := make([]cltypes.ResourceFile, 0, len(request.FilePaths))
	for _, filePath := range request.FilePaths {
		filePath = strings.TrimSpace(filePath)
		if filePath == "" {
			continue
		}
		files = append(files, cltypes.ResourceFile{
			Name:     path.Base(filePath),
			Path:     strings.TrimPrefix(filePath, "/"),
			FileType: detectFileType(filePath),
			FileSize: "1.0MB",
		})
	}
	item := cltypes.CommunityContent{
		ContentType:     cltypes.ContentTypeResource,
		Title:           strings.TrimSpace(request.Title),
		Desc:            firstNonEmpty(strings.TrimSpace(request.Desc), strings.TrimSpace(request.Title)),
		Status:          cltypes.StatusReviewing,
		PublisherUserID: principal.UserID,
		Contact:         strings.TrimSpace(request.Contact),
		TypePayload: marshalPayload(cltypes.ResourcePayload{
			Category:   strings.TrimSpace(request.Category),
			CourseName: strings.TrimSpace(request.CourseName),
			Files:      files,
			FileSize:   firstFileSize(files),
			FileType:   firstFileType(files),
			Views:      0,
			Likes:      0,
		}),
		CreatedBy: principal.UserID,
		UpdatedBy: principal.UserID,
	}
	item, err := s.repository.Save(ctx, item)
	if err != nil {
		return nil, httpx.Internal("保存资料失败", err)
	}
	s.recordAudit(ctx, principal, "campus_life.resource.publish", "resource", item.ID.Hex(), "资料发布成功", map[string]any{
		"status":      item.Status,
		"category":    request.Category,
		"course_name": request.CourseName,
	})
	return map[string]any{"id": item.ID.Hex()}, nil
}

func (s *Service) DeleteResource(ctx context.Context, principal auth.Principal, id string) error {
	_, err := s.repository.Update(ctx, id, func(item *cltypes.CommunityContent) error {
		if item.ContentType != cltypes.ContentTypeResource {
			return httpx.NotFound("资料不存在", nil)
		}
		if item.PublisherUserID != principal.UserID {
			return httpx.Forbidden("只有发布者可以下架", nil)
		}
		if item.Status != cltypes.StatusPublished && item.Status != cltypes.StatusReviewing && item.Status != cltypes.StatusRejected {
			return httpx.BadRequest("当前状态不允许下架", nil)
		}
		item.Status = cltypes.StatusOffline
		item.UpdatedBy = principal.UserID
		return nil
	})
	if err != nil {
		if errors.Is(err, clrepo.ErrNotFound) {
			return httpx.NotFound("资料不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("下架资料失败", err)
	}
	s.recordAudit(ctx, principal, "campus_life.resource.delete", "resource", id, "资料下架成功", nil)
	return nil
}
