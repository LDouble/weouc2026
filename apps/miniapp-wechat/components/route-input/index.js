import { getHistoryAddresses, saveHistoryAddress, getFavoriteAddresses, addFavoriteAddress, removeFavoriteAddress } from '../../utils/addressStore';

Component({
  properties: {
    startValue: { type: String, value: '' },
    endValue: { type: String, value: '' },
    startPlaceholder: { type: String, value: '出发地' },
    endPlaceholder: { type: String, value: '目的地' },
  },

  data: {
    activeField: '',
    favoriteAddresses: [],
    historyAddresses: [],
    showAddressManager: false,
    newFavoriteAddress: '',
  },

  lifetimes: {
    attached() {
      this.refreshAddresses();
    },
  },

  pageLifetimes: {
    show() {
      this.refreshAddresses();
    },
  },

  methods: {
    refreshAddresses() {
      this.setData({
        favoriteAddresses: getFavoriteAddresses(),
        historyAddresses: getHistoryAddresses(),
      });
    },
    onStartInput(e) {
      this.triggerEvent('startChange', { value: e.detail.value || '' });
    },
    onEndInput(e) {
      this.triggerEvent('endChange', { value: e.detail.value || '' });
    },
    onInputFocus(e) {
      const { field } = e.currentTarget.dataset;
      if (!field) return;
      this.refreshAddresses();
      this.setData({ activeField: field });
    },
    onInputBlur(e) {
      const { field } = e.currentTarget.dataset;
      setTimeout(() => {
        if (this.data.activeField === field) {
          this.setData({ activeField: '' });
        }
      }, 150);
    },
    onAddressTap(e) {
      const { address } = e.currentTarget.dataset;
      const { activeField } = this.data;
      if (!activeField || !address) return;
      if (activeField === 'start') {
        this.triggerEvent('startChange', { value: address });
      } else {
        this.triggerEvent('endChange', { value: address });
      }
      this.setData({ activeField: '' });
    },
    onManageFavorite() {
      this.refreshAddresses();
      this.setData({ showAddressManager: true });
    },
    onCloseManager() {
      this.setData({ showAddressManager: false, newFavoriteAddress: '' });
    },
    onNewFavoriteInput(e) {
      this.setData({ newFavoriteAddress: e.detail.value || '' });
    },
    onAddFavorite() {
      const { newFavoriteAddress } = this.data;
      const trimmed = (newFavoriteAddress || '').trim();
      if (!trimmed) {
        wx.showToast({ title: '请输入地址', icon: 'none' });
        return;
      }
      const success = addFavoriteAddress(trimmed);
      if (success) {
        this.setData({
          favoriteAddresses: getFavoriteAddresses(),
          newFavoriteAddress: '',
        });
      }
    },
    onRemoveFavorite(e) {
      const { address } = e.currentTarget.dataset;
      if (!address) return;
      removeFavoriteAddress(address);
      this.setData({ favoriteAddresses: getFavoriteAddresses() });
    },
    saveHistory(start, end) {
      if (start && start.trim()) saveHistoryAddress(start.trim());
      if (end && end.trim()) saveHistoryAddress(end.trim());
      this.refreshAddresses();
    },
  },
});
