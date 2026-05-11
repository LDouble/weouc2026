# modules

当前已落地：

- `system`：系统健康检查、就绪检查、会话画像等基础能力
- `iam`：微信登录、Bearer Token 会话、当前用户资料、教务绑定，以及 `PostgreSQL + Redis` 持久化
- `campus_life`：首页动态、二手、跑腿、资料、失物招领、拼车的基础交互，以及 `memory / postgres` 双后端
- `file_center`：腾讯 COS 临时上传凭证、下载预签名与文件路径引用能力

预留或待增强业务域：

- `portal`
- `academic`
- `notification`
- `analytics`
