# contracts

## 目标

作为跨端统一契约层，沉淀接口、事件、字典和生成 SDK 所需源文件。

## 建议内容

```text
contracts/
├── openapi/
├── events/
├── dictionary/
├── sdk-js/
└── sdk-dart/
```

## 约束

- 先改契约，再改实现
- 破坏性变更必须显式标注
- 客户端优先升级 SDK，不手写分叉接口层

## 当前契约

- [openapi/api-server.yaml](/Users/liangluo/code/weouc2026/packages/contracts/openapi/api-server.yaml)：后端基础工程系统接口基线

## SDK 生成流程（最小闭环）

1. 修改 `openapi/api-server.yaml` 并完成后端契约相关测试。
2. 在仓库根目录执行 `make generate-sdk` 生成 JS / Dart SDK。
3. 根据生成结果更新调用端（`admin-web`、`miniapp-wechat`、`mobile-flutter`）。
4. 将契约变更与消费改动在同一轮提交中说明。

### 命令入口

```bash
cd /Users/liangluo/code/weouc2026
make generate-sdk
```

说明：

- 生成脚本：`packages/contracts/scripts/generate-sdks.sh`
- 默认使用镜像：`openapitools/openapi-generator-cli:v7.16.0`
- 如需替换镜像，可设置环境变量：`OPENAPI_GENERATOR_IMAGE=<your-image>`

### 产物目录

- JS SDK：`packages/contracts/sdk-js/api-server/`
- Dart SDK：`packages/contracts/sdk-dart/api_server/`
