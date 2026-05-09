# Flutter App

## 技术选型

- `Flutter`
- `Riverpod`
- `go_router`
- `dio`
- `json_serializable`

## 目标

承载更完整、更沉浸、更可持续演进的移动端体验。

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

