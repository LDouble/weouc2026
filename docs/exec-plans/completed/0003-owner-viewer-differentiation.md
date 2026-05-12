# 0003 校园生活主态/客态操作行为差异化

> **状态：已归档**（2026-05-13）— 核心工作流 A/B/C/D 已完成，剩余 A6（编辑接口）、A8（Carpool 加入拼车）、C6（Carpool 详情页）作为后续迭代项。

## 目标

为拼车、闲置（二手）、组局、失物招领、跑腿五个校园生活模块统一落地"主态/客态"操作行为差异化：后端在详情接口中返回完整的 `user_role` + `can_xxx` 布尔字段，客户端直接消费不做本地推断，确保发布者（主态）与浏览者（客态）看到不同的操作入口。

## 背景

### 当前问题

1. **主态/客态区分不统一**：仅 Meetup 和 Errand 返回 `user_role`，Market、Resource、LostFound、Carpool 完全没有身份区分。
2. **can_xxx 布尔字段缺失**：除 Meetup 外，其他类型均未返回 `can_edit`、`can_delete`、`can_cancel_publish` 等操作权限字段。
3. **客户端违反约束 3**：Errand 小程序详情页在本地推断 `canAccept`、`canCancelPublish`、`canCancelAccept`，业务规则未收敛到后端。
4. **联系方式裁剪不一致**：LostFound 详情接口遗漏 `can_view_contact`；Resource 详情接口完全没有联系方式裁剪。
5. **发布者无法管理自己的内容**：Market、Resource、LostFound、Carpool 的发布者无法下架/删除自己发布的内容。
6. **OpenAPI 契约严重滞后**：所有详情接口的响应结构未定义，关键字段在契约层面不可见。

### 参考实现

Meetup 的 `buildMeetupPayload` 是当前最符合架构约束的实现，返回完整的 `user_role` + `can_join` + `can_cancel_join` + `can_cancel_publish`，客户端直接消费不推断。本次改造以此为标杆。

## 主态/客态行为矩阵

### 通用字段（所有类型统一）

| 字段 | 类型 | 说明 |
|------|------|------|
| `user_role` | string | `publisher` / `participant`(或`acceptor`) / `viewer` |
| `is_owner` | bool | 是否为发布者（`user_role == "publisher"` 的便捷布尔值） |
| `can_view_contact` | bool | 是否可查看联系方式（已绑定教务 或 是发布者本人） |
| `can_edit` | bool | 是否可编辑（仅发布者 且 状态允许编辑时） |
| `can_delete` | bool | 是否可删除/下架（仅发布者 且 状态允许下架时） |

### 各类型特有字段

| 类型 | 主态（publisher）操作 | 客态（viewer/acceptor/participant）操作 | 特有 can_xxx |
|------|---------------------|--------------------------------------|-------------|
| **Market（闲置/二手）** | 编辑、下架 | 查看联系方式、收藏 | `can_favorite` |
| **Errand（跑腿）** | 编辑、取消发布 | 查看联系方式、接单 | `can_accept`, `can_cancel_accept` |
| **Resource（资料）** | 编辑、下架 | 查看联系方式、下载 | `can_download` |
| **LostFound（失物招领）** | 编辑、下架、标记已找到/已认领 | 查看联系方式 | `can_mark_resolved` |
| **Carpool（拼车）** | 编辑、取消发布 | 查看联系方式 | `can_join_carpool` |
| **Meetup（组局）** | 编辑、取消发布 | 查看联系方式、报名、取消报名 | `can_join`, `can_cancel_join` |

### 各类型主态/客态操作详细规则

#### Market（闲置/二手）

| 操作 | 条件 | 说明 |
|------|------|------|
| `can_edit` | `is_owner && status in [published, reviewing]` | 发布者可编辑已发布或审核中的商品 |
| `can_delete` | `is_owner && status in [published, reviewing, rejected]` | 发布者可下架自己的商品 |
| `can_favorite` | `!is_owner && authenticated` | 非发布者可收藏 |

#### Errand（跑腿）

| 操作 | 条件 | 说明 |
|------|------|------|
| `can_edit` | `is_owner && status in [published, reviewing]` | 发布者可编辑 |
| `can_delete` | `is_owner && status in [published, reviewing, rejected]` | 发布者可取消发布 |
| `can_accept` | `!is_owner && status == published && authenticated` | 客态可接单 |
| `can_cancel_accept` | `user_role == acceptor && status == accepted` | 接单者可取消接单 |

