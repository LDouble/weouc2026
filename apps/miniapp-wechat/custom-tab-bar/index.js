const app = getApp();

Component({
  data: {
    value: '',
    unreadNum: 0,
    tabs: [
      {
        icon: 'home',
        value: 'home',
        label: '发现',
        url: '/pages/home/index',
      },
      {
        icon: 'chat',
        value: 'message',
        label: '消息',
        url: '/pages/message/index',
      },
      {
        icon: 'user',
        value: 'my',
        label: '我的',
        url: '/pages/my/index',
      },
    ],
  },

  lifetimes: {
    ready() {
      const pages = getCurrentPages();
      const curPage = pages[pages.length - 1];
      if (curPage) {
        const nameRe = /pages\/([^/]+)\/index/.exec(curPage.route);
        if (nameRe && nameRe[1]) {
          this.setData({ value: nameRe[1] });
        }
      }

      this.setUnreadNum(app.globalData.unreadNum || 0);
      app.eventBus.on('unread-num-change', (unreadNum) => {
        this.setUnreadNum(unreadNum);
      });
    },
  },

  methods: {
    handleTap(e) {
      const { value, url } = e.currentTarget.dataset;
      if (!url || value === this.data.value) return;
      wx.switchTab({ url });
    },

    setUnreadNum(unreadNum) {
      this.setData({ unreadNum });
    },
  },
});
