<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-800">用户管理</h2>
      <a-button type="primary" @click="showModal = true">
        <component :is="icons.Plus" />
        添加用户
      </a-button>
    </div>
    
    <a-table :columns="columns" :data-source="users" :pagination="pagination">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" @click="editUser(record)">
              <component :is="icons.Edit" />
            </a-button>
            <a-button size="small" danger @click="deleteUser(record.id)">
              <component :is="icons.Delete" />
            </a-button>
          </a-space>
        </template>
        <template v-else-if="column.key === 'status'">
          <a-badge :status="record.status === 'active' ? 'success' : 'warning'" />
          <span>{{ record.status === 'active' ? '活跃' : '禁用' }}</span>
        </template>
        <template v-else-if="column.key === 'academicBound'">
          <a-tag :color="record.academicBound ? 'green' : 'red'">
            {{ record.academicBound ? '已绑定' : '未绑定' }}
          </a-tag>
        </template>
      </template>
    </a-table>
    
    <a-modal
      v-model:open="showModal"
      :title="editingUser ? '编辑用户' : '添加用户'"
      @ok="handleOk"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="用户名">
          <a-input v-model:value="form.username" placeholder="请输入用户名" />
        </a-form-item>
        <a-form-item label="姓名">
          <a-input v-model:value="form.name" placeholder="请输入姓名" />
        </a-form-item>
        <a-form-item label="邮箱">
          <a-input v-model:value="form.email" placeholder="请输入邮箱" />
        </a-form-item>
        <a-form-item label="角色">
          <a-select v-model:value="form.role" placeholder="请选择角色">
            <a-select-option value="student">学生</a-select-option>
            <a-select-option value="teacher">教师</a-select-option>
            <a-select-option value="admin">管理员</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { PlusOutlined as Plus, EditOutlined as Edit, DeleteOutlined as Delete } from '@ant-design/icons-vue'

interface User {
  id: string
  username: string
  name: string
  email: string
  role: string
  status: string
  academicBound: boolean
  createdAt: string
}

const icons = { Plus, Edit, Delete }
const showModal = ref(false)
const editingUser = ref<User | null>(null)

const form = reactive({
  username: '',
  name: '',
  email: '',
  role: ''
})

const users = ref<User[]>([
  { id: '1', username: 'student001', name: '张三', email: 'zhangsan@edu.cn', role: 'student', status: 'active', academicBound: true, createdAt: '2024-01-15' },
  { id: '2', username: 'teacher001', name: '李四', email: 'lisi@edu.cn', role: 'teacher', status: 'active', academicBound: true, createdAt: '2024-01-16' },
  { id: '3', username: 'admin001', name: '王五', email: 'wangwu@edu.cn', role: 'admin', status: 'active', academicBound: true, createdAt: '2024-01-17' },
  { id: '4', username: 'student002', name: '赵六', email: 'zhaoliu@edu.cn', role: 'student', status: 'disabled', academicBound: false, createdAt: '2024-01-18' }
])

const pagination = computed(() => ({
  total: users.value.length,
  showSizeChanger: true,
  showQuickJumper: true
}))

const columns = [
  { title: '用户名', dataIndex: 'username', key: 'username' },
  { title: '姓名', dataIndex: 'name', key: 'name' },
  { title: '邮箱', dataIndex: 'email', key: 'email' },
  { title: '角色', dataIndex: 'role', key: 'role' },
  { title: '状态', key: 'status' },
  { title: '教务绑定', key: 'academicBound' },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
  { title: '操作', key: 'actions' }
]

function editUser(user: User) {
  editingUser.value = user
  form.username = user.username
  form.name = user.name
  form.email = user.email
  form.role = user.role
  showModal.value = true
}

function deleteUser(id: string) {
  users.value = users.value.filter(u => u.id !== id)
}

function handleOk() {
  if (editingUser.value) {
    const index = users.value.findIndex(u => u.id === editingUser.value!.id)
    if (index !== -1) {
      users.value[index] = { ...users.value[index], ...form }
    }
  } else {
    users.value.push({
      id: Date.now().toString(),
      ...form,
      status: 'active',
      academicBound: false,
      createdAt: new Date().toISOString().split('T')[0]
    })
  }
  showModal.value = false
  editingUser.value = null
  Object.keys(form).forEach(key => form[key as keyof typeof form] = '')
}
</script>