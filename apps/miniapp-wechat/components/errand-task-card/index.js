Component({
  properties: {
    task: {
      type: Object,
      value: {},
    },
  },

  methods: {
    onSelect() {
      this.triggerEvent('select', { task: this.data.task });
    },

    onAccept() {
      this.triggerEvent('accept', { task: this.data.task });
    },
  },
});
