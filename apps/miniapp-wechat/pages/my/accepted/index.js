import { fetchErrandList } from '~/api/modules/errand';

Page({
  data: {
    headerHeight: 0,
    activeCategory: 'all',
    categoryFilters: [
      { key: 'all', label: '全部' },
      { key: 'parcel', label: '取快递' },
      { key: 'food', label: '代买饭' },
      { key: 'seat', label: '帮占座' },
      { key: 'print', label: '代打印' },
      { key: 'other', label: '其他' },
    ],
    allRecords: [],
    filteredRecords: [],
  },

  onLoad() {
    this.loadRecords();
  },

  onShow() {
    this.loadRecords();
  },

  async loadRecords() {
    try {
      const res = await fetchErrandList({ page: 1, pageSize: 50, user_role: 'acceptor' });
      const data = res.data || res;
      const list = data.list || [];
      const allRecords = list.map((item) => ({
        id: item.id,
        title: item.title || '',
        category: item.category || '',
        desc: item.desc || '',
        reward: item.reward || '0',
        rewardText: Number(item.reward || 0).toFixed(2),
        status: item.status || '',
        deadline: item.deadline || '',
        publisher: item.publisher || '',
        created_at: item.created_at || '',
      }));
      this.setData({ allRecords }, () => {
        this.filterRecords();
      });
    } catch (e) {
      console.error(e);
    }
  },

  filterRecords() {
    const { activeCategory, allRecords } = this.data;
    const filteredRecords =
      activeCategory === 'all'
        ? allRecords
        : allRecords.filter((item) => item.category === activeCategory);
    this.setData({ filteredRecords });
  },

  onCategoryFilter(e) {
    const { key } = e.currentTarget.dataset;
    if (key === this.data.activeCategory) return;
    this.setData({ activeCategory: key }, () => {
      this.filterRecords();
    });
  },

  onRecordClick(e) {
    const { id } = e.currentTarget.dataset;
    wx.navigateTo({ url: `/pages/errand/detail/index?id=${id}` });
  },

  onBack() {
    wx.navigateBack();
  },

  onHeaderHeightChange(e) {
    this.setData({ headerHeight: e.detail.height });
  },
});
