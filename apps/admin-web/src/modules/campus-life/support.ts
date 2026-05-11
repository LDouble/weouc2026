import { campusLifeApi } from '@/api'

export type CampusLifeContentType = 'market' | 'errand' | 'resource' | 'lostFound' | 'meetup'
export type CampusLifeReviewStatus = 'reviewing' | 'published' | 'rejected' | 'offline'
export type CampusLifeListFilterStatus = 'all' | CampusLifeReviewStatus

export interface CampusLifeReviewRecord {
  key: string
  content_id: string
  content_type: CampusLifeContentType
  title: string
  desc: string
  publisher: string
  review_status: CampusLifeReviewStatus
  created_at: string
  created_at_text: string
  extra: Record<string, any>
}

export interface CampusLifeDetailField {
  label: string
  value: string
}

export interface CampusLifeDetailSection {
  title: string
  fields: CampusLifeDetailField[]
}

export interface CampusLifeStatusAction {
  label: string
  reviewStatus: Exclude<CampusLifeReviewStatus, 'reviewing'>
  successText: string
  confirmText: string
  type?: 'primary' | 'default'
  danger?: boolean
}

export const reviewStatusFilterOptions: Array<{
  label: string
  value: CampusLifeListFilterStatus
}> = [
  { label: '全部状态', value: 'all' },
  { label: '待审核', value: 'reviewing' },
  { label: '已发布', value: 'published' },
  { label: '已拒绝', value: 'rejected' },
  { label: '已下线', value: 'offline' }
]

export function mapCampusLifeReviewRecord(raw: Record<string, any>): CampusLifeReviewRecord {
  const reviewStatus = normalizeReviewStatus(raw.review_status)
  const contentID = String(raw.content_id || raw.id || '')

  return {
    key: contentID,
    content_id: contentID,
    content_type: raw.content_type as CampusLifeContentType,
    title: String(raw.title || '-'),
    desc: String(raw.desc || ''),
    publisher: String(raw.publisher || '-'),
    review_status: reviewStatus,
    created_at: String(raw.created_at || ''),
    created_at_text: formatDateTime(raw.created_at),
    extra: isPlainObject(raw.extra) ? raw.extra : {}
  }
}

export function getReviewStatusName(status: string): string {
  const names: Record<string, string> = {
    reviewing: '待审核',
    published: '已发布',
    rejected: '已拒绝',
    offline: '已下线'
  }
  return names[normalizeReviewStatus(status)] || status || '-'
}

export function getReviewStatusColor(status: string): string {
  const colors: Record<string, string> = {
    reviewing: 'processing',
    published: 'success',
    rejected: 'error',
    offline: 'default'
  }
  return colors[normalizeReviewStatus(status)] || 'default'
}

export function getContentStatusName(status: string): string {
  const names: Record<string, string> = {
    reviewing: '待审核',
    published: '已发布',
    rejected: '已拒绝',
    offline: '已下线',
    accepted: '已接单',
    cancelled: '已取消',
    open: '报名中',
    full: '已满员'
  }
  const normalized = String(status || '').trim().toLowerCase()
  return names[normalized] || status || '-'
}

export function getContentStatusColor(status: string): string {
  const colors: Record<string, string> = {
    reviewing: 'processing',
    published: 'success',
    rejected: 'error',
    offline: 'default',
    accepted: 'blue',
    cancelled: 'default',
    open: 'green',
    full: 'orange'
  }
  const normalized = String(status || '').trim().toLowerCase()
  return colors[normalized] || 'default'
}

export function getCampusLifeTypeName(type: string): string {
  const names: Record<string, string> = {
    market: '二手交易',
    errand: '跑腿服务',
    resource: '资料共享',
    lostFound: '失物招领',
    meetup: '组局活动'
  }
  return names[type] || type || '内容'
}

export function toDisplayText(value: unknown): string {
  if (Array.isArray(value)) {
    const parts = value
      .map((item) => toDisplayText(item))
      .filter((item) => item !== '-')
    return parts.length > 0 ? parts.join('、') : '-'
  }

  if (typeof value === 'boolean') {
    return value ? '是' : '否'
  }

  if (value === null || value === undefined) {
    return '-'
  }

  const text = String(value).trim()
  return text === '' ? '-' : text
}

