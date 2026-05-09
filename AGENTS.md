# AGENTS.md

本仓库默认面向“人类开发者 + AI 编码代理”共同协作。

## 阅读顺序

1. [ARCHITECTURE.md](/Users/liangluo/code/weouc2026/ARCHITECTURE.md)
2. [docs/product-specs/index.md](/Users/liangluo/code/weouc2026/docs/product-specs/index.md)
3. [docs/PLANS.md](/Users/liangluo/code/weouc2026/docs/PLANS.md)
4. 进入具体目录前，先读对应 `README.md`

## 仓库导航

- [apps/admin-web/README.md](/Users/liangluo/code/weouc2026/apps/admin-web/README.md)：管理员后台
- [apps/miniapp-wechat/README.md](/Users/liangluo/code/weouc2026/apps/miniapp-wechat/README.md)：微信小程序
- [apps/mobile-flutter/README.md](/Users/liangluo/code/weouc2026/apps/mobile-flutter/README.md)：Flutter App
- [services/api-server/README.md](/Users/liangluo/code/weouc2026/services/api-server/README.md)：后端主服务
- [packages/contracts/README.md](/Users/liangluo/code/weouc2026/packages/contracts/README.md)：接口契约与模型
- [docs/design-docs/core-beliefs.md](/Users/liangluo/code/weouc2026/docs/design-docs/core-beliefs.md)：核心工程信条

## 不可违反的架构约束

- 业务规则只在后端落地；客户端只做展示、交互编排、轻量本地缓存。
- 所有对外接口先定义契约，再实现服务，再消费 SDK。
- 任意数据库读写只能经过 `repo` 层；`handler/page/viewmodel` 不允许直连数据库。
- 任意外部系统接入只能经过 `providers/connectors`；不允许散落在业务代码中。
- 小程序、Flutter、管理后台必须复用同一套业务语义，禁止复制状态机规则。
- 管理后台和客户端共享后端能力，但权限隔离必须在后端完成，不能依赖前端隐藏按钮。
- 重大变更必须同步更新：
  - [ARCHITECTURE.md](/Users/liangluo/code/weouc2026/ARCHITECTURE.md)
  - [docs/PLANS.md](/Users/liangluo/code/weouc2026/docs/PLANS.md)
  - 对应模块的 `README.md`

## 后端分层约束

后端按如下方向依赖：

```text
types -> config -> repo -> service -> runtime -> transport
```

说明：

- `types`：领域实体、DTO、枚举、错误码、事件模型
- `config`：配置模型和模块装配
- `repo`：数据库和缓存访问
- `service`：领域规则、用例、权限判断
- `runtime`：定时任务、异步作业、跨域编排、同步任务
- `transport`：HTTP、管理端接口、回调入口、webhook

允许跨层访问的只有 `providers`，用于统一封装：

- 微信开放能力
- 短信服务
- 对象存储
- 校园统一认证
- 教务/学工等第三方系统

## 客户端约束

### 管理后台

- 按业务域拆分模块，不按页面随意堆逻辑
- 使用服务端返回的权限码控制路由和按钮
- 表单、表格、字典项优先通过配置和 schema 复用

### 微信小程序

- 使用原生小程序，不引入 Taro、Uni-app 等跨端框架
- 优先组件化和分包，严控主包体积
- 数据更新使用局部 `setData`，避免整树刷新

### Flutter App

- 采用 `feature-first` 目录组织
- 通过 `ViewModel + Repository` 实现单向数据流
- 离线缓存与同步逻辑放在数据层，不放在页面组件

## 文档规则

- 文档是系统事实来源，不是事后补录
- 每个设计决策应能追溯到文档、任务或计划
- 计划中的状态变更要及时更新到 [docs/PLANS.md](/Users/liangluo/code/weouc2026/docs/PLANS.md)
- 架构变化较大时，新增 `ADR` 或在现有文档中记录“为什么变更”

## 验证规则

新增代码后至少补齐以下检查：

- 后端：单元测试、接口契约校验、结构化日志、错误码一致性
- 管理后台：类型检查、路由权限校验、核心页面交互验证
- 小程序：分包体积检查、登录链路验证、关键页面性能检查
- Flutter：`analyze`、核心流程测试、离线缓存与恢复验证

## 提交前自检

- 是否违反了层级依赖方向
- 是否把业务规则写到了客户端
- 是否新增了未文档化的技术决策
- 是否遗漏了权限、审计、可观测性
- 是否更新了相关计划和说明

