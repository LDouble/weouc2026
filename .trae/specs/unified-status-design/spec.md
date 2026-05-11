# 校园生活统一状态设计 - 产品需求文档

## Overview

- **Summary**: 将审核状态（review\_status）整合到主状态（status）中，形成统一的状态模型，消除状态割裂问题
- **Purpose**: 简化状态管理逻辑，使状态流转更加清晰直观，避免业务逻辑中需要同时判断两个状态字段
- **Target Users**: 后端开发、前端开发、产品经理

## Goals

- 统一状态模型：将审核状态作为主状态的子集
- 简化业务逻辑：移除对 review\_status 的独立判断
- 向后兼容：确保现有接口行为不受影响

## Non-Goals (Out of Scope)

- 不改变现有的审核工作流程
- 不改变数据库表结构（暂时）

## Background & Context

当前校园生活模块存在两种状态字段：

- `status`: 业务状态（如 open, full, cancelled, accepted 等）
- `review_status`: 审核状态（reviewing, published, rejected, offline）

这种设计导致：

1. 状态判断逻辑复杂，需要同时检查两个字段
2. 前端需要合并两个状态来展示最终状态
3. 状态流转不直观

用户期望的状态模型：审核是主状态的一种，如：审核中、审核不通过、正常状态、取消状态

## Functional Requirements

- **FR-1**: 定义统一的状态枚举，包含审核相关状态
- **FR-2**: 修改服务层逻辑，使用统一状态替代两个独立状态
- **FR-3**: 前端适配统一状态模型
- **FR-4**: 更新接口返回数据结构，返回单一状态字段

## Non-Functional Requirements

- **NFR-1**: 保持 API 向后兼容性
- **NFR-2**: 状态流转规则清晰可追溯

## Constraints

- **Technical**: 需要同步修改后端服务和前端小程序
- **Dependencies**: 涉及多个模块（market, errand, resource, lostFound, carpool, meetup）

## Assumptions

- 审核状态是内容可见性的主要控制因素
- 原有业务状态（如 full, cancelled）与审核状态是互斥的

## Acceptance Criteria

### AC-1: 定义统一状态枚举

- **Given**: 需要处理校园生活内容状态
- **When**: 设计状态模型
- **Then**: 状态枚举包含：draft(草稿), reviewing(审核中), published(已发布/正常), rejected(审核不通过), offline(已下线), cancelled(已取消), full(已满)
- **Verification**: `programmatic`

### AC-2: 服务层状态判断逻辑简化

- **Given**: 业务逻辑中需要判断内容状态
- **When**: 检查内容是否可见/可操作
- **Then**: 只需检查单一状态字段，无需同时判断 review\_status 和 status
- **Verification**: `programmatic`

### AC-3: 前端状态展示逻辑简化

- **Given**: 前端展示内容卡片
- **When**: 需要显示状态标签
- **Then**: 直接使用后端返回的单一状态字段，无需二次合并
- **Verification**: `human-judgment`

### AC-4: 接口返回结构优化

- **Given**: 调用内容列表/详情接口
- **When**: 获取内容数据
- **Then**: 返回统一的 `status` 字段，包含审核状态信息
- **Verification**: `programmatic`

## Open Questions

- [ ] 不需要保留 review\_status 字段作为过渡
- [ ] 不需要需要迁移历史数据的状态

