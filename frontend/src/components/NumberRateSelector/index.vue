<template>
  <div class="number-rate-selector">
    <div
      v-for="num in numbers"
      :key="num"
      class="rate-item"
      :class="{
        active: isActive(num),
        'active-before': isActiveBefore(num),
        hover: isHover(num),
        'hover-before': isHoverBefore(num),
        disabled: isDisabled(num),
      }"
      @click="handleSelect(num)"
      @mouseenter="handleMouseEnter(num)"
      @mouseleave="handleMouseLeave"
    >
      {{ num }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

defineOptions({ name: 'NumberRateSelector' })

interface Props {
  modelValue?: number | number[] | null // 选中的值（单选或多选）
  min?: number // 最小值
  max?: number // 最大值
  length?: number // 数字个数
  multiple?: boolean // 是否多选
  disabled?: boolean // 是否禁用
  disabledValues?: number[] // 禁用的特定值
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: undefined,
  min: 1,
  max: undefined,
  length: 5,
  multiple: false,
  disabled: false,
  disabledValues: () => [],
})

interface Emits {
  (e: 'update:modelValue', value: number | number[]): void
  (e: 'change', value: number | number[]): void
}

const emit = defineEmits<Emits>()

// 记录当前hover的数字
const hoveredNum = ref<number | null>(null)

// 生成数字数组
const numbers = computed(() => {
  if (props.max !== undefined) {
    return Array.from({ length: props.max - props.min + 1 }, (_, i) => props.min + i)
  } else {
    // 否则使用length生成从min开始的数字
    return Array.from({ length: props.length }, (_, i) => props.min + i)
  }
})

// 判断是否选中
const isActive = (num: number) => {
  if (props.multiple) {
    return Array.isArray(props.modelValue) && props.modelValue.includes(num)
  } else {
    return props.modelValue === num
  }
}

// 判断是否在选中数字之前（包括选中数字本身）
const isActiveBefore = (num: number) => {
  if (props.multiple) {
    // 多选模式下不应用此效果
    return false
  } else {
    // 确保在单选模式下modelValue是数字类型
    return typeof props.modelValue === 'number' && num <= props.modelValue
  }
}

// 判断是否禁用
const isDisabled = (num: number) => {
  return props.disabled || props.disabledValues.includes(num)
}

// 判断是否hover
const isHover = (num: number) => {
  return hoveredNum.value === num
}

// 判断是否在hover数字之前（包括hover数字本身）
const isHoverBefore = (num: number) => {
  return hoveredNum.value !== null && num <= hoveredNum.value
}

// 处理鼠标进入
const handleMouseEnter = (num: number) => {
  hoveredNum.value = num
}

// 处理鼠标离开
const handleMouseLeave = () => {
  hoveredNum.value = null
}

// 处理选择
const handleSelect = (num: number) => {
  if (isDisabled(num)) return

  let newValue: number | number[]

  if (props.multiple) {
    // 多选模式
    const currentValue = Array.isArray(props.modelValue) ? props.modelValue : []
    if (currentValue.includes(num)) {
      newValue = currentValue.filter((v) => v !== num)
    } else {
      newValue = [...currentValue, num]
    }
  } else {
    // 单选模式
    newValue = num
  }
  console.log('newValue', newValue)
  emit('update:modelValue', newValue)
  emit('change', newValue)
}
</script>

<style scoped lang="less">
.number-rate-selector {
  display: inline-flex;
  gap: 4px;

  .rate-item {
    min-width: 16px;
    height: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #e8e8e8;
    border-radius: 2px;
    font-size: 10px;
    font-weight: 500;
    color: #999;
    cursor: pointer;
    transition: all 0.3s;
    user-select: none;

    // hover效果：当前及之前元素
    &.hover-before:not(.disabled) {
      background: rgb(123, 176, 236);
      color: #fff;
    }

    &.active {
      background: #4a7ff7;
      color: #ffffff;

      &:hover:not(.disabled) {
        background: #3a6fe7;
      }
    }

    // 激活效果：当前及之前元素
    &.active-before:not(.disabled) {
      background: #4a7ff7;
      color: #ffffff;
    }

    &.disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  }
}
</style>
