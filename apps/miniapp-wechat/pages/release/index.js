import { getMenuButtonSafeArea } from '../../utils/navigation';
import { publishMarket } from '../../api/modules/market';
import { getUploadResultPath, uploadFile } from '../../api/modules/upload';
import { getMarketCategoryValueMap, getPaletteOptions, getReleaseScene } from '../../stores/config';

const PUBLISH_REDIRECTS = {
  errand: '/pages/errand/publish/index',
  meetup: '/pages/meetup/publish/index',
  lostFound: '/pages/lost-found/publish/index',
  resource: '/pages/resource/publish/index',
};

const PUBLISH_SCENE_OPTIONS = [
  { label: '发布闲置', value: 'market' },
  { label: '发布跑腿', value: 'errand' },
  { label: '发起组局', value: 'meetup' },
  { label: '发布登记', value: 'lostFound' },
  { label: '上传资料', value: 'resource' },
];

Page({
  data: {
    activeType: 'market',
    currentScene: getReleaseScene('market'),
    canPublish: false,
    menuTop: 40,
    menuSafeRight: 16,
    menuButtonHeight: 32,
    navHeight: 88,
    imageFiles: [],
    paletteOptions: getPaletteOptions(),
    form: {
      desc: '',
      price: '',
      originalPrice: '',
      category: '数码电子',
      condition: '99新',
      tradeMode: '支持校园面交',
      contact: '',
      cardBg: 'blue',
      startPlace: '',
      endPlace: '',
      deadline: '',
      eventTime: '',
      itemFeature: '',
      courseName: '',
      resourceScope: '',
    },
  },

  onLoad(options = {}) {
    this.applyNavigationSafeArea();
    this.bootstrapScene(options.type || '');
  },

  applyNavigationSafeArea() {
    const { right, top, height } = getMenuButtonSafeArea(10);
    this.setData({
      menuTop: top,
      menuSafeRight: right,
      menuButtonHeight: height,
      navHeight: top + height + 16,
    });
  },

  setActiveType(activeType) {
    const currentScene = getReleaseScene(activeType);
    this.setData(
      {
        activeType: currentScene.key,
        currentScene,
        paletteOptions: getPaletteOptions(),
        'form.category': currentScene.defaults.category,
        'form.condition': currentScene.defaults.condition,
        'form.tradeMode': currentScene.tradeMode,
      },
      () => {
        this.updatePublishState();
      },
    );
  },

  bootstrapScene(type) {
    if (type && PUBLISH_REDIRECTS[type]) {
      wx.redirectTo({ url: PUBLISH_REDIRECTS[type] });
      return;
    }

    if (!type) {
      this.openPublishSelector();
      return;
    }

    this.setActiveType('market');
  },

  openPublishSelector() {
    wx.showActionSheet({
      itemList: PUBLISH_SCENE_OPTIONS.map((item) => item.label),
      success: (res) => {
        const selected = PUBLISH_SCENE_OPTIONS[res.tapIndex] || PUBLISH_SCENE_OPTIONS[0];
        if (!selected) return;
        if (selected.value === 'market') {
          this.setActiveType('market');
          return;
        }
        wx.redirectTo({ url: PUBLISH_REDIRECTS[selected.value] });
      },
      fail: () => {
        this.setActiveType('market');
      },
    });
  },

  onBack() {
    if (getCurrentPages().length > 1) {
      wx.navigateBack({ delta: 1 });
      return;
    }
    wx.switchTab({ url: '/pages/home/index' });
  },

  getEventValue(e) {
    if (e.detail && typeof e.detail === 'object' && 'value' in e.detail) return e.detail.value;
    return e.detail;
  },

  onFieldChange(e) {
    const { field } = e.currentTarget.dataset;
    if (!field) return;
    this.setData(
      {
        [`form.${field}`]: this.getEventValue(e),
      },
      () => {
        this.updatePublishState();
      },
    );
  },

  onOptionChange(e) {
    const { field, value } = e.currentTarget.dataset;
    if (!field || !value) return;
    this.setData({
      [`form.${field}`]: value,
    });
  },

  onPaletteChange(e) {
    this.setData({
      'form.cardBg': e.currentTarget.dataset.value,
    });
  },

  chooseImages() {
    const remain = Math.max(0, 4 - this.data.imageFiles.length);
    if (!remain) return;

    wx.chooseMedia({
      count: remain,
      mediaType: ['image'],
      sourceType: ['album', 'camera'],
      success: (res) => {
        const nextFiles = res.tempFiles.map((item) => ({
          url: item.tempFilePath,
        }));
        this.setData({
          imageFiles: this.data.imageFiles.concat(nextFiles).slice(0, 4),
        });
      },
    });
  },

  previewImage(e) {
    const index = Number(e.currentTarget.dataset.index || 0);
    const urls = this.data.imageFiles.map((item) => item.url);
    if (!urls.length) return;
    wx.previewImage({
      current: urls[index],
      urls,
    });
  },

  removeImage(e) {
    const index = Number(e.currentTarget.dataset.index || 0);
    const imageFiles = this.data.imageFiles.filter((_, fileIndex) => fileIndex !== index);
    this.setData({ imageFiles });
  },

  chooseLocation() {
    wx.showToast({
      title: '位置选择待接入',
      icon: 'none',
      duration: 1500,
    });
  },

  updatePublishState() {
    const { form, activeType } = this.data;
    const hasDesc = form.desc.trim().length >= 4;
    const needsPrice = activeType === 'market' || activeType === 'errand';
    const hasPrice = !needsPrice || Boolean(form.price);
    this.setData({
      canPublish: hasDesc && hasPrice,
    });
  },

  confirmPublishWithoutImages() {
    return new Promise((resolve) => {
      wx.showModal({
        title: '图片上传暂不可用',
        content: '当前服务端还未开放图片上传接口，可先继续发布无图版本。',
        confirmText: '继续发布',
        cancelText: '返回修改',
        success: (res) => resolve(Boolean(res.confirm)),
        fail: () => resolve(false),
      });
    });
  },

  async uploadMarketImages(imageFiles) {
    if (!imageFiles.length) return [];

    try {
      const uploadResults = await Promise.all(
        imageFiles.map((file) => uploadFile(file.url, { scene: 'market' })),
      );
      return uploadResults.map(getUploadResultPath).filter(Boolean);
    } catch (error) {
      const confirmed = await this.confirmPublishWithoutImages();
      if (!confirmed) return null;
      return [];
    }
  },

  async release() {
    const { form, currentScene, activeType, imageFiles, canPublish } = this.data;

    if (!form.desc.trim()) {
      wx.showToast({ title: '先描述一下内容', icon: 'none' });
      return;
    }

    if (!canPublish) {
      wx.showToast({ title: currentScene.key === 'market' ? '请填写出手价格' : '请补充必要信息', icon: 'none' });
      return;
    }

    if (activeType === 'market') {
      try {
        wx.showLoading({ title: '发布中' });
        const imageUrls = await this.uploadMarketImages(imageFiles);
        if (imageUrls === null) return;

        const marketCategoryValueMap = getMarketCategoryValueMap();
        const res = await publishMarket({
          title: form.desc.trim().split('\n')[0].slice(0, 42) || '刚发布的校园闲置',
          desc: form.desc.trim(),
          price: form.price || '0',
          original_price: form.originalPrice || '',
          category: marketCategoryValueMap[form.category] || form.category,
          condition: form.condition,
          trade_mode: form.tradeMode,
          contact: form.contact || '',
          images: imageUrls,
        });

        const data = res.data || res || {};
        const productId = data.id || '';
        wx.redirectTo({
          url: `/pages/market/detail/index?id=${productId}`,
        });
      } catch (error) {
        wx.showToast({ title: error.message || '发布失败，请重试', icon: 'none' });
      } finally {
        wx.hideLoading();
      }
      return;
    }

    wx.reLaunch({ url: '/pages/home/index?oper=release' });
  },
});
