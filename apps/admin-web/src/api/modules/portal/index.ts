import api from '../../index'
import type { 
  GenericObjectResponseEnvelope, 
  ListResponseEnvelope,
  IDResponseEnvelope,
  PortalNotice,
  PortalNoticePublishRequest
} from '../../types'

export interface BannerInfo {
  id: string
  title: string
  description: string
  image_url: string
  action_url: string
  sort: number
  created_at: string
}

const portalApi = {
  getHome: async (): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get('/portal/home')
    return response.data
  },

  listNotices: async (params?: {
    page?: number
    pageSize?: number
    keyword?: string
  }): Promise<ListResponseEnvelope> => {
    const response = await api.get('/portal/notices', { params })
    return response.data
  },

  getNoticeDetail: async (noticeId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get(`/portal/notices/${noticeId}`)
    return response.data
  },

  publishNotice: async (data: PortalNoticePublishRequest): Promise<IDResponseEnvelope> => {
    const response = await api.post('/admin/portal/notices/publish', data)
    return response.data
  },

  updateNotice: async (noticeId: string, data: Partial<PortalNotice>): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.put(`/admin/portal/notices/${noticeId}`, data)
    return response.data
  },

  deleteNotice: async (noticeId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.delete(`/admin/portal/notices/${noticeId}`)
    return response.data
  },

  listBanners: async (params?: {
    page?: number
    pageSize?: number
    keyword?: string
  }): Promise<ListResponseEnvelope> => {
    const response = await api.get('/admin/portal/banners', { params })
    return response.data
  },

  getBannerDetail: async (bannerId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get(`/admin/portal/banners/${bannerId}`)
    return response.data
  },

  createBanner: async (data: Partial<BannerInfo>): Promise<IDResponseEnvelope> => {
    const response = await api.post('/admin/portal/banners', data)
    return response.data
  },

  updateBanner: async (bannerId: string, data: Partial<BannerInfo>): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.put(`/admin/portal/banners/${bannerId}`, data)
    return response.data
  },

  deleteBanner: async (bannerId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.delete(`/admin/portal/banners/${bannerId}`)
    return response.data
  }
}

export default portalApi
