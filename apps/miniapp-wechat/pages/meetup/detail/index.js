import { cancelMeetupJoin, cancelMeetupPublish, joinMeetup } from '../../../api/modules/meetup';
import { loadMeetupDetail } from '../../../services/meetupService';

Page({
  data: {
    headerHeight: 128,
    loading: true,
    detail: {},
  },

  onLoad(options = {}) {
    this.loadDetail(options.id);
  },

  async loadDetail(id) {
    if (!id) {
      this.setData({ loading: false });
      return;
    }

    this.setData({ loading: true });
    try {
      const detail = await loadMeetupDetail(id);
      this.setData({ detail, loading: false });
    } catch (error) {
      this.setData({ loading: false });
      wx.showToast({ title: error.message || '加载失败，请重试', icon: 'none' });
    }
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

  onContact() {
    const { contact, canViewContact } = this.data.detail;
    if (canViewContact === false) {
      wx.showModal({
        title: '无法查看联系方式',
        content: '绑定教务后即可查看发起人的联系方式，是否前往绑定？',
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
      wx.showToast({ title: '当前暂无可复制的联系方式', icon: 'none' });
      return;
    }
    wx.setClipboardData({
      data: contact,
      success: () => wx.showToast({ title: '联系方式已复制', icon: 'success' }),
    });
  },

  async onJoinTap() {
    const { detail } = this.data;
    if (!detail.id || !detail.canJoin) return;

    try {
      await joinMeetup(detail.id);
      await this.loadDetail(detail.id);
      wx.showToast({ title: '报名成功', icon: 'success' });
    } catch (error) {
      wx.showToast({ title: error.message || '报名失败，请重试', icon: 'none' });
    }
  },

  async onCancelJoinTap() {
    const { detail } = this.data;
    if (!detail.id || !detail.canCancelJoin) return;

    try {
      await cancelMeetupJoin(detail.id);
      await this.loadDetail(detail.id);
      wx.showToast({ title: '已取消报名', icon: 'success' });
    } catch (error) {
      wx.showToast({ title: error.message || '取消失败，请重试', icon: 'none' });
    }
  },

  onCancelPublishTap() {
    const { detail } = this.data;
    if (!detail.id || !detail.canCancelPublish) return;

    wx.showModal({
      title: '取消组局',
      content: '取消后公开用户将无法再看到这场组局，确认继续？',
      confirmText: '确认取消',
      confirmColor: '#dc2626',
      success: async (res) => {
        if (!res.confirm) return;
        try {
          await cancelMeetupPublish(detail.id);
          await this.loadDetail(detail.id);
          wx.showToast({ title: '已取消组局', icon: 'success' });
        } catch (error) {
          wx.showToast({ title: error.message || '取消失败，请重试', icon: 'none' });
        }
      },
    });
  },

  onShareAppMessage() {
    const { detail } = this.data;
    return {
      title: detail.title || '校园组局',
      path: `/pages/meetup/detail/index?id=${detail.id || ''}`,
    };
  },
});
