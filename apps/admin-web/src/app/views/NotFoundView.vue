<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-50">
    <a-result
      status="404"
      title="页面未找到"
      sub-title="您访问的页面不存在，请检查地址是否正确"
    >
      <template #extra>
        <a-button type="primary" @click="goHome">
          返回首页
        </a-button>
      </template>
    </a-result>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

function goHome() {
  router.push('/')
}
</script>