import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

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
        component: () => import('../modules/analytics/views/DashboardView.vue')
      },
      {
        path: '/users',
        name: 'users',
        component: () => import('../modules/iam/views/UserListView.vue')
      },
      {
        path: '/roles',
        name: 'roles',
        component: () => import('../modules/iam/views/RoleListView.vue')
      },
      {
        path: '/permissions',
        name: 'permissions',
        component: () => import('../modules/iam/views/PermissionListView.vue')
      },
      {
        path: '/portal/articles',
        name: 'articles',
        component: () => import('../modules/portal/views/ArticleListView.vue')
      },
      {
        path: '/portal/banners',
        name: 'banners',
        component: () => import('../modules/portal/views/BannerListView.vue')
      },
      {
        path: '/portal/notices',
        name: 'notices',
        component: () => import('../modules/portal/views/NoticeListView.vue')
      },
      {
        path: '/campus-life/errands',
        name: 'errands',
        component: () => import('../modules/campus-life/views/ErrandListView.vue')
      },
      {
        path: '/campus-life/meetups',
        name: 'meetups',
        component: () => import('../modules/campus-life/views/MeetupListView.vue')
      },
      {
        path: '/campus-life/listings',
        name: 'listings',
        component: () => import('../modules/campus-life/views/ListingListView.vue')
      },
      {
        path: '/campus-life/resources',
        name: 'resources',
        component: () => import('../modules/campus-life/views/ResourceListView.vue')
      },
      {
        path: '/campus-life/lost-items',
        name: 'lostItems',
        component: () => import('../modules/campus-life/views/LostItemListView.vue')
      },
      {
        path: '/moderation/pending',
        name: 'moderationPending',
        component: () => import('../modules/moderation/views/PendingListView.vue')
      },
      {
        path: '/moderation/history',
        name: 'moderationHistory',
        component: () => import('../modules/moderation/views/HistoryView.vue')
      },
      {
        path: '/analytics/audit-logs',
        name: 'auditLogs',
        component: () => import('../modules/analytics/views/AuditLogView.vue')
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

router.beforeEach((to, _from, next) => {
  const isLoggedIn = localStorage.getItem('adminToken') !== null
  if (to.meta.requiresAuth && !isLoggedIn) {
    next('/login')
  } else if (to.path === '/login' && isLoggedIn) {
    next('/')
  } else {
    next()
  }
})

export default router