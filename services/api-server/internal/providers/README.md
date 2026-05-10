# providers

外部系统接入统一放在本目录下，禁止在业务模块中直接拼接第三方 HTTP 调用。

当前阶段预留以下 provider：

- `wechat_provider`：当前已提供 mock 登录身份交换实现
- `sso_provider`
- `academic_provider`：当前已提供 mock 验证码和学生信息装配实现
- `sms_provider`
- `storage_provider`
- `push_provider`
