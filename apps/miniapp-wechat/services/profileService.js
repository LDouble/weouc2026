import { sendEduCaptcha } from '~/api/modules/edu';
import { fetchErrandList } from '~/api/modules/errand';
import { fetchFeedList } from '~/api/modules/feed';
import { createStudentProfile, getStudentProfile, updateStudentProfile } from '~/api/modules/student';
import { getSessionState, setSessionProfile } from '~/stores/session';
import { formatDateTime } from '~/utils/date';
import { unwrapPayload } from './shared';

function isMissingProfileError(error) {
  return error && error.statusCode === 404;
}

export function mapStudentProfile(raw = {}) {
  const studentId = raw.student_id || '';
  const isBound = Boolean(raw.is_bound || studentId);

  return {
    name: raw.name || raw.nickname || '',
    avatarUrl: raw.avatar_url || raw.avatar || '',
    major: raw.major || '',
    studentId,
    college: raw.college || '',
    grade: raw.grade || '',
    isBound,
    bindingUpdatedAt: raw.updated_at || raw.created_at || '',
    bindingUpdatedLabel: formatDateTime(raw.updated_at || raw.created_at),
    raw,
  };
}

function buildPersonalInfo(profile) {
  const viewer = getSessionState().viewer || {};
  const finalProfile = profile || {};
  const name = finalProfile.name || viewer.nickname || '微信用户';

  return {
    name,
    image: finalProfile.avatarUrl || viewer.avatarUrl || '',
    major: finalProfile.major || '',
    sid: finalProfile.studentId || '',
    college: finalProfile.college || '',
    grade: finalProfile.grade || '',
    isBound: Boolean(finalProfile.isBound),
    identityLabel: finalProfile.isBound ? '已绑定教务' : '未绑定教务',
    secondaryLabel: finalProfile.major || '微信登录用户',
  };
}

async function fetchStats() {
  const [feedResult, errandResult] = await Promise.all([
    fetchFeedList({ page: 1, pageSize: 1 }).catch(() => ({})),
    fetchErrandList({ page: 1, pageSize: 1 }).catch(() => ({})),
  ]);

  const feedData = unwrapPayload(feedResult);
  const errandData = unwrapPayload(errandResult);

  return {
    published: Number(feedData.total || 0),
    accepted: Number(errandData.total || 0),
    favorite: 0,
  };
}

export async function loadCurrentProfile() {
  const response = await getStudentProfile();
  const profile = mapStudentProfile(unwrapPayload(response));
  setSessionProfile(profile);
  return profile;
}

export async function tryLoadCurrentProfile() {
  try {
    return await loadCurrentProfile();
  } catch (error) {
    if (isMissingProfileError(error)) {
      setSessionProfile(null);
      return null;
    }

    throw error;
  }
}

export async function loadMyPageModel() {
  const profile = await tryLoadCurrentProfile();
  const stats = await fetchStats();

  return {
    profile,
    personalInfo: buildPersonalInfo(profile),
    stats,
  };
}

export async function loadAcademicBindingModel() {
  const profile = await tryLoadCurrentProfile();

  return {
    isBound: Boolean(profile && profile.isBound),
    bindInfo: {
      sid: (profile && profile.studentId) || '',
      name: (profile && profile.name) || '',
      bindTime: (profile && profile.bindingUpdatedLabel) || '',
    },
  };
}

export async function sendAcademicCaptcha(studentId) {
  return sendEduCaptcha({
    sid: studentId,
  });
}

export async function bindAcademicAccount(payload = {}) {
  const response = await createStudentProfile({
    student_id: payload.studentId,
    password: payload.password,
    captcha: payload.captcha,
  });
  const profile = mapStudentProfile(unwrapPayload(response));
  setSessionProfile(profile);
  return profile;
}

export async function unbindAcademicAccount() {
  await updateStudentProfile({ is_bound: false });
  const currentProfile = getSessionState().profile || {};

  const nextProfile = {
    ...currentProfile,
    isBound: false,
    studentId: '',
    bindingUpdatedAt: '',
    bindingUpdatedLabel: '',
    raw: {
      ...(currentProfile.raw || {}),
      is_bound: false,
      student_id: '',
    },
  };

  setSessionProfile(nextProfile);
  return nextProfile;
}
