<template>
  <PersonPanel
    v-model="defaultModelValue"
    :allow-multiple="formData.settings.allowMultiple"
    @handle-blur="handleCellEdit"
    @handle-enter="handleCellEdit"
    @handle-esc="handleKeyDown"
  />
</template>

<script setup lang="ts">
import PersonPanel from '@/components/personPanel/index.vue'
import { ObjType } from '@/types'

interface Props {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  editingValue: any
  field: ObjType
}
const props = defineProps<Props>()
const formData = computed(() => props.field)
const emit = defineEmits(['update:editingValue', 'handleCellEdit', 'handleCancelCellEdit'])

const defaultModelValue = computed({
  get() {
    return props.field.settings.allowMultiple ? props.editingValue : props.editingValue[0] || '' // 多选返回原数组，单选则从数组中取第一个
  },
  set(value) {
    emit('update:editingValue', value)
  },
})

// 处理单元格编辑
const handleCellEdit = () => {
  emit('handleCellEdit')
}

// 处理按键esc
const handleKeyDown = () => {
  emit('handleCancelCellEdit')
}

defineOptions({ name: 'PersonComponent' })
</script>

<style scoped></style>
