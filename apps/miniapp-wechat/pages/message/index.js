import {
  loadNotificationList,
  readNotification,
} from '../../services/notificationService';

Page({
  data: {
    messageList: [],
    loading: false,
    loadingMore: false,
    hasMore: true,
    page: 1,
    pageSize: 20,
    total: 0,
    emptyText: '暂无消息',
  },

  onLoad() {},

  onShow() {
    this.refreshMessages();
  },

  async refreshMessages() {
    if (this.data.loading) return;
    this.setData({ loading: true });

    try {
      const result = await loadNotificationList({
        page: 1,
        pageSize: this.data.pageSize,
      });
      this.setData({
        messageList: result.list,
        loading: false,
        loadingMore: false,
        page: 1,
        total: result.total,
        hasMore: this.data.pageSize < result.total,
        emptyText: '暂无消息',
      });
    } catch (error) {
      this.setData({
        loading: false,
        loadingMore: false,
        emptyText: '消息加载失败，请稍后重试',
      });
      wx.showToast({ title: error.message || '加载失败，请重试', icon: 'none' });
    }
  },

  async onLoadMore() {
    const {
      loading,
      loadingMore,
      hasMore,
      page,
      pageSize,
      messageList,
    } = this.data;
    if (loading || loadingMore || !hasMore) return;

    const nextPage = page + 1;
    this.setData({ loadingMore: true });

    try {
      const result = await loadNotificationList({
        page: nextPage,
        pageSize,
      });
      const nextList = messageList.concat(result.list);
      this.setData({
        messageList: nextList,
        loadingMore: false,
        page: nextPage,
        total: result.total,
        hasMore: nextPage * pageSize < result.total,
      });
    } catch (error) {
      this.setData({ loadingMore: false });
      wx.showToast({ title: error.message || '加载失败，请重试', icon: 'none' });
    }
  },

  async markAsRead(messageId) {
    if (!messageId) return;
    try {
      await readNotification(messageId);
      const nextList = this.data.messageList.map((item) => {
        if (item.id !== messageId) return item;
        return { ...item, read: true };
      });
      this.setData({ messageList: nextList });
    } catch (_) {}
  },

  onOpenMessage(event) {
    const { item } = event.currentTarget.dataset;
    if (!item || !item.id) return;

    if (!item.read) {
      this.markAsRead(item.id);
    }

    if (!item.actionUrl || !item.actionUrl.startsWith('/pages/')) return;

    if (item.actionUrl === '/pages/home/index' || item.actionUrl === '/pages/my/index') {
      wx.switchTab({ url: item.actionUrl });
      return;
    }

    wx.navigateTo({ url: item.actionUrl });
  },
});
