## Task: 后端存储与状态驱动重构

**ID:** api-server-storage-state-refactor
**Label:** API Server: 存储与状态驱动重构
**Description:** 基于 `MySQL + MongoDB + Redis + BMFS` 重新设计后端存储边界、社区领域模型与状态驱动机制，不直接沿用当前 `PostgreSQL` 数据结构
**Type:** Enhancement
**Status:** Backlog
**Priority:** High
**Created:** 2026-05-13
**Updated:** 2026-05-13
**PRD:** [api-server-storage-state-refactor.md](/Users/liangluo/code/weouc2026/.agent/PRDS/api-server-storage-state-refactor.md)

**Progress:** 已完成设计任务建档、目标架构划分、执行计划草案与文档同步；设计已确认“审核态并入统一 `status`，不迁移历史数据、不兼容旧链路、旧代码逻辑可直接删除”，待确认 `BMFS` 的具体 SDK/仓库后进入契约改造与分阶段实现。