#### Resource（资料）

| 操作 | 条件 | 说明 |
|------|------|------|
| `can_edit` | `is_owner && status in [published, reviewing]` | 发布者可编辑 |
| `can_delete` | `is_owner && status in [published, reviewing, rejected]` | 发布者可下架 |
| `can_download` | `authenticated && can_view_contact` | 绑定教务后可下载 |

#### LostFound（失物招领）

| 操作 | 条件 | 说明 |
|------|------|------|
| `can_edit` | `is_owner && status in [published, reviewing]` | 发布者可编辑 |
| `can_delete` | `is_owner && status in [published, reviewing, rejected]` | 发布者可下架 |
| `can_mark_resolved` | `is_owner && status == published` | 发布者可标记已找到/已认领 |

#### Carpool（拼车）

| 操作 | 条件 | 说明 |
|------|------|------|
| `can_edit` | `is_owner && status in [published, reviewing] && travel_at > now` | 出发前可编辑 |
| `can_delete` | `is_owner && status in [published, reviewing, rejected]` | 发布者可取消发布 |
| `can_join_carpool` | `!is_owner && status == published && authenticated` | 客态可加入拼车 |

#### Meetup（组局）

| 操作 | 条件 | 说明 |
|------|------|------|
| `can_edit` | `is_owner && status in [open, reviewing]` | 发布者可编辑 |
| `can_delete` | `is_owner && status != cancelled` | 发布者可取消发布（已有 `can_cancel_publish`，对齐为 `can_delete`） |
| `can_join` | `!is_owner && status == open && remaining_seats > 0 && before_deadline` | 客态可报名 |
| `can_cancel_join` | `user_role == participant && status != cancelled` | 参与者可取消报名 |

## 工作流

### 工作流 A：后端 service 层统一改造

**目标**：为所有六个类型的详情接口补齐 `user_role`、`is_owner`、`can_xxx` 字段，将业务规则收敛到后端。

1. 提取通用 `userRole` 计算函数，统一为所有类型计算 `publisher / acceptor / participant / viewer`
2. 为每个类型新增 `can_edit`、`can_delete` 计算逻辑
3. 为每个类型新增特有 `can_xxx` 计算逻辑
4. 修改各类型的 `Get*Detail` 方法，在返回 payload 中包含完整的主态/客态字段
5. 修改各类型的列表接口，在列表项中也返回 `is_owner` 和关键 `can_xxx`（用于列表页操作按钮渲染）
6. 补齐 LostFound 详情的 `can_view_contact` 字段
7. 补齐 Resource 详情的联系方式裁剪逻辑
8. 新增各类型的"下架/删除"接口（`DELETE /api/v1/campus-life/{type}/{id}` 或 `POST .../offline`）
9. 新增各类型的"编辑"接口（`PUT /api/v1/campus-life/{type}/{id}`）
10. 新增 LostFound 的"标记已找到/已认领"接口
11. 新增 Carpool 的"加入拼车"接口
12. 将 Errand 的 `canAccept`、`canCancelPublish`、`canCancelAccept` 逻辑从客户端移到后端

当前进展（截至 `2026-05-13`）：

- ✅ 已为所有六个类型补齐 `user_role`、`is_owner`、`can_view_contact`、`can_edit`、`can_delete` 字段（详情+列表）。
- ✅ 已为各类型补齐特有 `can_xxx`：Market `can_favorite`、Errand `can_accept`/`can_cancel_accept`、Resource `can_download`、LostFound `can_mark_resolved`、Carpool `can_join_carpool`、Meetup `can_edit`/`can_delete`。
- ✅ 已补齐 LostFound 详情 `can_view_contact` 和 Resource 联系方式裁剪。
- ✅ 已新增 Market/Resource/LostFound/Carpool 下架接口（`POST /{type}/delete/:id`）。
- ✅ 已新增 LostFound 标记已找到/已认领接口（`POST /lostFound/resolve/:id`）。
- ✅ 已将 Errand 小程序详情页的本地推断逻辑移除，改为消费后端 `can_xxx`。
- ✅ 已为 Market/LostFound 详情页添加主态操作（下架/标记已找到）。
- ✅ 已为 Meetup 补齐 `can_edit`/`can_delete` 字段映射。
- ✅ `go test ./...` 和 `make check` 全量通过。
- 🔲 A6（编辑接口）待后续迭代。
- 🔲 A8（Carpool 加入拼车接口）待后续迭代。

