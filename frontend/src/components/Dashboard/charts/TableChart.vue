<template>
  <div class="table-chart-container">
    <div v-if="config.title" class="chart-title">{{ config.title }}</div>
    <div class="table-wrapper">
      <t-table
        :data="tableData"
        :columns="tableColumns"
        :max-height="maxHeight"
        row-key="id"
        stripe
        hover
        size="small"
        :pagination="paginationConfig"
      >
        <!-- 空数据状态 -->
        <template #empty>
          <div class="empty-state">
            <t-icon name="inbox" size="48px" style="color: #ddd;" />
            <p>暂无数据</p>
          </div>
        </template>
      </t-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'

interface Column {
  colKey: string
  title: string
  width?: number
  align?: 'left' | 'center' | 'right'
  ellipsis?: boolean
}

interface Props {
  config: {
    title?: string
    showPagination?: boolean
    pageSize?: number
  }
  data: {
    columns?: Column[]
    rows?: any[]
  }
  chartId: string
}

const props = withDefaults(defineProps<Props>(), {
  config: () => ({}),
  data: () => ({ columns: [], rows: [] }),
})

// 默认列配置
const defaultColumns: Column[] = [
  { colKey: 'name', title: '姓名', width: 120, align: 'left' },
  { colKey: 'status', title: '状态', width: 100, align: 'center' },
  { colKey: 'value', title: '数值', width: 100, align: 'right' },
  { colKey: 'date', title: '日期', width: 120, align: 'center' },
]

// 默认数据
const defaultRows = [
  { id: 1, name: '张三', status: '进行中', value: 1250, date: '2024-01-15' },
  { id: 2, name: '李四', status: '已完成', value: 890, date: '2024-01-16' },
  { id: 3, name: '王五', status: '未开始', value: 0, date: '2024-01-17' },
  { id: 4, name: '赵六', status: '进行中', value: 2340, date: '2024-01-18' },
  { id: 5, name: '钱七', status: '已完成', value: 1560, date: '2024-01-19' },
]

// 表格列
const tableColumns = computed(() => {
  if (props.data.columns && props.data.columns.length > 0) {
    return props.data.columns
  }
  return defaultColumns
})

// 表格数据
const tableData = computed(() => {
  let rows = props.data.rows && props.data.rows.length > 0 ? props.data.rows : defaultRows
  // 确保每行都有id
  return rows.map((row, index) => ({
    ...row,
    id: row.id || index + 1,
  }))
})

// 最大高度（根据容器自适应）
const maxHeight = ref<number | string>('auto')

// 分页配置
const paginationConfig = computed(() => {
  if (props.config.showPagination === false) {
    return undefined
  }
  return {
    defaultCurrent: 1,
    defaultPageSize: props.config.pageSize || 10,
    total: tableData.value.length,
    showJumper: true,
    size: 'small' as const,
  }
})

// ResizeObserver监听容器大小
let resizeObserver: ResizeObserver | null = null
const containerRef = ref<HTMLElement>()

const updateMaxHeight = () => {
  if (!containerRef.value) return
  const containerHeight = containerRef.value.clientHeight
  const titleHeight = props.config.title ? 40 : 0
  const paginationHeight = props.config.showPagination !== false ? 48 : 0
  maxHeight.value = containerHeight - titleHeight - paginationHeight - 32
}

onMounted(() => {
  const container = document.querySelector('.table-chart-container') as HTMLElement
  if (container) {
    containerRef.value = container
    resizeObserver = new ResizeObserver(() => {
      updateMaxHeight()
    })
    resizeObserver.observe(container)
    updateMaxHeight()
  }
})

onUnmounted(() => {
  resizeObserver?.disconnect()
})

watch(() => [props.config.title, props.config.showPagination], () => {
  updateMaxHeight()
})
</script>

<style scoped lang="less">
.table-chart-container {
  width: 100%;
  height: 100%;
  padding: 16px;
  display: flex;
  flex-direction: column;

  .chart-title {
    font-size: 16px;
    font-weight: 500;
    color: #333;
    margin-bottom: 12px;
  }

  .table-wrapper {
    flex: 1;
    min-height: 0;
    overflow: hidden;

    :deep(.t-table) {
      font-size: 12px;

      .t-table__header {
        background: #f5f7fa;

        th {
          font-weight: 600;
          color: #333;
          background: #f5f7fa !important;
        }
      }

      .t-table__body {
        tr {
          &:hover {
            background: #f5f7fa;
          }
        }
      }
    }

    .empty-state {
      padding: 40px 0;
      text-align: center;
      color: #999;

      p {
        margin-top: 12px;
        font-size: 14px;
      }
    }
  }
}
</style>
