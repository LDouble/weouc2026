import useToastBehavior from '~/behaviors/useToast';
import { getStudentProfile } from '~/api/modules/student';
import { fetchFeedList } from '~/api/modules/feed';
import { fetchErrandList } from '~/api/modules/errand';

Page({
  behaviors: [useToastBehavior],

  data: {
    isLoad: false,
    headerHeight: 0,
    personalInfo: {},
    stats: [
      { key: 'published', icon: 'chart-bar', value: 0, label: '发布', bg: 'indigo' },
      { key: 'accepted', icon: 'check-circle', value: 0, label: '接单', bg: 'purple' },
      { key: 'fav', icon: 'star', value: 0, label: '收藏', bg: 'amber' },
    ],
    serviceList: [
      { name: '我的发布', icon: 'root-list', type: 'publish', url: '/pages/my/publish/index' },
      { name: '我的接单', icon: 'check-circle', type: 'accepted', url: '/pages/my/accepted/index' },
      { name: '星标收藏', icon: 'star', type: 'fav', url: '' },
    ],
    settingList: [
      { name: '教务绑定', icon: 'education', type: 'eduBind', url: '/pages/edu-bind/index' },
      { name: '消息中心', icon: 'notification', type: 'message', url: '', badge: 2 },
      { name: '联系客服', icon: 'service', type: 'service', url: '' },
      { name: '设置', icon: 'setting', type: 'setting', url: '/pages/setting/index' },
    ],
  },

  onLoad() {},

  onBack() {
    wx.navigateBack();
  },

  onHeaderHeightChange(e) {
    this.setData({
      headerHeight: e.detail.height,
    });
  },

  async onShow() {
    const Token = wx.getStorageSync('access_token');
    if (!Token) {
      this.setData({
        isLoad: false,
        personalInfo: {},
        stats: this.data.stats.map((item) => ({ ...item, value: 0 })),
      });
      return;
    }

    try {
      const res = await getStudentProfile();
      const profile = res.data || res;
      const personalInfo = {
        name: profile.name || '',
        image: profile.avatar_url || '',
        major: profile.major || '',
        sid: profile.student_id || '',
        college: profile.college || '',
        grade: profile.grade || '',
        isBound: profile.is_bound || false,
      };

      this.setData({
        isLoad: true,
        personalInfo,
      });

      this.loadStats();
    } catch (e) {
      console.error(e);
    }
  },

  async loadStats() {
    try {
      const [feedRes, errandRes] = await Promise.all([
        fetchFeedList({ page: 1, pageSize: 1, user_role: 'publisher' }),
        fetchErrandList({ page: 1, pageSize: 1, user_role: 'acceptor' }),
      ]);
      const feedData = feedRes.data || feedRes;
      const errandData = errandRes.data || errandRes;
      this.setData({
        'stats[0].value': feedData.total || 0,
        'stats[1].value': errandData.total || 0,
      });
    } catch (e) {
      console.error(e);
    }
  },

  onLogin(e) {
    wx.navigateTo({
      url: '/pages/login/login',
    });
  },

  onNavigateTo() {
    wx.navigateTo({ url: `/pages/my/info-edit/index` });
  },

  onEleClick(e) {
    const { name, url } = e.currentTarget.dataset.data;
    if (url) {
      wx.navigateTo({ url });
      return;
    }
    this.onShowToast('#t-toast', name);
  },
});
