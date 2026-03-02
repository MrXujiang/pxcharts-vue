<template>
  <div class="fieldPop">
    <!-- 搜索栏 -->
    <div class="search">
      <t-input v-model="searchValue" borderless clearable placeholder="请输入字段名称">
        <template #prefix-icon>
          <t-icon name="search" />
        </template>
      </t-input>
    </div>
    <!--固定字段列表 - 单独拆分第一个元素 - 第一个固定元素（不参与拖拽） -->
    <div v-if="filteredOptionsData.fixedFields.length > 0" class="first-fixed-item pointer">
      <div class="flex-ce-sta field-item">
        <span class="move-handle flx-ce-ce move-handle-disable">
          <!-- 无拖拽图标 -->
        </span>
        <span class="flex1 item-content flx-ce-sta" style="gap: 12px">
          <t-icon
            class="field-icon"
            :name="getIconNameByType(filteredOptionsData.fixedFields[0].type)"
          />
          <span class="field-label">{{ filteredOptionsData.fixedFields[0].title }}</span>
        </span>
        <span class="flx-ce-sta" style="gap: 4px">
          <t-tooltip content="首字段是每条数据的索引，不支持删除、移动或隐藏">
            <t-icon name="lockOn" />
          </t-tooltip>
          <t-dropdown
            :min-column-width="200"
            trigger="click"
            :options="fixedDropdownOptions"
            @click="handleFieldMoreOptions($event, filteredOptionsData.fixedFields[0])"
          >
            <t-icon name="more" />
          </t-dropdown>
        </span>
      </div>
    </div>
    <!-- 剩余固定字段（参与拖拽） -->
    <div ref="fixedRef">
      <div
        v-for="(item, index) in filteredOptionsData.fixedFields.slice(1)"
        :key="item.id"
        class="flex-ce-sta field-item pointer"
      >
        <span class="move-handle flx-ce-ce">
          <t-icon name="gripVertical" />
        </span>
        <span class="flex1 item-content flx-ce-sta" style="gap: 12px">
          <t-icon class="field-icon" :name="getIconNameByType(item.type)" />
          <span class="field-label">{{ item.title }}</span>
        </span>
        <span class="flx-ce-sta" style="gap: 4px">
          <t-tooltip :content="item.isShow ? '隐藏字段' : '显示字段'">
            <t-icon name="browse" v-if="item.isShow" @click="handleFieldVis(item.id)" />
            <t-icon name="browseOff" v-else @click="handleFieldVis(item.id)" />
          </t-tooltip>
          <t-dropdown
            :min-column-width="200"
            trigger="click"
            :options="getDropdownOptions(index)"
            @click="handleFieldMoreOptions($event, item)"
          >
            <t-icon name="more" />
          </t-dropdown>
        </span>
      </div>
    </div>

    <t-divider
      v-if="
        filteredOptionsData.fixedFields.length > 0 || filteredOptionsData.nonFixedFields.length > 0
      "
    />
    <!-- 非固定字段列表 -->
    <div ref="noFixedRef" class="no-fixed-fields">
      <div
        v-for="(item, index) in filteredOptionsData.nonFixedFields"
        :key="item.id"
        class="flex-ce-sta field-item pointer"
      >
        <span class="move-handle flx-ce-ce">
          <t-icon name="gripVertical" />
        </span>
        <span class="flex1 item-content flx-ce-sta" style="gap: 12px">
          <t-icon class="field-icon" :name="getIconNameByType(item.type)" />
          <span class="field-label">{{ item.title }}</span>
        </span>
        <span class="flx-ce-sta" style="gap: 4px">
          <t-tooltip :content="item.isShow ? '隐藏字段' : '显示字段'">
            <t-icon name="browse" v-if="item.isShow" @click="handleFieldVis(item.id)" />
            <t-icon name="browseOff" v-else @click="handleFieldVis(item.id)" />
          </t-tooltip>

          <t-dropdown
            :min-column-width="200"
            trigger="click"
            :options="getDropdownOptions(index)"
            @click="handleFieldMoreOptions($event, item)"
          >
            <t-icon name="more" />
          </t-dropdown>
        </span>
      </div>
    </div>
    <t-divider
      v-if="
        filteredOptionsData.fixedFields.length > 0 || filteredOptionsData.nonFixedFields.length > 0
      "
    />
    <!-- 底部 添加字段 -->
    <div class="footer flx-ce-sta pointer" style="gap: 6px" @click="addField">
      <t-icon name="add" />
      <span style="font-size: 14px">添加字段</span>
    </div>
  </div>
