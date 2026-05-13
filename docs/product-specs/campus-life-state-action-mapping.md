# PRD：校园生活模块状态→操作映射规范化

## 1. 执行摘要

**决策**：为跑腿、组局、二手交易、资料、失物招领、拼车六个校园生活模块建立统一的"状态→操作"映射规范，消除当前前端本地推断业务规则的问题，确保所有操作权限由后端 `can_xxx` 字段驱动，客户端零推断。

**用户**：小程序前端开发者、Flutter 开发者、后端开发者

**范围**：六个校园生活模块的详情页和列表页的状态展示与操作按钮渲染逻辑

**成功标准**：前端代码中不存在任何基于 `status` 本地推断 `can_xxx` 的逻辑；所有模块的状态枚举值与展示文案有统一映射；后端 `resolved` 状态不再被归一化丢失

---

## 2. 问题陈述

### 2.1 核心问题

当前小程序前端存在"状态与操作不匹配"的系统性问题，具体表现为：

1. **客户端违反约束 3**：errandService.js 列表页本地计算 `canAccept`（`!isPublisher && !isAccepted && status === 'published'`），业务规则未收敛到后端
2. **状态归一化丢失信息**：后端 `normalizeReviewStatus` 将 LostFound 的 `resolved` 状态归一化为 `published`，导致前端无法区分"已发布"和"已找到"
3. **`canCancelPublish` 与 `canDelete` 语义混淆**：errand 详情页 `canCancelPublish` 复用 `canDelete`，meetup 同时返回 `can_delete` 和 `can_cancel_publish`（后者标注"兼容旧客户端"），语义不统一
4. **列表页与详情页状态展示不一致**：LostFound 列表页硬编码中文状态（`'寻找中'`/`'待认领'`），详情页使用后端枚举值
5. **meetup `getActionText` 包含本地状态判断**：根据 `status` 本地推断操作文案，而非消费后端返回的语义
6. **Carpool/Resource 缺少详情页**：can_xxx 字段虽已由后端返回，但前端无页面消费

### 2.2 影响范围

| 影响面 | 说明 |
|--------|------|
| 用户体验 | 同一内容在不同页面显示不一致的状态标签和操作入口 |
| 开发效率 | 前端开发者需要理解后端状态合并逻辑才能正确渲染 |
| 多端一致性 | 小程序和 Flutter 如果各自推断，规则必然漂移 |
| 可维护性 | 状态枚举散落在多个文件，修改时容易遗漏 |

---

## 3. 现状分析

### 3.1 当前状态枚举全景

#### 通用审核状态（所有类型共用）

| 枚举值 | 含义 | 当前前端展示 |
|--------|------|-------------|
| `reviewing` | 审核中 | ✅ 详情页标题旁 Tag |
| `published` | 已发布/审核通过 | ✅ 默认态，无 Tag |
| `rejected` | 审核未通过 | ✅ 详情页标题旁 Tag |
| `offline` | 已下架 | ✅ 详情页标题旁 Tag |

#### Errand 专有业务状态

| 枚举值 | 含义 | 当前前端展示 |
|--------|------|-------------|
| `accepted` | 已接单 | ❌ 无状态 Tag |
| `cancelled` | 已取消 | ✅ 详情页标题旁 Tag |

#### Meetup 专有业务状态

| 枚举值 | 含义 | 当前前端展示 |
|--------|------|-------------|
| `open` | 报名中 | ✅ meetupService STATUS_META |
| `full` | 人数已满 | ✅ meetupService STATUS_META |
| `cancelled` | 已取消 | ✅ 详情页标题旁 Tag |

#### LostFound 专有审核状态

| 枚举值 | 含义 | 当前前端展示 |
|--------|------|-------------|
| `resolved` | 已找到/已认领 | ⚠️ 详情页标题旁 Tag，但 `normalizeReviewStatus` 将其归一化为 `published` |

### 3.2 当前状态→操作映射问题矩阵

