<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-800">审核历史</h2>
      <a-select v-model:value="statusFilter" placeholder="筛选状态">
        <a-select-option value="all">全部</a-select-option>
        <a-select-option value="approved">已通过</a-select-option>
        <a-select-option value="rejected">已拒绝</a-select-option>
      </a-select>
    </div>
    
    <a-table :columns="columns" :data-source="historyItems" :pagination="pagination">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'contentType'">
          <a-tag :color="getTypeColor(record.contentType)">
            {{ getTypeName(record.contentType) }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'result'">
          <a-tag :color="record.result === 'approved' ? 'green' : 'red'">
            {{ record.result === 'approved' ? '通过' : '拒绝' }}
          </a-tag>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface HistoryItem {
  id: string
  title: string
  contentType: string
  result: string
  reviewer: string
  reviewTime: string
  reason: string
}

const statusFilter = ref('all')

const historyItems = ref<HistoryItem[]>([
  { id: '1', title: '帮取快递', contentType: 'errand', result: 'approved', reviewer: 'admin', reviewTime: '2024-01-15 16:00', reason: '内容合规' },
  { id: '2', title: '周末聚餐', contentType: 'meetup', result: 'approved', reviewer: 'admin', reviewTime: '2024-01-15 15:30', reason: '内容合规' },
  { id: '3', title: '违规广告', contentType: 'listing', result: 'rejected', reviewer: 'admin', reviewTime: '2024-01-15 14:00', reason: '内容不符合规定' },
  { id: '4', title: '考研资料', contentType: 'resource', result: 'approved', reviewer: 'admin', reviewTime: '2024-01-15 13:00', reason: '内容合规' },
  { id: '5', title: '校园卡', contentType: 'lost_item', result: 'approved', reviewer: 'admin', reviewTime: '2024-01-15 12:00', reason: '内容合规' }
])

const filteredItems = computed(() => {
  if (statusFilter.value === 'all') return historyItems.value
  return historyItems.value.filter(item => item.result === statusFilter.value)
})

const pagination = computed(() => ({
  total: filteredItems.value.length,
  showSizeChanger: true,
  showQuickJumper: true
}))

const columns = [
  { title: '标题', dataIndex: 'title', key: 'title' },
  { title: '类型', key: 'contentType' },
  { title: '审核结果', key: 'result' },
  { title: '审核人', dataIndex: 'reviewer', key: 'reviewer' },
  { title: '审核时间', dataIndex: 'reviewTime', key: 'reviewTime' },
  { title: '审核理由', dataIndex: 'reason', key: 'reason' }
]

function getTypeColor(type: string): string {
  const colors: Record<string, string> = {
    errand: 'blue',
    meetup: 'green',
    listing: 'purple',
    resource: 'cyan',
    lost_item: 'orange'
  }
  return colors[type] || 'gray'
}

function getTypeName(type: string): string {
  const names: Record<string, string> = {
    errand: '跑腿',
    meetup: '组局',
    listing: '二手',
    resource: '资料',
    lost_item: '失物招领'
  }
  return names[type] || type
}
</script>