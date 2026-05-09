import { loginWithWechat } from '~/services/sessionService';

Page({
  data: {
    isLoading: false,
    isCheck: false,
  },

  onBack() {
    if (getCurrentPages().length > 1) {
      wx.navigateBack({ delta: 1 });
      return;
    }

    wx.switchTab({ url: '/pages/home/index' });
  },

  onCheckChange() {
    this.setData({
      isCheck: !this.data.isCheck,
    });
  },

  async onWxLogin() {
    if (!this.data.isCheck) {
      wx.showToast({ title: '请先同意用户协议', icon: 'none' });
      return;
    }

    if (this.data.isLoading) return;
    this.setData({ isLoading: true });

    try {
      const result = await loginWithWechat();
      if (result && result.session && result.session.token) {
        wx.switchTab({ url: '/pages/home/index' });
      }
    } catch (err) {
      wx.showToast({ title: '登录失败，请重试', icon: 'none' });
    } finally {
      this.setData({ isLoading: false });
    }
  },
});
