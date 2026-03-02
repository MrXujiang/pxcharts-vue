<template>
  <div class="user-page-container">
    <!-- Logo -->
    <div class="page-logo">
      <div class="logo-icon">AI</div>
      <div class="logo-text">迈维AI表格</div>
    </div>
    
    <!-- 主内容 -->
    <div class="page-content">
      <div class="user-card">
        <div class="user-card-header">
          <h1 class="title">迈维表格<span class="highlight">登录</span></h1>
          <p class="subtitle">——一张多维表，装下所有业务</p>
        </div>
        
        <t-form
          ref="formRef"
          :data="formData"
          :rules="rules"
          :label-width="0"
          class="user-form"
          @submit="onSubmit"
        >
          <t-form-item name="email">
            <t-input
              v-model="formData.email"
              placeholder="请输入用户名/邮箱"
              size="large"
              clearable
            />
          </t-form-item>
          
          <t-form-item name="password">
            <t-input
              v-model="formData.password"
              type="password"
              placeholder="请输入密码"
              size="large"
              clearable
            />
          </t-form-item>
          
          <div class="form-footer">
            <t-link theme="primary" hover="color" @click="goToForget">忘记密码？</t-link>
          </div>
          
          <t-form-item>
            <t-button
              theme="primary"
              type="submit"
              block
              size="large"
              :loading="loading"
            >
              登录
            </t-button>
          </t-form-item>
        </t-form>
        
        <div class="register-link">
          还没有账号？
          <t-link theme="primary" hover="color" @click="goToRegister">立即注册</t-link>
        </div>
        
        <div class="agreement-text">
          登录即代表您已阅读并同意
          <t-link theme="primary" hover="color" @click="showAgreement('user')">《用户协议》</t-link>
          和
          <t-link theme="primary" hover="color" @click="showAgreement('privacy')">《隐私政策》</t-link>
        </div>
      </div>
      
      <!-- 更多产品 -->
      <div class="more-products">
        <div class="products-title">探索更多产品</div>
        <div class="products-list">
          <a href="https://jitword.com" target="_blank" class="product-item">
            <div class="product-info">
              <div class="product-name">jitword</div>
              <div class="product-desc">AI协同文档</div>
            </div>
            <span class="product-arrow">→</span>
          </a>
          <a href="https://orange.turntip.cn" target="_blank" class="product-item">
            <div class="product-info">
              <div class="product-name">橙子轻文档</div>
              <div class="product-desc">轻量级文档工具</div>
            </div>
            <span class="product-arrow">→</span>
          </a>
          <a href="https://mindlink.turntip.cn" target="_blank" class="product-item">
            <div class="product-info">
              <div class="product-name">灵语文档</div>
              <div class="product-desc">智能协作平台</div>
            </div>
            <span class="product-arrow">→</span>
          </a>
          <a href="https://pxcharts.turntip.cn" target="_blank" class="product-item">
            <div class="product-info">
              <div class="product-name">pxcharts</div>
              <div class="product-desc">超级表格工具</div>
            </div>
            <span class="product-arrow">→</span>
          </a>
        </div>
      </div>
    </div>
    
    <!-- 底部波浪 -->
    <div class="wave-decoration"></div>
    
    <!-- 协议弹窗 -->
    <AgreementDialog v-model="agreementVisible" :type="agreementType" />
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import type { FormInstanceFunctions, FormRule, SubmitContext } from 'tdesign-vue-next'
import AgreementDialog from './components/AgreementDialog.vue'
import { login as loginApi } from '@/api'
import { useUserStore } from '@/stores/user'

defineOptions({
  name: 'UserLogin',
})

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstanceFunctions>()
const loading = ref(false)
const agreementVisible = ref(false)
const agreementType = ref<'user' | 'privacy'>('user')

const formData = reactive({
  email: '',
  password: '',
})

const rules: Record<string, FormRule[]> = {
  email: [
    { required: true, message: '请输入邮箱', type: 'error' },
    { email: true, message: '请输入正确的邮箱格式', type: 'error' },
  ],
  password: [
    { required: true, message: '请输入密码', type: 'error' },
    { min: 6, message: '密码长度不能少于6位', type: 'error' },
  ],
}

const onSubmit = async (ctx: SubmitContext) => {
  if (ctx.validateResult !== true) {
    return
  }
  
  loading.value = true
  try {
    // 调用登录接口
    const response:any = await loginApi({
      email: formData.email,
      password: formData.password,
    })
    
    // 假设接口返回格式: { code: 0, data: { token: string, user: {...} } }
    if (response) {
      const { token, userInfo: user } = response
      
      // 保存用户信息到 store 和 localStorage
      userStore.login(token, user)
      
      MessagePlugin.success('登录成功')
      // 跳转到主页
      router.push('/home')
    }
  } catch (error: any) {
    console.error('登录失败:', error)
    MessagePlugin.error(error.message || '登录失败，请检查邮箱和密码')
  } finally {
    loading.value = false
  }
}

const goToRegister = () => {
  router.push('/user/register')
}

const goToForget = () => {
  router.push('/user/forget')
}

const showAgreement = (type: 'user' | 'privacy') => {
  agreementType.value = type
  agreementVisible.value = true
}
</script>

<style scoped lang="less">
@import './user.less';
</style>
