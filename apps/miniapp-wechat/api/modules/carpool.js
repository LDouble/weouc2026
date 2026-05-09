import { get, post } from '~/api/request';

export function fetchCarpoolList(params = {}) {
  const { category, keyword, page = 1, pageSize = 20 } = params;
  const query = { page, pageSize };
  if (category) query.category = category;
  if (keyword) query.keyword = keyword;
  return get('/carpool/list', query);
}

export function fetchCarpoolDetail(id) {
  return get(`/carpool/detail/${id}`);
}

export function publishCarpool(data) {
  return post('/carpool/publish', data);
}
