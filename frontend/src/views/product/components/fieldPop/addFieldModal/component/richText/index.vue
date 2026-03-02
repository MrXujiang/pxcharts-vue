<template>
  <t-space size="small" class="w-100" direction="vertical">
    <div>默认内容</div>
    <t-input v-model="formData.defaultValue" class="w-100" />
  </t-space>
</template>

<script setup lang="ts">
const props = defineProps({
  formData: {
    type: Object,
    required: true,
  },
})
const { formData } = toRefs(props)
defineOptions({ name: 'RichTextComponent' })
</script>

<style lang="less" scoped>
.content {
  height: 32px;
  padding: 0 8px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  align-items: center;
}
</style>
