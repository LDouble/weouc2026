import config from '~/config';
import { ensureLogin, clearToken } from './auth';

const { baseUrl } = config;

let _loginPromise = null;

function createRequestError(message, options = {}) {
  const error = new Error(message || '请求失败');
  error.statusCode = options.statusCode || 0;
  error.code = options.code || '';
  error.details = options.details || null;
  error.response = options.response;
  return error;
}

function extractErrorInfo(body, fallbackMessage) {
  const errorObject = body && typeof body === 'object' ? body.error : null;
  const message = (errorObject && errorObject.message)
    || (body && body.message)
    || (typeof errorObject === 'string' ? errorObject : '')
    || fallbackMessage
    || '请求失败';
  const code = (errorObject && errorObject.code)
    || (body && body.code)
    || '';

  return {
    message,
    code,
    details: errorObject && typeof errorObject === 'object' ? errorObject.details || null : null,
  };
}

function isBusinessSuccess(body) {
  if (body && typeof body === 'object' && body.error) {
    return false;
  }

  if (!body || typeof body !== 'object' || !Object.prototype.hasOwnProperty.call(body, 'code')) {
    return true;
  }

  const code = Number(body.code);
  if (Number.isNaN(code)) {
    return true;
  }

  return code >= 200 && code < 300;
}

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
          if (res.statusCode === 401 && !isRetry && !options.skipAuth) {
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
            const body = res.data;
            if (isBusinessSuccess(body)) {
              resolve(body);
              return;
            }

            const errorInfo = extractErrorInfo(body, `请求失败(${res.statusCode})`);
            reject(
              createRequestError(errorInfo.message, {
                statusCode: res.statusCode,
                code: errorInfo.code,
                details: errorInfo.details,
                response: res,
              }),
            );
            return;
          } else {
            const body = res.data || {};
            const errorInfo = extractErrorInfo(body, `请求失败(${res.statusCode})`);
            reject(
              createRequestError(errorInfo.message, {
                statusCode: res.statusCode,
                code: errorInfo.code,
                details: errorInfo.details,
                response: res,
              }),
            );
          }
        },
        fail(err) {
          if (!isRetry) {
            doRequest(true);
          } else {
            reject(
              createRequestError('网络异常，请稍后再试', {
                response: err,
              }),
            );
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