| 模块 | 状态 | 应有操作 | 当前操作 | 问题 |
|------|------|---------|---------|------|
| Errand | `accepted` | 接单者：取消接单；发布者：查看进度 | 接单者：取消接单 ✅；发布者：无操作 ❌ | 发布者在已接单状态下无任何操作入口 |
| Errand | `cancelled` | 无操作（终态） | 无操作 ✅ | - |
| Errand 列表 | `published` | 客态：接单 | 客态：本地推断 canAccept ❌ | 违反约束 3 |
| Meetup | `full` | 发布者：取消组局；参与者：取消报名 | 发布者：取消组局 ✅；参与者：取消报名 ✅ | - |
| Meetup 列表 | `open` | 客态：报名 | 本地推断 actionText ⚠️ | getActionText 包含本地状态判断 |
| Market | `published` | 发布者：编辑/下架；客态：收藏/查看联系方式 | 发布者：下架 ✅/编辑 ❌；客态：收藏 ✅/联系方式 ✅ | 编辑按钮未渲染 |
| LostFound | `resolved` | 发布者：下架；客态：无操作 | 发布者：下架 ✅；客态：无操作 ✅ | 但后端归一化丢失 resolved 状态 |
| LostFound 列表 | `published` | 显示"寻找中"/"待认领" | 硬编码中文 ❌ | 应使用后端 status + type 映射 |
| Carpool | `published` | 发布者：编辑/取消；客态：加入/查看联系方式 | 无详情页 ❌ | can_xxx 未被消费 |
| Resource | `published` | 发布者：编辑/下架；客态：下载/查看联系方式 | 无详情页 ❌ | can_xxx 未被消费 |

---

## 4. 解决方案

### 4.1 设计原则

1. **后端是唯一真相源**：所有 `can_xxx` 由后端 service 层计算，客户端零推断
2. **状态枚举统一对外**：后端合并后的 `status` 字段是客户端唯一消费的状态值
3. **展示映射集中管理**：前端维护一个统一的 `status → 展示文案 + tone` 映射表，不做业务判断
4. **操作权限只看 can_xxx**：前端渲染操作按钮的唯一依据是后端返回的 `can_xxx` 布尔字段

### 4.2 后端改造

#### 4.2.1 修复 `normalizeReviewStatus` 识别 `resolved`

**当前问题**：`normalizeReviewStatus` 不识别 `resolved`，将其归一化为 `published`，导致 LostFound 的"已找到"状态信息丢失。

**改造方案**：在 `normalizeReviewStatus` 中新增 `resolved` 分支：

```go
func normalizeReviewStatus(status string) string {
    switch strings.TrimSpace(strings.ToLower(status)) {
    case statusReviewing:  return "reviewing"
    case statusRejected:   return "rejected"
    case statusOffline:    return "offline"
    case statusResolved:   return "resolved"    // 新增
    default:               return "published"
    }
}
```

需同步新增常量：

```go
const statusResolved = "resolved"
```

**影响范围**：
- `mergeErrandStatus`：不受影响（errand 无 resolved 状态）
- `mergeMeetupStatus`：不受影响（meetup 无 resolved 状态）
- LostFound 列表/详情：`status` 字段将正确返回 `"resolved"` 而非 `"published"`
- `canEditContent`：需确认 `resolved` 状态下 `can_edit` 应为 false（已找到的失物不应再编辑）
- `canDeleteContent`：需确认 `resolved` 状态下 `can_delete` 应为 true（发布者仍可下架已找到的失物）
- `shouldExposeContent`：`resolved` 应视为可公开状态（同 `published`）

#### 4.2.2 统一 `can_cancel_publish` 语义

**当前问题**：
- Errand 详情页前端将 `canCancelPublish` 映射为 `canDelete`，语义混淆
- Meetup 同时返回 `can_delete` 和 `can_cancel_publish`（后者标注"兼容旧客户端"）
- OpenAPI 契约中 Meetup 的 `can_cancel_publish` 描述为"兼容旧客户端，等同于 can_delete"

**改造方案**：

统一为 `can_delete`，废弃 `can_cancel_publish`：

