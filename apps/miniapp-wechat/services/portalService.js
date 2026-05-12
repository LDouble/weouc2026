import { fetchPortalHome, fetchPortalNoticeDetail, fetchPortalNotices } from '~/api/modules/portal';
import { formatRelativeTime } from '~/utils/date';
import { unwrapPayload } from './shared';

const HERO_THEMES = ['rose', 'indigo', 'mint'];

function mapBanner(item = {}, index = 0) {
  return {
    id: item.id || `banner-${index}`,
    type: 'banner',
    tag: '校园推荐',
    title: item.title || '校园服务',
    desc: item.description || '',
    image: item.image_url || '',
    actionUrl: item.action_url || '',
    theme: HERO_THEMES[index % HERO_THEMES.length],
  };
}

function mapNotice(item = {}, index = 0) {
  return {
    id: item.id || `notice-${index}`,
    type: 'notice',
    tag: item.pinned ? '置顶公告' : '校园公告',
    title: item.title || '校园公告',
    desc: item.summary || '',
    publisher: item.publisher || '校园服务',
    timeLabel: formatRelativeTime(item.published_at || item.created_at),
    tags: item.tags || [],
    pinned: Boolean(item.pinned),
    theme: HERO_THEMES[index % HERO_THEMES.length],
    actionUrl: `/pages/message/index?noticeId=${item.id || ''}`,
    raw: item,
  };
}

export async function loadPortalHome() {
  const response = await fetchPortalHome();
  const data = unwrapPayload(response);
  return {
    banners: (data.banners || []).map(mapBanner),
    notices: (data.notices || []).map(mapNotice),
  };
}

export async function loadPortalNotices(params = {}) {
  const response = await fetchPortalNotices(params);
  const data = unwrapPayload(response);
  return {
    items: (data.list || []).map(mapNotice),
    total: Number(data.total || 0),
  };
}

export async function loadPortalNoticeDetail(id) {
  const response = await fetchPortalNoticeDetail(id);
  const item = unwrapPayload(response);
  return {
    ...mapNotice(item),
    content: item.content || item.summary || '',
  };
}
