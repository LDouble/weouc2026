# weouc2026

校园综合应用单仓库，目标形态为：

- `Golang` 后端
- 微信原生小程序
- `Flutter` App
- 管理员后台

## 项目目标

- 一套业务中台，支撑多端一致能力
- 一套数据契约，减少多端接口漂移
- 一套工程约束，提升人和 AI 协作效率

## 先读文档

1. [AGENTS.md](/Users/liangluo/code/weouc2026/AGENTS.md)
2. [ARCHITECTURE.md](/Users/liangluo/code/weouc2026/ARCHITECTURE.md)
3. [docs/product-specs/campus-super-app-v1.md](/Users/liangluo/code/weouc2026/docs/product-specs/campus-super-app-v1.md)
4. [docs/exec-plans/active/0001-foundation.md](/Users/liangluo/code/weouc2026/docs/exec-plans/active/0001-foundation.md)

## 仓库结构

```text
.
├── .agent/                  # 任务和 PRD
├── apps/                    # 客户端应用
│   ├── admin-web/
│   ├── miniapp-wechat/
│   └── mobile-flutter/
├── docs/                    # 系统设计、规范、计划、参考资料
├── ops/                     # 部署、环境、运维相关说明
├── packages/                # 契约、SDK、共享工具
└── services/                # 后端服务
    └── api-server/
```

## 当前关键决策

- 后端采用 `Go 模块化单体` 起步，不直接上微服务
- 对外接口以 `OpenAPI` 为单一事实来源
- 微信端采用 `原生小程序`，不引入跨端框架
- `Flutter` 采用 `feature-first + MVVM + repository`
- 管理后台采用 `Vue 3 + TypeScript + Vite + Ant Design Vue`
- 重大设计先写文档，再改代码，并同步更新计划

## 下一步建议

1. 初始化 `services/api-server` 的目录骨架与基础中间件
2. 先完成 `认证、权限、内容门户、通知消息` 四个基础域
3. 建立 `OpenAPI -> JS SDK / Dart SDK` 的生成链路
4. 再分别启动管理员后台、小程序、Flutter App 的外壳工程

## 最小校验入口

```bash
# 使用本地依赖执行最小校验
make check

# 使用与 CI 一致的安装+校验流程
make ci-check
```

## 契约生成入口

```bash
# 基于 OpenAPI 生成 JS / Dart SDK
make generate-sdk
```
