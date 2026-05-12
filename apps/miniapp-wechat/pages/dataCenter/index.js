import { loadAcademicDashboardModel } from '../../services/academicService';

Page({
  data: {
    loading: true,
    placeholder: true,
    totalSituationDataList: [],
    completeRateDataList: [],
    interactionSituationDataList: [],
  },

  onLoad() {
    this.init();
  },

  init() {
    this.loadData();
  },

  async loadData() {
    this.setData({ loading: true });

    try {
      const model = await loadAcademicDashboardModel();
      this.setData({
        loading: false,
        placeholder: false,
        totalSituationDataList: model.totalSituationDataList,
        completeRateDataList: model.completeRateDataList,
        interactionSituationDataList: model.interactionSituationDataList,
      });
    } catch (error) {
      this.setData({
        loading: false,
        placeholder: true,
      });
      wx.showToast({
        title: error.message || '教务数据读取失败',
        icon: 'none',
      });
    }
  },
});
