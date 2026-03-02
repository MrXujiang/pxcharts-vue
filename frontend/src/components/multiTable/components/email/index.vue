<template>
  <t-input
    ref="inputRef"
    v-model="defaultModelValue"
    @enter="handleCellEdit"
    @blur="handleCellEdit"
    @keydown="handleKeyDown"
  />
</template>

<script setup lang="ts">
import { InputValue } from 'tdesign-vue-next'

interface Props {
  editingValue: string | number
}

const props = defineProps<Props>()
const emit = defineEmits(['handleCellEdit', 'handleCancelCellEdit', 'update:editingValue'])
const inputRef = ref()

const defaultModelValue = computed({
  get() {
    return props.editingValue
  },
  set(value) {
    emit('update:editingValue', value)
  },
})

const handleKeyDown = (value: InputValue, context: { e: KeyboardEvent }) => {
  if (context.e.code === 'Escape') {
    console.log('Cell edit cancelled')
    emit('handleCancelCellEdit')
  }
}
// 邮箱格式校验正则表达式
const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

// 校验邮箱格式
const validateEmail = (email: string): boolean => {
  return emailRegex.test(email)
}
const handleCellEdit = () => {
  if (defaultModelValue.value && !validateEmail(defaultModelValue.value as string)) {
    MessagePlugin.warning('邮箱格式不正确')
    return
  }
  emit('handleCellEdit')
}

onMounted(() => {
  nextTick(() => {
    inputRef.value.focus()
  })
})
defineOptions({ name: 'EmailComponent' })
</script>

<style scoped></style>
