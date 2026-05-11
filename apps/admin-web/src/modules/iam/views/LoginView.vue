<template>
  <div class="flex items-center justify-center min-h-screen bg-gradient-to-br from-blue-500 to-purple-600">
    <div class="bg-white rounded-2xl shadow-2xl p-8 w-full max-w-md">
      <div class="text-center mb-8">
        <h1 class="text-2xl font-bold text-gray-800 mb-2">校园管理后台</h1>
        <p class="text-gray-500">请登录您的账户</p>
      </div>
      
      <a-form
        :model="form"
        layout="vertical"
        @submit.prevent="handleSubmit"
      >
        <a-form-item
          label="用户名"
          :rules="[{ required: true, message: '请输入用户名' }]"
        >
          <a-input
            v-model:value="form.username"
            placeholder="请输入用户名"
            size="large"
          >
            <template #prefix>
              <component :is="icons.User" />
            </template>
          </a-input>
        </a-form-item>
        
        <a-form-item
          label="密码"
          :rules="[{ required: true, message: '请输入密码' }]"
        >
          <a-input-password
            v-model:value="form.password"
            placeholder="请输入密码"
            size="large"
          >
            <template #prefix>
              <component :is="icons.Lock" />
            </template>
          </a-input-password>
        </a-form-item>
        
        <a-form-item class="mb-6">
          <a-checkbox v-model:checked="form.remember">
            记住我
          </a-checkbox>
        </a-form-item>
        
        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            size="large"
            class="w-full"
            :loading="loading"
          >
            登录
          </a-button>
        </a-form-item>
      </a-form>
      
      <div class="text-center text-gray-400 text-sm">
        忘记密码？联系管理员重置
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { UserOutlined as User, LockOutlined as Lock } from '@ant-design/icons-vue'
import { useAuthStore } from '../../../stores/auth'
import { iamApi } from '@/api'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)

const icons = { User, Lock }

const form = reactive({
  username: '',
  password: '',
  remember: false
})

async function handleSubmit() {
  loading.value = true
  
  try {
    const response = await iamApi.adminLogin({
      username: form.username,
      password: form.password
    })
    
    console.log('[Login] API response:', response)
    
    authStore.login(response.token, {
      id: response.user_id,
      username: response.username,
      role: response.roles?.[0] || 'admin',
      permissions: [
        'iam:user:view', 'iam:user:manage',
        'iam:role:view', 'iam:role:manage',
        'iam:permission:view',
        'portal:publish', 'portal:view',
        'campus_life:moderate', 'campus_life:view',
        'moderation:review',
        'analytics:view'
      ]
    })
    
    console.log('[Login] Token saved to localStorage:', localStorage.getItem('adminToken'))
    
    message.success('登录成功')
    router.push('/')
  } catch (error: any) {
    console.log('[Login] Error:', error)
    message.error(error.response?.data?.message || '登录失败，请检查用户名和密码')
  } finally {
    loading.value = false
  }
}
</script>