import { CARPOOL_TIME_FILTERS } from './data';
import { deleteCarpool, fetchCarpoolDetail, fetchCarpoolList } from '../../api/modules/carpool';
import { getStatusDisplay, getReviewStatus, normalizeContactFields } from '../../services/shared';

const CARPOOL_STATUS_OVERRIDES = {
  published: { label: '可拼' },
};

function mapCarpoolItem(item) {
  const extra = item.extra || {};
  const status = getReviewStatus(item);
  const statusDisplay = getStatusDisplay(status, CARPOOL_STATUS_OVERRIDES);
  const contactFields = normalizeContactFields(item);
  return {
    id: item.id,
    category: item.category || extra.category || '',
    from: item.from || extra.from || '',
    to: item.to || extra.to || '',
    time: item.time || extra.time || '',
    type: item.type || extra.type || '',
    seatsText: item.seats_text || extra.seats_text || '',
    price: item.price || extra.price || '',
    status: statusDisplay.label,
    statusTone: statusDisplay.tone,
    sponsor: item.publisher || '',
    sponsorInitial: item.publisher_initial || '',
    sponsorTag: item.created_at || '',
    note: item.note || extra.note || '',
    tags: item.tags || extra.tags || [],
    contact: contactFields.contact,
    canViewContact: contactFields.canViewContact,
    contactHiddenReason: contactFields.contactHiddenReason,
    isOwner: Boolean(item.is_owner),
    canDelete: Boolean(item.can_delete),
  };
}

function moveItemToTop(items, id) {
  if (!id) return items;
  const target = items.find((item) => item.id === id);
  if (!target) return items;
  return [target].concat(items.filter((item) => item.id !== id));
}

function dedupeById(items) {
  const seen = {};
  return items.filter((item) => {
    if (!item.id || seen[item.id]) return false;
    seen[item.id] = true;
    return true;
  });
}

Page({
  data: {
    timeFilters: CARPOOL_TIME_FILTERS,
    activeTimeFilter: 'all',
    tripList: [],
    visibleTrips: [],
    searchKeyword: '',
    headerHeight: 194,
    page: 1,
    pageSize: 20,
    total: 0,
    loading: false,
    insertId: '',
  },

  onLoad(options = {}) {
    if (options.insertId) {
      this.setData({ insertId: options.insertId });
    }
  },

  onShow() {
    this.refreshTrips();
  },

  async forceInsertTrip(items) {
    const { insertId } = this.data;
    if (!insertId) return items;

    const moved = moveItemToTop(items, insertId);
    if (moved[0] && moved[0].id === insertId) return moved;

    try {
      const res = await fetchCarpoolDetail(insertId);
      const detail = mapCarpoolItem(res.data || res);
      return [detail].concat(items.filter((item) => item.id !== insertId));
    } catch (e) {
      return items;
    }
  },

  async refreshTrips() {
    const { activeTimeFilter, searchKeyword, pageSize } = this.data;
    this.setData({ loading: true, page: 1 });

    const params = { page: 1, pageSize };
    if (activeTimeFilter !== 'all') params.category = activeTimeFilter;
    if (searchKeyword) params.keyword = searchKeyword.trim();

    try {
      const res = await fetchCarpoolList(params);
      const list = (res.data && res.data.list) || [];
      const tripList = await this.forceInsertTrip(list.map(mapCarpoolItem));
      this.setData({
        tripList,
        visibleTrips: tripList,
        total: Math.max((res.data && res.data.total) || 0, tripList.length),
        page: 1,
        loading: false,
      });
    } catch (e) {
      console.error('[carpool] refreshTrips failed:', e);
      this.setData({ loading: false });
      wx.showToast({ title: '加载失败，请重试', icon: 'none' });
    }
  },

  filterTrips() {
    const { tripList, searchKeyword } = this.data;
    if (!searchKeyword) {
      this.setData({ visibleTrips: tripList });
      return;
    }
    const keyword = searchKeyword.trim().toLowerCase();
    const visibleTrips = tripList.filter(
      (item) =>
        (item.from && item.from.toLowerCase().includes(keyword)) ||
        (item.to && item.to.toLowerCase().includes(keyword)) ||
        (item.sponsor && item.sponsor.toLowerCase().includes(keyword))
    );
    this.setData({ visibleTrips });
  },

  async onLoadMore() {
    const { loading, page, pageSize, total, tripList, activeTimeFilter, searchKeyword } = this.data;
    if (loading) return;
    if (tripList.length >= total) return;

    const nextPage = page + 1;
    this.setData({ loading: true });

    const params = { page: nextPage, pageSize };
    if (activeTimeFilter !== 'all') params.category = activeTimeFilter;
    if (searchKeyword) params.keyword = searchKeyword.trim();

    try {
      const res = await fetchCarpoolList(params);
      const list = (res.data && res.data.list) || [];
      const newTrips = list.map(mapCarpoolItem);
      const merged = dedupeById(tripList.concat(newTrips));
      this.setData({
        tripList: merged,
        visibleTrips: merged,
        total: (res.data && res.data.total) || 0,
        page: nextPage,
        loading: false,
      });
    } catch (e) {
      console.error('[carpool] onLoadMore failed:', e);
      this.setData({ loading: false });
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
    this.setData({ searchKeyword: e.detail.value }, () => {
      this.filterTrips();
    });
  },

  onSearchConfirm() {
    this.refreshTrips();
  },

  onSearchClear() {
    this.setData({ searchKeyword: '' }, () => {
      this.refreshTrips();
    });
  },

  onTimeFilterChange(e) {
    const { value } = e.detail;
    if (!value || value === this.data.activeTimeFilter) return;
    this.setData({ activeTimeFilter: value }, () => {
      this.refreshTrips();
    });
  },

  onHeaderHeightChange(e) {
    const { height } = e.detail;
    if (height) this.setData({ headerHeight: height });
  },

  onCopyContact(e) {
    const { contact, canViewContact } = e.currentTarget.dataset;
    if (canViewContact === false || canViewContact === 'false') {
      wx.showModal({
        title: '无法查看联系方式',
        content: '绑定教务后即可查看联系方式，是否前往绑定？',
        confirmText: '前往绑定',
        cancelText: '暂不需要',
        success: (res) => {
          if (res.confirm) {
            wx.navigateTo({ url: '/pages/edu-bind/index' });
          }
        },
      });
      return;
    }
    if (!contact || contact === '站内私信') {
      wx.showToast({ title: '暂无联系方式', icon: 'none' });
      return;
    }
    wx.setClipboardData({
      data: contact,
      success: () => {
        wx.showToast({ title: '联系方式已复制', icon: 'success' });
      },
    });
  },

  onCancelPublish(e) {
    const { id } = e.currentTarget.dataset;
    if (!id) return;
    wx.showModal({
      title: '取消发布',
      content: '确定要取消这条拼车信息吗？',
      confirmText: '确定取消',
      confirmColor: '#dc2626',
      success: async (res) => {
        if (!res.confirm) return;
        try {
          await deleteCarpool(id);
          wx.showToast({ title: '已取消发布', icon: 'success' });
          this.refreshTrips();
        } catch (err) {
          console.error('[carpool] cancel failed:', err);
          wx.showToast({ title: '取消失败，请重试', icon: 'none' });
        }
      },
    });
  },

  onCreateTap() {
    wx.navigateTo({ url: '/pages/carpool/publish/index' });
  },
});
