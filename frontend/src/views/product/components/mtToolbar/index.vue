<template>
  <div class="layout-toolbar w-100 flx-ce-bet">
    <div class="flx-ce-sta">
      <t-button variant="text" @click="handleAddRecord">
        <template #icon><t-icon style="color: #0052d9" name="addCircleFilled" /></template>
        添加记录
      </t-button>

      <t-divider layout="vertical" />
      <t-popup placement="bottom-left" v-model:visible="visible" trigger="click">
        <t-button variant="text">
          <template #icon><t-icon style="color: #6e7286" name="setting1" /></template>
          字段管理
        </t-button>
        <template #content>
          <FieldPop
            @handleAddField="addField"
            @handleEditField="handleEditField"
            @handleClosePop="visible = false"
          />
        </template>
      </t-popup>

      <t-divider layout="vertical" />
      <t-popup placement="bottom-left" v-model:visible="filterVis" variant="text" trigger="click">
        <t-button variant="text">
          <template #icon><t-icon style="color: #6e7286" name="filter" /></template>
          筛选
        </t-button>
        <template #content>
          <FilterPopup v-if="filterVis" :filterConfig="filterConfig" />
        </template>
      </t-popup>

      <t-popup placement="bottom-left" trigger="click">
        <t-button variant="text">
          <template #icon><t-icon style="color: #6e7286" name="viewAgenda" /></template>
          分组
        </t-button>
        <template #content>
          <t-card size="small" :bordered="false">
            <t-select
              style="width: 240px"
              v-model="groupByField"
              placeholder="请选择分组字段"
              @change="handleGroupFieldChange"
              clearable
            >
              <template #prefixIcon>
                <span style="width: 16px; height: 16px" class="flx-ce-ce">
                  <t-icon
                    :name="
                      getIconNameByType(fields.find((field) => field.id === groupByField)?.type)
                    "
                  />
                </span>
              </template>
              <template #valueDisplay="{ value }">
                <div class="flx-ce-sta gap-4">
                  <!-- <t-icon :name="getIconName(value)" /> -->
                  <span class="gray-col">{{ getGroupSelectValue(value) }}</span>
                </div>
              </template>
              <t-option v-for="field in fields" :key="field.id" :value="field.id">
                <div class="flx-ce-sta gap-8">
                  <t-icon size="16" :name="getIconNameByType(field.type)" />
                  <div>{{ field.title }}</div>
                </div>
              </t-option>
            </t-select>
          </t-card>
        </template>
      </t-popup>
      <t-popup placement="bottom-left" trigger="click">
        <t-button variant="text">
          <template #icon><t-icon style="color: #6e7286" name="arrowUpDown1" /></template>
          排序
        </t-button>
        <template #content>
          <SortPopup :sortConfig="sortConfig" />
        </template>
      </t-popup>
      <t-popup placement="bottom-left" trigger="click">
        <t-button variant="text">
          <template #icon><t-icon style="color: #6e7286" name="expandVertical" /></template>
          行高
        </t-button>
        <template #content>
          <div class="row-height-options">
            <div
              v-for="option in rowHeightOptions"
              :key="option.value"
              class="row-height-item"
              :class="{ 'row-height-item--active': rowHeight === option.value }"
              @click="selectRowHeight(option.value)"
            >
              {{ option.content }}
            </div>
          </div>
        </template>
      </t-popup>

      <t-button variant="text">
        <template #icon><t-icon style="color: #6e7286" name="fillColor1" /></template>
        填色
      </t-button>
    </div>
    <div class="flx-ce-ce">
      <t-tooltip content="搜索">
        <t-button variant="text">
          <template #icon><t-icon style="color: #6e7286" name="search" /></template>
        </t-button>
      </t-tooltip>
    </div>
  </div>
  <AddFieldModal
    v-if="isAddVisible"
    v-model:isAddVisible="isAddVisible"
    :fieldConfig="fieldConfig"
  />
  <AddRecord v-if="isAddRecordVisible" v-model:visible="isAddRecordVisible" />
</template>

<script setup lang="ts">
import FieldPop from '../fieldPop/index.vue'
import AddFieldModal from '../fieldPop/addFieldModal/index.vue'
import AddRecord from '@/components/multiTable/components/addRecord/index.vue'
import { useMtTableStore } from '@/stores/mtTable'
import { getIconNameByType, clearReactiveObject } from '@/utils'
import { handleRowChangeValue } from '@/utils/mtable'
import { SelectProps } from 'tdesign-vue-next'
import { updateViewTable, insertRecordApi, getTableData } from '@/api'
import { ObjType } from '@/types'
import { cloneDeep } from 'lodash'
import FilterPopup from './component/filter/index.vue'
import SortPopup from './component/sort/index.vue'
import emitter from '@/utils/mitt'
import { rowHeightOptions } from '@/modal/options'

interface Props {
  viewConfig: ObjType
}
const props = defineProps<Props>()

