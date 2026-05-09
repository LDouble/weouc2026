import { LOST_FOUND_CATEGORIES } from './data';
import { fetchLostFoundList } from '../../api/modules/lostFound';

const CATEGORY_LABEL_MAP = LOST_FOUND_CATEGORIES.reduce((map, item) => {
  if (item.value && item.value !== 'all') {
    map[item.value] = item.label;
  }
  return map;
}, {});

function mapLostFoundItem(raw) {
  const extra = raw.extra || {};
  const type = extra.type || raw.feed_type || 'lost';
  const isLost = type === 'lost';
  const category = extra.category || '';
  const categoryLabel = CATEGORY_LABEL_MAP[category] || category || '其他';
  return {
    id: raw.id,
    type,
    title: raw.title || '',
    desc: raw.desc || '',
    category,
    categoryLabel,
    location: extra.location || '',
    time: extra.event_time || '',
    contact: extra.contact || '',
    status: isLost ? '寻找中' : '待认领',
    statusTone: isLost ? 'amber' : 'green',
    sponsor: raw.publisher || '',
    sponsorInitial: raw.publisher_initial || (raw.publisher || '').charAt(0),
    sponsorTag: '',
    avatarColor: ['amber', 'rose', 'blue', 'purple', 'green'][Math.abs(raw.id || 0) % 5],
    note: extra.item_feature || '',
    tags: [categoryLabel, isLost ? '急寻' : '可认领'].filter(Boolean),
    initial: raw.publisher_initial || (raw.publisher || '').charAt(0),
  };
}

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
      const res = await fetchLostFoundList(params);
      const list = (res.data && res.data.list) || [];
      const allItems = list.map(mapLostFoundItem);
      this.setData({
        allItems,
        visibleItems: allItems,
        page: 1,
        total: (res.data && res.data.total) || 0,
        loading: false,
      });
    } catch (error) {
      this.setData({ loading: false });
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
      const res = await fetchLostFoundList(params);
      const list = (res.data && res.data.list) || [];
      const moreItems = list.map(mapLostFoundItem);
      const merged = allItems.concat(moreItems);
      this.setData({
        allItems: merged,
        visibleItems: merged,
        page: nextPage,
        total: (res.data && res.data.total) || 0,
        loading: false,
      });
    } catch (error) {
      this.setData({ loading: false });
    }
  },

  filterItems() {
    const { allItems, searchKeyword } = this.data;
    let filtered = allItems;
    if (searchKeyword) {
      const kw = searchKeyword.toLowerCase();
      filtered = filtered.filter(
        (item) =>
          (item.title || '').toLowerCase().includes(kw) ||
          (item.desc || '').toLowerCase().includes(kw) ||
          (item.location || '').toLowerCase().includes(kw),
      );
    }
    this.setData({ visibleItems: filtered });
  },

  onSearchInput(e) {
    this.setData({ searchKeyword: e.detail.value || '' }, () => this.loadItems());
  },

  onSearchClear() {
    this.setData({ searchKeyword: '' }, () => this.loadItems());
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
    if (!contact || contact === '站内私信') {
      wx.showToast({ title: '可通过站内私信联系', icon: 'none' });
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
