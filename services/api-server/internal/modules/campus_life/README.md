# campus_life

当前已落地：

- 首页动态流
- 二手：列表、详情、发布、收藏
- 跑腿：列表、详情、发布、接单、取消发布、取消接单
- 资料：列表、详情、发布
- 失物招领：列表、详情、发布、标记已找到
- 拼车：列表、详情、发布
- 组局：列表、详情、发布、报名、取消报名、取消组局
- 审核列表与审核状态更新接口（接收 Action 字段：`review_approve` / `review_reject` / `offline_by_admin`）
- 管理端校园生活列表复用真实审核分页数据

存储架构：

- 统一 `community_content` MongoDB 集合，以聚合根组织
- "公共字段 + 类型载荷"模型：公共字段（_id, content_type, title, desc, status, publisher_user_id, contact, images, tags, 审计字段）+ `type_payload`（6 种结构化载荷）
- `_id` 使用 MongoDB 自动生成的 ObjectID
- 统一状态模型：单一 `status` 字段，包含 reviewing/published/rejected/offline/cancelled/accepted/open/full/resolved
- 统一审计字段：created_at/updated_at/created_by/updated_by/deleted_at
- `ext_json` 保留用于未来扩展，严格限制使用范围
- 软删除通过 `deleted_at` 字段实现
- 状态转换日志写入 `state_transition_logs` 集合

BMFS 状态驱动（已落地）：

- 通用状态机引擎 `internal/platform/bmfs`：声明式定义状态、动作、转换规则；支持 guard 守卫和 onTransition 钩子；自动派生 `can_xxx`
- 6 类内容状态机定义在 `types/statemachines.go`，通过 `GetMachine(contentType)` 工厂获取
- 所有状态变更操作通过 `bmfs.Execute()` 执行，service 层不再直接设置 `item.Status`（发布创建初始状态除外）
- `can_xxx` 从 `bmfs.AvailableActions()` 派生，不再手写条件分支
- 客户端只能提交动作命令，不能直接提交目标状态

约束：

- 业务规则只在 service 层
- 数据访问只经 repo
- 文件访问只保存稳定对象路径，展示时由后端按需签 URL
- 公开列表与详情只展示 published 状态内容；发布者与审核员可继续访问待审内容
- 发布者信息（昵称、首字母）不冗余存储，读取时从 principal 实时解析
- 不在客户端复制状态机规则

MongoDB 索引：

- `{content_type, status, created_at}` 按类型+状态列表查询
- `{status, created_at}` Feed 流查询
- `{publisher_user_id, created_at}` 我的发布查询
- `{deleted_at}` sparse 软删除过滤
- `state_transition_logs`：`{content_type, content_id, created_at}` 和 `{actor_user_id, created_at}`
