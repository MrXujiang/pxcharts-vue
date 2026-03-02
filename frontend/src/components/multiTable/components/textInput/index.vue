<template>
  <div>
    <t-input
      ref="inputRef"
      v-model="defaultModelValue"
      @enter="handleCellEdit"
      @blur="handleCellEdit"
      @keydown="handleKeyDown"
    />
  </div>
</template>

<script setup lang="ts">
import { InputValue } from 'tdesign-vue-next'

interface Props {
  editingValue: string | number
}

const props = defineProps<Props>()
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
defineOptions({ name: 'TextInputComponent' })
</script>

<style scoped></style>
