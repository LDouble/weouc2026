import api from '../../index'
import type { 
  GenericObjectResponseEnvelope, 
  ListResponseEnvelope
} from '../../types'

const analyticsApi = {
  getDashboard: async (): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get('/admin/analytics/dashboard')
    return response.data
  },

  listAuditLogs: async (params?: {
    page?: number
    pageSize?: number
    actor_id?: string
    action?: string
    resource_type?: string
    resource_id?: string
  }): Promise<ListResponseEnvelope> => {
    const response = await api.get('/admin/analytics/audit-logs', { params })
    return response.data
  },

  getAuditLogDetail: async (logId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get(`/admin/analytics/audit-logs/${logId}`)
    return response.data
  }
}

export default analyticsApi