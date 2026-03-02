<!-- eslint-disable vue/no-parsing-error -->
<template>
  <div class="virtual-table" ref="tableContainer">
    <!-- 工具栏 -->
    <!-- <div class="table-toolbar">
      <div class="toolbar-left">
        <h3>多维表格 ({{ totalRecords.toLocaleString() }} 条记录)</h3>
        <button @click="generateTestData(1000)" class="btn">生成1K数据</button>
        <button @click="generateTestData(10000)" class="btn">生成1W数据</button>
        <button @click="generateTestData(100000)" class="btn">生成10W数据</button>
        <label class="config-option">
          <input type="checkbox" v-model="tableConfig.stickyHeader" />
          固定表头
        </label>
        <label class="config-option">
          <input type="checkbox" v-model="tableConfig.enableCrossGroupDrag" />
          跨分组拖拽
        </label>
      </div>
      <div class="toolbar-right">
        <div class="batch-actions" v-if="selectedRows.size > 0">
          <span class="selected-count">已选中 {{ selectedRows.size }} 行</span>
          <button @click="batchDelete" class="btn btn-danger">批量删除</button>
          <button @click="clearSelection" class="btn">取消选择</button>
        </div>
        <select v-model="groupByField" class="group-select">
          <option value="">无分组</option>
          <option v-for="field in fields" :key="field.id" :value="field.id">
            按{{ field.title }}分组
          </option>
        </select>
        <button @click="exportData" class="btn">导出数据</button>
      </div>
    </div> -->

    <!-- 表格容器 -->
    <div class="table-wrapper">
      <!-- 表头区域 -->
      <div class="table-header-wrapper" :class="{ 'header-sticky': tableConfig.stickyHeader }">
        <!-- 固定列表头 -->
        <div class="fixed-header-section" :style="{ width: getFixedColumnsWidth() + 'px' }">
          <!-- 全选复选框 -->
          <div
            v-if="isMultipleSelection"
            class="header-selector"
            :style="{ borderRadius: isGrouped ? '8px 0 0 8px' : '' }"
          >
            <input
              type="checkbox"
              :checked="isAllSelected"
              :indeterminate="isPartiallySelected"
              @change="toggleSelectAll"
            />
          </div>
          <!-- 固定列头 -->
          <div
            v-for="field in fixedFields"
            :key="field.id"
            class="column-header fixed-column-header"
            :style="{
              width: field.width + 'px',
              left: getFixedColumnStickyLeft(fields.indexOf(field)) + 'px',
              borderBottom: isGrouped ? '1px solid #e1e5e9' : '',
            }"
            :data-field-id="field.id"
            :data-index="fields.indexOf(field)"
            draggable="true"
            @dragstart="handleHeaderDragStart"
            @dragover="handleHeaderDragOver"
            @drop="handleHeaderDrop"
          >
            <div class="header-content">
              <t-icon class="header-icon" :name="getIconNameByType(field.type)" />
              <span class="header-title">{{ field.title }} </span>
            </div>
            <div class="column-resizer" @mousedown="startColumnResize($event, field)"></div>
          </div>
        </div>

        <!-- 可滚动列表头 -->
        <div
          class="scrollable-header-section"
          ref="scrollableHeaderRef"
          :style="{
            width: getScrollableColumnsWidthWithAdd() + 'px',
            left: getFixedColumnsWidth() + 'px',
          }"
        >
          <div class="scrollable-header-container" ref="scrollableHeaderContainer">
            <div
              v-for="(field, index) in scrollableFields"
              :key="field.id"
              class="column-header scrollable-column-header"
              :style="{
                width: field.width + 'px',
                left: getScrollableColumnLeftByScrollableIndexWithAdd(index) + 'px',
                borderBottom: isGrouped ? '1px solid #e1e5e9' : '',
              }"
              :data-field-id="field.id"
              :data-index="fields.indexOf(field)"
              draggable="true"
              @dragstart="handleHeaderDragStart"
              @dragover="handleHeaderDragOver"
              @drop="handleHeaderDrop"
            >
              <div class="header-content">
                <t-icon class="header-icon" :name="getIconNameByType(field.type)" />
                <span class="header-title">{{ field.title }} </span>
              </div>
              <div class="column-resizer" @mousedown="startColumnResize($event, field)"></div>
            </div>
            <!-- 添加列表头 -->
            <!-- <t-tooltip content="添加字段/列"> -->
            <div
              class="column-header scrollable-column-header add-column-header"
              :class="isGrouped ? 'isGroupedAddedColumn' : ''"
              style="cursor: pointer"
              :style="{
                width: 56 + 'px',
                left:
                  getScrollableColumnLeftByScrollableIndexWithAdd(scrollableFields.length) + 'px',
              }"
              @click="handleAddColumnClick"
            >
              <div v-if="isEditor || !isAssociated" class="header-content flx-ce-ce h-100">
                <span class="header-title">
                  <t-icon class="icon" name="add" />
                </span>
              </div>
            </div>
            <!-- </t-tooltip> -->
          </div>
        </div>
      </div>

      <!-- 表格主体区域 -->
      <div class="table-body-wrapper">
        <!-- 固定列数据区域 -->
        <div class="fixed-columns-wrapper" :style="{ width: getFixedColumnsWidth() + 'px' }">
          <div class="fixed-columns-container" ref="fixedBodyContainer" @scroll="handleFixedScroll">
            <!-- 固定列虚拟滚动体 -->
            <div
              class="fixed-virtual-body"
              @contextmenu="handleContextMenu"
              :style="{ height: tableHeight + 'px' }"
            >
              <!-- 渲染固定列的可见行 -->
              <div
                v-for="item in visibleItems"
                :key="'fixed-' + getRowKey(item)"
                :data-row-key="'fixed-' + getRowKey(item)"
                class="fixed-row"
                :class="{
                  'row-selected': selectedRows.has(getRowKey(item)),
                  'row-group-header': item.isGroupHeader,
                  'group-last-row': item.isLastInGroup,
                  'is-add-row': item.isAddRow,
                  'data-row-last': item.isLastDataRow && !isGrouped,
                }"
                :style="{
                  transform: `translateY(${item.translateY}px)`,
                  height: itemHeight + 'px',
                  width: getFixedColumnsWidth() + 'px',
                  borderBottom: collapsedGroups.has(item.groupKey) ? '1px solid #e1e5e9' : '',
                  borderRadius: collapsedGroups.has(item.groupKey) ? '8px 0 0 8px' : '',
                }"
                @click="handleRowClick(item, $event)"
                draggable="true"
                @dragstart="handleRowDragStart($event, item)"
                @dragend="handleRowDragEnd"
                @dragover="handleRowDragOver"
                @drop="handleRowDrop($event, item)"
              >
                <!-- 分组头部 -->
                <template v-if="item.isGroupHeader">
                  <div class="group-header-content">
                    <t-icon
                      style="font-size: 22px"
                      class="icon pointer"
                      :name="
                        collapsedGroups.has(item.groupKey) ? 'caretRightSmall' : 'caretDownSmall'
                      "
                      @click="toggleGroup(item.groupKey)"
                    />
                    <component
                      v-if="
                        !['text', 'number', 'email', 'idCard', 'richText', 'findRef'].includes(
                          getFieldTypeByGroupByField(),
                        )
                      "
                      :is="isActiveGroupComponent(getFieldTypeByGroupByField())"
                      :value="item.groupKey"
                    />

                    <span v-else class="group-title">{{ item.groupTitle }}11</span>
                    <span class="group-count">({{ item.groupCount }})</span>
                  </div>
                </template>

                <!-- 添加行 -->
                <template v-else-if="item.isAddRow">
                  <!-- 合并的固定列添加行单元格 -->
                  <div
                    class="table-cell fixed-cell is-add-row-cell"
                    :style="{
                      width: getFixedColumnsWidth() + 'px',
                      left: '0px',
                    }"
                    @click="handleAddRowClick(item.groupKey, $event)"
                    @mouseenter="handleCellMouseEnter(item)"
                    @mouseleave="handleCellMouseLeave(item)"
                  >
                    <span class="cell-content">
                      <t-icon class="icon" name="add" />
                    </span>
                  </div>
                </template>

                <!-- 固定列数据行 -->
                <template v-else>
                  <!-- 行选择器 -->
                  <div class="row-selector" @click.stop>
                    <input
                      type="checkbox"
                      :checked="selectedRows.has(getRowKey(item))"
                      @change="toggleRowSelection(item)"
                      @click="$event.stopPropagation()"
                    />
                  </div>

                  <!-- 固定列数据  fixed-cell-->
                  <div
                    v-for="field in fixedFields"
                    :key="field.id"
                    class="table-cell"
                    :class="isGrouped ? 'fixed-cell' : ''"
                    :style="{
                      width: field.width + 'px',
                      left: getFixedColumnLeft(fields.indexOf(field)) + 'px',
                    }"
                    @click="handleCellClick(item, field, $event)"
                    @dblclick="startCellEdit(item, field)"
                    @mouseenter="handleCellMouseEnter(item)"
                    @mouseleave="handleCellMouseLeave(item)"
                  >
                    <!-- 特殊字段不区分显示or编辑模式 -->
                    <template v-if="isEditor && isSpecialField(field)">
                      <component
                        :is="isActiveComponent(field?.type || '')"
                        :field="field"
                        :item="item"
                        @handle-special-change="handleSpecialChange"
                      />
                    </template>

                    <template v-else>
                      <!-- 固定类单元格编辑模式 -->
                      <div
                        v-if="
                          isEditor &&
                          editingCell &&
                          editingCell.rowKey == getRowKey(item) &&
                          editingCell.fieldId == field.id
                        "
                        style="width: 100%"
                      >
                        <component
                          :is="isActiveComponent(field?.type || '')"
                          v-model:editingValue="editingValue"
                          :field="field"
                          @handleCellEdit="handleCellEdit"
                          @handleCancelCellEdit="handleCancelCellEdit"
                        />
                      </div>

                      <!-- 显示模式 -->
                      <template v-else>
                        <DisplayModeComponent :field="field" :item="item" :isEditor="isEditor" />
                      </template>
                    </template>
                  </div>
                </template>
              </div>
            </div>
          </div>
        </div>
        <!-- 可滚动列数据区域 -->
        <div
          class="table-container"
          ref="scrollContainer"
          :style="{ paddingLeft: getFixedColumnsWidth() + 'px' }"
          @scroll="handleScroll"
          @dragleave="handleDragLeave"
        >
          <!-- 虚拟滚动体 -->
          <div
            class="table-body"
            :style="{
              height: tableHeight + 'px',
              width: getScrollableColumnsWidthWithAdd() + 'px',
              minWidth: '100%',
            }"
            ref="bodyContainer"
          >
            <!-- 渲染可滚动列的可见行 -->
            <div
              v-for="item in visibleItems"
              :key="'scrollable-' + getRowKey(item)"
              :data-row-key="getRowKey(item)"
              class="table-row"
              :class="{
                'row-selected': selectedRows.has(getRowKey(item)),
                'row-group-header': item.isGroupHeader,
                'group-last-row': item.isLastInGroup,
                'is-add-row': item.isAddRow,
              }"
              :style="{
                transform: `translateY(${item.translateY}px)`,
                height: itemHeight + 'px',
                width: getScrollableColumnsWidthWithAdd() + 'px',
                borderBottom: collapsedGroups.has(item.groupKey) ? '1px solid #e1e5e9' : '',
              }"
              @click="handleRowClick(item, $event)"
              draggable="true"
              @dragstart="handleRowDragStart($event, item)"
              @dragend="handleRowDragEnd"
              @dragover="handleRowDragOver"
              @drop="handleRowDrop($event, item)"
            >
              <!-- 分组头部 -->
              <template v-if="item.isGroupHeader">
                <!-- <div class="group-header-content">
                  <span class="group-title">{{ item.groupTitle }}</span>
                  <span class="group-count">({{ item.groupCount }})</span>
                </div> -->
              </template>

              <!-- 添加行 -->
              <template v-else-if="item.isAddRow">
                <!-- 合并的可滚动列添加行单元格 -->
                <div
                  class="table-cell scrollable-cell"
                  :style="{
                    width: getScrollableColumnsWidthWithAdd() - 56 + 'px',
                    left: getScrollableColumnLeftByScrollableIndexWithAdd(0) + 'px',
                    background: 'white',
                    borderTop: '1px solid #e5e5e5',
                    borderRight: 'none',
                    borderRadius: '0',
                  }"
                  @click="handleAddRowClick(item.groupKey, $event)"
                  @mouseenter="handleCellMouseEnter(item)"
                  @mouseleave="handleCellMouseLeave(item)"
                ></div>
                <!-- 右侧添加列空白区域 -->
                <div
                  class="table-cell scrollable-cell add-column-cell"
                  :style="{
                    width: '56px',
                    left:
                      getScrollableColumnLeftByScrollableIndexWithAdd(scrollableFields.length) +
                      'px',
                    background: 'white',
                  }"
                  @click="handleAddRowClick(item.groupKey, $event)"
                  @mouseenter="handleCellMouseEnter(item)"
                  @mouseleave="handleCellMouseLeave(item)"
                ></div>
              </template>

              <!-- 可滚动列数据行 -->
              <template v-else>
                <!-- 可滚动列数据 -->
                <div
                  v-for="field in scrollableFields"
                  :key="field.id"
                  class="table-cell"
                  :class="isGrouped ? 'scrollable-cell' : ''"
                  :style="{
                    width: field.width + 'px',
                    left:
                      getScrollableColumnLeftByScrollableIndexWithAdd(
                        scrollableFields.indexOf(field),
                      ) + 'px',
                  }"
                  @click="handleCellClick(item, field, $event)"
                  @dblclick="startCellEdit(item, field)"
                  @mouseenter="handleCellMouseEnter(item)"
                  @mouseleave="handleCellMouseLeave(item)"
                >
                  <!-- 特殊字段不区分显示or编辑模式 -->
                  <template v-if="isSpecialField(field) && item">
                    <component
                      :is="isActiveComponent(field?.type || '')"
                      :field="field"
                      :item="item"
                      :isEditor="isEditor"
                      @handle-special-change="handleSpecialChange"
                    />
                  </template>
                  <template v-else>
                    <!-- 非固定类单元格编辑模式 -->
                    <div
                      style="width: 100%"
                      v-if="
                        isEditor &&
                        editingCell &&
                        editingCell.rowKey == getRowKey(item) &&
                        editingCell.fieldId == field.id
                      "
                    >
                      <component
                        :is="isActiveComponent(field?.type || '')"
                        v-model:editingValue="editingValue"
                        :field="field"
                        :item="item"
                        @handleCellEdit="handleCellEdit"
                        @handleCancelCellEdit="handleCancelCellEdit"
                      />
                    </div>

                    <!-- 显示模式 -->
                    <template v-else>
                      <DisplayModeComponent :field="field" :item="item" :isEditor="isEditor" />
                    </template>
                  </template>
                </div>
              </template>
            </div>
            <!-- 添加列的完整空白块 -->
            <div
              class="table-cell scrollable-cell add-column-cell"
              :style="{
                position: 'absolute',
                zIndex: 2,
                background: 'white',
                top: 0,
                left:
                  getScrollableColumnLeftByScrollableIndexWithAdd(scrollableFields.length) + 'px',
                width: 56 + 'px',
                height: tableHeight + 'px',
              }"
            ></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 性能统计 -->
    <!-- <div class="performance-stats">
      <span>渲染耗时: {{ renderTime }}ms</span>
      <span>纵向滚动: {{ Math.round(scrollTop) }}px</span>
      <span>横向滚动: {{ Math.round(scrollLeft) }}px</span>
      <span>固定列: {{ fixedFields.length }} 个</span>
    </div> -->
  </div>
  <ContextMenu
    :menu-items="menuItems"
    :visible="contextMenuVisible"
    @select="handleMenuItemSelect"
    @close="contextMenuVisible = false"
    ref="contextMenuRef"
  />
