# modules

当前已落地：

- `system`：系统健康检查、就绪检查、会话画像等基础能力
- `iam`：微信登录、Bearer Token 会话、当前用户资料、教务绑定，以及 `PostgreSQL + Redis` 持久化
- `academic`：已绑定学生的学期、课程表、考试与成绩查询，当前经 mock provider 联调
- `portal`：门户首页聚合、公告列表/详情与管理端公告发布，支持 `memory / postgres` 双后端
- `notification`：当前用户通知列表、未读统计、已读回执与管理端通知发布，支持 `memory / postgres` 双后端
- `analytics`：审计日志列表、基础运营看板与统一审计存储，支持 `memory / postgres` 双后端
- `campus_life`：首页动态、二手、跑腿、资料、失物招领、拼车的基础交互，以及 `memory / postgres` 双后端
- `file_center`：腾讯 COS 临时上传凭证、下载预签名与文件路径引用能力

预留或待增强业务域：

- 更完整的 `academic` 同步编排与持久化快照
