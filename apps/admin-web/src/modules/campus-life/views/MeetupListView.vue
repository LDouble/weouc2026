<template>
  <CampusLifeManagementView
    title="组局活动"
    description="当前页面展示真实组局内容，可直接查看人数、时间、状态，并执行审核和上下线操作。"
    content-type="meetup"
    :columns="columns"
  >
    <template #cell="{ column, record, text }">
      <template v-if="column.key === 'category'">
        <a-tag color="green">{{ toDisplayText(record.extra.category) }}</a-tag>
      </template>
      <template v-else-if="column.key === 'startAt'">
        {{ formatDateTime(record.extra.start_at) }}
      </template>
      <template v-else-if="column.key === 'participants'">
        {{ formatParticipants(record.extra.joined_count, record.extra.max_participants) }}
      </template>
      <template v-else-if="column.key === 'businessStatus'">
        <a-tag :color="getContentStatusColor(record.extra.status)">
          {{ getContentStatusName(record.extra.status) }}
        </a-tag>
      </template>
      <template v-else>
        {{ toDisplayText(text) }}
      </template>
    </template>
  </CampusLifeManagementView>
</template>

<script setup lang="ts">
import CampusLifeManagementView from '../components/CampusLifeManagementView.vue'
import { formatDateTime, getContentStatusColor, getContentStatusName, toDisplayText } from '../support'

const columns = [
  { title: '活动名称', dataIndex: 'title', key: 'title', ellipsis: true },
  { title: '类型', key: 'category', width: 120 },
  { title: '地点', dataIndex: ['extra', 'location'], key: 'location', ellipsis: true },
  { title: '开始时间', key: 'startAt', width: 170 },
  { title: '参与人数', key: 'participants', width: 120 },
  { title: '业务状态', key: 'businessStatus', width: 120 },
  { title: '审核状态', key: 'reviewStatus', width: 120 },
  { title: '组织者', dataIndex: 'publisher', key: 'publisher', width: 140 },
  { title: '提交时间', key: 'createdAt', width: 170 },
  { title: '操作', key: 'actions', width: 220, fixed: 'right' }
]

function formatParticipants(joinedCount?: number, maxParticipants?: number): string {
  const joined = toDisplayText(joinedCount)
  const max = toDisplayText(maxParticipants)
  return `${joined} / ${max}`
}
</script>
