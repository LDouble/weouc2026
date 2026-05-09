# 管理员后台

## 技术选型

- `Vue 3`
- `TypeScript`
- `Vite`
- `Pinia`
- `Vue Router`
- `Ant Design Vue`

## 目标

承载内容运营、用户权限、活动管理、工单处理、统计看板等管理能力。

## 推荐结构

```text
src/
├── app/
├── modules/
│   ├── iam/
│   ├── portal/
│   ├── activity/
│   ├── service-desk/
│   └── analytics/
└── shared/
```

## 约束

- 所有接口通过共享 SDK 调用
- 业务权限以后端权限码为准
- 通用页面能力优先抽象复用

