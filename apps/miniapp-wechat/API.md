# 微信小程序依赖接口说明

本文档只描述 `apps/miniapp-wechat` 当前已经依赖的接口，不覆盖后台或 Flutter 侧尚未消费的能力。

## 1. 通用约定

### 1.1 Base URL

- 开发环境默认使用 `http://localhost:8080/api`
- 小程序当前通过 [config.js](/Users/liangluo/code/weouc2026/apps/miniapp-wechat/config.js) 读取 `baseUrl`

### 1.2 鉴权方式

- 默认使用 `Authorization: Bearer <token>`
- 登录接口和少量公开接口可跳过鉴权
- 小程序请求层在收到 `401` 时会尝试重新建立会话；若恢复失败，再跳转登录页

### 1.3 响应兼容规则

小程序请求层当前兼容三类响应：

1. 直接返回业务对象或列表对象
2. 返回 `api-server` 当前默认包裹结构：

```json
{
  "request_id": "req-123",
  "data": {}
}
```

3. 返回历史标准包裹结构：

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

约定：

- 若存在 `code` 字段，则 `200-299` 视为成功
- 若存在 `error` 对象，则按失败处理，并优先读取 `error.message`
- 列表接口统一推荐返回：

```json
{
  "list": [],
  "total": 0,
  "page": 1,
  "pageSize": 20
}
```

### 1.4 权限与受限字段约束

- 联系方式、教务状态等受限信息必须由后端裁剪后再返回
- 前端不允许根据本地状态自行推断联系方式是否可见
- 详情接口推荐显式返回 `can_view_contact`，列表接口推荐在无权场景直接返回空联系方式
- 校园生活新发布内容默认进入 `reviewing`，公开列表只显示 `published`；“我的发布”依赖 `review_status` 展示待审/已发布/已下架

## 2. 会话与身份

### 2.1 微信登录

- 方法：`POST /auth/wechat/login`
- 使用位置：登录页、会话恢复、请求层 `401` 后重登

请求体：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `code` | `string` | 是 | `wx.login` 返回的临时凭证 |
| `app_id` | `string` | 是 | 小程序 `AppID` |

成功响应示例：

```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "jwt-token",
    "openid": "oAbCdEf",
    "userInfo": {
      "userId": 1001,
      "nickname": "海大同学",
      "avatarUrl": "https://example.com/avatar.png"
    }
  }
}
```

前端消费字段：

- `token`
- `openid`
- `userInfo.nickname`
- `userInfo.avatarUrl`

### 2.2 获取当前用户资料

- 方法：`GET /student`
- 使用位置：我的页、教务绑定页、会话恢复后的资料同步

成功响应建议字段：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `name` | `string` | 用户姓名或展示名 |
| `avatar_url` | `string` | 头像 |
| `student_id` | `string` | 学号 |
| `major` | `string` | 专业 |
| `college` | `string` | 学院 |
| `grade` | `string` | 年级 |
| `is_bound` | `boolean` | 是否已完成教务绑定 |
| `updated_at` | `string` | 最近绑定或更新资料时间 |

说明：

- 若用户已登录但尚未建立学生资料，当前前端兼容 `404`
- 若返回 `student_id`，前端会默认视为已绑定教务

### 2.3 发送教务绑定验证码

- 方法：`POST /edu/send-captcha`
- 使用位置：教务绑定页

请求体：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `sid` | `string` | 是 | 学号 |

成功响应：

```json
{
  "code": 200,
  "message": "验证码已发送"
}
```

### 2.4 提交教务绑定

- 方法：`POST /student`
- 使用位置：教务绑定页

说明：

- 当前前端把该接口当作“提交教务绑定”的入口消费
- 后续建议在契约层明确为更清晰的绑定语义接口，避免和通用资料创建混淆

请求体：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `student_id` | `string` | 是 | 学号 |
| `password` | `string` | 是 | 教务系统密码 |
| `captcha` | `string` | 是 | 验证码 |

成功响应至少应包含：

- `student_id`
- `name`
- `major`
- `college`
- `grade`
- `is_bound`
- `updated_at`

