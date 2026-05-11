import { fetchFeedList } from '~/api/modules/feed';
import { formatRelativeTime } from '~/utils/date';
import { splitAlternateColumns, unwrapPayload } from './shared';

function getFeedTargetUrl(feedType, id) {
  if (!id) return '';

  const detailRoutes = {
    market: '/pages/market/detail/index',
    errand: '/pages/errand/detail/index',
    lostFound: '/pages/lost-found/detail/index',
    meetup: '/pages/meetup/detail/index',
  };

  if (detailRoutes[feedType]) {
    return `${detailRoutes[feedType]}?id=${id}`;
  }

  const listRoutes = {
    resource: '/pages/resource/index',
    carpool: '/pages/carpool/index',
  };

  if (listRoutes[feedType]) {
    return `${listRoutes[feedType]}?insertId=${id}`;
  }

  return '';
}

function mapFeedItem(item = {}) {
  const extra = item.extra || {};
  const image = item.image || (extra.images && extra.images[0]) || '';

  return {
    id: item.id || '',
    feedType: item.feed_type || '',
    name: item.publisher || '校园用户',
    time: formatRelativeTime(item.created_at),
    badge: item.feed_type_label || '',
    badgeTone: 'hot',
    tone: 'indigo',
    avatarIcon: 'user',
    title: item.title || '',
    desc: item.desc || '',
    image,
    targetUrl: getFeedTargetUrl(item.feed_type, item.id),
    likes: Number(extra.likes || item.likes || 0),
    comments: Number(extra.comments || item.comments || 0),
  };
}

export async function loadHomeFeeds(params = {}) {
  const response = await fetchFeedList(params);
  const data = unwrapPayload(response);
  const items = (data.list || []).map(mapFeedItem);

  return {
    items,
    columns: splitAlternateColumns(items),
    total: Number(data.total || 0),
  };
}
