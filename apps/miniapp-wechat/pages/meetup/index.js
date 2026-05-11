import { LIST_CATEGORIES } from './data';
import { loadMeetupList, mergeMeetupItems } from '../../services/meetupService';

Page({
  data: {
    categoryList: LIST_CATEGORIES,
    activeCategory: 'all',
    meetupList: [],
    visibleMeetups: [],
    searchKeyword: '',
    headerHeight: 194,
    page: 1,
    pageSize: 20,
    total: 0,
    loading: false,
  },

  onShow() {
    this.refreshMeetups();
  },

  async refreshMeetups() {
    this.setData({ loading: true, page: 1 });
    const { activeCategory, searchKeyword, pageSize } = this.data;
    const params = { page: 1, pageSize };
    if (activeCategory !== 'all') params.category = activeCategory;
    if (searchKeyword && searchKeyword.trim()) params.keyword = searchKeyword.trim();

    try {
      const result = await loadMeetupList(params);
      this.setData({
        meetupList: result.items,
        visibleMeetups: result.items,
        total: result.total,
        page: 1,
        loading: false,
      });
    } catch (error) {
      this.setData({ loading: false });
      wx.showToast({ title: error.message || '加载失败，请重试', icon: 'none' });
    }
  },

  onBack() {
    if (getCurrentPages().length > 1) {
      wx.navigateBack({ delta: 1 });
      return;
    }
    wx.switchTab({ url: '/pages/home/index' });
  },

  onSearchInput(e) {
    this.setData({ searchKeyword: e.detail.value || '' });
  },

  onSearchConfirm() {
    this.refreshMeetups();
  },

  onSearchClear() {
    this.setData({ searchKeyword: '' }, () => {
      this.refreshMeetups();
    });
  },

  onCategoryChange(e) {
    const { value } = e.detail;
    if (!value || value === this.data.activeCategory) return;
    this.setData({ activeCategory: value }, () => {
      this.refreshMeetups();
    });
  },

  onHeaderHeightChange(e) {
    const { height } = e.detail;
    if (height) this.setData({ headerHeight: height });
  },

  async onLoadMore() {
    const { loading, page, pageSize, total, activeCategory, searchKeyword } = this.data;
    if (loading) return;
    if ((page * pageSize) >= total) return;

    const nextPage = page + 1;
    this.setData({ loading: true });
    const params = { page: nextPage, pageSize };
    if (activeCategory !== 'all') params.category = activeCategory;
    if (searchKeyword && searchKeyword.trim()) params.keyword = searchKeyword.trim();

    try {
      const result = await loadMeetupList(params);
      const merged = mergeMeetupItems(this.data.meetupList, result.items);
      this.setData({
        meetupList: merged,
        visibleMeetups: merged,
        page: nextPage,
        loading: false,
      });
    } catch (error) {
      this.setData({ loading: false });
      wx.showToast({ title: error.message || '加载失败，请重试', icon: 'none' });
    }
  },

  onOpenMeetup(e) {
    const { id } = e.currentTarget.dataset;
    if (!id) return;
    wx.navigateTo({ url: `/pages/meetup/detail/index?id=${id}` });
  },

  onCreateTap() {
    wx.navigateTo({ url: '/pages/meetup/publish/index' });
  },
});
