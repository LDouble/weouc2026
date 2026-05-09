import config from '~/config';
import request from './request';

const TOKEN_KEY = 'access_token';
const OPENID_KEY = 'openid';

let _loginPromise = null;

export function getToken() {
  return wx.getStorageSync(TOKEN_KEY) || '';
}

export function getOpenId() {
  return wx.getStorageSync(OPENID_KEY) || '';
}

export function isLoggedIn() {
  return !!getToken();
}

export function setToken(token, openid) {
  wx.setStorageSync(TOKEN_KEY, token);
  if (openid) wx.setStorageSync(OPENID_KEY, openid);
}

export function clearToken() {
  wx.removeStorageSync(TOKEN_KEY);
  wx.removeStorageSync(OPENID_KEY);
}

export function wxLogin() {
  if (_loginPromise) return _loginPromise;

  _loginPromise = new Promise((resolve, reject) => {
    wx.login({
      success: (loginRes) => {
        if (!loginRes.code) {
          _loginPromise = null;
          reject(new Error('wx.login failed'));
          return;
        }
        request('/auth/wechat/login', 'POST', {
          code: loginRes.code,
          app_id: config.appId,
        }, { skipAuth: true })
          .then((res) => {
            const data = res.data || res;
            const token = data.token || (data.data && data.data.token) || '';
            const openid = data.openid || (data.data && data.data.openid) || '';
            if (token) {
              setToken(token, openid);
              resolve({ token, openid, data });
            } else {
              _loginPromise = null;
              reject(new Error('No token in login response'));
            }
          })
          .catch((err) => {
            _loginPromise = null;
            reject(err);
          })
          .finally(() => {
            _loginPromise = null;
          });
      },
      fail: (err) => {
        _loginPromise = null;
        reject(err);
      },
    });
  });

  return _loginPromise;
}

export function ensureLogin(force = false) {
  if (!force && isLoggedIn()) return Promise.resolve({ token: getToken() });
  return wxLogin();
}

export function silentLogin() {
  return ensureLogin().catch(() => {});
}
