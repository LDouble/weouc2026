import {
  getMarketPageConfig,
  loadMarketList,
  toggleMarketFavorite,
} from '../../services/marketService';

Page({
  data: {
    activeCategory: 'all',
    searchKeyword: '',
    productList: [],
    categoryList: getMarketPageConfig().categoryList,
    productColumns: {
      left: [],
      right: [],
    },
    headerHeight: 260,
    page: 1,
    pageSize: 20,
    total: 0,
    loading: false,
  },

  onShow() {
    this.refreshMarketConfig();
    this.refreshProducts();
  },

  refreshMarketConfig() {
    const categoryList = getMarketPageConfig().categoryList;
    const hasActiveCategory = categoryList.some((item) => item.value === this.data.activeCategory);
    const defaultCategory = categoryList[0] ? categoryList[0].value : 'all';

    this.setData({
      categoryList,
      activeCategory: hasActiveCategory ? this.data.activeCategory : defaultCategory,
    });
  },

  async refreshProducts() {
    const { activeCategory, searchKeyword, pageSize } = this.data;
    const category = activeCategory === 'all' ? '' : activeCategory;
    const keyword = searchKeyword.trim();

    this.setData({ loading: true, page: 1 });

    try {
      const result = await loadMarketList({ category, keyword, page: 1, pageSize });
      this.setData(
        {
          productList: result.items,
          total: result.total,
          loading: false,
        },
        () => {
          this.filterProducts();
        },
      );
    } catch (error) {
      this.setData({ loading: false });
      wx.showToast({ title: '加载失败', icon: 'none' });
    }
  },

  async onLoadMore() {
    const { loading, page, pageSize, total, productList, activeCategory, searchKeyword } = this.data;
    if (loading) return;
    if (productList.length >= total) return;

    const nextPage = page + 1;
    const category = activeCategory === 'all' ? '' : activeCategory;
    const keyword = searchKeyword.trim();

    this.setData({ loading: true });

    try {
      const result = await loadMarketList({ category, keyword, page: nextPage, pageSize });
      this.setData(
        {
          productList: productList.concat(result.items),
          page: nextPage,
          total: result.total,
          loading: false,
        },
        () => {
          this.filterProducts();
        },
      );
    } catch (error) {
      this.setData({ loading: false });
    }
  },

  onReachBottom() {
    this.onLoadMore();
  },

  onBack() {
    if (getCurrentPages().length > 1) {
      wx.navigateBack({ delta: 1 });
      return;
    }
    wx.switchTab({ url: '/pages/home/index' });
  },

  onSearchInput(e) {
    this.setData({ searchKeyword: e.detail.value });
  },

  onSearchConfirm() {
    this.refreshProducts();
  },

  onSearchClear() {
    this.setData({ searchKeyword: '' }, () => {
      this.refreshProducts();
    });
  },

  onCategoryChange(e) {
    this.setData({ activeCategory: e.detail.value }, () => {
      this.refreshProducts();
    });
  },

  onHeaderHeightChange(e) {
    const { height } = e.detail;
    if (height) this.setData({ headerHeight: height });
  },

  filterProducts() {
    const left = [];
    const right = [];
    this.data.productList.forEach((item, index) => {
      if (index % 2 === 0) {
        left.push(item);
      } else {
        right.push(item);
      }
    });
    this.setData({ productColumns: { left, right } });
  },

  onProductSelect(e) {
    const { product } = e.detail;
    wx.navigateTo({ url: `/pages/market/detail/index?id=${product.id}` });
  },

  async onFavorite(e) {
    const { product } = e.detail;

    try {
      const productList = this.data.productList.map((item) => {
        if (item.id !== product.id) return item;
        return toggleMarketFavorite(item);
      });
      const resolvedList = await Promise.all(productList);
      this.setData({ productList: resolvedList }, () => {
        this.filterProducts();
      });
    } catch (error) {
      wx.showToast({ title: '操作失败', icon: 'none' });
    }
  },

  goRelease() {
    wx.navigateTo({ url: '/pages/release/index?type=market' });
  },
});
