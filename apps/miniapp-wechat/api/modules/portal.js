import { get } from '~/api/request';

export function fetchPortalHome() {
  return get('/portal/home');
}

export function fetchPortalNotices(params = {}) {
  const { page = 1, pageSize = 20, keyword = '' } = params;
  const query = { page, pageSize };
  if (keyword && keyword.trim()) query.keyword = keyword.trim();
  return get('/portal/notices', query);
}

export function fetchPortalNoticeDetail(id) {
  return get(`/portal/notices/${id}`);
}
