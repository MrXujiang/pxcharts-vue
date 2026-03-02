<template>
  <div class="date-filter">
    <t-select
      v-model="defaultModelValue"
      placeholder="请选择日期范围"
      @change="handleDateChange"
      :options="dateOptions"
      clearable
    />
  </div>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'
import { Select as TSelect, type SelectValue } from 'tdesign-vue-next'

interface Props {
  editingValue: string
  field: ObjType
}

const props = defineProps<Props>()
// 定义组件事件
const emit = defineEmits(['update:editingValue'])

const defaultModelValue = computed({
  get() {
    return props.editingValue
  },
  set(value) {
    emit('update:editingValue', value)
  },
})

// 日期选项计算
const dateOptions = [
  { label: '今天', value: 'today' },
  { label: '明天', value: 'tomorrow' },
  { label: '昨天', value: 'yesterday' },
  { label: '本周', value: 'thisWeek' },
  { label: '上周', value: 'lastWeek' },
  { label: '下周', value: 'nextWeek' },
  { label: '本月', value: 'thisMonth' },
  { label: '上个月', value: 'lastMonth' },
  { label: '下个月', value: 'nextMonth' },
  { label: '今年', value: 'thisYear' },
  { label: '去年', value: 'lastYear' },
  { label: '明年', value: 'nextYear' },
]

// 处理日期选择变化
const handleDateChange = (value: SelectValue) => {
  console.log('value00000', value)
  defaultModelValue.value = value as string
}

defineOptions({
  name: 'DateFilter',
})
</script>

<style scoped>
.date-filter {
  width: 100%;
}
</style>
