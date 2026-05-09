import { wxLogin, isLoggedIn } from '~/api/auth';

Page({
  data: {
    isLoading: false,
    isCheck: false,
    canLogin: false,
  },

  onCheckChange(e) {
    const { value } = e.detail;
    this.setData({
      isCheck: value === 'agree',
      canLogin: value === 'agree',
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
      const result = await wxLogin();
      if (result && result.token) {
        wx.switchTab({ url: '/pages/home/index' });
      }
    } catch (err) {
      wx.showToast({ title: '登录失败，请重试', icon: 'none' });
    } finally {
      this.setData({ isLoading: false });
    }
  },
});
