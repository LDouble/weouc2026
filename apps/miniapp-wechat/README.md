# 微信小程序

## 技术选型

- 微信原生小程序
- JavaScript
- 原生组件与分包

## 目标

承载高频校园入口、微信生态传播与校园生活主线能力，当前以跑腿、组局（含拼车等轻社交撮合）、二手交易、资料、失物招领为主，后续再承接教务功能。

## 推荐结构

```text
.
├── api/
├── services/
├── stores/
├── pages/
├── components/
├── behaviors/
├── utils/
└── subpackages/
```

## 约束

- 不使用跨端框架
- 默认使用微信小程序登录
- 主包控制体积
- 服务层统一处理鉴权与错误映射
- 未绑定教务的账号不能查看联系方式，页面只按服务端返回结果渲染

## 分层职责

- `api/`：底层传输封装与接口模块，只描述路径、参数、鉴权头和返回包结构。
- `services/`：按业务域组织用例，负责 DTO 到页面语义的映射，不把接口字段拼装散落到页面。
- `stores/`：存放运行期配置与会话状态，只维护客户端必要状态，不承载业务规则。
- `pages/`：只负责交互编排与局部 `setData`，不直接推断绑定状态、联系方式可见性等规则。

## 当前已对齐的核心流

- 微信登录与会话恢复
- 我的页与教务绑定状态展示
- 首页动态流
- 门户公告卡片读取与通知未读角标
- 消息中心通知列表与已读回执
- 教务数据页读取 `academic` 契约接口
- 二手、跑腿、资料、失物招领列表页

## 接口文档

- [API.md](/Users/liangluo/code/weouc2026/apps/miniapp-wechat/API.md)：微信小程序当前依赖的后端接口说明

## 提测前 COS 检查项

- `uploadFile` 合法域名：微信公众平台需配置后端 API 域名，以及腾讯 COS 上传域名（例如 `https://<bucket>.cos.<region>.myqcloud.com`）。
- `downloadFile` 合法域名：微信公众平台需配置预签名 URL 所属 COS 下载域名，确保资料、图片回显和下载可访问。
- COS bucket CORS：允许小程序直传所需的 `PUT` / `OPTIONS` 请求、必要请求头和业务侧访问来源。
- STS 权限范围：临时凭证只允许写入后端返回的 `path_prefix` 范围，场景、用户与日期前缀需要和后端配置一致。
- 预签名 URL 有效期：确认 `/upload/presigned-get` 返回的下载地址有效期覆盖页面预览、详情查看和资料下载的提测场景。

## 本地校验

```bash
cd /Users/liangluo/code/weouc2026/apps/miniapp-wechat
npm run check:syntax
```

说明：

- 该校验会递归检查 `api/`、`services/`、`stores/`、`pages/`、`components/` 等目录下的 `.js` 语法正确性。
