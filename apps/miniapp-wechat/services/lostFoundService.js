import { fetchLostFoundList } from '~/api/modules/lostFound';
import { LOST_FOUND_CATEGORIES } from '~/constants/campus';
import { unwrapPayload } from './shared';

const CATEGORY_LABEL_MAP = LOST_FOUND_CATEGORIES.reduce((map, item) => {
  if (item.value && item.value !== 'all') {
    map[item.value] = item.label;
  }

  return map;
}, {});

const AVATAR_COLORS = ['amber', 'rose', 'blue', 'purple', 'green'];

function stableIndex(value, modulo) {
  const text = `${value || ''}`;
  let total = 0;

  for (let index = 0; index < text.length; index += 1) {
    total += text.charCodeAt(index);
  }

  return total % modulo;
}

function mapLostFoundItem(raw = {}) {
  const extra = raw.extra || {};
  const type = extra.type || raw.feed_type || 'lost';
  const isLost = type === 'lost';
  const category = extra.category || '';
  const categoryLabel = CATEGORY_LABEL_MAP[category] || category || '其他';

  return {
    id: raw.id,
    type,
    title: raw.title || '',
    desc: raw.desc || '',
    category,
    categoryLabel,
    location: extra.location || '',
    time: extra.event_time || '',
    contact: extra.contact || '',
    status: isLost ? '寻找中' : '待认领',
    statusTone: isLost ? 'amber' : 'green',
    sponsor: raw.publisher || '',
    sponsorInitial: raw.publisher_initial || (raw.publisher || '').charAt(0),
    sponsorTag: '',
    avatarColor: AVATAR_COLORS[stableIndex(raw.id, AVATAR_COLORS.length)],
    note: extra.item_feature || '',
    tags: [categoryLabel, isLost ? '急寻' : '可认领'].filter(Boolean),
    initial: raw.publisher_initial || (raw.publisher || '').charAt(0),
  };
}

export async function loadLostFoundList(params = {}) {
  const response = await fetchLostFoundList(params);
  const data = unwrapPayload(response);
  const items = (data.list || []).map(mapLostFoundItem);

  return {
    items,
    total: Number(data.total || 0),
  };
}
