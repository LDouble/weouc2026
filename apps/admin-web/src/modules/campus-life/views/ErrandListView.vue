<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-800">跑腿服务</h2>
      <a-button type="primary" @click="showModal = true">
        <component :is="icons.Plus" />
        添加跑腿
      </a-button>
    </div>
    
    <a-table :columns="columns" :data-source="errands" :pagination="pagination">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusName(record.status) }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" @click="viewDetail(record)">
              <component :is="icons.Eye" />
            </a-button>
            <a-button size="small" danger @click="deleteErrand(record.id)">
              <component :is="icons.Delete" />
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>
    
    <a-modal
      v-model:open="showModal"
      :title="editingErrand ? '编辑跑腿' : '添加跑腿'"
      @ok="handleOk"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="标题">
          <a-input v-model:value="form.title" placeholder="请输入跑腿标题" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model:value="form.description" placeholder="请输入跑腿描述" :rows="3" />
        </a-form-item>
        <a-form-item label="报酬">
          <a-input v-model:value="form.reward" placeholder="请输入报酬金额" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="form.status" placeholder="请选择状态">
            <a-select-option value="pending">待接单</a-select-option>
            <a-select-option value="accepted">已接单</a-select-option>
            <a-select-option value="completed">已完成</a-select-option>
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

interface Errand {
  id: string
  title: string
  description: string
  reward: string
  status: string
  publisher: string
  createdAt: string
}

const icons = { Plus, Eye, Delete }
const showModal = ref(false)
const editingErrand = ref<Errand | null>(null)

const form = reactive({
  title: '',
  description: '',
  reward: '',
  status: 'pending'
})

const errands = ref<Errand[]>([
  { id: '1', title: '帮取快递', description: '帮忙从菜鸟驿站取快递', reward: '10元', status: 'completed', publisher: 'student001', createdAt: '2024-01-15' },
  { id: '2', title: '食堂带饭', description: '帮忙从三食堂带一份饭', reward: '5元', status: 'accepted', publisher: 'student002', createdAt: '2024-01-14' },
  { id: '3', title: '图书馆占座', description: '早上8点帮忙占座', reward: '8元', status: 'pending', publisher: 'student003', createdAt: '2024-01-13' }
])

const pagination = computed(() => ({
  total: errands.value.length,
  showSizeChanger: true,
  showQuickJumper: true
}))

const columns = [
  { title: '标题', dataIndex: 'title', key: 'title' },
  { title: '描述', dataIndex: 'description', key: 'description' },
  { title: '报酬', dataIndex: 'reward', key: 'reward' },
  { title: '状态', key: 'status' },
  { title: '发布者', dataIndex: 'publisher', key: 'publisher' },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
  { title: '操作', key: 'actions' }
]

function getStatusColor(status: string): string {
  const colors: Record<string, string> = {
    pending: 'orange',
    accepted: 'blue',
    completed: 'green'
  }
  return colors[status] || 'gray'
}

function getStatusName(status: string): string {
  const names: Record<string, string> = {
    pending: '待接单',
    accepted: '已接单',
    completed: '已完成'
  }
  return names[status] || status
}

function viewDetail(errand: Errand) {
  message.info(`查看跑腿详情: ${errand.title}`)
}

function deleteErrand(id: string) {
  errands.value = errands.value.filter(e => e.id !== id)
}

function handleOk() {
  if (editingErrand.value) {
    const index = errands.value.findIndex(e => e.id === editingErrand.value!.id)
    if (index !== -1) {
      errands.value[index] = { ...errands.value[index], ...form }
    }
  } else {
    errands.value.push({
      id: Date.now().toString(),
      ...form,
      publisher: 'admin',
      createdAt: new Date().toISOString().split('T')[0]
    })
  }
  showModal.value = false
  editingErrand.value = null
  form.title = ''
  form.description = ''
  form.reward = ''
  form.status = 'pending'
}
</script>