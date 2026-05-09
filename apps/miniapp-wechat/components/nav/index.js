Component({
  options: {
    styleIsolation: 'shared',
  },
  properties: {
    titleText: String,
  },
  data: {
    visible: false,
    sidebar: [
      {
        title: '首页',
        url: 'pages/home/index',
        isSidebar: true,
      },
      {
        title: '二手交易',
        url: 'pages/market/index',
        isSidebar: false,
      },
      {
        title: '跑腿帮办',
        url: 'pages/errand/index',
        isSidebar: false,
      },
      {
        title: '失物招领',
        url: 'pages/lost-found/index',
        isSidebar: false,
      },
      {
        title: '资料汇总',
        url: 'pages/resource/index',
        isSidebar: false,
      },
      {
        title: '发布信息',
        url: 'pages/release/index',
        isSidebar: false,
      },
      {
        title: '搜索',
        url: 'pages/search/index',
        isSidebar: false,
      },
      {
        title: '消息',
        url: 'pages/message/index',
        isSidebar: true,
      },
      {
        title: '我的',
        url: 'pages/my/index',
        isSidebar: true,
      },
      {
        title: '设置',
        url: 'pages/setting/index',
        isSidebar: false,
      },
    ],
    statusHeight: 0,
  },
  lifetimes: {
    ready() {
      const statusHeight = wx.getWindowInfo().statusBarHeight;
      this.setData({ statusHeight });
    },
  },
  methods: {
    openDrawer() {
      this.setData({
        visible: true,
      });
    },
    itemClick(e) {
      const that = this;
      const { isSidebar, url } = e.detail.item;
      if (isSidebar) {
        wx.switchTab({
          url: `/${url}`,
        }).then(() => {
          that.setData({
            visible: false,
          });
        });
      } else {
        wx.navigateTo({
          url: `/${url}`,
        }).then(() => {
          that.setData({
            visible: false,
          });
        });
      }
    },
  },
});
