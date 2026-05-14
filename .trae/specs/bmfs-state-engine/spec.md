# BMFS 状态驱动引擎 Spec

## Why

当前 campus\_life 模块的状态流转逻辑以手写条件分支的形式散落在 6 个 service 文件和 helpers.go 中，can\_xxx 派生语义同样零散内联。每新增内容类型都需要复制一整套状态判断，状态漂移越来越严重。通过 BMFS 收口状态机，保证多端只消费同一套业务语义。

## What Changes

- 新增 `internal/platform/bmfs` 包，实现通用状态机引擎
- 为 6 类社区内容（errand / meetup / market / lost\_found / carpool / resource）声明式定义状态机
- service 层通过 BMFS 执行动作（而非直接设置 status 字段）
- 所有 can\_xxx 从 BMFS 派生，移除 helpers.go 和各 service 中的手写 can\_xxx 逻辑
- 客户端只能提交动作命令（如 `accept`、`cancel`），不允许直接提交目标状态
- 审核阶段并入主状态机，BMFS 统一管理审核流转

## Impact

- Affected specs: unified-status-design（已完成，BMFS 是其工程化落地）
- Affected code:
  - `internal/platform/bmfs/`（新增）
  - `internal/modules/campus_life/service/` 下所有 service 文件
  - `internal/modules/campus_life/service/helpers.go` 中的 can\_xxx 函数
  - `internal/modules/campus_life/types/models.go` 中的状态常量
  - `internal/modules/campus_life/transport/` 中 handler 对动作命令的接收

## ADDED Requirements

### Requirement: BMFS 状态机引擎核心

系统 SHALL 提供一个通用状态机引擎 `bmfs`，支持以下能力：

1. **声明式状态机定义**：通过 Go 代码声明状态（State）、动作（Action）、转换规则（Transition）
2. **守卫条件**：每个转换可附加 guard 函数，运行时根据上下文决定转换是否允许
3. **副作用钩子**：每个转换可附加 onTransition 回调，用于触发副作用（如清空 acceptor、更新参与人数）
4. **can\_xxx 派生**：给定当前状态和上下文，引擎能计算所有可用动作，派生 can\_xxx 布尔值

#### Scenario: 定义状态机并执行动作

- **WHEN** 开发者为某内容类型定义了状态机，包含状态 A → 状态 B 的动作 `do_something`
- **AND** 某条内容当前处于状态 A
- **AND** service 层调用 `bmfs.Execute(ctx, machine, item, "do_something", context)`
- **THEN** BMFS 检查 guard 条件通过后，将状态转为 B，并执行 onTransition 回调
- **AND** 返回新状态和派生的 can\_xxx 集合

#### Scenario: 守卫条件拒绝非法转换

- **WHEN** 某条内容当前处于状态 A
- **AND** service 层调用 `bmfs.Execute(ctx, machine, item, "do_something", context)`
- **AND** guard 条件不满足（如用户不是发布者）
- **THEN** BMFS 返回错误，状态不变

#### Scenario: 查询可用动作

- **WHEN** service 层调用 `bmfs.AvailableActions(ctx, machine, item, context)`
- **THEN** BMFS 返回当前状态下所有 guard 通过的动作列表
- **AND** 派生 can\_xxx 映射（如 `can_cancel: true`, `can_accept: false`）

### Requirement: 六类社区内容状态机定义

系统 SHALL 为以下 6 类内容定义各自的状态机：

#### Errand（跑腿）状态机

| 当前状态      | 动作                 | 目标状态      | 守卫条件        |
| --------- | ------------------ | --------- | ----------- |
| reviewing | review\_approve    | published | 调用方为审核员     |
| reviewing | review\_reject     | rejected  | 调用方为审核员     |
| reviewing | cancel             | cancelled | 调用方为发布者     |
| published | accept             | accepted  | 调用方非发布者且已认证 |
| published | cancel             | cancelled | 调用方为发布者     |
| published | offline\_by\_admin | offline   | 调用方为审核员     |
| accepted  | cancel\_accept     | published | 调用方为接单者     |
| accepted  | offline\_by\_admin | offline   | 调用方为审核员     |
| rejected  | cancel             | cancelled | 调用方为发布者     |
| rejected  | review\_reapprove  | reviewing | 调用方为发布者     |

#### Meetup（组局）状态机

| 当前状态      | 动作                 | 目标状态        | 守卫条件                               |
| --------- | ------------------ | ----------- | ---------------------------------- |
| reviewing | review\_approve    | open / full | 调用方为审核员；根据参与人数决定 open 或 full       |
| reviewing | review\_reject     | rejected    | 调用方为审核员                            |
| reviewing | cancel             | cancelled   | 调用方为发布者                            |
| open      | join               | open / full | 调用方非发布者且已认证；剩余名额 > 0；未过截止时间；未过开始时间 |
| open      | cancel             | cancelled   | 调用方为发布者                            |
| open      | offline\_by\_admin | offline     | 调用方为审核员                            |
| full      | cancel\_join       | open        | 调用方为参与者                            |
| full      | cancel             | cancelled   | 调用方为发布者                            |
| full      | offline\_by\_admin | offline     | 调用方为审核员                            |
| rejected  | cancel             | cancelled   | 调用方为发布者                            |
| rejected  | review\_reapprove  | reviewing   | 调用方为发布者                            |

#### Market（二手）状态机

| 当前状态      | 动作                 | 目标状态      | 守卫条件    |
| --------- | ------------------ | --------- | ------- |
| reviewing | review\_approve    | published | 调用方为审核员 |
| reviewing | review\_reject     | rejected  | 调用方为审核员 |
| reviewing | delete             | offline   | 调用方为发布者 |
| published | delete             | offline   | 调用方为发布者 |
| published | offline\_by\_admin | offline   | 调用方为审核员 |
| rejected  | delete             | offline   | 调用方为发布者 |
| rejected  | review\_reapprove  | reviewing | 调用方为发布者 |

