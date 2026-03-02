<template>
  <div class="image-list">
    <img v-for="(item, index) in listData" :key="index" :src="item.url" alt="" />
  </div>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'

interface Props {
  editingValue: Array<ObjType>
  field: ObjType
}

const props = defineProps<Props>()
const listData = computed(() => {
  return props.editingValue || []
})

defineOptions({
  name: 'imagePreviewCom',
})
</script>

<style scoped lang="less">
.image-list {
  display: flex;
  align-items: center;
  height: 22px;
  gap: 6px;
  overflow-x: auto;
  padding: 0 4px;
  border-radius: 4px;
  scrollbar-width: none;
  &::-webkit-scrollbar {
    display: none;
  }
  &::-webkit-scrollbar {
    display: none;
  }
  -ms-overflow-style: none;
  img {
    height: 100%;
    border-radius: 4px;
    width: auto;
    min-width: 0;
    object-fit: contain;
    flex-shrink: 0;
  }
}
</style>
