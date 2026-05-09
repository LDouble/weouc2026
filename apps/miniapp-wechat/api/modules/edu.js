import { post } from '~/api/request';

export function sendEduCaptcha(data) {
  return post('/edu/send-captcha', data);
}
