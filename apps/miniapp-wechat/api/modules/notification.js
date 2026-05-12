import { get, post } from '~/api/request';

export function fetchNotificationList(params = {}) {
  const {
    page = 1,
    pageSize = 20,
    category = '',
    unreadOnly = false,
  } = params;

  const query = { page, pageSize };
  if (category && category.trim()) query.category = category.trim();
  if (unreadOnly) query.unread_only = true;

  return get('/notification/list', query);
}

export function fetchUnreadNotificationCount() {
  return get('/notification/unread-count');
}

export function markNotificationRead(messageId) {
  return post('/notification/read', { message_id: messageId });
}
