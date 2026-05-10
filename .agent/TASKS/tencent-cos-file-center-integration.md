## Task: 腾讯 COS 文件管理接入

**ID:** tencent-cos-file-center-integration
**Label:** API Server: 腾讯 COS 文件管理接入
**Description:** 为 `services/api-server` 接入腾讯 COS，提供 STS 临时凭证、下载预签名、对象存储健康探针，并把小程序上传链路改为“只存对象路径、读取时签 URL”
**Type:** Feature
**Status:** Completed
**Priority:** High
**Created:** 2026-05-10
**Updated:** 2026-05-10
**PRD:** [tencent-cos-file-center-integration.md](/Users/liangluo/code/weouc2026/.agent/PRDS/tencent-cos-file-center-integration.md)

**Progress:** 已完成 `storage_provider + file_center` 落地，`campus_life` 与微信小程序上传/读取链路已切换为“只存对象路径、读取时签 URL”，并补齐 Docker、OpenAPI 与就绪检查。
