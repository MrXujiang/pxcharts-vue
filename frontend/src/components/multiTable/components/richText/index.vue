<template>
  <div class="w-100">
    <div class="w-100 pointer flx-ce-sta rich-text-container" @dblclick="handleDoubleClick">
      {{ defaultModelValue }}
    </div>
  </div>
  <t-dialog
    width="50%"
    v-model:visible="visible"
    attach="body"
    :on-close="handleClose"
    :on-confirm="handleConfirm"
  >
    <RichTextEditor
      v-model="basicContent"
      placeholder="请输入内容，支持图片、链接、文件上传等功能..."
      @change="handleBasicChange"
    />
  </t-dialog>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'
import { htmlToText } from '@/utils'
import { getRichTextContent } from '@/api'

interface Props {
  field: ObjType
  item: ObjType
  isEditor: boolean // 是否处于编辑状态
}

const props = defineProps<Props>()
const basicContent = ref('') // 富文本内容
const visible = ref(false) // 是否显示富文本编辑器对话框

// 计算属性：将HTML内容转换为纯文本并截取前十个字符
const defaultModelValue = computed(() => {
  const htmlContent = props.item[props.field.id] || ''
  const textContent = htmlToText(htmlContent)
  // 如果文本长度大于10，则截取前10个字符并添加省略号，否则显示完整文本
  return textContent.length > 10 ? textContent.substring(0, 10) + '...' : textContent
})

const emit = defineEmits(['handleSpecialChange'])

const handleDoubleClick = async () => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const res: any = await getRichTextContent({ recordId: props.item.rowId, fieldId: props.field.id })
  basicContent.value = res.content
  visible.value = true
}

const handleBasicChange = (html: string) => {
  // 更新basicContent的值
  basicContent.value = html
}

const handleClose = () => {
  visible.value = false
  // 重置basicContent为初始值
  basicContent.value = props.item[props.field.id] || ''
}

const handleConfirm = () => {
  visible.value = false
  emit('handleSpecialChange', {
    field: props.field,
    item: props.item,
    value: basicContent.value,
  })
}

defineOptions({ name: 'RichTextComponent' })
</script>

<style lang="less" scoped>
.rich-text-container {
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-height: 32px;
}
</style>