const route = useRoute()
const viewId = computed(() => props.viewConfig.id) // 当前激活项视图id
const visible = ref<boolean>(false) // 字段管理popup是否可见
const filterVis = ref<boolean>(false) // 筛选popup是否可见
const { getFields, updateTableSettings, currentTable, insertRecord, getRecords, updateRecords } =
  useMtTableStore()
// 获取字段列表，过滤掉image、attachment类型
const fields = computed(() =>
  getFields().filter(
    (field) =>
      !['image', 'attachment', 'autoNumber', 'singleAssociation', 'doubleAssociation'].includes(
        field.type,
      ),
  ),
)
// 获取分组字段
const groupByField = computed(() =>
  currentTable.settings.tableConfig.groupConfig?.length > 0
    ? currentTable.settings.tableConfig.groupConfig[0]
    : '',
)

// 获取筛选配置数据
const filterConfig = computed(() => currentTable.settings.tableConfig?.filterConfig || [])
const sortConfig = computed(() => currentTable.settings.tableConfig?.sortConfig || [])

const fieldConfig = reactive<ObjType>({}) // 字段配置

const isAddVisible = ref<boolean>(false) // 添加字段模态框是否可见
const isAddRecordVisible = ref<boolean>(false) // 添加记录模态框是否可见

// 当前行高，默认为常规
const rowHeight = ref<string>(currentTable.settings.tableConfig.rowHeight || '1')
// 监听添加字段事件
emitter.on('handleAddColumnClick', () => {
  isAddVisible.value = true
})

// 监听handleSort事件，触发排序
emitter.on('handleSort', (async (list: ObjType[]) => {
  await updateViewTable({
    viewId: viewId.value,
    // 需要return 出 list 中的除了radioGroup以外的值 不希望用删除来操作
    sortConfig: list.map((item) => {
      return {
        fieldId: item.fieldId,
        order: item.order,
      }
    }),
  })
  const obj = {
    tableConfig: {
      sortConfig: list.map((item) => {
        return {
          fieldId: item.fieldId,
          order: item.order,
        }
      }),
    },
  }
  updateTableSettings(obj) // 更新当前表配置
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const res: any = await getTableData({ viewId: viewId.value as string })
  console.log('res', res)
  updateRecords(res.records)
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
}) as any)

// 监听handleFilter事件，触发筛选
emitter.on('handleFilter', (async (filterList: ObjType[]) => {
  console.log('filterList', filterList)
  await updateViewTable({
    viewId: viewId.value,
    filterConfig: (filterList as ObjType[]) || [],
  })
  const obj = {
    tableConfig: {
      filterConfig: (filterList as ObjType[]) || [],
    },
  }
  console.log('obj', obj)
  updateTableSettings(obj) // 更新分组字段
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const res: any = await getTableData({ viewId: viewId.value as string })
  console.log('res', res)
  updateRecords(res.records)
  //   setCurrentTable(res) // 设置当前表格数据
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
}) as any)

// 处理行高改变
// 选择行高
const selectRowHeight = async (value: string) => {
  rowHeight.value = value
  await updateViewTable({
    viewId: viewId.value,
    rowHeight: Number(value),
  })
  const obj = {
    tableConfig: {
      rowHeight: Number(value),
    },
  }
  updateTableSettings(obj) // 更新当前表配置
  emitter.emit('handleRowHeightChange', value)
}

// 处理分组字段改变
const getGroupSelectValue = (value: string) => {
  return fields.value.find((option) => option.id === value)?.title || ''
}

// 处理字段管理编辑操作
const handleEditField = (field: ObjType) => {
  Object.assign(fieldConfig, { ...cloneDeep(field) })
  console.log('fieldConfig', fieldConfig)
  visible.value = false // 隐藏字段管理popup
  isAddVisible.value = true
}

// 添加字段
const addField = () => {
  clearReactiveObject(fieldConfig)
  console.log('fieldConfig', fieldConfig)
  visible.value = false // 隐藏字段管理popup
  isAddVisible.value = true
}
// 添加记录
const handleAddRecord = async () => {
  // 获取字段配置
  const fields = ref(cloneDeep(getFields()))
  // itemObj 为 fields id值为key，defaultValue 字段值为value
  const fieldItem = fields.value.reduce((acc, cur) => {
    acc[cur.id] = cur.defaultValue
    return acc
  }, {})

  const rawData = getRecords() // 获取记录数据

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const res: any = await insertRecordApi({
    prevId: rawData[0]?.id,
    tableSchemaId: route.params.tableSchemaId as string,
  })
  insertRecord({ ...fieldItem, ...res })
  isAddRecordVisible.value = true
}

// 处理分组字段改变
const handleGroupFieldChange: SelectProps['onChange'] = async (value) => {
  console.log('value', value)
  await updateViewTable({
    viewId: viewId.value,
    groupConfig: value ? [value as string] : [],
  })
  const obj = {
    tableConfig: {
      groupConfig: value ? [value as string] : [],
    },
  }
  updateTableSettings(obj) // 更新分组字段
}
defineOptions({
  name: 'MtToolbar',
})
</script>

<style lang="less" scoped>
@import './index.less';
</style>
