# 腾讯 COS 文件管理接入 PRD

## 背景

当前仓库已经为微信小程序预留了 COS 直传客户端逻辑，并约定了：

- `GET /api/upload/cos-sts`
- `POST /api/upload/presigned-get`

但 `api-server` 侧仍未提供真实实现，现状存在三个结构性问题：

- 小程序上传链路只能降级或失败，无法完成真实文件管理。
- 二手/跑腿发布当前直接把临时下载 URL 写入业务数据，天然会过期。
- 资料发布虽然传的是 `file_paths`，但后端读取时仍拼接示例域名，不是真实 COS 文件访问。

因此需要把 COS 接入提升为正式的文件管理能力，而不是继续停留在“客户端已有 SDK、服务端未实现”的半链路状态。

## 目标

在保持现有架构约束的前提下，完成基于腾讯 COS 的文件管理接入，使以下能力可真实跑通：

- 服务端签发受限路径的 COS STS 临时上传凭证
- 服务端按对象路径生成 GET 预签名下载链接
- `readyz` 对对象存储给出真实健康状态
- 微信小程序上传后只持久化对象路径，不再持久化临时下载 URL
- 二手、跑腿、资料读取时由后端动态签发可访问 URL

## 非目标

- 本轮不做完整的文件元数据持久化表设计
- 本轮不做管理端文件列表、删除、回收站等后台能力
- 本轮不做断点续传、多段上传调优
- 本轮不接入其他云厂商对象存储

## 约束

- 外部存储接入必须经 `internal/providers/storage_provider`
- HTTP 上传接口能力应放入 `internal/modules/file_center`
- 业务模块不能直连 COS SDK
- 业务数据层只存稳定对象路径，不存会过期的签名 URL

## 范围

### 1. 配置与运行时

- 在 `AppConfig` 中增加 COS 配置
- 支持通过环境变量注入：
  - `API_SERVER_COS_ENABLED`
  - `API_SERVER_COS_SECRET_ID`
  - `API_SERVER_COS_SECRET_KEY`
  - `API_SERVER_COS_BUCKET`
  - `API_SERVER_COS_APP_ID`
  - `API_SERVER_COS_REGION`
  - `API_SERVER_COS_PATH_PREFIX`
  - `API_SERVER_COS_STS_DURATION`
  - `API_SERVER_COS_PRESIGNED_GET_TTL`
- Docker compose 与 `.env.example` 同步补齐

### 2. Storage Provider

- 新增 `storage_provider`
- 提供：
  - 生成上传临时凭证
  - 生成下载预签名 URL
  - 对象存储健康检查
- 默认支持真实 COS provider
- 未启用 COS 时提供明确的降级错误

### 3. File Center 模块

- 新增 `file_center` 模块最小闭环：
  - `GET /api/upload/cos-sts`
  - `POST /api/upload/presigned-get`
- 路由需要登录保护
- `path_prefix` 应带用户隔离维度，避免所有用户共享上传前缀

### 4. 业务链路改造

- 小程序二手/跑腿发布改为传对象 `path`
- 后端二手/跑腿/资料读取时签 URL 再返回
- 保持已有种子数据和历史 URL 兼容：若值已经是 `http/https`，则直接透传

## 验收标准

- `GET /api/upload/cos-sts` 可返回小程序 COS SDK 可直接使用的字段
- `POST /api/upload/presigned-get` 可基于对象路径返回可访问下载链接
- `readyz` 中 `object_storage` 不再固定为 `skipped`
- 二手、跑腿、资料链路不再持久化临时访问 URL
- `apps/miniapp-wechat/API.md`、`services/api-server/README.md`、`docs/PLANS.md` 同步更新

## 风险

- 腾讯云 STS 权限配置错误会导致客户端上传失败
- 若继续把签名 URL 入库，问题会在数小时后才暴露，联调期容易误判
- COS 桶策略、地域、AppID/BucketName 配置不一致会直接导致签名或访问失败

## 本轮实施顺序

1. 先补 COS 配置与 `storage_provider`
2. 再落 `file_center` 接口
3. 然后改 `campus_life + miniapp`
4. 最后补文档和 Docker 配置

## 实施结果

- 已新增 `internal/providers/storage_provider`，落地腾讯 COS STS 临时凭证、GET 预签名与桶健康检查
- 已新增 `internal/modules/file_center`，提供 `GET /api/upload/cos-sts` 与 `POST /api/upload/presigned-get`
- 已将 `readyz` 的 `object_storage` 从静态占位改为真实探测
- 已将二手、跑腿、资料链路改为持久化对象路径，读取阶段再由服务端签发访问 URL
- 已同步更新小程序上传逻辑、Docker 运行配置、OpenAPI 契约与相关 README/计划文档
