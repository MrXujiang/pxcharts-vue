<template>
  <div>
    <RateEditComponent :field="field" :item="null" :defaultValue="defaultValue" :isEditor="false" />
  </div>
</template>

<script setup lang="ts">
import RateEditComponent from '@/components/multiTable/components/rate/index.vue'
import { useMtTableStore } from '@/stores/mtTable'

interface Props {
  value: string
}

const props = defineProps<Props>()

const defaultValue = computed(() => props.value)

const mtTableStore = useMtTableStore()
const groupByField = computed(() =>
  mtTableStore.currentTable.settings.tableConfig.groupConfig?.length > 0
    ? mtTableStore.currentTable.settings.tableConfig.groupConfig[0]
    : '',
)
const fieldAll = computed(() => mtTableStore.getFields())

const field = computed(() => fieldAll.value.find((item) => item.id === groupByField.value))

defineOptions({
  name: 'RateGroupHeadCom',
})
</script>

<style scoped></style>
