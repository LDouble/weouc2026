package types

type PresignedGetRequest struct {
	Path string `json:"path"`
}

type UploadSTSCredentials struct {
	TmpSecretID  string `json:"tmp_secret_id"`
	TmpSecretKey string `json:"tmp_secret_key"`
	SessionToken string `json:"session_token"`
	StartTime    int64  `json:"start_time"`
	ExpiredTime  int64  `json:"expired_time"`
	Bucket       string `json:"bucket"`
	Region       string `json:"region"`
	PathPrefix   string `json:"path_prefix"`
}

type PresignedGetResponse struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}
