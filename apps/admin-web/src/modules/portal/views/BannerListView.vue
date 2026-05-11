<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-6">
      <div>
        <h2 class="text-lg font-semibold text-gray-800">轮播管理</h2>
        <p class="text-sm text-gray-500 mt-1">当前页面已切换到真实后端接口，维护门户首页轮播内容。</p>
      </div>
      <a-button type="primary" @click="openCreateModal">
        <component :is="icons.Plus" />
        添加轮播
      </a-button>
    </div>

    <a-table :columns="columns" :data-source="banners" :pagination="tablePagination" @change="handleTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'image'">
          <img :src="record.image_url" :alt="record.title" class="w-20 h-12 object-cover rounded" />
        </template>
        <template v-else-if="column.key === 'sort'">
          <span>{{ record.sort }}</span>
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" @click="editBanner(record)">
              <component :is="icons.Edit" />
            </a-button>
            <a-button size="small" danger @click="deleteBanner(record.id)">
              <component :is="icons.Delete" />
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="showModal" :title="editingBanner ? '编辑轮播' : '添加轮播'" @ok="handleOk">
      <a-form :model="form" layout="vertical">
        <a-form-item label="标题">
          <a-input v-model:value="form.title" placeholder="请输入轮播标题" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model:value="form.description" placeholder="请输入轮播描述" :rows="3" />
        </a-form-item>
        <a-form-item label="图片链接">
          <a-input v-model:value="form.image_url" placeholder="请输入图片 URL" />
        </a-form-item>
        <a-form-item label="跳转链接">
          <a-input v-model:value="form.action_url" placeholder="请输入跳转链接" />
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model:value="form.sort" :min="0" class="w-full" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { DeleteOutlined as Delete, EditOutlined as Edit, PlusOutlined as Plus } from '@ant-design/icons-vue'
import { portalApi } from '@/api'

interface Banner {
  id: string
  title: string
  description: string
  image_url: string
  action_url: string
  sort: number
  created_at: string
}

const icons = { Plus, Edit, Delete }
const showModal = ref(false)
const editingBanner = ref<Banner | null>(null)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const banners = ref<Banner[]>([])

const form = reactive({
  title: '',
  description: '',
  image_url: '',
  action_url: '',
  sort: 0
})

const tablePagination = computed(() => ({
  current: currentPage.value,
  pageSize: pageSize.value,
  total: total.value,
  showSizeChanger: true,
  showQuickJumper: true
}))

const columns = [
  { title: '标题', dataIndex: 'title', key: 'title' },
  { title: '图片', key: 'image' },
  { title: '描述', dataIndex: 'description', key: 'description' },
  { title: '跳转链接', dataIndex: 'action_url', key: 'action_url' },
  { title: '排序', key: 'sort' },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
  { title: '操作', key: 'actions' }
]

function resetForm() {
  form.title = ''
  form.description = ''
  form.image_url = ''
  form.action_url = ''
  form.sort = 0
}

async function fetchBanners() {
  try {
    const response = await portalApi.listBanners({
      page: currentPage.value,
      pageSize: pageSize.value
    })
    banners.value = (response.data?.list || []).map((item: any) => ({
      id: item.id,
      title: item.title,
      description: item.description || '',
      image_url: item.image_url || '',
      action_url: item.action_url || '',
      sort: Number(item.sort || 0),
      created_at: item.created_at ? new Date(item.created_at).toLocaleString('zh-CN') : '-'
    }))
    total.value = Number(response.data?.total || 0)
  } catch (error) {
    console.error('Failed to fetch banners:', error)
    message.error('获取轮播列表失败')
  }
}

function handleTableChange(pagination: any) {
  currentPage.value = pagination.current
  pageSize.value = pagination.pageSize
  fetchBanners()
}

function openCreateModal() {
  editingBanner.value = null
  resetForm()
  showModal.value = true
}

function editBanner(banner: Banner) {
  editingBanner.value = banner
  form.title = banner.title
  form.description = banner.description
  form.image_url = banner.image_url
  form.action_url = banner.action_url
  form.sort = banner.sort
  showModal.value = true
}

async function deleteBanner(id: string) {
  try {
    await portalApi.deleteBanner(id)
    message.success('删除成功')
    await fetchBanners()
  } catch (error) {
    console.error('Failed to delete banner:', error)
    message.error('删除失败')
  }
}

async function handleOk() {
  try {
    const payload = {
      title: form.title,
      description: form.description,
      image_url: form.image_url,
      action_url: form.action_url,
      sort: form.sort
    }

    if (editingBanner.value) {
      await portalApi.updateBanner(editingBanner.value.id, payload)
      message.success('更新成功')
    } else {
      await portalApi.createBanner(payload)
      message.success('创建成功')
    }

    showModal.value = false
    editingBanner.value = null
    resetForm()
    await fetchBanners()
  } catch (error) {
    console.error('Failed to save banner:', error)
    message.error('操作失败')
  }
}

onMounted(() => {
  fetchBanners()
})
</script>
