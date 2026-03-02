<template>
  <t-input
    ref="inputRef"
    v-model="defaultModelValue"
    @enter="handleCellEdit"
    @blur="handleCellEdit"
    @focus="handleFocus"
    @change="handleInputChange"
    @keydown="handleKeyDown"
  />
</template>

<script setup lang="ts">
import { ObjType } from '@/types'
import { InputProps, InputValue } from 'tdesign-vue-next'
import { NumberFormatOptions, formatNumber, parseDisplayToRaw } from '@/utils/numberFormatter'

interface Props {
  editingValue: string | number
  field: ObjType
}

const props = defineProps<Props>()
const emit = defineEmits(['update:editingValue', 'handleCellEdit', 'handleCancelCellEdit'])
const inputRef = ref()
const formData = computed(() => props.field)

const getFormatOptions = (): NumberFormatOptions => ({
  displayFormat: formData.value.settings.displayFormat as 'integer' | 'decimal' | 'percentage', // 显示格式
  decimalPlaces: formData.value.settings.decimalPlaces as number, // 小数点位数
  useThousandSeparator: formData.value.settings.useThousandSeparator, // 是否使用千位分隔符
  thousandSeparator: formData.value.settings.thousandSeparator as 'comma' | 'dot', // 千分位分隔符
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

// ===== 输入处理 =====
const handleInputChange: InputProps['onChange'] = (val) => {
  if (!/^-?[\d.,%\s]*$/.test(val as string)) {
    MessagePlugin.warning('仅支持输入数字、小数点、逗号、负号和百分号')
    defaultModelValue.value =
      rawValue.value !== null ? formatNumber(rawValue.value as number, getFormatOptions()) : ''
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
  defaultModelValue.value = parseDisplayToRaw(props.editingValue as string, getFormatOptions()) as
    | string
    | number
}
// 失焦时格式化
function handleCellEdit() {
  if (rawValue.value !== null) {
    defaultModelValue.value = formatNumber(rawValue.value as number, getFormatOptions())
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

defineOptions({ name: 'NumberComponent' })
</script>

<style scoped></style>
