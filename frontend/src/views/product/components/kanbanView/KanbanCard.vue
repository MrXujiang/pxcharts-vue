<template>
  <div class="kanban-card" :data-id="card.rowId">
    <div class="card-header">
      <div class="card-title">{{ getCardTitle(card) }}</div>
      <div class="card-actions">
        <t-icon name="more" class="action-icon" />
      </div>
    </div>
    <div class="card-body">
      <div v-for="field in displayFields" :key="field.id" class="card-field">
        <span class="field-label">{{ field.title }}:</span>
        <span class="field-value">{{ formatFieldValue(card[field.id], field) }}</span>
      </div>
    </div>
    <div v-if="showFooter" class="card-footer">
      <div class="card-meta">
        <span class="card-id">#{{ card.idText || card.rowId }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  card: Record<string, any>
  fields: Array<{
    id: string
    title: string
    type?: string
    isShow?: boolean
  }>
  titleField?: string
  showFooter?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  titleField: 'name',
  showFooter: true,
})

// 显示的字段（排除标题字段和一些系统字段）
const displayFields = computed(() => {
  return props.fields.filter(
    (field) =>
      field.id !== props.titleField &&
      field.id !== 'idText' &&
      field.id !== 'rowId' &&
      field.isShow !== false
  ).slice(0, 4) // 最多显示4个字段
})

// 获取卡片标题
const getCardTitle = (card: Record<string, any>) => {
  return card[props.titleField] || '未命名'
}

// 格式化字段值
const formatFieldValue = (value: any, field: { type?: string }) => {
  if (value === null || value === undefined || value === '') {
    return '-'
  }
  
  switch (field.type) {
    case 'number':
      return typeof value === 'number' ? value.toLocaleString() : value
    case 'date':
      return value ? new Date(value).toLocaleDateString('zh-CN') : '-'
    default:
      return value
  }
}
</script>

<style scoped lang="less">
.kanban-card {
  background: #ffffff;
  border-radius: 8px;
  border: 1px solid #e8e8e8;
  padding: 12px;
  margin-bottom: 12px;
  cursor: grab;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  user-select: none;

  &:hover {
    border-color: #4a7ff7;
    box-shadow: 0 2px 8px rgba(74, 127, 247, 0.15);
    transform: translateY(-1px);
  }
  
  &:active {
    cursor: grabbing;
  }

  &.sortable-ghost {
    opacity: 0.4;
    background: #f0f5ff;
  }

  &.sortable-drag {
    opacity: 0.8;
    transform: rotate(2deg);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 12px;

    .card-title {
      font-size: 14px;
      font-weight: 500;
      color: #1f2329;
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .card-actions {
      display: flex;
      gap: 4px;
      opacity: 0;
      transition: opacity 0.2s;

      .action-icon {
        font-size: 16px;
        color: #8f959e;
        cursor: pointer;
        padding: 2px;
        border-radius: 4px;

        &:hover {
          background: #f0f0f0;
          color: #4a7ff7;
        }
      }
    }
  }

  &:hover .card-actions {
    opacity: 1;
  }

  .card-body {
    .card-field {
      display: flex;
      align-items: baseline;
      margin-bottom: 8px;
      font-size: 13px;

      &:last-child {
        margin-bottom: 0;
      }

      .field-label {
        color: #8f959e;
        margin-right: 8px;
        flex-shrink: 0;
        min-width: 60px;
      }

      .field-value {
        color: #1f2329;
        flex: 1;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }

  .card-footer {
    margin-top: 12px;
    padding-top: 12px;
    border-top: 1px solid #f0f0f0;

    .card-meta {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .card-id {
        font-size: 12px;
        color: #8f959e;
      }
    }
  }
}
</style>