### 2.5 解绑教务

- 方法：`PUT /student`
- 使用位置：教务绑定页

请求体：

```json
{
  "is_bound": false
}
```

说明：

- 当前前端只依赖解绑动作成功，不强依赖返回体
- 若服务端返回最新资料对象，建议保持和 `GET /student` 同构

## 3. 首页动态

### 3.1 获取动态流

- 方法：`GET /feed/list`
- 使用位置：首页

查询参数：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `page` | `number` | 否 | 页码，默认 `1` |
| `pageSize` | `number` | 否 | 每页条数，默认 `10` 或 `20` |
| `feed_types` | `string[]` 或 `string` | 否 | 指定动态类型 |
| `keyword` | `string` | 否 | 搜索关键字 |
| `user_role` | `string` | 否 | 我的页统计场景会传 `publisher` |

列表项建议字段：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `string` | 动态 ID |
| `feed_type` | `string` | `market` / `errand` / `resource` / `lostFound` / `carpool` / `meetup` |
| `feed_type_label` | `string` | 动态类型展示文案 |
| `title` | `string` | 标题 |
| `desc` | `string` | 摘要 |
| `publisher` | `string` | 发布人 |
| `created_at` | `string` | 创建时间 |
| `review_status` | `string` | `reviewing` / `published` / `rejected` / `offline` |
| `image` | `string` | 封面图 |
| `extra.images` | `string[]` | 备用图片列表 |

## 4. 二手交易

### 4.1 获取列表

- 方法：`GET /market/list`
- 使用位置：二手交易列表页

查询参数：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `category` | `string` | 分类值，如 `digital`、`book`、`wanted` |
| `keyword` | `string` | 搜索关键词 |
| `page` | `number` | 页码 |
| `pageSize` | `number` | 每页条数 |

列表项建议字段：

- `id`
- `title`
- `desc`
- `publisher`
- `publisher_initial`
- `image`
- `likes`
- `liked`
- `extra.category`
- `extra.price`
- `extra.original_price`
- `extra.condition`
- `extra.images`
- `extra.is_favorited`

### 4.2 获取详情

- 方法：`GET /market/detail/{id}`
- 使用位置：二手交易详情页

详情额外字段建议：

- `can_view_contact`
- `extra.contact`
- `extra.trade_mode`

说明：

- 若 `can_view_contact=false`，前端会展示“前往教务绑定”的引导
- 详情页的“想要”按钮已经对接 `POST /market/favorite`

### 4.3 收藏/取消收藏

- 方法：`POST /market/favorite`
- 使用位置：二手交易列表卡片

请求体：

```json
{
  "product_id": "market-1001",
  "action": "add"
}
```

### 4.4 发布

- 方法：`POST /market/publish`
- 使用位置：统一发布页

请求体字段：

- `title`
- `desc`
- `price`
- `original_price`
- `category`
- `condition`
- `trade_mode`
- `contact`
- `images`

说明：

- `images` 现在传的是 COS 对象路径数组，不再传临时下载 URL
- 详情与列表中的图片 URL 由后端读取时动态签发

## 5. 跑腿

### 5.1 获取列表

- 方法：`GET /errand/list`
- 使用位置：跑腿列表页、我的页统计

查询参数：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `category` | `string` | 跑腿类型 |
| `keyword` | `string` | 搜索关键词 |
| `user_role` | `string` | `publisher` / `acceptor` |
| `page` | `number` | 页码 |
| `pageSize` | `number` | 每页条数 |

列表项建议字段：

- `id`
- `title`
- `desc`
- `category`
- `route_start`
- `route_end`
- `deadline`
- `reward`
- `publisher`
- `publisher_initial`
- `created_at`
- `status`
- `user_role`
- `is_accepted`

### 5.2 获取详情

- 方法：`GET /errand/detail/{id}`

### 5.3 接单

- 方法：`POST /errand/accept`

请求体：

```json
{
  "task_id": "errand-1001"
}
```

### 5.4 发布

- 方法：`POST /errand/publish`

说明：

- `images` 现在传的是 COS 对象路径数组，不再传临时下载 URL