当前状态：A1-A5/A7/A9 已完成；A6、A8 待后续迭代。

### 工作流 B：OpenAPI 契约更新

**目标**：为所有详情接口定义具体的响应 schema，使 `user_role`、`can_xxx` 等字段在契约层面可见。

1. 为每个类型定义详情响应 schema（`MarketDetailResponse`、`ErrandDetailResponse` 等）
2. 在 schema 中明确 `user_role`、`is_owner`、`can_view_contact`、`can_edit`、`can_delete` 及各类型特有 `can_xxx` 字段
3. 为新增的下架/删除、编辑、标记已找到等接口定义请求/响应 schema
4. 重新生成 JS/Dart SDK

当前进展（截至 `2026-05-13`）：

- ✅ 为每个类型定义了详情响应 schema：`MarketDetail`、`ErrandDetail`、`ResourceDetail`、`LostFoundDetail`、`CarpoolDetail`、`MeetupDetail`，以及对应的 `*ResponseEnvelope`。
- ✅ 在 schema 中明确了 `user_role`、`is_owner`、`can_view_contact`、`can_edit`、`can_delete` 及各类型特有 `can_xxx` 字段。
- ✅ 新增通用 `CampusLifeOwnerContext` schema，描述所有类型共用的主态/客态上下文字段。
- ✅ 为新增的下架/删除接口定义了路径和响应：`/api/market/delete/{id}`、`/api/resource/delete/{id}`、`/api/lostFound/delete/{id}`、`/api/carpool/delete/{id}`。
- ✅ 为 LostFound 标记已找到接口定义了路径和响应：`/api/lostFound/resolve/{id}`。
- ✅ 将 6 个详情接口的响应从 `GenericObjectResponseEnvelope` 改为具体 schema。
- 🔲 重新生成 JS/Dart SDK（待后续迭代）。

当前状态：B1-B3 已完成；B4（SDK 生成）待后续迭代。

### 工作流 C：小程序详情页主态/客态 UI 改造

**目标**：小程序各详情页根据后端返回的 `can_xxx` 字段渲染不同的操作按钮，移除所有本地推断逻辑。

1. **Errand 详情页**：移除本地 `canAccept`/`canCancelPublish`/`canCancelAccept` 推断，改为消费后端返回的 `can_accept`/`can_delete`/`can_cancel_accept`
2. **Market 详情页**：新增主态操作栏（编辑、下架按钮），根据 `can_edit`/`can_delete` 渲染
3. **LostFound 详情页**：新增主态操作栏（编辑、下架、标记已找到按钮），根据 `can_edit`/`can_delete`/`can_mark_resolved` 渲染
4. **Carpool**：新增详情页，区分主态（编辑、取消发布）和客态（查看联系方式、加入拼车）
5. **Resource 详情页**：新增主态操作栏（编辑、下架按钮），补齐联系方式可见性控制
6. **Meetup 详情页**：将 `can_cancel_publish` 对齐为 `can_delete`，新增 `can_edit` 按钮

当前进展（截至 `2026-05-13`）：

- ✅ Errand 详情页：移除本地 `canAccept`/`canCancelPublish`/`canCancelAccept` 推断，改为消费后端返回的 `can_accept`/`can_delete`/`can_cancel_accept`。WXML 底部栏将 `canCancelPublish` 改为 `canDelete`，修复语义混淆。添加"已取消"状态横幅。
- ✅ Market 详情页：JS 层新增 `is_owner`/`can_edit`/`can_delete`/`can_favorite`/`status` 字段映射 + `onDelete` 事件处理。WXML 底部栏新增主态下架按钮（`isOwner` 时显示下架，客态显示想要+联系方式）。添加"已下架"/"审核未通过"/"审核中"状态横幅。
- ✅ LostFound 详情页：JS 层新增 `is_owner`/`can_edit`/`can_delete`/`can_mark_resolved`/`status` 字段映射 + `onDelete`/`onMarkResolved` 事件处理。修复 `status` 字段被中文标签覆盖的问题（新增 `statusLabel` 用于显示，`status` 保留后端原始值）。WXML 底部栏新增主态操作区（已找到/已认领 + 下架按钮），客态保留联系方式按钮。添加"已下架"/"审核未通过"/"审核中"/"已找到"/"已认领"状态横幅。
- ✅ Meetup 详情页：service 层补齐 `is_owner`/`can_edit`/`can_delete` 字段映射。WXML 底部栏将 `canCancelPublish` 改为 `canDelete`。添加"已取消"状态横幅。
- ✅ API 模块：4 个模块新增 `deleteXxx` 和 `resolveLostFound` API 封装。
- ✅ 样式：Market/LostFound/Errand/Meetup 新增状态横幅样式和 owner 按钮样式。
- 🔲 Carpool 详情页：待后续迭代（当前仅有列表页）。
- 🔲 编辑按钮渲染：所有页面的 `canEdit` 已映射但 WXML 未渲染编辑按钮（待编辑接口 A6 完成后补齐）。

