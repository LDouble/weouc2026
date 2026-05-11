import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: () => import('../modules/iam/views/LoginView.vue')
  },
  {
    path: '/',
    name: 'layout',
    component: () => import('../app/layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'dashboard',
        component: () => import('../modules/analytics/views/DashboardView.vue'),
        meta: { title: '仪表盘', permission: 'analytics:view' }
      },
      {
        path: '/users',
        name: 'users',
        component: () => import('../modules/iam/views/UserListView.vue'),
        meta: { title: '用户管理', permission: 'iam:user:view' }
      },
      {
        path: '/roles',
        name: 'roles',
        component: () => import('../modules/iam/views/RoleListView.vue'),
        meta: { title: '角色管理', permission: 'iam:role:view' }
      },
      {
        path: '/permissions',
        name: 'permissions',
        component: () => import('../modules/iam/views/PermissionListView.vue'),
        meta: { title: '权限管理', permission: 'iam:permission:view' }
      },
      {
        path: '/portal/articles',
        name: 'articles',
        component: () => import('../modules/portal/views/ArticleListView.vue'),
        meta: { title: '文章管理', permission: 'portal:view' }
      },
      {
        path: '/portal/banners',
        name: 'banners',
        component: () => import('../modules/portal/views/BannerListView.vue'),
        meta: { title: '轮播管理', permission: 'portal:view' }
      },
      {
        path: '/portal/notices',
        name: 'notices',
        component: () => import('../modules/portal/views/NoticeListView.vue'),
        meta: { title: '公告管理', permission: 'portal:view' }
      },
      {
        path: '/campus-life/errands',
        name: 'errands',
        component: () => import('../modules/campus-life/views/ErrandListView.vue'),
        meta: { title: '跑腿服务', permission: 'campus_life:view' }
      },
      {
        path: '/campus-life/meetups',
        name: 'meetups',
        component: () => import('../modules/campus-life/views/MeetupListView.vue'),
        meta: { title: '组局活动', permission: 'campus_life:view' }
      },
      {
        path: '/campus-life/listings',
        name: 'listings',
        component: () => import('../modules/campus-life/views/ListingListView.vue'),
        meta: { title: '二手交易', permission: 'campus_life:view' }
      },
      {
        path: '/campus-life/resources',
        name: 'resources',
        component: () => import('../modules/campus-life/views/ResourceListView.vue'),
        meta: { title: '资料共享', permission: 'campus_life:view' }
      },
      {
        path: '/campus-life/lost-items',
        name: 'lostItems',
        component: () => import('../modules/campus-life/views/LostItemListView.vue'),
        meta: { title: '失物招领', permission: 'campus_life:view' }
      },
      {
        path: '/moderation/pending',
        name: 'moderationPending',
        component: () => import('../modules/moderation/views/PendingListView.vue'),
        meta: { title: '待审核', permission: 'campus_life:moderate' }
      },
      {
        path: '/moderation/history',
        name: 'moderationHistory',
        component: () => import('../modules/moderation/views/HistoryView.vue'),
        meta: { title: '审核历史', permission: 'campus_life:moderate' }
      },
      {
        path: '/analytics/audit-logs',
        name: 'auditLogs',
        component: () => import('../modules/analytics/views/AuditLogView.vue'),
        meta: { title: '审计日志', permission: 'analytics:view' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'notFound',
    component: () => import('../app/views/NotFoundView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

const protectedPaths = [
  { path: '/', permission: 'analytics:view' },
  { path: '/users', permission: 'iam:user:view' },
  { path: '/roles', permission: 'iam:role:view' },
  { path: '/permissions', permission: 'iam:permission:view' },
  { path: '/portal/articles', permission: 'portal:view' },
  { path: '/portal/banners', permission: 'portal:view' },
  { path: '/portal/notices', permission: 'portal:view' },
  { path: '/campus-life/errands', permission: 'campus_life:view' },
  { path: '/campus-life/meetups', permission: 'campus_life:view' },
  { path: '/campus-life/listings', permission: 'campus_life:view' },
  { path: '/campus-life/resources', permission: 'campus_life:view' },
  { path: '/campus-life/lost-items', permission: 'campus_life:view' },
  { path: '/moderation/pending', permission: 'campus_life:moderate' },
  { path: '/moderation/history', permission: 'campus_life:moderate' },
  { path: '/analytics/audit-logs', permission: 'analytics:view' }
]

function firstAccessiblePath() {
  const authStore = useAuthStore()
  const matched = protectedPaths.find((item) => authStore.hasPermission(item.permission))
  return matched?.path || '/login'
}

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next('/login')
    return
  }

  if (to.path === '/login' && authStore.isLoggedIn) {
    next(firstAccessiblePath())
    return
  }

  if (to.meta.requiresAuth && !authStore.user) {
    authStore.logout()
    next('/login')
    return
  }

  const requiredPermission = typeof to.meta.permission === 'string' ? to.meta.permission : ''
  if (requiredPermission && !authStore.hasPermission(requiredPermission)) {
    next(firstAccessiblePath())
    return
  }

  next()
})

export default router
