import {
  getAcademicSchedule,
  listAcademicExams,
  listAcademicGrades,
} from '~/api/modules/academic';
import { unwrapPayload } from './shared';

function roundTo(value, digits = 2) {
  const base = 10 ** digits;
  return Math.round(Number(value || 0) * base) / base;
}

function normalizePercentage(value) {
  const number = Number(value || 0);
  if (number <= 0) return 0;
  if (number >= 100) return 100;
  return roundTo(number, 0);
}

export async function loadAcademicDashboardModel() {
  const [scheduleResponse, examResponse, gradeResponse] = await Promise.all([
    getAcademicSchedule(),
    listAcademicExams(),
    listAcademicGrades(),
  ]);

  const schedulePayload = unwrapPayload(scheduleResponse);
  const examPayload = unwrapPayload(examResponse);
  const gradePayload = unwrapPayload(gradeResponse);

  const scheduleList = Array.isArray(schedulePayload.list) ? schedulePayload.list : [];
  const examList = Array.isArray(examPayload.list) ? examPayload.list : [];
  const gradeList = Array.isArray(gradePayload.list) ? gradePayload.list : [];
  const gradeSummary = gradePayload.summary || {};
  const scheduleSummary = schedulePayload.summary || {};

  const passedCount = Number(gradeSummary.passed_count || 0);
  const gradeCount = Number(gradeSummary.course_count || gradeList.length || 0);
  const excellentCount = gradeList.filter((item) => Number(item.score || 0) >= 85).length;
  const passRate = gradeCount > 0 ? (passedCount * 100) / gradeCount : 0;
  const excellentRate = gradeCount > 0 ? (excellentCount * 100) / gradeCount : 0;

  return {
    totalSituationDataList: [
      { name: '本学期课程', number: scheduleList.length },
      { name: '考试安排', number: examList.length },
      { name: '成绩条目', number: gradeCount },
    ],
    interactionSituationDataList: [
      { name: '每周上课天数', number: Number(scheduleSummary.teaching_days || 0) },
      { name: '平均成绩', number: roundTo(gradeSummary.average_score || 0, 1) },
      { name: '平均绩点', number: roundTo(gradeSummary.average_grade_point || 0, 2) },
    ],
    completeRateDataList: [
      { time: '课程通过率', percentage: normalizePercentage(passRate) },
      { time: '优秀率(>=85)', percentage: normalizePercentage(excellentRate) },
    ],
  };
}
