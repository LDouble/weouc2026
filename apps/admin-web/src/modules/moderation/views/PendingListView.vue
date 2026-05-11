<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-800">待审核内容</h2>
      <a-select v-model:value="contentType" placeholder="筛选类型" @change="fetchPendingItems">
        <a-select-option value="all">全部</a-select-option>
        <a-select-option value="errand">跑腿</a-select-option>
        <a-select-option value="meetup">组局</a-select-option>
        <a-select-option value="market">二手</a-select-option>
        <a-select-option value="resource">资料</a-select-option>
        <a-select-option value="lostFound">失物招领</a-select-option>
        <a-select-option value="carpool">拼车</a-select-option>
      </a-select>
    </div>
    
    <a-table :columns="columns" :data-source="pendingItems" :pagination="tablePagination" @change="handleTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'contentType'">
          <a-tag :color="getTypeColor(record.contentType)">
            {{ getTypeName(record.contentType) }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" type="primary" @click="approveItem(record)">
              通过
            </a-button>
            <a-button size="small" danger @click="rejectItem(record)">
              拒绝
            </a-button>
            <a-button size="small" @click="viewDetail(record)">
              详情
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { campusLifeApi } from '@/api'

interface PendingItem {
  id: string
  content_id: string
  content_type: string
  title: string
  content: string
  publisher: string
  created_at: string
}

const contentType = ref('all')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const pendingItems = ref<PendingItem[]>([])

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
  { title: '内容摘要', dataIndex: 'content', key: 'content' },
  { title: '发布者', dataIndex: 'publisher', key: 'publisher' },
  { title: '提交时间', dataIndex: 'created_at', key: 'created_at' },
  { title: '操作', key: 'actions' }
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

const fetchPendingItems = async () => {
  try {
    const params: any = {
      page: currentPage.value,
      pageSize: pageSize.value,
      review_status: 'reviewing'
    }
    if (contentType.value !== 'all') {
      params.content_type = contentType.value
    }
    
    const response = await campusLifeApi.listReviews(params)
    if (response.data) {
      pendingItems.value = response.data.list.map((item: any) => ({
        id: item.id,
        content_id: item.content_id,
        content_type: item.content_type,
        title: item.title,
        content: item.desc,
        publisher: item.publisher || '-',
        created_at: item.created_at ? new Date(item.created_at).toLocaleString('zh-CN') : '-'
      }))
      total.value = response.data.total
    }
  } catch (error) {
    console.error('Failed to fetch pending items:', error)
    message.error('获取待审核列表失败')
  }
}

const handleTableChange = (pagination: any) => {
  currentPage.value = pagination.current
  pageSize.value = pagination.pageSize
  fetchPendingItems()
}

const approveItem = async (item: PendingItem) => {
  try {
    await campusLifeApi.updateReviewStatus({
      content_type: item.content_type as any,
      content_id: item.content_id,
      review_status: 'published'
    })
    message.success('审核通过')
    fetchPendingItems()
  } catch (error) {
    console.error('Failed to approve item:', error)
    message.error('审核失败')
  }
}

const rejectItem = async (item: PendingItem) => {
  try {
    await campusLifeApi.updateReviewStatus({
      content_type: item.content_type as any,
      content_id: item.content_id,
      review_status: 'rejected',
      reason: '不符合审核要求'
    })
    message.success('已拒绝')
    fetchPendingItems()
  } catch (error) {
    console.error('Failed to reject item:', error)
    message.error('操作失败')
  }
}

const viewDetail = (item: PendingItem) => {
  message.info(`查看详情: ${item.title}`)
}

onMounted(() => {
  fetchPendingItems()
})
</script>
