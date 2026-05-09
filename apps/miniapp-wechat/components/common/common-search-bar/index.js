Component({
  options: {
    styleIsolation: 'shared',
  },

  externalClasses: ['custom-class', 'input-class'],

  properties: {
    value: {
      type: String,
      value: '',
    },
    placeholder: {
      type: String,
      value: '搜索',
    },
    disabled: {
      type: Boolean,
      value: false,
    },
    clearable: {
      type: Boolean,
      value: true,
    },
    confirmType: {
      type: String,
      value: 'search',
    },
    icon: {
      type: String,
      value: 'search',
    },
    iconSize: {
      type: String,
      value: '28rpx',
    },
    clearIcon: {
      type: String,
      value: 'close-circle-filled',
    },
    clearIconSize: {
      type: String,
      value: '28rpx',
    },
    variant: {
      type: String,
      value: 'default',
    },
    background: {
      type: String,
      value: 'rgba(241, 245, 249, 0.96)',
    },
    color: {
      type: String,
      value: '#1d293d',
    },
    iconColor: {
      type: String,
      value: '#90a1b9',
    },
  },

  methods: {
    onInput(e) {
      this.triggerEvent('input', { value: e.detail.value || '' });
    },

    onConfirm(e) {
      this.triggerEvent('confirm', { value: e.detail.value || this.data.value || '' });
    },

    onClear() {
      this.triggerEvent('clear', { value: '' });
    },

    onTap() {
      this.triggerEvent('tap');
    },
  },
});
