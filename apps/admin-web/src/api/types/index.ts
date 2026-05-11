export interface RequestEnvelope {
  request_id: string
}

export interface GenericData {
  [key: string]: any
}

export interface ErrorDetail {
  code: string
  message: string
  details?: GenericData
}

export interface ErrorResponseEnvelope extends RequestEnvelope {
  error: ErrorDetail
}

export interface GenericObjectResponseEnvelope extends RequestEnvelope {
  data: GenericData
}

export interface IDPayload {
  id: string
}

export interface IDResponseEnvelope extends RequestEnvelope {
  data: IDPayload
}

export interface ListPayload {
  list: GenericData[]
  total: number
  page: number
  pageSize: number
}

export interface ListResponseEnvelope extends RequestEnvelope {
  data: ListPayload
}

export interface StudentProfile {
  name: string
  avatar_url: string
  student_id: string
  major: string
  college: string
  grade: string
  is_bound: boolean
  updated_at: string
}

export interface StudentProfileResponseEnvelope extends RequestEnvelope {
  data: StudentProfile
}

export interface PortalNotice {
  id: string
  title: string
  summary: string
  content: string
  audience: string
  tags: string[]
  pinned: boolean
  status: string
  created_at: string
  updated_at: string
}

export interface PortalNoticePublishRequest {
  title: string
  summary: string
  content: string
  audience: string
  tags: string[]
  pinned?: boolean
}

export interface NotificationPublishRequest {
  title: string
  content: string
  category: string
  target_scope: string
  action_url?: string
}

export interface ReviewUpdateRequest {
  content_type: 'market' | 'errand' | 'resource' | 'lostFound' | 'carpool' | 'meetup'
  content_id: string
  review_status: 'published' | 'rejected' | 'offline'
  reason?: string
}

export interface AuditLog {
  id: string
  actor_id: string
  action: string
  resource_type: string
  resource_id: string
  details: GenericData
  created_at: string
}

export interface DashboardStats {
  user_count: number
  active_users: number
  total_posts: number
  pending_reviews: number
  notifications_sent: number
}

export interface UserInfo {
  id: string
  name: string
  avatar_url: string
  roles: string[]
  permissions: string[]
  is_bound: boolean
  created_at: string
}

export interface RoleInfo {
  id: string
  name: string
  description: string
  permissions: string[]
  created_at: string
}

export interface PermissionInfo {
  id: string
  name: string
  description: string
  resource: string
  action: string
}