import { initAppConfigStore } from './store/config';
import createBus from './utils/eventBus';
import { silentLogin } from './api/auth';

App({
  onLaunch() {
    initAppConfigStore();
    silentLogin();

    const updateManager = wx.getUpdateManager();

    updateManager.onCheckForUpdate(() => {});

    updateManager.onUpdateReady(() => {
      wx.showModal({
        title: '更新提示',
        content: '新版本已经准备好，是否重启应用？',
        success(res) {
          if (res.confirm) {
            updateManager.applyUpdate();
          }
        },
      });
    });
  },

  globalData: {
    userInfo: null,
  },

  eventBus: createBus(),
});
