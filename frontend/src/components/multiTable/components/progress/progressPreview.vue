<template>
  <div class="progress-container">
    <t-progress
      :percentage="defaultModelValue"
      :color="colorOptions.find((option) => option.value === formData.settings.color)?.colors"
    >
      <template #label>
        <span
          >{{ formatNumberByDecimal(defaultModelValue, formData.settings.decimalPlaces)
          }}{{ formData.settings.numberFormat == 'percentage' ? '%' : '' }}</span
        >
      </template>
    </t-progress>
  </div>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'
import { colorOptions } from '@/views/product/components/fieldPop/addFieldModal/component/progress/options'

interface Props {
  editingValue: number
  field: ObjType
}
const props = defineProps<Props>()
const defaultModelValue = computed(() => props.editingValue)
const formData = computed(() => props.field)
/**
 * 根据小数位数格式化数字（仅截断，不四舍五入）
 * @param {number|string} num - 待格式化的数字（支持数字或字符串类型）
 * @param {number} decimal - 小数位数（对应 decimalOptions 的 value）
 * @returns {string} 格式化后的结果
 */
function formatNumberByDecimal(num, decimal) {
  // 先将输入转为数字类型（处理字符串输入）
  const number = Number(num)

  // 校验数字有效性
  if (isNaN(number)) {
    console.warn('输入不是有效的数字')
    return ''
  }

  // 核心逻辑：仅截断，不四舍五入
  const multiplier = Math.pow(10, decimal) // 10的n次方（n为小数位数）
  const truncatedNum = Math.trunc(number * multiplier) / multiplier // 截断处理

  // 补零并转为字符串（确保小数位数符合要求）
  const formatted = truncatedNum.toFixed(decimal)

  // 整数优化：去除末尾的 .0（可选，可按需删除）
  return decimal === 0 ? formatted.replace('.0', '') : formatted
}

defineOptions({ name: 'ProgressComponent' })
</script>

<style scoped>
.progress-container {
  width: 100%;
  padding: 0 10px;
}
</style>
