import { getMenuButtonSafeArea } from '../../../utils/navigation';
import { getNetworkConfirmMessage } from '../../../utils/networkError';
import {
  fetchErrandDetail,
  acceptErrand,
  cancelErrandAccept,
  cancelErrandPublish,
} from '../../../api/modules/errand';
import { ERRAND_CATEGORIES } from '../data';

function formatTimeAgo(dateStr) {
  const now = Date.now();
  const then = new Date(dateStr).getTime();
  const diff = now - then;
  const minutes = Math.floor(diff / 60000);
  if (minutes < 1) return '刚刚';
  if (minutes < 60) return `${minutes}分钟前`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}小时前`;
  const days = Math.floor(hours / 24);
  return `${days}天前`;
}

function formatDeadline(dateStr) {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  const now = new Date();
  const isToday = date.toDateString() === now.toDateString();
  const tomorrow = new Date(now.getFullYear(), now.getMonth(), now.getDate() + 1);
  const isTomorrow = date.toDateString() === tomorrow.toDateString();
  const hours = date.getHours() < 10 ? `0${date.getHours()}` : `${date.getHours()}`;
  const minutes = date.getMinutes() < 10 ? `0${date.getMinutes()}` : `${date.getMinutes()}`;
  if (isToday) return `今天 ${hours}:${minutes} 前`;
  if (isTomorrow) return `明天 ${hours}:${minutes} 前`;
  return `${date.getMonth() + 1}月${date.getDate()}日 ${hours}:${minutes} 前`;
}

const TONE_MAP = {
  parcel: 'blue',
  food: 'orange',
  seat: 'purple',
  print: 'green',
  other: 'gray',
};

function mapErrandDetail(item, userRole, canViewContact, canEdit, canDelete, canAccept, canCancelAccept) {
  const extra = item.extra || {};
  const status = item.status || extra.status || '';
  const isAccepted = !!item.is_accepted || status === 'accepted';
  const isPublisher = userRole === 'publisher';
  const isAcceptor = userRole === 'acceptor';
  let primaryActionText = '立即抢单';
  if (isPublisher) primaryActionText = isAccepted ? '已被接单' : '等待接单';
  else if (isAccepted) primaryActionText = '已锁定订单';
  const images = (item.images || extra.images || []).map((url, index) => ({
    key: `img-${item.id}-${index}`,
    url,
  }));
  const hasContact = Boolean(item.contact || extra.contact);
  return {
    id: item.id,
    category: item.category || extra.category,
    type: (ERRAND_CATEGORIES.find((c) => c.value === (item.category || extra.category)) || {}).label || '其他',
    title: item.title || '',
    desc: item.desc || '',
    detail: item.desc || '',
    routeStart: item.route_start || extra.route_start || '',
    routeEnd: item.route_end || extra.route_end || '',
    deadline: formatDeadline(item.deadline || extra.deadline),
    reward: item.reward || extra.reward || '0',
    rewardText: item.reward || extra.reward || '0',
    urgent: Boolean(item.urgent || extra.urgent),
    images,
    views: 0,
    publisher: item.publisher || '匿名',
    publisherInitial: item.publisher_initial || '匿',
    publisherMeta: formatTimeAgo(item.created_at),
    time: formatTimeAgo(item.created_at),
    credit: '信用良好',
    creditShort: '良好',
    avatarTone: TONE_MAP[item.category || extra.category] || 'gray',
    contact: canViewContact && hasContact ? (item.contact || extra.contact || '') : '',
    status,
    userRole,
    isPublisher,
    isAcceptor,
    isOwner: isPublisher,
    accepted: isAccepted,
    canEdit: Boolean(canEdit),
    canDelete: Boolean(canDelete),
    canAccept: Boolean(canAccept),
    canCancelPublish: Boolean(canDelete),
    canCancelAccept: Boolean(canCancelAccept),
    canCopyContact: canViewContact && hasContact,
    canViewContact,
    primaryActionText,
  };
}

Page({
  data: {
    task: {},
    accepted: false,
    loading: true,
    menuTop: 40,
    menuSafeRight: 16,
    menuButtonHeight: 32,
    navHeight: 88,
  },

  onLoad(options = {}) {
    this.applyNavigationSafeArea();
    this.loadTask(options.id);
  },

  async loadTask(id) {
    if (!id) {
      this.setData({ loading: false });
      return;
    }
    this.setData({ loading: true });
    try {
      const res = await fetchErrandDetail(id);
      const data = res.data || {};
      const item = data.item || data;
      const task = mapErrandDetail(
        item,
        data.user_role || 'viewer',
        data.can_view_contact || false,
        data.can_edit || false,
        data.can_delete || false,
        data.can_accept || false,
        data.can_cancel_accept || false,
      );
      this.setData({ task, accepted: task.accepted, loading: false });
    } catch (e) {
      this.setData({ loading: false });
      wx.showToast({ title: (e && e.message) || '加载失败，请重试', icon: 'none' });
    }
  },

  applyNavigationSafeArea() {
    const { right, top, height } = getMenuButtonSafeArea(10);
    this.setData({
      menuTop: top,
      menuSafeRight: right,
      menuButtonHeight: height,
      navHeight: top + height + 16,
    });
  },

  onBack() {
    if (getCurrentPages().length > 1) {
      wx.navigateBack({ delta: 1 });
      return;
    }
    wx.navigateTo({ url: '/pages/errand/index' });
  },

  onContact() {
    const { contact, canCopyContact, canViewContact } = this.data.task;
    if (!canViewContact) {
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
    if (!canCopyContact) {
      wx.showToast({ title: '发布者暂未留下联系方式', icon: 'none' });
      return;
    }
    if (!contact) {
      wx.showToast({ title: '发布者暂未留下联系方式', icon: 'none' });
      return;
    }

    wx.setClipboardData({
      data: contact,
      success: () => {
        wx.showToast({ title: '联系方式已复制', icon: 'success' });
      },
    });
  },

  onPreviewImages(e) {
    const images = this.data.task.images || [];
    if (!images.length) return;
    const index = Number(e.currentTarget.dataset.index || 0);
    wx.previewImage({
      current: images[index].url,
      urls: images.map((item) => item.url),
    });
  },

  async onAccept() {
    if (!this.data.task.canAccept) return;
    try {
      await acceptErrand(this.data.task.id);
      await this.loadTask(this.data.task.id);
      wx.showToast({ title: '已为你锁定订单', icon: 'success' });
    } catch (err) {
      wx.showToast({ title: getNetworkConfirmMessage(err, (err && err.message) || '接单失败，请重试'), icon: 'none' });
    }
  },

  onCancelPublish() {
    wx.showModal({
      title: '取消发布',
      content: '取消后该跑腿将不再展示，确认取消发布？',
      confirmText: '取消发布',
      confirmColor: '#dc2626',
      success: async (res) => {
        if (!res.confirm) return;
        try {
          await cancelErrandPublish(this.data.task.id);
          wx.showToast({ title: '已取消发布', icon: 'success' });
          setTimeout(() => {
            if (getCurrentPages().length > 1) {
              wx.navigateBack({ delta: 1 });
              return;
            }
            wx.redirectTo({ url: '/pages/errand/index' });
          }, 500);
        } catch (err) {
          wx.showToast({ title: getNetworkConfirmMessage(err, (err && err.message) || '取消失败，请重试'), icon: 'none' });
        }
      },
    });
  },

  onCancelAccept() {
    wx.showModal({
      title: '取消接单',
      content: '取消后该任务会重新开放给其他同学，确认取消接单？',
      confirmText: '取消接单',
      confirmColor: '#dc2626',
      success: async (res) => {
        if (!res.confirm) return;
        try {
          await cancelErrandAccept(this.data.task.id);
          await this.loadTask(this.data.task.id);
          wx.showToast({ title: '已取消接单', icon: 'success' });
        } catch (err) {
          wx.showToast({ title: getNetworkConfirmMessage(err, (err && err.message) || '取消失败，请重试'), icon: 'none' });
        }
      },
    });
  },

  onPublisherHome() {
    wx.showToast({ title: '个人主页接入中', icon: 'none' });
  },

  onShareAppMessage() {
    const { task } = this.data;
    return {
      title: task.title || '校园跑腿',
      path: `/pages/errand/detail/index?id=${task.id || ''}`,
    };
  },
});
