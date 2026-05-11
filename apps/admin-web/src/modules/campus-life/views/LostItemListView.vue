<template>
  <CampusLifeManagementView
    title="失物招领"
    description="当前页面展示真实失物招领审核数据，可查看类型、地点、发生时间并执行状态操作。"
    content-type="lostFound"
    :columns="columns"
  >
    <template #cell="{ column, record, text }">
      <template v-if="column.key === 'type'">
        <a-tag :color="record.extra.type === 'found' ? 'green' : 'orange'">
          {{ record.extra.type === 'found' ? '失物招领' : '寻物启事' }}
        </a-tag>
      </template>
      <template v-else-if="column.key === 'category'">
        <a-tag color="orange">{{ toDisplayText(record.extra.category) }}</a-tag>
      </template>
      <template v-else-if="column.key === 'eventTime'">
        {{ formatDateTime(record.extra.event_time) }}
      </template>
      <template v-else>
        {{ toDisplayText(text) }}
      </template>
    </template>
  </CampusLifeManagementView>
</template>

<script setup lang="ts">
import CampusLifeManagementView from '../components/CampusLifeManagementView.vue'
import { formatDateTime, toDisplayText } from '../support'

const columns = [
  { title: '类型', key: 'type', width: 120 },
  { title: '物品名称', dataIndex: 'title', key: 'title', ellipsis: true },
  { title: '分类', key: 'category', width: 120 },
  { title: '地点', dataIndex: ['extra', 'location'], key: 'location', ellipsis: true },
  { title: '发生时间', key: 'eventTime', width: 170 },
  { title: '审核状态', key: 'reviewStatus', width: 120 },
  { title: '发布者', dataIndex: 'publisher', key: 'publisher', width: 140 },
  { title: '提交时间', key: 'createdAt', width: 170 },
  { title: '操作', key: 'actions', width: 220, fixed: 'right' }
]
</script>
