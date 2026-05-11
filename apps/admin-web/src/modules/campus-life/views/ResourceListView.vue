<template>
  <CampusLifeManagementView
    title="资料共享"
    description="当前页面直接读取真实资料审核数据，可查看课程名、文件类型与大小，并执行状态管理。"
    content-type="resource"
    :columns="columns"
  >
    <template #cell="{ column, record, text }">
      <template v-if="column.key === 'category'">
        <a-tag color="cyan">{{ toDisplayText(record.extra.category) }}</a-tag>
      </template>
      <template v-else-if="column.key === 'courseName'">
        {{ toDisplayText(record.extra.course_name) }}
      </template>
      <template v-else-if="column.key === 'fileType'">
        {{ toDisplayText(record.extra.file_type) }}
      </template>
      <template v-else-if="column.key === 'fileSize'">
        {{ toDisplayText(record.extra.file_size) }}
      </template>
      <template v-else>
        {{ toDisplayText(text) }}
      </template>
    </template>
  </CampusLifeManagementView>
</template>

<script setup lang="ts">
import CampusLifeManagementView from '../components/CampusLifeManagementView.vue'
import { toDisplayText } from '../support'

const columns = [
  { title: '资料名称', dataIndex: 'title', key: 'title', ellipsis: true },
  { title: '分类', key: 'category', width: 120 },
  { title: '课程', key: 'courseName', ellipsis: true },
  { title: '文件类型', key: 'fileType', width: 160 },
  { title: '文件大小', key: 'fileSize', width: 120 },
  { title: '审核状态', key: 'reviewStatus', width: 120 },
  { title: '上传者', dataIndex: 'publisher', key: 'publisher', width: 140 },
  { title: '提交时间', key: 'createdAt', width: 170 },
  { title: '操作', key: 'actions', width: 220, fixed: 'right' }
]
</script>
