<template>
  <div class="w-100">
    <t-space size="8px" class="w-100" direction="vertical">
      <!-- 顶部预览栏 -->
      <div class="flx-ce-bet nav-bar">
        <div>显示格式</div>
      </div>
      <div>
        <t-select v-model="formData.format" :options="formatOptions" />
      </div>

      <div class="w-100 flex-col gap-8">
        <t-checkbox v-model="formData.settings.isDefaultCreateTime">默认填充创建时间</t-checkbox>
      </div>
    </t-space>
  </div>
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

defineOptions({ name: 'DateComponent' })
</script>

<style lang="less" scoped>
@import './index.less';
</style>
