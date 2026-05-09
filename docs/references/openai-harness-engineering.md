# OpenAI Harness Engineering 参考摘要

参考文章：

- [OpenAI Harness Engineering](https://openai.com/zh-Hans-CN/index/harness-engineering/)

## 对本仓库最重要的启发

### 1. `AGENTS.md` 应是目录和约束，不是超长手册

因此本仓库将 `AGENTS.md` 用于：

- 指向关键文档
- 固化硬性架构边界
- 约束提交前自检

### 2. 文档要成为系统记录

因此本仓库将：

- 用 `docs/` 管理架构、规格、计划、质量
- 用 `.agent/` 管理任务与 PRD
- 要求变更同步更新文档

### 3. 清晰架构比“聪明提示词”更重要

因此本仓库优先做：

- 明确分层
- 明确目录
- 明确边界
- 明确契约

### 4. 计划和质量评分要持续维护

因此本仓库建立：

- [docs/PLANS.md](/Users/liangluo/code/weouc2026/docs/PLANS.md)
- [docs/QUALITY_SCORE.md](/Users/liangluo/code/weouc2026/docs/QUALITY_SCORE.md)

