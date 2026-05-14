# Tasks

- [x] Task 1: 定义统一数据模型（types 层）
  - [x] 1.1 定义 `ContentType` 枚举常量（market / errand / resource / lost_found / carpool / meetup）
  - [x] 1.2 定义统一 `Status` 枚举常量（reviewing / published / rejected / offline / cancelled / accepted / open / full / resolved）
  - [x] 1.3 定义 `CommunityContent` 公共结构体，包含所有公共字段 + `TypePayload` + `ExtJSON`（`_id` 使用 MongoDB ObjectID，不存储 `publisher`/`publisher_initial`）
  - [x] 1.4 定义六种类型载荷结构体：`MarketPayload` / `ErrandPayload` / `ResourcePayload` / `LostFoundPayload` / `CarpoolPayload` / `MeetupPayload`
  - [x] 1.5 定义统一查询结构：`ContentFilter` / `FeedFilter`
  - [x] 1.6 保留现有发布请求结构（transport 层 DTO）
  - [x] 1.7 移除旧的六个 Item 结构体和 Extra 结构体

- [x] Task 2: 重构 Repository 接口和 MongoDB 实现（repo 层）
  - [x] 2.1 定义统一 `Repository` 接口：`Save` / `GetByID` / `Update` / `ListByType` / `ListForFeed`（移除 `NextID`）
  - [x] 2.2 实现 `MongoRepository`：单集合 `community_content`，`_id` 使用 MongoDB 自动 ObjectID
  - [x] 2.3 移除 `mongoEnvelope` 信封结构和 `campus_life_sequences` 集合
  - [x] 2.4 实现软删除过滤（`deleted_at == null`）
  - [x] 2.5 实现 `type_payload` 的 BSON 序列化/反序列化（使用 `bson.M` 适配多态）
  - [x] 2.6 创建 MongoDB 索引：`{content_type, status, created_at}` / `{status, created_at}` / `{publisher_user_id, created_at}` / `{deleted_at}` sparse

- [x] Task 3: 重构 Service 层
  - [x] 3.1 统一发布逻辑：根据 `content_type` 构造 `CommunityContent`，初始 `status` 为 `reviewing`，`_id` 由 MongoDB 自动生成
  - [x] 3.2 统一状态流转逻辑：用 `status` 单字段替代 `review_status` + `status` 双字段
  - [x] 3.3 移除 `mergeErrandStatus` / `mergeMeetupStatus` / `normalizeReviewStatus` 合并函数
  - [x] 3.4 统一可见性判断：`shouldExposeContent` 改为检查 `status` 是否为可见状态
  - [x] 3.5 统一联系方式裁剪逻辑：`contact` 从公共字段读取
  - [x] 3.6 统一审核逻辑：`UpdateReviewStatus` 直接修改 `status` 字段
  - [x] 3.7 类型特化逻辑保留在 service 层（跑腿接单、组局报名等），通过 `type_payload` 读写类型特有数据
  - [x] 3.8 统一 Feed 流逻辑：使用 `ListByType` 按类型查询，按 `created_at` 排序
  - [x] 3.9 发布者信息实时解析：service 层根据 `publisher_user_id` 从 principal 解析昵称和首字母

- [x] Task 4: 适配 Transport 层
  - [x] 4.1 Handler 方法无需修改（service 方法签名保持不变）
  - [x] 4.2 保持现有 API 路径不变
  - [x] 4.3 响应结构保持与现有行为一致（ID 格式从 `{type}-{序号}` 变为 ObjectID）

- [x] Task 5: 更新 OpenAPI 契约
  - [x] 5.1 更新 `packages/contracts/openapi/api-server.yaml` 中的 CampusLife 相关 Schema
  - [x] 5.2 `id` 格式改为 ObjectID hex string，`status` 枚举统一

- [x] Task 6: 数据迁移脚本
  - [x] 6.1 编写迁移脚本：将六个旧集合数据转换并写入 `community_content`
  - [x] 6.2 `review_status` → `status` 映射逻辑
  - [x] 6.3 `Extra` → `type_payload` 映射逻辑
  - [x] 6.4 补全审计字段（`updated_at` / `updated_by` / `deleted_at`）
  - [x] 6.5 旧 `{type}-{序号}` ID 映射：迁移后使用新 ObjectID，旧 ID 不再保证兼容

- [x] Task 7: 更新模块文档
  - [x] 7.1 更新 `campus_life/README.md` 反映新集合设计
  - [x] 7.2 更新 `ARCHITECTURE.md` 中的 campus_life 存储架构描述

# Task Dependencies

- [Task 2] depends on [Task 1]
- [Task 3] depends on [Task 1, Task 2]
- [Task 4] depends on [Task 3]
- [Task 5] depends on [Task 1]
- [Task 6] depends on [Task 1, Task 2]
- [Task 7] depends on [Task 1, Task 2, Task 3]
- Task 1 可独立开始
- Task 5 可与 Task 2/3 并行
