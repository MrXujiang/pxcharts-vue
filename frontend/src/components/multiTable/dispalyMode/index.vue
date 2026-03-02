<template>
  <div class="w-100">
    <span v-if="isText" class="cell-content">
      <!-- 货币 -->
      <span v-if="field?.type == 'currency'"
        >{{ currencyMap[field.settings.currency]
        }}{{ formatCellValue(item[field.id], field) }}</span
      >
      <!-- 链接 -->
      <template v-else-if="field?.type == 'link'">
        <div v-if="isEditor" class="flx-ce-sta gap-2">
          <t-icon name="link" style="color: blue" @click="handleClickLink" />
          <div>
            {{
              formatCellValue(item[field.id], field)?.linkTitle
                ? formatCellValue(item[field.id], field)?.linkTitle
                : formatCellValue(item[field.id], field)?.linkUrl
            }}
          </div>
        </div>

        <span v-else>{{
          formatCellValue(item[field.id], field)?.linkTitle
            ? formatCellValue(item[field.id], field)?.linkTitle
            : formatCellValue(item[field.id], field)?.linkUrl
        }}</span>
      </template>
      <template v-else-if="field?.type == 'date'">
        <span>{{
          formatCellValue(item[field.id], field)
            ? dayjs(formatCellValue(item[field.id], field)).format(field.format)
            : ''
        }}</span>
      </template>
      <span v-else>{{ formatCellValue(item[field.id], field) }}</span>
    </span>

    <!-- 进度条显示模式 -->
    <ProgressPreviewComponent
      v-if="field?.type == 'progress'"
      :field="field"
      :editingValue="formatCellValue(item[field.id], field)"
    />
    <!-- 单选、多选显示模式 -->
    <SingleMultipleChoiceTagPreviewComponent
      v-if="['radio', 'selectMultiple'].includes(field?.type || '')"
      :field="field"
      :item="item"
    />
    <PersonPreviewComponent
      v-if="field?.type == 'person'"
      :field="field"
      :editingValue="formatCellValue(item[field.id], field)"
    />
    <!-- 关联显示模式 -->
    <AssociationTagPreviewComponent
      v-if="['singleAssociation', 'doubleAssociation'].includes(field?.type || '')"
      :field="field"
      :item="item"
    />
  </div>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'
import dayjs from 'dayjs'
import { currencyMap } from '@/views/product/components/fieldPop/addFieldModal/component/currency/options'
import { formatCellValue } from '@/utils/mtable'
import SingleMultipleChoiceTagPreviewComponent from '@/components/multiTable/components/singleMultipleChoice/tagPreview.vue'
const ProgressPreviewComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/progress/progressPreview.vue'),
)
const PersonPreviewComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/person/prewView.vue'),
)
const AssociationTagPreviewComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/association/tagPreview.vue'),
)

interface Props {
  field: ObjType
  item: ObjType
  isEditor: boolean // 是否处于编辑状态
}

const props = defineProps<Props>()

const field = computed(() => props.field)
const item = computed(() => props.item)

const isText = computed(() => {
  return ['text', 'number', 'date', 'currency', 'link', 'email', 'idCard', 'richText'].includes(
    field.value?.type || '',
  )
})

// 处理点击链接
const handleClickLink = () => {
  window.open('https://' + formatCellValue(item.value[field.value.id], field.value)?.linkUrl)
}

defineOptions({ name: 'DisplayModeComponent' })
</script>

<style scoped></style>
