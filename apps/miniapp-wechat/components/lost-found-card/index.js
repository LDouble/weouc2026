Component({
  properties: {
    item: {
      type: Object,
      value: {},
    },
  },

  methods: {
    onSelect() {
      this.triggerEvent('select', { item: this.data.item });
    },

    onContact() {
      this.triggerEvent('contact', { item: this.data.item });
    },
  },
});
