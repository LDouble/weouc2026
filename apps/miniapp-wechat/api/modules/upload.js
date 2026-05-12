import { get, post } from '~/api/request';

const COS = require('cos-wx-sdk-v5');

const REQUIRED_STS_FIELDS = [
  'bucket',
  'region',
  'path_prefix',
  'tmp_secret_id',
  'tmp_secret_key',
  'session_token',
  'start_time',
  'expired_time',
];

function normalizeCOSSTS(response) {
  const data = getUploadResultData(response) || {};
  const missing = REQUIRED_STS_FIELDS.filter(
    (field) => data[field] === undefined || data[field] === null || data[field] === '',
  );
  if (missing.length) {
    throw new Error('上传凭证字段不完整，请稍后重试');
  }
  return data;
}

function normalizePresignedGET(response) {
  const data = getUploadResultData(response) || {};
  if (!data.path || !data.url) {
    throw new Error('文件访问地址生成失败，请稍后重试');
  }
  return data;
}

function pickExt(name, filePath) {
  const n = name || (filePath && filePath.split('/').pop()) || '';
  const i = n.lastIndexOf('.');
  return i > 0 ? n.slice(i).toLowerCase() : '';
}

function generateObjectKey(pathPrefix, filePath) {
  const timestamp = Date.now().toString(36);
  const random = Math.random().toString(36).substr(2, 8);
  const ext = pickExt('', filePath);
  return `${pathPrefix}${timestamp}/${timestamp}${random}${ext}`;
}

function fetchCOSSTS(scene = '') {
  const query = {};
  if (scene) query.scene = scene;
  return get('/upload/cos-sts', query).then(normalizeCOSSTS);
}

function fetchPresignedGET(objectKey) {
  return post('/upload/presigned-get', { path: objectKey }).then(normalizePresignedGET);
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
 * 依赖 npm：cos-wx-sdk-v5，并在微信开发者工具中执行「构建 npm」。
 */
export function uploadFile(filePath, options = {}) {
  return new Promise((resolve, reject) => {
    if (!filePath) {
      reject(new Error('文件路径无效，请重新选择文件'));
      return;
    }

    const run = async () => {
      const sts = await fetchCOSSTS(options.scene || '');
      const objectKey = generateObjectKey(sts.path_prefix, filePath);

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
