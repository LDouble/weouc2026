import useToastBehavior from '~/behaviors/useToast';
import { getStudentProfile, createStudentProfile, updateStudentProfile } from '~/api/modules/student';
import { post } from '~/api/request';

Page({
  behaviors: [useToastBehavior],

  data: {
    isBound: false,
    headerHeight: 0,
    sid: '',
    password: '',
    captcha: '',
    countdown: 0,
    bindInfo: {
      sid: '',
      name: '',
      bindTime: '',
    },
    unbindDialogVisible: false,
    isSubmitDisabled: true,
  },

  onLoad() {
    this.getBindStatus();
  },

  onBack() {
    wx.navigateBack();
  },

  onHeaderHeightChange(e) {
    this.setData({
      headerHeight: e.detail.height,
    });
  },

  async getBindStatus() {
    try {
      const res = await getStudentProfile();
      const profile = res.data || res;
      if (profile.is_bound) {
        this.setData({
          isBound: true,
          bindInfo: {
            sid: profile.student_id || '',
            name: profile.name || '',
            bindTime: profile.updated_at || '',
          },
        });
      }
    } catch (e) {
      console.error(e);
    }
  },

  updateSubmitState() {
    const { sid, password, captcha } = this.data;
    this.setData({
      isSubmitDisabled: !sid || !password || !captcha,
    });
  },

  onSidChange(e) {
    this.setData({ sid: e.detail.value }, () => this.updateSubmitState());
  },

  onPasswordChange(e) {
    this.setData({ password: e.detail.value }, () => this.updateSubmitState());
  },

  onCaptchaChange(e) {
    this.setData({ captcha: e.detail.value }, () => this.updateSubmitState());
  },

  onSendCaptcha() {
    if (this.data.countdown > 0) return;
    if (!this.data.sid) {
      this.onShowToast('#t-toast', '请先输入学号');
      return;
    }
    post('/edu/send-captcha', { sid: this.data.sid })
      .then(() => {
        this.onShowToast('#t-toast', '验证码已发送');
        this.startCountdown();
      })
      .catch(() => {
        this.onShowToast('#t-toast', '验证码发送失败');
      });
  },

  startCountdown() {
    this.setData({ countdown: 60 });
    this._timer = setInterval(() => {
      if (this.data.countdown <= 1) {
        clearInterval(this._timer);
        this.setData({ countdown: 0 });
      } else {
        this.setData({ countdown: this.data.countdown - 1 });
      }
    }, 1000);
  },

  async onBind() {
    const { sid, password, captcha, isSubmitDisabled } = this.data;
    if (isSubmitDisabled) return;

    try {
      const data = { student_id: sid, password, captcha };
      const res = await createStudentProfile(data);
      const profile = res.data || res;
      this.setData({
        isBound: true,
        bindInfo: {
          sid: profile.student_id || sid,
          name: profile.name || '',
          bindTime: profile.updated_at || '',
        },
      });
      this.onShowToast('#t-toast', '绑定成功');
    } catch (e) {
      this.onShowToast('#t-toast', '绑定失败，请检查信息');
    }
  },

  async onUnbind() {
    this.setData({
      unbindDialogVisible: true,
    });
  },

  async onUnbindConfirm() {
    try {
      await updateStudentProfile({ is_bound: false });
      this.setData({
        isBound: false,
        sid: '',
        password: '',
        captcha: '',
        bindInfo: { sid: '', name: '', bindTime: '' },
        unbindDialogVisible: false,
        isSubmitDisabled: true,
      });
      this.onShowToast('#t-toast', '已解绑');
    } catch (e) {
      this.onShowToast('#t-toast', '解绑失败');
    }
  },

  onUnbindCancel() {
    this.setData({
      unbindDialogVisible: false,
    });
  },
});
