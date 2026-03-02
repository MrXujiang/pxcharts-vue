<template>
  <!-- 单选字段 -->
  <div class="tag-preview-container flx-ce-sta">
    <t-tag
      v-if="field.type === 'radio' && item[field.id]"
      size="small"
      shape="round"
      :color="getTagInfo(field, item[field.id]).color"
    >
      {{ getTagInfo(field, item[field.id]).label }}
    </t-tag>

    <!-- 多选字段 -->
    <template v-else-if="field.type === 'selectMultiple'">
      <t-tag
        v-for="(value, index) in Array.isArray(item[field.id]) ? item[field.id] : []"
        :key="index"
        size="small"
        shape="round"
        :color="getTagInfo(field, value).color"
        style="margin-right: 4px"
      >
        {{ getTagInfo(field, value).label }}
      </t-tag>
    </template>
    <template v-else>
      <span></span>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'

interface Props {
  field: ObjType // 当前字段信息
  item: ObjType // 当前操作行数据
}

defineProps<Props>()

// 获取标签信息的函数
const getTagInfo = (field, value) => {
  // 处理空值情况
  if (!value && value !== 0) {
    return { color: '', label: '' }
  }

  // 找到对应的选项
  const foundItem = field.options.find((option) => option.id == value)

  // 如果找到了匹配的选项，返回其颜色和标签
  if (foundItem) {
    return {
      color: foundItem.color,
      label: foundItem.label,
    }
  }

  // 如果没有找到匹配的选项，返回默认值
  return { color: 'default', label: value }
}

defineOptions({ name: 'SingleMultipleChoiceTagPreviewComponent' })
</script>

<style lang="less" scoped>
.tag-preview-container {
  width: 100%;
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