</template>

<script setup lang="ts">
import { useMtTableStore } from '@/stores/mtTable'
import { getIconNameByType } from '@/utils'
import {
  getRowKey,
  isActiveComponent,
  isSpecialField,
  handleRowChangeValue,
  isActiveGroupComponent,
} from '@/utils/mtable'
import DisplayModeComponent from '@/components/multiTable/dispalyMode/index.vue'
import { ObjType } from '@/types'
import emitter from '@/utils/mitt'
import { updateCellData, sortField, insertRecordApi, updateFieldApi } from '@/api'
import ContextMenu from '@/components/contextMenu/index.vue'
import type { MenuItem } from '@/components/contextMenu/index.vue' // 导入菜单项类型
import { deleteRecord } from '@/api'

// 监听行高改变
emitter.on('handleRowHeightChange', (data) => {
  tableConfig.rowHeight = handleRowChangeValue(data)
})

const store = useMtTableStore()
const { updateFields, getFields, updateRecords, currentTable, insertRecord, delRecord } = store

interface Props {
  mvTableData: ObjType[]
  fixedFieldsData: ObjType[]
  scrollableFieldsData: ObjType[]
  isEditor?: boolean
  isAssociated?: boolean // 是否关联
  isMultipleSelection?: boolean // 是否允许多选
  selectedKeys?: string[] // 选中的行
}
const props = withDefaults(defineProps<Props>(), {
  mvTableData: () => [],
  fixedFieldsData: () => [],
  scrollableFieldsData: () => [],
  isEditor: true, // 默认可编辑
  isAssociated: false, // 默认不关联
  isMultipleSelection: true, // 默认允许多选
  selectedKeys: () => [],
})