</template>

<script setup lang="tsx">
import Sortable from 'sortablejs'
import { getIconNameByType } from '@/utils'
import { fieldManageMoreOptions } from '@/modal/options'
import { useMtTableStore, TableColumn } from '@/stores/mtTable'
import { FieldsResult } from './type'
import { ObjType } from '@/types'
import { Button } from 'tdesign-vue-next'
import { deleteField, sortField } from '@/api'

const route = useRoute()
const store = useMtTableStore()
const { updateFields, getGroupedFields, toggleFieldVisibility, delField } = store
let sortableFixed: Sortable | null = null
let sortableNoFixed: Sortable | null = null
const fixedRef = ref<HTMLElement | null>(null)
const noFixedRef = ref<HTMLElement | null>(null)
const searchValue = ref<string>('')

// 计算过滤后的选项数据
const filteredOptionsData = computed<FieldsResult>(() => {
  const allOptionsData = {
    fixedFields: getGroupedFields().fixed,
    nonFixedFields: getGroupedFields().scrollable,
  }

  if (!searchValue.value.trim()) {
    return allOptionsData
  }

  const searchTerm = searchValue.value.toLowerCase().trim()

  const filteredFixedFields = allOptionsData.fixedFields.filter((field) =>
    field.title.toLowerCase().includes(searchTerm),
  )

  const filteredNonFixedFields = allOptionsData.nonFixedFields.filter((field) =>
    field.title.toLowerCase().includes(searchTerm),
  )

  return {
    fixedFields: filteredFixedFields,
    nonFixedFields: filteredNonFixedFields,
  }
})

const optionsData = computed<FieldsResult>(() => {
  return {
    fixedFields: getGroupedFields().fixed,
    nonFixedFields: getGroupedFields().scrollable,
  }
})

const emits = defineEmits(['handleAddField', 'handleEditField', 'handleClosePop'])

// 处理字段更多操作
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const handleFieldMoreOptions = (option: any, field: ObjType) => {
  console.log('option', option)
  console.log('field', field)

  switch (option.value) {
    case 'edit':
      emits('handleEditField', field)
      break
    case 'deleteField':
      emits('handleClosePop')
      const confirmDia = DialogPlugin({
        header: '删除字段/列',
        body: `确认删除 「${field.title}」字段吗？`,
        confirmBtn: () => (
          <Button
            theme="danger"
            onClick={async () => {
              await deleteField({
                id: field.id,
              })
              delField(field.id as string)
              MessagePlugin.success('删除成功')
              confirmDia.hide()
            }}
          >
            删除
          </Button>
        ),
        cancelBtn: '取消',
        onClose: () => {
          confirmDia.hide()
        },
      })
    default:
      break
  }
}

// 添加字段
const addField = () => {
  emits('handleAddField')
}
// 处理字段可见性
const handleFieldVis = (fieldId: string) => {
  toggleFieldVisibility(fieldId)
}
// 获取固定字段第一个字段的更多操作
const fixedDropdownOptions = computed(() => {
  const dropdownOptions = fieldManageMoreOptions?.filter(
    (item) => item.value !== 'moveToTop' && item.value !== 'deleteField',
  )
  return dropdownOptions
})
// 获取字段更多操作
const getDropdownOptions = (index: number) => {
  const dropdownOptions =
    index == 0
      ? fieldManageMoreOptions?.filter((item) => item.value !== 'moveToTop')
      : fieldManageMoreOptions
  return dropdownOptions
}

