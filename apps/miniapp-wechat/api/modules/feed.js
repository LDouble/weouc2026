import { get } from '~/api/request';

export function fetchFeedList(params = {}) {
  const { feed_types, keyword, page = 1, pageSize = 20 } = params;
  const query = { page, pageSize };
  if (feed_types) query.feed_types = feed_types;
  if (keyword) query.keyword = keyword;
  return get('/feed/list', query);
}

export function fetchReviewStatus(id) {
  return get(`/feed/review-status/${id}`);
}
