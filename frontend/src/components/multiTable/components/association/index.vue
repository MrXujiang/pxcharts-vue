<template>
  <div
    ref="associationContainerRef"
    class="association-container w-100 h-100 flx-ce-sta"
    style="min-height: 32px"
    @dblclick="handleDoubleClick"
  >
    <template v-if="defaultModelValue && defaultModelValue.length > 0">
      <t-tag
        v-for="(value, index) in defaultModelValue || []"
        :key="index"
        shape="round"
        style="margin-right: 4px"
        closable
        @close="handleCloseTag(index)"
      >
        {{ value.value }}
      </t-tag>
    </template>
    <template v-else>
      <div class="w-100 gray-col" style="font-size: 14px; text-align: center">双击编辑</div>
    </template>
  </div>
  <t-dialog
    v-if="visible"
    dialogClassName="association-dialog"
    v-model:visible="visible"
    attach="body"
    header="未命名记录"
    width="80%"
    top="100px"
    @close-btn-click="handleCancel"
  >
    <div class="flx-ce-sta gray-col subTitle gap-8">
      <div>已关联</div>
      <t-tooltip content="点击跳转到关联表「这里是表名」">
        <div class="flx-ce-sta gap-2 pointer">
          <t-icon size="12" name="database" />
          <span>这里是表名哦</span>
        </div>
      </t-tooltip>
    </div>
    <t-space direction="vertical" class="w-100">
      <t-tabs v-model="activeTab" @change="handleTabChange">
        <t-tab-panel :value="1" label="选择要关联的记录"> </t-tab-panel>
        <t-tab-panel :value="2" :label="`已选 ${selectedKeys.length} 条记录`"> </t-tab-panel>
      </t-tabs>
      <!-- 搜索 -->
      <div class="w-100 flx-ce-sta" style="padding: 0 8px">
        <t-input
          v-model="searchText"
          placeholder="请输入搜索内容"
          clearable
          @change="handleSearch"
        />
      </div>
      <div class="w-100" style="height: 500px; overflow: hidden; overflow-y: auto">
        <MultiTable
          :mvTableData="tableData"
          :fixedFieldsData="fixedFields"
          :scrollableFieldsData="scrollableFields"
          :isAssociated="true"
          :isEditor="false"
          :isMultipleSelection="isMultipleSelection"
          :selectedKeys="selectedKeys.map((key) => key.recordId)"
          @handle-toggle-row-selection="handleToggleRowSelection"
        />
      </div>
    </t-space>

    <template #footer>
      <div class="flx-ce-end">
        <!-- <t-button theme="default">添加记录</t-button> -->
        <t-space size="small" class="flx-ce-sta">
          <div class="gray-col">已选 {{ selectedKeys.length }} 条记录</div>
          <t-button theme="default" @click="handleCancel">取消</t-button>
          <t-button theme="primary" @click="handleConfirm">确定</t-button>
        </t-space>
      </div>
    </template>
  </t-dialog>
</template>

<script setup lang="ts">
import { useMtTableStore } from '@/stores/mtTable'
import { ObjType } from '@/types'
import { TabsProps } from 'tdesign-vue-next'
import { getTableAllRecords, updateCellData, getRecord } from '@/api'

interface Props {
  field: ObjType
  editingValue: ObjType[]
  isEditor: boolean
  item: ObjType
}
const props = defineProps<Props>()
console.log('props-association', props)
const mtTableStore = useMtTableStore()
const emit = defineEmits(['update:editingValue', 'handleCellEdit', 'handleCancelCellEdit'])

const isMultipleSelection = computed(() =>
  props.field.type === 'singleAssociation' ? false : true,
) // 是否允许多选
const associationContainerRef = ref(null)
const route = useRoute()
const visible = ref<boolean>(false)
const activeTab = ref<TabsProps['value']>(1) // 当前选中的标签页，1表示选择要关联的记录，2表示已选 X 条记录
const selectedKeys = ref<ObjType[]>(props.editingValue || []) // 选中的键
const searchText = ref<string>('') // 搜索文本
// const cacheValue = ref<ObjType[]>(props.editingValue || []) // 缓存值，用作合并
// 表格数据
const tableData = ref<ObjType[]>([]) // 表格数据
const originalTableData = ref<ObjType[]>([]) // 原始表格数据
const fixedFields = ref<ObjType[]>([]) // 字段数据
const scrollableFields = ref<ObjType[]>([]) // 滚动字段数据

