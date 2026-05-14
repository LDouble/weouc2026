# Tasks

- [x] Task 1: 扩展 types 层 — 更新 FeedFilter、ContentFilter、FeedQuery 等类型定义
  - [x] SubTask 1.1: FeedFilter 新增 PublisherUserID、AcceptorUserID、ParticipantUserID、VisibleStatuses、IncludeAllStatus 字段，移除 UserRole
  - [x] SubTask 1.2: ContentFilter 新增 Statuses []string（替代单一 Status）、Keyword、PublisherUserID、AcceptorUserID、ParticipantUserID、VisibleStatuses、IncludeAllStatus、CurrentUserID 字段，移除 UserRole
  - [x] SubTask 1.3: FeedQuery 移除 UserRole 字段
  - [x] SubTask 1.4: ErrandQuery、MeetupQuery 移除 UserRole 字段
  - [x] SubTask 1.5: 新增 VisibleStatuses 变量

- [x] Task 2: 扩展 repo 接口和 MongoDB 实现 — 支持数据库层过滤
  - [x] SubTask 2.1: Repository 接口 ListForFeed 签名变更为返回 `([]CommunityContent, int64, error)`
  - [x] SubTask 2.2: Repository 接口 ListByType 签名变更为返回 `([]CommunityContent, int64, error)`
  - [x] SubTask 2.3: MongoRepository.ListForFeed 实现：根据 FeedFilter 构建 MongoDB 查询条件（feedTypes、keyword $regex、主态/客态 $or 查询、分页、countDocuments）
  - [x] SubTask 2.4: MongoRepository.ListByType 实现：根据 ContentFilter 构建 MongoDB 查询条件（statuses $in、keyword $regex、publisher_user_id、主态/客态 $or 查询、分页、countDocuments）
  - [x] SubTask 2.5: EnsureIndexes 新增索引：`{ publisher_user_id: 1, status: 1, created_at: -1 }`
  - [x] SubTask 2.6: 现有索引 `{ content_type: 1, status: 1, created_at: -1 }` 已覆盖客态查询

- [x] Task 3: 重构 service 层 — 拆分文件 + 移除内存过滤
  - [x] SubTask 3.1: 创建 `feed_service.go`，将 ListFeed 方法移入，改为调用 repo.ListForFeed，移除逐类型 ListByType + 内存合并逻辑
  - [x] SubTask 3.2: 创建 `market_service.go`，将 ListMarket、GetMarketDetail、PublishMarket、FavoriteMarket、DeleteMarket 移入，ListMarket 改用增强的 ListByType
  - [x] SubTask 3.3: 创建 `errand_service.go`，将 ListErrands、GetErrandDetail、PublishErrand、AcceptErrand、CancelErrandPublish、CancelErrandAccept 移入，ListErrands 改用增强的 ListByType
  - [x] SubTask 3.4: 创建 `resource_service.go`，将 ListResources、GetResourceDetail、PublishResource、DeleteResource 移入
  - [x] SubTask 3.5: 创建 `lostfound_service.go`，将 ListLostFound、GetLostFoundDetail、PublishLostFound、DeleteLostFound、MarkLostFoundResolved 移入
  - [x] SubTask 3.6: 创建 `carpool_service.go`，将 ListCarpools、GetCarpoolDetail、PublishCarpool、DeleteCarpool 移入
  - [x] SubTask 3.7: 创建 `meetup_service.go`，将 ListMeetups、GetMeetupDetail、PublishMeetup、JoinMeetup、CancelMeetupJoin、CancelMeetupPublish 移入
  - [x] SubTask 3.8: 创建 `review_service.go`，将 ListReviewQueue、UpdateReviewStatus 移入
  - [x] SubTask 3.9: 创建 `helpers.go`，将所有纯函数辅助方法移入
  - [x] SubTask 3.10: 保留 `service.go` 仅包含 Service struct 定义、New 构造函数、marshalPayload、unmarshalPayload、recordAudit
  - [x] SubTask 3.11: 所有列表方法移除内存过滤逻辑，改用 repo 层查询条件
  - [x] SubTask 3.12: 所有列表方法使用 repo 返回的 total 计数替代 len(filtered)

- [x] Task 4: 实现主态/客态可见性逻辑
  - [x] SubTask 4.1: 新增 `buildVisibilityFilter` 辅助函数
  - [x] SubTask 4.2: ListFeed 使用 buildVisibilityFilter 构建 FeedFilter
  - [x] SubTask 4.3: ListErrands 使用 buildVisibilityFilter 构建 ContentFilter（额外支持 acceptor_user_id 主态）
  - [x] SubTask 4.4: ListMeetups 使用 buildVisibilityFilter 构建 ContentFilter（额外支持 participant_user_ids 主态）
  - [x] SubTask 4.5: 其他类型列表（ListMarket、ListResources 等）使用 buildVisibilityFilter
  - [x] SubTask 4.6: 管理员（campus_life:moderate）跳过状态过滤，IncludeAllStatus = true

- [x] Task 5: 更新 transport 层 — 移除 user_role 参数
  - [x] SubTask 5.1: Handler.ListFeed 移除 `c.Query("user_role")` 解析
  - [x] SubTask 5.2: Handler.ListErrands 移除 `c.Query("user_role")` 解析
  - [x] SubTask 5.3: Handler.ListMeetups 移除 `c.Query("user_role")` 解析

- [x] Task 6: 更新小程序客户端 — 移除 user_role 参数
  - [x] SubTask 6.1: `apps/miniapp-wechat/api/modules/feed.js` 移除 user_role 参数
  - [x] SubTask 6.2: 搜索并更新小程序中其他传递 user_role 的 API 调用（errand.js、meetup.js、profileService.js、publish/index.js、accepted/index.js）

- [x] Task 7: 验证与测试
  - [x] SubTask 7.1: 编译通过（`go build ./internal/modules/campus_life/...`）
  - [x] SubTask 7.2: go vet 通过
  - [x] SubTask 7.3: 验证 ListFeed 接口返回数据正确（主态/客态逻辑）
  - [x] SubTask 7.4: 验证各类型列表接口返回数据正确
  - [x] SubTask 7.5: 验证详情接口 user_role 响应字段仍正确返回

# Task Dependencies

- [Task 2] depends on [Task 1]（repo 接口依赖 types 定义）
- [Task 3] depends on [Task 2]（service 层依赖 repo 接口变更）
- [Task 4] depends on [Task 3]（主态/客态逻辑在 service 层实现）
- [Task 5] depends on [Task 1]（handler 依赖 types 定义变更）
- [Task 6] depends on [Task 5]（客户端依赖后端接口变更）
- [Task 7] depends on [Task 4, Task 5, Task 6]（验证依赖所有变更完成）
- [Task 1] and [Task 5] 可并行（types 和 handler 变更互不依赖）
