import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 30000
})

api.interceptors.request.use((config) => {
  let token = localStorage.getItem('adminToken')
  
  console.log('[axios] Request config:', {
    url: config.url,
    method: config.method,
    tokenExists: !!token,
    token: token ? token.substring(0, 20) + '...' : 'none'
  })
  
  if (!token) {
    try {
      const authStore = useAuthStore()
      token = authStore.token
    } catch (e) {
      console.log('[axios] Error getting authStore:', e)
    }
  }
  
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export { default as iamApi } from './modules/iam'
export { default as portalApi } from './modules/portal'
export { default as campusLifeApi } from './modules/campus-life'
export { default as analyticsApi } from './modules/analytics'
export { default as notificationApi } from './modules/notification'

export type {
  GenericObjectResponseEnvelope,
  ListResponseEnvelope,
  IDResponseEnvelope,
  StudentProfile,
  PortalNotice,
  NotificationPublishRequest,
  ReviewUpdateRequest,
  AuditLog,
  DashboardStats,
  UserInfo,
  RoleInfo,
  PermissionInfo
} from './types'

export default api