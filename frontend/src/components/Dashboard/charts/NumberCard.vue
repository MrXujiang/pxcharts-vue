<template>
  <div class="number-card" :style="{ background: cardBackground }">
    <div class="card-header">
      <span class="card-title">{{ config.title || '指标名称' }}</span>
    </div>
    <div class="card-content">
      <div class="number-value">
        <span class="value" :style="{ color: valueColor }">{{ formattedValue }}</span>
        <span v-if="data.trend" class="trend-icon" :class="data.trend">
          <t-icon :name="data.trend === 'up' ? 'caret-up' : 'caret-down'" />
        </span>
      </div>
      <div v-if="data.compare" class="compare-text">
        环比 {{ data.compare }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  config: {
    title?: string
    bgColor?: string
    unit?: string
  }
  data: {
    value: number
    trend?: 'up' | 'down'
    compare?: string
  }
  themeColors?: string[]
}

const props = withDefaults(defineProps<Props>(), {
  config: () => ({}),
  data: () => ({ value: 0 }),
  themeColors: () => ['#3b82f6', '#8b5cf6', '#22c55e', '#eab308', '#f97316', '#ec4899'],
})

const formattedValue = computed(() => {
  const { value } = props.data
  const { unit = '' } = props.config
  
  // 格式化大数字
  if (value >= 10000) {
    return `${(value / 10000).toFixed(1)}万${unit}`
  }
  return `${value}${unit}`
})

// 使用主题颜色作为数值颜色
const valueColor = computed(() => {
  return props.themeColors?.[0] || '#3b82f6'
})

// 使用主题颜色的淡化版本作为背景
const cardBackground = computed(() => {
  if (props.config.bgColor) {
    return props.config.bgColor
  }
  // 将主题颜色转换为淡化版本
  const color = props.themeColors?.[0] || '#3b82f6'
  // 使用淡化的渐变背景
  return `linear-gradient(135deg, ${color}10 0%, ${color}20 100%)`
})
</script>

<style scoped lang="less">
.number-card {
  width: 100%;
  height: 100%;
  padding: 20px;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;

  .card-header {
    .card-title {
      font-size: 14px;
      color: #666;
      font-weight: 400;
    }
  }

  .card-content {
    .number-value {
      display: flex;
      align-items: baseline;
      gap: 8px;
      margin-bottom: 8px;

      .value {
        font-size: 36px;
        font-weight: 600;
        line-height: 1;
        transition: color 0.3s;
      }

      .trend-icon {
        font-size: 20px;
        
        &.up {
          color: #52c41a;
        }
        
        &.down {
          color: #f5222d;
        }
      }
    }

    .compare-text {
      font-size: 12px;
      color: #999;
    }
  }
}
</style>
