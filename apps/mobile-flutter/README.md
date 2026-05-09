# Flutter App

## 技术选型

- `Flutter`
- `Riverpod`
- `go_router`
- `dio`
- `json_serializable`

## 目标

承载与小程序一致业务语义的更完整移动端体验，当前以跑腿、组局（含拼车等轻社交撮合）、二手交易、资料、失物招领为主，后续再扩展教务功能。

## 推荐结构

```text
lib/
├── features/
├── shared/
└── main.dart
```

## 约束

- 采用 `feature-first + MVVM + repository`
- 业务规则不在页面层实现
- 缓存和同步逻辑进入数据层
- 联系方式可见性必须复用后端基于教务绑定状态的裁决结果
