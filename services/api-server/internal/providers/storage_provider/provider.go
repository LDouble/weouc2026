package storage_provider

import "context"

type Provider interface {
	IssueUploadCredentials(ctx context.Context, input IssueUploadCredentialsInput) (UploadCredentials, error)
	PresignGetURL(ctx context.Context, objectPath string) (string, error)
	Check(ctx context.Context) error
}

type IssueUploadCredentialsInput struct {
	UserID string
	Scene  string
}

type UploadCredentials struct {
	TmpSecretID  string
	TmpSecretKey string
	SessionToken string
	StartTime    int64
	ExpiredTime  int64
	Bucket       string
	Region       string
	PathPrefix   string
}