const getTableData = async () => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const res: any = await getTableAllRecords({ tableSchemaId: route.params.tableSchemaId as string })

  fixedFields.value = res.fields.slice(0, 1) || []
  scrollableFields.value = res.fields.slice(1) || []
  originalTableData.value = res.records || []
  tableData.value = originalTableData.value
}

const defaultModelValue = computed({
  get() {
    return props.editingValue
  },
  set(value) {
    emit('update:editingValue', value)
  },
})
console.log('defaultModelValue.value', defaultModelValue.value)

// 处理双击
const handleDoubleClick = async () => {
  await getTableData()
  nextTick(() => {
    visible.value = true
  })
}

// 处理标签页切换
const handleTabChange: TabsProps['onChange'] = (newValue) => {
  activeTab.value = newValue
  if (newValue === 2) {
    const valueSet = new Set(selectedKeys.value?.map((i) => i?.recordId).filter(Boolean))
    console.log('valueSet', valueSet)
    tableData.value = tableData.value.filter((row) => row?.rowId && valueSet.has(row.rowId))
    console.log('tableData.value', tableData.value)
  } else {
    tableData.value = [...originalTableData.value]
  }
}

// 处理行选择切换
const handleToggleRowSelection = (selectedRowKeys: ObjType[]) => {
  const selectedRowIdSet = new Set(selectedRowKeys)

  selectedKeys.value = [...selectedRowIdSet]
  //   通过第一列字段类型 找到数据中对应的值
  const firstColumnField = fixedFields.value[0].id
  //   selectedKeys.value 在tableData中找到rowId对应的值并返回
  const selectedValues = selectedKeys.value.map((key) => {
    const foundRow = tableData.value.find((row) => row.rowId === key)
    return {
      value: foundRow ? foundRow[firstColumnField] : '未命名记录',
      recordId: foundRow ? foundRow.rowId : '',
    }
  })
  selectedKeys.value = [...selectedValues]
}
// 处理取消
const handleCancel = () => {
  visible.value = false
  emit('handleCancelCellEdit')
}

// 更新单元格数据逻辑
const updateCellDataLogic = async () => {
  await updateCellData({
    fieldId: props.field.id,
    recordId: props.item.rowId,
    value: selectedKeys.value.map((key) => key.recordId),
  })
  const res = await getRecord({
    tableSchemaId: route.params.tableSchemaId as string,
    recordId: props.item.rowId,
  })
  mtTableStore.updateRecord(props.item.rowId, { ...res })
}
// 处理确定
const handleConfirm = async () => {
  console.log('selectedKeys.value', selectedKeys.value)
  console.log('props', props)

  await updateCellDataLogic() // 更新单元格数据
  emit('handleCancelCellEdit')
  visible.value = false
}
// 处理搜索
const handleSearch = () => {
  console.log('Search Text:', searchText.value)
}

const handleCloseTag = async (index: number) => {
  // 创建新数组，移除指定索引的元素
  const newDefaultModelValue = [...defaultModelValue.value]
  newDefaultModelValue.splice(index, 1)

  // 通过计算属性的 setter 触发更新
  defaultModelValue.value = newDefaultModelValue

  // 同样处理 selectedKeys
  const newSelectedKeys = [...selectedKeys.value]
  newSelectedKeys.splice(index, 1)
  selectedKeys.value = newSelectedKeys

  await updateCellDataLogic()
}

// 全局点击监听
const handleGlobalClick = (e: MouseEvent) => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const container: any = associationContainerRef.value
  if (!container) return
  if (visible.value) return
  const isClickInContainer = container.contains(e.target as Node)

  // 4. 若既不在容器内，也不在弹窗内 → 取消编辑
  if (!isClickInContainer) {
    emit('handleCellEdit')
    handleCancel()
  }
}
// 挂载时监听全局点击
onMounted(() => {
  document.addEventListener('click', handleGlobalClick)
})

// 卸载时移除监听（避免内存泄漏）
onUnmounted(() => {
  document.removeEventListener('click', handleGlobalClick)
})
defineOptions({ name: 'AssociationComponent' })
</script>

<style lang="less" scoped>
.association-container {
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
.subTitle {
  font-size: 12px;
}
</style>
