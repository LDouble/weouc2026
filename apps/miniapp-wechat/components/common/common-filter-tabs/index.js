Component({
  options: {
    styleIsolation: 'shared',
  },

  externalClasses: ['custom-class'],

  properties: {
    items: {
      type: Array,
      value: [],
    },
    active: {
      type: String,
      value: '',
    },
    valueKey: {
      type: String,
      value: 'value',
    },
    labelKey: {
      type: String,
      value: 'label',
    },
    variant: {
      type: String,
      value: 'pill',
    },
    scrollable: {
      type: Boolean,
      value: true,
    },
    justify: {
      type: String,
      value: 'start',
    },
    itemMinWidth: {
      type: Number,
      value: 0,
    },
    activeColor: {
      type: String,
      value: '#615fff',
    },
    activeTextColor: {
      type: String,
      value: '#ffffff',
    },
    itemBackground: {
      type: String,
      value: 'rgba(241, 245, 249, 0.96)',
    },
    textColor: {
      type: String,
      value: '#52627a',
    },
  },

  data: {
    normalizedItems: [],
  },

  observers: {
    'items, valueKey, labelKey': function normalizeItems() {
      const { items, valueKey, labelKey } = this.data;
      const normalizedItems = (items || []).map((item, index) => ({
        value: item ? item[valueKey] : '',
        label: item ? item[labelKey] : '',
        raw: item,
        index,
      }));
      this.setData({ normalizedItems });
    },
  },

  lifetimes: {
    attached() {
      const { items, valueKey, labelKey } = this.data;
      const normalizedItems = (items || []).map((item, index) => ({
        value: item ? item[valueKey] : '',
        label: item ? item[labelKey] : '',
        raw: item,
        index,
      }));
      this.setData({ normalizedItems });
    },
  },

  methods: {
    onTap(e) {
      const { value, index } = e.currentTarget.dataset;
      if (value === undefined || value === this.data.active) return;

      const itemIndex = Number(index);
      this.triggerEvent('change', {
        value,
        item: this.data.normalizedItems[itemIndex] ? this.data.normalizedItems[itemIndex].raw : undefined,
        index: itemIndex,
      });
    },
  },
});
