import { loadNotifications, loadUnreadNotificationCount, publishUnreadCount, readNotification } from '../../services/notificationService';
import { loadPortalNoticeDetail, loadPortalNotices } from '../../services/portalService';


function navigateInternal(url) {
  if (!url) return;
  if (['/pages/home/index', '/pages/message/index', '/pages/my/index'].includes(url)) {
    wx.switchTab({ url });
    return;
  }
  wx.navigateTo({ url });
}

Page({
  data: {
    activeTab: 'notifications',
    messageList: [],
    noticeList: [],
    loading: false,
    noticeLoading: false,
    emptyText: '暂无通知',
    page: 1,
    pageSize: 20,
    hasMore: true,
    unreadOnly: false,
    noticeDetail: null,
  },

  onLoad(options = {}) {
    if (options.noticeId) {
      this.setData({ activeTab: 'notices' });
      this.openNoticeById(options.noticeId);
    }
    this.getMessageList(true);
    this.getNoticeList();
  },

  onShow() {
    this.getMessageList(true, { silent: true });
  },

  onPullDownRefresh() {
    Promise.all([this.getMessageList(true, { silent: true }), this.getNoticeList({ silent: true })])
      .finally(() => wx.stopPullDownRefresh());
  },

  onReachBottom() {
    if (this.data.activeTab === 'notifications') this.loadMoreMessages();
  },

  switchTab(e) {
    const { tab } = e.currentTarget.dataset;
    if (!tab || tab === this.data.activeTab) return;
    this.setData({ activeTab: tab, noticeDetail: null });
  },

  toggleUnreadOnly() {
    this.setData({ unreadOnly: !this.data.unreadOnly }, () => this.getMessageList(true));
  },

  async getMessageList(reset = true, options = {}) {
    if (this.data.loading) return null;
    const page = reset ? 1 : this.data.page + 1;
    if (!reset && !this.data.hasMore) return null;

    if (!options.silent) wx.showNavigationBarLoading();
    this.setData({ loading: true });

    try {
      const result = await loadNotifications({
        page,
        pageSize: this.data.pageSize,
        unreadOnly: this.data.unreadOnly,
      });
      const nextList = reset ? result.items : this.data.messageList.concat(result.items);
      loadUnreadNotificationCount().then(publishUnreadCount).catch(() => {
        publishUnreadCount(nextList.filter((item) => !item.read).length);
      });
      this.setData({
        messageList: nextList,
        page,
        hasMore: page * this.data.pageSize < result.total,
        emptyText: this.data.unreadOnly ? '没有未读通知' : '暂无通知',
      });
      return result;
    } catch (error) {
      if (!options.silent) wx.showToast({ title: error.message || '通知加载失败', icon: 'none' });
      return null;
    } finally {
      this.setData({ loading: false });
      if (!options.silent) wx.hideNavigationBarLoading();
    }
  },

  loadMoreMessages() {
    this.getMessageList(false);
  },

  async getNoticeList(options = {}) {
    this.setData({ noticeLoading: true });
    try {
      const result = await loadPortalNotices({ page: 1, pageSize: 20 });
      this.setData({ noticeList: result.items });
      return result;
    } catch (error) {
      if (!options.silent) wx.showToast({ title: error.message || '公告加载失败', icon: 'none' });
      return null;
    } finally {
      this.setData({ noticeLoading: false });
    }
  },

  async openNoticeById(noticeId) {
    if (!noticeId) return;
    wx.showLoading({ title: '加载公告' });
    try {
      const detail = await loadPortalNoticeDetail(noticeId);
      this.setData({ activeTab: 'notices', noticeDetail: detail });
    } catch (error) {
      wx.showToast({ title: error.message || '公告不存在', icon: 'none' });
    } finally {
      wx.hideLoading();
    }
  },

  openNotice(e) {
    const { id } = e.currentTarget.dataset;
    this.openNoticeById(id);
  },

  closeNotice() {
    this.setData({ noticeDetail: null });
  },

  async openMessage(e) {
    const { id } = e.currentTarget.dataset;
    const current = this.data.messageList.find((item) => item.id === id);
    if (!current) return;

    if (!current.read) {
      const nextList = this.data.messageList.map((item) => (
        item.id === id ? { ...item, read: true } : item
      ));
      this.setData({ messageList: nextList });
      publishUnreadCount(nextList.filter((item) => !item.read).length);
      readNotification(id)
        .then(() => loadUnreadNotificationCount().then(publishUnreadCount))
        .catch(() => this.getMessageList(true, { silent: true }));
    }

    wx.showModal({
      title: current.title,
      content: current.content || '暂无内容',
      confirmText: current.actionUrl ? '去查看' : '知道了',
      cancelText: '关闭',
      showCancel: Boolean(current.actionUrl),
      success: (res) => {
        if (res.confirm && current.actionUrl) navigateInternal(current.actionUrl);
      },
    });
  },
});
