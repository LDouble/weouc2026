<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-800">失物招领</h2>
      <a-button type="primary" @click="showModal = true">
        <component :is="icons.Plus" />
        发布失物招领
      </a-button>
    </div>
    
    <a-table :columns="columns" :data-source="lostItems" :pagination="pagination">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'type'">
          <a-tag :color="record.type === 'lost' ? 'orange' : 'green'">
            {{ record.type === 'lost' ? '寻物启事' : '失物招领' }}
          </a-tag>
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
            <a-button size="small" danger @click="deleteLostItem(record.id)">
              <component :is="icons.Delete" />
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>
    
    <a-modal
      v-model:open="showModal"
      :title="editingLostItem ? '编辑失物招领' : '发布失物招领'"
      @ok="handleOk"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="类型">
          <a-select v-model:value="form.type" placeholder="请选择类型">
            <a-select-option value="lost">寻物启事</a-select-option>
            <a-select-option value="found">失物招领</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="物品名称">
          <a-input v-model:value="form.title" placeholder="请输入物品名称" />
        </a-form-item>
        <a-form-item label="物品描述">
          <a-textarea v-model:value="form.description" placeholder="请输入物品描述" :rows="3" />
        </a-form-item>
        <a-form-item label="地点">
          <a-input v-model:value="form.location" placeholder="请输入丢失/发现地点" />
        </a-form-item>
        <a-form-item label="时间">
          <a-input v-model:value="form.time" placeholder="请输入丢失/发现时间" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="form.status" placeholder="请选择状态">
            <a-select-option value="pending">寻找中</a-select-option>
            <a-select-option value="claimed">已认领</a-select-option>
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

interface LostItem {
  id: string
  type: string
  title: string
  description: string
  location: string
  time: string
  status: string
  publisher: string
  createdAt: string
}

const icons = { Plus, Eye, Delete }
const showModal = ref(false)
const editingLostItem = ref<LostItem | null>(null)

const form = reactive({
  type: 'lost',
  title: '',
  description: '',
  location: '',
  time: '',
  status: 'pending'
})

const lostItems = ref<LostItem[]>([
  { id: '1', type: 'lost', title: '学生证', description: '蓝色外壳，姓名张三', location: '图书馆三楼', time: '2024-01-15 14:00', status: 'pending', publisher: 'student001', createdAt: '2024-01-15' },
  { id: '2', type: 'found', title: '苹果耳机', description: '白色AirPods Pro', location: '食堂二楼', time: '2024-01-14 12:00', status: 'claimed', publisher: 'student002', createdAt: '2024-01-14' },
  { id: '3', type: 'lost', title: '笔记本电脑', description: '银色MacBook Pro 14寸', location: '教学楼A座', time: '2024-01-13 18:00', status: 'pending', publisher: 'student003', createdAt: '2024-01-13' }
])

const pagination = computed(() => ({
  total: lostItems.value.length,
  showSizeChanger: true,
  showQuickJumper: true
}))

const columns = [
  { title: '类型', key: 'type' },
  { title: '物品名称', dataIndex: 'title', key: 'title' },
  { title: '描述', dataIndex: 'description', key: 'description' },
  { title: '地点', dataIndex: 'location', key: 'location' },
  { title: '时间', dataIndex: 'time', key: 'time' },
  { title: '状态', key: 'status' },
  { title: '发布者', dataIndex: 'publisher', key: 'publisher' },
  { title: '发布时间', dataIndex: 'createdAt', key: 'createdAt' },
  { title: '操作', key: 'actions' }
]

function getStatusColor(status: string): string {
  const colors: Record<string, string> = {
    pending: 'orange',
    claimed: 'green'
  }
  return colors[status] || 'gray'
}

function getStatusName(status: string): string {
  const names: Record<string, string> = {
    pending: '寻找中',
    claimed: '已认领'
  }
  return names[status] || status
}

function viewDetail(lostItem: LostItem) {
  message.info(`查看详情: ${lostItem.title}`)
}

function deleteLostItem(id: string) {
  lostItems.value = lostItems.value.filter(l => l.id !== id)
}

function handleOk() {
  if (editingLostItem.value) {
    const index = lostItems.value.findIndex(l => l.id === editingLostItem.value!.id)
    if (index !== -1) {
      lostItems.value[index] = { ...lostItems.value[index], ...form }
    }
  } else {
    lostItems.value.push({
      id: Date.now().toString(),
      ...form,
      publisher: 'admin',
      createdAt: new Date().toISOString().split('T')[0]
    })
  }
  showModal.value = false
  editingLostItem.value = null
  form.type = 'lost'
  form.title = ''
  form.description = ''
  form.location = ''
  form.time = ''
  form.status = 'pending'
}
</script>