| 模块 | 当前 | 改造后 |
|------|------|--------|
| Errand | `can_delete`（前端映射为 `canCancelPublish`） | `can_delete`（前端直接使用，按钮文案"取消发布"） |
| Meetup | `can_delete` + `can_cancel_publish` | 仅 `can_delete`（移除 `can_cancel_publish`） |
| Market | `can_delete` | `can_delete`（按钮文案"下架"） |
| Resource | `can_delete` | `can_delete`（按钮文案"下架"） |
| LostFound | `can_delete` | `can_delete`（按钮文案"下架"） |
| Carpool | `can_delete` | `can_delete`（按钮文案"取消发布"） |

**按钮文案由前端根据模块类型决定**，不由 can_xxx 字段名决定。后端只管"能不能"，前端只管"怎么展示"。

#### 4.2.3 Errand 列表接口补齐 `can_accept`

**当前问题**：errandService.js 列表页本地计算 `canAccept`。

**改造方案**：后端 Errand 列表接口的每个列表项中返回 `can_accept` 字段（与详情接口一致），前端列表页直接消费。

#### 4.2.4 各类型列表接口统一返回 `is_owner` 和关键 `can_xxx`

**当前状态**：0003 计划已完成此项，但需验证列表项中的 `can_xxx` 是否完整。

**验证清单**：

| 模块 | 列表项应返回的 can_xxx | 当前状态 |
|------|----------------------|---------|
| Market | `is_owner`, `can_favorite` | 需验证 |
| Errand | `is_owner`, `can_accept` | 需补齐 `can_accept` |
| Resource | `is_owner`, `can_download` | 需验证 |
| LostFound | `is_owner`, `can_mark_resolved` | 需验证 |
| Carpool | `is_owner`, `can_join_carpool` | 需验证 |
| Meetup | `is_owner`, `can_join` | 需验证 |

### 4.3 前端改造

#### 4.3.1 建立统一状态展示映射表

在 `services/shared.js` 中新增统一的状态展示映射函数：

```javascript
const STATUS_DISPLAY_MAP = {
  reviewing:  { label: '审核中',     tone: 'amber' },
  published:  { label: '已发布',     tone: 'green' },
  rejected:   { label: '审核未通过', tone: 'red'   },
  offline:    { label: '已下架',     tone: 'red'   },
  cancelled:  { label: '已取消',     tone: 'red'   },
  accepted:   { label: '已接单',     tone: 'blue'  },
  open:       { label: '报名中',     tone: 'green' },
  full:       { label: '人数已满',   tone: 'purple'},
  resolved:   { label: '已找到',     tone: 'green' },
};

export function getStatusDisplay(status, overrides = {}) {
  const entry = STATUS_DISPLAY_MAP[status];
  if (!entry) return { label: status, tone: 'default' };
  return {
    label: overrides[status]?.label || entry.label,
    tone:  overrides[status]?.tone  || entry.tone,
  };
}
```

**模块级覆盖**：

```javascript
// LostFound 根据类型覆盖
const lostFoundOverrides = {
  published: { label: '寻找中' },  // type=lost
  resolved:  { label: '已找到' },  // type=lost
};
// 或
const lostFoundOverrides = {
  published: { label: '待认领' },  // type=found
  resolved:  { label: '已认领' },  // type=found
};

// Errand 覆盖
const errandOverrides = {
  published: { label: '待接单' },
};
```

#### 4.3.2 移除所有本地 can_xxx 推断

| 文件 | 当前问题 | 改造方案 |
|------|---------|---------|
| `services/errandService.js:60` | `canAccept: !isPublisher && !isAccepted && status === 'published'` | 改为 `canAccept: Boolean(item.can_accept)` |
| `services/meetupService.js:28-36` | `getActionText` 根据 status 推断文案 | 改为根据 `can_join`/`can_cancel_join`/`can_delete` 决定文案 |
| `pages/errand/detail/index.js:92` | `canCancelPublish: Boolean(canDelete)` | 移除 `canCancelPublish`，直接使用 `canDelete` |
| `services/lostFoundService.js:46-47` | `status: isLost ? '寻找中' : '待认领'` | 改为使用 `getStatusDisplay(item.status, overrides)` |

#### 4.3.3 详情页操作按钮渲染规范化

**统一渲染规则**：

```
操作按钮渲染 = f(can_xxx)         // 只看 can_xxx，不看 status
标题旁状态 Tag = f(status)        // 只看 status，不做业务判断
```

