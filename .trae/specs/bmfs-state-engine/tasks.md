# Tasks

- [x] Task 1: 实现 BMFS 引擎核心包
  - [x] SubTask 1.1: 定义核心类型（State、Action、Transition、Machine、ActionContext、GuardFunc、OnTransitionFunc）
  - [x] SubTask 1.2: 实现 Machine 构建器（NewMachine、AddState、AddTransition）
  - [x] SubTask 1.3: 实现 Execute 方法（查找转换、执行 guard、更新状态、触发 onTransition）
  - [x] SubTask 1.4: 实现 AvailableActions 方法（遍历当前状态的转换、执行 guard、返回可用动作和 can_xxx 映射）
  - [x] SubTask 1.5: 编写 BMFS 核心单元测试（正常转换、guard 拒绝、无效动作、can_xxx 派生）

- [x] Task 2: 定义 6 类社区内容状态机
  - [x] SubTask 2.1: 在 campus_life/types 下新增 statemachines.go，定义 errandStateMachine
  - [x] SubTask 2.2: 定义 meetupStateMachine（含 join/cancel_join 的动态目标状态逻辑）
  - [x] SubTask 2.3: 定义 marketStateMachine
  - [x] SubTask 2.4: 定义 lostFoundStateMachine
  - [x] SubTask 2.5: 定义 carpoolStateMachine
  - [x] SubTask 2.6: 定义 resourceStateMachine
  - [x] SubTask 2.7: 提供 GetMachine(contentType string) *bmfs.Machine 工厂函数

- [x] Task 3: 实现状态转换日志
  - [x] SubTask 3.1: 在 repo 层新增 StateTransitionLog 模型和 WriteTransitionLog 方法
  - [x] SubTask 3.2: 在 BMFS onTransition 钩子中调用 repo 写入日志

- [x] Task 4: 改造 errand_service 接入 BMFS
  - [x] SubTask 4.1: AcceptErrand 改为调用 bmfs.Execute(action="accept")
  - [x] SubTask 4.2: CancelErrandPublish 改为调用 bmfs.Execute(action="cancel")
  - [x] SubTask 4.3: CancelErrandAccept 改为调用 bmfs.Execute(action="cancel_accept")
  - [x] SubTask 4.4: buildErrandPayload 中 can_xxx 改为从 bmfs.AvailableActions 派生

- [x] Task 5: 改造 meetup_service 接入 BMFS
  - [x] SubTask 5.1: JoinMeetup 改为调用 bmfs.Execute(action="join")
  - [x] SubTask 5.2: CancelMeetupJoin 改为调用 bmfs.Execute(action="cancel_join")
  - [x] SubTask 5.3: CancelMeetupPublish 改为调用 bmfs.Execute(action="cancel")
  - [x] SubTask 5.4: buildMeetupPayload 中 can_join/can_cancel_join/can_cancel_publish 改为从 BMFS 派生

- [x] Task 6: 改造 market/lostfound/carpool/resource service 接入 BMFS
  - [x] SubTask 6.1: market_service 的 DeleteMarket 改为 bmfs.Execute(action="delete")
  - [x] SubTask 6.2: lostfound_service 的 DeleteLostFound 改为 bmfs.Execute(action="delete")，MarkResolved 改为 bmfs.Execute(action="mark_resolved")
  - [x] SubTask 6.3: carpool_service 的 DeleteCarpool 改为 bmfs.Execute(action="delete")
  - [x] SubTask 6.4: resource_service 的 DeleteResource 改为 bmfs.Execute(action="delete")

- [x] Task 7: 改造 review_service 接入 BMFS
  - [x] SubTask 7.1: UpdateReviewStatus 改为接收 Action 字段（review_approve / review_reject）
  - [x] SubTask 7.2: 审核通过时通过 BMFS 执行 review_approve 动作
  - [x] SubTask 7.3: meetup 审核通过后的 open/full 判断通过 BMFS guard/onTransition 处理

- [x] Task 8: 清理 helpers.go 中的旧状态判断函数
  - [x] SubTask 8.1: canEditContent 保留（edit 不是状态转换，是 UI 功能）
  - [x] SubTask 8.2: 移除 canDeleteContent 函数，改用 BMFS 派生
  - [x] SubTask 8.3: 移除 isSupportedReviewStatus 函数（review_service 已改用 action 验证）

- [x] Task 9: 更新 transport 层动作命令
  - [x] SubTask 9.1: ReviewUpdateRequest.ReviewStatus 改为 Action 字段
  - [x] SubTask 9.2: review_service 内部将 action 映射为 BMFS 动作名

- [x] Task 10: 编写集成测试
  - [x] SubTask 10.1: 验证 6 类内容的完整状态流转路径
  - [x] SubTask 10.2: 验证 guard 条件拒绝非法操作
  - [x] SubTask 10.3: 验证 can_xxx 派生结果正确
  - [x] SubTask 10.4: 验证 GetMachine 工厂函数和 guard 函数

# Task Dependencies

- [Task 2] depends on [Task 1]
- [Task 3] depends on [Task 1]
- [Task 4] depends on [Task 1, Task 2]
- [Task 5] depends on [Task 1, Task 2]
- [Task 6] depends on [Task 1, Task 2]
- [Task 7] depends on [Task 1, Task 2]
- [Task 8] depends on [Task 4, Task 5, Task 6]
- [Task 9] depends on [Task 7]
- [Task 10] depends on [Task 4, Task 5, Task 6, Task 7, Task 8, Task 9]

# Parallelizable Work

- Task 4, Task 5, Task 6 可并行执行（各自独立改造不同 service）
- Task 3 可与 Task 2 并行执行
