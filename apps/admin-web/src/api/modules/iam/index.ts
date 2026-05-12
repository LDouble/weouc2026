import api from '../../index'
import type { 
  GenericObjectResponseEnvelope, 
  StudentProfileResponseEnvelope
} from '../../types'

export interface WeChatLoginRequest {
  code: string
  app_id: string
}

export interface WeChatLoginResponse {
  token: string
  openid: string
  userInfo: {
    userId: string
    nickname: string
    avatarUrl: string
  }
}

export interface AdminLoginRequest {
  username: string
  password: string
}

export interface AdminLoginResponse {
  token: string
  user_id: string
  username: string
  roles: string[]
  permissions: string[]
}

export interface SendCaptchaRequest {
  sid: string
}

export interface BindStudentRequest {
  student_id: string
  password: string
  captcha: string
}

export interface UpdateStudentRequest {
  is_bound: boolean
}

const iamApi = {
  login: async (data: WeChatLoginRequest): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.post('/auth/wechat/login', data)
    return response.data
  },

  adminLogin: async (data: AdminLoginRequest): Promise<AdminLoginResponse> => {
    const response = await api.post('/auth/admin/login', data)
    return response.data.data
  },

  getCurrentStudent: async (): Promise<StudentProfileResponseEnvelope> => {
    const response = await api.get('/student')
    return response.data
  },

  bindStudent: async (data: BindStudentRequest): Promise<StudentProfileResponseEnvelope> => {
    const response = await api.post('/student', data)
    return response.data
  },

  unbindStudent: async (): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.put('/student', { is_bound: false })
    return response.data
  },

  sendCaptcha: async (data: SendCaptchaRequest): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.post('/edu/send-captcha', data)
    return response.data
  }
}

export default iamApi
