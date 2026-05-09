import createStore from './createStore';

const DEFAULT_SESSION_STATE = {
  token: '',
  openId: '',
  authenticated: false,
  viewer: null,
  profile: null,
  profileLoaded: false,
  profileLoading: false,
  lastProfileSyncAt: 0,
};

const sessionStore = createStore(DEFAULT_SESSION_STATE);

function normalizeViewer(viewer) {
  if (!viewer || typeof viewer !== 'object') return null;

  const nickname = viewer.nickname || viewer.name || '';
  const avatarUrl = viewer.avatarUrl || viewer.avatar_url || viewer.avatar || '';

  if (!nickname && !avatarUrl) return null;

  return {
    nickname,
    avatarUrl,
  };
}

export function getSessionState() {
  return sessionStore.getState();
}

export function subscribeSession(listener) {
  return sessionStore.subscribe(listener);
}

export function hydrateSessionState(payload = {}) {
  return sessionStore.setState({
    token: payload.token || '',
    openId: payload.openId || '',
    authenticated: Boolean(payload.token),
  });
}

export function setSessionViewer(viewer) {
  return sessionStore.setState({
    viewer: normalizeViewer(viewer),
  });
}

export function setSessionProfileLoading(profileLoading) {
  return sessionStore.setState({
    profileLoading: Boolean(profileLoading),
  });
}

export function setSessionProfile(profile) {
  const currentState = sessionStore.getState();
  const viewer = profile
    ? normalizeViewer({
      nickname: profile.name,
      avatarUrl: profile.avatarUrl,
    }) || currentState.viewer
    : currentState.viewer;

  return sessionStore.setState({
    profile: profile || null,
    viewer,
    profileLoaded: true,
    profileLoading: false,
    lastProfileSyncAt: Date.now(),
  });
}

export function clearSessionProfile() {
  return sessionStore.setState({
    profile: null,
    profileLoaded: true,
    profileLoading: false,
    lastProfileSyncAt: Date.now(),
  });
}

export function resetSessionState() {
  return sessionStore.reset();
}
