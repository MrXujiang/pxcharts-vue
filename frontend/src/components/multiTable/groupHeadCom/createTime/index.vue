<template>
  <div class="time-text">{{ defaultModelValue !== '未分组' ? defaultModelValue : '' }}</div>
</template>

<script setup lang="ts">
import { useMtTableStore } from '@/stores/mtTable'
import dayjs from 'dayjs'

interface Props {
  value: string
}

const props = defineProps<Props>()
const mtTableStore = useMtTableStore()
const groupByField = computed(() =>
  mtTableStore.currentTable.settings.tableConfig.groupConfig?.length > 0
    ? mtTableStore.currentTable.settings.tableConfig.groupConfig[0]
    : '',
)

const fieldAll = computed(() => mtTableStore.getFields())

const field = computed(() => fieldAll.value.find((item) => item.id === groupByField.value))

const defaultModelValue = computed(() => dayjs(props.value).format(field.value?.format))

defineOptions({
  name: 'CreateTimeGroupHeadCom',
})
</script>

<style lang="less" scoped>
.time-text {
  font-size: 12px;
  color: #404258;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