当前状态：C1-C3/C5 已完成（含 WXML 模板改造和状态横幅）；C4（编辑按钮渲染）待 A6 完成后补齐；C6（Carpool 详情页）待后续迭代。

### 工作流 D：管理后台适配

**目标**：确保管理后台的校园生活管理页与新增的下架/编辑接口兼容。

1. 管理后台已有审核下线/重新发布能力，不需要主态/客态区分
2. 确认管理后台的详情展示与新增字段不冲突
3. 如有需要，在管理后台详情抽屉中展示 `user_role` 信息

当前进展（截至 `2026-05-13`）：

- ✅ 管理后台已有审核下线/重新发布能力，不需要主态/客态区分，确认与新增字段不冲突。

当前状态：已完成。

## 验收标准

### 后端验收

- [x] 所有六个类型的详情接口返回 `user_role`、`is_owner`、`can_view_contact`、`can_edit`、`can_delete`
- [x] 各类型特有 `can_xxx` 字段正确返回（Errand: `can_accept`/`can_cancel_accept`；Meetup: `can_join`/`can_cancel_join`；LostFound: `can_mark_resolved`；Carpool: `can_join_carpool`；Resource: `can_download`；Market: `can_favorite`）
- [x] 发布者调用下架/删除接口成功，非发布者调用返回 403
- [ ] 发布者调用编辑接口成功，非发布者调用返回 403（待 A6 完成）
- [x] LostFound 详情返回 `can_view_contact`
- [x] Resource 详情联系方式由后端裁剪
- [x] `go test ./...` 全量通过

### 契约验收

- [x] OpenAPI 中定义了各类型详情响应的具体 schema
- [x] 新增接口的请求/响应 schema 已定义
- [ ] JS/Dart SDK 可正常生成（待 B4 完成）

### 小程序验收

- [x] Errand 详情页不再本地推断 `can_xxx`，全部消费后端返回值
- [x] Market 详情页主态显示下架按钮，客态显示收藏/查看联系方式
- [x] LostFound 详情页主态显示下架/标记已找到按钮
- [ ] Carpool 有独立详情页，区分主态/客态（待 C6 完成）
- [ ] Resource 详情页主态显示编辑/下架按钮（待 C4 完成）
- [x] 所有详情页的操作按钮由后端 `can_xxx` 控制，无本地推断（Errand/Market/LostFound/Meetup 已完成，WXML 已改造）

## 风险

- 各类型状态枚举不完全一致（Errand 有 `accepted`，Meetup 有 `open`/`full`），`can_edit`/`can_delete` 的条件需逐类型校验
- 新增编辑接口需考虑部分字段不可编辑（如已接单的跑腿不可修改联系方式）
- Carpool 当前无详情页，需从零创建
- 小程序 Errand 详情页移除本地推断后，需确保后端返回的 `can_xxx` 与当前行为完全一致

## 风险应对

- 先完成后端改造和单测，再推进前端消费
- 每个类型独立改造、独立验证，不跨类型耦合
- 保留 Errand 现有 `user_role` 字段，在其基础上补充 `can_xxx`，不破坏现有前端
- Carpool 详情页优先级可适当降低，先保证列表页主态/客态区分

## 退出条件

- 所有六个类型的后端详情接口返回完整的 `user_role` + `can_xxx` 字段
- 小程序所有详情页根据后端 `can_xxx` 渲染操作按钮，无本地推断
- OpenAPI 契约与实际实现一致
- `make check` 全量通过
