<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex flex-col gap-4 mb-6 lg:flex-row lg:items-end lg:justify-between">
      <div>
        <h2 class="text-lg font-semibold text-gray-800">{{ title }}</h2>
        <p class="mt-1 text-sm text-gray-500">{{ description }}</p>
      </div>

      <div class="flex flex-col gap-3 md:flex-row md:items-center">
        <a-input-search
          v-model:value="keyword"
          allow-clear
          class="w-full md:w-72"
          placeholder="按标题、摘要、发布者搜索"
          @search="handleFilterChange"
        />
        <a-select v-model:value="statusFilter" class="w-full md:w-40" @change="handleFilterChange">
          <a-select-option
            v-for="option in reviewStatusFilterOptions"
            :key="option.value"
            :value="option.value"
          >
            {{ option.label }}
          </a-select-option>
        </a-select>
        <a-button @click="resetFilters">重置</a-button>
      </div>
    </div>

    <a-table
      :columns="columns"
      :data-source="items"
      :loading="loading"
      :pagination="tablePagination"
      :row-key="(record: CampusLifeReviewRecord) => record.content_id"
      :scroll="{ x: 1180 }"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record, text }">
        <template v-if="column.key === 'reviewStatus'">
          <a-tag :color="getReviewStatusColor(record.review_status)">
            {{ getReviewStatusName(record.review_status) }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'createdAt'">
          {{ record.created_at_text }}
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-space wrap>
            <a-button size="small" @click="openDetail(record)">
              <component :is="icons.Eye" />
              详情
            </a-button>

            <a-popconfirm
              v-for="action in getReviewStatusActions(record.review_status)"
              :key="`${record.content_id}-${action.reviewStatus}`"
              :title="action.confirmText"
              ok-text="确认"
              cancel-text="取消"
              @confirm="changeReviewStatus(record, action)"
            >
              <a-button
                size="small"
                :type="action.type"
                :danger="action.danger"
                :loading="pendingActionKey === actionKey(record.content_id, action.reviewStatus)"
              >
                {{ action.label }}
              </a-button>
            </a-popconfirm>
          </a-space>
        </template>
        <template v-else>
          <slot name="cell" :column="column" :record="record" :text="text">
            <span class="whitespace-pre-wrap break-words">{{ toDisplayText(text) }}</span>
          </slot>
        </template>
      </template>
    </a-table>

    <a-drawer v-model:open="detailOpen" :title="detailTitle" :width="640" placement="right">
      <template #extra>
        <a-tag v-if="detailReviewStatus" :color="getReviewStatusColor(detailReviewStatus)">
          {{ getReviewStatusName(detailReviewStatus) }}
        </a-tag>
      </template>

      <a-skeleton v-if="detailLoading" active :paragraph="{ rows: 10 }" />

      <template v-else>
        <p v-if="detailSubtitle" class="mb-4 text-sm text-gray-500">
          {{ detailSubtitle }}
        </p>

        <div v-for="section in detailSections" :key="section.title" class="mb-6">
          <h3 class="mb-3 text-sm font-semibold text-gray-700">
            {{ section.title }}
          </h3>
          <a-descriptions bordered :column="1" size="small">
            <a-descriptions-item
              v-for="field in section.fields"
              :key="`${section.title}-${field.label}`"
              :label="field.label"
            >
              <div class="whitespace-pre-wrap break-all">
                {{ field.value }}
              </div>
            </a-descriptions-item>
          </a-descriptions>
        </div>
      </template>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { message } from 'ant-design-vue'
import { EyeOutlined as Eye } from '@ant-design/icons-vue'
import { campusLifeApi } from '@/api'
import {
  buildCampusLifeDetailSections,
  buildCampusLifeDetailTitle,
  fetchCampusLifeDetail,
  getReviewStatusActions,
  getReviewStatusColor,
  getReviewStatusName,
  mapCampusLifeReviewRecord,
  reviewStatusFilterOptions,
  toDisplayText,
  type CampusLifeContentType,
  type CampusLifeDetailSection,
  type CampusLifeListFilterStatus,
  type CampusLifeReviewRecord,
  type CampusLifeReviewStatus,
  type CampusLifeStatusAction
} from '../support'

const props = defineProps<{
  title: string
  description: string
  contentType: CampusLifeContentType
  columns: Array<Record<string, any>>
}>()

const icons = { Eye }
const statusFilter = ref<CampusLifeListFilterStatus>('all')
const keyword = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const loading = ref(false)
const items = ref<CampusLifeReviewRecord[]>([])
const pendingActionKey = ref('')

const detailOpen = ref(false)
const detailLoading = ref(false)
const detailTitle = ref('')
const detailSubtitle = ref('')
const detailReviewStatus = ref<CampusLifeReviewStatus | ''>('')
const detailSections = ref<CampusLifeDetailSection[]>([])
const activeDetailRecord = ref<CampusLifeReviewRecord | null>(null)

const tablePagination = computed(() => ({
  current: currentPage.value,
  pageSize: pageSize.value,
  total: total.value,
  showSizeChanger: true,
  showQuickJumper: true
}))

async function fetchItems() {
  loading.value = true

  try {
    const response = await campusLifeApi.listReviews({
      page: currentPage.value,
      pageSize: pageSize.value,
      keyword: keyword.value.trim() || undefined,
      content_type: props.contentType,
      review_status: statusFilter.value === 'all' ? undefined : statusFilter.value
    })

    const list = Array.isArray(response.data?.list) ? response.data.list : []
    items.value = list.map((item: Record<string, any>) => mapCampusLifeReviewRecord(item))
    total.value = Number(response.data?.total || items.value.length)
  } catch (error) {
    console.error('Failed to fetch campus life management list:', error)
    message.error('获取内容列表失败')
  } finally {
    loading.value = false
  }
}

async function openDetail(record: CampusLifeReviewRecord) {
  detailOpen.value = true
  detailLoading.value = true
  activeDetailRecord.value = record
  detailTitle.value = record.title
  detailSubtitle.value = `${record.publisher} · ${record.created_at_text}`
  detailReviewStatus.value = record.review_status

  try {
    const payload = await fetchCampusLifeDetail(props.contentType, record.content_id)
    detailTitle.value = buildCampusLifeDetailTitle(props.contentType, payload)
    detailSections.value = buildCampusLifeDetailSections(props.contentType, payload)
  } catch (error) {
    console.error('Failed to fetch campus life detail:', error)
    message.error('获取详情失败')
  } finally {
    detailLoading.value = false
  }
}

async function changeReviewStatus(record: CampusLifeReviewRecord, action: CampusLifeStatusAction) {
  const key = actionKey(record.content_id, action.reviewStatus)
  pendingActionKey.value = key

  try {
    await campusLifeApi.updateReviewStatus({
      content_type: props.contentType,
      content_id: record.content_id,
      review_status: action.reviewStatus
    })
    message.success(action.successText)

    if (activeDetailRecord.value?.content_id === record.content_id) {
      await openDetail({
        ...record,
        review_status: action.reviewStatus
      })
    }

    await fetchItems()
  } catch (error) {
    console.error('Failed to update campus life review status:', error)
    message.error('状态更新失败')
  } finally {
    pendingActionKey.value = ''
  }
}

function handleTableChange(pagination: Record<string, any>) {
  currentPage.value = Number(pagination.current || 1)
  pageSize.value = Number(pagination.pageSize || 20)
  fetchItems()
}

function handleFilterChange() {
  currentPage.value = 1
  fetchItems()
}

function resetFilters() {
  statusFilter.value = 'all'
  keyword.value = ''
  currentPage.value = 1
  fetchItems()
}

function actionKey(contentID: string, nextStatus: CampusLifeReviewStatus): string {
  return `${contentID}:${nextStatus}`
}

onMounted(() => {
  fetchItems()
})
</script>
