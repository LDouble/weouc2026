<template>
  <div class="bg-white rounded-lg shadow p-6">
    <a-alert
      type="info"
      show-icon
      message="文章管理尚未接入后端真实能力"
      description="当前 portal 模块已支持轮播维护与公告发布/编辑/删除，文章模型仍未在后端定义。本页面保留为占位入口，避免继续请求不存在的接口。"
    />

    <div class="mt-6 flex gap-3">
      <a-button type="primary" @click="router.push('/portal/notices')">前往公告管理</a-button>
      <a-button @click="router.push('/portal/banners')">前往轮播管理</a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>
