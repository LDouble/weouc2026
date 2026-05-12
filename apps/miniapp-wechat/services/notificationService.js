import {
  fetchNotificationList,
  fetchUnreadNotificationCount,
  markNotificationRead,
} from '~/api/modules/notification';
import { formatRelativeTime } from '~/utils/date';
import { unwrapPayload } from './shared';

const CATEGORY_LABELS = {
  system: '系统通知',
  activity: '活动通知',
  review: '审核通知',
  course: '教务通知',
};

function mapCategoryLabel(category) {
  return CATEGORY_LABELS[category] || '站内消息';
}

function mapNotificationItem(raw = {}) {
  return {
    id: raw.id || '',
    title: raw.title || '站内通知',
    content: raw.content || '',
    category: raw.category || 'system',
    categoryLabel: mapCategoryLabel(raw.category),
    createdAt: raw.created_at || '',
    timeLabel: formatRelativeTime(raw.created_at),
    read: Boolean(raw.read),
    actionUrl: raw.action_url || '',
  };
}

export async function loadNotificationList(params = {}) {
  const response = await fetchNotificationList(params);
  const payload = unwrapPayload(response);
  const list = Array.isArray(payload.list) ? payload.list.map(mapNotificationItem) : [];

  return {
    list,
    total: Number(payload.total || 0),
    page: Number(payload.page || params.page || 1),
    pageSize: Number(payload.pageSize || params.pageSize || 20),
  };
}

export async function loadNotificationUnreadCount() {
  const response = await fetchUnreadNotificationCount();
  const payload = unwrapPayload(response);
  return Number(payload.count || 0);
}

export async function readNotification(messageId) {
  if (!messageId) return;
  await markNotificationRead(messageId);
}
