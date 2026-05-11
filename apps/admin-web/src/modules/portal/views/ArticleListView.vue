<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-800">文章管理</h2>
      <a-button type="primary" @click="showModal = true">
        <component :is="icons.Plus" />
        发布文章
      </a-button>
    </div>
    
    <a-table :columns="columns" :data-source="articles" :pagination="tablePagination" @change="handleTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="record.status === 'published' ? 'green' : record.status === 'reviewing' ? 'orange' : 'gray'">
            {{ getStatusName(record.status) }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" @click="editArticle(record)">
              <component :is="icons.Edit" />
            </a-button>
            <a-button size="small" @click="previewArticle(record)">
              <component :is="icons.Eye" />
            </a-button>
            <a-button size="small" danger @click="deleteArticle(record.id)">
              <component :is="icons.Delete" />
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>
    
    <a-modal
      v-model:open="showModal"
      :title="editingArticle ? '编辑文章' : '发布文章'"
      :width="800"
      @ok="handleOk"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="标题">
          <a-input v-model:value="form.title" placeholder="请输入文章标题" />
        </a-form-item>
        <a-form-item label="摘要">
          <a-textarea v-model:value="form.summary" placeholder="请输入文章摘要" :rows="3" />
        </a-form-item>
        <a-form-item label="内容">
          <a-textarea v-model:value="form.content" placeholder="请输入文章内容" :rows="6" />
        </a-form-item>
        <a-form-item label="标签">
          <a-input v-model:value="form.tags" placeholder="多个标签用逗号分隔" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { PlusOutlined as Plus, EditOutlined as Edit, EyeOutlined as Eye, DeleteOutlined as Delete } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { portalApi } from '@/api'

interface Article {
  id: string
  title: string
  summary: string
  content: string
  tags: string[]
  status: string
  created_at: string
}

const icons = { Plus, Edit, Eye, Delete }
const showModal = ref(false)
const editingArticle = ref<Article | null>(null)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const form = reactive({
  title: '',
  summary: '',
  content: '',
  tags: ''
})

const articles = ref<Article[]>([])

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
  { title: '状态', key: 'status' },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
  { title: '操作', key: 'actions' }
]

function getStatusName(status: string): string {
  const names: Record<string, string> = {
    draft: '草稿',
    published: '已发布',
    reviewing: '审核中',
    archived: '已归档'
  }
  return names[status] || status
}

const fetchArticles = async () => {
  try {
    const response = await portalApi.listArticles({
      page: currentPage.value,
      pageSize: pageSize.value
    })
    if (response.data) {
      articles.value = response.data.list.map((item: any) => ({
        id: item.id,
        title: item.title,
        summary: item.summary,
        content: item.content,
        tags: item.tags || [],
        status: item.status || 'draft',
        created_at: item.created_at ? new Date(item.created_at).toLocaleString('zh-CN') : '-'
      }))
      total.value = response.data.total
    }
  } catch (error) {
    console.error('Failed to fetch articles:', error)
    message.error('获取文章列表失败')
  }
}

const handleTableChange = (pagination: any) => {
  currentPage.value = pagination.current
  pageSize.value = pagination.pageSize
  fetchArticles()
}

function editArticle(article: Article) {
  editingArticle.value = article
  form.title = article.title
  form.summary = article.summary
  form.content = article.content
  form.tags = article.tags.join(',')
  showModal.value = true
}

function previewArticle(article: Article) {
  message.info(`预览文章: ${article.title}`)
}

const deleteArticle = async (id: string) => {
  try {
    await portalApi.deleteArticle(id)
    message.success('删除成功')
    fetchArticles()
  } catch (error) {
    console.error('Failed to delete article:', error)
    message.error('删除失败')
  }
}

const handleOk = async () => {
  try {
    if (editingArticle.value) {
      await portalApi.updateArticle(editingArticle.value.id, {
        title: form.title,
        summary: form.summary,
        content: form.content,
        tags: form.tags.split(',').filter(t => t.trim())
      })
      message.success('更新成功')
    } else {
      await portalApi.createArticle({
        title: form.title,
        summary: form.summary,
        content: form.content,
        tags: form.tags.split(',').filter(t => t.trim())
      })
      message.success('发布成功')
    }
    showModal.value = false
    editingArticle.value = null
    form.title = ''
    form.summary = ''
    form.content = ''
    form.tags = ''
    fetchArticles()
  } catch (error) {
    console.error('Failed to save article:', error)
    message.error('操作失败')
  }
}

onMounted(() => {
  fetchArticles()
})
</script>