// 统一处理所有拖拽结束逻辑
const handleDragEnd = async (evt: Sortable.SortableEvent) => {
  const { from, to, oldIndex, newIndex } = evt

  // 如果没有发生位置变化（如拖拽后放回原处）
  if (from === to && oldIndex === newIndex) return

  // 检查当前是否在搜索状态，如果是，则不允许拖拽
  if (searchValue.value.trim()) {
    MessagePlugin.warning('搜索状态下不允许拖拽排序')
    return
  }

  // 获取当前操作的是哪个列表：固定字段剩余部分 或 非固定字段
  const isFromFixed = from === fixedRef.value
  const isToFixed = to === fixedRef.value

  // 提取 fixedFirst（第一个固定字段，始终不变）
  const fixedFirst = optionsData.value.fixedFields[0]
  const restFixedFields = optionsData.value.fixedFields.slice(1)
  const nonFixedFields = optionsData.value.nonFixedFields

  // 构建可变字段数组（用于排序映射）
  const newRestFixed = [...restFixedFields]
  const newNonFixed = [...nonFixedFields]

  // 根据拖拽来源更新数组
  if (isFromFixed && oldIndex !== null) {
    const movedItem = newRestFixed.splice(oldIndex as number, 1)[0]
    if (isToFixed && newIndex !== null) {
      // 固定 → 固定（内部排序）
      newRestFixed.splice(newIndex as number, 0, movedItem)
    } else if (!isToFixed && to === noFixedRef.value && newIndex !== null) {
      // 固定 → 非固定（移出 fixed）
      movedItem.fixed = false
      newNonFixed.splice(newIndex as number, 0, movedItem)
    }
  } else if (!isFromFixed && from === noFixedRef.value && oldIndex !== null) {
    const movedItem = newNonFixed.splice(oldIndex as number, 1)[0]
    if (!isToFixed && newIndex !== null) {
      // 非固定 → 非固定（内部排序）
      newNonFixed.splice(newIndex as number, 0, movedItem)
    } else if (isToFixed && newIndex !== null) {
      // 非固定 → 固定（加入 fixed）
      movedItem.fixed = true
      newRestFixed.splice(newIndex as number, 0, movedItem)
    }
  }
  // 合并最终字段顺序：  1. 第一个固定字段 2. 其余固定字段（已排序） 3. 非固定字段（已排序）

  const orderedFields = [{ ...fixedFirst }, ...newRestFixed, ...newNonFixed] as TableColumn[]
  console.log('orderedFields', orderedFields)
  await sortField({
    fieldIds: orderedFields.map((field) => field.id),
    tableSchemaId: route.params.tableSchemaId as string,
  })
  updateFields(orderedFields)
}

onMounted(() => {
  nextTick(() => {
    if (noFixedRef.value) {
      sortableNoFixed = new Sortable(noFixedRef.value, {
        group: 'fields',
        animation: 150, // 拖拽动画时间
        ghostClass: 'sortable-ghost', // 拖拽时的样式类
        chosenClass: 'sortable-chosen', // 选中项样式
        handle: '.move-handle', // 可拖拽的把手（可选）
        onEnd: handleDragEnd,
      })
    }
    if (fixedRef.value) {
      sortableFixed = new Sortable(fixedRef.value, {
        group: 'fields',
        animation: 150, // 拖拽动画时间
        ghostClass: 'sortable-ghost', // 拖拽时的样式类
        chosenClass: 'sortable-chosen', // 选中项样式
        handle: '.move-handle', // 可拖拽的把手（可选）
        // 当没有固定字段时，不需要过滤不可拖拽的元素
        filter: optionsData.value.fixedFields.length > 0 ? '.move-handle-disable' : undefined,
        onEnd: handleDragEnd,
      })
    }
  })
})
onUnmounted(() => {
  if (sortableFixed) sortableFixed.destroy()
  if (sortableNoFixed) sortableNoFixed.destroy()
})
defineOptions({
  name: 'FieldPop',
})
</script>

<style scoped lang="less">
@import './index.less';
</style>
