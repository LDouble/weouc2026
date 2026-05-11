<template>
  <CampusLifeManagementView
    title="跑腿服务"
    description="当前页面直接展示真实跑腿审核数据，可查看详情并执行通过、拒绝、下线和重新发布操作。"
    content-type="errand"
    :columns="columns"
  >
    <template #cell="{ column, record, text }">
      <template v-if="column.key === 'category'">
        <a-tag color="blue">{{ toDisplayText(record.extra.category) }}</a-tag>
      </template>
      <template v-else-if="column.key === 'route'">
        {{ formatRoute(record.extra.route_start, record.extra.route_end) }}
      </template>
      <template v-else-if="column.key === 'reward'">
        <span class="font-medium text-rose-500">{{ toDisplayText(record.extra.reward) }}</span>
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
import { getContentStatusColor, getContentStatusName, toDisplayText } from '../support'

const columns = [
  { title: '标题', dataIndex: 'title', key: 'title', ellipsis: true },
  { title: '分类', key: 'category', width: 120 },
  { title: '路线', key: 'route', ellipsis: true },
  { title: '报酬', key: 'reward', width: 120 },
  { title: '业务状态', key: 'businessStatus', width: 120 },
  { title: '审核状态', key: 'reviewStatus', width: 120 },
  { title: '发布者', dataIndex: 'publisher', key: 'publisher', width: 140 },
  { title: '提交时间', key: 'createdAt', width: 170 },
  { title: '操作', key: 'actions', width: 220, fixed: 'right' }
]

function formatRoute(routeStart?: string, routeEnd?: string): string {
  const start = toDisplayText(routeStart)
  const end = toDisplayText(routeEnd)

  if (start === '-' && end === '-') {
    return '-'
  }

  return `${start} -> ${end}`
}
</script>
