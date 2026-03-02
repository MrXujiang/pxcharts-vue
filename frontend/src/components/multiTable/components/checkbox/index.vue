<template>
  <div class="w-100 flx-ce-ce">
    <t-checkbox :disabled="!isEditor" v-model="checkVal" @change="handleCheckboxChange" />
  </div>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'
import { CheckboxProps } from 'tdesign-vue-next'

interface Props {
  field: ObjType
  item: ObjType
  isEditor: boolean // 是否处于编辑状态
}

const props = defineProps<Props>()
const checkVal = ref<boolean>(false)
const emit = defineEmits(['handleSpecialChange'])

const handleCheckboxChange: CheckboxProps['onChange'] = (val) => {
  checkVal.value = val
  emit('handleSpecialChange', {
    field: props.field,
    item: props.item,
    value: val,
  })
}
watch(
  () => props.item,
  (val) => {
    if (val) {
      checkVal.value = props.item[props.field.id]
    }
  },
  { immediate: true },
)
defineOptions({ name: 'CheckboxComponent' })
</script>

<style scoped></style>
