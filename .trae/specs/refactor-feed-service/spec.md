# Feed 服务重构 Spec

## Why

当前 `service.go` 文件长达 2271 行，混合了 6 种内容类型的业务逻辑、内存过滤、分页和辅助函数，导致：
- 难以维护和定位逻辑
- `ListFeed` 逐类型调用 `ListByType` 后在内存中过滤、排序、分页，而非利用 MongoDB 索引直接在数据库层完成过滤
- `user_role` 作为客户端传入的过滤参数，违反"权限判断只在后端完成"的约束，且语义不清晰

## What Changes

- 将 `service.go` 按内容类型拆分为独立文件（`feed_service.go`、`market_service.go`、`errand_service.go`、`resource_service.go`、`lostfound_service.go`、`carpool_service.go`、`meetup_service.go`、`helpers.go`）
- 将过滤逻辑下沉到 repo 层，利用 MongoDB 查询条件和索引直接在数据库层过滤，service 层不再做内存过滤
- **BREAKING**：移除 `user_role` 请求参数，替换为服务端"主态/客态"可见性逻辑
- 增强 `ListForFeed` repo 方法，支持完整的查询条件（feedTypes、keyword、status 过滤、publisher 过滤）
- 新增复合索引支持主态/客态查询模式

## Impact

- Affected specs: `unified-community-content`（repo 接口变更）、`unified-status-design`（状态过滤逻辑变更）
- Affected code:
  - `services/api-server/internal/modules/campus_life/service/service.go` — 拆分为多个文件
  - `services/api-server/internal/modules/campus_life/repo/repository.go` — 接口扩展
  - `services/api-server/internal/modules/campus_life/repo/mongo.go` — 查询逻辑增强
  - `services/api-server/internal/modules/campus_life/transport/handler.go` — 移除 user_role 参数解析
  - `services/api-server/internal/modules/campus_life/types/models.go` — FeedFilter/FeedQuery/ContentFilter 移除 UserRole 字段，新增 ViewMode 字段
  - `apps/miniapp-wechat/api/modules/feed.js` — 移除 user_role 参数
  - `apps/miniapp-wechat/api/modules/errand.js` — 移除 user_role 参数（如存在）
  - `apps/miniapp-wechat/api/modules/meetup.js` — 移除 user_role 参数（如存在）

***

## ADDED Requirements

### Requirement: Service 按内容类型拆分文件

系统 SHALL 将 `service.go`（2271 行）按内容类型拆分为独立文件，每个文件负责一种内容类型的 CRUD 和列表逻辑。

#### Scenario: 文件拆分结构

- **WHEN** 开发者查看 campus_life/service 目录
- **THEN** 目录结构 MUST 为：
  - `service.go` — Service struct 定义、New 构造函数、通用工具函数（marshalPayload、unmarshalPayload）
  - `feed_service.go` — ListFeed 方法
  - `market_service.go` — ListMarket、GetMarketDetail、PublishMarket、FavoriteMarket、DeleteMarket
  - `errand_service.go` — ListErrands、GetErrandDetail、PublishErrand、AcceptErrand、CancelErrandPublish、CancelErrandAccept
  - `resource_service.go` — ListResources、GetResourceDetail、PublishResource、DeleteResource
  - `lostfound_service.go` — ListLostFound、GetLostFoundDetail、PublishLostFound、DeleteLostFound、MarkLostFoundResolved
  - `carpool_service.go` — ListCarpools、GetCarpoolDetail、PublishCarpool、DeleteCarpool
  - `meetup_service.go` — ListMeetups、GetMeetupDetail、PublishMeetup、JoinMeetup、CancelMeetupJoin、CancelMeetupPublish
  - `review_service.go` — ListReviewQueue、UpdateReviewStatus
  - `helpers.go` — 所有纯函数辅助方法（shouldExposeContent、canViewContact、matchKeyword、listEnvelope、paginateRows 等）

#### Scenario: 所有文件在同一 package

- **WHEN** 编译 campus_life 模块
- **THEN** 所有拆分后的文件 MUST 仍在 `package service` 中，无包名变更

### Requirement: 过滤逻辑下沉到 Repo 层

系统 SHALL 将当前在 service 层通过内存遍历实现的过滤逻辑（状态过滤、类型过滤、关键词过滤、发布者过滤）下沉到 repo 层，利用 MongoDB 查询条件和索引直接在数据库层完成。

