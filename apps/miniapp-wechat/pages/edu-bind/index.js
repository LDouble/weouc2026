import useToastBehavior from '~/behaviors/useToast';
import {
  bindAcademicAccount,
  loadAcademicBindingModel,
  sendAcademicCaptcha,
  unbindAcademicAccount,
} from '~/services/profileService';

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
    this.loadBindingState();
  },

  onShow() {
    this.loadBindingState();
  },

  onUnload() {
    this.clearCountdownTimer();
  },

  onBack() {
    wx.navigateBack();
  },

  onHeaderHeightChange(e) {
    this.setData({
      headerHeight: e.detail.height,
    });
  },

  async loadBindingState() {
    try {
      const model = await loadAcademicBindingModel();
      this.setData({
        isBound: model.isBound,
        bindInfo: model.bindInfo,
      });
    } catch (e) {
      this.onShowToast('#t-toast', '教务状态加载失败');
    }
  },

  updateSubmitState() {
    const { sid, password, captcha } = this.data;
    this.setData({
      isSubmitDisabled: !sid.trim() || !password.trim() || !captcha.trim(),
    });
  },

  onSidChange(e) {
    this.setData({ sid: e.detail.value || '' });
    this.updateSubmitState();
  },

  onPasswordChange(e) {
    this.setData({ password: e.detail.value || '' });
    this.updateSubmitState();
  },

  onCaptchaChange(e) {
    this.setData({ captcha: e.detail.value || '' });
    this.updateSubmitState();
  },

  onSendCaptcha() {
    if (this.data.countdown > 0) return;
    if (!this.data.sid.trim()) {
      this.onShowToast('#t-toast', '请先输入学号');
      return;
    }
    sendAcademicCaptcha(this.data.sid.trim())
      .then(() => {
        this.onShowToast('#t-toast', '验证码已发送');
        this.startCountdown();
      })
      .catch(() => {
        this.onShowToast('#t-toast', '验证码发送失败');
      });
  },

  startCountdown() {
    this.clearCountdownTimer();
    this.setData({ countdown: 60 });
    this._timer = setInterval(() => {
      if (this.data.countdown <= 1) {
        this.clearCountdownTimer();
        this.setData({ countdown: 0 });
      } else {
        this.setData({ countdown: this.data.countdown - 1 });
      }
    }, 1000);
  },

  clearCountdownTimer() {
    if (this._timer) {
      clearInterval(this._timer);
      this._timer = null;
    }
  },

  async onBind() {
    const { sid, password, captcha, isSubmitDisabled } = this.data;
    if (isSubmitDisabled) return;

    try {
      const profile = await bindAcademicAccount({
        studentId: sid.trim(),
        password: password.trim(),
        captcha: captcha.trim(),
      });
      this.setData({
        isBound: true,
        bindInfo: {
          sid: profile.studentId || sid.trim(),
          name: profile.name || '',
          bindTime: profile.bindingUpdatedLabel || '',
        },
        password: '',
        captcha: '',
        isSubmitDisabled: true,
      });
      this.clearCountdownTimer();
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
      await unbindAcademicAccount();
      this.setData({
        isBound: false,
        sid: '',
        password: '',
        captcha: '',
        bindInfo: { sid: '', name: '', bindTime: '' },
        unbindDialogVisible: false,
        isSubmitDisabled: true,
      });
      this.clearCountdownTimer();
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
