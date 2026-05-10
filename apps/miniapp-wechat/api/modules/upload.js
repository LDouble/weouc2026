import { get, post } from '~/api/request';

const COS = require('cos-wx-sdk-v5');
const { md5 } = require('js-md5');

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

function fetchCOSSTS() {
  return get('/upload/cos-sts').then(getUploadResultData);
}

function fetchPresignedGET(objectKey) {
  return post('/upload/presigned-get', { path: objectKey }).then(getUploadResultData);
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
      const error = new Error(msg);
      error.statusCode = (e && e.statusCode) || 0;
      error.code = (e && e.code) || '';
      error.response = e && e.response;
      reject(error);
    });
  });
}
