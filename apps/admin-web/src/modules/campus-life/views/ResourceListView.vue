<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-800">资料共享</h2>
      <a-button type="primary" @click="showModal = true">
        <component :is="icons.Plus" />
        上传资料
      </a-button>
    </div>
    
    <a-table :columns="columns" :data-source="resources" :pagination="pagination">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'downloads'">
          <span class="text-blue-500">{{ record.downloads }}次下载</span>
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" @click="downloadResource(record)">
              <component :is="icons.Download" />
            </a-button>
            <a-button size="small" danger @click="deleteResource(record.id)">
              <component :is="icons.Delete" />
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>
    
    <a-modal
      v-model:open="showModal"
      :title="editingResource ? '编辑资料' : '上传资料'"
      @ok="handleOk"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="资料名称">
          <a-input v-model:value="form.title" placeholder="请输入资料名称" />
        </a-form-item>
        <a-form-item label="资料分类">
          <a-select v-model:value="form.category" placeholder="请选择分类">
            <a-select-option value="course">课程资料</a-select-option>
            <a-select-option value="exam">考试资料</a-select-option>
            <a-select-option value="thesis">论文资料</a-select-option>
            <a-select-option value="other">其他</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="资料描述">
          <a-textarea v-model:value="form.description" placeholder="请输入资料描述" :rows="3" />
        </a-form-item>
        <a-form-item label="文件链接">
          <a-input v-model:value="form.fileUrl" placeholder="请输入文件链接" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { PlusOutlined as Plus, DownloadOutlined as Download, DeleteOutlined as Delete } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'

interface Resource {
  id: string
  title: string
  category: string
  description: string
  fileUrl: string
  downloads: number
  uploader: string
  createdAt: string
}

const icons = { Plus, Download, Delete }
const showModal = ref(false)
const editingResource = ref<Resource | null>(null)

const form = reactive({
  title: '',
  category: '',
  description: '',
  fileUrl: ''
})

const resources = ref<Resource[]>([
  { id: '1', title: '高等数学期末复习资料', category: 'exam', description: '包含历年真题和复习要点', fileUrl: '/files/math-review.pdf', downloads: 156, uploader: 'student001', createdAt: '2024-01-15' },
  { id: '2', title: '操作系统课件', category: 'course', description: '操作系统课程完整课件', fileUrl: '/files/os-slides.pdf', downloads: 89, uploader: 'student002', createdAt: '2024-01-14' },
  { id: '3', title: '毕业设计模板', category: 'thesis', description: '本科毕业设计模板', fileUrl: '/files/thesis-template.docx', downloads: 234, uploader: 'student003', createdAt: '2024-01-13' }
])

const pagination = computed(() => ({
  total: resources.value.length,
  showSizeChanger: true,
  showQuickJumper: true
}))

const columns = [
  { title: '资料名称', dataIndex: 'title', key: 'title' },
  { title: '分类', dataIndex: 'category', key: 'category' },
  { title: '描述', dataIndex: 'description', key: 'description' },
  { title: '下载次数', key: 'downloads' },
  { title: '上传者', dataIndex: 'uploader', key: 'uploader' },
  { title: '上传时间', dataIndex: 'createdAt', key: 'createdAt' },
  { title: '操作', key: 'actions' }
]

function downloadResource(resource: Resource) {
  message.info(`下载资料: ${resource.title}`)
}

function deleteResource(id: string) {
  resources.value = resources.value.filter(r => r.id !== id)
}

function handleOk() {
  if (editingResource.value) {
    const index = resources.value.findIndex(r => r.id === editingResource.value!.id)
    if (index !== -1) {
      resources.value[index] = { ...resources.value[index], ...form }
    }
  } else {
    resources.value.push({
      id: Date.now().toString(),
      ...form,
      downloads: 0,
      uploader: 'admin',
      createdAt: new Date().toISOString().split('T')[0]
    })
  }
  showModal.value = false
  editingResource.value = null
  form.title = ''
  form.category = ''
  form.description = ''
  form.fileUrl = ''
}
</script>