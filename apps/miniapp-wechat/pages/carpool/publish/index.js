import { getMenuButtonSafeArea } from '../../../utils/navigation';
import { publishCarpool } from '../../../api/modules/carpool';
import { saveHistoryAddress } from '../../../utils/addressStore';

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

function inferCategory(dateText) {
  if (!dateText) return 'today';

  const [yearText, monthText, dayText] = dateText.split('-');
  const targetDate = new Date(Number(yearText), Number(monthText) - 1, Number(dayText));
  const now = new Date();
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
  const tomorrow = new Date(now.getFullYear(), now.getMonth(), now.getDate() + 1);

  if (targetDate.getTime() === today.getTime() || targetDate.getTime() === tomorrow.getTime()) {
    return 'today';
  }
  return 'return';
}

Page({
  data: {
    canPublish: false,
    menuTop: 40,
    menuButtonHeight: 32,
    navHeight: 88,
    dateRange: getDateRange(),
    form: {
      from: '',
      to: '',
      date: '',
      time: '',
      seats: '',
      price: '',
      contact: '',
      note: '',
    },
  },

  onLoad(options = {}) {
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
    wx.redirectTo({ url: '/pages/carpool/index' });
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

  onStartChange(e) {
    this.setData({ 'form.from': e.detail.value }, () => {
      this.updatePublishState();
    });
  },

  onEndChange(e) {
    this.setData({ 'form.to': e.detail.value }, () => {
      this.updatePublishState();
    });
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
    const hasRoute = form.from.trim() && form.to.trim();
    const hasTime = form.date && form.time;
    const hasSeats = form.seats.trim();
    const hasPrice = form.price.trim();
    const hasContact = form.contact.trim().length >= 2;

    this.setData({
      canPublish: Boolean(hasRoute && hasTime && hasSeats && hasPrice && hasContact),
    });
  },

  async release() {
    const { form } = this.data;
    const missing = [];
    if (!form.from.trim()) missing.push('出发地');
    if (!form.to.trim()) missing.push('目的地');
    if (!form.date) missing.push('日期');
    if (!form.time) missing.push('时间');
    if (!form.seats.trim()) missing.push('人数/余座');
    if (!form.price.trim()) missing.push('费用说明');
    if (form.contact.trim().length < 2) missing.push('联系方式');

    if (missing.length) {
      wx.showToast({ title: `请补充${missing.join('、')}`, icon: 'none' });
      return;
    }

    const category = inferCategory(form.date);
    const typeMap = { today: '今日顺路', return: '返校专线', leave: '出校拼车', longterm: '长期通勤' };
    const tagsMap = {
      longterm: ['固定路线', '支持长期沟通'],
      return: ['返校优先', '刚刚发布'],
      leave: ['周末可约', '刚刚发布'],
      today: ['今日可拼', '刚刚发布'],
    };

    const payload = {
      category,
      from: form.from.trim(),
      to: form.to.trim(),
      time: `${formatDisplayDate(form.date)} ${form.time}`,
      type: typeMap[category] || '今日顺路',
      seats_text: form.seats.trim(),
      price: form.price.trim(),
      note: form.note.trim() || `可通过 ${form.contact.trim()} 联系我，细节可继续沟通。`,
      tags: tagsMap[category] || ['今日可拼', '刚刚发布'],
      contact: form.contact.trim(),
    };

    try {
      await publishCarpool(payload);
      saveHistoryAddress(form.from.trim());
      saveHistoryAddress(form.to.trim());
      wx.showToast({ title: '发布成功', icon: 'success' });
      setTimeout(() => {
        if (getCurrentPages().length > 1) {
          wx.navigateBack({ delta: 1 });
          return;
        }
        wx.redirectTo({ url: '/pages/carpool/index' });
      }, 360);
    } catch (e) {
      wx.showToast({ title: '发布失败，请重试', icon: 'none' });
    }
  },
});
