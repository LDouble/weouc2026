<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-lg font-semibold text-gray-800">角色管理</h2>
      <a-button @click="refreshView">刷新</a-button>
    </div>

    <a-alert
      class="mb-4"
      type="info"
      show-icon
      message="后端当前未开放角色管理写接口，页面展示当前管理员角色与权限快照。"
    />

    <a-table
      row-key="id"
      :columns="columns"
      :data-source="rows"
      :pagination="false"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'permissions'">
          <a-space wrap>
            <a-tag v-for="permission in record.permissions" :key="permission" color="purple">
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

interface SessionRoleRow {
  id: string
  name: string
  description: string
  permissions: string[]
}

const authStore = useAuthStore()

const columns = [
  { title: '角色编码', dataIndex: 'id', key: 'id', width: 240 },
  { title: '角色名称', dataIndex: 'name', key: 'name', width: 220 },
  { title: '说明', dataIndex: 'description', key: 'description', width: 280 },
  { title: '权限集', key: 'permissions' },
]

const rows = computed<SessionRoleRow[]>(() => {
  if (!authStore.user) {
    return []
  }

  return [{
    id: authStore.user.role || 'unknown',
    name: authStore.user.role || 'unknown',
    description: '来源于当前管理员登录会话',
    permissions: authStore.user.permissions || [],
  }]
})

function refreshView() {
  window.location.reload()
}
</script>
