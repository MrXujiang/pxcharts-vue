<template>
  <div class="w-100">
    <t-space size="small" direction="vertical" class="w-100">
      <div class="flx-ce-bet">
        <div>邮箱地址</div>
      </div>
      <t-input
        v-model="formData.defaultValue"
        :status="inputStatus"
        :tips="tips"
        placeholder="请输入邮箱地址"
        @blur="handleBlur"
      />
    </t-space>
  </div>
</template>

<script setup lang="ts">
import { InputProps } from 'tdesign-vue-next'

const props = defineProps({
  formData: {
    type: Object,
    required: true,
  },
})
const { formData } = toRefs(props)
defineOptions({ name: 'EmailComponent' })

// 邮箱格式校验正则表达式
const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

// 校验邮箱格式
const validateEmail = (email: string): boolean => {
  return emailRegex.test(email)
}

const inputStatus = computed<InputProps['status']>(() => {
  if (formData.value.defaultValue && !validateEmail(formData.value.defaultValue)) {
    return 'error'
  }
  return 'default'
})

const tips = computed<InputProps['tips']>(() => {
  if (formData.value.defaultValue && !validateEmail(formData.value.defaultValue)) {
    return '邮箱格式不正确'
  }
  return ''
})

// 失去焦点时校验，如果格式不正确则清空内容
const handleBlur = () => {
  if (formData.value.defaultValue && !validateEmail(formData.value.defaultValue)) {
    // 格式不正确，置空内容
    formData.value.defaultValue = ''
  }
}
</script>

<style lang="less" scoped></style>
