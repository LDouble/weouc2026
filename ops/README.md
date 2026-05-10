# 运维与环境

## 目标

管理环境配置、部署规范、备份恢复和监控接入方式。

## 后续建议

- `environments/`：开发、测试、预发、生产配置模板
- `docker/`：本地联调与基础镜像
  当前已落地 [ops/docker/api-server/README.md](/Users/liangluo/code/weouc2026/ops/docker/api-server/README.md)，用于拉起 `api-server + postgres + redis`
- `k8s/`：生产部署清单
- `runbooks/`：故障处理手册
