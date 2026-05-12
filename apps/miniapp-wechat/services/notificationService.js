import {
  fetchNotificationList,
  fetchUnreadNotificationCount,
  markNotificationRead,
} from '~/api/modules/notification';
import { formatRelativeTime } from '~/utils/date';
import { unwrapPayload } from './shared';

const CATEGORY_META = {
  system: { label: '系统', icon: 'notification', tone: 'indigo' },
  audit: { label: '审核', icon: 'check-circle', tone: 'green' },
  campus_life: { label: '校园生活', icon: 'app', tone: 'orange' },
  portal: { label: '公告', icon: 'sound', tone: 'rose' },
};

function normalizeActionUrl(url = '') {
  if (!url) return '';
  if (url.startsWith('/pages/')) return url;
  return '';
}

function mapNotificationItem(item = {}) {
  const category = item.category || 'system';
  const meta = CATEGORY_META[category] || CATEGORY_META.system;

  return {
    id: item.id || '',
    title: item.title || '校园通知',
    content: item.content || '',
    category,
    categoryLabel: meta.label,
    icon: meta.icon,
    tone: meta.tone,
    publisher: item.publisher || '校园服务',
    timeLabel: formatRelativeTime(item.created_at),
    read: Boolean(item.read),
    actionUrl: normalizeActionUrl(item.action_url || ''),
    raw: item,
  };
}

export async function loadNotifications(params = {}) {
  const response = await fetchNotificationList(params);
  const data = unwrapPayload(response);
  return {
    items: (data.list || []).map(mapNotificationItem),
    total: Number(data.total || 0),
    page: Number(data.page || params.page || 1),
    pageSize: Number(data.pageSize || params.pageSize || 20),
  };
}

export async function loadUnreadNotificationCount() {
  const response = await fetchUnreadNotificationCount();
  const data = unwrapPayload(response);
  return Number(data.count || 0);
}

export async function readNotification(messageId) {
  await markNotificationRead(messageId);
}

export function publishUnreadCount(count) {
  const app = getApp();
  app.globalData.unreadNum = count;
  if (app.eventBus) app.eventBus.emit('unread-num-change', count);
}
