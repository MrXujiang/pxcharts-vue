<template>
  <div class="kanban-column">
    <div class="column-header">
      <div class="column-title">
        <span class="title-text">{{ group.label }}</span>
        <t-tag size="small" theme="default" variant="outline" class="count-tag">
          {{ cards.length }}
        </t-tag>
      </div>
      <div class="column-actions">
        <t-icon name="add" class="action-icon" @click="handleAddCard" />
      </div>
    </div>
    <div
      ref="cardListRef"
      class="column-body"
      :data-group-value="group.value"
    >
      <KanbanCard
        v-for="card in cards"
        :key="card.rowId"
        :card="card"
        :fields="fields"
        :title-field="titleField"
      />
      <div v-if="cards.length === 0" class="empty-placeholder">
        <t-icon name="inbox" size="32px" class="empty-icon" />
        <p class="empty-text">暂无数据</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import Sortable from 'sortablejs'
import KanbanCard from './KanbanCard.vue'

interface Props {
  group: {
    label: string
    value: any
  }
  cards: Array<Record<string, any>>
  fields: Array<{
    id: string
    title: string
    type?: string
    isShow?: boolean
  }>
  titleField?: string
}

interface Emits {
  (e: 'card-move', event: {
    fromGroup: any
    toGroup: any
    oldIndex: number
    newIndex: number
    card: Record<string, any>
  }): void
  (e: 'add-card', group: any): void
}

const props = withDefaults(defineProps<Props>(), {
  titleField: 'name',
})

const emit = defineEmits<Emits>()

const cardListRef = ref<HTMLElement | null>(null)
let sortableInstance: Sortable | null = null

// 处理添加卡片
const handleAddCard = () => {
  emit('add-card', props.group.value)
}

// 初始化拖拽
onMounted(() => {
  if (!cardListRef.value) return

  sortableInstance = Sortable.create(cardListRef.value, {
    group: 'kanban-cards',
    animation: 200, // 增加动画时长，更流畅
    ghostClass: 'sortable-ghost',
    chosenClass: 'sortable-chosen', // 添加选中样式
    dragClass: 'sortable-drag',
    handle: '.kanban-card',
    forceFallback: false, // 使用原生HTML5拖拽
    fallbackOnBody: true,
    swapThreshold: 0.65,
    invertSwap: false,
    
    // 拖拽开始
    onStart: (evt) => {
      document.body.classList.add('is-dragging')
    },
    
    // 拖拽结束
    onEnd: (evt) => {
      document.body.classList.remove('is-dragging')
      
      const { oldIndex, newIndex, from, to } = evt

      if (oldIndex === undefined || newIndex === undefined) return

      const fromGroup = from.getAttribute('data-group-value')
      const toGroup = to.getAttribute('data-group-value')

      if (fromGroup === null || toGroup === null) return

      // 获取被移动的卡片
      const movedCard = props.cards[oldIndex]

      emit('card-move', {
        fromGroup,
        toGroup,
        oldIndex,
        newIndex,
        card: movedCard,
      })
    },
  })
})

onBeforeUnmount(() => {
  sortableInstance?.destroy()
})
</script>

<style scoped lang="less">
.kanban-column {
  background: #f5f7fa;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  min-width: 300px;
  max-width: 320px;
  flex-shrink: 0;
  height: 100%;
  position: relative;

  .column-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px;
    border-bottom: 1px solid #e8e8e8;
    background: #ffffff;
    border-radius: 8px 8px 0 0;
    flex-shrink: 0;
    position: relative;
    z-index: 1;

    .column-title {
      display: flex;
      align-items: center;
      gap: 8px;

      .title-text {
        font-size: 14px;
        font-weight: 500;
        color: #1f2329;
      }

      .count-tag {
        font-size: 12px;
      }
    }

    .column-actions {
      display: flex;
      gap: 4px;

      .action-icon {
        font-size: 16px;
        color: #8f959e;
        cursor: pointer;
        padding: 4px;
        border-radius: 4px;
        transition: all 0.2s;

        &:hover {
          background: #f0f0f0;
          color: #4a7ff7;
        }
      }
    }
  }

  .column-body {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
    padding: 12px;
    min-height: 200px;
    background: #f5f7fa;

    &::-webkit-scrollbar {
      width: 6px;
    }

    &::-webkit-scrollbar-thumb {
      background: #d0d0d0;
      border-radius: 3px;

      &:hover {
        background: #a0a0a0;
      }
    }

    .empty-placeholder {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 40px 20px;
      text-align: center;

      .empty-icon {
        color: #d0d0d0;
        margin-bottom: 12px;
      }

      .empty-text {
        font-size: 13px;
        color: #8f959e;
        margin: 0;
      }
    }
  }
}

// 拖拽样式
:deep(.sortable-ghost) {
  opacity: 0.4;
  background: #e3f2fd;
  border: 2px dashed #4a7ff7;
}

:deep(.sortable-chosen) {
  cursor: move;
  transform: rotate(2deg);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

:deep(.sortable-drag) {
  opacity: 1;
  cursor: move;
  transform: rotate(3deg);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
}
</style>