**各模块详情页底部操作栏规范**：

##### Errand（跑腿）

```
┌──────────────────────────────────────────────────┐
│ 发布者 (is_owner)                                │
│   can_delete  → "取消发布" 按钮                   │
│   can_edit    → "编辑" 按钮                       │
│                                                  │
│ 接单者 (user_role=acceptor)                      │
│   can_cancel_accept → "取消接单" 按钮             │
│                                                  │
│ 浏览者 (user_role=viewer)                        │
│   can_accept  → "接单" 按钮（主操作）             │
│   !can_accept → "已锁定" 禁用态                   │
│                                                  │
│ 通用                                             │
│   can_view_contact → 联系方式区域                 │
└──────────────────────────────────────────────────┘
```

##### Meetup（组局）

```
┌──────────────────────────────────────────────────┐
│ 发布者 (is_owner)                                │
│   can_delete  → "取消组局" 按钮                   │
│   can_edit    → "编辑" 按钮                       │
│                                                  │
│ 参与者 (user_role=participant)                   │
│   can_cancel_join → "取消报名" 按钮               │
│                                                  │
│ 浏览者 (user_role=viewer)                        │
│   can_join    → "报名" 按钮（主操作）             │
│   !can_join   → "报名已截止"/"人数已满" 禁用态    │
│                                                  │
│ 通用                                             │
│   can_view_contact → 联系方式区域                 │
└──────────────────────────────────────────────────┘
```

##### Market（二手）

```
┌──────────────────────────────────────────────────┐
│ 发布者 (is_owner)                                │
│   can_delete  → "下架" 按钮                      │
│   can_edit    → "编辑" 按钮                       │
│                                                  │
│ 浏览者 (!is_owner)                               │
│   can_favorite → "想要" 按钮                      │
│                                                  │
│ 通用                                             │
│   can_view_contact → 联系方式区域                 │
└──────────────────────────────────────────────────┘
```

##### LostFound（失物招领）

```
┌──────────────────────────────────────────────────┐
│ 发布者 (is_owner)                                │
│   can_mark_resolved → "已找到"/"已认领" 按钮      │
│   can_delete        → "下架" 按钮                 │
│   can_edit          → "编辑" 按钮                 │
│                                                  │
│ 浏览者 (!is_owner)                               │
│   can_view_contact → 联系方式区域                 │
└──────────────────────────────────────────────────┘
```

##### Carpool（拼车）

```
┌──────────────────────────────────────────────────┐
│ 发布者 (is_owner)                                │
│   can_delete  → "取消发布" 按钮                   │
│   can_edit    → "编辑" 按钮                       │
│                                                  │
│ 浏览者 (!is_owner)                               │
│   can_join_carpool → "加入拼车" 按钮              │
│                                                  │
│ 通用                                             │
│   can_view_contact → 联系方式区域                 │
└──────────────────────────────────────────────────┘
```

##### Resource（资料）

```
┌──────────────────────────────────────────────────┐
│ 发布者 (is_owner)                                │
│   can_delete  → "下架" 按钮                      │
│   can_edit    → "编辑" 按钮                       │
│                                                  │
│ 浏览者 (!is_owner)                               │
│   can_download → "下载" 按钮                      │
│                                                  │
│ 通用                                             │
│   can_view_contact → 联系方式区域                 │
└──────────────────────────────────────────────────┘
```

---

## 5. 完整状态→操作映射表

### 5.1 Errand（跑腿）

| 合并后 status | 发布者操作 | 接单者操作 | 浏览者操作 | 标题旁 Tag |
|--------------|-----------|-----------|-----------|---------|
| `reviewing` | can_edit=true, can_delete=true | - | - | "审核中" |
| `published` | can_edit=true, can_delete=true | - | can_accept=true | 无 |
| `accepted` | can_edit=false, can_delete=true | can_cancel_accept=true | - | "已接单" |
| `rejected` | can_delete=true | - | - | "审核未通过" |
| `offline` | - | - | - | "已下架" |
| `cancelled` | - | - | - | "已取消" |

### 5.2 Meetup（组局）

