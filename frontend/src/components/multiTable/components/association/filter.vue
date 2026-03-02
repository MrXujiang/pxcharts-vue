<template>
  <t-select
    v-model="defaultModelValue"
    :options="optionsData"
    :min-collapsed-num="1"
    placeholder="请选择"
    :keys="{ label: 'value', value: 'id' }"
    multiple
  />
</template>

<script setup lang="ts">
import { getRecordOptions } from '@/api'
import { ObjType } from '@/types'

interface Props {
  editingValue: string | number
  field: ObjType
}
const props = defineProps<Props>()
const emit = defineEmits(['handleCellEdit', 'handleCancelCellEdit', 'update:editingValue'])

console.log('props', props)

const defaultModelValue = computed({
  get() {
    return props.editingValue
  },
  set(value) {
    emit('update:editingValue', value)
  },
})

const optionsData = ref([])

// 获取当前引用表选项数据
const handleGetOptionsData = async () => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const res: any = await getRecordOptions({
    tableSchemaId: props.field.settings.sourceRef.tableSchemaId,
  })
  optionsData.value = res.options || []
}

onMounted(() => {
  handleGetOptionsData()
})

defineOptions({ name: 'AssociationFilter' })
</script>

<style scoped></style>
