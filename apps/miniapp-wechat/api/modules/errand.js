import { get, post } from '~/api/request';

export function fetchErrandList(params = {}) {
  const { category, keyword, page = 1, pageSize = 20 } = params;
  const query = { page, pageSize };
  if (category) query.category = category;
  if (keyword) query.keyword = keyword;
  return get('/errand/list', query);
}

export function fetchErrandDetail(id) {
  return get(`/errand/detail/${id}`);
}

export function publishErrand(data) {
  return post('/errand/publish', data);
}

export function acceptErrand(taskId) {
  return post('/errand/accept', { task_id: taskId });
}

export function cancelErrandPublish(taskId) {
  return post('/errand/cancel-publish', { task_id: taskId });
}

export function cancelErrandAccept(taskId) {
  return post('/errand/cancel-accept', { task_id: taskId });
}
