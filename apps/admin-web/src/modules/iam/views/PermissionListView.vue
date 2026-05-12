<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-lg font-semibold text-gray-800">权限管理</h2>
      <a-button @click="refreshView">刷新</a-button>
    </div>

    <a-alert
      class="mb-4"
      type="info"
      show-icon
      message="后端当前未开放权限目录接口，页面展示当前管理员登录会话中的权限快照。"
    />

    <a-table
      row-key="id"
      :columns="columns"
      :data-source="permissionRows"
      :pagination="false"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'category'">
          <a-tag :color="getCategoryColor(record.category)">
            {{ getCategoryName(record.category) }}
          </a-tag>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '../../../stores/auth'

interface PermissionRow {
  id: string
  code: string
  name: string
  description: string
  category: string
}

const authStore = useAuthStore()

const columns = [
  { title: '权限编码', dataIndex: 'code', key: 'code', width: 300 },
  { title: '权限名称', dataIndex: 'name', key: 'name', width: 220 },
  { title: '描述', dataIndex: 'description', key: 'description' },
  { title: '分类', key: 'category', width: 150 },
]

const permissionRows = computed<PermissionRow[]>(() => {
  const permissions = authStore.user?.permissions || []
  return permissions.map((permission: string, index: number) => {
    const segments = permission.split(':')
    const category = segments[0] || 'unknown'
    const action = segments[1] || 'view'
    const displayAction = action.replace(/_/g, ' ')

    return {
      id: `${permission}-${index}`,
      code: permission,
      name: `${category} ${displayAction}`.trim(),
      description: '来源于当前管理员登录会话',
      category,
    }
  })
})

function getCategoryColor(category: string): string {
  const colors: Record<string, string> = {
    iam: 'blue',
    portal: 'green',
    campus_life: 'purple',
    moderation: 'orange',
    analytics: 'cyan',
    notification: 'gold',
  }
  return colors[category] || 'default'
}

function getCategoryName(category: string): string {
  const names: Record<string, string> = {
    iam: '身份权限',
    portal: '门户内容',
    campus_life: '校园生活',
    moderation: '审核管理',
    analytics: '数据分析',
    notification: '站内通知',
  }
  return names[category] || category
}

function refreshView() {
  window.location.reload()
}
</script>