const emit = defineEmits(['handleToggleRowSelection'])
const isAssociated = computed(() => props.isAssociated) // 是否为单双向关联
const isMultipleSelection = computed(() => props.isMultipleSelection) // 是否允许多选

const fields = computed(
  () =>
    (!isAssociated.value
      ? getFields()
      : [...props.fixedFieldsData, ...props.scrollableFieldsData]) || [],
)

const route = useRoute()
// ===== 响应式数据 =====
const tableContainer = ref() // 表格容器
const scrollContainer = ref() // 滚动容器
const scrollableHeaderRef = ref() // 可滚动列表头
const scrollableHeaderContainer = ref() // 可滚动列容器
const fixedBodyContainer = ref() // 固定列容器

const contextMenuRef = ref() // 右键菜单
const contextMenuVisible = ref(false) // 右键菜单是否可见
// 定义菜单数据
const menuItems = ref<MenuItem[]>([])

const isEditor = ref(props.isEditor) // 是否处于编辑状态 false-非编辑状态 true-编辑状态

// 表格基础配置
const itemHeight = ref(handleRowChangeValue(currentTable.settings.tableConfig.rowHeight)) // 每行高度
const containerHeightDefault = 600 // 容器高度
const overscan = 5 // 预渲染行数
const groupSpace = 16 // 分组间隔

// 数据
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const rawData = ref<any[]>([]) // 原始数据
const groupByField = computed(() =>
  isAssociated.value
    ? ''
    : currentTable.settings.tableConfig &&
        currentTable.settings.tableConfig.groupConfig &&
        currentTable.settings.tableConfig.groupConfig.length > 0
      ? currentTable.settings.tableConfig.groupConfig[0]
      : '',
) // 分组字段
const collapsedGroups = ref(new Set()) // 折叠的分组
const selectedRows = ref(props.selectedKeys.length > 0 ? new Set(props.selectedKeys) : new Set()) // 选中的行
const groupOrder = ref<string[]>([]) // 新增：记录分组的原始顺序

// 表格配置选项
const tableConfig = reactive({
  stickyHeader: true, // 表头是否固定
  enableCrossGroupDrag: true, // 是否启用跨分组拖拽
  rowHeight: handleRowChangeValue(currentTable.settings.tableConfig.rowHeight), // 行高，默认为常规高度
})

// 获取固定列和可滚动列
const fixedFields = computed(() => props.fixedFieldsData)
const scrollableFields = computed(() => props.scrollableFieldsData)
// 编辑状态
const editingCell = ref()
const editingValue = ref<string | number>('')

// 滚动状态
const scrollTop = ref(0)
const scrollLeft = ref(0) // 添加横向滚动位置
const renderTime = ref(0)

const recordsData = computed(() => props.mvTableData) // 获取记录数据

// ===== 计算属性 =====
const totalRecords = computed(() => rawData.value.length)

// 检查是否分组
const isGrouped = computed(() => !!groupByField.value)

// 计算滚动条高度
const tableHeight = computed(() => {
  if (!groupByField.value) {
    // 未分组模式下需要包含添加行的高度
    return (totalRecords.value + 1) * itemHeight.value
  } else {
    // 计算分组间距
    const groupSpacing = groupSpace // 每个分组之间的间距
    const groupCount = groupOrder.value.length
    // 包含第一个分组与表头之间的间距
    const totalSpacing = groupCount > 0 ? groupCount * groupSpacing : 0
    return processedData.value.length * itemHeight.value + totalSpacing
  }
})