| 合并后 status | 发布者操作 | 参与者操作 | 浏览者操作 | 标题旁 Tag |
|--------------|-----------|-----------|-----------|---------|
| `reviewing` | can_edit=true, can_delete=true | - | - | "审核中" |
| `open` | can_edit=true, can_delete=true | can_cancel_join=true | can_join=true | 无 |
| `full` | can_edit=true, can_delete=true | can_cancel_join=true | can_join=false | "人数已满" |
| `rejected` | can_delete=true | - | - | "审核未通过" |
| `offline` | - | - | - | "已下架" |
| `cancelled` | - | - | - | "已取消" |

### 5.3 Market（二手）

| 合并后 status | 发布者操作 | 浏览者操作 | 标题旁 Tag |
|--------------|-----------|-----------|---------|
| `reviewing` | can_edit=true, can_delete=true | - | "审核中" |
| `published` | can_edit=true, can_delete=true | can_favorite=true, can_view_contact | 无 |
| `rejected` | can_delete=true | - | "审核未通过" |
| `offline` | - | - | "已下架" |

### 5.4 Resource（资料）

| 合并后 status | 发布者操作 | 浏览者操作 | 标题旁 Tag |
|--------------|-----------|-----------|---------|
| `reviewing` | can_edit=true, can_delete=true | - | "审核中" |
| `published` | can_edit=true, can_delete=true | can_download=true, can_view_contact | 无 |
| `rejected` | can_delete=true | - | "审核未通过" |
| `offline` | - | - | "已下架" |

### 5.5 LostFound（失物招领）

| 合并后 status | 发布者操作 | 浏览者操作 | 标题旁 Tag |
|--------------|-----------|-----------|---------|
| `reviewing` | can_edit=true, can_delete=true | - | "审核中" |
| `published` | can_edit=true, can_delete=true, can_mark_resolved=true | can_view_contact | 无（列表显示"寻找中"/"待认领"） |
| `resolved` | can_delete=true | - | "已找到"/"已认领" |
| `rejected` | can_delete=true | - | "审核未通过" |
| `offline` | - | - | "已下架" |

### 5.6 Carpool（拼车）

| 合并后 status | 发布者操作 | 浏览者操作 | 标题旁 Tag |
|--------------|-----------|-----------|---------|
| `reviewing` | can_edit=true, can_delete=true | - | "审核中" |
| `published` | can_edit=true, can_delete=true | can_join_carpool=true, can_view_contact | 无 |
| `rejected` | can_delete=true | - | "审核未通过" |
| `offline` | - | - | "已下架" |

---

## 6. 实施工作流

### 工作流 A：后端状态归一化修复

1. 在 `service.go` 新增 `statusResolved` 常量
2. 修改 `normalizeReviewStatus` 识别 `resolved`
3. 修改 `shouldExposeContent` 将 `resolved` 视为可公开状态
4. 修改 `canEditContent` 在 `resolved` 状态下返回 false
5. 修改 `canDeleteContent` 在 `resolved` 状态下返回 true（发布者仍可下架）
6. Errand 列表接口补齐 `can_accept` 字段
7. 移除 Meetup 的 `can_cancel_publish` 字段，统一使用 `can_delete`
8. 验证各类型列表接口的 `can_xxx` 返回完整性
9. `go test ./...` 全量通过

### 工作流 B：OpenAPI 契约同步

1. 在 `CampusLifeOwnerContext` 的 `status` 枚举中新增 `resolved` 和 `accepted`
2. 移除 Meetup 的 `can_cancel_publish` 字段
3. Errand 列表响应 schema 新增 `can_accept`
4. 重新生成 JS/Dart SDK

### 工作流 C：前端统一改造

1. 在 `services/shared.js` 新增 `STATUS_DISPLAY_MAP` 和 `getStatusDisplay` 函数
2. 修改 `errandService.js`：移除本地 `canAccept` 推断，改为消费 `item.can_accept`
3. 修改 `meetupService.js`：移除 `getActionText` 中的本地状态判断，改为根据 `can_xxx` 决定文案
4. 修改 `lostFoundService.js`：列表页 status 使用 `getStatusDisplay` 替代硬编码中文
5. 修改 `errand/detail/index.js`：移除 `canCancelPublish`，直接使用 `canDelete`
6. 修改 `meetup/detail/index.js`：移除 `canCancelPublish`，直接使用 `canDelete`
7. Errand 详情页新增 `accepted` 状态 Tag
8. 将所有详情页的全宽状态横幅替换为标题旁内联 Tag（移除 `detail-status-banner`/`lf-status-banner`/`meetup-status-banner` 样式，新增统一 `status-tag` 样式）
9. 各详情页确认编辑按钮在 `can_edit=true` 时渲染（待编辑接口 A6 完成后启用）
10. `npm run check:syntax` 通过

