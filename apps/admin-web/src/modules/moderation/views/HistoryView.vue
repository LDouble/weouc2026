<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-6">
      <div>
        <h2 class="text-lg font-semibold text-gray-800">审核历史</h2>
        <p class="text-sm text-gray-500 mt-1">当前页面展示已发布、已拒绝和已下线内容的真实审核结果。</p>
      </div>
      <a-select v-model:value="statusFilter" placeholder="筛选状态" @change="handleStatusChange" class="w-44">
        <a-select-option value="all">全部</a-select-option>
        <a-select-option value="published">已通过</a-select-option>
        <a-select-option value="rejected">已拒绝</a-select-option>
        <a-select-option value="offline">已下线</a-select-option>
      </a-select>
    </div>

    <a-table :columns="columns" :data-source="historyItems" :pagination="tablePagination" @change="handleTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'contentType'">
          <a-tag :color="getTypeColor(record.content_type)">
            {{ getTypeName(record.content_type) }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'reviewStatus'">
          <a-tag :color="getReviewStatusColor(record.review_status)">
            {{ getReviewStatusName(record.review_status) }}
          </a-tag>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { message } from 'ant-design-vue'
import { campusLifeApi } from '@/api'

interface HistoryItem {
  content_id: string
  content_type: string
  title: string
  desc: string
  publisher: string
  review_status: string
  created_at: string
}

const statusFilter = ref<'all' | 'published' | 'rejected' | 'offline'>('all')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const historyItems = ref<HistoryItem[]>([])

const tablePagination = computed(() => ({
  current: currentPage.value,
  pageSize: pageSize.value,
  total: total.value,
  showSizeChanger: true,
  showQuickJumper: true
}))

const columns = [
  { title: '标题', dataIndex: 'title', key: 'title' },
  { title: '类型', key: 'contentType' },
  { title: '审核结果', key: 'reviewStatus' },
  { title: '内容摘要', dataIndex: 'desc', key: 'desc' },
  { title: '发布者', dataIndex: 'publisher', key: 'publisher' },
  { title: '提交时间', dataIndex: 'created_at', key: 'created_at' }
]

function getTypeColor(type: string): string {
  const colors: Record<string, string> = {
    errand: 'blue',
    meetup: 'green',
    market: 'purple',
    resource: 'cyan',
    lostFound: 'orange',
    carpool: 'magenta'
  }
  return colors[type] || 'gray'
}

function getTypeName(type: string): string {
  const names: Record<string, string> = {
    errand: '跑腿',
    meetup: '组局',
    market: '二手',
    resource: '资料',
    lostFound: '失物招领',
    carpool: '拼车'
  }
  return names[type] || type
}

function getReviewStatusColor(status: string): string {
  const colors: Record<string, string> = {
    published: 'green',
    rejected: 'red',
    offline: 'gray'
  }
  return colors[status] || 'blue'
}

function getReviewStatusName(status: string): string {
  const names: Record<string, string> = {
    published: '已通过',
    rejected: '已拒绝',
    offline: '已下线'
  }
  return names[status] || status
}

async function fetchHistoryItems() {
  try {
    if (statusFilter.value === 'all') {
      const responses = await Promise.all([
        campusLifeApi.listReviews({ page: 1, pageSize: 200, review_status: 'published' }),
        campusLifeApi.listReviews({ page: 1, pageSize: 200, review_status: 'rejected' }),
        campusLifeApi.listReviews({ page: 1, pageSize: 200, review_status: 'offline' })
      ])

      const merged = responses
        .flatMap((response) => response.data?.list || [])
        .sort((a: any, b: any) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())

      total.value = merged.length
      const start = (currentPage.value - 1) * pageSize.value
      const end = start + pageSize.value
      historyItems.value = merged.slice(start, end).map((item: any) => ({
        content_id: item.content_id,
        content_type: item.content_type,
        title: item.title,
        desc: item.desc || '',
        publisher: item.publisher || '-',
        review_status: item.review_status,
        created_at: item.created_at ? new Date(item.created_at).toLocaleString('zh-CN') : '-'
      }))
      return
    }

    const response = await campusLifeApi.listReviews({
      page: currentPage.value,
      pageSize: pageSize.value,
      review_status: statusFilter.value as 'published' | 'rejected' | 'offline'
    })

    historyItems.value = (response.data?.list || []).map((item: any) => ({
      content_id: item.content_id,
      content_type: item.content_type,
      title: item.title,
      desc: item.desc || '',
      publisher: item.publisher || '-',
      review_status: item.review_status,
      created_at: item.created_at ? new Date(item.created_at).toLocaleString('zh-CN') : '-'
    }))
    total.value = Number(response.data?.total || historyItems.value.length)
  } catch (error) {
    console.error('Failed to fetch review history:', error)
    message.error('获取审核历史失败')
  }
}

function handleTableChange(pagination: any) {
  currentPage.value = pagination.current
  pageSize.value = pagination.pageSize
  fetchHistoryItems()
}

function handleStatusChange() {
  currentPage.value = 1
  fetchHistoryItems()
}

onMounted(() => {
  fetchHistoryItems()
})
</script>
