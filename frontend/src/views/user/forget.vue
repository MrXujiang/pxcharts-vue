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
          <h1 class="title">迈维表格<span class="highlight">找回密码</span></h1>
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
              placeholder="请输入注册时的邮箱"
              size="large"
              clearable
            />
          </t-form-item>
          
          <t-form-item name="verifyCode">
            <div class="verify-code-wrapper">
              <t-input
                v-model="formData.verifyCode"
                placeholder="请输入验证码"
                size="large"
                clearable
              />
              <t-button
                :disabled="countdown > 0"
                theme="primary"
                size="large"
                @click="sendVerifyCode"
              >
                {{ countdown > 0 ? `发送验证码(${countdown})` : '发送验证码' }}
              </t-button>
            </div>
          </t-form-item>
          
          <t-form-item name="password">
            <t-input
              v-model="formData.password"
              type="password"
              placeholder="请输入新密码（至少6位）"
              size="large"
              clearable
            />
          </t-form-item>
          
          <t-form-item name="confirmPassword">
            <t-input
              v-model="formData.confirmPassword"
              type="password"
              placeholder="请再次输入新密码"
              size="large"
              clearable
            />
          </t-form-item>
          
          <t-form-item>
            <t-button
              theme="primary"
              type="submit"
              block
              size="large"
              :loading="loading"
            >
              重置密码
            </t-button>
          </t-form-item>
        </t-form>
        
        <div class="register-link">
          想起密码了？
          <t-link theme="primary" hover="color" @click="goToLogin">返回登录</t-link>
        </div>
      </div>
    </div>
    
    <!-- 底部波浪 -->
    <div class="wave-decoration"></div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import type { FormInstanceFunctions, FormRule, SubmitContext } from 'tdesign-vue-next'
import { resetPassword, sendEmailCode } from '@/api'

defineOptions({
  name: 'UserForget',
})

const router = useRouter()
const formRef = ref<FormInstanceFunctions>()
const loading = ref(false)
const countdown = ref(0)

const formData = reactive({
  email: '',
  verifyCode: '',
  password: '',
  confirmPassword: '',
})

const rules: Record<string, FormRule[]> = {
  email: [
    { required: true, message: '请输入邮箱', type: 'error' },
    { email: true, message: '请输入正确的邮箱格式', type: 'error' },
  ],
  verifyCode: [
    { required: true, message: '请输入验证码', type: 'error' },
    { len: 6, message: '验证码为6位', type: 'error' },
  ],
  password: [
    { required: true, message: '请输入新密码', type: 'error' },
    { min: 6, message: '密码长度不能少于6位', type: 'error' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', type: 'error' },
    {
      validator: (val) => val === formData.password,
      message: '两次输入的密码不一致',
      type: 'error',
    },
  ],
}

const sendVerifyCode = async () => {
  const result = await formRef.value?.validate({ fields: ['email'] })
  if (result !== true) {
    return
  }
  
  try {
    // 调用发送验证码接口
    await sendEmailCode({ email: formData.email })
    
    MessagePlugin.success('验证码已发送到您的邮箱')
    countdown.value = 60
    
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (error: any) {
    console.error('发送验证码失败:', error)
    MessagePlugin.error(error.message || '发送验证码失败，请稍后重试')
  }
}

const onSubmit = async (ctx: SubmitContext) => {
  if (ctx.validateResult !== true) {
    return
  }
  
  loading.value = true
  try {
    // 调用重置密码接口
    await resetPassword({
      email: formData.email,
      verifyCode: formData.verifyCode,
      password: formData.password,
    })
    
    MessagePlugin.success('密码重置成功，请登录')
    router.push('/user/login')
  } catch (error: any) {
    console.error('密码重置失败:', error)
    MessagePlugin.error(error.message || '密码重置失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

const goToLogin = () => {
  router.push('/user/login')
}
</script>

<style scoped lang="less">
@import './user.less';

.verify-code-wrapper {
  display: flex;
  gap: 12px;
  width: 100%;
  
  :deep(.t-input) {
    flex: 1;
    min-width: 0;
  }
  
  :deep(.t-button) {
    flex-shrink: 0;
    width: 130px;
    white-space: nowrap;
  }
}
</style>
