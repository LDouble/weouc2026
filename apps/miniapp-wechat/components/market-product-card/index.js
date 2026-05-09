Component({
  properties: {
    product: {
      type: Object,
      value: {},
    },
  },

  methods: {
    onSelect() {
      this.triggerEvent('select', { product: this.data.product });
    },

    onFavorite() {
      this.triggerEvent('favorite', { product: this.data.product });
    },
  },
});
