# AGENTS.md

本仓库面向"人类开发者 + AI 编码代理"共同协作。

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

---

## ⚠️ 生成代码前必做：约束激活检查

每次生成或修改代码前，逐条确认以下约束，并在输出中说明实现方式：

- [ ] 数据库访问 → 只经过 repo 层（见约束 1）
- [ ] 外部系统调用 → 只经过 providers（见约束 2）
- [ ] 业务规则 → 只在后端 service 层落地（见约束 3）
- [ ] 权限判断 → 只在后端完成（见约束 4）
- [ ] 接口变更 → 先改契约再改实现（见约束 5）
- [ ] 多端语义 → 复用后端返回的状态，不在客户端重写规则（见约束 6）
- [ ] 重大变更 → 同步更新 ARCHITECTURE.md / PLANS.md / 模块 README.md

---

## 架构约束

### 约束 1：数据库访问只在 repo 层

所有数据库读写（SQL、ORM、缓存）只放在 `repo` 层。`handler`、`page`、`viewmodel`、`service` 通过调用 repo 接口间接访问数据。

```text
✅ 正确：transport → service → repo → database
❌ 错误：transport → database（跳过 repo）
❌ 错误：service → database（跳过 repo）
```

```go
// ✅ 正确：service 调用 repo 接口
func (s *ErrandService) List(ctx context.Context, filter ErrandFilter) ([]Errand, error) {
    return s.repo.FindByFilter(ctx, filter)
}

// ❌ 错误：handler 直接操作数据库
func (h *ErrandHandler) List(c *gin.Context) {
    var items []Errand
    db.Where("status = ?", 1).Find(&items) // 禁止！
}
```

### 约束 2：外部系统调用只在 providers 层

微信 SDK、短信服务、对象存储、校园统一认证、教务/学工系统等所有外部调用，统一封装在 `providers` 层。业务模块通过 provider 接口调用，不在业务代码中直接拼 URL 或调 SDK。

```text
✅ 正确：service → wechat_provider.SendTemplateMsg()
❌ 错误：handler → wechat SDK 直接调用
❌ 错误：service → http.Post("https://api.weixin.qq.com/...")
```

```go
// ✅ 正确：service 调用 provider 接口
func (s *NotificationService) SendWechat(ctx context.Context, msg TemplateMsg) error {
    return s.wechatProvider.SendTemplateMsg(ctx, msg)
}

// ❌ 错误：service 直接调微信 SDK
func (s *NotificationService) SendWechat(ctx context.Context, msg TemplateMsg) error {
    client := wechat.NewClient(appID, secret) // 禁止！
    client.SendTemplate(...)
}
```

### 约束 3：业务规则只在后端落地

状态流转、校验规则、权限判断、数据裁剪等业务逻辑只写在后端 `service` 层。客户端只负责：展示、交互编排、轻量本地缓存。

```text
✅ 正确：后端返回 { status: "pending", canCancel: true }，客户端直接渲染
❌ 错误：客户端根据 status 本地计算 canCancel = status === "pending"
```

```javascript
// ✅ 正确：小程序/Flutter 直接使用后端返回的 canCancel
this.setData({ canCancel: res.can_cancel })

// ❌ 错误：客户端自行推断业务规则
const canCancel = this.data.status === 'pending' && this.data.creator_id === userId // 禁止！
```

### 约束 4：权限隔离只在后端完成

权限裁决在后端完成，前端只根据后端返回的权限码控制展示。联系方式等受限字段由后端在读取时裁剪，前端不自行放开。

```text
✅ 正确：后端返回 { permissions: ["campus_life:moderate"], contact: null }（未绑定教务时联系方式为 null）
❌ 错误：前端 v-if="isAdmin" 隐藏按钮（仅做展示优化，不做权限裁决）
❌ 错误：前端根据本地状态决定是否显示联系方式
```

### 约束 5：接口变更先改契约再改实现

所有对外接口遵循：先定义契约（OpenAPI / contracts）→ 再实现后端服务 → 客户端消费生成的 SDK。

```text
✅ 正确：修改 openapi.yaml → 重新生成 SDK → 更新后端实现 → 客户端升级 SDK
❌ 错误：先改后端接口 → 再补契约文档
❌ 错误：客户端手写接口调用，不走生成 SDK
```

### 约束 6：多端复用同一套业务语义

小程序、Flutter、管理后台必须复用后端返回的业务语义（状态枚举、权限码、错误码），不在客户端复制状态机规则。

```text
✅ 正确：三端共用后端返回的 status 字段和 can_xxx 布尔值
❌ 错误：小程序写一套状态判断，Flutter 写一套状态判断，管理后台又写一套
```

### 约束 7：重大变更同步更新文档

重大变更必须同步更新以下文件：
- [ARCHITECTURE.md](/Users/liangluo/code/weouc2026/ARCHITECTURE.md)
- [docs/PLANS.md](/Users/liangluo/code/weouc2026/docs/PLANS.md)
- 对应模块的 `README.md`

