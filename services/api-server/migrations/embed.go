package migrations

import "embed"

// Files 包含服务端内置迁移脚本。
//
//go:embed *.sql
var Files embed.FS
