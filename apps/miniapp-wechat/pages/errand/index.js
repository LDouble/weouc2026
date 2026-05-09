import {
  acceptErrandTask,
  getErrandCategoryList,
  loadErrandList,
} from '../../services/errandService';

Page({
  data: {
    activeCategory: 'all',
    categoryList: getErrandCategoryList(),
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
      const result = await loadErrandList(params);
      this.setData({ taskList: result.items, total: result.total, page: 1, loading: false }, () => {
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
      const result = await loadErrandList(params);
      this.setData(
        {
          taskList: taskList.concat(result.items),
          total: result.total || total,
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
      const taskList = await Promise.all(this.data.taskList.map(async (item) => {
        if (item.id !== task.id) return item;
        return acceptErrandTask(item);
      }));
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
