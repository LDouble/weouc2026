<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-800">权限管理</h2>
    </div>
    
    <a-table :columns="columns" :data-source="permissions" :pagination="pagination">
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
import { ref, computed } from 'vue'

interface Permission {
  id: string
  code: string
  name: string
  description: string
  category: string
}

const permissions = ref<Permission[]>([
  { id: '1', code: 'iam:user:view', name: '查看用户', description: '查看用户列表和详情', category: 'iam' },
  { id: '2', code: 'iam:user:manage', name: '管理用户', description: '创建、编辑、删除用户', category: 'iam' },
  { id: '3', code: 'iam:role:view', name: '查看角色', description: '查看角色列表', category: 'iam' },
  { id: '4', code: 'iam:role:manage', name: '管理角色', description: '创建、编辑、删除角色', category: 'iam' },
  { id: '5', code: 'iam:permission:view', name: '查看权限', description: '查看权限列表', category: 'iam' },
  { id: '6', code: 'portal:publish', name: '发布内容', description: '发布文章、公告、轮播', category: 'portal' },
  { id: '7', code: 'portal:view', name: '查看内容', description: '查看内容列表', category: 'portal' },
  { id: '8', code: 'campus_life:moderate', name: '审核校园生活', description: '审核跑腿、组局、二手等内容', category: 'campus_life' },
  { id: '9', code: 'campus_life:view', name: '查看校园生活', description: '查看校园生活内容', category: 'campus_life' },
  { id: '10', code: 'moderation:review', name: '审核内容', description: '审核待审核内容', category: 'moderation' },
  { id: '11', code: 'analytics:view', name: '查看数据', description: '查看统计数据和审计日志', category: 'analytics' }
])

const pagination = computed(() => ({
  total: permissions.value.length,
  showSizeChanger: true,
  showQuickJumper: true
}))

const columns = [
  { title: '权限编码', dataIndex: 'code', key: 'code' },
  { title: '权限名称', dataIndex: 'name', key: 'name' },
  { title: '描述', dataIndex: 'description', key: 'description' },
  { title: '分类', key: 'category' }
]

function getCategoryColor(category: string): string {
  const colors: Record<string, string> = {
    iam: 'blue',
    portal: 'green',
    campus_life: 'purple',
    moderation: 'orange',
    analytics: 'cyan'
  }
  return colors[category] || 'gray'
}

function getCategoryName(category: string): string {
  const names: Record<string, string> = {
    iam: '用户权限',
    portal: '内容管理',
    campus_life: '校园生活',
    moderation: '审核管理',
    analytics: '数据分析'
  }
  return names[category] || category
}
</script>