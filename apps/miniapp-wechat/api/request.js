import config from '~/config';
import { ensureLogin, clearToken } from './auth';

const { baseUrl } = config;

let _loginPromise = null;

function request(url, method = 'GET', data = {}, options = {}) {
  const header = {
    'content-type': 'application/json',
  };

  if (!options.skipAuth) {
    const token = wx.getStorageSync('access_token');
    if (token) {
      header.Authorization = `Bearer ${token}`;
    }
  }

  return new Promise((resolve, reject) => {
    const doRequest = (isRetry) => {
      wx.request({
        url: baseUrl + url,
        method,
        data,
        dataType: 'json',
        header,
        success(res) {
          if (res.statusCode === 401 && !isRetry) {
            if (!_loginPromise) {
              _loginPromise = ensureLogin(true)
                .then((result) => {
                  _loginPromise = null;
                  return result;
                })
                .catch((error) => {
                  _loginPromise = null;
                  clearToken();
                  wx.redirectTo({ url: '/pages/login/login' });
                  throw error;
                });
            }

            _loginPromise
              .then(() => {
                const newToken = wx.getStorageSync('access_token');
                if (newToken) {
                  header.Authorization = `Bearer ${newToken}`;
                } else {
                  delete header.Authorization;
                }
                doRequest(true);
              })
              .catch((err) => {
                reject(err);
              });
            return;
          }

          if (res.statusCode >= 200 && res.statusCode < 300) {
            resolve(res.data);
          } else {
            reject(res);
          }
        },
        fail(err) {
          if (!isRetry) {
            doRequest(true);
          } else {
            reject(err);
          }
        },
      });
    };

    doRequest(false);
  });
}

export function get(url, data, options) {
  return request(url, 'GET', data, options);
}

export function post(url, data, options) {
  return request(url, 'POST', data, options);
}

export function put(url, data, options) {
  return request(url, 'PUT', data, options);
}

export function del(url, data, options) {
  return request(url, 'DELETE', data, options);
}

export default request;
