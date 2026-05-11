import api from '../../index'
import type { 
  GenericObjectResponseEnvelope, 
  ListResponseEnvelope,
  IDResponseEnvelope,
  NotificationPublishRequest
} from '../../types'

export interface NotificationItem {
  id: string
  title: string
  content: string
  category: string
  target_scope: string
  action_url?: string
  status: string
  created_at: string
}

const notificationApi = {
  listNotifications: async (params?: {
    page?: number
    pageSize?: number
    category?: string
    unread_only?: boolean
  }): Promise<ListResponseEnvelope> => {
    const response = await api.get('/notification/list', { params })
    return response.data
  },

  getUnreadCount: async (): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get('/notification/unread-count')
    return response.data
  },

  markAsRead: async (data: { notification_ids: string[] }): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.post('/notification/read', data)
    return response.data
  },

  publishNotification: async (data: NotificationPublishRequest): Promise<IDResponseEnvelope> => {
    const response = await api.post('/admin/notification/publish', data)
    return response.data
  },

  listAdminNotifications: async (params?: {
    page?: number
    pageSize?: number
    keyword?: string
    category?: string
  }): Promise<ListResponseEnvelope> => {
    const response = await api.get('/admin/notification/list', { params })
    return response.data
  },

  getNotificationDetail: async (notificationId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get(`/admin/notification/${notificationId}`)
    return response.data
  },

  updateNotification: async (notificationId: string, data: Partial<NotificationItem>): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.put(`/admin/notification/${notificationId}`, data)
    return response.data
  },

  deleteNotification: async (notificationId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.delete(`/admin/notification/${notificationId}`)
    return response.data
  }
}

export default notificationApi