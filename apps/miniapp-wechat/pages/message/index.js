Page({
  data: {
    messageList: [],
    loading: false,
    emptyText: '功能开发中',
  },

  onLoad() {},

  onShow() {},

  getMessageList() {
    this.setData({ loading: false });
  },

  toChat(event) {
    const { userId } = event.currentTarget.dataset;
    wx.navigateTo({ url: `/pages/chat/index?userId=${userId}` });
  },
});
