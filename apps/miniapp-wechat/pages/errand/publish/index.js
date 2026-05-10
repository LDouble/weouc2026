import { getMenuButtonSafeArea } from '../../../utils/navigation';
import { publishErrand } from '../../../api/modules/errand';
import { getUploadResultPath, uploadFile } from '../../../api/modules/upload';
import { saveHistoryAddress } from '../../../utils/addressStore';
import { PUBLISH_CATEGORIES } from '../data';

const DEADLINE_OPTIONS = [
  { label: '1小时未接单自动关闭', shortLabel: '1小时', hours: 1 },
  { label: '2小时', shortLabel: '2小时', hours: 2 },
  { label: '3小时', shortLabel: '3小时', hours: 3 },
  { label: '4小时', shortLabel: '4小时', hours: 4 },
  { label: '8小时', shortLabel: '8小时', hours: 8 },
  { label: '12小时', shortLabel: '12小时', hours: 12 },
  { label: '24小时', shortLabel: '24小时', hours: 24 },
];

function getCategoryLabel(value) {
  const category = PUBLISH_CATEGORIES.find((item) => item.value === value) || PUBLISH_CATEGORIES[0];
  return category.label;
}

function formatReward(value) {
  const reward = Number(value);
  if (!Number.isFinite(reward) || reward <= 0) return '0.00';
  return reward.toFixed(2);
}

function normalizeText(value) {
  if (value === null || value === undefined) return '';
  if (typeof value === 'object') return normalizeText(value.value);
  return String(value).trim();
}

function confirmWithoutImages() {
  return new Promise((resolve) => {
    wx.showModal({
      title: '图片上传暂不可用',
      content: '当前服务端还未开放图片上传接口，可先继续发布无图任务。',
      confirmText: '继续发布',
      cancelText: '返回修改',
      success: (res) => resolve(Boolean(res.confirm)),
      fail: () => resolve(false),
    });
  });
}

Page({
  data: {
    categoryList: PUBLISH_CATEGORIES,
    deadlineOptions: DEADLINE_OPTIONS,
    canPublish: false,
    submitting: false,
    menuTop: 40,
    menuButtonHeight: 32,
    navHeight: 88,
    imageFiles: [],
    form: {
      category: 'parcel',
      todo: '',
      startPlace: '',
      endPlace: '',
      contact: '',
      deadlineHours: 1,
      reward: '5.00',
      urgent: false,
    },
  },

  onLoad() {
    this.applyNavigationSafeArea();
    this.updatePublishState();
  },

  onShow() {},

  applyNavigationSafeArea() {
    const { top, height } = getMenuButtonSafeArea(10);
    this.setData({
      menuTop: top,
      menuButtonHeight: height,
      navHeight: top + height + 16,
    });
  },

  onBack() {
    if (getCurrentPages().length > 1) {
      wx.navigateBack({ delta: 1 });
      return;
    }
    wx.navigateTo({ url: '/pages/errand/index' });
  },

  getEventValue(e) {
    if (e.detail && typeof e.detail === 'object' && 'value' in e.detail) return e.detail.value;
    return e.detail;
  },

  onCategorySelect(e) {
    const { value } = e.currentTarget.dataset;
    if (!value) return;

    this.setData({
      'form.category': value,
    });
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

  onUrgentToggle() {
    this.setData({
      'form.urgent': !this.data.form.urgent,
    });
  },

  onDeadlineSelect(e) {
    const hours = Number(e.currentTarget.dataset.hours);
    if (!hours) return;

    this.setData({
      'form.deadlineHours': hours,
    });
  },

  onStartPlaceChange(e) {
    this.setData({ 'form.startPlace': this.getEventValue(e) }, () => {
      this.updatePublishState();
    });
  },

  onEndPlaceChange(e) {
    this.setData({ 'form.endPlace': this.getEventValue(e) }, () => {
      this.updatePublishState();
    });
  },

  chooseImages() {
    const remain = Math.max(0, 3 - this.data.imageFiles.length);
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
          imageFiles: this.data.imageFiles.concat(nextFiles).slice(0, 3),
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

  updatePublishState() {
    const { form } = this.data;
    const hasTodo = normalizeText(form.todo).length >= 4;
    const hasRoute = normalizeText(form.startPlace) && normalizeText(form.endPlace);
    const hasContact = normalizeText(form.contact).length >= 2;
    const hasReward = Number(form.reward) > 0;
    this.setData({ canPublish: hasTodo && hasRoute && hasContact && hasReward });
  },

  async release() {
    if (!this.data.canPublish) {
      wx.showToast({ title: '请补充事项、路线、联系方式和赏金', icon: 'none' });
      return;
    }
    if (this.data.submitting) return;

    const { form, imageFiles } = this.data;
    const deadlineAt = new Date(Date.now() + Number(form.deadlineHours) * 60 * 60 * 1000);
    const title = normalizeText(form.todo);
    const routeStart = normalizeText(form.startPlace);
    const routeEnd = normalizeText(form.endPlace);
    const contact = normalizeText(form.contact);

    const payload = {
      title,
      desc: title,
      category: form.category,
      route_start: routeStart,
      route_end: routeEnd,
      deadline: deadlineAt.toISOString(),
      reward: formatReward(form.reward),
      contact,
      urgent: form.urgent,
    };

    this.setData({ submitting: true });
    try {
      if (imageFiles.length) {
        try {
          const uploadResults = await Promise.all(
            imageFiles.map((file) => uploadFile(file.url, { scene: 'errand' })),
          );
          payload.images = uploadResults.map(getUploadResultPath).filter(Boolean);
        } catch (error) {
          const confirmed = await confirmWithoutImages();
          if (!confirmed) return;
          payload.images = [];
        }
      }

      const res = await publishErrand(payload);
      const data = res.data || res || {};
      const id = data.id || '';
      saveHistoryAddress(routeStart);
      saveHistoryAddress(routeEnd);
      wx.showToast({ title: '发布成功', icon: 'success' });
      setTimeout(() => {
        wx.redirectTo({ url: `/pages/errand/detail/index?id=${id}` });
      }, 1500);
    } catch (err) {
      wx.showToast({ title: err.message || '发布失败，请重试', icon: 'none' });
    } finally {
      this.setData({ submitting: false });
    }
  },
});
