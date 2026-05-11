import api from '../../index'
import type { 
  GenericObjectResponseEnvelope, 
  StudentProfileResponseEnvelope,
  UserInfo,
  RoleInfo,
  PermissionInfo
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
  },

  listUsers: async (params?: {
    page?: number
    pageSize?: number
    keyword?: string
  }): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get('/admin/iam/users', { params })
    return response.data
  },

  getUserDetail: async (userId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get(`/admin/iam/users/${userId}`)
    return response.data
  },

  createUser: async (data: Partial<UserInfo>): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.post('/admin/iam/users', data)
    return response.data
  },

  updateUser: async (userId: string, data: Partial<UserInfo>): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.put(`/admin/iam/users/${userId}`, data)
    return response.data
  },

  deleteUser: async (userId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.delete(`/admin/iam/users/${userId}`)
    return response.data
  },

  listRoles: async (params?: {
    page?: number
    pageSize?: number
    keyword?: string
  }): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get('/admin/iam/roles', { params })
    return response.data
  },

  getRoleDetail: async (roleId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get(`/admin/iam/roles/${roleId}`)
    return response.data
  },

  createRole: async (data: Partial<RoleInfo>): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.post('/admin/iam/roles', data)
    return response.data
  },

  updateRole: async (roleId: string, data: Partial<RoleInfo>): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.put(`/admin/iam/roles/${roleId}`, data)
    return response.data
  },

  deleteRole: async (roleId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.delete(`/admin/iam/roles/${roleId}`)
    return response.data
  },

  listPermissions: async (params?: {
    page?: number
    pageSize?: number
    keyword?: string
  }): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get('/admin/iam/permissions', { params })
    return response.data
  },

  getPermissionDetail: async (permissionId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.get(`/admin/iam/permissions/${permissionId}`)
    return response.data
  },

  createPermission: async (data: Partial<PermissionInfo>): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.post('/admin/iam/permissions', data)
    return response.data
  },

  updatePermission: async (permissionId: string, data: Partial<PermissionInfo>): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.put(`/admin/iam/permissions/${permissionId}`, data)
    return response.data
  },

  deletePermission: async (permissionId: string): Promise<GenericObjectResponseEnvelope> => {
    const response = await api.delete(`/admin/iam/permissions/${permissionId}`)
    return response.data
  }
}

export default iamApi