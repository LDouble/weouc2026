import { get, post, put } from '~/api/request';

export function getStudentProfile() {
  return get('/student');
}

export function createStudentProfile(data) {
  return post('/student', data);
}

export function updateStudentProfile(data) {
  return put('/student', data);
}
