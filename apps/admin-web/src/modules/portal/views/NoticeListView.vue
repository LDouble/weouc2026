<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-800">公告管理</h2>
      <a-button type="primary" @click="showModal = true">
        <component :is="icons.Plus" />
        发布公告
      </a-button>
    </div>
    
    <a-table :columns="columns" :data-source="notices" :pagination="tablePagination" @change="handleTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'pinned'">
          <a-tag v-if="record.pinned" color="gold">置顶</a-tag>
        </template>
        <template v-else-if="column.key === 'status'">
          <a-tag :color="record.status === 'published' ? 'green' : record.status === 'reviewing' ? 'orange' : 'gray'">
            {{ getStatusName(record.status) }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" @click="previewNotice(record)">
              <component :is="icons.Eye" />
            </a-button>
            <a-button size="small" @click="editNotice(record)">
              <component :is="icons.Edit" />
            </a-button>
            <a-button size="small" danger @click="deleteNotice(record.id)">
              <component :is="icons.Delete" />
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>
    
    <a-modal
      v-model:open="showModal"
      :title="editingNotice ? '编辑公告' : '发布公告'"
      @ok="handleOk"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="标题">
          <a-input v-model:value="form.title" placeholder="请输入公告标题" />
        </a-form-item>
        <a-form-item label="摘要">
          <a-input v-model:value="form.summary" placeholder="请输入公告摘要" />
        </a-form-item>
        <a-form-item label="内容">
          <a-textarea v-model:value="form.content" placeholder="请输入公告内容" :rows="5" />
        </a-form-item>
        <a-form-item label="受众">
          <a-select v-model:value="form.audience" placeholder="请选择受众">
            <a-select-option value="all">全部用户</a-select-option>
            <a-select-option value="student">学生</a-select-option>
            <a-select-option value="teacher">教师</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="标签">
          <a-input v-model:value="form.tags" placeholder="多个标签用逗号分隔" />
        </a-form-item>
        <a-form-item label="是否置顶">
          <a-switch v-model:checked="form.pinned" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined as Plus, EditOutlined as Edit, DeleteOutlined as Delete, EyeOutlined as Eye } from '@ant-design/icons-vue'
import { portalApi } from '@/api'

interface Notice {
  id: string
  title: string
  summary: string
  content: string
  audience: string
  tags: string[]
  pinned: boolean
  status: string
  created_at: string
}

const icons = { Plus, Edit, Delete, Eye }
const showModal = ref(false)
const editingNotice = ref<Notice | null>(null)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const form = reactive({
  title: '',
  summary: '',
  content: '',
  audience: 'all',
  tags: '',
  pinned: false
})

const notices = ref<Notice[]>([])

const tablePagination = computed(() => ({
  current: currentPage.value,
  pageSize: pageSize.value,
  total: total.value,
  showSizeChanger: true,
  showQuickJumper: true
}))

const columns = [
  { title: '标题', dataIndex: 'title', key: 'title' },
  { title: '摘要', dataIndex: 'summary', key: 'summary' },
  { title: '置顶', key: 'pinned' },
  { title: '状态', key: 'status' },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
  { title: '操作', key: 'actions' }
]

function getStatusName(status: string): string {
  const names: Record<string, string> = {
    published: '已发布',
    reviewing: '审核中',
    draft: '草稿'
  }
  return names[status] || status
}

const fetchNotices = async () => {
  try {
    const response = await portalApi.listNotices({
      page: currentPage.value,
      pageSize: pageSize.value
    })
    if (response.data) {
      notices.value = response.data.list.map((item: any) => ({
        id: item.id,
        title: item.title,
        summary: item.summary,
        content: item.content,
        audience: item.audience,
        tags: item.tags || [],
        pinned: item.pinned || false,
        status: item.status || 'published',
        created_at: item.created_at ? new Date(item.created_at).toLocaleString('zh-CN') : '-'
      }))
      total.value = response.data.total
    }
  } catch (error) {
    console.error('Failed to fetch notices:', error)
    message.error('获取公告列表失败')
  }
}

const handleTableChange = (pagination: any) => {
  currentPage.value = pagination.current
  pageSize.value = pagination.pageSize
  fetchNotices()
}

async function previewNotice(notice: Notice) {
  try {
    const response = await portalApi.getNoticeDetail(notice.id)
    const data = response.data || {}
    message.info(`公告内容：${data.content || notice.summary || notice.title}`)
  } catch (error) {
    console.error('Failed to preview notice:', error)
    message.error('获取公告详情失败')
  }
}

async function editNotice(notice: Notice) {
  try {
    const response = await portalApi.getNoticeDetail(notice.id)
    const data = response.data || {}

    editingNotice.value = notice
    form.title = data.title || notice.title
    form.summary = data.summary || notice.summary
    form.content = data.content || ''
    form.audience = data.audience || notice.audience
    form.tags = Array.isArray(data.tags) ? data.tags.join(',') : notice.tags.join(',')
    form.pinned = Boolean(data.pinned ?? notice.pinned)
    showModal.value = true
  } catch (error) {
    console.error('Failed to fetch notice detail:', error)
    message.error('获取公告详情失败')
  }
}

const deleteNotice = async (id: string) => {
  try {
    await portalApi.deleteNotice(id)
    message.success('删除成功')
    fetchNotices()
  } catch (error) {
    console.error('Failed to delete notice:', error)
    message.error('删除失败')
  }
}

const handleOk = async () => {
  try {
    if (editingNotice.value) {
      await portalApi.updateNotice(editingNotice.value.id, {
        title: form.title,
        summary: form.summary,
        content: form.content,
        audience: form.audience,
        tags: form.tags.split(',').filter(t => t.trim()),
        pinned: form.pinned
      })
      message.success('更新成功')
    } else {
      await portalApi.publishNotice({
        title: form.title,
        summary: form.summary,
        content: form.content,
        audience: form.audience,
        tags: form.tags.split(',').filter(t => t.trim()),
        pinned: form.pinned
      })
      message.success('发布成功')
    }
    showModal.value = false
    editingNotice.value = null
    form.title = ''
    form.summary = ''
    form.content = ''
    form.audience = 'all'
    form.tags = ''
    form.pinned = false
    fetchNotices()
  } catch (error) {
    console.error('Failed to save notice:', error)
    message.error('操作失败')
  }
}

onMounted(() => {
  fetchNotices()
})
</script>
