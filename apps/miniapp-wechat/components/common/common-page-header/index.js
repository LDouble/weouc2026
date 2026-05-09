import { getMenuButtonSafeArea } from '../../../utils/navigation';

Component({
  options: {
    multipleSlots: true,
    styleIsolation: 'shared',
  },

  externalClasses: ['custom-class', 'title-class', 'subtitle-class'],

  properties: {
    title: {
      type: String,
      value: '',
    },
    subtitle: {
      type: String,
      value: '',
    },
    useTitleSlot: {
      type: Boolean,
      value: false,
    },
    showBack: {
      type: Boolean,
      value: true,
    },
    backIcon: {
      type: String,
      value: 'arrow-left',
    },
    rightIcon: {
      type: String,
      value: '',
    },
    rightText: {
      type: String,
      value: '',
    },
    iconSize: {
      type: String,
      value: '32rpx',
    },
    background: {
      type: String,
      value: 'rgba(255, 255, 255, 0.86)',
    },
    borderColor: {
      type: String,
      value: 'rgba(241, 245, 249, 0.9)',
    },
    buttonBackground: {
      type: String,
      value: 'rgba(241, 245, 249, 0.96)',
    },
    buttonColor: {
      type: String,
      value: '#314158',
    },
    titleColor: {
      type: String,
      value: '#1d293d',
    },
    subtitleColor: {
      type: String,
      value: '#52627a',
    },
  },

  data: {
    menuTop: 40,
    menuSafeRight: 16,
    menuButtonHeight: 32,
  },

  lifetimes: {
    attached() {
      this.applySafeArea();
    },

    ready() {
      this.measureHeight();
    },
  },

  methods: {
    applySafeArea() {
      const { right, top, height } = getMenuButtonSafeArea(12);
      this.setData(
        {
          menuTop: top,
          menuSafeRight: right,
          menuButtonHeight: height,
        },
        () => this.measureHeight(),
      );
    },

    measureHeight() {
      wx.nextTick(() => {
        this.createSelectorQuery()
          .in(this)
          .select('.common-page-header')
          .boundingClientRect((rect) => {
            if (!rect || !rect.height) return;
            this.triggerEvent('heightchange', { height: rect.height });
          })
          .exec();
      });
    },

    onBack() {
      this.triggerEvent('back');
    },

    onRightTap() {
      this.triggerEvent('righttap');
    },
  },
});