### 工作流 D：验收与回归

1. 逐模块验证标题旁状态 Tag 与操作按钮的对应关系
2. 验证 LostFound `resolved` 状态在列表和详情页正确展示
3. 验证 Errand 列表页 `can_accept` 来自后端
4. 验证 Meetup 不再返回 `can_cancel_publish`
5. 验证"我的发布"页面状态展示与详情页一致
6. `make check` 全量通过

---

## 7. 验收标准

### 后端验收

- [ ] `normalizeReviewStatus` 识别 `resolved`，LostFound 详情/列表返回 `status: "resolved"`
- [ ] Errand 列表接口每个列表项返回 `can_accept` 布尔字段
- [ ] Meetup 详情/列表不再返回 `can_cancel_publish`
- [ ] `resolved` 状态下 `can_edit=false`、`can_delete=true`
- [ ] `go test ./...` 全量通过

### 契约验收

- [ ] OpenAPI `status` 枚举包含 `resolved`、`accepted`、`open`、`full`、`cancelled`
- [ ] Meetup schema 不包含 `can_cancel_publish`
- [ ] JS/Dart SDK 可正常生成

### 前端验收

- [ ] `errandService.js` 无本地 `canAccept` 推断，消费 `item.can_accept`
- [ ] `meetupService.js` 的 `getActionText` 无本地 status 判断，基于 `can_xxx` 决定文案
- [ ] `lostFoundService.js` 列表页 status 使用 `getStatusDisplay`，无硬编码中文
- [ ] Errand/Meetup 详情页无 `canCancelPublish`，统一使用 `canDelete`
- [ ] Errand 详情页有 `accepted` 状态 Tag
- [ ] 所有详情页操作按钮仅根据 `can_xxx` 渲染，无基于 `status` 的条件判断
- [ ] 所有详情页状态展示使用标题旁内联 Tag，无全宽横幅
- [ ] `npm run check:syntax` 通过

### 端到端验收

- [ ] LostFound 标记已找到后，列表页显示"已找到"/"已认领"标签
- [ ] LostFound 标记已找到后，详情页标题旁显示"已找到"/"已认领" Tag，且发布者仍可下架
- [ ] Errand 被接单后，详情页标题旁显示"已接单" Tag
- [ ] Meetup 不再出现 `can_cancel_publish` 字段

---

## 8. 风险与应对

| 风险 | 影响 | 应对 |
|------|------|------|
| `normalizeReviewStatus` 修改影响已有逻辑 | 现有 Market/Resource 等模块的 `can_edit`/`can_delete` 计算可能受影响 | `resolved` 仅 LostFound 使用，其他模块不会进入此分支；逐模块单测覆盖 |
| 移除 `can_cancel_publish` 可能影响旧版小程序 | 旧版客户端仍依赖此字段 | 保留一个过渡期（2周），在过渡期内同时返回 `can_delete` 和 `can_cancel_publish`，之后移除 |
| Errand 列表接口新增 `can_accept` 增加响应体积 | 列表接口响应变大 | `can_accept` 为布尔值，体积影响可忽略 |
| `getStatusDisplay` 集中管理可能导致模块间文案冲突 | 不同模块对同一状态值需要不同展示 | 通过 `overrides` 参数支持模块级覆盖 |

---

## 9. 不纳入本次范围

- Carpool 详情页创建（需独立迭代）
- Resource 详情页创建（需独立迭代）
- 编辑接口实现（A6，需独立迭代）
- Carpool 加入拼车接口实现（A8，需独立迭代）
- Flutter App 状态消费改造（需独立迭代）
- 管理后台适配（当前已满足，无需改动）
