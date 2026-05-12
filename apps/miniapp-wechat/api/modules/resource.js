import { get, post } from '~/api/request';

export function fetchResourceList(params = {}) {
  const { category, keyword, page = 1, pageSize = 20 } = params;
  const query = { page, pageSize };
  if (category) query.category = category;
  if (keyword) query.keyword = keyword;
  return get('/resource/list', query);
}

export function fetchResourceDetail(id) {
  return get(`/resource/detail/${id}`);
}

export function publishResource(data) {
  return post('/resource/publish', data);
}

export function deleteResource(id) {
  return post(`/resource/delete/${id}`);
}

export function favoriteResource(resourceId, action) {
  return post('/resource/favorite', { resource_id: resourceId, action });
}
