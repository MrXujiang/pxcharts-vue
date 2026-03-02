<template>
  <t-tooltip content="系统自动填充，不可修改">
    <div class="time-text">{{ defaultModelValue || '' }}</div>
  </t-tooltip>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'
import dayjs from 'dayjs'

interface Props {
  field: ObjType
  item: ObjType
  editingValue: string | number
}

const props = defineProps<Props>()
const defaultModelValue = computed(() => dayjs(props.item.createdAt).format(props.field.format))

defineOptions({ name: 'CreateTimeComponent' })
</script>

<style lang="less" scoped>
.time-text {
  font-size: 14px;
}
</style>
