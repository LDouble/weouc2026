import { fetchMeetupDetail, fetchMeetupList } from '~/api/modules/meetup';
import { MEETUP_CATEGORIES } from '~/constants/campus';
import { formatDateTime, formatRelativeTime } from '~/utils/date';
import { dedupeById, getStatusDisplay, getReviewStatus, unwrapPayload } from './shared';

const CATEGORY_LABEL_MAP = MEETUP_CATEGORIES.reduce((map, item) => {
  if (item.value && item.value !== 'all') {
    map[item.value] = item.label;
  }
  return map;
}, {});

const MEETUP_STATUS_OVERRIDES = {
  published: { label: '报名中' },
  open: { label: '报名中' },
};

function getActionText(item) {
  if (item.canDelete) return '取消组局';
  if (item.canCancelJoin) return '取消报名';
  if (item.canJoin) return '立即报名';
  return '查看详情';
}

function mapMeetupItem(item = {}) {
  const extra = item.extra || {};
  const category = item.category || extra.category || 'other';
  const status = getReviewStatus(item) || 'open';
  const statusDisplay = getStatusDisplay(status, MEETUP_STATUS_OVERRIDES);
  const joinedCount = Number(item.joined_count || extra.joined_count || 1);
  const remainingSeats = Number(item.remaining_seats || extra.remaining_seats || 0);
  const startAt = item.start_at || extra.start_at || '';
  const deadlineAt = item.deadline_at || extra.deadline_at || '';
  const contactFields = normalizeContactFields(item);
  const canJoin = Boolean(item.can_join);
  const canCancelJoin = Boolean(item.can_cancel_join);
  const canDelete = Boolean(item.can_delete);

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
    statusText: statusDisplay.label,
    statusTone: statusDisplay.tone,
    publisher: item.publisher || '',
    publisherInitial: item.publisher_initial || (item.publisher ? item.publisher.charAt(0) : ''),
    publisherMeta: formatRelativeTime(item.created_at),
    createdAt: item.created_at || '',
    userRole: item.user_role || 'viewer',
    isOwner: Boolean(item.is_owner),
    joined: Boolean(item.joined),
    canJoin,
    canCancelJoin,
    canEdit: Boolean(item.can_edit),
    canDelete,
    actionText: getActionText({
      canDelete,
      canCancelJoin,
      canJoin,
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
