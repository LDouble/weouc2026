# iam

当前已落地：

- 微信小程序登录
- 管理员密码登录
- Bearer Token 会话解析
- 当前用户资料查询
- 教务验证码发送、绑定与解绑
- `PostgreSQL` 用户资料持久化
- `Redis` 会话与验证码持久化
- `GORM + MySQL` 用户仓储骨架
- `mysql_redis` IAM 后端装配入口与最小测试

说明：

- 当前仍使用 mock 微信与 mock 教务 provider
- 用户模型已扩展支持 username/password_hash 字段用于管理员账号
- 角色与权限当前以基础学生语义为主，后续再扩展管理员与组织模型
- 当前 `mysql_redis` 已具备运行时装配、`GORM AutoMigrate` 建表、健康探测与最小测试

目标重构方向（设计阶段）：

- 主事实源从当前 `PostgreSQL` 调整为 `MySQL`
- `users / accounts / roles / permissions / academic_bindings` 作为核心主数据留在 `MySQL`
- `Redis` 继续承担会话、验证码、幂等键等短期状态
- 联系方式与权限裁决的事实源仍以 `iam` 主数据为准，不下放到社区文档模型