// 处理分组数据
const processedData = computed(() => {
  const startTime = performance.now()

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let result: any[] = []

  if (!groupByField.value) {
    result = rawData.value.map((item, index) => ({
      ...item,
      originalIndex: index,
      isGroupHeader: false,
      isLastDataRow: index === rawData.value.length - 1, // 标记是否为未分组模式下的最后一行数据
    }))
    // 在未分组模式下添加一个添加行
    result.push({
      isAddRow: true,
      groupKey: 'ungrouped', // 未分组模式下的特殊标识
      originalIndex: -2, // 特殊标记为添加行
      isLastDataRow: false,
    })
  } else {
    // 按字段分组，保持分组顺序稳定
    const groups = new Map()
    console.log('groupByField.value', groupByField.value)
    rawData.value.forEach((item, index) => {
      // 使用专门的函数获取分组值，确保link类型使用linkTitle进行分组
      const groupValue = getGroupValue(item, groupByField.value)

      if (!groups.has(groupValue)) {
        groups.set(groupValue, [])
        // 只有当组不存在于稳定顺序中时才添加

        if (!groupOrder.value.includes(groupValue)) {
          groupOrder.value.push(groupValue)
        }
      }
      // 保留原始数据项，确保子组件可以访问完整信息
      groups.get(groupValue).push({ ...item, originalIndex: index })
    })
    console.log('groupOrder.value----', groupOrder.value)
    console.log('groups', groups)
    // 按稳定顺序生成分组数据
    groupOrder.value.forEach((groupKey, groupIndex) => {
      if (groups.has(groupKey)) {
        const items = groups.get(groupKey)

        // 对于link类型，我们需要获取原始的link对象作为groupKey
        let actualGroupKey = groupKey
        const groupField = fields.value.find((f) => f.id === groupByField.value)
        if (groupField && groupField.type === 'link') {
          const representativeItem = items.find((item) => {
            const itemFieldValue = item[groupByField.value]
            return (
              itemFieldValue &&
              typeof itemFieldValue === 'object' &&
              (itemFieldValue.linkTitle === groupKey || itemFieldValue.linkUrl === groupKey)
            )
          })
          actualGroupKey = representativeItem ? representativeItem[groupByField.value] : groupKey
        }

        if (groupField && groupField.type === 'creator') {
          const representativeItem = items.find((item) => {
            const creatorValue = item['creator']
            return creatorValue && typeof creatorValue === 'object' && creatorValue.id === groupKey
          })
          actualGroupKey = representativeItem ? representativeItem['creator'] : groupKey
        }

        // 添加分组头
        result.push({
          isGroupHeader: true,
          groupKey: actualGroupKey, // 使用实际的link对象作为groupKey
          groupTitle: groupKey, // 显示标题仍使用分组键
          groupCount: items.length,
          originalIndex: -1,
          groupIndex, // 添加分组索引用于计算间距
        })

        // 添加分组数据（如果未折叠）
        if (!collapsedGroups.value.has(groupKey)) {
          const groupItems = items.map((item, index) => ({
            ...item,
            isGroupHeader: false,
            groupKey: actualGroupKey, // 使用实际的link对象作为groupKey
            groupIndex, // 为分组内的项目也添加分组索引
            isLastInGroup: index === items.length - 1, // 标记是否为分组内最后一行
          }))

          // 为分组内最后一行添加特殊标记
          if (groupItems.length > 0) {
            groupItems[groupItems.length - 1].isLastInGroup = true
          }

          result.push(...groupItems)

          // 在每个分组的最后一行后添加一个+"行
          result.push({
            isAddRow: true,
            groupKey: actualGroupKey, // 使用实际的link对象作为groupKey
            groupIndex,
            originalIndex: -2, // 特殊标记为添加行
          })
        }
      }
    })
  }

  // 计算渲染时间并返回结果
  const endTime = performance.now()
  nextTick(() => {
    renderTime.value = Math.round(endTime - startTime)
  })
  console.log('result', result)
  return result
})

// 根据groupByField值反查字段中类型
const getFieldTypeByGroupByField = () => {
  const field = fields.value.find((field) => field.id === groupByField.value)
  return field ? field.type : ''
}

// 获取用于分组的值，对于link类型仅使用linkTitle
const getGroupValue = (item, fieldName) => {
  const fieldValue = item[fieldName]
  const field = fields.value.find((f) => f.id === fieldName)

  console.log('getGroupValue - fieldName:', fieldName)
  console.log('getGroupValue - field:', field)
  console.log('getGroupValue - fieldValue:', fieldValue)
  console.log('getGroupValue - typeof fieldValue:', typeof fieldValue)

  if (
    field &&
    field.type === 'link' &&
    fieldValue &&
    typeof fieldValue === 'object' &&
    fieldValue !== null
  ) {
    return fieldValue.linkTitle || fieldValue.linkUrl || '未分组'
  }

  if (field && field.type === 'creator') {
    const creatorValue = item['creator']
    console.log('getGroupValue - creatorValue:', creatorValue)
    if (creatorValue && typeof creatorValue === 'object' && creatorValue !== null) {
      return creatorValue.id || '未分组'
    }
    return '未分组'
  }

  if (Array.isArray(fieldValue)) {
    return [...fieldValue].sort().join(',')
  }

  return fieldValue ? fieldValue : '未分组'
}

// 添加列点击处理
const handleAddColumnClick = () => {
  emitter.emit('handleAddColumnClick')
}

// 可见项目计算
const visibleItems = computed(() => {
  if (!processedData.value.length) {
    return []
  }

  // 动态获取容器高度 - 使用主滚动容器
  const containerH = scrollContainer.value
    ? scrollContainer.value.clientHeight
    : containerHeightDefault

  // 计算带间距的项目位置
  const getItemPosition = (index) => {
    if (!groupByField.value) {
      return index * itemHeight.value
    }

    let position = 0
    const groupSpacing = groupSpace

    for (let i = 0; i < index; i++) {
      const item = processedData.value[i]
      position += itemHeight.value

      // 如果当前项是分组头部，添加间距（包括第一个分组）
      if (item.isGroupHeader) {
        position += groupSpacing
      }
    }

    return position
  }

  // 根据滚动位置计算起始索引
  let startIndex = 0
  for (let i = 0; i < processedData.value.length; i++) {
    if (getItemPosition(i + 1) > scrollTop.value) {
      startIndex = i
      break
    }
  }

  // 计算结束索引
  let endIndex = processedData.value.length
  for (let i = startIndex; i < processedData.value.length; i++) {
    if (getItemPosition(i) > scrollTop.value + containerH + overscan * itemHeight.value) {
      endIndex = i
      break
    }
  }

  const visibleStart = Math.max(0, startIndex - overscan)
  const visibleEnd = Math.min(endIndex + overscan, processedData.value.length)

  //   console.log(
  //     'update',
  //     processedData.value.slice(visibleStart, visibleEnd).map((item, index) => {
  //       const actualIndex = visibleStart + index
  //       return {
  //         ...item,
  //         translateY: getItemPosition(actualIndex),
  //       }
  //     }),
  //   )
  const visibleData = processedData.value.slice(visibleStart, visibleEnd).map((item, index) => {
    const actualIndex = visibleStart + index
    return {
      ...item,
      translateY: getItemPosition(actualIndex),
    }
  })

  return visibleData
})
// 获取固定列总宽度
const getFixedColumnsWidth = () => {
  let width = 40 // 选择器宽度
  fixedFields.value.forEach((field) => {
    width += field.width as number
  })
  return width
}
// 处理菜单项选择
const handleMenuItemSelect = async (item: MenuItem) => {
  console.log('选择了菜单项:', item)
  switch (item.key) {
    case 'delete':
      const rowIds = Array.from(selectedRows.value) as string[]
      await deleteRecord({
        recordIds: rowIds,
      })
      delRecord(rowIds)
      // 删除后重置选中状态
      clearSelection()
      break
    default:
      break
  }
}

// 处理右键事件
const handleContextMenu = (e: MouseEvent) => {
  e.preventDefault()
  if (selectedRows.value.size === 0) return
  // 设置菜单位置
  if (contextMenuRef.value) {
    contextMenuRef.value.setPosition(e.clientX, e.clientY)
  }
  menuItems.value = [
    {
      key: 'delete',
      label: selectedRows.value.size > 0 ? `删除 ${selectedRows.value.size} 行` : '删除',
      iconName: 'delete',
    }, // 禁用状态测试
  ]
  // 显示菜单
  contextMenuVisible.value = true
}

