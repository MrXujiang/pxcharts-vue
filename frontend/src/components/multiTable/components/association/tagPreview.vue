<template>
  <div class="association-container w-100 flx-ce-sta">
    <t-tag
      v-for="(value, index) in tagsData"
      :key="index"
      size="small"
      shape="round"
      style="margin-right: 4px"
    >
      {{ value.value }}
    </t-tag>
  </div>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'

interface Props {
  field: ObjType
  item: ObjType
}

const props = defineProps<Props>()
const tagsData = computed(() => {
  return props.item[props.field.id] || []
})
</script>

<style lang="less" scoped>
.association-container {
  overflow-x: auto;
  scrollbar-width: none;
  &::-webkit-scrollbar {
    display: none;
  }
  &::-webkit-scrollbar {
    display: none;
  }
  -ms-overflow-style: none;
}
</style>
