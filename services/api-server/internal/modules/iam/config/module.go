package config

import "time"

type ModuleConfig struct {
	DefaultRoles      []string
	CaptchaTTL        time.Duration
	AccessTokenTTL    time.Duration
	StudentPermission string
	AllowPasswordLogin bool
}

func New(accessTokenTTL time.Duration) ModuleConfig {
	return ModuleConfig{
		DefaultRoles:      []string{"student"},
		CaptchaTTL:        5 * time.Minute,
		AccessTokenTTL:    accessTokenTTL,
		StudentPermission: "contact:view",
		AllowPasswordLogin: true,
	}
}
