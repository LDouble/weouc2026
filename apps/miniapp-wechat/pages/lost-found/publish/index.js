import { getMenuButtonSafeArea } from '../../../utils/navigation';
import { LOST_FOUND_CATEGORIES, LOST_FOUND_TYPES } from '../data';
import { publishLostFound } from '../../../api/modules/lostFound';
import { saveHistoryAddress } from '../../../utils/addressStore';

const PUBLISH_CATEGORIES = LOST_FOUND_CATEGORIES.filter((item) => item.value !== 'all');

function padNumber(value) {
  return value < 10 ? `0${value}` : `${value}`;
}

function formatPickerDate(date) {
  return `${date.getFullYear()}-${padNumber(date.getMonth() + 1)}-${padNumber(date.getDate())}`;
}

function formatDisplayDate(dateText) {
  if (!dateText) return '';

  const [yearText, monthText, dayText] = dateText.split('-');
  const year = Number(yearText);
  const month = Number(monthText);
  const day = Number(dayText);
  const targetDate = new Date(year, month - 1, day);
  const now = new Date();
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
  const tomorrow = new Date(now.getFullYear(), now.getMonth(), now.getDate() + 1);

  if (targetDate.getTime() === today.getTime()) return '今天';
  if (targetDate.getTime() === tomorrow.getTime()) return '明天';
  return `${month}月${day}日`;
}

function getDateRange() {
  const start = new Date();
  const end = new Date(start.getTime() + 180 * 24 * 60 * 60 * 1000);
  return {
    start: formatPickerDate(start),
    end: formatPickerDate(end),
  };
}

function getCategoryLabel(value) {
  const category = LOST_FOUND_CATEGORIES.find((item) => item.value === value);
  return category ? category.label : '';
}

Page({
  data: {
    canPublish: false,
    menuTop: 40,
    menuButtonHeight: 32,
    navHeight: 88,
    dateRange: getDateRange(),
    categoryList: PUBLISH_CATEGORIES,
    form: {
      type: 'lost',
      title: '',
      category: '',
      desc: '',
      location: '',
      date: '',
      time: '',
      contact: '',
    },
  },

  onLoad() {
    this.applyNavigationSafeArea();
    this.updatePublishState();
  },

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
    wx.redirectTo({ url: '/pages/lost-found/index' });
  },

  getEventValue(e) {
    if (e.detail && typeof e.detail === 'object' && 'value' in e.detail) return e.detail.value;
    return e.detail;
  },

  onTypeSelect(e) {
    const { type } = e.currentTarget.dataset;
    if (!type) return;

    this.setData({ 'form.type': type }, () => {
      this.updatePublishState();
    });
  },

  onCategorySelect(e) {
    const { value } = e.currentTarget.dataset;
    if (!value) return;

    const nextCategory = this.data.form.category === value ? '' : value;
    this.setData({ 'form.category': nextCategory }, () => {
      this.updatePublishState();
    });
  },

  onInputChange(e) {
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

  onDateChange(e) {
    this.setData(
      {
        'form.date': this.getEventValue(e),
      },
      () => {
        this.updatePublishState();
      },
    );
  },

  onTimeChange(e) {
    this.setData(
      {
        'form.time': this.getEventValue(e),
      },
      () => {
        this.updatePublishState();
      },
    );
  },

  updatePublishState() {
    const { form } = this.data;
    const hasTitle = form.title.trim();
    const hasCategory = form.category.trim();
    const hasLocation = form.location.trim();
    const hasDate = form.date && form.time;
    const hasContact = form.contact.trim().length >= 2;

    this.setData({
      canPublish: Boolean(hasTitle && hasCategory && hasLocation && hasDate && hasContact),
    });
  },

  async release() {
    const { form } = this.data;
    const missing = [];
    if (!form.title.trim()) missing.push('物品名称');
    if (!form.category) missing.push('物品类别');
    if (!form.location.trim()) missing.push('地点');
    if (!form.date) missing.push('日期');
    if (!form.time) missing.push('时间');
    if (form.contact.trim().length < 2) missing.push('联系方式');

    if (missing.length) {
      wx.showToast({ title: `请补充${missing.join('、')}`, icon: 'none' });
      return;
    }

    try {
      const res = await publishLostFound({
        type: form.type,
        category: form.category,
        title: form.title.trim(),
        desc: form.desc.trim(),
        location: form.location.trim(),
        event_time: `${form.date} ${form.time || ''}`.trim(),
        item_feature: '',
        contact: form.contact.trim(),
        reward: '',
      });
      const data = res.data || res || {};
      const detailId = data.id || '';
      saveHistoryAddress(form.location.trim());
      wx.showToast({ title: '发布成功', icon: 'success' });
      setTimeout(() => {
        if (detailId) {
          wx.redirectTo({ url: `/pages/lost-found/detail/index?id=${detailId}` });
          return;
        }
        wx.redirectTo({ url: '/pages/lost-found/index' });
      }, 360);
    } catch (error) {
      wx.showToast({ title: error.message || '发布失败，请重试', icon: 'none' });
    }
  },
});
