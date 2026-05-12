import Message from 'tdesign-miniprogram/message/index';
import { loadHomeFeeds } from '../../services/homeService';
import { loadNotificationUnreadCount } from '../../services/notificationService';
import { DEFAULT_HERO_CARDS, loadPortalHeroCards } from '../../services/portalService';
import { getSessionState } from '../../stores/session';
import { getMenuButtonSafeArea } from '../../utils/navigation';

Page({
  data: {
    enable: false,
    menuSafeRight: 16,
    refreshTexts: ['下拉刷新动态', '松手更新首页', '正在同步校园动态', '刷新完成'],
    refreshLoadingProps: {
      size: '36rpx',
      text: '正在同步校园动态',
      theme: 'circular',
    },
    isLoadingMore: false,
    hasMore: true,
    feedPage: 1,
    feedPageSize: 10,
    feedTotal: 0,
    notifyUnreadCount: 0,
    loadingFeed: false,
    heroCurrent: 1,
    weekText: '第 8 周',
    heroCards: DEFAULT_HERO_CARDS,
    quickLinks: [
      { label: '资料', icon: 'book', tone: 'amber', url: '/pages/resource/index' },
      { label: '闲置', icon: 'shop', tone: 'cyan', url: '/pages/market/index' },
      { label: '拼车', icon: 'car', tone: 'teal', url: '/pages/carpool/index' },
      { label: '组局', icon: 'chat', tone: 'green', url: '/pages/meetup/index' },
      { label: '失物', icon: 'search', tone: 'rose', url: '/pages/lost-found/index' },
      { label: '跑腿', icon: 'send', tone: 'orange', url: '/pages/errand/index' },
      { label: '我的', icon: 'user', tone: 'pink', tab: '/pages/my/index' },
      { label: '日程', icon: 'calendar', tone: 'green', disabled: true },
      { label: '成绩', icon: 'chart', tone: 'blue', disabled: true },
      { label: '考试', icon: 'file-copy', tone: 'indigo', disabled: true },
      { label: '自习', icon: 'time', tone: 'purple', disabled: true },
      { label: '校历', icon: 'calendar-1', tone: 'teal', disabled: true },
    ],
    courseList: [
      { id: 'network', period: '1-2', start: '08:00', title: '计算机网络基础', place: '教二-301 · 张教授', status: '上课中', active: true },
      { id: 'math', period: '3-4', start: '10:00', title: '高等数学 A', place: '理学楼-205' },
      { id: 'english', period: '7-8', start: '15:30', title: '大学英语听说', place: '外语楼-412' },
    ],
    campusFeedLeft: [],
    campusFeedRight: [],
    canIUseGetUserProfile: false,
  },

  onLoad(option) {
    this.applyNavigationSafeArea();
    if (wx.getUserProfile) this.setData({ canIUseGetUserProfile: true });
    if (option && option.oper) {
      const content = option.oper === 'release' ? '发布成功' : '保存成功';
      this.showOperMsg(content);
    }
    this.loadHeroCards();
    this.refreshNotificationBadge();
    this.loadFeed();
  },

  onShow() {
    this.refreshNotificationBadge();
  },

  applyNavigationSafeArea() {
    const { right } = getMenuButtonSafeArea();
    this.setData({ menuSafeRight: right });
  },

  loadFeed() {
    this.setData({ loadingFeed: true });
    loadHomeFeeds({ page: 1, pageSize: this.data.feedPageSize })
      .then((result) => {
        this.setData({
          campusFeedLeft: result.columns.left,
          campusFeedRight: result.columns.right,
          feedPage: 1,
          feedTotal: result.total,
          hasMore: result.total > this.data.feedPageSize,
          loadingFeed: false,
        });
      })
      .catch(() => {
        this.setData({ loadingFeed: false });
      });
  },

  loadHeroCards() {
    loadPortalHeroCards()
      .then((heroCards) => {
        if (Array.isArray(heroCards) && heroCards.length) {
          this.setData({ heroCards });
        }
      })
      .catch(() => {});
  },

  refreshNotificationBadge() {
    if (!getSessionState().authenticated) {
      this.setData({ notifyUnreadCount: 0 });
      return;
    }

    loadNotificationUnreadCount()
      .then((count) => {
        this.setData({ notifyUnreadCount: Math.max(0, Number(count || 0)) });
      })
      .catch(() => {
        this.setData({ notifyUnreadCount: 0 });
      });
  },

  onRefresh() {
    this.setData({ enable: true });
    loadHomeFeeds({ page: 1, pageSize: this.data.feedPageSize })
      .then((result) => {
        this.setData({
          enable: false,
          campusFeedLeft: result.columns.left,
          campusFeedRight: result.columns.right,
          feedPage: 1,
          feedTotal: result.total,
          hasMore: result.total > this.data.feedPageSize,
          isLoadingMore: false,
        });
      })
      .catch(() => {
        this.setData({ enable: false });
      });
  },

  onLoadMore() {
    if (this.data.isLoadingMore || !this.data.hasMore) return;
    const nextPage = this.data.feedPage + 1;
    this.setData({ isLoadingMore: true });

    loadHomeFeeds({ page: nextPage, pageSize: this.data.feedPageSize })
      .then((result) => {
        this.setData({
          isLoadingMore: false,
          feedPage: nextPage,
          feedTotal: result.total,
          hasMore: nextPage * this.data.feedPageSize < result.total,
          campusFeedLeft: this.data.campusFeedLeft.concat(result.columns.left),
          campusFeedRight: this.data.campusFeedRight.concat(result.columns.right),
        });
      })
      .catch(() => {
        this.setData({ isLoadingMore: false });
      });
  },

  onHeroChange(e) {
    this.setData({ heroCurrent: e.detail.current });
  },

  goQuickLink(e) {
    const { url, tab, label, disabled } = e.currentTarget.dataset;
    if (disabled) { this.showInfo(`${label}即将上线`); return; }
    if (tab) { wx.switchTab({ url: tab }); return; }
    if (url) { wx.navigateTo({ url }); }
  },

  goNotify() {
    if (!getSessionState().authenticated) {
      wx.navigateTo({ url: '/pages/login/login' });
      return;
    }
    wx.navigateTo({ url: '/pages/message/index' });
  },
  goExploreMore() { this.showInfo('更多服务正在接入中'); },
  goFeedMore() { this.showInfo('更多校园动态正在接入中'); },
  goPost(e) {
    const { item } = e.currentTarget.dataset;
    if (item && item.targetUrl) {
      wx.navigateTo({ url: item.targetUrl });
      return;
    }
    this.showInfo('暂时无法打开该动态');
  },

  showInfo(content) {
    Message.info({ context: this, offset: [120, 32], duration: 2200, content });
  },

  showOperMsg(content) {
    Message.success({ context: this, offset: [120, 32], duration: 3000, content });
  },
});