---

## 后端分层约束

依赖方向（只允许从左到右）：

```text
types → config → repo → service → runtime → transport
```

| 层 | 职责 | 允许 import | 禁止 import |
|---|---|---|---|
| `types` | 实体、DTO、枚举、错误码、事件 | 无外部依赖 | config, repo, service, runtime, transport |
| `config` | 配置模型、模块装配 | types | repo, service, runtime, transport |
| `repo` | 数据库访问、缓存读写 | types, config | service, runtime, transport |
| `service` | 领域规则、用例、权限判断 | types, config, repo, providers | runtime, transport |
| `runtime` | 定时任务、异步作业、跨域编排 | types, config, repo, service | transport |
| `transport` | HTTP handler、webhook、回调 | types, service | repo（必须通过 service） |

`providers` 是唯一允许跨层访问的模块，统一封装外部系统：

```text
providers/
├── wechat_provider/    # 微信开放能力
├── sms_provider/       # 短信服务
├── storage_provider/   # 对象存储
├── sso_provider/       # 校园统一认证
└── academic_provider/  # 教务/学工等第三方系统
```

新增模块的标准文件结构：

```text
internal/modules/<domain>/
├── types/       # 实体、DTO、枚举、错误码
├── config/      # 模块配置、装配入口
├── repo/        # 数据库访问接口和实现
├── service/     # 领域规则、用例编排
├── runtime/     # 定时任务、异步消费
└── transport/   # HTTP handler、webhook
```

---

## 客户端约束

### 管理后台（Vue 3 + TypeScript）

- 按业务域组织模块（`modules/iam`、`modules/portal`），不按页面堆目录
- 路由和按钮可见性由后端返回的权限码控制（`v-if="hasPermission('campus_life:moderate')"`）
- 表单、表格、字典项通过配置和 schema 复用，不逐页手写
- 接口调用使用 contracts 生成的 SDK，不手写 axios 请求

### 微信小程序（原生）

- 使用原生小程序，不引入 Taro、Uni-app 等跨端框架
- 主包只保留登录、首页、信息流、个人中心；跑腿/组局/二手等放入分包
- `api/` 只负责路径和传输；页面消费 `services/` 输出的业务语义
- 使用局部 `setData` 更新数据，不整树刷新

```javascript
// ✅ 正确：局部更新
this.setData({ 'list[0].status': 'done' })

// ❌ 错误：整树刷新
this.setData({ list: newList })
```

- 联系方式可见性完全由后端返回，前端只根据 `contact_visible` 渲染

### Flutter App（Riverpod + go_router）

- 采用 `feature-first` 目录组织（`features/auth`、`features/campus_life`）
- 通过 `ViewModel + Repository` 实现单向数据流
- 离线缓存与同步逻辑放在 `data` 层，不放在页面组件

```dart
// ✅ 正确：ViewModel 调用 Repository
class ErrandListViewModel extends ChangeNotifier {
  final ErrandRepository _repo;
  Future<void> load() async {
    state = await _repo.fetchErrands();
    notifyListeners();
  }
}

// ❌ 错误：页面组件直接调 API
class ErrandListPage extends StatelessWidget {
  Widget build(BuildContext context) {
    dio.get('/errands'); // 禁止！
  }
}
```

- 状态机规则复用后端语义，不在 Flutter 重新实现

---

## 文档规则

- 文档是系统事实来源，不是事后补录
- 每个设计决策应能追溯到文档、任务或计划
- 计划中的状态变更要及时更新到 [docs/PLANS.md](/Users/liangluo/code/weouc2026/docs/PLANS.md)
- 架构变化较大时，新增 `ADR` 或在现有文档中记录"为什么变更"

---

## 验证规则

新增代码后至少补齐以下检查：

- 后端：单元测试、接口契约校验、结构化日志、错误码一致性
- 管理后台：类型检查、路由权限校验、核心页面交互验证
- 小程序：分包体积检查、登录链路验证、关键页面性能检查
- Flutter：`analyze`、核心流程测试、离线缓存与恢复验证

---

## 提交规则

- 每次完成一轮可独立说明的改动后，立即提交一个 `commit`
- `commit message` 准确描述改动范围；同时涉及代码与文档时需在说明中体现
- 开始下一轮改动前，确认当前工作区已提交完毕

---

## 提交前自检清单

逐条确认，全部通过后再提交：

- [ ] 层级依赖方向正确（transport 不直接 import repo，service 不直接操作数据库）
- [ ] 业务规则只写在后端 service 层（客户端无状态机、无权限推断逻辑）
- [ ] 外部系统调用只经过 providers（无散落的 SDK 调用或 HTTP 拼接）
- [ ] 权限判断在后端完成（前端仅做展示优化，不依赖隐藏按钮做权限裁决）
- [ ] 未遗漏权限、审计、可观测性
- [ ] 技术决策已文档化
- [ ] 相关计划（PLANS.md / ARCHITECTURE.md / README.md）已同步更新
