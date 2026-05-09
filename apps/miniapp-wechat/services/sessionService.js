import { clearToken, getOpenId, getToken, wxLogin } from '~/api/auth';
import {
  clearSessionProfile,
  getSessionState,
  hydrateSessionState,
  resetSessionState,
  setSessionProfileLoading,
  setSessionViewer,
} from '~/stores/session';
import { tryLoadCurrentProfile } from './profileService';

function extractViewer(loginResult) {
  const data = (loginResult && loginResult.data) || {};
  const userInfo = data.userInfo || {};

  if (!userInfo.userId && !userInfo.nickname && !userInfo.avatarUrl) {
    return null;
  }

  return {
    nickname: userInfo.nickname || '',
    avatarUrl: userInfo.avatarUrl || '',
  };
}

function syncSessionFromStorage() {
  hydrateSessionState({
    token: getToken(),
    openId: getOpenId(),
  });
  return getSessionState();
}

async function syncProfileSilently() {
  if (!getToken()) {
    resetSessionState();
    return null;
  }

  setSessionProfileLoading(true);

  try {
    return await tryLoadCurrentProfile();
  } catch (error) {
    clearSessionProfile();
    return null;
  }
}

export async function bootstrapSession() {
  const currentState = syncSessionFromStorage();

  if (!currentState.token) return currentState;

  await syncProfileSilently();
  return getSessionState();
}

export async function loginWithWechat() {
  const loginResult = await wxLogin();
  syncSessionFromStorage();

  const viewer = extractViewer(loginResult);
  if (viewer) {
    setSessionViewer(viewer);
  }

  await syncProfileSilently();

  return {
    loginResult,
    session: getSessionState(),
  };
}

export async function refreshSessionProfile() {
  return syncProfileSilently();
}

export function logoutSession() {
  clearToken();
  resetSessionState();
}
