Component({
  options: {
    styleIsolation: 'shared',
  },

  externalClasses: ['custom-class'],

  properties: {
    icon: {
      type: String,
      value: 'add',
    },
    iconSize: {
      type: String,
      value: '40rpx',
    },
    text: {
      type: String,
      value: '',
    },
    right: {
      type: Number,
      value: 32,
    },
    bottom: {
      type: Number,
      value: 44,
    },
    background: {
      type: String,
      value: 'linear-gradient(135deg, #615fff, #4f46e5)',
    },
    shadow: {
      type: String,
      value: '0 16rpx 30rpx rgba(97, 95, 255, 0.32)',
    },
  },

  methods: {
    onTap() {
      this.triggerEvent('tap');
    },
  },
});