## 6. 资料

### 6.1 获取列表

- 方法：`GET /resource/list`
- 使用位置：资料列表页

查询参数：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `category` | `string` | 分类值 |
| `keyword` | `string` | 搜索关键词 |
| `page` | `number` | 页码 |
| `pageSize` | `number` | 每页条数 |

列表项建议字段：

- `id`
- `title`
- `desc`
- `publisher`
- `publisher_initial`
- `created_at`
- `extra.category`
- `extra.course_name`
- `extra.views`
- `extra.likes`
- `extra.files`
- `extra.files[].name`
- `extra.files[].url`
- `extra.files[].file_type`
- `extra.files[].file_size`

### 6.2 获取详情

- 方法：`GET /resource/detail/{id}`
- 使用位置：资料列表页插顶场景

### 6.3 收藏

- 方法：`POST /resource/favorite`

### 6.4 发布

- 方法：`POST /resource/publish`

请求体字段：

- `title`
- `desc`
- `category`
- `course_name`
- `contact`
- `file_paths`

说明：

- 当前资料发布强依赖上传接口先拿到 `file_paths`
- 若上传接口未部署，前端会明确提示“文件上传暂不可用”，不会再伪装成本地发布成功
- `file_paths` 为稳定对象路径，真正可访问链接由后端详情/列表接口动态签发

## 7. 失物招领

### 7.1 获取列表

- 方法：`GET /lostFound/list`
- 使用位置：失物招领列表页

查询参数：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `type` | `string` | `lost` / `found` |
| `category` | `string` | 分类值 |
| `keyword` | `string` | 搜索关键词 |
| `page` | `number` | 页码 |
| `pageSize` | `number` | 每页条数 |

列表项建议字段：

- `id`
- `title`
- `desc`
- `publisher`
- `publisher_initial`
- `extra.type`
- `extra.category`
- `extra.location`
- `extra.event_time`
- `extra.item_feature`
- `extra.contact`

### 7.2 获取详情

- 方法：`GET /lostFound/detail/{id}`

说明：

- 当前前端优先消费裁剪后的 `extra.contact`
- 若后端后续补充 `can_view_contact`，详情页会直接消费该字段并给出绑定引导

### 7.3 发布

- 方法：`POST /lostFound/publish`

## 8. 校园拼车

### 8.1 获取列表

- 方法：`GET /carpool/list`
- 使用位置：拼车列表页

查询参数：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `category` | `string` | `today` / `tomorrow` / `week` / `longterm` |
| `keyword` | `string` | 搜索出发地、目的地或发起人 |
| `page` | `number` | 页码 |
| `pageSize` | `number` | 每页条数 |

列表项建议字段：

- `id`
- `category`
- `from`
- `to`
- `time`
- `type`
- `seats_text`
- `price`
- `note`
- `tags`
- `contact`
- `review_status`
- `publisher`
- `publisher_initial`
- `created_at`

说明：

- `contact` 必须继续由后端按教务绑定状态裁剪
- `time` 由后端基于稳定出发时间动态格式化为“今天 18:30 / 明天 09:00 / 5月18日 14:00”
- 新发布拼车默认进入 `reviewing`；小程序发布成功后会带 `insertId` 回列表页，依赖详情接口插顶展示自己的待审记录

### 8.2 获取详情

- 方法：`GET /carpool/detail/{id}`
- 使用位置：拼车列表页插顶回流场景

### 8.3 发布

- 方法：`POST /carpool/publish`
- 使用位置：拼车发布页

请求体字段：

- `category`
- `travel_date`
- `travel_time`
- `from`
- `to`
- `type`
- `seats_text`
- `price`
- `note`
- `tags`
- `contact`

说明：

- 小程序现在会传稳定的 `travel_date + travel_time`，不再只传中文展示时间
- 后端会据此统一计算 `today / tomorrow / week / longterm` 的筛选结果与展示文案

## 9. 校园组局

### 9.1 获取列表

- 方法：`GET /meetup/list`
- 使用位置：校园组局列表页

