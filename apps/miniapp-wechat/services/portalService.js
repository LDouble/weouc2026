import { fetchPortalHome } from '~/api/modules/portal';
import { formatRelativeTime } from '~/utils/date';
import { unwrapPayload } from './shared';

const HERO_THEMES = ['rose', 'indigo', 'mint'];

export const DEFAULT_HERO_CARDS = [
  { id: 'cet', tag: '重要通知', title: '四六级考试报名进行中', desc: '请于本周五前完成报名缴费', theme: 'rose' },
  { id: 'job', tag: '职业发展', title: '春季校园招聘会', desc: '300+ 企业现场招聘', theme: 'indigo' },
  { id: 'seat', tag: '学习资源', title: '图书馆新增自习座位', desc: '支持线上预约，先到先得', theme: 'mint' },
];

function buildNoticeDesc(item = {}) {
  if (item.summary) return item.summary;
  const publisher = item.publisher || '校园平台';
  const relativeTime = formatRelativeTime(item.published_at || item.created_at);
  if (relativeTime) return `${publisher} · ${relativeTime}`;
  return publisher;
}

function mapNoticeToHeroCard(item = {}, index = 0) {
  return {
    id: item.id || `notice-${index}`,
    tag: item.pinned ? '置顶公告' : '校园公告',
    title: item.title || '校园通知',
    desc: buildNoticeDesc(item),
    theme: HERO_THEMES[index % HERO_THEMES.length],
  };
}

export async function loadPortalHeroCards() {
  const response = await fetchPortalHome();
  const payload = unwrapPayload(response);
  const notices = Array.isArray(payload.notices) ? payload.notices : [];

  if (!notices.length) {
    return DEFAULT_HERO_CARDS;
  }

  return notices.slice(0, 3).map(mapNoticeToHeroCard);
}
