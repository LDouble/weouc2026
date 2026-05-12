import { get, post } from '~/api/request';

export function fetchLostFoundList(params = {}) {
  const { category, keyword, type, page = 1, pageSize = 20 } = params;
  const query = { page, pageSize };
  if (category) query.category = category;
  if (keyword) query.keyword = keyword;
  if (type) query.type = type;
  return get('/lostFound/list', query);
}

export function fetchLostFoundDetail(id) {
  return get(`/lostFound/detail/${id}`);
}

export function publishLostFound(data) {
  return post('/lostFound/publish', data);
}

export function deleteLostFound(id) {
  return post(`/lostFound/delete/${id}`);
}

export function resolveLostFound(id) {
  return post(`/lostFound/resolve/${id}`);
}
