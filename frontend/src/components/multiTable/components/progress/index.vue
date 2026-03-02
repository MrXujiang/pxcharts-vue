<template>
  <t-input-number
    ref="inputRef"
    v-model="defaultModelValue"
    theme="normal"
    class="w-100"
    :decimal-places="formData.settings.decimalPlaces"
    @blur="handleCellEdit"
    @enter="handleCellEdit"
    @keydown="handleKeyDown"
  />
</template>

<script setup lang="ts">
import { ObjType } from '@/types'
import { InputValue } from 'tdesign-vue-next'

interface Props {
  editingValue: number
  field: ObjType
}
const props = defineProps<Props>()
const formData = computed(() => props.field)
const emit = defineEmits(['handleCellEdit', 'handleCancelCellEdit', 'update:editingValue'])
const inputRef = ref()

const defaultModelValue = computed({
  get() {
    return props.editingValue
  },
  set(value) {
    emit('update:editingValue', value)
  },
})
const handleKeyDown = (value: InputValue, context: { e: KeyboardEvent }) => {
  if (context.e.code === 'Escape') {
    console.log('Cell edit cancelled')
    emit('handleCancelCellEdit')
  }
}

const handleCellEdit = () => {
  console.log('Cell edited', defaultModelValue.value)
  emit('handleCellEdit')
}

onMounted(() => {
  nextTick(() => {
    inputRef.value.focus()
  })
})

defineOptions({ name: 'ProgressPreviewComponent' })
</script>

<style scoped></style>