#### LostFound（失物招领）状态机

| 当前状态      | 动作                 | 目标状态      | 守卫条件    |
| --------- | ------------------ | --------- | ------- |
| reviewing | review\_approve    | published | 调用方为审核员 |
| reviewing | review\_reject     | rejected  | 调用方为审核员 |
| reviewing | delete             | offline   | 调用方为发布者 |
| published | mark\_resolved     | resolved  | 调用方为发布者 |
| published | delete             | offline   | 调用方为发布者 |
| published | offline\_by\_admin | offline   | 调用方为审核员 |
| rejected  | delete             | offline   | 调用方为发布者 |
| rejected  | review\_reapprove  | reviewing | 调用方为发布者 |

#### Carpool（拼车）状态机

| 当前状态      | 动作                 | 目标状态      | 守卫条件    |
| --------- | ------------------ | --------- | ------- |
| reviewing | review\_approve    | published | 调用方为审核员 |
| reviewing | review\_reject     | rejected  | 调用方为审核员 |
| reviewing | delete             | offline   | 调用方为发布者 |
| published | delete             | offline   | 调用方为发布者 |
| published | offline\_by\_admin | offline   | 调用方为审核员 |
| rejected  | delete             | offline   | 调用方为发布者 |
| rejected  | review\_reapprove  | reviewing | 调用方为发布者 |

#### Resource（资料）状态机

| 当前状态      | 动作                 | 目标状态      | 守卫条件    |
| --------- | ------------------ | --------- | ------- |
| reviewing | review\_approve    | published | 调用方为审核员 |
| reviewing | review\_reject     | rejected  | 调用方为审核员 |
| reviewing | delete             | offline   | 调用方为发布者 |
| published | delete             | offline   | 调用方为发布者 |
| published | offline\_by\_admin | offline   | 调用方为审核员 |
| rejected  | delete             | offline   | 调用方为发布者 |
| rejected  | review\_reapprove  | reviewing | 调用方为发布者 |

### Requirement: BMFS 上下文与守卫

系统 SHALL 定义 `ActionContext` 结构体，包含守卫判断所需的全部上下文信息：

```go
type ActionContext struct {
    Principal       auth.Principal
    IsOwner         bool
    UserRole        string  // publisher / acceptor / participant / viewer
    Now             time.Time
    Extra           map[string]any  // 类型特有上下文（如 remaining_seats, deadline_at）
}
```

守卫函数签名：`func(ctx ActionContext) bool`

### Requirement: can\_xxx 派生规则

系统 SHALL 通过 BMFS 自动派生 can\_xxx 布尔值：

- 对每个状态机中定义的动作，BMFS 根据当前状态和 ActionContext 计算该动作是否可用
- 动作名 `cancel` → `can_cancel`，动作名 `accept` → `can_accept`，以此类推
- 派生结果以 `map[string]bool` 形式返回，service 层直接注入到响应 payload 中

### Requirement: service 层接入 BMFS

系统 SHALL 修改 campus\_life 各 service 文件，将所有手写状态判断替换为 BMFS 调用：

- 发布操作：创建内容时初始状态为 `reviewing`，无需 BMFS 执行动作
- 状态变更操作：调用 `bmfs.Execute()` 替代直接设置 `item.Status`
- can\_xxx 计算：调用 `bmfs.AvailableActions()` 替代 helpers.go 中的手写函数
- 删除 helpers.go 中的 `canEditContent`、`canDeleteContent` 等函数

### Requirement: transport 层动作命令

系统 SHALL 修改 transport 层，客户端提交动作命令而非目标状态：

- 接口路径从 `/api/errand/:id/cancel` 等语义化路径保持不变
- handler 内部将 HTTP 请求映射为 BMFS 动作名（如 `cancel`、`accept`、`join`）
- 审核接口 `ReviewUpdateRequest.ReviewStatus` 改为 `Action` 字段，值为 `review_approve` 或 `review_reject`

### Requirement: 状态转换日志

系统 SHALL 在每次状态转换成功后记录转换日志：

- 记录内容：content\_id、content\_type、from\_status、to\_status、action、actor\_user\_id、timestamp
- 日志写入 MongoDB `state_transition_logs` 集合
- 日志记录通过 BMFS 的 onTransition 钩子触发，不侵入 service 层

## MODIFIED Requirements

### Requirement: helpers.go 中的状态判断函数

`canEditContent`、`canDeleteContent`、`shouldExposeContent`、`shouldExposeMeetupState` 等函数 SHALL 改为从 BMFS 派生结果中读取，不再手写条件分支。`canViewContact`、`canModerateCampusLife` 等权限判断函数保留在 helpers.go 中（它们属于权限裁决而非状态流转）。

### Requirement: ReviewUpdateRequest

`ReviewUpdateRequest.ReviewStatus` 字段 SHALL 改为 `Action` 字段，值为 `review_approve` 或 `review_reject`，不再接受 `reviewing`/`published`/`rejected`/`offline` 等状态值。

## REMOVED Requirements

### Requirement: helpers.go 中散落的 can\_xxx 手写逻辑

**Reason**: BMFS 统一派生 can\_xxx，手写逻辑不再需要
**Migration**: 逐步替换，先保留旧函数标记 deprecated，BMFS 接入完成后再删除

### Requirement: service 中直接设置 item.Status 的代码

**Reason**: 状态变更必须经过 BMFS，不允许 service 层绕过状态机直接修改状态
**Migration**: 所有 `item.Status = xxx` 替换为 `bmfs.Execute()` 调用
