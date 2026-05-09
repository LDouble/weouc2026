import { getMenuButtonSafeArea } from '../../../utils/navigation';
import { fetchMarketDetail } from '../../../api/modules/market';
import { getMarketCategories } from '../../../stores/config';

function buildCategoryLabelMap() {
  const categories = getMarketCategories();
  return categories.reduce((map, item) => {
    if (item.value && item.value !== 'all') {
      map[item.value] = item.label;
    }
    return map;
  }, {});
}

function mapMarketDetail(item, canViewContact = false) {
  const extra = item.extra || {};
  const categoryLabelMap = buildCategoryLabelMap();
  const isWanted = extra.category === 'wanted';
  const imageUrls = extra.images && extra.images.length ? extra.images : [];
  if (!imageUrls.length && item.image) imageUrls.push(item.image);
  const images = imageUrls.map((url, index) => ({
    key: `img-${item.id}-${index}`,
    url,
    variant: 'center',
  }));
  const hasContact = Boolean(extra.contact);

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
    cardBg: 'blue',
    detail: item.desc || '',
    contact: canViewContact && hasContact ? (extra.contact || '') : '',
    canViewContact,
    tags: [
      { label: categoryLabelMap[extra.category] || extra.category || '', icon: isWanted ? 'search' : 'shop', tone: 'slate' },
      { label: extra.trade_mode || '', icon: 'location', tone: 'blue' },
    ].filter((t) => t.label),
    sellerMeta: {
      verified: true,
      rating: '优秀',
      onSale: '1件',
      deals: '0单',
    },
  };
}

Page({
  data: {
    product: {
      images: [],
      tags: [],
      sellerMeta: {},
      likes: 0,
    },
    currentImageIndex: 0,
    wanted: false,
    menuTop: 40,
    menuSafeRight: 16,
    menuButtonHeight: 32,
    loading: false,
  },

  onLoad(options = {}) {
    this.applyNavigationSafeArea();
    this.loadProduct(options.id);
  },

  applyNavigationSafeArea() {
    const { right, top, height } = getMenuButtonSafeArea(8);
    this.setData({
      menuTop: top,
      menuSafeRight: right,
      menuButtonHeight: height,
    });
  },

  async loadProduct(id) {
    if (!id) return;
    this.setData({ loading: true });
    try {
      const res = await fetchMarketDetail(id);
      const data = res.data || res;
      const canViewContact = data.can_view_contact || false;
      const product = mapMarketDetail(data, canViewContact);
      this.setData({
        product,
        currentImageIndex: 0,
        wanted: product.liked || false,
        loading: false,
      });
    } catch (error) {
      this.setData({ loading: false });
      wx.showToast({ title: '加载失败', icon: 'none' });
    }
  },

  onBack() {
    if (getCurrentPages().length > 1) {
      wx.navigateBack({ delta: 1 });
      return;
    }
    wx.navigateTo({ url: '/pages/market/index' });
  },

  onToggleWanted() {
    const wanted = !this.data.wanted;
    this.setData({
      wanted,
      product: {
        ...this.data.product,
        likes: this.data.product.likes + (wanted ? 1 : -1),
      },
    });
  },

  onImageChange(e) {
    this.setData({ currentImageIndex: e.detail.current });
  },

  onPreviewImages(e) {
    const { product } = this.data;
    const images = product.images || [];
    if (!images.length) return;

    const index = e.currentTarget.dataset.index || 0;
    wx.previewImage({
      current: images[index].url,
      urls: images.map((item) => item.url),
    });
  },

  onContact() {
    const { contact, canViewContact } = this.data.product;
    if (!canViewContact) {
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
    if (!contact) {
      wx.showToast({ title: '卖家暂未留下联系方式', icon: 'none' });
      return;
    }
    wx.setClipboardData({
      data: contact,
      success: () => {
        wx.showToast({ title: '联系方式已复制', icon: 'success' });
      },
    });
  },

  onSellerHome() {
    wx.showToast({ title: '个人主页接入中', icon: 'none' });
  },

  onShareAppMessage() {
    const { product } = this.data;
    return {
      title: product.title || '校园跳蚤市场',
      path: `/pages/market/detail/index?id=${product.id || ''}`,
      imageUrl: product.images && product.images[0] ? product.images[0].url : product.image || '',
    };
  },
});
