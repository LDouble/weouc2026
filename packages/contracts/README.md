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
