package storage_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcprofile "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"
	cos "github.com/tencentyun/cos-go-sdk-v5"

	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
)

var unsafePathSegmentPattern = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

type COSProvider struct {
	config    appconfig.COSConfig
	stsClient *sts.Client
	cosClient *cos.Client
}

func NewCOSProvider(cfg appconfig.COSConfig) (*COSProvider, error) {
	bucketURL, err := cos.NewBucketURL(cfg.Bucket, cfg.Region, true)
	if err != nil {
		return nil, fmt.Errorf("build bucket url failed: %w", err)
	}

	credential := common.NewCredential(cfg.SecretID, cfg.SecretKey)
	clientProfile := tcprofile.NewClientProfile()
	stsClient, err := sts.NewClient(credential, cfg.Region, clientProfile)
	if err != nil {
		return nil, fmt.Errorf("create sts client failed: %w", err)
	}

	cosClient := cos.NewClient(&cos.BaseURL{BucketURL: bucketURL}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})

	return &COSProvider{
		config:    cfg,
		stsClient: stsClient,
		cosClient: cosClient,
	}, nil
}

func (p *COSProvider) IssueUploadCredentials(ctx context.Context, input IssueUploadCredentialsInput) (UploadCredentials, error) {
	scene := sanitizePathSegment(input.Scene, "general")
	userID := sanitizePathSegment(input.UserID, "anonymous")
	pathPrefix := buildObjectPrefix(p.config.PathPrefix, scene, userID, time.Now().UTC())

	policyJSON, err := p.buildPolicyJSON(pathPrefix)
	if err != nil {
		return UploadCredentials{}, err
	}

	request := sts.NewGetFederationTokenRequest()
	request.Name = stringPtr("weouccos")
	request.Policy = stringPtr(policyJSON)
	request.DurationSeconds = uint64Ptr(uint64(p.config.STSDuration.Seconds()))

	response, err := p.stsClient.GetFederationTokenWithContext(ctx, request)
	if err != nil {
		return UploadCredentials{}, fmt.Errorf("issue cos federation token failed: %w", err)
	}

	if response == nil || response.Response == nil || response.Response.Credentials == nil {
		return UploadCredentials{}, fmt.Errorf("cos federation token response is empty")
	}

	credentials := response.Response.Credentials
	expiredAt := valueUint64(response.Response.ExpiredTime)
	return UploadCredentials{
		TmpSecretID:  valueString(credentials.TmpSecretId),
		TmpSecretKey: valueString(credentials.TmpSecretKey),
		SessionToken: valueString(credentials.Token),
		StartTime:    time.Now().UTC().Unix(),
		ExpiredTime:  int64(expiredAt),
		Bucket:       p.config.Bucket,
		Region:       p.config.Region,
		PathPrefix:   pathPrefix,
	}, nil
}

func (p *COSProvider) PresignGetURL(ctx context.Context, objectPath string) (string, error) {
	objectPath = normalizeObjectPath(objectPath)
	if objectPath == "" {
		return "", fmt.Errorf("object path is empty")
	}

	presignedURL, err := p.cosClient.Object.GetPresignedURL(
		ctx,
		http.MethodGet,
		objectPath,
		p.config.SecretID,
		p.config.SecretKey,
		p.config.PresignedGETTTL,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("sign cos object url failed: %w", err)
	}

	return presignedURL.String(), nil
}

func (p *COSProvider) Check(ctx context.Context) error {
	_, err := p.cosClient.Bucket.Head(ctx)
	if err != nil {
		return fmt.Errorf("head cos bucket failed: %w", err)
	}

	return nil
}

func (p *COSProvider) buildPolicyJSON(pathPrefix string) (string, error) {
	resource := fmt.Sprintf(
		"qcs::cos:%s:uid/%s:%s/%s*",
		p.config.Region,
		p.config.BucketAppID(),
		p.config.Bucket,
		normalizePolicyPrefix(pathPrefix),
	)

	policy := map[string]any{
		"version": "2.0",
		"statement": []map[string]any{
			{
				"action": []string{
					"name/cos:PutObject",
					"name/cos:PostObject",
					"name/cos:InitiateMultipartUpload",
					"name/cos:UploadPart",
					"name/cos:ListParts",
					"name/cos:CompleteMultipartUpload",
					"name/cos:AbortMultipartUpload",
				},
				"effect":   "allow",
				"resource": []string{resource},
			},
		},
	}

	raw, err := json.Marshal(policy)
	if err != nil {
		return "", fmt.Errorf("marshal cos policy failed: %w", err)
	}

	return string(raw), nil
}

func buildObjectPrefix(basePrefix, scene, userID string, now time.Time) string {
	parts := make([]string, 0, 4)
	if trimmed := normalizeObjectPath(basePrefix); trimmed != "" {
		parts = append(parts, trimmed)
	}
	parts = append(parts, scene, userID, now.Format("20060102"))

	return strings.Trim(path.Join(parts...), "/") + "/"
}

func normalizePolicyPrefix(prefix string) string {
	trimmed := strings.Trim(normalizeObjectPath(prefix), "/")
	if trimmed == "" {
		return ""
	}

	return trimmed + "/"
}

func normalizeObjectPath(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	return strings.TrimPrefix(value, "/")
}

func sanitizePathSegment(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}

	sanitized := unsafePathSegmentPattern.ReplaceAllString(value, "-")
	sanitized = strings.Trim(sanitized, "-")
	if sanitized == "" {
		return fallback
	}

	return sanitized
}

func valueString(pointer *string) string {
	if pointer == nil {
		return ""
	}

	return *pointer
}

func valueUint64(pointer *uint64) uint64 {
	if pointer == nil {
		return 0
	}

	return *pointer
}

func stringPtr(value string) *string {
	return &value
}

func uint64Ptr(value uint64) *uint64 {
	return &value
}
