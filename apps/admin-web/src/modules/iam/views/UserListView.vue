<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-lg font-semibold text-gray-800">用户管理</h2>
      <a-button @click="refreshView">刷新</a-button>
    </div>

    <a-alert
      class="mb-4"
      type="info"
      show-icon
      message="后端当前未开放管理员用户目录接口，页面仅展示当前登录管理员会话。"
    />

    <a-table
      row-key="id"
      :columns="columns"
      :data-source="rows"
      :pagination="false"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag color="green">活跃</a-tag>
        </template>
        <template v-else-if="column.key === 'permissions'">
          <a-space wrap>
            <a-tag v-for="permission in record.permissions" :key="permission" color="blue">
              {{ permission }}
            </a-tag>
          </a-space>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '../../../stores/auth'

interface SessionUserRow {
  id: string
  username: string
  role: string
  status: 'active'
  permissions: string[]
}

const authStore = useAuthStore()

const columns = [
  { title: '用户 ID', dataIndex: 'id', key: 'id', width: 180 },
  { title: '用户名', dataIndex: 'username', key: 'username', width: 220 },
  { title: '角色', dataIndex: 'role', key: 'role', width: 220 },
  { title: '状态', key: 'status', width: 120 },
  { title: '权限集', key: 'permissions' },
]

const rows = computed<SessionUserRow[]>(() => {
  if (!authStore.user) {
    return []
  }

  return [{
    id: authStore.user.id,
    username: authStore.user.username,
    role: authStore.user.role,
    status: 'active',
    permissions: authStore.user.permissions || [],
  }]
})

function refreshView() {
  window.location.reload()
}
</script>