// 获取包含添加列的可滚动列总宽度
const getScrollableColumnsWidthWithAdd = () => {
  let width = 0
  scrollableFields.value.forEach((field) => {
    width += field.width as number
  })
  // 添加列宽度
  width += 56
  return width
}

// 获取固定列在容器中的位置
const getFixedColumnLeft = (fieldIndex) => {
  let left = 40 // 选择器宽度
  for (let i = 0; i < fieldIndex; i++) {
    if (fields.value[i].fixed) {
      left += fields.value[i].width as number
    }
  }
  return left
}

// 获取固定列的sticky left值（用于position: sticky）
const getFixedColumnStickyLeft = (fieldIndex) => {
  const field = fields.value[fieldIndex]
  if (!field.fixed) return 0

  let left = 40 // 选择器宽度
  for (let i = 0; i < fieldIndex; i++) {
    if (fields.value[i].fixed) {
      left += fields.value[i].width as number
    }
  }
  return left
}

// 获取包含添加列的可滚动列在其区域内的位置
const getScrollableColumnLeftByScrollableIndexWithAdd = (scrollableIndex) => {
  let left = 0
  for (let i = 0; i < scrollableIndex; i++) {
    // 使用整数像素值避免小数点精度问题
    left += Math.round(scrollableFields.value[i].width as number)
  }
  return Math.round(left)
}

// 分离架构的滚动处理
// 统一滚动架构的滚动处理
const handleScroll = (e) => {
  scrollTop.value = e.target.scrollTop
  scrollLeft.value = e.target.scrollLeft

  // 使用 requestAnimationFrame 优化性能
  requestAnimationFrame(() => {
    // 同步可滚动表头的横向滚动
    const scrollLeftInt = Math.round(scrollLeft.value)

    if (scrollableHeaderRef.value) {
      const headerContainer = scrollableHeaderRef.value.querySelector(
        '.scrollable-header-container',
      )
      if (headerContainer) {
        headerContainer.style.transform = `translateX(-${scrollLeftInt}px)`
      }
    }

    // 同步固定列的纵向滚动
    if (fixedBodyContainer.value && e.target !== fixedBodyContainer.value) {
      fixedBodyContainer.value.scrollTop = scrollTop.value
    }

    // 动态调整固定列容器的高度以补偿滚动条
    if (fixedBodyContainer.value && scrollContainer.value) {
      const scrollBarHeight =
        scrollContainer.value.offsetHeight - scrollContainer.value.clientHeight
      fixedBodyContainer.value.style.height = `calc(100% - ${scrollBarHeight}px)`
    }
  })
}

// 固定列的滚动处理
const handleFixedScroll = (e) => {
  scrollTop.value = e.target.scrollTop

  // 同步主滚动容器的纵向滚动
  requestAnimationFrame(() => {
    if (scrollContainer.value && e.target !== scrollContainer.value) {
      scrollContainer.value.scrollTop = scrollTop.value
    }

    // 动态调整固定列容器的高度以补偿滚动条
    if (fixedBodyContainer.value && scrollContainer.value) {
      const scrollBarHeight =
        scrollContainer.value.offsetHeight - scrollContainer.value.clientHeight
      fixedBodyContainer.value.style.height = `calc(100% - ${scrollBarHeight}px)`
    }
  })
}

// 更新虚拟列表 - 用于数据变化后强制更新虚拟滚动
const updateVirtualList = () => {
  // 通过微小改变滚动位置来触发虚拟列表的重新计算
  if (scrollContainer.value) {
    const currentScrollTop = scrollContainer.value.scrollTop
    const currentScrollLeft = scrollContainer.value.scrollLeft

    // 先保存当前滚动位置
    setTimeout(() => {
      if (scrollContainer.value) {
        // 恢复原始滚动位置并触发滚动事件
        scrollContainer.value.scrollTo(currentScrollLeft, currentScrollTop)
        handleScroll({ target: scrollContainer.value })
      }
    }, 0)
  }
}

// 行点击
const handleRowClick = (item, event) => {
  if (item.isGroupHeader || item.isAddRow) return

  if (event.ctrlKey || event.metaKey) {
    toggleRowSelection(item)
  } else {
    clearSelection()
    selectedRows.value.add(getRowKey(item))
  }
}

// 添加行点击
const handleAddRowClick = async (groupKey, event) => {
  if (event) {
    event.stopPropagation()
  }
  if (groupKey === 'ungrouped') {
    console.log('在未分组模式下添加新行')
    const fields = getFields()
    const newRecord = fields.reduce((acc, cur) => {
      acc[cur.id] = cur.defaultValue
      return acc
    }, {})
    console.log('newRecord', newRecord)
    const response = await insertRecordApi({
      prevId: rawData.value[rawData.value.length - 1]?.rowId,
      tableSchemaId: route.params.tableSchemaId as string,
    })
    insertRecord({ ...newRecord, ...response }, 'end')
  } else {
    console.log(`在分组 ${groupKey} 中添加新行`)
    // 在当前分组中添加新行,新行内容为当前分组最后一条数据
    // 找到当前分组的所有原始数据项
    const groupItems = rawData.value.filter((item) => {
      const itemGroupValue = item[groupByField.value] || '未分组'
      return itemGroupValue === groupKey
    })

    let newRecord
    let lastItemId
    if (groupItems.length > 0) {
      // 复制最后一个数据项
      const lastItem = { ...groupItems[groupItems.length - 1] }
      lastItemId = lastItem.rowId
      newRecord = lastItem
    } else {
      // 如果该分组没有数据，则创建默认值的新行
      const fields = getFields()
      newRecord = fields.reduce((acc, cur) => {
        acc[cur.id] = cur.defaultValue
        return acc
      }, {})
    }
    // 设置分组字段值（确保新行属于当前分组）
    newRecord[groupByField.value] = groupKey
    const response = await insertRecordApi({
      prevId: lastItemId,
      tableSchemaId: route.params.tableSchemaId as string,
      rowData: {
        [groupByField.value]: groupKey,
      },
    })
    // 插入新记录（可选：插入到分组末尾，这里统一用 'end'）
    insertRecord({ ...newRecord, ...response }, 'end')
  }
}

// 切换行选择
const toggleRowSelection = (item) => {
  // 如果 isMultipleSelection 为 false，则只允许单选
  if (!isMultipleSelection.value) {
    clearSelection()
  }
  // 关闭右键菜单
  if (contextMenuVisible.value) {
    contextMenuVisible.value = false
  }
  const key = getRowKey(item)
  if (selectedRows.value.has(key)) {
    selectedRows.value.delete(key)
  } else {
    selectedRows.value.add(key)
  }
  emit('handleToggleRowSelection', selectedRows.value)
  console.log('selectedRows.value', selectedRows.value)
  console.log('切换行选择:', key, '当前选中:', selectedRows.value.size)
}

// 切换分组
const toggleGroup = (groupKey) => {
  if (collapsedGroups.value.has(groupKey)) {
    collapsedGroups.value.delete(groupKey)
  } else {
    collapsedGroups.value.add(groupKey)
  }
}

// 单元格点击
const handleCellClick = (item, field, event) => {
  if (item.isGroupHeader) return
  event.stopPropagation()
}

