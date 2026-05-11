import api from '../../index'
import type { 
  GenericObjectResponseEnvelope, 
  ListResponseEnvelope,
  ReviewUpdateRequest
} from '../../types'

export interface MarketItem {
  id: string
  title: string
  desc: string
  category: string
  price: number
  original_price?: number
  images: string[]
  status: string
  review_status: string
  contact: string
  publisher_id: string
  created_at: string
  updated_at: string
}

export interface ErrandItem {
  id: string
  title: string
  desc: string
  category: string
  reward: number
  location_from: string
  location_to: string
  deadline_at: string
  status: string
  review_status: string
  contact: string
  publisher_id: string
  acceptor_id?: string
  created_at: string
  updated_at: string
}

export interface ResourceItem {
  id: string
  title: string
  desc: string
  category: string
  file_path: string
  status: string
  review_status: string
  contact: string
  publisher_id: string
  created_at: string
}

export interface LostFoundItem {
  id: string
  type: 'lost' | 'found'
  title: string
  desc: string
  category: string
  location: string
  found_at: string
  images: string[]
  status: string
  review_status: string
  contact: string
  publisher_id: string
  created_at: string
}

export interface CarpoolItem {
  id: string
  title: string
  desc: string
  category: string
  location_from: string
  location_to: string
  departure_at: string
  seats: number
  fee_text: string
  status: string
  review_status: string
  contact: string
  publisher_id: string
  created_at: string
}

export interface MeetupItem {
  id: string
  category: string
  title: string
  desc: string
  location: string
  start_at: string
  deadline_at: string
  max_participants: number
  fee_text: string
  tags: string[]
  status: string
  review_status: string
  contact: string
  publisher_id: string
  participant_count: number
  created_at: string
}

const campusLifeApi = {
  listFeeds: async (params?: {
    page?: number
    pageSize?: number
    keyword?: string
  }): Promise<ListResponseEnvelope> => {
    const response = await api.get('/feed/list', { params })
    return response.data
  },

  listMarkets: async (params?: {
    page?: number
    pageSize?: number
    keyword?: string
    category?: string
  }): Promise<ListResponseEnvelope> => {
    const response = await api.get('/market/list', { params })
    return response.data
  },

  getMarketDetail: async (marketId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get(`/market/detail/${marketId}`)
    return response.data
  },

  listErrands: async (params?: {
    page?: number
    pageSize?: number
    keyword?: string
    category?: string
  }): Promise<ListResponseEnvelope> => {
    const response = await api.get('/errand/list', { params })
    return response.data
  },

  getErrandDetail: async (errandId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get(`/errand/detail/${errandId}`)
    return response.data
  },

  listResources: async (params?: {
    page?: number
    pageSize?: number
    keyword?: string
    category?: string
  }): Promise<ListResponseEnvelope> => {
    const response = await api.get('/resource/list', { params })
    return response.data
  },

  getResourceDetail: async (resourceId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get(`/resource/detail/${resourceId}`)
    return response.data
  },

  listLostFound: async (params?: {
    page?: number
    pageSize?: number
    keyword?: string
    category?: string
    type?: 'lost' | 'found'
  }): Promise<ListResponseEnvelope> => {
    const response = await api.get('/lostFound/list', { params })
    return response.data
  },

  getLostFoundDetail: async (lostFoundId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get(`/lostFound/detail/${lostFoundId}`)
    return response.data
  },

  listCarpools: async (params?: {
    page?: number
    pageSize?: number
    keyword?: string
    category?: string
  }): Promise<ListResponseEnvelope> => {
    const response = await api.get('/carpool/list', { params })
    return response.data
  },

  getCarpoolDetail: async (carpoolId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get(`/carpool/detail/${carpoolId}`)
    return response.data
  },

  listMeetups: async (params?: {
    page?: number
    pageSize?: number
    keyword?: string
    category?: string
  }): Promise<ListResponseEnvelope> => {
    const response = await api.get('/meetup/list', { params })
    return response.data
  },

  getMeetupDetail: async (meetupId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get(`/meetup/detail/${meetupId}`)
    return response.data
  },

  listReviews: async (params?: {
    page?: number
    pageSize?: number
    keyword?: string
    content_type?: 'market' | 'errand' | 'resource' | 'lostFound' | 'carpool' | 'meetup'
    review_status?: 'reviewing' | 'published' | 'rejected' | 'offline'
  }): Promise<ListResponseEnvelope> => {
    const response = await api.get('/admin/campus-life/review/list', { params })
    return response.data
  },

  updateReviewStatus: async (data: ReviewUpdateRequest): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.post('/admin/campus-life/review/update', data)
    return response.data
  }
}

export default campusLifeApi