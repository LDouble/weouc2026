# 管理员后台

基于 Vue 3 + TypeScript + Vite + Ant Design Vue 的校园管理后台系统。

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **TypeScript** - 类型安全
- **Vite** - 下一代前端构建工具
- **Pinia** - Vue 状态管理
- **Vue Router** - Vue 路由管理
- **Ant Design Vue** - UI 组件库

## 项目结构

```text
src/
├── app/                    # 应用层
│   ├── layouts/            # 布局组件
│   │   └── MainLayout.vue  # 主布局（侧边栏 + 头部）
│   └── views/              # 公共视图
│       └── NotFoundView.vue
├── modules/                # 业务域模块
│   ├── iam/                # 身份权限管理
│   │   └── views/
│   │       ├── LoginView.vue
│   │       ├── UserListView.vue
│   │       ├── RoleListView.vue
│   │       └── PermissionListView.vue
│   ├── portal/             # 内容门户管理
│   │   └── views/
│   │       ├── ArticleListView.vue
│   │       ├── BannerListView.vue
│   │       └── NoticeListView.vue
│   ├── campus-life/        # 校园生活管理
│   │   └── views/
│   │       ├── ErrandListView.vue
│   │       ├── MeetupListView.vue
│   │       ├── ListingListView.vue
│   │       ├── ResourceListView.vue
│   │       └── LostItemListView.vue
│   ├── moderation/         # 内容审核管理
│   │   └── views/
│   │       ├── PendingListView.vue
│   │       └── HistoryView.vue
│   └── analytics/          # 数据分析
│       └── views/
│           ├── DashboardView.vue
│           └── AuditLogView.vue
├── router/                 # 路由配置
│   └── index.ts
├── stores/                 # 状态管理
│   └── auth.ts
├── api/                    # API 封装
│   └── index.ts
└── App.vue                 # 根组件
```

## 功能模块

| 模块 | 功能 |
|------|------|
| **IAM** | 用户管理、角色管理、权限管理、登录认证 |
| **Portal** | 文章入口预留、轮播管理、公告管理 |
| **Campus Life** | 跑腿服务、组局活动、二手交易、资料共享、失物招领 |
| **Moderation** | 待审核列表、审核历史 |
| **Analytics** | 数据仪表盘、审计日志 |

## 开发

```bash
# 安装依赖
npm install

# 开发模式
npm run dev

# 构建生产版本
npm run build

# 预览生产版本
npm run preview
```

## 路由配置

路由定义在 `src/router/index.ts`，包含：
- `/login` - 登录页
- `/` - 仪表盘
- `/users` - 用户管理
- `/roles` - 角色管理
- `/permissions` - 权限管理
- `/portal/articles` - 文章管理预留页
- `/portal/banners` - 轮播管理
- `/portal/notices` - 公告管理
- `/campus-life/errands` - 跑腿服务
- `/campus-life/meetups` - 组局活动
- `/campus-life/listings` - 二手交易
- `/campus-life/resources` - 资料共享
- `/campus-life/lost-items` - 失物招领
- `/moderation/pending` - 待审核
- `/moderation/history` - 审核历史
- `/analytics/audit-logs` - 审计日志

## 权限控制

- 使用 `meta.requiresAuth` 标记需要认证的路由
- 登录状态通过 `localStorage` 存储 Token
- 权限码通过 `stores/auth.ts` 的 `hasPermission()` 方法校验

## 当前真实联调范围

- 轮播管理：已接入后端真实列表、创建、更新、删除接口
- 公告管理：已接入后端真实列表、详情、发布、更新、删除接口
- 校园生活管理：跑腿、组局、二手、资料、失物招领已切到真实审核数据源，并支持详情查看与状态操作
- 待审核内容：已接入真实审核列表与审核更新接口
- 审核历史：已切换为真实审核结果查询
- 文章管理：后端模型尚未实现，当前页面保留为预留入口，不再请求不存在的接口

## 架构约束

遵循 `/Users/liangluo/code/weouc2026/AGENTS.md` 规范：
- 业务规则只在后端落地；客户端只做展示、交互编排、轻量本地缓存
- 使用服务端返回的权限码控制路由和按钮
- 按业务域拆分模块，不按页面随意堆逻辑
