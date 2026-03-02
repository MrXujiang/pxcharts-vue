<template>
  <t-date-picker
    ref="datePickerRef"
    v-model="defaultModelValue"
    class="w-100"
    :readonly="false"
    :enable-time-picker="timePicker.includes(formData.format)"
    :format="formData.format"
    @change="handleCellEdit"
    @blur="handleCellEdit"
  />
</template>

<script setup lang="ts">
import { ObjType } from '@/types'

interface Props {
  editingValue: string
  field: ObjType
}

const props = defineProps<Props>()
const emit = defineEmits(['update:editingValue', 'handleCellEdit', 'handleCancelCellEdit'])
const datePickerRef = ref()

const timePicker = [
  'YYYY-MM-DD HH:mm',
  'YYYY-MM-DD HH:mm:ss',
  'YYYY/MM/DD HH:mm',
  'YYYY/MM/DD HH:mm:ss',
]
const formData = computed(() => props.field)
const defaultModelValue = computed({
  get() {
    return props.editingValue
  },
  set(value) {
    emit('update:editingValue', value)
  },
})

const handleCellEdit = () => {
  emit('handleCellEdit')
}

onMounted(() => {
  // 组件挂载后，自动聚焦到日期选择器以打开面板
  nextTick(() => {
    if (datePickerRef.value) {
      // 尝试让日期选择器自动获得焦点，这通常会打开面板
      setTimeout(() => {
        const inputElement = datePickerRef.value.$el.querySelector('input')
        if (inputElement) {
          inputElement.click()
        }
      }, 100)
    }
  })
})
defineOptions({ name: 'DateComponent' })
</script>

<style scoped></style>
