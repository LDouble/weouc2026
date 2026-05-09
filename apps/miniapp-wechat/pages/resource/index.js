import { RESOURCE_CATEGORIES } from './data';
import {
  ensureResourcePinned,
  loadResourceList,
  mergeResourceItems,
} from '../../services/resourceService';

Page({
  data: {
    activeCategory: 'all',
    categoryList: RESOURCE_CATEGORIES,
    resourceList: [],
    visibleResources: [],
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
    this.refreshResources();
  },

  async forceInsertResource(items) {
    const { insertId } = this.data;
    return ensureResourcePinned(items, insertId);
  },

  async refreshResources() {
    this.setData({ loading: true, page: 1 });
    const { activeCategory, searchKeyword, pageSize } = this.data;
    const params = { page: 1, pageSize };
    if (activeCategory !== 'all') params.category = activeCategory;
    if (searchKeyword && searchKeyword.trim()) params.keyword = searchKeyword.trim();

    try {
      const result = await loadResourceList(params);
      const mapped = await this.forceInsertResource(result.items);
      this.setData({
        resourceList: mapped,
        visibleResources: mapped,
        total: Math.max(result.total, mapped.length),
        page: 1,
        loading: false,
      });
    } catch (e) {
      this.setData({ loading: false });
    }
  },

  filterResources() {
    this.refreshResources();
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
    this.filterResources();
  },

  onSearchClear() {
    this.setData({ searchKeyword: '' }, () => {
      this.filterResources();
    });
  },

  onCategoryChange(e) {
    const { value } = e.detail;
    if (!value || value === this.data.activeCategory) return;

    this.setData({ activeCategory: value }, () => {
      this.filterResources();
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
      const result = await loadResourceList(params);
      const newVisible = mergeResourceItems(this.data.visibleResources, result.items);
      this.setData({
        resourceList: mergeResourceItems(this.data.resourceList, result.items),
        visibleResources: newVisible,
        page: nextPage,
        loading: false,
      });
    } catch (e) {
      this.setData({ loading: false });
    }
  },

  goRelease() {
    wx.navigateTo({ url: '/pages/resource/publish/index' });
  },

  onOpenResource(e) {
    const resource = e.currentTarget.dataset.resource || {};
    const files = resource.files || [];
    const firstFile = files[0] || {};
    const url = firstFile.url || resource.downloadUrl;
    if (!url) {
      wx.showToast({ title: '暂无可打开的文件', icon: 'none' });
      return;
    }

    if (resource.fileType === '图片') {
      const imageUrls = files
        .filter((file) => (file.file_type || '').includes('image'))
        .map((file) => file.url);
      wx.previewImage({
        current: url,
        urls: imageUrls.length ? imageUrls : [url],
      });
      return;
    }

    wx.showLoading({ title: '打开中' });
    wx.downloadFile({
      url,
      success: (res) => {
        if (res.statusCode !== 200) {
          wx.showToast({ title: '下载失败', icon: 'none' });
          return;
        }
        wx.openDocument({
          filePath: res.tempFilePath,
          showMenu: true,
          fail: () => wx.showToast({ title: '暂不支持打开该文件', icon: 'none' }),
        });
      },
      fail: () => wx.showToast({ title: '下载失败', icon: 'none' }),
      complete: () => wx.hideLoading(),
    });
  },
});
