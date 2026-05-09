# 校园综合应用基础架构 PRD

## 目标

为校园综合应用建立清晰、可扩展、AI 友好的工程架构，覆盖：

- Golang 后端
- 微信小程序
- Flutter App
- 管理员后台

## 关键决策

- 单仓库管理
- Go 模块化单体后端
- 原生微信小程序
- Flutter 采用 `feature-first + MVVM`
- 管理后台采用 `Vue 3 + TypeScript`
- 契约优先，OpenAPI 作为跨端事实来源

## 交付物

- [README.md](/Users/liangluo/code/weouc2026/README.md)
- [AGENTS.md](/Users/liangluo/code/weouc2026/AGENTS.md)
- [ARCHITECTURE.md](/Users/liangluo/code/weouc2026/ARCHITECTURE.md)
- [docs/product-specs/campus-super-app-v1.md](/Users/liangluo/code/weouc2026/docs/product-specs/campus-super-app-v1.md)
- [docs/exec-plans/active/0001-foundation.md](/Users/liangluo/code/weouc2026/docs/exec-plans/active/0001-foundation.md)

## 验收标准

- 能明确说明为什么当前阶段不直接上微服务
- 能明确说明四端的职责边界和技术选型
- 能明确说明后端模块边界、契约层和外部系统接入策略
- 文档能支持后续 AI 和人工协作继续落地工程

