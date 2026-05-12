import { get, post } from '~/api/request';

export function fetchNotificationList(params = {}) {
  const {
    page = 1,
    pageSize = 20,
    category,
    unreadOnly,
    unread_only: unreadOnlyRaw,
  } = params;
  const query = { page, pageSize };
  if (category) query.category = category;
  const onlyUnread = unreadOnlyRaw !== undefined ? unreadOnlyRaw : unreadOnly;
  if (onlyUnread !== undefined) query.unread_only = Boolean(onlyUnread);
  return get('/notification/list', query);
}

export function fetchUnreadNotificationCount() {
  return get('/notification/unread-count');
}

export function markNotificationRead(messageId) {
  return post('/notification/read', { message_id: messageId });
}
