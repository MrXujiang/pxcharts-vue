<template>
  <t-input
    ref="inputRef"
    v-model="defaultModelValue"
    placeholder="请输入数字"
    style="width: 100%"
    @blur="handleBlur"
    @focus="handleFocus"
    @change="handleInputChange"
    @keydown="handleKeyDown"
  />
</template>

<script setup lang="ts">
import { ObjType } from '@/types'
import { NumberFormatOptions, formatNumber, parseDisplayToRaw } from '@/utils/numberFormatter'
import { InputProps, InputValue } from 'tdesign-vue-next'

interface Props {
  editingValue: string | number
  field: ObjType
}

const props = defineProps<Props>()
const emit = defineEmits(['update:editingValue', 'handleCellEdit', 'handleCancelCellEdit'])
const inputRef = ref()
const formData = computed(() => props.field)
const getFormatOptions = (): NumberFormatOptions => ({
  currency: formData.value.settings.currency,
  decimalPlaces: formData.value.settings.decimalPlaces as number,
  useThousandSeparator: formData.value.settings.useThousandSeparator,
  thousandSeparator: formData.value.settings.thousandSeparator,
  largeNumberAbbreviation: formData.value.settings.largeNumberAbbreviation,
  disallowNegative: formData.value.settings.disallowNegative,
})
const rawValue = ref(parseDisplayToRaw(props.editingValue as string, getFormatOptions())) // 真实数据
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
    emit('handleCancelCellEdit')
  }
}

// 输入框内容变化时处理
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const handleInputChange: InputProps['onChange'] = (val: any) => {
  console.log('val', val)
  const reg = /^-?[\d.,%\s]*$/
  if (!reg.test(val)) {
    MessagePlugin.warning('仅支持输入数字、小数点、逗号、负号和百分号')
    defaultModelValue.value =
      rawValue.value !== null ? formatNumber(rawValue.value, getFormatOptions()) : ''
    return
  }

  const parsed = parseDisplayToRaw(val as string, getFormatOptions())
  if (parsed !== null) {
    rawValue.value = parsed
    defaultModelValue.value = val as string
  } else {
    defaultModelValue.value = val as string
  }
}
// 输入框聚焦时处理
const handleFocus = () => {
  defaultModelValue.value = rawValue.value as number
}
// 输入框失焦时处理
function handleBlur() {
  if (rawValue.value !== null) {
    defaultModelValue.value = formatNumber(rawValue.value, getFormatOptions())
  } else {
    defaultModelValue.value = ''
  }
  emit('handleCellEdit')
}

onMounted(() => {
  nextTick(() => {
    inputRef.value.focus()
  })
})
defineOptions({ name: 'CurrencyComponent' })
</script>

<style scoped></style>
