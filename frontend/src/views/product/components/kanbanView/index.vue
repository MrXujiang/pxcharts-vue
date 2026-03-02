<template>
  <div class="kanban-view">
    <!-- 工具栏 -->
    <div class="kanban-toolbar">
      <div class="toolbar-left">
        <t-select
          v-model="selectedGroupField"
          :options="groupFieldOptions"
          placeholder="选择分组字段"
          style="width: 200px"
          @change="handleGroupFieldChange"
        >
          <template #prefixIcon>
            <t-icon name="layers" />
          </template>
        </t-select>
        <t-button variant="outline" @click="handleRefresh">
          <template #icon>
            <t-icon name="refresh" />
          </template>
          刷新
        </t-button>
      </div>
      <div class="toolbar-right">
        <span class="total-count">共 {{ totalCards }} 条记录</span>
      </div>
    </div>

    <!-- 看板主体 -->
    <div class="kanban-container">
      <div class="kanban-scroll">
        <KanbanColumn
          v-for="group in groups"
          :key="group.value"
          :group="group"
          :cards="getGroupCards(group.value)"
          :fields="fields"
          :title-field="titleField"
          @card-move="handleCardMove"
          @add-card="handleAddCard"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { useMtTableStore } from '@/stores'
import KanbanColumn from './KanbanColumn.vue'

defineOptions({
  name: 'KanbanView',
})

const tableStore = useMtTableStore()

// 当前选中的分组字段
const selectedGroupField = ref('department')

// 标题字段
const titleField = ref('name')

// 分组顺序缓存（保持分组顺序不变）
const groupOrderCache = ref<Map<string, number>>(new Map())

// 表格字段
const fields = computed(() => tableStore.getFields())

// 表格数据（本地副本，用于拖拽排序）
const kanbanData = ref<Array<Record<string, any>>>([])

// 初始化数据
const initData = () => {
  kanbanData.value = JSON.parse(JSON.stringify(tableStore.getRecords()))
}

initData()

// 可用于分组的字段选项
const groupFieldOptions = computed(() => {
  return fields.value
    .filter((field) => 
      field.type === 'text' || 
      field.type === 'select' || 
      !field.type
    )
    .map((field) => ({
      label: field.title,
      value: field.id,
    }))
})

// 获取所有分组
const groups = computed(() => {
  const groupMap = new Map<any, { value: any; label: string; order: number }>()
  
  // 收集所有分组值
  kanbanData.value.forEach((record) => {
    const value = record[selectedGroupField.value]
    if (value !== null && value !== undefined && value !== '') {
      const key = String(value)
      if (!groupMap.has(value)) {
        // 如果缓存中有顺序，使用缓存的顺序；否则使用当前最大顺序+1
        const order = groupOrderCache.value.has(key) 
          ? groupOrderCache.value.get(key)! 
          : groupOrderCache.value.size
        
        groupMap.set(value, {
          value,
          label: key,
          order,
        })
        
        // 更新缓存
        if (!groupOrderCache.value.has(key)) {
          groupOrderCache.value.set(key, order)
        }
      }
    }
  })

  // 添加"未分组"
  const ungroupKey = '__ungroup__'
  const ungroupOrder = groupOrderCache.value.has(ungroupKey)
    ? groupOrderCache.value.get(ungroupKey)!
    : groupOrderCache.value.size
  
  groupMap.set('__ungroup__', {
    value: '__ungroup__',
    label: '未分组',
    order: ungroupOrder,
  })
  
  if (!groupOrderCache.value.has(ungroupKey)) {
    groupOrderCache.value.set(ungroupKey, ungroupOrder)
  }

  // 按顺序排序后返回
  return Array.from(groupMap.values())
    .sort((a, b) => a.order - b.order)
    .map(({ value, label }) => ({ value, label }))
})

// 获取指定分组的卡片
const getGroupCards = (groupValue: any) => {
  return kanbanData.value.filter((record) => {
    const value = record[selectedGroupField.value]
    
    // 处理未分组
    if (groupValue === '__ungroup__') {
      return value === null || value === undefined || value === ''
    }
    
    return value === groupValue
  })
}

// 总卡片数
const totalCards = computed(() => kanbanData.value.length)

// 处理分组字段变更
const handleGroupFieldChange = () => {
  // 清空分组顺序缓存，让新字段的分组重新排序
  groupOrderCache.value.clear()
  MessagePlugin.success(`已切换到 ${fields.value.find(f => f.id === selectedGroupField.value)?.title} 分组`)
}

// 处理卡片移动
const handleCardMove = (event: {
  fromGroup: any
  toGroup: any
  oldIndex: number
  newIndex: number
  card: Record<string, any>
}) => {
  const { fromGroup, toGroup, card } = event

  // 如果跨分组移动，更新卡片的分组字段值
  if (fromGroup !== toGroup) {
    const cardIndex = kanbanData.value.findIndex((c) => c.rowId === card.rowId)
    if (cardIndex !== -1) {
      // 更新分组字段的值
      if (toGroup === '__ungroup__') {
        kanbanData.value[cardIndex][selectedGroupField.value] = ''
      } else {
        kanbanData.value[cardIndex][selectedGroupField.value] = toGroup
      }
    }
  }

  // 同步到 store
  tableStore.updateRecords(kanbanData.value)
  
  MessagePlugin.success('卡片移动成功')
}

// 处理添加卡片
const handleAddCard = (groupValue: any) => {
  MessagePlugin.info(`添加卡片到 ${groupValue} 分组`)
  // TODO: 实现添加卡片逻辑
}

// 刷新数据
const handleRefresh = () => {
  initData()
  MessagePlugin.success('数据已刷新')
}

// 监听 store 数据变化
watch(
  () => tableStore.currentTable.records,
  () => {
    initData()
  },
  { deep: true }
)
</script>

<style scoped lang="less">
.kanban-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #ffffff;
  position: relative;

  .kanban-toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 20px;
    border-bottom: 1px solid #e8e8e8;
    background: #ffffff;
    flex-shrink: 0;
    z-index: 10;

    .toolbar-left {
      display: flex;
      align-items: center;
      gap: 12px;
    }

    .toolbar-right {
      .total-count {
        font-size: 13px;
        color: #8f959e;
      }
    }
  }

  .kanban-container {
    flex: 1;
    overflow: hidden;
    padding: 20px;
    background: #f5f7fa;

    .kanban-scroll {
      display: flex;
      gap: 16px;
      height: 100%;
      overflow-x: auto;
      overflow-y: hidden;
      padding-bottom: 8px;

      &::-webkit-scrollbar {
        height: 8px;
      }

      &::-webkit-scrollbar-track {
        background: #f5f5f5;
        border-radius: 4px;
      }

      &::-webkit-scrollbar-thumb {
        background: #d0d0d0;
        border-radius: 4px;

        &:hover {
          background: #a0a0a0;
        }
      }
    }
  }
}
</style>
