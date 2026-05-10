import { LOST_FOUND_CATEGORIES } from './data';
import { loadLostFoundList } from '../../services/lostFoundService';

Page({
  data: {
    headerHeight: 280,
    searchKeyword: '',
    activeType: 'lost',
    typeTabs: [
      { label: '失物', value: 'lost' },
      { label: '招领', value: 'found' },
    ],
    activeCategory: 'all',
    categoryList: LOST_FOUND_CATEGORIES,
    allItems: [],
    visibleItems: [],
    page: 1,
    pageSize: 20,
    total: 0,
    loading: false,
  },

  onLoad() {
    this.loadItems();
  },

  onShow() {
    this.loadItems();
  },

  async loadItems() {
    const { activeType, activeCategory, searchKeyword, pageSize } = this.data;
    this.setData({ loading: true });
    try {
      const params = { type: activeType, page: 1, pageSize };
      if (activeCategory !== 'all') params.category = activeCategory;
      if (searchKeyword) params.keyword = searchKeyword;
      const result = await loadLostFoundList(params);
      this.setData({
        allItems: result.items,
        visibleItems: result.items,
        page: 1,
        total: result.total,
        loading: false,
      });
    } catch (error) {
      this.setData({ loading: false });
      wx.showToast({ title: error.message || '加载失败，请重试', icon: 'none' });
    }
  },

  async onLoadMore() {
    const { activeType, activeCategory, searchKeyword, page, pageSize, total, loading, allItems } = this.data;
    if (loading) return;
    if (allItems.length >= total) return;
    const nextPage = page + 1;
    this.setData({ loading: true });
    try {
      const params = { type: activeType, page: nextPage, pageSize };
      if (activeCategory !== 'all') params.category = activeCategory;
      if (searchKeyword) params.keyword = searchKeyword;
      const result = await loadLostFoundList(params);
      const merged = allItems.concat(result.items);
      this.setData({
        allItems: merged,
        visibleItems: merged,
        page: nextPage,
        total: result.total,
        loading: false,
      });
    } catch (error) {
      this.setData({ loading: false });
      wx.showToast({ title: error.message || '加载失败，请重试', icon: 'none' });
    }
  },

  onSearchInput(e) {
    this.setData({ searchKeyword: e.detail.value || '' });
  },

  onSearchClear() {
    this.setData({ searchKeyword: '' }, () => this.loadItems());
  },

  onSearchConfirm() {
    this.loadItems();
  },

  onTypeChange(e) {
    const { value: type } = e.detail;
    if (!type) return;
    this.setData({ activeType: type }, () => this.loadItems());
  },

  onCategoryChange(e) {
    const { value } = e.detail;
    if (!value) return;
    this.setData({ activeCategory: value }, () => this.loadItems());
  },

  onHeaderHeightChange(e) {
    const { height } = e.detail;
    if (height) this.setData({ headerHeight: height });
  },

  onItemSelect(e) {
    const { item } = e.detail;
    wx.navigateTo({ url: `/pages/lost-found/detail/index?id=${item.id}` });
  },

  onCopyContact(e) {
    const { item } = e.detail;
    const { contact } = item;
    if (!contact) {
      wx.showToast({ title: '登记人未公开联系方式', icon: 'none' });
      return;
    }
    wx.setClipboardData({
      data: contact,
      success: () => wx.showToast({ title: '联系方式已复制', icon: 'success' }),
    });
  },

  onCreateTap() {
    wx.navigateTo({ url: '/pages/lost-found/publish/index' });
  },

  onBack() {
    if (getCurrentPages().length > 1) {
      wx.navigateBack({ delta: 1 });
      return;
    }
    wx.switchTab({ url: '/pages/home/index' });
  },
});
