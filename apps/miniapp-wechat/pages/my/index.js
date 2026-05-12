import useToastBehavior from '~/behaviors/useToast';
import { loadNotificationUnreadCount } from '~/services/notificationService';
import { loadMyPageModel } from '~/services/profileService';
import { getSessionState, subscribeSession } from '~/stores/session';

function createStats(values = {}) {
  return [
    { key: 'published', icon: 'chart-bar', value: values.published || 0, label: '发布', bg: 'indigo' },
    { key: 'accepted', icon: 'check-circle', value: values.accepted || 0, label: '接单', bg: 'purple' },
    { key: 'fav', icon: 'star', value: values.favorite || 0, label: '收藏', bg: 'amber' },
  ];
}

function createSettingList(isBound = false, unreadCount = 0) {
  const badge = Math.max(0, Number(unreadCount || 0));
  return [
    { name: isBound ? '教务已绑定' : '教务绑定', icon: 'education', type: 'eduBind', url: '/pages/edu-bind/index' },
    { name: '消息中心', icon: 'notification', type: 'message', url: '/pages/message/index', badge },
    { name: '联系客服', icon: 'service', type: 'service', url: '' },
    { name: '设置', icon: 'setting', type: 'setting', url: '/pages/setting/index' },
  ];
}

function buildSessionFallbackInfo(sessionState) {
  const viewer = sessionState.viewer || {};

  return {
    name: viewer.nickname || '微信用户',
    image: viewer.avatarUrl || '',
    major: '',
    sid: '',
    college: '',
    grade: '',
    isBound: false,
    identityLabel: '未绑定教务',
    secondaryLabel: '微信登录用户',
  };
}

Page({
  behaviors: [useToastBehavior],

  data: {
    hasSession: false,
    isLoad: false,
    headerHeight: 0,
    personalInfo: {},
    stats: createStats(),
    serviceList: [
      { name: '我的发布', icon: 'root-list', type: 'publish', url: '/pages/my/publish/index' },
      { name: '我的接单', icon: 'check-circle', type: 'accepted', url: '/pages/my/accepted/index' },
      { name: '星标收藏', icon: 'star', type: 'fav', url: '' },
    ],
    settingList: createSettingList(false),
  },

  onLoad() {
    this.unsubscribeSession = subscribeSession((state) => {
      if (!state.authenticated) {
        this.applyGuestState();
      }
    });
  },

  onUnload() {
    if (this.unsubscribeSession) {
      this.unsubscribeSession();
      this.unsubscribeSession = null;
    }
  },

  onBack() {
    wx.navigateBack();
  },

  onHeaderHeightChange(e) {
    this.setData({
      headerHeight: e.detail.height,
    });
  },

  applyGuestState() {
    this.setData({
      hasSession: false,
      isLoad: false,
      personalInfo: {},
      stats: createStats(),
      settingList: createSettingList(false, 0),
    });
  },

  applyAuthenticatedState(model) {
    this.setData({
      hasSession: true,
      isLoad: true,
      personalInfo: model.personalInfo,
      stats: createStats(model.stats),
      settingList: createSettingList(Boolean(model.personalInfo.isBound), model.unreadCount),
    });
  },

  async onShow() {
    const sessionState = getSessionState();
    if (!sessionState.authenticated) {
      this.applyGuestState();
      return;
    }

    const unreadPromise = loadNotificationUnreadCount().catch(() => 0);

    this.applyAuthenticatedState({
      personalInfo: buildSessionFallbackInfo(sessionState),
      stats: {
        published: 0,
        accepted: 0,
        favorite: 0,
      },
      unreadCount: 0,
    });

    try {
      const pageModel = await loadMyPageModel();
      const unreadCount = await unreadPromise;
      this.applyAuthenticatedState({ ...pageModel, unreadCount });
    } catch (e) {
      const unreadCount = await unreadPromise;
      this.setData({
        settingList: createSettingList(Boolean(this.data.personalInfo.isBound), unreadCount),
      });
      this.onShowToast('#t-toast', '个人页加载失败');
    }
  },

  onLogin() {
    wx.navigateTo({
      url: '/pages/login/login',
    });
  },

  onNavigateTo() {
    wx.navigateTo({ url: `/pages/my/info-edit/index` });
  },

  onEleClick(e) {
    const { name, url } = e.currentTarget.dataset.data;
    if (url) {
      wx.navigateTo({ url });
      return;
    }
    this.onShowToast('#t-toast', name);
  },
});