查询参数：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `category` | `string` | `study` / `sports` / `food` / `game` / `activity` / `other` / `all` |
| `keyword` | `string` | 搜索主题、地点或发起人 |
| `user_role` | `string` | `publisher` / `participant` / `viewer` |
| `page` | `number` | 页码 |
| `pageSize` | `number` | 每页条数 |

列表项建议字段：

- `id`
- `category`
- `title`
- `desc`
- `location`
- `start_at`
- `deadline_at`
- `max_participants`
- `joined_count`
- `remaining_seats`
- `fee_text`
- `tags`
- `contact`
- `status`
- `review_status`
- `user_role`
- `can_join`
- `can_cancel_join`
- `can_cancel_publish`
- `publisher`
- `publisher_initial`
- `created_at`

说明：

- 公开列表只展示 `published` 的组局；发起人和审核员可继续查看待审内容
- `status` 会综合 `review_status` 与人数、截止时间、开始时间，常见值包括 `open` / `full` / `cancelled` / `reviewing`
- `contact` 必须继续由后端按教务绑定状态裁剪

### 9.2 获取详情

- 方法：`GET /meetup/detail/{id}`
- 使用位置：组局详情页、发布成功后的回流页

说明：

- 详情接口建议继续返回和列表同构的核心字段，并显式返回 `can_view_contact`
- 发起人可以在详情页看到自己刚发布但仍处于 `reviewing` 的记录

### 9.3 发布

- 方法：`POST /meetup/publish`
- 使用位置：组局发布页

请求体字段：

- `category`
- `title`
- `desc`
- `location`
- `start_at`
- `deadline_at`
- `max_participants`
- `fee_text`
- `tags`
- `contact`

说明：

- `start_at` 与 `deadline_at` 使用 `RFC3339` 时间字符串
- `deadline_at` 可为空；为空时默认与 `start_at` 一致
- 新发布内容默认进入 `reviewing`，发布成功后前端会用返回的 `id` 跳转详情页

### 9.4 报名

- 方法：`POST /meetup/join`
- 使用位置：组局详情页

请求体：

```json
{
  "meetup_id": "meetup-101"
}
```

### 9.5 取消报名

- 方法：`POST /meetup/cancel-join`
- 使用位置：组局详情页

请求体：

```json
{
  "meetup_id": "meetup-101"
}
```

### 9.6 取消组局

- 方法：`POST /meetup/cancel-publish`
- 使用位置：组局详情页，仅发起人可见

请求体：

```json
{
  "meetup_id": "meetup-101"
}
```

## 10. 文件上传

### 10.1 获取 COS 临时凭证

- 方法：`GET /upload/cos-sts`
- 使用位置：闲置发布页图片上传、跑腿图片上传、资料文件上传

查询参数：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `scene` | `string` | 否 | `market` / `errand` / `resource`，用于服务端隔离对象前缀 |

成功响应建议字段：

- `tmp_secret_id`
- `tmp_secret_key`
- `session_token`
- `expired_time`
- `start_time`
- `bucket`
- `region`
- `path_prefix`

说明：

- 前端会用 `path_prefix + hash/...` 生成实际对象键
- `path_prefix` 已包含业务场景、用户隔离和日期维度

### 10.2 获取预签名访问地址

- 方法：`POST /upload/presigned-get`
- 使用位置：文件直传完成后回显

请求体：

```json
{
  "path": "miniapp/market/ab/cdef1234.png"
}
```

成功响应建议字段：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "url": "https://example.com/presigned-url"
  }
}
```

## 11. 后续契约收敛建议

- 统一把列表接口收敛为同一分页包结构，避免小程序为每个业务单独兼容
- `POST /student` 当前承担“教务绑定”语义，建议后续在 `packages/contracts` 中拆出更明确的绑定接口
- `lostFound` 路径建议后续统一命名风格，例如与业务域命名对齐为 `lost-found`
- 所有受限联系方式接口建议统一返回显式可见性字段，而不是让前端依赖“有值/没值”猜测权限
- 当前已改为“业务仅持久化对象路径、后端读取时签 URL”，后续若引入文件元数据中心，应继续保持这一约束
