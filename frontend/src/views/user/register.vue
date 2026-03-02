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
          <h1 class="title">迈维表格<span class="highlight">注册</span></h1>
          <p class="subtitle">——一张多维表，装下所有业务</p>
        </div>

        <!-- Tab 切换 -->
        <div class="register-tabs">
          <div
            class="tab-item"
            :class="{ active: registerType === 'email' }"
            @click="switchRegisterType('email')"
          >
            邮箱注册
          </div>
          <div
            class="tab-item"
            :class="{ active: registerType === 'invite' }"
            @click="switchRegisterType('invite')"
          >
            邀请码注册
          </div>
        </div>

        <t-form
          ref="formRef"
          :data="formData"
          :rules="rules"
          :label-width="0"
          class="user-form"
          @submit="onSubmit"
        >
          <!-- 邮箱注册 -->
          <t-form-item v-if="registerType === 'email'" name="email">
            <t-input
              v-model="formData.email"
              placeholder="请输入邮箱"
              size="large"
              clearable
            >
              <template #prefix-icon>
                <MailIcon style="color: #909399" />
              </template>
            </t-input>
          </t-form-item>

          <t-form-item v-if="registerType === 'email'" name="verifyCode">
            <div class="verify-code-wrapper">
              <t-input
                v-model="formData.verifyCode"
                placeholder="请输入验证码"
                size="large"
                clearable
              >
                <template #prefix-icon>
                  <MobileIcon style="color: #909399" />
                </template>
              </t-input>
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

          <!-- 邀请码注册 -->
          <t-form-item v-if="registerType === 'invite'" name="email">
            <t-input
              v-model="formData.email"
              placeholder="请输入邮箱"
              size="large"
              clearable
            >
              <template #prefix-icon>
                <MailIcon style="color: #909399" />
              </template>
            </t-input>
          </t-form-item>

          <t-form-item v-if="registerType === 'invite'" name="inviteCode">
            <t-input
              v-model="formData.inviteCode"
              placeholder="请输入邀请码"
              size="large"
              clearable
            >
              <template #prefix-icon>
                <GiftIcon style="color: #909399" />
              </template>
            </t-input>
          </t-form-item>

          <t-form-item name="password">
            <t-input
              v-model="formData.password"
              type="password"
              placeholder="请输入密码"
              size="large"
              clearable
            >
              <template #prefix-icon>
                <LockOnIcon style="color: #909399" />
              </template>
            </t-input>
          </t-form-item>

          <t-form-item name="confirmPassword">
            <t-input
              v-model="formData.confirmPassword"
              type="password"
              placeholder="请确认密码"
              size="large"
              clearable
            >
              <template #prefix-icon>
                <LockOnIcon style="color: #909399" />
              </template>
            </t-input>
          </t-form-item>

          <t-form-item>
            <t-button
              theme="primary"
              type="submit"
              block
              size="large"
              :loading="loading"
            >
              注册
            </t-button>
          </t-form-item>
        </t-form>

        <div class="register-link">
          已有账号？
          <t-link theme="primary" hover="color" @click="goToLogin">立即登录</t-link>
        </div>

        <div class="agreement-text">
          注册即代表您已阅读并同意
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
import { reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import type { FormInstanceFunctions, FormRule, SubmitContext } from 'tdesign-vue-next'
import { MailIcon, GiftIcon, LockOnIcon, MobileIcon } from 'tdesign-icons-vue-next'
import AgreementDialog from './components/AgreementDialog.vue'
import { registerByEmailCode, registerByInvitationCode, sendEmailCode } from '@/api'
import { useUserStore } from '@/stores/user'

defineOptions({
  name: 'UserRegister',
})

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstanceFunctions>()
const loading = ref(false)
const countdown = ref(0)
const registerType = ref<'email' | 'invite'>('email')
const agreementVisible = ref(false)
const agreementType = ref<'user' | 'privacy'>('user')

const formData = reactive({
  email: '',
  verifyCode: '',
  inviteCode: '',
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
  inviteCode: [
    { required: true, message: '请输入邀请码', type: 'error' },
  ],
  password: [
    { required: true, message: '请输入密码', type: 'error' },
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

const switchRegisterType = (type: 'email' | 'invite') => {
  registerType.value = type
  // 清空表单数据
  formData.email = ''
  formData.verifyCode = ''
  formData.inviteCode = ''
  formData.password = ''
  formData.confirmPassword = ''
  countdown.value = 0
  // 清空校验状态
  formRef.value?.clearValidate()
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
    let response

    if (registerType.value === 'email') {
      // 邮箱验证码注册
      response = await registerByEmailCode({
        email: formData.email,
        verifyCode: formData.verifyCode,
        password: formData.password,
      })
    } else {
      // 邀请码注册
      response = await registerByInvitationCode({
        email: formData.email,
        inviteCode: formData.inviteCode,
        password: formData.password,
      })
    }

    // 假设接口返回格式: { code: 0, data: { token: string, user: {...} } }
    if (response) {
      const { token, userInfo: user } = response

      // 保存用户信息
      userStore.login(token, user)

      MessagePlugin.success('注册成功')
      // 跳转到标签选择页面
      router.push('/user/entry')
    }
  } catch (error: any) {
    console.error('注册失败:', error)
    MessagePlugin.error(error.message || '注册失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

const goToLogin = () => {
  router.push('/user/login')
}

const showAgreement = (type: 'user' | 'privacy') => {
  agreementType.value = type
  agreementVisible.value = true
}
</script>

<style scoped lang="less">
@import './user.less';

.register-tabs {
  display: flex;
  gap: 0;
  margin-bottom: 24px;
  border-bottom: 2px solid #f0f0f0;

  .tab-item {
    flex: 1;
    text-align: center;
    padding: 12px 0;
    font-size: 15px;
    color: #666;
    cursor: pointer;
    position: relative;
    transition: all 0.3s;

    &:hover {
      color: #4a7ff7;
    }

    &.active {
      color: #4a7ff7;
      font-weight: 500;

      &::after {
        content: '';
        position: absolute;
        bottom: -2px;
        left: 0;
        right: 0;
        height: 2px;
        background: #4a7ff7;
      }
    }
  }
}

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
