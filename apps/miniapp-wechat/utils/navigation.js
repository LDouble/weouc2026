export function getMenuButtonSafeArea(extraGap = 8) {
  const fallback = {
    right: 16,
    top: 40,
    height: 32,
  };

  try {
    const windowInfo = wx.getWindowInfo ? wx.getWindowInfo() : wx.getSystemInfoSync();
    const menuButton = wx.getMenuButtonBoundingClientRect ? wx.getMenuButtonBoundingClientRect() : null;

    if (!windowInfo || !menuButton || !menuButton.left) return fallback;

    return {
      right: Math.max(fallback.right, windowInfo.windowWidth - menuButton.left + extraGap),
      top: menuButton.top || fallback.top,
      height: menuButton.height || fallback.height,
    };
  } catch (error) {
    return fallback;
  }
}
