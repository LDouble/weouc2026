import Message from 'tdesign-miniprogram/message/index';
import { getMenuButtonSafeArea } from '../../utils/navigation';
import { fetchFeedList } from '../../api/modules/feed';

function formatRelativeTime(dateStr) {
  if (!dateStr) return '';
  const now = Date.now();
  const time = new Date(dateStr).getTime();
  const diff = now - time;
  if (diff < 60000) return '刚刚';
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`;
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`;
  if (diff < 604800000) return `${Math.floor(diff / 86400000)} 天前`;
  return dateStr.slice(0, 10);
}

function getFeedTargetUrl(feedType, id) {
  if (!id) return '';
  const detailRoutes = {
    market: '/pages/market/detail/index',
    errand: '/pages/errand/detail/index',
    lostFound: '/pages/lost-found/detail/index',
  };
  if (detailRoutes[feedType]) return `${detailRoutes[feedType]}?id=${id}`;

  const listRoutes = {
    resource: '/pages/resource/index',
    carpool: '/pages/carpool/index',
  };
  if (listRoutes[feedType]) return `${listRoutes[feedType]}?insertId=${id}`;

  return '';
}

function mapFeedItem(item) {
  const extra = item.extra || {};
  const image = item.image || (extra.images && extra.images[0]) || '';
  return {
    id: item.id || '',
    feedType: item.feed_type || '',
    name: item.publisher || '',
    time: formatRelativeTime(item.created_at),
    badge: item.feed_type_label || '',
    badgeTone: 'hot',
    tone: 'indigo',
    avatarIcon: 'user',
    title: item.title || '',
    desc: item.desc || '',
    image,
    targetUrl: getFeedTargetUrl(item.feed_type, item.id),
    likes: 0,
    comments: 0,
  };
}

function distributeToColumns(items) {
  const left = [];
  const right = [];
  items.forEach((item, index) => {
    if (index % 2 === 0) left.push(item);
    else right.push(item);
  });
  return { left, right };
}

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
    loadingFeed: false,
    heroCurrent: 1,
    weekText: '第 8 周',
    heroCards: [
      { id: 'cet', tag: '重要通知', title: '四六级考试报名进行中', desc: '请于本周五前完成报名缴费', theme: 'rose' },
      { id: 'job', tag: '职业发展', title: '春季校园招聘会', desc: '300+ 企业现场招聘', theme: 'indigo' },
      { id: 'seat', tag: '学习资源', title: '图书馆新增自习座位', desc: '支持线上预约，先到先得', theme: 'mint' },
    ],
    quickLinks: [
      { label: '资料', icon: 'book', tone: 'amber', url: '/pages/resource/index' },
      { label: '闲置', icon: 'shop', tone: 'cyan', url: '/pages/market/index' },
      { label: '拼车', icon: 'car', tone: 'teal', url: '/pages/carpool/index' },
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
    this.loadFeed();
  },

  applyNavigationSafeArea() {
    const { right } = getMenuButtonSafeArea();
    this.setData({ menuSafeRight: right });
  },

  loadFeed() {
    this.setData({ loadingFeed: true });
    fetchFeedList({ page: 1, pageSize: this.data.feedPageSize })
      .then((res) => {
        const data = res.data || res;
        const list = (data.list || []).map(mapFeedItem);
        const columns = distributeToColumns(list);
        this.setData({
          campusFeedLeft: columns.left,
          campusFeedRight: columns.right,
          feedPage: 1,
          feedTotal: data.total || 0,
          hasMore: (data.total || 0) > this.data.feedPageSize,
          loadingFeed: false,
        });
      })
      .catch(() => {
        this.setData({ loadingFeed: false });
      });
  },

  onRefresh() {
    this.setData({ enable: true });
    fetchFeedList({ page: 1, pageSize: this.data.feedPageSize })
      .then((res) => {
        const data = res.data || res;
        const list = (data.list || []).map(mapFeedItem);
        const columns = distributeToColumns(list);
        this.setData({
          enable: false,
          campusFeedLeft: columns.left,
          campusFeedRight: columns.right,
          feedPage: 1,
          feedTotal: data.total || 0,
          hasMore: (data.total || 0) > this.data.feedPageSize,
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

    fetchFeedList({ page: nextPage, pageSize: this.data.feedPageSize })
      .then((res) => {
        const data = res.data || res;
        const list = (data.list || []).map(mapFeedItem);
        const columns = distributeToColumns(list);
        this.setData({
          isLoadingMore: false,
          feedPage: nextPage,
          feedTotal: data.total || 0,
          hasMore: nextPage * this.data.feedPageSize < (data.total || 0),
          campusFeedLeft: this.data.campusFeedLeft.concat(columns.left),
          campusFeedRight: this.data.campusFeedRight.concat(columns.right),
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

  goNotify() { this.showInfo('通知中心正在接入中'); },
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
