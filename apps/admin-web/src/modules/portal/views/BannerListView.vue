<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-800">轮播管理</h2>
      <a-button type="primary" @click="showModal = true">
        <component :is="icons.Plus" />
        添加轮播
      </a-button>
    </div>
    
    <a-table :columns="columns" :data-source="banners" :pagination="pagination">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'image'">
          <img :src="record.image" :alt="record.title" class="w-20 h-12 object-cover rounded" />
        </template>
        <template v-else-if="column.key === 'status'">
          <a-tag :color="record.status === 'active' ? 'green' : 'gray'">
            {{ record.status === 'active' ? '启用' : '禁用' }}
          </a-tag>
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
    
    <a-modal
      v-model:open="showModal"
      :title="editingBanner ? '编辑轮播' : '添加轮播'"
      @ok="handleOk"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="标题">
          <a-input v-model:value="form.title" placeholder="请输入轮播标题" />
        </a-form-item>
        <a-form-item label="图片链接">
          <a-input v-model:value="form.image" placeholder="请输入图片URL" />
        </a-form-item>
        <a-form-item label="跳转链接">
          <a-input v-model:value="form.link" placeholder="请输入跳转链接" />
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model:value="form.sortOrder" :min="0" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="form.status" placeholder="请选择状态">
            <a-select-option value="active">启用</a-select-option>
            <a-select-option value="disabled">禁用</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { PlusOutlined as Plus, EditOutlined as Edit, DeleteOutlined as Delete } from '@ant-design/icons-vue'

interface Banner {
  id: string
  title: string
  image: string
  link: string
  sortOrder: number
  status: string
  createdAt: string
}

const icons = { Plus, Edit, Delete }
const showModal = ref(false)
const editingBanner = ref<Banner | null>(null)

const form = reactive({
  title: '',
  image: '',
  link: '',
  sortOrder: 0,
  status: 'active'
})

const banners = ref<Banner[]>([
  { id: '1', title: '新生入学指南', image: 'https://via.placeholder.com/300x150', link: '/guide', sortOrder: 1, status: 'active', createdAt: '2024-01-15' },
  { id: '2', title: '校园招聘季', image: 'https://via.placeholder.com/300x150', link: '/career', sortOrder: 2, status: 'active', createdAt: '2024-01-14' },
  { id: '3', title: '运动会报名', image: 'https://via.placeholder.com/300x150', link: '/sports', sortOrder: 3, status: 'disabled', createdAt: '2024-01-13' }
])

const pagination = computed(() => ({
  total: banners.value.length,
  showSizeChanger: true,
  showQuickJumper: true
}))

const columns = [
  { title: '标题', dataIndex: 'title', key: 'title' },
  { title: '图片', key: 'image' },
  { title: '跳转链接', dataIndex: 'link', key: 'link' },
  { title: '排序', dataIndex: 'sortOrder', key: 'sortOrder' },
  { title: '状态', key: 'status' },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
  { title: '操作', key: 'actions' }
]

function editBanner(banner: Banner) {
  editingBanner.value = banner
  form.title = banner.title
  form.image = banner.image
  form.link = banner.link
  form.sortOrder = banner.sortOrder
  form.status = banner.status
  showModal.value = true
}

function deleteBanner(id: string) {
  banners.value = banners.value.filter(b => b.id !== id)
}

function handleOk() {
  if (editingBanner.value) {
    const index = banners.value.findIndex(b => b.id === editingBanner.value!.id)
    if (index !== -1) {
      banners.value[index] = { ...banners.value[index], ...form }
    }
  } else {
    banners.value.push({
      id: Date.now().toString(),
      ...form,
      createdAt: new Date().toISOString().split('T')[0]
    })
  }
  showModal.value = false
  editingBanner.value = null
  form.title = ''
  form.image = ''
  form.link = ''
  form.sortOrder = 0
  form.status = 'active'
}
</script>