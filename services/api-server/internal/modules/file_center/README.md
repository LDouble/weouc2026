# file_center

当前已落地最小文件管理闭环：

- `GET /api/upload/cos-sts`：签发腾讯 COS 临时上传凭证
- `POST /api/upload/presigned-get`：按对象路径签发 GET 下载地址

约束：

- 业务层只持久化稳定对象路径
- 面向客户端返回可访问 URL 时，由后端按需签名
- COS 接入只允许经 `storage_provider`
