export const NETWORK_RESULT_CONFIRM_MESSAGE = '网络异常，请到我的发布/消息列表确认结果';
export const NETWORK_MESSAGE_CONFIRM_MESSAGE = '网络异常，请到消息列表确认结果';

export function isNetworkFail(error) {
  return Boolean(error && (error.networkFail || error.code === 'NETWORK_FAIL'));
}

export function getNetworkConfirmMessage(error, fallbackMessage) {
  if (isNetworkFail(error)) return NETWORK_RESULT_CONFIRM_MESSAGE;
  return fallbackMessage || (error && error.message) || '操作失败，请重试';
}
