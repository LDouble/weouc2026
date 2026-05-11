<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-800">角色管理</h2>
      <a-button type="primary" @click="showModal = true">
        <component :is="icons.Plus" />
        添加角色
      </a-button>
    </div>
    
    <a-table :columns="columns" :data-source="roles" :pagination="pagination">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'permissions'">
          <a-tag-group>
            <a-tag v-for="p in record.permissions.slice(0, 3)" :key="p" color="blue">
              {{ p }}
            </a-tag>
            <span v-if="record.permissions.length > 3">
              +{{ record.permissions.length - 3 }}
            </span>
          </a-tag-group>
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" @click="editRole(record)">
              <component :is="icons.Edit" />
            </a-button>
            <a-button size="small" danger @click="deleteRole(record.id)">
              <component :is="icons.Delete" />
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>
    
    <a-modal
      v-model:open="showModal"
      :title="editingRole ? '编辑角色' : '添加角色'"
      @ok="handleOk"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="角色名称">
          <a-input v-model:value="form.name" placeholder="请输入角色名称" />
        </a-form-item>
        <a-form-item label="角色描述">
          <a-textarea v-model:value="form.description" placeholder="请输入角色描述" />
        </a-form-item>
        <a-form-item label="权限列表">
          <a-select v-model:value="form.permissions" mode="multiple" placeholder="请选择权限">
            <a-select-option v-for="p in availablePermissions" :key="p" :value="p">
              {{ p }}
            </a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { PlusOutlined as Plus, EditOutlined as Edit, DeleteOutlined as Delete } from '@ant-design/icons-vue'

interface Role {
  id: string
  name: string
  description: string
  permissions: string[]
  createdAt: string
}

const icons = { Plus, Edit, Delete }
const showModal = ref(false)
const editingRole = ref<Role | null>(null)

const availablePermissions = [
  'iam:user:view', 'iam:user:manage',
  'iam:role:view', 'iam:role:manage',
  'iam:permission:view',
  'portal:publish', 'portal:view',
  'campus_life:moderate', 'campus_life:view',
  'analytics:view'
]

const form = reactive({
  name: '',
  description: '',
  permissions: [] as string[]
})

const roles = ref<Role[]>([
  { id: '1', name: '超级管理员', description: '拥有所有权限', permissions: availablePermissions, createdAt: '2024-01-01' },
  { id: '2', name: '内容运营', description: '负责内容发布和审核', permissions: ['portal:publish', 'portal:view', 'campus_life:view', 'campus_life:moderate'], createdAt: '2024-01-02' },
  { id: '3', name: '用户管理员', description: '负责用户管理', permissions: ['iam:user:view', 'iam:user:manage'], createdAt: '2024-01-03' },
  { id: '4', name: '数据分析员', description: '查看数据统计', permissions: ['analytics:view', 'portal:view', 'campus_life:view'], createdAt: '2024-01-04' }
])

const pagination = computed(() => ({
  total: roles.value.length,
  showSizeChanger: true,
  showQuickJumper: true
}))

const columns = [
  { title: '角色名称', dataIndex: 'name', key: 'name' },
  { title: '描述', dataIndex: 'description', key: 'description' },
  { title: '权限', key: 'permissions' },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
  { title: '操作', key: 'actions' }
]

function editRole(role: Role) {
  editingRole.value = role
  form.name = role.name
  form.description = role.description
  form.permissions = [...role.permissions]
  showModal.value = true
}

function deleteRole(id: string) {
  roles.value = roles.value.filter(r => r.id !== id)
}

function handleOk() {
  if (editingRole.value) {
    const index = roles.value.findIndex(r => r.id === editingRole.value!.id)
    if (index !== -1) {
      roles.value[index] = { ...roles.value[index], ...form }
    }
  } else {
    roles.value.push({
      id: Date.now().toString(),
      ...form,
      createdAt: new Date().toISOString().split('T')[0]
    })
  }
  showModal.value = false
  editingRole.value = null
  form.name = ''
  form.description = ''
  form.permissions = []
}
</script>
