import { get } from '~/api/request';

export function fetchPortalHome() {
  return get('/portal/home', {}, { skipAuth: true });
}

export function fetchPortalNotices(params = {}) {
  const { page = 1, pageSize = 20, keyword } = params;
  const query = { page, pageSize };
  if (keyword) query.keyword = keyword;
  return get('/portal/notices', query, { skipAuth: true });
}

export function fetchPortalNoticeDetail(id) {
  return get(`/portal/notices/${id}`, {}, { skipAuth: true });
}
