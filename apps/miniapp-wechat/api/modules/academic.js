import { get } from '~/api/request';

export function listAcademicSemesters() {
  return get('/academic/semesters');
}

export function getAcademicSchedule(params = {}) {
  const query = {};
  if (params.semesterId) query.semester_id = params.semesterId;
  return get('/academic/schedule', query);
}

export function listAcademicExams(params = {}) {
  const query = {};
  if (params.semesterId) query.semester_id = params.semesterId;
  return get('/academic/exams', query);
}

export function listAcademicGrades(params = {}) {
  const query = {};
  if (params.semesterId) query.semester_id = params.semesterId;
  return get('/academic/grades', query);
}
