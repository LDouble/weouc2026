import { get, post } from '~/api/request';

export function fetchMeetupList(params = {}) {
  const { category, keyword, user_role, page = 1, pageSize = 20 } = params;
  const query = { page, pageSize };
  if (category) query.category = category;
  if (keyword) query.keyword = keyword;
  if (user_role) query.user_role = user_role;
  return get('/meetup/list', query);
}

export function fetchMeetupDetail(id) {
  return get(`/meetup/detail/${id}`);
}

export function publishMeetup(data) {
  return post('/meetup/publish', data);
}

export function joinMeetup(meetupId) {
  return post('/meetup/join', { meetup_id: meetupId });
}

export function cancelMeetupJoin(meetupId) {
  return post('/meetup/cancel-join', { meetup_id: meetupId });
}

export function cancelMeetupPublish(meetupId) {
  return post('/meetup/cancel-publish', { meetup_id: meetupId });
}
