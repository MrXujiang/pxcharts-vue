<template>
  <t-tooltip content="系统自动填充，不可修改">
    <div class="text">{{ defaultValue }}</div>
  </t-tooltip>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'
import { customNumberFormat } from '@/utils/mtable'

interface Props {
  field: ObjType
  item: ObjType
}
const props = defineProps<Props>()
const defaultValue = computed(() =>
  props.field.settings.idType == 'autoIncrement'
    ? props.item[props.field.id]
    : customNumberFormat(props.item[props.field.id], props.field),
)
defineOptions({ name: 'AutoNumberComponent' })
</script>

<style lang="less" scoped>
.text {
  font-size: 14px;
}
</style>
