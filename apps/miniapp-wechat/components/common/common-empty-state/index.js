Component({
  options: {
    styleIsolation: 'shared',
  },

  externalClasses: ['custom-class'],

  properties: {
    icon: {
      type: String,
      value: 'search',
    },
    iconSize: {
      type: String,
      value: '48rpx',
    },
    title: {
      type: String,
      value: '暂无数据',
    },
    description: {
      type: String,
      value: '',
    },
    actionText: {
      type: String,
      value: '',
    },
    minHeight: {
      type: Number,
      value: 620,
    },
    iconColor: {
      type: String,
      value: '#615fff',
    },
    iconBackground: {
      type: String,
      value: '#eef2ff',
    },
    actionBackground: {
      type: String,
      value: 'linear-gradient(135deg, #615fff, #4f46e5)',
    },
  },

  methods: {
    onAction() {
      this.triggerEvent('action');
    },
  },
});
