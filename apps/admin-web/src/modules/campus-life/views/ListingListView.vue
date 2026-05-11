<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-800">二手交易</h2>
      <a-button type="primary" @click="showModal = true">
        <component :is="icons.Plus" />
        发布二手
      </a-button>
    </div>
    
    <a-table :columns="columns" :data-source="listings" :pagination="pagination">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'image'">
          <img :src="record.image" :alt="record.title" class="w-16 h-12 object-cover rounded" />
        </template>
        <template v-if="column.key === 'price'">
          <span class="text-red-500 font-medium">{{ record.price }}</span>
        </template>
        <template v-else-if="column.key === 'status'">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusName(record.status) }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" @click="viewDetail(record)">
              <component :is="icons.Eye" />
            </a-button>
            <a-button size="small" danger @click="deleteListing(record.id)">
              <component :is="icons.Delete" />
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>
    
    <a-modal
      v-model:open="showModal"
      :title="editingListing ? '编辑二手' : '发布二手'"
      @ok="handleOk"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="物品名称">
          <a-input v-model:value="form.title" placeholder="请输入物品名称" />
        </a-form-item>
        <a-form-item label="物品图片">
          <a-input v-model:value="form.image" placeholder="请输入图片URL" />
        </a-form-item>
        <a-form-item label="物品描述">
          <a-textarea v-model:value="form.description" placeholder="请输入物品描述" :rows="3" />
        </a-form-item>
        <a-form-item label="价格">
          <a-input v-model:value="form.price" placeholder="请输入价格" />
        </a-form-item>
        <a-form-item label="物品分类">
          <a-select v-model:value="form.category" placeholder="请选择分类">
            <a-select-option value="electronics">电子产品</a-select-option>
            <a-select-option value="books">图书教材</a-select-option>
            <a-select-option value="clothing">服饰</a-select-option>
            <a-select-option value="other">其他</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="form.status" placeholder="请选择状态">
            <a-select-option value="available">在售</a-select-option>
            <a-select-option value="sold">已售出</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { PlusOutlined as Plus, EyeOutlined as Eye, DeleteOutlined as Delete } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'

interface Listing {
  id: string
  title: string
  image: string
  description: string
  price: string
  category: string
  status: string
  seller: string
  createdAt: string
}

const icons = { Plus, Eye, Delete }
const showModal = ref(false)
const editingListing = ref<Listing | null>(null)

const form = reactive({
  title: '',
  image: '',
  description: '',
  price: '',
  category: '',
  status: 'available'
})

const listings = ref<Listing[]>([
  { id: '1', title: 'iPhone 14 Pro', image: 'https://via.placeholder.com/200x200', description: '95新，无磕碰', price: '5000元', category: 'electronics', status: 'available', seller: 'student001', createdAt: '2024-01-15' },
  { id: '2', title: '高等数学教材', image: 'https://via.placeholder.com/200x200', description: '几乎全新，附赠笔记', price: '30元', category: 'books', status: 'sold', seller: 'student002', createdAt: '2024-01-14' },
  { id: '3', title: '羽绒服', image: 'https://via.placeholder.com/200x200', description: 'M码，穿过两次', price: '200元', category: 'clothing', status: 'available', seller: 'student003', createdAt: '2024-01-13' }
])

const pagination = computed(() => ({
  total: listings.value.length,
  showSizeChanger: true,
  showQuickJumper: true
}))

const columns = [
  { title: '图片', key: 'image' },
  { title: '物品名称', dataIndex: 'title', key: 'title' },
  { title: '分类', dataIndex: 'category', key: 'category' },
  { title: '价格', key: 'price' },
  { title: '状态', key: 'status' },
  { title: '卖家', dataIndex: 'seller', key: 'seller' },
  { title: '发布时间', dataIndex: 'createdAt', key: 'createdAt' },
  { title: '操作', key: 'actions' }
]

function getStatusColor(status: string): string {
  const colors: Record<string, string> = {
    available: 'green',
    sold: 'gray'
  }
  return colors[status] || 'gray'
}

function getStatusName(status: string): string {
  const names: Record<string, string> = {
    available: '在售',
    sold: '已售出'
  }
  return names[status] || status
}

function viewDetail(listing: Listing) {
  message.info(`查看二手详情: ${listing.title}`)
}

function deleteListing(id: string) {
  listings.value = listings.value.filter(l => l.id !== id)
}

function handleOk() {
  if (editingListing.value) {
    const index = listings.value.findIndex(l => l.id === editingListing.value!.id)
    if (index !== -1) {
      listings.value[index] = { ...listings.value[index], ...form }
    }
  } else {
    listings.value.push({
      id: Date.now().toString(),
      ...form,
      seller: 'admin',
      createdAt: new Date().toISOString().split('T')[0]
    })
  }
  showModal.value = false
  editingListing.value = null
  form.title = ''
  form.image = ''
  form.description = ''
  form.price = ''
  form.category = ''
  form.status = 'available'
}
</script>