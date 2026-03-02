<template>
  <div class="entry-page">
    <!-- Logo -->
    <div class="page-logo">
      <div class="logo-icon">AI</div>
      <div class="logo-text">迈维AI表格</div>
    </div>
    
    <!-- 主内容 -->
    <div class="page-content">
      <h1 class="page-title">请输入您感兴趣的标签，系统将为您推荐更符合您的模版</h1>
      
      <div class="tags-container">
        <div 
          v-for="tag in tags" 
          :key="tag.value"
          class="tag-item"
          :class="{ active: selectedTags.includes(tag.value) }"
          @click="toggleTag(tag.value)"
        >
          {{ tag.label }}
        </div>
      </div>
      
      <div class="action-button">
        <t-button 
          theme="primary" 
          size="large"
          :loading="loading"
          @click="handleSubmit"
        >
          立即体验
        </t-button>
      </div>
    </div>
    
    <!-- 底部波浪装饰 -->
    <div class="wave-decoration"></div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import { updateUserInfo } from '@/api'
import { useUserStore } from '@/stores/user'

defineOptions({
  name: 'UserEntry',
})

const router = useRouter()
const userStore = useUserStore()
const selectedTags = ref<string[]>([])
const loading = ref(false)

const tags = [
  { label: '互联网/科技', value: 'internet' },
  { label: '企业服务/To B', value: 'enterprise' },
  { label: '金融/财经', value: 'finance' },
  { label: '教育/培训', value: 'education' },
  { label: '电商/零售', value: 'ecommerce' },
  { label: '新能源/汽车', value: 'newenergy' },
  { label: '项目管理', value: 'project' },
  { label: '客户/人员管理', value: 'customer' },
  { label: '运营分析', value: 'operation' },
  { label: '资产管理', value: 'asset' },
  { label: '个人/爱好', value: 'personal' },
]

const toggleTag = (value: string) => {
  const index = selectedTags.value.indexOf(value)
  if (index > -1) {
    selectedTags.value.splice(index, 1)
  } else {
    selectedTags.value.push(value)
  }
}

const handleSubmit = async () => {
  if (selectedTags.value.length === 0) {
    MessagePlugin.warning('请至少选择一个感兴趣的标签')
    return
  }
  
  loading.value = true
  try {
    // 调用更新用户信息接口
    await updateUserInfo({
      tags: selectedTags.value,
    })
    
    // 更新本地 store 中的用户信息
    userStore.updateUserInfo({
      tags: selectedTags.value,
    })
    
    MessagePlugin.success('保存成功')
    // 跳转到主页
    router.push('/home')
  } catch (error: any) {
    console.error('保存标签失败:', error)
    MessagePlugin.error(error.message || '保存失败，请稍后重试')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="less">
.entry-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #ffffff;
  position: relative;
  overflow: hidden;
  margin: 0;
  padding: 0;
  
  // Logo区域
  .page-logo {
    position: absolute;
    top: 30px;
    left: 40px;
    display: flex;
    align-items: center;
    gap: 8px;
    z-index: 100;
    
    .logo-icon {
      width: 32px;
      height: 32px;
      background: #4a7ff7;
      border-radius: 6px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: bold;
      color: #fff;
      font-size: 18px;
    }
    
    .logo-text {
      font-size: 16px;
      font-weight: 500;
      color: #333;
    }
  }
  
  // 主内容区域
  .page-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 100px 40px 200px;
    position: relative;
    z-index: 10;
    
    .page-title {
      font-size: 20px;
      font-weight: 400;
      color: #333;
      margin: 0 0 60px;
      text-align: center;
    }
    
    .tags-container {
      display: flex;
      flex-wrap: wrap;
      gap: 20px;
      max-width: 800px;
      justify-content: center;
      margin-bottom: 80px;
      
      .tag-item {
        padding: 12px 32px;
        background: #f5f5f5;
        border-radius: 8px;
        font-size: 15px;
        color: #666;
        cursor: pointer;
        transition: all 0.3s;
        user-select: none;
        
        &:hover {
          background: #e8f0ff;
          color: #4a7ff7;
        }
        
        &.active {
          background: #4a7ff7;
          color: #ffffff;
          font-weight: 500;
        }
      }
    }
    
    .action-button {
      :deep(.t-button) {
        min-width: 160px;
        height: 48px;
        font-size: 16px;
        border-radius: 8px;
        background: #4a7ff7;
        border-color: #4a7ff7;
        
        &:hover {
          background: #6691ff;
          border-color: #6691ff;
        }
      }
    }
  }
  
  // 底部波浪装饰
  .wave-decoration {
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    height: 400px;
    z-index: 1;
    background-image: 
      url('data:image/svg+xml;utf8,<svg viewBox="0 0 1440 400" xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"><path d="M0,160 C320,100 420,200 720,180 C1020,160 1120,60 1440,100 L1440,400 L0,400 Z" fill="%234a7ff7" opacity="0.3"/><path d="M0,200 C360,140 460,240 760,220 C1060,200 1160,100 1440,140 L1440,400 L0,400 Z" fill="%234a7ff7" opacity="0.5"/><path d="M0,240 C400,180 500,280 800,260 C1100,240 1200,140 1440,180 L1440,400 L0,400 Z" fill="%234a7ff7"/></svg>');
    background-size: cover;
    background-position: bottom;
    background-repeat: no-repeat;
  }
}
</style>
