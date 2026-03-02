<template>
  <div>
    <span v-if="value == '未分组'" class="no-content">无内容</span>
    <span v-else class="date-text">{{ dayjs(value).format(dataFormat) }}</span>
  </div>
</template>

<script setup lang="ts">
import dayjs from 'dayjs'
import { useMtTableStore } from '@/stores/mtTable'

interface Props {
  value: string
}

const mtTableStore = useMtTableStore()
const props = defineProps<Props>()
console.log('props', props)

const field = computed(() => mtTableStore.getFields())
const groupByField = computed(() =>
  mtTableStore.currentTable.settings.tableConfig.groupConfig?.length > 0
    ? mtTableStore.currentTable.settings.tableConfig.groupConfig[0]
    : '',
)
const dataFormat = computed(() => {
  // 根据groupByField 在field中找到id相等的 format 值
  return field.value.find((item) => item.id === groupByField.value)?.format
})
console.log('dataFormat', dataFormat.value)
defineOptions({ name: 'GroupHeadDate' })
</script>

<style scoped>
.date-text {
  color: #999;
  font-size: 12px;
}
</style>
