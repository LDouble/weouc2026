import { fetchMarketList, favoriteMarket } from '../../api/modules/market';
import { getMarketCategories } from '../../store/config';

function buildCategoryLabelMap() {
  const categories = getMarketCategories();
  return categories.reduce((map, item) => {
    if (item.value && item.value !== 'all') {
      map[item.value] = item.label;
    }
    return map;
  }, {});
}

function mapMarketItem(item, categoryLabelMap) {
  const extra = item.extra || {};
  const isWanted = extra.category === 'wanted';
  const imageUrls = extra.images && extra.images.length ? extra.images : [];
  if (!imageUrls.length && item.image) imageUrls.push(item.image);
  const images = imageUrls.map((url, index) => ({
    key: `img-${item.id}-${index}`,
    url,
    variant: 'center',
  }));

  return {
    id: item.id,
    category: extra.category || 'life',
    categoryName: categoryLabelMap[extra.category] || extra.category || '',
    mode: isWanted ? 'wanted' : '',
    title: item.title || '',
    price: extra.price || '0',
    originalPrice: extra.original_price || '',
    condition: extra.condition || '',
    seller: item.publisher || '',
    sellerInitial: item.publisher_initial || '',
    avatarTone: 'green',
    likes: extra.likes || item.likes || 0,
    liked: extra.is_favorited || item.liked || false,
    image: item.image || (images[0] ? images[0].url : ''),
    images,
    imageHeight: images.length ? 342 : 354,
    cardBg: 'blue',
    detail: item.desc || '',
    contact: extra.contact || '',
  };
}

Page({
  data: {
    activeCategory: 'all',
    searchKeyword: '',
    productList: [],
    categoryList: getMarketCategories(),
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
    const categoryList = getMarketCategories();
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
      const res = await fetchMarketList({ category, keyword, page: 1, pageSize });
      const categoryLabelMap = buildCategoryLabelMap();
      const list = (res.data && res.data.list) || [];
      const productList = list.map((item) => mapMarketItem(item, categoryLabelMap));
      this.setData(
        {
          productList,
          total: (res.data && res.data.total) || 0,
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
      const res = await fetchMarketList({ category, keyword, page: nextPage, pageSize });
      const categoryLabelMap = buildCategoryLabelMap();
      const list = (res.data && res.data.list) || [];
      const newProducts = list.map((item) => mapMarketItem(item, categoryLabelMap));
      this.setData(
        {
          productList: productList.concat(newProducts),
          page: nextPage,
          total: (res.data && res.data.total) || 0,
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
    const { productList } = this.data;
    const productColumns = productList.reduce(
      (columns, item, index) => {
        const target = item.column || (index % 2 === 0 ? 'left' : 'right');
        columns[target].push(item);
        return columns;
      },
      { left: [], right: [] },
    );
    this.setData({ productColumns });
  },

  onProductSelect(e) {
    const { product } = e.detail;
    wx.navigateTo({ url: `/pages/market/detail/index?id=${product.id}` });
  },

  async onFavorite(e) {
    const { product } = e.detail;
    const action = product.liked ? 'remove' : 'add';

    try {
      await favoriteMarket(product.id, action);
      const productList = this.data.productList.map((item) => {
        if (item.id !== product.id) return item;
        const liked = !item.liked;
        return {
          ...item,
          liked,
          likes: item.likes + (liked ? 1 : -1),
        };
      });
      this.setData({ productList }, () => {
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
