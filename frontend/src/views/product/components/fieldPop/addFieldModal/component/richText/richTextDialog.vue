<template>
  <t-dialog
    top="100px"
    v-model:visible="visible"
    header="默认值"
    width="50%"
    :on-confirm="handleConfirmCreate"
    :on-cancel="handleCancel"
  >
    <RichTextEditor
      v-model="basicContent"
      placeholder="请输入内容，支持图片、链接、文件上传等功能..."
      @change="handleBasicChange"
    />
  </t-dialog>
</template>
<!-- <RichTextDialog v-model:visible="visible" @confirmCreate="handleConfirmCreate" /> -->
<script setup lang="ts">
defineOptions({ name: 'RichTextDialog' })
interface Props {
  visible: boolean
}

const props = defineProps<Props>()
const emit = defineEmits(['update:visible', 'confirmCreate'])
const visible = computed(() => props.visible) // 是否显示对话框
const basicContent = ref('') // 富文本内容

// 处理富文本内容变化
const handleBasicChange = (html: string) => {
  console.log('内容更新:', html)
}

// 处理确认创建
const handleConfirmCreate = () => {
  emit('confirmCreate', basicContent.value)
}

// 处理取消
const handleCancel = () => {
  emit('update:visible', false)
}
</script>

<style scoped></style>