// 单元格鼠标进入
const handleCellMouseEnter = (item) => {
  if (item.isGroupHeader) return
  // 为整行添加高亮样式
  const rowKey = getRowKey(item)
  const fixedRow = document.querySelector(`.fixed-row[data-row-key="fixed-${rowKey}"]`)
  const scrollableRow = document.querySelector(`.table-row[data-row-key="${rowKey}"]`)

  if (fixedRow) {
    fixedRow.classList.add('row-hover')
  }
  if (scrollableRow) {
    scrollableRow.classList.add('row-hover')
  }
}

// 单元格鼠标离开
const handleCellMouseLeave = (item) => {
  if (item.isGroupHeader) return
  // 移除整行高亮样式
  const rowKey = getRowKey(item)
  const fixedRow = document.querySelector(`.fixed-row[data-row-key="fixed-${rowKey}"]`)
  const scrollableRow = document.querySelector(`.table-row[data-row-key="${rowKey}"]`)

  if (fixedRow) {
    fixedRow.classList.remove('row-hover')
  }
  if (scrollableRow) {
    scrollableRow.classList.remove('row-hover')
  }
}
// 特殊字段不区分显示or编辑模式回调方法
const handleSpecialFieldChange = (obj: ObjType) => {
  const { field, item, value } = obj

  editingCell.value = {
    rowKey: getRowKey(item),
    fieldId: field.id,
  }
  console.log('editingCell.value', editingCell.value)
  if (field.type == 'checkbox') {
    editingValue.value = !value ? 0 : 1
  } else {
    editingValue.value = value
  }
  finishCellEdit()
}
// 复选框类型字段、图片类型字段、附件类型字段、评分类型字段改变回调事件
const handleSpecialChange = (obj: ObjType) => {
  console.log('obj', obj)
  handleSpecialFieldChange(obj)
}

// 开始编辑单元格
const startCellEdit = async (item, field) => {
  if (item.isGroupHeader) return
  if (['autoNumber', 'createTime', 'creator'].includes(field.type)) return

  editingCell.value = {
    rowKey: getRowKey(item),
    fieldId: field.id,
  }
  if (field.type == 'person') {
    console.log('field', field)
    console.log('item[field.id]', item[field.id])
    editingValue.value = field.settings?.allowMultiple ? item[field.id] || [] : [item[field.id]]
    console.log('editingValue.value+++++', editingValue.value)
  } else {
    editingValue.value = item[field.id] || ''
  }
  console.log('editingValue.value', editingValue.value)
  await nextTick()
}

// 编辑组件单元格编辑
const handleCellEdit = () => {
  console.log('更改后---editingValue.value', editingValue.value)
  finishCellEdit()
}

// 编辑组件单元格取消编辑
const handleCancelCellEdit = () => {
  cancelCellEdit()
}

// 单元格完成编辑
const finishCellEdit = async () => {
  if (!editingCell.value) return

  const item = processedData.value.find((item) => getRowKey(item) == editingCell.value.rowKey)
  if (item && !item.isGroupHeader) {
    item[editingCell.value.fieldId] = editingValue.value
    // 更新原始数据
    const originalItem = rawData.value.find((d) => d.rowId == item.rowId)
    if (originalItem) {
      originalItem[editingCell.value.fieldId] = editingValue.value
      await updateCellData({
        fieldId: editingCell.value.fieldId,
        recordId: item.rowId,
        value: editingValue.value,
      })
      updateRecords(rawData.value) // 更新records数据
    }
  }

  editingCell.value = null
  editingValue.value = ''
}

// 单元格取消编辑
const cancelCellEdit = () => {
  editingCell.value = null
  editingValue.value = ''
}

// 列宽调整
let resizing = false
let resizeStartX = 0
let resizeField = null
let initialWidth = 0 // 记录初始宽度

const startColumnResize = (event, field) => {
  event.preventDefault()
  event.stopPropagation()

  resizing = true
  resizeStartX = event.clientX
  resizeField = field
  initialWidth = field.width // 记录开始时的宽度
  console.log('resizeField', resizeField)
  document.addEventListener('mousemove', handleColumnResize)
  document.addEventListener('mouseup', endColumnResize)
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'
}

const handleColumnResize = (event) => {
  if (!resizing || !resizeField) return

  const deltaX = event.clientX - resizeStartX
  const newWidth = Math.max(60, initialWidth + deltaX) // 使用初始宽度+偏移量
  if (resizeField) (resizeField as ObjType).width = newWidth as number
}

// 列宽调整结束
const endColumnResize = async () => {
  if (isAssociated.value) return // 单双向关联下禁止拖拽
  document.removeEventListener('mousemove', handleColumnResize)
  document.removeEventListener('mouseup', endColumnResize)
  document.body.style.cursor = ''
  document.body.style.userSelect = ''

  //   获取拖动的具体字段
  const fieldObj = resizeField
    ? (fields.value.find((field) => field.id === resizeField.id) as ObjType | undefined)
    : undefined

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  await updateFieldApi({ ...(fieldObj as any) }) // 更新字段

  //   本地更新 width
  updateFields(fields.value as any) // eslint-disable-line @typescript-eslint/no-explicit-any

  resizing = false
  resizeField = null // 重置 resizeField
  initialWidth = 0
}

// 列拖拽排序
let draggedIndex = -1
let draggedRowItem = null

//  列拖拽排序 拖拽开始
const handleHeaderDragStart = (event) => {
  const index = parseInt(event.target.dataset.index)
  if (index === 0) {
    event.preventDefault() // 完全阻止拖拽
    return
  }
  draggedIndex = index
  event.dataTransfer.effectAllowed = 'move'
}
// 列拖拽排序 拖拽结束
const handleHeaderDragOver = (event) => {
  event.preventDefault()
  event.dataTransfer.dropEffect = 'move'
}
//  列拖拽
const handleHeaderDrop = async (event) => {
  if (isAssociated.value) return // 单双向关联下禁止拖拽
  event.preventDefault()
  const dropIndex = parseInt(event.target.closest('.column-header').dataset.index)
  console.log('dropIndex', dropIndex)
  //   第一列禁止拖拽
  if (dropIndex == 0) return

  if (draggedIndex !== -1 && draggedIndex !== dropIndex) {
    const draggedField = fields.value.splice(draggedIndex, 1)[0]
    fields.value.splice(dropIndex, 0, draggedField)
    console.log('fields.value', fields.value)
    await sortField({
      fieldIds: fields.value.map((field) => field.id),
      tableSchemaId: route.params.tableSchemaId as string,
    })
    updateFields(fields.value as any) // eslint-disable-line @typescript-eslint/no-explicit-any
  }
  draggedIndex = -1
}

// 行拖拽排序
// 拖拽相关变量
let dropIndicator: HTMLElement | null = null
let dragStartIndex = -1

// 更新行高
const updateRowHeight = () => {
  itemHeight.value = Number(tableConfig.rowHeight)
  // 更新虚拟列表以适应新的行高
  updateVirtualList()
}

