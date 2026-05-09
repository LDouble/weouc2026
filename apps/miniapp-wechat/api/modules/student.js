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

export function getCourses() {
  return get('/student/courses');
}

export function syncCourses(data) {
  return post('/student/courses', data);
}

export function getGradeDetails() {
  return get('/grade-detail');
}

export function syncGradeDetails(data) {
  return post('/grade-detail', data);
}

export function getRankingConfig(params) {
  return get('/ranking/config', params);
}

export function calculateRanking(data) {
  return post('/ranking/dynamic', data);
}
