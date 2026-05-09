import { ERRAND_CATEGORIES } from './data';
import { fetchErrandList, acceptErrand } from '../../api/modules/errand';

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

function mapErrandItem(item) {
  const status = item.status || '';
  const userRole = item.user_role || 'viewer';
  const isAccepted = !!item.is_accepted || status === 'accepted';
  const isPublisher = userRole === 'publisher';
  const isAcceptor = userRole === 'acceptor';
  let actionText = '接单';
  if (isPublisher) actionText = '我的发布';
  else if (isAcceptor) actionText = '已接单';
  else if (isAccepted) actionText = '已被接';
  return {
    id: item.id,
    category: item.category,
    type: (ERRAND_CATEGORIES.find((c) => c.value === item.category) || {}).label || '其他',
    title: item.title || '',
    desc: item.desc || '',
    routeStart: item.route_start || '',
    routeEnd: item.route_end || '',
    deadline: formatDeadline(item.deadline),
    reward: item.reward || '0',
    urgent: false,
    views: 0,
    publisher: item.publisher || '匿名',
    publisherInitial: item.publisher_initial || '匿',
    publisherMeta: formatTimeAgo(item.created_at),
    time: formatTimeAgo(item.created_at),
    credit: '信用良好',
    creditShort: '良好',
    avatarTone: TONE_MAP[item.category] || 'gray',
    status,
    userRole,
    accepted: isAccepted,
    canAccept: !isPublisher && !isAccepted && status !== 'cancelled',
    actionText,
  };
}

Page({
  data: {
    activeCategory: 'all',
    categoryList: ERRAND_CATEGORIES,
    taskList: [],
    visibleTasks: [],
    headerHeight: 138,
    searchKeyword: '',
    page: 1,
    pageSize: 20,
    total: 0,
    loading: false,
  },

  onShow() {
    this.refreshTasks();
  },

  async refreshTasks() {
    this.setData({ loading: true, page: 1 });
    try {
      const { activeCategory, searchKeyword, pageSize } = this.data;
      const params = { page: 1, pageSize };
      if (activeCategory !== 'all') params.category = activeCategory;
      if (searchKeyword) params.keyword = searchKeyword;
      const res = await fetchErrandList(params);
      const list = (res.data && res.data.list) || [];
      const total = (res.data && res.data.total) || 0;
      const taskList = list.map(mapErrandItem);
      this.setData({ taskList, total, page: 1, loading: false }, () => {
        this.filterTasks();
      });
    } catch (e) {
      this.setData({ loading: false });
      wx.showToast({ title: '加载失败，请重试', icon: 'none' });
    }
  },

  async onLoadMore() {
    const { loading, page, pageSize, total, activeCategory, searchKeyword, taskList } = this.data;
    if (loading) return;
    if (taskList.length >= total) return;
    const nextPage = page + 1;
    this.setData({ loading: true });
    try {
      const params = { page: nextPage, pageSize };
      if (activeCategory !== 'all') params.category = activeCategory;
      if (searchKeyword) params.keyword = searchKeyword;
      const res = await fetchErrandList(params);
      const list = (res.data && res.data.list) || [];
      const newItems = list.map(mapErrandItem);
      this.setData(
        {
          taskList: taskList.concat(newItems),
          total: (res.data && res.data.total) || total,
          page: nextPage,
          loading: false,
        },
        () => {
          this.filterTasks();
        },
      );
    } catch (e) {
      this.setData({ loading: false });
      wx.showToast({ title: '加载失败，请重试', icon: 'none' });
    }
  },

  filterTasks() {
    const { taskList } = this.data;
    this.setData({ visibleTasks: taskList });
  },

  onBack() {
    if (getCurrentPages().length > 1) {
      wx.navigateBack({ delta: 1 });
      return;
    }
    wx.switchTab({ url: '/pages/home/index' });
  },

  onCategoryChange(e) {
    const { value } = e.detail;
    if (!value || value === this.data.activeCategory) return;
    this.setData({ activeCategory: value }, () => {
      this.refreshTasks();
    });
  },

  onSearchInput(e) {
    this.setData({ searchKeyword: e.detail.value });
  },

  onSearchConfirm() {
    this.refreshTasks();
  },

  onSearchClear() {
    this.setData({ searchKeyword: '' }, () => {
      this.refreshTasks();
    });
  },

  onHeaderHeightChange(e) {
    const { height } = e.detail;
    if (height) this.setData({ headerHeight: height });
  },

  onTaskSelect(e) {
    const { task } = e.detail;
    wx.navigateTo({ url: `/pages/errand/detail/index?id=${task.id}` });
  },

  async onAcceptTask(e) {
    const { task } = e.detail;
    if (!task.canAccept) {
      this.onTaskSelect(e);
      return;
    }
    try {
      await acceptErrand(task.id);
      const taskList = this.data.taskList.map((item) => {
        if (item.id !== task.id) return item;
        return { ...item, accepted: true, userRole: 'acceptor', canAccept: false, actionText: '已接单' };
      });
      this.setData({ taskList }, () => {
        this.filterTasks();
        wx.showToast({ title: '已为你锁定订单', icon: 'success' });
        setTimeout(() => {
          wx.navigateTo({ url: `/pages/errand/detail/index?id=${task.id}` });
        }, 450);
      });
    } catch (err) {
      wx.showToast({ title: '接单失败，请重试', icon: 'none' });
    }
  },

  goRelease() {
    wx.navigateTo({ url: '/pages/errand/publish/index' });
  },
});
