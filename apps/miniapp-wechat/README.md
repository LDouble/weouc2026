# 微信小程序

## 技术选型

- 微信原生小程序
- JavaScript
- 原生组件与分包

## 目标

承载高频校园入口、微信生态传播与轻量业务办理。

## 推荐结构

```text
.
├── pages/
├── components/
├── services/
├── stores/
├── behaviors/
├── utils/
└── subpackages/
```

## 约束

- 不使用跨端框架
- 主包控制体积
- 服务层统一处理鉴权与错误映射

