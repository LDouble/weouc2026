# migrations

数据库结构迁移脚本统一放在本目录，后续按时间戳或版本号管理。

当前约定：

- `*.sql` 文件按文件名字典序执行
- 服务启动时可通过 `API_SERVER_AUTO_MIGRATE=true` 自动执行未应用迁移
- 当前已内置 IAM 用户表迁移
