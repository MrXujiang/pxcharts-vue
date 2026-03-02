<template>
  <t-space size="8px" class="w-100" direction="vertical">
    <!-- 顶部预览栏 -->
    <div>
      <div>显示格式</div>
    </div>
    <div>
      <t-select v-model="formData.format" :options="formatOptions" />
    </div>
  </t-space>
</template>

<script setup lang="ts">
import { formatOptions } from './options'

const props = defineProps({
  formData: {
    type: Object,
    required: true,
  },
})
const { formData } = toRefs(props)

defineOptions({ name: 'CreateTimeComponent' })
</script>

<style lang="less" scoped></style>
