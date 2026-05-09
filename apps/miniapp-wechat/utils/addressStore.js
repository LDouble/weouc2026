export const HISTORY_ADDRESSES_KEY = 'carpoolHistoryAddresses';
export const FAVORITE_ADDRESSES_KEY = 'carpoolFavoriteAddresses';

export function getHistoryAddresses() {
  try {
    const addresses = wx.getStorageSync(HISTORY_ADDRESSES_KEY);
    return Array.isArray(addresses) ? addresses : [];
  } catch (error) {
    return [];
  }
}

export function saveHistoryAddress(address) {
  if (!address) return;
  try {
    const addresses = getHistoryAddresses();
    const nextAddresses = [address].concat(addresses.filter((item) => item !== address)).slice(0, 10);
    wx.setStorageSync(HISTORY_ADDRESSES_KEY, nextAddresses);
  } catch (error) {
    wx.setStorageSync(HISTORY_ADDRESSES_KEY, [address]);
  }
}

export function getFavoriteAddresses() {
  try {
    const addresses = wx.getStorageSync(FAVORITE_ADDRESSES_KEY);
    return Array.isArray(addresses) ? addresses : [];
  } catch (error) {
    return [];
  }
}

export function addFavoriteAddress(address) {
  if (!address) return false;
  try {
    const addresses = getFavoriteAddresses();
    if (addresses.includes(address)) return true;
    if (addresses.length >= 8) {
      wx.showToast({ title: '常用地址最多8个', icon: 'none' });
      return false;
    }
    const nextAddresses = [address].concat(addresses);
    wx.setStorageSync(FAVORITE_ADDRESSES_KEY, nextAddresses);
    return true;
  } catch (error) {
    wx.setStorageSync(FAVORITE_ADDRESSES_KEY, [address]);
    return true;
  }
}

export function removeFavoriteAddress(address) {
  try {
    const addresses = getFavoriteAddresses();
    const nextAddresses = addresses.filter((item) => item !== address);
    wx.setStorageSync(FAVORITE_ADDRESSES_KEY, nextAddresses);
  } catch (error) {}
}
