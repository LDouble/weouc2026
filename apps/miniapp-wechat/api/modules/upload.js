import config from '~/config';
import { clearToken, ensureLogin } from '~/api/auth';

const COS = require('cos-wx-sdk-v5');
const { md5 } = require('js-md5');

let _loginPromise = null;

function getAuthHeader() {
  const token = wx.getStorageSync('access_token');
  return token ? { Authorization: `Bearer ${token}` } : {};
}

function pickExt(name, filePath) {
  const n = name || (filePath && filePath.split('/').pop()) || '';
  const i = n.lastIndexOf('.');
  return i > 0 ? n.slice(i).toLowerCase() : '';
}

function readFileArrayBuffer(filePath) {
  return new Promise((resolve, reject) => {
    wx.getFileSystemManager().readFile({
      filePath,
      success: (res) => resolve(res.data),
      fail: reject,
    });
  });
}

function jsonRequest(method, path, data, isRetry) {
  return new Promise((resolve, reject) => {
    wx.request({
      url: `${config.baseUrl}${path}`,
      method,
      data: method === 'GET' ? undefined : data,
      header: {
        'Content-Type': 'application/json',
        ...getAuthHeader(),
      },
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
              jsonRequest(method, path, data, true).then(resolve).catch(reject);
            })
            .catch(reject);
          return;
        }

        const body = res.data || {};
        if (res.statusCode >= 200 && res.statusCode < 300 && body.code === 200) {
          resolve(body);
          return;
        }
        const msg = (body && body.message) || `请求失败(${res.statusCode})`;
        const err = new Error(msg);
        err.response = res;
        reject(err);
      },
      fail: reject,
    });
  });
}

function fetchCOSSTS() {
  return jsonRequest('GET', '/upload/cos-sts', null, false).then((body) => body.data);
}

function fetchPresignedGET(objectKey) {
  return jsonRequest('POST', '/upload/presigned-get', { path: objectKey }, false).then((body) => body.data);
}

export function getUploadResultData(uploadResponse) {
  if (!uploadResponse) return null;
  return uploadResponse.data || uploadResponse;
}

/** COS GET 预签名临时链接 */
export function getUploadResultUrl(uploadResponse) {
  const data = getUploadResultData(uploadResponse);
  if (!data) return '';
  return data.url || '';
}

export function getUploadResultId(uploadResponse) {
  const data = getUploadResultData(uploadResponse);
  if (!data) return '';
  return data.id || data.file_id || '';
}

/** 持久化用的 COS 对象键 */
export function getUploadResultPath(uploadResponse) {
  const data = getUploadResultData(uploadResponse);
  if (!data) return '';
  return data.path || '';
}

/**
 * 使用 STS 临时密钥直传 COS，成功后向后端申请 GET 预签名用于展示。
 * 依赖 npm：cos-wx-sdk-v5、js-md5，并在微信开发者工具中执行「构建 npm」。
 */
export function uploadFile(filePath, options = {}) {
  return new Promise((resolve, reject) => {
    if (!filePath) {
      reject(new Error('文件路径无效，请重新选择文件'));
      return;
    }

    const run = async () => {
      const sts = await fetchCOSSTS();
      const buf = await readFileArrayBuffer(filePath);
      const hash = md5(buf);
      const ext = pickExt(options.name, filePath);
      const objectKey = `${sts.path_prefix}${hash.slice(0, 2)}/${hash}${ext}`;

      const cos = new COS({
        getAuthorization(opts, callback) {
          callback({
            TmpSecretId: sts.tmp_secret_id,
            TmpSecretKey: sts.tmp_secret_key,
            SecurityToken: sts.session_token,
            XCosSecurityToken: sts.session_token,
            ExpiredTime: sts.expired_time,
            StartTime: sts.start_time,
          });
        },
      });

      await new Promise((res, rej) => {
        cos.putObject(
          {
            Bucket: sts.bucket,
            Region: sts.region,
            Key: objectKey,
            FilePath: filePath,
          },
          (err) => {
            if (err) rej(err);
            else res();
          },
        );
      });

      const presign = await fetchPresignedGET(objectKey);
      const name = (options.name || objectKey.split('/').pop() || '').trim();

      resolve({
        code: 200,
        message: 'success',
        data: {
          path: objectKey,
          url: presign.url,
          name,
          md5: hash,
        },
      });
    };

    run().catch((e) => {
      const msg = (e && e.message) || String(e);
      reject(new Error(msg));
    });
  });
}
