import { getCourses, getGradeDetails } from '~/api/modules/student'

Page({
  data: {
    loading: true,
    placeholder: true,
    totalSituationDataList: null,
    totalSituationKeyList: null,
    completeRateDataList: null,
    complete_rate_keyList: null,
    interactionSituationDataList: null,
    interaction_situation_keyList: null,
    areaDataList: null,
    areaDataKeysList: null,
    memberitemWidth: null,
    smallitemWidth: null,
  },

  onLoad() {
    this.init();
  },

  init() {
    this.loadData();
  },

  loadData() {
    Promise.all([getCourses(), getGradeDetails()])
      .then(([coursesRes, gradesRes]) => {
        this.setData({ loading: false, placeholder: false });
      })
      .catch(() => {
        this.setData({ loading: false, placeholder: true });
      });
  },
});