#### Scenario: ListFeed 使用 ListForFeed 而非逐类型 ListByType

- **WHEN** 调用 ListFeed
- **THEN** service 层 MUST 调用 repo 的 `ListForFeed` 方法，一次性从数据库获取跨类型数据
- **THEN** 不再逐类型调用 `ListByType` 后在内存中合并

#### Scenario: 状态过滤在数据库层完成

- **WHEN** 查询内容列表
- **THEN** 状态过滤（published、reviewing 等）MUST 通过 MongoDB 查询条件完成，不在 service 层遍历过滤

#### Scenario: 类型过滤在数据库层完成

- **WHEN** 查询 feed 列表且指定了 feedTypes
- **THEN** 类型过滤 MUST 通过 `content_type: { $in: feedTypes }` 查询条件完成

#### Scenario: 发布者过滤在数据库层完成

- **WHEN** 查询"我发布的"内容
- **THEN** 发布者过滤 MUST 通过 `publisher_user_id` 查询条件完成，不在 service 层遍历匹配

#### Scenario: 关键词过滤在数据库层完成

- **WHEN** 查询带关键词的内容列表
- **THEN** 关键词过滤 MUST 通过 MongoDB `$or` + `$regex` 查询条件完成（对 title、desc 字段），不在 service 层遍历匹配

#### Scenario: 分页在数据库层完成

- **WHEN** 查询列表
- **THEN** 分页 MUST 通过 MongoDB `skip`/`limit` 完成，不在 service 层内存分页
- **THEN** `total` 计数 MUST 通过 MongoDB `countDocuments` 完成，不在 service 层 `len(filtered)`

#### Scenario: 各类型列表接口同样下沉过滤

- **WHEN** 调用 ListMarket、ListErrands、ListResources 等类型列表接口
- **THEN** category、keyword、status 等过滤条件 MUST 通过 repo 层查询条件完成，不在 service 层遍历过滤

### Requirement: 主态/客态可见性逻辑

系统 SHALL 移除客户端传入的 `user_role` 过滤参数，替换为服务端自动判断的"主态/客态"可见性逻辑。

#### Scenario: 主态定义

- **WHEN** 用户查看与自己相关的内容（自己是发布者、接单者、参与者）
- **THEN** 系统返回该用户作为相关方的所有内容，包括：`reviewing`、`published`、`rejected`、`offline`、`cancelled`、`accepted`、`open`、`full`、`resolved` 等所有状态
- **THEN** 同时返回其他用户已发布（`published`/`open`/`accepted`/`full`/`resolved`）的在线内容

#### Scenario: 客态定义

- **WHEN** 用户查看与自己无关的内容（非发布者、非接单者、非参与者）
- **THEN** 系统只返回审核通过且在线的内容，状态限定为：`published`、`open`、`accepted`、`full`、`resolved`

#### Scenario: Feed 列表的主态/客态逻辑

- **WHEN** 已认证用户请求 Feed 列表
- **THEN** 系统通过一次查询返回两部分内容的并集：
  - 该用户作为发布者/接单者/参与者的所有状态内容（主态）
  - 其他用户已发布在线的内容（客态）
- **THEN** 查询条件 MUST 为：
  ```
  $or: [
    { publisher_user_id: currentUserID },  // 主态：自己的所有内容
    { status: { $in: visibleStatuses } }    // 客态：已发布在线的内容
  ]
  ```
  其中 `visibleStatuses` = `["published", "open", "accepted", "full", "resolved"]`

#### Scenario: 未认证用户

- **WHEN** 未认证用户请求 Feed 列表
- **THEN** 系统只返回客态内容（`status: { $in: visibleStatuses }`）

#### Scenario: 管理员

- **WHEN** 拥有 `campus_life:moderate` 权限的管理员请求 Feed 列表
- **THEN** 系统返回所有状态的内容（不限制状态过滤）

#### Scenario: 类型列表接口的主态/客态逻辑

- **WHEN** 已认证用户请求特定类型列表（如 ListErrands、ListMeetups）
- **THEN** 系统同样应用主态/客态逻辑：
  - 主态：该用户作为相关方（发布者/接单者/参与者）的所有状态内容
  - 客态：其他用户已发布在线的内容

#### Scenario: 移除 user_role 请求参数

- **WHEN** 客户端调用 Feed 列表或类型列表接口
- **THEN** 不再接受 `user_role` 查询参数
- **THEN** 服务端根据请求用户的 principal 自动判断可见性

