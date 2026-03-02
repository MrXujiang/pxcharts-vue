<template>
  <div class="color-picker-panel">
    <div v-for="(row, rowIndex) in colorRows" :key="rowIndex" class="color-row">
      <div
        v-for="(color, colIndex) in row"
        :key="colIndex"
        class="color-item"
        :style="{ backgroundColor: color }"
        @click="selectColor(color)"
      >
        <t-icon v-if="color === selectedColor" name="check" class="check-icon" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { colorRows } from '@/modal/options'

const props = defineProps({
  colVal: {
    type: String,
    default: '',
  },
})
const emit = defineEmits(['update:colVal'])
const selectedColor = computed(() => props.colVal)

const selectColor = (color: string) => {
  emit('update:colVal', color)
}
defineOptions({
  name: 'ColorPanel',
})
</script>

<style lang="less" scoped>
@import './index.less';
</style>
