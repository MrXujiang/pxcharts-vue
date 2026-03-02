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
// 身份证格式校验正则表达式
const idCardRegex = /^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[0-9Xx]$/

// 校验邮箱格式
const validateIdCard = (idCard: string): boolean => {
  return idCardRegex.test(idCard)
}
const handleCellEdit = () => {
  if (defaultModelValue.value && !validateIdCard(defaultModelValue.value as string)) {
    defaultModelValue.value = ''
    MessagePlugin.warning('身份证格式不正确')
    return
  }
  emit('handleCellEdit')
}

onMounted(() => {
  nextTick(() => {
    inputRef.value.focus()
  })
})
defineOptions({ name: 'IdCardComponent' })
</script>

<style scoped></style>