const handleRowDragStart = (event, item) => {
  if (isAssociated.value) return // 单双向关联下禁止拖拽
  if (item.isGroupHeader) {
    event.preventDefault()
    return
  }

  draggedRowItem = item
  dragStartIndex = rawData.value.findIndex((data) => data.rowId === item.rowId)
  event.dataTransfer.effectAllowed = 'move'

  // 只创建插入指示器，不创建拖拽预览
  if (!dropIndicator) {
    dropIndicator = document.createElement('div')
    dropIndicator.className = 'drop-indicator'
    dropIndicator.style.cssText = `
      position: absolute;
      left: 0;
      height: 3px;
      background: linear-gradient(90deg, #007bff, #0056b3);
      border-radius: 2px;
      z-index: 1000;
      pointer-events: none;
      display: none;
      box-shadow: 0 0 8px rgba(0,123,255,0.4);
    `
    // 将标示线添加到表格主体容器中，让它覆盖整个表格宽度（包括固定列）
    const tableWrapper = document.querySelector('.table-body-wrapper') as HTMLElement
    if (tableWrapper) {
      // 设置标示线宽度为整个表格区域宽度
      const tableWidth = tableWrapper.offsetWidth
      dropIndicator.style.width = `${tableWidth}px`
      tableWrapper.appendChild(dropIndicator)
    }
  }

  console.log('开始拖拽:', item.name, '当前分组:', item.groupKey || '无分组')
}

const handleRowDragOver = (event) => {
  event.preventDefault()
  event.dataTransfer.dropEffect = 'move'

  if (!draggedRowItem || !dropIndicator) return

  const targetRow = event.currentTarget
  const targetItem = visibleItems.value.find(
    (item) => getRowKey(item) === targetRow.getAttribute('data-row-key'),
  )

  // 清除之前的高亮
  document.querySelectorAll('.table-row').forEach((row) => {
    row.classList.remove('drag-over', 'cross-group-target', 'same-group-target')
  })

  if (targetItem && (targetItem as any)?.rowId !== (draggedRowItem as any)?.rowId) {
    const targetRect = targetRow.getBoundingClientRect()

    // 计算插入位置（上半部还是下半部）
    const relativeY = event.clientY - targetRect.top
    const insertBefore = relativeY < targetRect.height / 2

    if (targetItem.isGroupHeader && tableConfig.enableCrossGroupDrag) {
      targetRow.classList.add('drag-over')
      dropIndicator.style.display = 'none'
    } else if (!targetItem.isGroupHeader) {
      const isCrossGroup =
        tableConfig.enableCrossGroupDrag &&
        (draggedRowItem as any)?.groupKey !== (targetItem as any)?.groupKey // eslint-disable-line @typescript-eslint/no-explicit-any
      targetRow.classList.add(isCrossGroup ? 'cross-group-target' : 'same-group-target')

      // 虚拟滚动中的标示线位置计算
      // 使用目标项在 visibleItems 中的 translateY 值来计算正确位置
      const targetItemIndex = visibleItems.value.findIndex(
        (item) => getRowKey(item) === getRowKey(targetItem),
      )
      if (targetItemIndex !== -1) {
        const targetVisibleItem = visibleItems.value[targetItemIndex]
        const indicatorTop = insertBefore
          ? targetVisibleItem.translateY
          : targetVisibleItem.translateY + itemHeight

        // 确保标示线位置相对于整个表格容器
        const tableBodyWrapper = document.querySelector('.table-body-wrapper')
        if (tableBodyWrapper) {
          // 获取滚动容器的滚动位置
          const scrollOffset = scrollTop.value
          // 计算标示线相对于表格主体容器的绝对位置
          dropIndicator.style.top = `${indicatorTop - scrollOffset}px`
          dropIndicator.style.display = 'block'
        }
      }
    }
  } else {
    dropIndicator.style.display = 'none'
  }
}

const handleRowDrop = (event, targetItem) => {
  event.preventDefault()

  if (!draggedRowItem) {
    return
  }

  // 清理拖拽状态
  cleanupDragState()

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  if ((draggedRowItem as any)?.rowId === (targetItem as any)?.rowId) {
    draggedRowItem = null
    return
  }

  // 计算插入位置
  const targetRect = event.currentTarget.getBoundingClientRect()
  const relativeY = event.clientY - targetRect.top
  const insertBefore = relativeY < targetRect.height / 2

  // 如果是分组头部，处理跨分组移动
  if (targetItem.isGroupHeader && tableConfig.enableCrossGroupDrag) {
    handleCrossGroupMove(draggedRowItem, targetItem.groupKey)
  }
  // 如果是数据行，处理同分组内移动或跨分组移动
  else if (!targetItem.isGroupHeader) {
    if (
      tableConfig.enableCrossGroupDrag ||
      (draggedRowItem as any)?.groupKey === (targetItem as any)?.groupKey
    ) {
      handleDataRowMove(draggedRowItem, targetItem, insertBefore)
    }
  }

  draggedRowItem = null
}

// 清理拖拽状态
const cleanupDragState = () => {
  // 清除所有拖拽样式
  document.querySelectorAll('.table-row').forEach((row) => {
    ;(row as HTMLElement).style.opacity = ''
    row.classList.remove('dragging', 'drag-over', 'cross-group-target', 'same-group-target')
  })

  // 隐藏插入指示器
  if (dropIndicator) {
    dropIndicator.style.display = 'none'
  }
}

// 拖拽结束处理
const handleRowDragEnd = () => {
  // 清理拖拽状态
  cleanupDragState()
  draggedRowItem = null
  dragStartIndex = -1

  console.log('拖拽结束')
}

// 监听鼠标事件，确保在拖拽过程中意外结束时清理状态
const handleGlobalMouseUp = () => {
  if (draggedRowItem) {
    cleanupDragState()
    draggedRowItem = null
    dragStartIndex = -1
  }
}

// 监听拖拽离开页面
const handleDragLeave = (event) => {
  // 如果离开了容器区域，清理状态
  if (
    !event.relatedTarget ||
    !document.querySelector('.table-container')?.contains(event.relatedTarget)
  ) {
    if (dropIndicator) {
      dropIndicator.style.display = 'none'
    }
    document.querySelectorAll('.table-row').forEach((row) => {
      row.classList.remove('drag-over', 'cross-group-target', 'same-group-target')
    })
  }
}

// 处理跨分组移动（移动到分组头部）
const handleCrossGroupMove = (draggedItem, targetGroupKey) => {
  const draggedOriginalIndex = rawData.value.findIndex((item) => item.rowId === draggedItem.rowId)

  if (draggedOriginalIndex !== -1) {
    const draggedData = rawData.value[draggedOriginalIndex]

    // 只更新被拖拽行的分组字段，不移动数据位置
    if (groupByField.value && targetGroupKey !== '未分组') {
      const oldGroupValue = draggedData[groupByField.value]
      draggedData[groupByField.value] = targetGroupKey

      // 确保目标分组在稳定顺序中
      if (!groupOrder.value.includes(targetGroupKey)) {
        // 如果是新分组，添加到原始分组位置附近
        const sourceGroupIndex = groupOrder.value.indexOf(oldGroupValue)
        if (sourceGroupIndex !== -1) {
          groupOrder.value.splice(sourceGroupIndex + 1, 0, targetGroupKey)
        } else {
          groupOrder.value.push(targetGroupKey)
        }
      }

      console.log(`跨分组移动: ${draggedData.name} 从 ${oldGroupValue} 移动到 ${targetGroupKey}`)
      console.log('分组顺序保持:', groupOrder.value)
    }
  }
}

