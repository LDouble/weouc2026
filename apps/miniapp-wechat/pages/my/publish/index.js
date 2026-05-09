import { fetchFeedList } from '~/api/modules/feed';

const TYPE_MAP = {
  market: '跳蚤市场',
  errand: '跑腿',
  lostFound: '失物招领',
  resource: '资源共享',
  carpool: '拼车',
};

const URL_MAP = {
  market: '/pages/market/detail/index',
  errand: '/pages/errand/detail/index',
  lostFound: '/pages/lost-found/detail/index',
  resource: '',
  carpool: '',
};

const STATUS_MAP = {
  pending: '审核中',
  reviewing: '审核中',
  published: '已发布',
  rejected: '已下架',
  offline: '已下架',
};

function normalizeStatus(status) {
  if (status === 'pending') return 'reviewing';
  if (status === 'rejected') return 'offline';
  return status || 'published';
}

Page({
  data: {
    headerHeight: 0,
    activeType: 'all',
    activeStatus: 'all',
    typeFilters: [
      { key: 'all', label: '全部' },
      { key: 'market', label: '跳蚤市场' },
      { key: 'errand', label: '跑腿' },
      { key: 'lostFound', label: '失物招领' },
      { key: 'resource', label: '资源共享' },
      { key: 'carpool', label: '拼车' },
    ],
    statusFilters: [
      { key: 'all', label: '全部' },
      { key: 'reviewing', label: '审核中' },
      { key: 'published', label: '已发布' },
      { key: 'offline', label: '已下架' },
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
      const res = await fetchFeedList({ page: 1, pageSize: 50, user_role: 'publisher' });
      const data = res.data || res;
      const list = data.list || [];
      const allRecords = list.map((item) => {
        const status = normalizeStatus(item.review_status);
        return {
          id: item.id,
          title: item.title || '',
          type: item.feed_type || '',
          typeName: TYPE_MAP[item.feed_type] || item.feed_type_label || '',
          status,
          statusText: STATUS_MAP[item.review_status] || STATUS_MAP[status] || '已发布',
          time: item.created_at || '',
          url: URL_MAP[item.feed_type] || '',
        };
      }).sort((a, b) => new Date(b.time).getTime() - new Date(a.time).getTime());

      this.setData({ allRecords }, () => {
        this.filterRecords();
      });
    } catch (e) {
      console.error(e);
    }
  },

  filterRecords() {
    const { allRecords, activeType, activeStatus } = this.data;
    const filtered = allRecords.filter((record) => {
      if (activeType !== 'all' && record.type !== activeType) return false;
      if (activeStatus !== 'all' && record.status !== activeStatus) return false;
      return true;
    });
    this.setData({ filteredRecords: filtered });
  },

  onTypeFilter(e) {
    const { key } = e.currentTarget.dataset;
    if (key === this.data.activeType) return;
    this.setData({ activeType: key }, () => {
      this.filterRecords();
    });
  },

  onStatusFilter(e) {
    const { key } = e.currentTarget.dataset;
    if (key === this.data.activeStatus) return;
    this.setData({ activeStatus: key }, () => {
      this.filterRecords();
    });
  },

  onRecordClick(e) {
    const { url, id } = e.currentTarget.dataset;
    if (!url) {
      wx.showToast({ title: '该类型暂不支持查看详情', icon: 'none' });
      return;
    }
    wx.navigateTo({ url: `${url}?id=${id}` });
  },

  onBack() {
    wx.navigateBack();
  },

  onHeaderHeightChange(e) {
    this.setData({ headerHeight: e.detail.height });
  },
});
