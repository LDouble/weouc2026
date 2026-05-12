import { get, post } from '~/api/request';

export function fetchMarketList(params = {}) {
  const { category, keyword, page = 1, pageSize = 20 } = params;
  const query = { page, pageSize };
  if (category) query.category = category;
  if (keyword) query.keyword = keyword;
  return get('/market/list', query);
}

export function fetchMarketDetail(id) {
  return get(`/market/detail/${id}`);
}

export function publishMarket(data) {
  return post('/market/publish', data);
}

export function favoriteMarket(productId, action) {
  return post('/market/favorite', { product_id: productId, action });
}

export function deleteMarket(id) {
  return post(`/market/delete/${id}`);
}
