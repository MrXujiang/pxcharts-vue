<template>
  <div>
    <t-rate
      :model-value="innerValue ?? undefined"
      @update:model-value="handleTRateChange"
      v-if="iconType !== 'number'"
      size="16px"
      :key="iconType"
      :default-value="0"
      :count="count"
      :color="iconType === 'heart' ? ['#FF0000', '#E3E6EB'] : ['#ED7B2F', '#E3E6EB']"
      :disabled="disabled"
    >
      <template #icon="iconProps">
        <t-icon v-bind="iconProps" :name="iconName" />
      </template>
    </t-rate>
    <NumberRateSelector
      v-else
      :model-value="innerValue"
      @update:model-value="handleNumberRateChange"
      :min="min"
      :max="max"
      :disabled="disabled"
    />
  </div>
</template>

<script setup lang="ts">
import NumberRateSelector from '@/components/NumberRateSelector/index.vue'

defineOptions({
  name: 'IconRate',
})

interface Props {
  modelValue?: number | null
  iconType: string
  iconName: string
  count: number
  min?: number
  max?: number
  disabled?: boolean // 是否禁用
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: null,
  min: 1,
  max: 5,
  disabled: false, // 是否禁用评分，默认为false不禁用
})

interface Emits {
  (e: 'update:modelValue', value: number | null): void
}

const emit = defineEmits<Emits>()

// 内部值，用于v-model
const innerValue = ref<number | null>(props.modelValue)

// 监听props变化，同步内部值
watch(
  () => props.modelValue,
  (newValue) => {
    innerValue.value = newValue
  },
)

// 处理t-rate组件值变化
const handleTRateChange = (value: number) => {
  innerValue.value = value
  emit('update:modelValue', value)
}

// 处理NumberRateSelector组件值变化
const handleNumberRateChange = (value: number | number[]) => {
  // 只取第一个值或直接使用数值
  const numericValue = Array.isArray(value) ? value[0] : value
  innerValue.value = numericValue
  emit('update:modelValue', numericValue)
}
</script>
