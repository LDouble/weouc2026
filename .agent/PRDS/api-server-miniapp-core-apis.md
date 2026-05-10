# 小程序核心接口落地 PRD

## 背景

`services/api-server` 已完成基础工程初始化，但 `apps/miniapp-wechat` 当前依赖的真实后端接口尚未落地，导致登录、我的页、教务绑定页、首页和校园生活列表页仍无法与后端联调。

## 目标

在保持模块化单体与契约优先约束的前提下，为小程序当前主链路补齐第一批可联调接口：

- 微信登录与会话鉴权
- 当前用户资料查询
- 教务验证码发送与教务绑定/解绑
- 首页动态流
- 跑腿、二手、资料、失物招领列表
- 二手、跑腿、资料、失物招领详情中的关键接口

## 非目标

- 不接入真实微信开放平台
- 不接入真实教务系统
- 不实现真实对象存储直传
- 不实现完整审核流、通知中心、管理后台接口
- 不实现数据库持久化，当前阶段允许使用内存仓储和种子数据

## 设计约束

- 登录和教务绑定的外部能力必须走 `providers`
- 鉴权必须支持小程序当前 `Authorization: Bearer <token>` 方式
- 联系方式等受限字段必须由后端按教务绑定状态裁剪
- 列表接口需支持最基本的筛选、分页和用户角色视图

## 本轮范围

### 1. IAM 会话与资料

- `POST /api/auth/wechat/login`
- Bearer Token 鉴权
- `GET /api/student`
- `POST /api/student`
- `PUT /api/student`
- `POST /api/edu/send-captcha`

### 2. 校园生活基础接口

- `GET /api/feed/list`
- `GET /api/market/list`
- `GET /api/market/detail/{id}`
- `POST /api/market/favorite`
- `GET /api/errand/list`
- `GET /api/errand/detail/{id}`
- `POST /api/errand/accept`
- `POST /api/errand/cancel-publish`
- `POST /api/errand/cancel-accept`
- `GET /api/resource/list`
- `GET /api/resource/detail/{id}`
- `GET /api/lostFound/list`
- `GET /api/lostFound/detail/{id}`

### 3. 发布入口

- `POST /api/market/publish`
- `POST /api/errand/publish`
- `POST /api/resource/publish`
- `POST /api/lostFound/publish`

说明：

- 资料发布先支持“后端接受已上传文件路径”的语义，不在本轮打通真实上传。

## 验收标准

- 小程序登录后可以用 Bearer Token 请求资料与列表接口
- 未绑定教务时看不到受限联系方式
- 完成教务绑定后，市场和跑腿详情可以看到联系方式
- 首页、跑腿、二手、资料、失物招领列表能返回稳定种子数据
- 关键详情接口能被页面消费
- 至少存在覆盖登录、绑定、联系方式裁剪的自动化测试

## 风险

- 当前前端上传依赖 COS 直传，若不接真实存储，带文件的发布流程只能部分联调
- 使用内存仓储时，重启服务后会话和发布数据会丢失
- 当前接口命名包含历史包袱，如 `POST /student` 同时承担绑定语义，后续需在契约层继续收敛
