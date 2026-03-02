<template>
  <div>
    <IconRate
      v-model="defaultValue"
      :icon-type="selectedIcon"
      :icon-name="iconData.filter((icon) => icon.type === selectedIcon)[0]?.iconName || ''"
      :count="rateCount"
      :min="formData.settings.scoreStart"
      :max="formData.settings.scoreEnd"
      :disabled="!isEditor"
    />
  </div>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'
import IconRate from '@/components/fieldRate/index.vue'
import { iconData } from '@/views/product/components/fieldPop/addFieldModal/component/rate/options'

interface Props {
  field: ObjType
  item: ObjType
  isEditor: boolean // 是否处于编辑状态
  defaultValue?: number // 默认值
}

const props = defineProps<Props>()
const emit = defineEmits(['handleSpecialChange'])
const formData = computed(() => props.field)
const selectedIcon = ref<string>(props.field.settings.icon)
const defaultValue = ref<number>(props.defaultValue || props.item[props.field.id])

// 计算rate的count值，即分值范围内的选项数量
const rateCount = computed(() => {
  // 确保结束值大于等于起始值
  if (formData.value.settings.scoreEnd >= formData.value.settings.scoreStart) {
    // 选项数量 = 结束值 - 起始值 + 1
    return formData.value.settings.scoreEnd - formData.value.settings.scoreStart + 1
  }
  return 5 // 默认返回5
})

watch(
  () => defaultValue.value,
  (newVal) => {
    emit('handleSpecialChange', {
      field: props.field,
      item: props.item,
      value: newVal,
    })
  },
)

defineOptions({ name: 'RateEditComponent' })
</script>

<style scoped></style>
