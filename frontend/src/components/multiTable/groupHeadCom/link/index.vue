<template>
  <div class="link">
    <div class="link-label">链接</div>
    <span
      v-if="(value && value.linkUrl) || value.linkTitle"
      class="link-item pointer"
      @click="handleOpenLink"
      >{{ value.linkTitle ? value.linkTitle : value.linkUrl }}</span
    >
    <span v-else class="no-content">无链接</span>
  </div>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'

interface Props {
  value: ObjType
}

const props = defineProps<Props>()
console.log('props', props)

const handleOpenLink = () => {
  window.open('https://' + props.value.linkUrl)
}

defineOptions({ name: 'GroupHeadLink' })
</script>

<style scoped>
.link {
  font-size: 10px;
  color: #2f75fd;
  display: flex;
  flex-direction: column;
  min-width: 0; /* 允许flex item收缩至内容宽度以下 */
}

.link-label {
  margin-bottom: 2px;
  color: #999;
}

.link-item {
  display: inline-block;
  color: #2f75fd;
  text-decoration: none;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
  min-width: 0;
  flex-shrink: 1;
}

.link-item:hover {
  text-decoration: underline;
}

.no-content {
  color: #999;
  font-size: 12px;
}
</style>
