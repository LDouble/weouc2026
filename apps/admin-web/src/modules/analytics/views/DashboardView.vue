<template>
  <div class="space-y-6">
    <div class="grid grid-cols-4 gap-4">
      <a-card class="bg-blue-500 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-blue-100 text-sm">总用户数</p>
            <p class="text-2xl font-bold">{{ stats.user_count }}</p>
          </div>
          <component :is="icons.User" class="text-4xl text-blue-200" />
        </div>
      </a-card>
      
      <a-card class="bg-green-500 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-green-100 text-sm">活跃用户</p>
            <p class="text-2xl font-bold">{{ stats.active_users }}</p>
          </div>
          <component :is="icons.Active" class="text-4xl text-green-200" />
        </div>
      </a-card>
      
      <a-card class="bg-purple-500 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-purple-100 text-sm">待审核</p>
            <p class="text-2xl font-bold">{{ stats.pending_reviews }}</p>
          </div>
          <component :is="icons.Review" class="text-4xl text-purple-200" />
        </div>
      </a-card>
      
      <a-card class="bg-orange-500 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-orange-100 text-sm">总发布</p>
            <p class="text-2xl font-bold">{{ stats.total_posts }}</p>
          </div>
          <component :is="icons.Post" class="text-4xl text-orange-200" />
        </div>
      </a-card>
    </div>
    
    <div class="grid grid-cols-2 gap-6">
      <a-card title="内容分类统计">
        <a-table :columns="categoryColumns" :data-source="categoryStats" :pagination="false" />
      </a-card>
      
      <a-card title="审核统计">
        <a-table :columns="reviewColumns" :data-source="reviewStats" :pagination="false" />
      </a-card>
    </div>
    
    <a-card title="最近操作记录">
      <a-table :columns="logColumns" :data-source="recentLogs" :pagination="pagination" />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { UserOutlined as User, FileDoneOutlined as Review, FileAddOutlined as Post, ClockCircleOutlined as Active } from '@ant-design/icons-vue'
import { analyticsApi } from '@/api'

const icons = { User, Active, Review, Post }

const stats = ref({
  user_count: 0,
  active_users: 0,
  total_posts: 0,
  pending_reviews: 0,
  notifications_sent: 0
})

const categoryStats = ref([
  { name: '跑腿服务', count: 0, percentage: '0%' },
  { name: '二手交易', count: 0, percentage: '0%' },
  { name: '组局活动', count: 0, percentage: '0%' },
  { name: '资料共享', count: 0, percentage: '0%' },
  { name: '失物招领', count: 0, percentage: '0%' }
])

const reviewStats = ref([
  { name: '已通过', count: 0, percentage: '0%' },
  { name: '已拒绝', count: 0, percentage: '0%' }
])

const recentLogs = ref<any[]>([])

const categoryColumns = [
  { title: '分类', dataIndex: 'name', key: 'name' },
  { title: '数量', dataIndex: 'count', key: 'count' },
  { title: '占比', dataIndex: 'percentage', key: 'percentage' }
]

const reviewColumns = [
  { title: '状态', dataIndex: 'name', key: 'name' },
  { title: '数量', dataIndex: 'count', key: 'count' },
  { title: '占比', dataIndex: 'percentage', key: 'percentage' }
]

const logColumns = [
  { title: '操作', dataIndex: 'action', key: 'action' },
  { title: '目标', dataIndex: 'target', key: 'target' },
  { title: '操作人', dataIndex: 'operator', key: 'operator' },
  { title: '时间', dataIndex: 'time', key: 'time' }
]

const pagination = computed(() => ({
  total: recentLogs.value.length,
  showSizeChanger: true,
  showQuickJumper: true
}))

const fetchDashboardData = async () => {
  try {
    const response = await analyticsApi.getDashboard()
    if (response.data) {
      stats.value = {
        user_count: response.data.user_count || 0,
        active_users: response.data.active_users || 0,
        total_posts: response.data.total_posts || 0,
        pending_reviews: response.data.pending_reviews || 0,
        notifications_sent: response.data.notifications_sent || 0
      }
      
      if (response.data.category_stats) {
        categoryStats.value = response.data.category_stats.map((item: any) => ({
          name: item.name,
          count: item.count,
          percentage: item.percentage
        }))
      }
      
      if (response.data.review_stats) {
        reviewStats.value = response.data.review_stats.map((item: any) => ({
          name: item.name,
          count: item.count,
          percentage: item.percentage
        }))
      }
    }
  } catch (error) {
    console.error('Failed to fetch dashboard data:', error)
  }
}

const fetchAuditLogs = async () => {
  try {
    const response = await analyticsApi.listAuditLogs({ page: 1, pageSize: 10 })
    if (response.data && response.data.list) {
      recentLogs.value = response.data.list.map((log: any) => ({
        id: log.id,
        action: log.action || '未知操作',
        target: log.resource_id || '-',
        operator: log.actor_id || 'unknown',
        time: log.created_at ? new Date(log.created_at).toLocaleString('zh-CN') : '-'
      }))
    }
  } catch (error) {
    console.error('Failed to fetch audit logs:', error)
  }
}

onMounted(() => {
  fetchDashboardData()
  fetchAuditLogs()
})
</script>