// 处理数据行间移动
const handleDataRowMove = (draggedItem, targetItem, insertBefore = true) => {
  console.log('rawData.value', rawData.value)
  const draggedOriginalIndex = rawData.value.findIndex((item) => item.rowId === draggedItem.rowId)
  const targetOriginalIndex = rawData.value.findIndex((item) => item.rowId === targetItem.rowId)

  if (draggedOriginalIndex === -1 || targetOriginalIndex === -1) {
    console.warn('无法找到拖拽或目标项')
    return
  }

  // 防止拖拽到自己
  if (draggedOriginalIndex === targetOriginalIndex) {
    console.log('无法拖拽到自己')
    return
  }

  // 如果启用了跨分组拖拽，且目标行有分组信息
  if (
    tableConfig.enableCrossGroupDrag &&
    targetItem.groupKey &&
    draggedItem.groupKey !== targetItem.groupKey
  ) {
    const draggedData = rawData.value[draggedOriginalIndex]
    // 更新分组字段
    if (groupByField.value) {
      const oldGroupValue = draggedData[groupByField.value]
      draggedData[groupByField.value] = targetItem[groupByField.value]

      // 确保目标分组在稳定顺序中
      const newGroupValue = targetItem[groupByField.value]
      if (!groupOrder.value.includes(newGroupValue)) {
        const sourceGroupIndex = groupOrder.value.indexOf(oldGroupValue)
        if (sourceGroupIndex !== -1) {
          groupOrder.value.splice(sourceGroupIndex + 1, 0, newGroupValue)
        } else {
          groupOrder.value.push(newGroupValue)
        }
      }

      console.log(
        `跨分组移动: ${draggedData.name} 从 ${oldGroupValue} 移动到 ${targetItem.groupKey}`,
      )
    }
  } else {
    // 同分组内排序，实现真正的位置交换
    const draggedData = rawData.value[draggedOriginalIndex]

    // 移除拖拽项
    rawData.value.splice(draggedOriginalIndex, 1)

    // 重新计算目标位置（因为移除了一项）
    let newTargetIndex = targetOriginalIndex
    if (draggedOriginalIndex < targetOriginalIndex) {
      newTargetIndex = targetOriginalIndex - 1
    }

    // 根据插入位置调整
    if (!insertBefore) {
      newTargetIndex += 1
    }

    // 边界检查
    newTargetIndex = Math.max(0, Math.min(rawData.value.length, newTargetIndex))

    // 插入到新位置
    rawData.value.splice(newTargetIndex, 0, draggedData)

    console.log(
      `同分组排序: ${draggedData.name} 从索引 ${draggedOriginalIndex} 移动到索引 ${newTargetIndex}`,
    )
  }

  // 触发更新
  nextTick(() => {
    // 刷新虚拟滚动，确保排序后立即显示新顺序
    updateVirtualList()

    // 如果拖拽项超出当前可视区域，滚动到相应位置
    if (dragStartIndex !== -1) {
      const newIndex = rawData.value.findIndex((item) => item.rowId === draggedItem.rowId)
      if (newIndex !== -1) {
        const targetY = newIndex * itemHeight.value
        const containerHeight = scrollContainer.value?.clientHeight || containerHeightDefault
        const currentScrollTop = scrollTop.value

        // 如果新位置不在可视区域内，滚动到相应位置
        if (targetY < currentScrollTop || targetY > currentScrollTop + containerHeight) {
          scrollContainer.value?.scrollTo({
            top: Math.max(0, targetY - containerHeight / 2),
            behavior: 'smooth',
          })
        }
      }
    }
  })
}

// 批量操作
const batchDelete = () => {
  if (selectedRows.value.size === 0) return

  if (confirm(`确定要删除选中的 ${selectedRows.value.size} 行数据吗？`)) {
    // 获取选中的ID列表
    const selectedIds = new Set()
    selectedRows.value.forEach((rowKey: unknown) => {
      if (rowKey.startsWith('row_')) {
        const id = parseInt(rowKey.replace('row_', ''))
        if (!isNaN(id)) {
          selectedIds.add(id)
        }
      }
    })

    // 从原始数据中删除
    rawData.value = rawData.value.filter((item) => !selectedIds.has(item.id))

    // 清空选中状态
    clearSelection()

    console.log(`已删除 ${selectedIds.size} 行数据`)
  }
}

// 清空选中状态
const clearSelection = () => {
  selectedRows.value.clear()
}

// 全选/取消全选
const toggleSelectAll = () => {
  // 关闭右键菜单
  if (contextMenuVisible.value) {
    contextMenuVisible.value = false
  }
  const dataRows = processedData.value.filter((item) => !item.isGroupHeader)
  if (selectedRows.value.size === dataRows.length && dataRows.length > 0) {
    // 当前全选，取消全选
    clearSelection()
    console.log('取消全选')
  } else {
    // 全选
    clearSelection()
    dataRows.forEach((item) => {
      selectedRows.value.add(getRowKey(item))
    })
    emit('handleToggleRowSelection', selectedRows.value)
    console.log('全选，共', dataRows.length, '行')
  }
}

const isAllSelected = computed(() => {
  if (!processedData.value || processedData.value.length === 0) return false
  const dataRows = processedData.value.filter((item) => !item.isGroupHeader && !item.isAddRow)
  return dataRows.length > 0 && selectedRows.value.size === dataRows.length
})

const isPartiallySelected = computed(() => {
  if (!processedData.value || processedData.value.length === 0) return false
  const dataRows = processedData.value.filter((item) => !item.isGroupHeader && !item.isAddRow)
  return (
    selectedRows.value.size > 0 && selectedRows.value.size < dataRows.length && dataRows.length > 0
  )
})
const exportData = () => {
  const dataToExport = rawData.value
  const csv = convertToCSV(dataToExport)
  downloadCSV(csv, 'table_data.csv')
}

const convertToCSV = (data) => {
  const headers = fields.value.map((f) => f.title).join(',')
  const rows = data
    .map((row) => fields.value.map((f) => `"${row[f.id] || ''}"`).join(','))
    .join('\n')
  return headers + '\n' + rows
}

const downloadCSV = (csv, filename) => {
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = filename
  link.click()
}

// ===== 生命周期 =====
onMounted(() => {
  rawData.value = recordsData.value.length > 0 ? recordsData.value : [] // 默认生成数据
  // 添加全局鼠标事件监听
  document.addEventListener('mouseup', handleGlobalMouseUp)
})

// 监听表头固定配置变化
watch(
  () => tableConfig.stickyHeader,
  () => {
    nextTick(() => {
      if (scrollContainer.value) {
        handleScroll({ target: scrollContainer.value })
      }
    })
  },
)

// 监听分组字段变化，重置分组顺序
watch(
  () => groupByField.value,
  () => {
    groupOrder.value = [] // 重置分组顺序
  },
)

// 监听数据变化，更新原始数据
watch(
  () => recordsData.value,
  () => {
    rawData.value = recordsData.value
    updateVirtualList()
  },
)

// 监听行高变化
watch(
  () => tableConfig.rowHeight,
  () => {
    updateRowHeight()
  },
)
// 组件销毁时清理资源
onUnmounted(() => {
  document.removeEventListener('mouseup', handleGlobalMouseUp)

  // 清理拖拽指示器
  if (dropIndicator && dropIndicator.parentNode) {
    dropIndicator.parentNode.removeChild(dropIndicator)
  }

  if (resizing) {
    endColumnResize()
  }
})
// 为组件设置多词名称，避免 ESLint 警告
defineOptions({
  name: 'MultiTable',
})
</script>

<style>
@import './index.css';
</style>
