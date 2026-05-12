import { fetchMeetupDetail, fetchMeetupList } from '~/api/modules/meetup';
import { MEETUP_CATEGORIES } from '~/constants/campus';
import { formatDateTime, formatRelativeTime } from '~/utils/date';
import { dedupeById, normalizeContactFields, unwrapPayload } from './shared';

const CATEGORY_LABEL_MAP = MEETUP_CATEGORIES.reduce((map, item) => {
  if (item.value && item.value !== 'all') {
    map[item.value] = item.label;
  }
  return map;
}, {});

const STATUS_META = {
  reviewing: { text: '审核中', tone: 'amber' },
  published: { text: '报名中', tone: 'green' },
  open: { text: '报名中', tone: 'green' },
  rejected: { text: '审核未通过', tone: 'red' },
  offline: { text: '已下架', tone: 'red' },
  full: { text: '人数已满', tone: 'purple' },
  cancelled: { text: '已取消', tone: 'red' },
};

function resolveStatusMeta(status) {
  return STATUS_META[status] || STATUS_META.published;
}

function getActionText(item) {
  if (item.canCancelPublish) return '取消组局';
  if (item.canCancelJoin) return '取消报名';
  if (item.canJoin) return '立即报名';
  if (item.userRole === 'participant') return '已报名';
  if (item.status === 'reviewing') return '等待审核';
  if (item.status === 'full') return '人数已满';
  if (item.status === 'cancelled') return '组局已取消';
  return '查看详情';
}

function mapMeetupItem(item = {}) {
  const extra = item.extra || {};
  const category = item.category || extra.category || 'other';
  const status = item.status || extra.status || 'open';
  const statusMeta = resolveStatusMeta(status);
  const joinedCount = Number(item.joined_count || extra.joined_count || 1);
  const remainingSeats = Number(item.remaining_seats || extra.remaining_seats || 0);
  const startAt = item.start_at || extra.start_at || '';
  const deadlineAt = item.deadline_at || extra.deadline_at || '';
  const contactFields = normalizeContactFields(item);

  return {
    id: item.id || '',
    category,
    categoryLabel: CATEGORY_LABEL_MAP[category] || '其他',
    title: item.title || '',
    desc: item.desc || '',
    location: item.location || extra.location || '',
    startAt,
    deadlineAt,
    startLabel: formatDateTime(startAt),
    deadlineLabel: deadlineAt ? `截止 ${formatDateTime(deadlineAt)}` : '',
    maxParticipants: Number(item.max_participants || extra.max_participants || 0),
    joinedCount,
    remainingSeats,
    feeText: item.fee_text || extra.fee_text || '',
    tags: item.tags || extra.tags || [],
    contact: contactFields.contact,
    canViewContact: contactFields.canViewContact,
    contactHiddenReason: contactFields.contactHiddenReason,
    status,
    statusText: statusMeta.text,
    statusTone: statusMeta.tone,
    publisher: item.publisher || '',
    publisherInitial: item.publisher_initial || (item.publisher ? item.publisher.charAt(0) : ''),
    publisherMeta: formatRelativeTime(item.created_at),
    createdAt: item.created_at || '',
    userRole: item.user_role || 'viewer',
    joined: Boolean(item.joined),
    canJoin: Boolean(item.can_join),
    canCancelJoin: Boolean(item.can_cancel_join),
    canCancelPublish: Boolean(item.can_cancel_publish),
    actionText: getActionText({
      canCancelPublish: Boolean(item.can_cancel_publish),
      canCancelJoin: Boolean(item.can_cancel_join),
      canJoin: Boolean(item.can_join),
      userRole: item.user_role || 'viewer',
      status,
    }),
  };
}

export function getMeetupCategoryList() {
  return MEETUP_CATEGORIES;
}

export async function loadMeetupList(params = {}) {
  const response = await fetchMeetupList(params);
  const data = unwrapPayload(response);
  return {
    items: (data.list || []).map(mapMeetupItem),
    total: Number(data.total || 0),
  };
}

export async function loadMeetupDetail(id) {
  const response = await fetchMeetupDetail(id);
  return mapMeetupItem(unwrapPayload(response));
}

export function mergeMeetupItems(currentItems, nextItems) {
  return dedupeById([].concat(currentItems || [], nextItems || []));
}
