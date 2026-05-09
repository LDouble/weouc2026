# 前端统一规范

适用范围：

- [apps/admin-web](/Users/liangluo/code/weouc2026/apps/admin-web)
- [apps/miniapp-wechat](/Users/liangluo/code/weouc2026/apps/miniapp-wechat)
- [apps/mobile-flutter](/Users/liangluo/code/weouc2026/apps/mobile-flutter)

## 统一原则

- 前端不保存业务真相，只保存展示态和局部缓存
- 所有前端都消费统一契约，不各自发明字段含义
- 权限控制必须后端裁决，前端只做体验层隐藏与禁用
- 错误码、空态、加载态、重试策略尽量统一

## 管理后台

- 优先面向运营效率和数据密度设计
- 表格、表单、字典项、上传流程做成共享能力

## 小程序

- 优先启动速度、分享转化、微信能力接入
- 主包轻量化，复杂流程拆分包

## Flutter App

- 优先长期留存、沉浸体验、离线体验与消息触达
- 页面状态和数据层严格分离

## 设计系统建议

- 建立统一颜色、间距、圆角、状态色规范
- 统一业务图标、空态图、错误提示语气
- 业务名词统一由产品和后端契约共同定义

