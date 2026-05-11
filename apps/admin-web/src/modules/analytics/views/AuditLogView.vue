<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-800">审计日志</h2>
      <div class="flex items-center gap-4">
        <a-input v-model:value="searchText" placeholder="搜索操作或目标" style="width: 200px" @change="fetchLogs" />
        <a-select v-model:value="actionFilter" placeholder="筛选操作类型" @change="fetchLogs">
          <a-select-option value="all">全部</a-select-option>
          <a-select-option value="login">登录</a-select-option>
          <a-select-option value="create">创建</a-select-option>
          <a-select-option value="update">更新</a-select-option>
          <a-select-option value="delete">删除</a-select-option>
          <a-select-option value="review">审核</a-select-option>
        </a-select>
      </div>
    </div>
    
    <a-table :columns="columns" :data-source="auditLogs" :pagination="tablePagination" @change="handleTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a-tag :color="getActionColor(record.action)">
            {{ getActionName(record.action) }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'ip'">
          <span class="text-sm text-gray-500">{{ record.ip }}</span>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { analyticsApi } from '@/api'

interface AuditLog {
  id: string
  action: string
  target: string
  targetType: string
  operator: string
  ip: string
  time: string
  details: string
}

const searchText = ref('')
const actionFilter = ref('all')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const auditLogs = ref<AuditLog[]>([])

const tablePagination = computed(() => ({
  current: currentPage.value,
  pageSize: pageSize.value,
  total: total.value,
  showSizeChanger: true,
  showQuickJumper: true
}))

const columns = [
  { title: '操作类型', key: 'action' },
  { title: '操作目标', dataIndex: 'target', key: 'target' },
  { title: '目标类型', dataIndex: 'targetType', key: 'targetType' },
  { title: '操作人', dataIndex: 'operator', key: 'operator' },
  { title: 'IP地址', key: 'ip' },
  { title: '操作时间', dataIndex: 'time', key: 'time' },
  { title: '详情', dataIndex: 'details', key: 'details' }
]

function getActionColor(action: string): string {
  const colors: Record<string, string> = {
    login: 'blue',
    create: 'green',
    update: 'orange',
    delete: 'red',
    review: 'purple'
  }
  return colors[action] || 'gray'
}

function getActionName(action: string): string {
  const names: Record<string, string> = {
    login: '登录',
    create: '创建',
    update: '更新',
    delete: '删除',
    review: '审核'
  }
  return names[action] || action
}

const fetchLogs = async () => {
  try {
    const params: any = {
      page: currentPage.value,
      pageSize: pageSize.value
    }
    if (actionFilter.value !== 'all') {
      params.action = actionFilter.value
    }
    if (searchText.value) {
      params.keyword = searchText.value
    }
    
    const response = await analyticsApi.listAuditLogs(params)
    if (response.data && response.data.list) {
      auditLogs.value = response.data.list.map((item: any) => ({
        id: item.id,
        action: item.action || 'unknown',
        target: item.resource_id || '-',
        targetType: item.resource_type || '-',
        operator: item.actor_id || 'unknown',
        ip: item.ip_address || '-',
        time: item.created_at ? new Date(item.created_at).toLocaleString('zh-CN') : '-',
        details: item.description || '-'
      }))
      total.value = response.data.total
    }
  } catch (error) {
    console.error('Failed to fetch audit logs:', error)
    message.error('获取审计日志失败')
  }
}

const handleTableChange = (pagination: any) => {
  currentPage.value = pagination.current
  pageSize.value = pagination.pageSize
  fetchLogs()
}

onMounted(() => {
  fetchLogs()
})
</script>