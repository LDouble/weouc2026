# providers

外部系统接入统一放在本目录下，禁止在业务模块中直接拼接第三方 HTTP 调用。

当前阶段预留以下 provider：

- `wechat_provider`：当前已提供 mock 登录身份交换实现
- `sso_provider`
- `academic_provider`：当前已提供 mock 验证码和学生信息装配实现
- `sms_provider`
- `storage_provider`：当前已提供腾讯 COS 上传临时凭证、下载预签名与健康检查实现
- `push_provider`
