import { acceptErrand, fetchErrandList } from '~/api/modules/errand';
import { ERRAND_CATEGORIES } from '~/constants/campus';
import { formatDeadline, formatRelativeTime } from '~/utils/date';
import { unwrapPayload } from './shared';

const TONE_MAP = {
  parcel: 'blue',
  food: 'orange',
  seat: 'purple',
  print: 'green',
  other: 'gray',
};

const CATEGORY_LABEL_MAP = ERRAND_CATEGORIES.reduce((map, item) => {
  if (item.value && item.value !== 'all') {
    map[item.value] = item.label;
  }

  return map;
}, {});

function mapErrandItem(item = {}) {
  const status = item.status || '';
  const userRole = item.user_role || 'viewer';
  const isAccepted = Boolean(item.is_accepted) || status === 'accepted';
  const isPublisher = userRole === 'publisher';
  const isAcceptor = userRole === 'acceptor';
  let actionText = '接单';

  if (isPublisher) actionText = '我的发布';
  else if (isAcceptor) actionText = '已接单';
  else if (isAccepted) actionText = '已被接';
  else if (status === 'reviewing') actionText = '审核中';
  else if (status === 'rejected') actionText = '审核未通过';
  else if (status === 'offline') actionText = '已下线';
  else if (status === 'cancelled') actionText = '已取消';

  return {
    id: item.id,
    category: item.category,
    type: CATEGORY_LABEL_MAP[item.category] || '其他',
    title: item.title || '',
    desc: item.desc || '',
    routeStart: item.route_start || '',
    routeEnd: item.route_end || '',
    deadline: formatDeadline(item.deadline),
    reward: item.reward || '0',
    urgent: false,
    views: Number(item.views || 0),
    publisher: item.publisher || '匿名',
    publisherInitial: item.publisher_initial || '匿',
    publisherMeta: formatRelativeTime(item.created_at),
    time: formatRelativeTime(item.created_at),
    credit: '信用良好',
    creditShort: '良好',
    avatarTone: TONE_MAP[item.category] || 'gray',
    status,
    userRole,
    accepted: isAccepted,
    canAccept: !isPublisher && !isAccepted && status === 'published',
    actionText,
  };
}

export function getErrandCategoryList() {
  return ERRAND_CATEGORIES;
}

export async function loadErrandList(params = {}) {
  const response = await fetchErrandList(params);
  const data = unwrapPayload(response);
  const items = (data.list || []).map(mapErrandItem);

  return {
    items,
    total: Number(data.total || 0),
  };
}

export async function acceptErrandTask(task) {
  await acceptErrand(task.id);

  return {
    ...task,
    accepted: true,
    userRole: 'acceptor',
    canAccept: false,
    actionText: '已接单',
  };
}
