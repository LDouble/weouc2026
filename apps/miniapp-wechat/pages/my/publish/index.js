import { fetchFeedList } from '~/api/modules/feed';
import { getReviewStatus } from '../../../services/shared';

const TYPE_MAP = {
  market: '跳蚤市场',
  errand: '跑腿',
  meetup: '组局',
  lostFound: '失物招领',
  resource: '资源共享',
  carpool: '拼车',
};

const ROUTE_MAP = {
  market: { url: '/pages/market/detail/index', queryKey: 'id' },
  errand: { url: '/pages/errand/detail/index', queryKey: 'id' },
  meetup: { url: '/pages/meetup/detail/index', queryKey: 'id' },
  lostFound: { url: '/pages/lost-found/detail/index', queryKey: 'id' },
  resource: { url: '/pages/resource/index', queryKey: 'insertId' },
  carpool: { url: '/pages/carpool/index', queryKey: 'insertId' },
};

const STATUS_MAP = {
  pending: '审核中',
  reviewing: '审核中',
  published: '已发布',
  open: '已发布',
  rejected: '已下架',
  offline: '已下架',
  cancelled: '已取消',
};

function normalizeStatus(status) {
  if (status === 'pending' || status === 'reviewing') return 'reviewing';
  if (status === 'published' || status === 'open') return 'published';
  if (status === 'rejected' || status === 'offline' || status === 'cancelled') return 'offline';
  return status || 'published';
}

function buildRecordUrl(type, id) {
  const route = ROUTE_MAP[type];
  if (!route || !route.url || !id) return '';
  return `${route.url}?${route.queryKey || 'id'}=${id}`;
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
      { key: 'meetup', label: '组局' },
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
        const reviewStatus = getReviewStatus(item);
        const status = normalizeStatus(reviewStatus);
        return {
          id: item.id,
          title: item.title || '',
          type: item.feed_type || '',
          typeName: TYPE_MAP[item.feed_type] || item.feed_type_label || '',
          status,
          statusText: STATUS_MAP[reviewStatus] || STATUS_MAP[status] || '已发布',
          time: item.created_at || '',
          url: buildRecordUrl(item.feed_type, item.id),
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
    const { url } = e.currentTarget.dataset;
    if (!url) {
      wx.showToast({ title: '该类型暂不支持查看详情', icon: 'none' });
      return;
    }
    wx.navigateTo({ url });
  },

  onBack() {
    wx.navigateBack();
  },

  onHeaderHeightChange(e) {
    this.setData({ headerHeight: e.detail.height });
  },
});