### Requirement: 增强的 Repo 接口

系统 SHALL 扩展 Repository 接口以支持数据库层过滤。

#### Scenario: ListForFeed 接口增强

- **WHEN** 调用 `ListForFeed`
- **THEN** FeedFilter MUST 支持以下字段：
  - `FeedTypes []string` — 内容类型过滤
  - `Keyword string` — 关键词过滤
  - `PublisherUserID string` — 发布者过滤（主态查询）
  - `AcceptorUserID string` — 接单者过滤（跑腿主态查询）
  - `ParticipantUserID string` — 参与者过滤（组局主态查询）
  - `VisibleStatuses []string` — 客态可见状态列表
  - `IncludeAllStatus bool` — 是否包含所有状态（管理员）
  - `Pagination` — 分页

#### Scenario: ListByType 接口增强

- **WHEN** 调用 `ListByType`
- **THEN** ContentFilter MUST 支持以下字段：
  - `Statuses []string` — 多状态过滤（替代单一 Status）
  - `Keyword string` — 关键词过滤
  - `PublisherUserID string` — 发布者过滤
  - `AcceptorUserID string` — 接单者过滤
  - `ParticipantUserID string` — 参与者过滤
  - `VisibleStatuses []string` — 客态可见状态列表
  - `IncludeAllStatus bool` — 管理员标记
  - `CurrentUserID string` — 当前用户 ID（用于主态/客态查询）

#### Scenario: 返回总数

- **WHEN** 调用 ListForFeed 或 ListByType
- **THEN** 方法签名 MUST 返回 `([]CommunityContent, int64, error)`，其中 `int64` 为满足过滤条件的总记录数

### Requirement: MongoDB 索引增强

系统 SHALL 新增索引以支持主态/客态查询模式。

#### Scenario: 新增索引

- **WHEN** 部署 community_content 集合
- **THEN** MUST 在现有索引基础上新增以下索引：
  - 复合索引：`{ publisher_user_id: 1, status: 1, created_at: -1 }`（主态查询：按发布者+状态+时间）
  - 复合索引：`{ status: 1, content_type: 1, created_at: -1 }`（客态查询：按状态+类型+时间，与现有索引合并）
  - 文本索引或复合索引：`{ title: "text", desc: "text" }`（关键词搜索，或使用 $regex + 复合索引）

### Requirement: 响应中保留 user_role 语义字段

系统 SHALL 在响应中保留 `user_role` 字段用于前端展示，但不再作为请求参数。

#### Scenario: 响应中计算 user_role

- **WHEN** 返回内容详情或列表项
- **THEN** 服务端根据当前用户与内容的关系计算 `user_role`（publisher/acceptor/participant/viewer），作为响应字段返回
- **THEN** 前端仅用于展示判断（如"编辑"按钮可见性），不做权限裁决

***

## MODIFIED Requirements

### Requirement: Feed 列表查询

Feed 列表查询 SHALL 通过 `ListForFeed` 方法一次性从数据库获取跨类型数据，不再逐类型调用 `ListByType` 后在内存合并。

- **WHEN** 查询首页动态流
- **THEN** 通过 `ListForFeed` 方法跨类型查询，应用主态/客态过滤条件，按 `created_at` 降序排列，数据库层分页

### Requirement: 类型列表查询

各类型列表查询 SHALL 通过增强的 `ListByType` 方法在数据库层完成过滤，不再在 service 层遍历过滤。

- **WHEN** 查询特定类型的内容列表
- **THEN** category、keyword、status、publisher 等过滤条件通过 MongoDB 查询条件完成

### Requirement: Handler 参数解析

Handler 层 SHALL 移除 `user_role` 查询参数的解析，不再传递给 service 层。

- **WHEN** 解析请求参数
- **THEN** 不再从 query string 中读取 `user_role`

***

## REMOVED Requirements

### Requirement: user_role 作为请求过滤参数

**Reason**: 违反"权限判断只在后端完成"约束，客户端不应决定可见性范围。替换为服务端主态/客态自动判断。
**Migration**: 客户端移除 `user_role` 查询参数传递；服务端根据 principal 自动判断可见性。

### Requirement: 内存过滤和分页

**Reason**: 性能差，无法利用数据库索引，数据量增长后不可扩展。替换为数据库层过滤和分页。
**Migration**: 所有过滤条件通过 MongoDB 查询条件传递，分页通过 skip/limit 完成，总数通过 countDocuments 获取。
