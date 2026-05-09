import { fetchLostFoundDetail } from '../../../api/modules/lostFound';
import { LOST_FOUND_CATEGORIES } from '../data';

const CATEGORY_LABEL_MAP = LOST_FOUND_CATEGORIES.reduce((map, item) => {
  if (item.value && item.value !== 'all') {
    map[item.value] = item.label;
  }
  return map;
}, {});

const AVATAR_COLORS = ['amber', 'rose', 'blue', 'purple', 'green'];

function formatRelativeTime(dateStr) {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  const time = date.getTime();
  if (Number.isNaN(time)) return dateStr;

  const diff = Date.now() - time;
  if (diff < 60000) return '刚刚登记';
  if (diff < 3600000) return `${Math.max(1, Math.floor(diff / 60000))} 分钟前登记`;
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前登记`;
  if (diff < 604800000) return `${Math.floor(diff / 86400000)} 天前登记`;

  const month = date.getMonth() + 1;
  const day = date.getDate();
  const hours = `${date.getHours()}`.padStart(2, '0');
  const minutes = `${date.getMinutes()}`.padStart(2, '0');
  return `${month}月${day}日 ${hours}:${minutes} 登记`;
}

function mapLostFoundDetail(raw = {}) {
  const extra = raw.extra || {};
  const type = extra.type || raw.feed_type || 'lost';
  const category = extra.category || '';
  const categoryLabel = CATEGORY_LABEL_MAP[category] || category || '其他';

  return {
    id: raw.id || '',
    type,
    typeLabel: type === 'lost' ? '失物登记' : '招领登记',
    status: type === 'lost' ? '寻找中' : '待认领',
    title: raw.title || '',
    desc: raw.desc || '',
    category,
    categoryLabel,
    location: extra.location || '',
    eventTime: extra.event_time || '',
    contact: extra.contact || '',
    note: extra.item_feature || '',
    publisher: raw.publisher || '匿名同学',
    publisherInitial: raw.publisher_initial || (raw.publisher || '匿').charAt(0),
    publisherMeta: formatRelativeTime(raw.created_at),
    avatarColor: AVATAR_COLORS[Math.abs(raw.id || 0) % AVATAR_COLORS.length],
  };
}

Page({
  data: {
    headerHeight: 0,
    loading: true,
    detail: {},
  },

  onLoad(options = {}) {
    this.loadDetail(options.id);
  },

  async loadDetail(id) {
    if (!id) {
      this.setData({ loading: false });
      return;
    }

    this.setData({ loading: true });
    try {
      const res = await fetchLostFoundDetail(id);
      const detail = mapLostFoundDetail(res.data || res);
      this.setData({ detail, loading: false });
    } catch (error) {
      this.setData({ loading: false });
      wx.showToast({ title: '加载失败，请重试', icon: 'none' });
    }
  },

  onHeaderHeightChange(e) {
    const { height } = e.detail;
    if (height) this.setData({ headerHeight: height });
  },

  onBack() {
    if (getCurrentPages().length > 1) {
      wx.navigateBack({ delta: 1 });
      return;
    }
    wx.redirectTo({ url: '/pages/lost-found/index' });
  },

  onCopyContact() {
    const { contact } = this.data.detail;
    if (!contact || contact === '站内私信') {
      wx.showToast({ title: '可通过站内私信联系', icon: 'none' });
      return;
    }

    wx.setClipboardData({
      data: contact,
      success: () => {
        wx.showToast({ title: '联系方式已复制', icon: 'success' });
      },
    });
  },

  onShareAppMessage() {
    const { detail } = this.data;
    return {
      title: detail.title || '失物招领',
      path: `/pages/lost-found/detail/index?id=${detail.id || ''}`,
    };
  },
});