export function formatDateTime(value: unknown): string {
  const text = typeof value === 'string' ? value.trim() : ''
  if (text === '') {
    return '-'
  }

  const date = new Date(text)
  if (Number.isNaN(date.getTime())) {
    return text
  }

  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

export function getReviewStatusActions(status: string): CampusLifeStatusAction[] {
  switch (normalizeReviewStatus(status)) {
    case 'reviewing':
      return [
        {
          label: '通过',
          reviewStatus: 'published',
          successText: '内容已发布',
          confirmText: '确认通过当前内容并立即发布吗？',
          type: 'primary'
        },
        {
          label: '拒绝',
          reviewStatus: 'rejected',
          successText: '内容已拒绝',
          confirmText: '确认拒绝当前内容吗？',
          danger: true
        }
      ]
    case 'published':
      return [
        {
          label: '下线',
          reviewStatus: 'offline',
          successText: '内容已下线',
          confirmText: '确认下线当前内容吗？',
          danger: true
        }
      ]
    case 'rejected':
    case 'offline':
      return [
        {
          label: '重新发布',
          reviewStatus: 'published',
          successText: '内容已重新发布',
          confirmText: '确认重新发布当前内容吗？',
          type: 'primary'
        }
      ]
    default:
      return []
  }
}

export async function fetchCampusLifeDetail(contentType: CampusLifeContentType, contentID: string): Promise<Record<string, any>> {
  switch (contentType) {
    case 'market': {
      const response = await campusLifeApi.getMarketDetail(contentID)
      return response.data || {}
    }
    case 'errand': {
      const response = await campusLifeApi.getErrandDetail(contentID)
      return response.data || {}
    }
    case 'resource': {
      const response = await campusLifeApi.getResourceDetail(contentID)
      return response.data || {}
    }
    case 'lostFound': {
      const response = await campusLifeApi.getLostFoundDetail(contentID)
      return response.data || {}
    }
    case 'meetup': {
      const response = await campusLifeApi.getMeetupDetail(contentID)
      return response.data || {}
    }
    default:
      return {}
  }
}

export function buildCampusLifeDetailTitle(contentType: CampusLifeContentType, payload: Record<string, any>): string {
  const root = detailRoot(contentType, payload)
  const title = toDisplayText(root.title)
  if (title !== '-') {
    return title
  }
  return `${getCampusLifeTypeName(contentType)}详情`
}

export function buildCampusLifeDetailSections(
  contentType: CampusLifeContentType,
  payload: Record<string, any>
): CampusLifeDetailSection[] {
  switch (contentType) {
    case 'market':
      return buildMarketSections(payload)
    case 'errand':
      return buildErrandSections(payload)
    case 'resource':
      return buildResourceSections(payload)
    case 'lostFound':
      return buildLostFoundSections(payload)
    case 'meetup':
      return buildMeetupSections(payload)
    default:
      return []
  }
}

function buildMarketSections(payload: Record<string, any>): CampusLifeDetailSection[] {
  const root = detailRoot('market', payload)
  const extra = root.extra || {}

  return [
    section('基础信息', [
      ['标题', root.title],
      ['描述', root.desc],
      ['发布者', root.publisher],
      ['发布时间', formatDateTime(root.created_at)],
      ['审核状态', getContentStatusName(root.status)]
    ]),
    section('交易信息', [
      ['分类', extra.category],
      ['价格', extra.price],
      ['原价', extra.original_price],
      ['成色', extra.condition],
      ['交易方式', extra.trade_mode],
      ['联系方式', extra.contact],
      ['图片', toDisplayText(extra.images)]
    ])
  ]
}

function buildErrandSections(payload: Record<string, any>): CampusLifeDetailSection[] {
  const root = detailRoot('errand', payload)
  const extra = root.extra || {}

  return [
    section('基础信息', [
      ['标题', root.title],
      ['描述', root.desc],
      ['发布者', root.publisher],
      ['发布时间', formatDateTime(root.created_at)],
      ['审核/业务状态', getContentStatusName(root.status)],
      ['我的角色', payload.user_role]
    ]),
    section('任务信息', [
      ['分类', extra.category || root.category],
      ['出发地', extra.route_start || root.route_start],
      ['目的地', extra.route_end || root.route_end],
      ['截止时间', formatDateTime(extra.deadline || root.deadline)],
      ['报酬', extra.reward || root.reward],
      ['联系方式', extra.contact || root.contact],
      ['图片', toDisplayText(extra.images || root.images)]
    ])
  ]
}

function buildResourceSections(payload: Record<string, any>): CampusLifeDetailSection[] {
  const root = detailRoot('resource', payload)
  const extra = root.extra || {}
  const files = Array.isArray(extra.files)
    ? extra.files.map((file: Record<string, any>) => {
        const name = toDisplayText(file?.name)
        const fileType = toDisplayText(file?.file_type)
        const fileSize = toDisplayText(file?.file_size)
        const url = toDisplayText(file?.url)
        return [name, fileType, fileSize, url].filter((value) => value !== '-').join(' | ')
      })
    : []

  return [
    section('基础信息', [
      ['标题', root.title],
      ['描述', root.desc],
      ['发布者', root.publisher],
      ['发布时间', formatDateTime(root.created_at)],
      ['审核状态', getContentStatusName(root.status)]
    ]),
    section('资料信息', [
      ['分类', extra.category],
      ['课程名称', extra.course_name],
      ['文件类型', extra.file_type],
      ['文件大小', extra.file_size],
      ['联系方式', extra.contact],
      ['文件清单', toDisplayText(files)],
      ['下载链接', extra.download_url]
    ])
  ]
}

function buildLostFoundSections(payload: Record<string, any>): CampusLifeDetailSection[] {
  const root = detailRoot('lostFound', payload)
  const extra = root.extra || {}
  const typeName = extra.type === 'found' ? '失物招领' : '寻物启事'

  return [
    section('基础信息', [
      ['标题', root.title],
      ['描述', root.desc],
      ['发布者', root.publisher],
      ['发布时间', formatDateTime(root.created_at)],
      ['审核状态', getContentStatusName(root.status)]
    ]),
    section('线索信息', [
      ['类型', typeName],
      ['分类', extra.category],
      ['地点', extra.location],
      ['发生时间', formatDateTime(extra.event_time)],
      ['物品特征', extra.item_feature],
      ['联系方式', extra.contact]
    ])
  ]
}

function buildMeetupSections(payload: Record<string, any>): CampusLifeDetailSection[] {
  const root = detailRoot('meetup', payload)
  const extra = root.extra || {}

  return [
    section('基础信息', [
      ['标题', root.title],
      ['描述', root.desc],
      ['发布者', root.publisher],
      ['发布时间', formatDateTime(root.created_at)],
      ['审核/业务状态', getContentStatusName(root.status)]
    ]),
    section('活动信息', [
      ['分类', extra.category || root.category],
      ['地点', extra.location || root.location],
      ['开始时间', formatDateTime(extra.start_at || root.start_at)],
      ['报名截止', formatDateTime(extra.deadline_at || root.deadline_at)],
      ['人数上限', extra.max_participants || root.max_participants],
      ['已报名人数', extra.joined_count || root.joined_count],
      ['剩余名额', extra.remaining_seats || root.remaining_seats],
      ['费用说明', extra.fee_text || root.fee_text],
      ['标签', toDisplayText(extra.tags || root.tags)],
      ['联系方式', extra.contact || root.contact]
    ])
  ]
}

function normalizeReviewStatus(status: unknown): CampusLifeReviewStatus {
  switch (String(status || '').trim().toLowerCase()) {
    case 'reviewing':
      return 'reviewing'
    case 'rejected':
      return 'rejected'
    case 'offline':
      return 'offline'
    default:
      return 'published'
  }
}

function detailRoot(contentType: CampusLifeContentType, payload: Record<string, any>): Record<string, any> {
  if (contentType === 'errand' && isPlainObject(payload.item)) {
    return payload.item
  }
  return payload
}

function section(title: string, rows: Array<[string, unknown]>): CampusLifeDetailSection {
  return {
    title,
    fields: rows.map(([label, value]) => ({
      label,
      value: toDisplayText(value)
    }))
  }
}

function isPlainObject(value: unknown): value is Record<string, any> {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}
