import { PUBLISH_CATEGORIES } from '../data';
import { publishMeetup } from '../../../api/modules/meetup';

function padNumber(value) {
  return value < 10 ? `0${value}` : `${value}`;
}

function formatPickerDate(date) {
  return `${date.getFullYear()}-${padNumber(date.getMonth() + 1)}-${padNumber(date.getDate())}`;
}

function getDateRange() {
  const start = new Date();
  const end = new Date(start.getTime() + 180 * 24 * 60 * 60 * 1000);
  return {
    start: formatPickerDate(start),
    end: formatPickerDate(end),
  };
}

function toISOString(dateText, timeText) {
  const target = new Date(`${dateText}T${timeText}:00+08:00`);
  return Number.isNaN(target.getTime()) ? '' : target.toISOString();
}

Page({
  data: {
    headerHeight: 128,
    dateRange: getDateRange(),
    categoryList: PUBLISH_CATEGORIES,
    canPublish: false,
    submitting: false,
    form: {
      category: 'study',
      title: '',
      desc: '',
      location: '',
      date: '',
      time: '',
      deadlineTime: '',
      maxParticipants: '4',
      feeText: '',
      contact: '',
    },
  },

  onLoad() {
    this.updatePublishState();
  },

  onHeaderHeightChange(e) {
    const { height } = e.detail;
    if (height) this.setData({ headerHeight: height });
  },

  onBack() {
    if (getCurrentPages().length > 1) {
      wx.navigateBack({ delta: 1 });
      return;
    }
    wx.redirectTo({ url: '/pages/meetup/index' });
  },

  getEventValue(e) {
    if (e.detail && typeof e.detail === 'object' && 'value' in e.detail) return e.detail.value;
    return e.detail;
  },

  onFieldChange(e) {
    const { field } = e.currentTarget.dataset;
    if (!field) return;

    this.setData({
      [`form.${field}`]: this.getEventValue(e),
    }, () => {
      this.updatePublishState();
    });
  },

  onCategorySelect(e) {
    const { value } = e.currentTarget.dataset;
    if (!value) return;

    this.setData({ 'form.category': value }, () => {
      this.updatePublishState();
    });
  },

  updatePublishState() {
    const { form } = this.data;
    const hasTitle = form.title.trim().length >= 2;
    const hasLocation = form.location.trim().length >= 2;
    const hasDate = form.date && form.time;
    const hasContact = form.contact.trim().length >= 2;
    const hasMembers = Number(form.maxParticipants) >= 2;

    this.setData({
      canPublish: Boolean(hasTitle && hasLocation && hasDate && hasContact && hasMembers),
    });
  },

  async onPublish() {
    const { form, canPublish, submitting } = this.data;
    if (submitting) return;

    if (!canPublish) {
      wx.showToast({ title: '请补充标题、地点、时间、人数和联系方式', icon: 'none' });
      return;
    }

    const startAt = toISOString(form.date, form.time);
    const deadlineAt = form.deadlineTime ? toISOString(form.date, form.deadlineTime) : '';
    if (!startAt) {
      wx.showToast({ title: '开始时间格式有误', icon: 'none' });
      return;
    }

    this.setData({ submitting: true });
    try {
      const res = await publishMeetup({
        category: form.category,
        title: form.title.trim(),
        desc: form.desc.trim(),
        location: form.location.trim(),
        start_at: startAt,
        deadline_at: deadlineAt,
        max_participants: Number(form.maxParticipants),
        fee_text: form.feeText.trim(),
        tags: [],
        contact: form.contact.trim(),
      });
      const data = (res && res.data) || res || {};
      const id = data.id || '';
      wx.showToast({ title: '发布成功', icon: 'success' });
      setTimeout(() => {
        wx.redirectTo({ url: `/pages/meetup/detail/index?id=${id}` });
      }, 360);
    } catch (error) {
      wx.showToast({ title: error.message || '发布失败，请重试', icon: 'none' });
    } finally {
      this.setData({ submitting: false });
    }
  },
});
