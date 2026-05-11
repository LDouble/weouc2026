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
  image_url: string
  link_url: string
  position: number
  status: string
  created_at: string
}

export interface ArticleInfo {
  id: string
  title: string
  summary: string
  content: string
  author: string
  category: string
  tags: string[]
  status: string
  created_at: string
  updated_at: string
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
  },

  listArticles: async (params?: {
    page?: number
    pageSize?: number
    keyword?: string
    category?: string
  }): Promise<ListResponseEnvelope> => {
    const response = await api.get('/admin/portal/articles', { params })
    return response.data
  },

  getArticleDetail: async (articleId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get(`/admin/portal/articles/${articleId}`)
    return response.data
  },

  createArticle: async (data: Partial<ArticleInfo>): Promise<IDResponseEnvelope> => {
    const response = await api.post('/admin/portal/articles', data)
    return response.data
  },

  updateArticle: async (articleId: string, data: Partial<ArticleInfo>): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.put(`/admin/portal/articles/${articleId}`, data)
    return response.data
  },

  deleteArticle: async (articleId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.delete(`/admin/portal/articles/${articleId}`)
    return response.data
  }
}

export default portalApi