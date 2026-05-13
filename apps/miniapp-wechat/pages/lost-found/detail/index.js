import { fetchLostFoundDetail, deleteLostFound, resolveLostFound } from '../../../api/modules/lostFound';
import { LOST_FOUND_CATEGORIES } from '../data';
import { getStatusDisplay, normalizeContactFields } from '../../../services/shared';

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

function getLostFoundDetailOverrides(type) {
  if (type === 'lost') {
    return {
      published: { label: '寻找中' },
      resolved: { label: '已找到' },
    };
  }
  return {
    published: { label: '待认领' },
    resolved: { label: '已认领' },
  };
}

function mapLostFoundDetail(raw = {}) {
  const extra = raw.extra || {};
  const type = extra.type || raw.feed_type || 'lost';
  const category = extra.category || '';
  const categoryLabel = CATEGORY_LABEL_MAP[category] || category || '其他';
  const contactFields = normalizeContactFields(raw);
  const status = raw.status || 'published';
  const statusDisplay = getStatusDisplay(status, getLostFoundDetailOverrides(type));

  return {
    id: raw.id || '',
    type,
    typeLabel: type === 'lost' ? '失物登记' : '招领登记',
    reviewStatus: status,
    statusLabel: statusDisplay.label,
    statusTone: statusDisplay.tone,
    status,
    title: raw.title || '',
    desc: raw.desc || '',
    category,
    categoryLabel,
    location: extra.location || '',
    eventTime: extra.event_time || '',
    contact: contactFields.contact,
    canViewContact: contactFields.canViewContact,
    contactHiddenReason: contactFields.contactHiddenReason,
    note: extra.item_feature || '',
    publisher: raw.publisher || '匿名同学',
    publisherInitial: raw.publisher_initial || (raw.publisher || '匿').charAt(0),
    publisherMeta: formatRelativeTime(raw.created_at),
    avatarColor: AVATAR_COLORS[stableIndex(raw.id, AVATAR_COLORS.length)],
    isOwner: Boolean(raw.is_owner),
    canEdit: Boolean(raw.can_edit),
    canDelete: Boolean(raw.can_delete),
    canMarkResolved: Boolean(raw.can_mark_resolved),
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
      wx.showToast({ title: error.message || '加载失败，请重试', icon: 'none' });
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
    const { contact, canViewContact } = this.data.detail;
    if (canViewContact === false) {
      wx.showModal({
        title: '无法查看联系方式',
        content: '绑定教务后即可查看联系方式，是否前往绑定？',
        confirmText: '前往绑定',
        cancelText: '暂不需要',
        success: (res) => {
          if (res.confirm) {
            wx.navigateTo({ url: '/pages/edu-bind/index' });
          }
        },
      });
      return;
    }
    if (!contact) {
      wx.showToast({ title: '登记人未公开联系方式', icon: 'none' });
      return;
    }

    wx.setClipboardData({
      data: contact,
      success: () => {
        wx.showToast({ title: '联系方式已复制', icon: 'success' });
      },
    });
  },

  onDelete() {
    const { detail } = this.data;
    if (!detail.canDelete) return;
    wx.showModal({
      title: '下架信息',
      content: '下架后该信息将不再展示，确认下架？',
      confirmText: '确认下架',
      confirmColor: '#dc2626',
      success: async (res) => {
        if (!res.confirm) return;
        try {
          await deleteLostFound(detail.id);
          wx.showToast({ title: '已下架', icon: 'success' });
          setTimeout(() => {
            if (getCurrentPages().length > 1) {
              wx.navigateBack({ delta: 1 });
              return;
            }
            wx.redirectTo({ url: '/pages/lost-found/index' });
          }, 500);
        } catch (error) {
          wx.showToast({ title: (error && error.message) || '下架失败，请重试', icon: 'none' });
        }
      },
    });
  },

  onMarkResolved() {
    const { detail } = this.data;
    if (!detail.canMarkResolved) return;
    const label = detail.type === 'lost' ? '已找到' : '已认领';
    wx.showModal({
      title: `标记${label}`,
      content: `确认标记为${label}？标记后信息将不再展示。`,
      confirmText: `确认${label}`,
      success: async (res) => {
        if (!res.confirm) return;
        try {
          await resolveLostFound(detail.id);
          wx.showToast({ title: `已标记${label}`, icon: 'success' });
          await this.loadDetail(detail.id);
        } catch (error) {
          wx.showToast({ title: (error && error.message) || '操作失败，请重试', icon: 'none' });
        }
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
