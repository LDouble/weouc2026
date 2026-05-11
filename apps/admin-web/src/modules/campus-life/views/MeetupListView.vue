<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-800">组局活动</h2>
      <a-button type="primary" @click="showModal = true">
        <component :is="icons.Plus" />
        创建活动
      </a-button>
    </div>
    
    <a-table :columns="columns" :data-source="meetups" :pagination="pagination">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusName(record.status) }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'participants'">
          <span>{{ record.participants }}人参与</span>
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" @click="viewDetail(record)">
              <component :is="icons.Eye" />
            </a-button>
            <a-button size="small" danger @click="deleteMeetup(record.id)">
              <component :is="icons.Delete" />
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>
    
    <a-modal
      v-model:open="showModal"
      :title="editingMeetup ? '编辑活动' : '创建活动'"
      @ok="handleOk"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="活动名称">
          <a-input v-model:value="form.title" placeholder="请输入活动名称" />
        </a-form-item>
        <a-form-item label="活动类型">
          <a-select v-model:value="form.type" placeholder="请选择活动类型">
            <a-select-option value="study">学习小组</a-select-option>
            <a-select-option value="sports">体育活动</a-select-option>
            <a-select-option value="carpool">拼车出行</a-select-option>
            <a-select-option value="other">其他</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="活动描述">
          <a-textarea v-model:value="form.description" placeholder="请输入活动描述" :rows="3" />
        </a-form-item>
        <a-form-item label="活动时间">
          <a-input v-model:value="form.time" placeholder="请输入活动时间" />
        </a-form-item>
        <a-form-item label="活动地点">
          <a-input v-model:value="form.location" placeholder="请输入活动地点" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="form.status" placeholder="请选择状态">
            <a-select-option value="pending">报名中</a-select-option>
            <a-select-option value="ongoing">进行中</a-select-option>
            <a-select-option value="ended">已结束</a-select-option>
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

interface Meetup {
  id: string
  title: string
  type: string
  description: string
  time: string
  location: string
  participants: number
  status: string
  organizer: string
  createdAt: string
}

const icons = { Plus, Eye, Delete }
const showModal = ref(false)
const editingMeetup = ref<Meetup | null>(null)

const form = reactive({
  title: '',
  type: '',
  description: '',
  time: '',
  location: '',
  status: 'pending'
})

const meetups = ref<Meetup[]>([
  { id: '1', title: '周末篮球友谊赛', type: 'sports', description: '周末下午篮球场约球', time: '2024-01-20 14:00', location: '西校区篮球场', participants: 8, status: 'pending', organizer: 'student001', createdAt: '2024-01-15' },
  { id: '2', title: '考研学习小组', type: 'study', description: '每天晚上一起复习考研', time: '每晚 19:00-22:00', location: '图书馆三楼', participants: 12, status: 'ongoing', organizer: 'student002', createdAt: '2024-01-14' },
  { id: '3', title: '假期拼车回家', type: 'carpool', description: '寒假拼车回家，还差2人', time: '2024-01-25 08:00', location: '学校东门', participants: 2, status: 'pending', organizer: 'student003', createdAt: '2024-01-13' }
])

const pagination = computed(() => ({
  total: meetups.value.length,
  showSizeChanger: true,
  showQuickJumper: true
}))

const columns = [
  { title: '活动名称', dataIndex: 'title', key: 'title' },
  { title: '类型', dataIndex: 'type', key: 'type' },
  { title: '时间', dataIndex: 'time', key: 'time' },
  { title: '地点', dataIndex: 'location', key: 'location' },
  { title: '参与人数', key: 'participants' },
  { title: '状态', key: 'status' },
  { title: '组织者', dataIndex: 'organizer', key: 'organizer' },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
  { title: '操作', key: 'actions' }
]

function getStatusColor(status: string): string {
  const colors: Record<string, string> = {
    pending: 'orange',
    ongoing: 'blue',
    ended: 'gray'
  }
  return colors[status] || 'gray'
}

function getStatusName(status: string): string {
  const names: Record<string, string> = {
    pending: '报名中',
    ongoing: '进行中',
    ended: '已结束'
  }
  return names[status] || status
}

function viewDetail(meetup: Meetup) {
  message.info(`查看活动详情: ${meetup.title}`)
}

function deleteMeetup(id: string) {
  meetups.value = meetups.value.filter(m => m.id !== id)
}

function handleOk() {
  if (editingMeetup.value) {
    const index = meetups.value.findIndex(m => m.id === editingMeetup.value!.id)
    if (index !== -1) {
      meetups.value[index] = { ...meetups.value[index], ...form }
    }
  } else {
    meetups.value.push({
      id: Date.now().toString(),
      ...form,
      participants: 0,
      organizer: 'admin',
      createdAt: new Date().toISOString().split('T')[0]
    })
  }
  showModal.value = false
  editingMeetup.value = null
  form.title = ''
  form.type = ''
  form.description = ''
  form.time = ''
  form.location = ''
  form.status = 'pending'
}
</script>