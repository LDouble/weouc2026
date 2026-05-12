import { getMenuButtonSafeArea } from '../../../utils/navigation';
import { fetchMarketDetail, deleteMarket } from '../../../api/modules/market';
import { toggleMarketFavorite } from '../../../services/marketService';
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

function mapMarketDetail(item, canViewContact, isOwner, canEdit, canDelete, canFavorite) {
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
    status: item.status || 'published',
    contact: canViewContact && hasContact ? (extra.contact || '') : '',
    canViewContact,
    isOwner: Boolean(isOwner),
    canEdit: Boolean(canEdit),
    canDelete: Boolean(canDelete),
    canFavorite: Boolean(canFavorite),
    canCopyContact: canViewContact && hasContact,
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
      const product = mapMarketDetail(
        data,
        data.can_view_contact || false,
        data.is_owner || false,
        data.can_edit || false,
        data.can_delete || false,
        data.can_favorite || false,
      );
      this.setData({
        product,
        currentImageIndex: 0,
        wanted: product.liked || false,
        loading: false,
      });
    } catch (error) {
      this.setData({ loading: false });
      wx.showToast({ title: error.message || '加载失败，请重试', icon: 'none' });
    }
  },

  onBack() {
    if (getCurrentPages().length > 1) {
      wx.navigateBack({ delta: 1 });
      return;
    }
    wx.navigateTo({ url: '/pages/market/index' });
  },

  async onToggleWanted() {
    try {
      const product = await toggleMarketFavorite(this.data.product);
      this.setData({
        wanted: product.liked,
        product,
      });
    } catch (error) {
      wx.showToast({ title: error.message || '操作失败，请重试', icon: 'none' });
    }
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

  onDelete() {
    const { product } = this.data;
    if (!product.canDelete) return;
    wx.showModal({
      title: '下架商品',
      content: '下架后该商品将不再展示，确认下架？',
      confirmText: '确认下架',
      confirmColor: '#dc2626',
      success: async (res) => {
        if (!res.confirm) return;
        try {
          await deleteMarket(product.id);
          wx.showToast({ title: '已下架', icon: 'success' });
          setTimeout(() => {
            if (getCurrentPages().length > 1) {
              wx.navigateBack({ delta: 1 });
              return;
            }
            wx.redirectTo({ url: '/pages/market/index' });
          }, 500);
        } catch (error) {
          wx.showToast({ title: (error && error.message) || '下架失败，请重试', icon: 'none' });
        }
      },
    });
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
