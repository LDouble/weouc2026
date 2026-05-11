# 校园生活统一状态设计 - 实现计划

## [x] Task 1: 定义统一状态枚举常量
- **Priority**: P0
- **Depends On**: None
- **Description**: 
  - 在 service.go 中定义统一的状态枚举常量
  - 状态值：reviewing, published, rejected, offline, cancelled, full, accepted, open
- **Acceptance Criteria Addressed**: [AC-1]
- **Test Requirements**:
  - `programmatic` TR-1.1: 状态常量定义正确且可访问
- **Notes**: 需要与现有状态值兼容

## [x] Task 2: 修改服务层状态判断逻辑
- **Priority**: P0
- **Depends On**: Task 1
- **Description**: 
  - 修改 service.go 中的状态判断逻辑
  - 将 review_status 和 status 的双重判断合并为单一状态判断
  - 添加 mergeMeetupStatus 和 mergeErrandStatus 函数
- **Acceptance Criteria Addressed**: [AC-2]
- **Test Requirements**:
  - `programmatic` TR-2.1: 状态判断逻辑正确，审核中内容对普通用户不可见
  - `programmatic` TR-2.2: 已发布内容对所有用户可见
- **Notes**: 需要确保权限控制逻辑不受影响

## [x] Task 3: 修改 meetup 模块状态处理
- **Priority**: P0
- **Depends On**: Task 2
- **Description**: 
  - 修改 buildMeetupPayload 函数，返回统一的 status 字段
  - 更新 JoinMeetup 函数使用统一状态判断
- **Acceptance Criteria Addressed**: [AC-2, AC-4]
- **Test Requirements**:
  - `programmatic` TR-3.1: meetup 列表接口返回统一的 status 字段
  - `programmatic` TR-3.2: meetup 详情接口返回统一的 status 字段
- **Notes**: 需要检查所有使用 review_status 的地方

## [x] Task 4: 修改其他模块状态处理（market/errand/resource/lostFound/carpool）
- **Priority**: P1
- **Depends On**: Task 2
- **Description**: 
  - 对 market、errand、resource、lostFound、carpool 模块应用相同的状态统一逻辑
  - 更新各模块的 List/Get 接口返回单一 status 字段
- **Acceptance Criteria Addressed**: [AC-2, AC-4]
- **Test Requirements**:
  - `programmatic` TR-4.1: 各模块列表接口返回统一的 status 字段
- **Notes**: 需要逐一检查每个模块的实现

## [x] Task 5: 更新前端小程序状态处理
- **Priority**: P0
- **Depends On**: Task 3
- **Description**: 
  - 修改 meetupService.js 中的 resolveStatusMeta 函数
  - 移除对 review_status 的独立判断
  - 更新 errandService.js 状态处理逻辑
- **Acceptance Criteria Addressed**: [AC-3]
- **Test Requirements**:
  - `human-judgment` TR-5.1: 状态标签正确显示
  - `human-judgment` TR-5.2: 状态颜色符合预期
- **Notes**: 需要同步修改相关页面组件

## [/] Task 6: 更新测试用例
- **Priority**: P1
- **Depends On**: Task 3, Task 4
- **Description**: 
  - 更新现有测试用例以适应新的状态模型
  - 添加状态流转测试
- **Acceptance Criteria Addressed**: [AC-2, AC-4]
- **Test Requirements**:
  - `programmatic` TR-6.1: 所有现有测试通过
  - `programmatic` TR-6.2: 状态流转测试覆盖主要场景
- **Notes**: 需要检查 miniapp_api_test.go 等测试文件