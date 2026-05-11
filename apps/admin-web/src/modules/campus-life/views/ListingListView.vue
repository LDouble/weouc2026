<template>
  <CampusLifeManagementView
    title="二手交易"
    description="当前页面使用真实二手审核数据，支持查看价格、成色以及执行审核与上下线操作。"
    content-type="market"
    :columns="columns"
  >
    <template #cell="{ column, record, text }">
      <template v-if="column.key === 'category'">
        <a-tag color="purple">{{ toDisplayText(record.extra.category) }}</a-tag>
      </template>
      <template v-else-if="column.key === 'price'">
        <span class="font-medium text-rose-500">{{ toDisplayText(record.extra.price) }}</span>
      </template>
      <template v-else-if="column.key === 'condition'">
        {{ toDisplayText(record.extra.condition) }}
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
  { title: '物品名称', dataIndex: 'title', key: 'title', ellipsis: true },
  { title: '分类', key: 'category', width: 120 },
  { title: '价格', key: 'price', width: 120 },
  { title: '成色', key: 'condition', width: 120 },
  { title: '审核状态', key: 'reviewStatus', width: 120 },
  { title: '卖家', dataIndex: 'publisher', key: 'publisher', width: 140 },
  { title: '提交时间', key: 'createdAt', width: 170 },
  { title: '操作', key: 'actions', width: 220, fixed: 'right' }
]
</script>
