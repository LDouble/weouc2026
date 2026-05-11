<template>
  <a-layout class="min-h-screen">
    <a-layout-sider v-model:collapsed="collapsed" width="200" class="bg-slate-900">
      <div class="h-16 flex items-center justify-center">
        <span class="text-white font-bold text-lg">校园管理后台</span>
      </div>
      <a-menu mode="inline" theme="dark" class="mt-4" :selected-keys="[currentKey]">
        <a-menu-item key="/" @click="navigate('/')">
          <component :is="icons.Dashboard" />
          <span>仪表盘</span>
        </a-menu-item>
        
        <a-sub-menu key="iam">
          <template #title>
            <component :is="icons.User" />
            <span>用户权限</span>
          </template>
          <a-menu-item key="/users" @click="navigate('/users')">用户管理</a-menu-item>
          <a-menu-item key="/roles" @click="navigate('/roles')">角色管理</a-menu-item>
          <a-menu-item key="/permissions" @click="navigate('/permissions')">权限管理</a-menu-item>
        </a-sub-menu>
        
        <a-sub-menu key="portal">
          <template #title>
            <component :is="icons.FileText" />
            <span>内容管理</span>
          </template>
          <a-menu-item key="/portal/articles" @click="navigate('/portal/articles')">文章管理</a-menu-item>
          <a-menu-item key="/portal/banners" @click="navigate('/portal/banners')">轮播管理</a-menu-item>
          <a-menu-item key="/portal/notices" @click="navigate('/portal/notices')">公告管理</a-menu-item>
        </a-sub-menu>
        
        <a-sub-menu key="campus-life">
          <template #title>
            <component :is="icons.Home" />
            <span>校园生活</span>
          </template>
          <a-menu-item key="/campus-life/errands" @click="navigate('/campus-life/errands')">跑腿服务</a-menu-item>
          <a-menu-item key="/campus-life/meetups" @click="navigate('/campus-life/meetups')">组局活动</a-menu-item>
          <a-menu-item key="/campus-life/listings" @click="navigate('/campus-life/listings')">二手交易</a-menu-item>
          <a-menu-item key="/campus-life/resources" @click="navigate('/campus-life/resources')">资料共享</a-menu-item>
          <a-menu-item key="/campus-life/lost-items" @click="navigate('/campus-life/lost-items')">失物招领</a-menu-item>
        </a-sub-menu>
        
        <a-sub-menu key="moderation">
          <template #title>
            <component :is="icons.CheckCircle" />
            <span>内容审核</span>
          </template>
          <a-menu-item key="/moderation/pending" @click="navigate('/moderation/pending')">待审核</a-menu-item>
          <a-menu-item key="/moderation/history" @click="navigate('/moderation/history')">审核历史</a-menu-item>
        </a-sub-menu>
        
        <a-sub-menu key="analytics">
          <template #title>
            <component :is="icons.BarChart" />
            <span>数据分析</span>
          </template>
          <a-menu-item key="/analytics/audit-logs" @click="navigate('/analytics/audit-logs')">审计日志</a-menu-item>
        </a-sub-menu>
      </a-menu>
    </a-layout-sider>
    
    <a-layout>
      <a-layout-header class="bg-white border-b flex items-center justify-between px-6">
        <div class="flex items-center gap-4">
          <a-button type="text" @click="collapsed = !collapsed">
            <component :is="collapsed ? icons.MenuUnfold : icons.MenuFold" />
          </a-button>
          <span class="text-gray-800 font-medium">{{ currentTitle }}</span>
        </div>
        <div class="flex items-center gap-4">
          <a-dropdown>
            <a-button type="text" class="flex items-center gap-2">
              <span>{{ authStore.user?.username || '管理员' }}</span>
              <component :is="icons.Down" />
            </a-button>
            <template #overlay>
              <a-menu>
                <a-menu-item @click="handleLogout">
                  <component :is="icons.Logout" />
                  <span>退出登录</span>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>
      
      <a-layout-content class="p-6 bg-gray-50">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { 
  DashboardOutlined as Dashboard,
  UserOutlined as User,
  FileTextOutlined as FileText,
  HomeOutlined as Home,
  CheckCircleOutlined as CheckCircle,
  BarChartOutlined as BarChart,
  MenuUnfoldOutlined as MenuUnfold,
  MenuFoldOutlined as MenuFold,
  DownOutlined as Down,
  LogoutOutlined as Logout
} from '@ant-design/icons-vue'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const collapsed = ref(false)

const icons = {
  Dashboard,
  User,
  FileText,
  Home,
  CheckCircle,
  BarChart,
  MenuUnfold,
  MenuFold,
  Down,
  Logout
}

const currentKey = computed(() => router.currentRoute.value.path)

const titleMap: Record<string, string> = {
  '/': '仪表盘',
  '/users': '用户管理',
  '/roles': '角色管理',
  '/permissions': '权限管理',
  '/portal/articles': '文章管理',
  '/portal/banners': '轮播管理',
  '/portal/notices': '公告管理',
  '/campus-life/errands': '跑腿服务',
  '/campus-life/meetups': '组局活动',
  '/campus-life/listings': '二手交易',
  '/campus-life/resources': '资料共享',
  '/campus-life/lost-items': '失物招领',
  '/moderation/pending': '待审核',
  '/moderation/history': '审核历史',
  '/analytics/audit-logs': '审计日志'
}

const currentTitle = computed(() => titleMap[router.currentRoute.value.path] || '')

function navigate(path: string) {
  router.push(path)
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>