Page({
  data: {
    emptyText: '功能开发中',
    userId: null,
    name: '',
    messages: [],
    input: '',
    anchor: '',
    keyboardHeight: 0,
  },

  onLoad(options) {
    if (options && options.userId) {
      this.setData({ userId: options.userId });
    }
  },

  onReady() {},

  onShow() {},

  onHide() {},

  onUnload() {},
});
