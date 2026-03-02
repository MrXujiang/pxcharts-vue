<template>
  <div>
    <t-select
      v-model="defaultModelValue"
      placeholder="请选择"
      :popup-props="{ overlayClassName: 'auto-number' }"
      :multiple="fieldObj.type === 'radio' ? false : true"
      @PopupVisibleChange="handlePopupVisibleChange"
      @keydown.esc="handleKeyDown"
    >
      <template #valueDisplay="{ value, onClose, label }">
        <template v-if="fieldObj.type != 'radio'">
          <div class="w-100 flex tag-display">
            <t-tag
              shape="round"
              size="small"
              v-for="(item, index) in value"
              :key="index"
              :closable="true"
              :color="fieldObj.options.find((option) => option.id === item.value)?.color"
              :on-close="
                ({ e }) => {
                  e.stopPropagation()
                  onClose(index)
                }
              "
            >
              {{ item.label }}
            </t-tag>
          </div>
        </template>
        <template v-else>
          <t-tag
            size="small"
            shape="round"
            closable
            :on-close="
              ({ e }) => {
                e.stopPropagation()
                fieldObj.defaultValue = ''
              }
            "
            :color="fieldObj.options.find((option) => option.id === value)?.color"
            >{{ label }}</t-tag
          >
        </template>
      </template>
      <t-option
        v-for="item in fieldObj.options"
        :key="item.id"
        :value="item.id"
        :label="item.label"
      >
        <t-tag size="small" :color="item.color" shape="round">{{ item.label }}</t-tag>
      </t-option>
    </t-select>
  </div>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'
interface Props {
  editingValue: string | number
  field: ObjType
}
const props = defineProps<Props>()
const emit = defineEmits(['update:editingValue', 'handleCellEdit', 'handleCancelCellEdit'])

const defaultModelValue = computed({
  get() {
    return props.editingValue
  },
  set(value) {
    emit('update:editingValue', value)
  },
})
const fieldObj = computed(() => props.field)

const handleCellEdit = () => {
  console.log('多选-单选 Cell edited', defaultModelValue.value)
  emit('handleCellEdit')
}

const handlePopupVisibleChange = (visible) => {
  if (!visible) {
    handleCellEdit()
  }
}

const handleKeyDown = () => {
  console.log('多选-单选 Key down')
  emit('handleCancelCellEdit')
}
defineOptions({ name: 'SingleMultipleChoiceComponent' })
</script>

<style lang="less" scoped>
.tag-display {
  width: 100%;
  padding: 0 25px 0 0;
  overflow: auto;
  white-space: nowrap;
  -ms-overflow-style: none;
  scrollbar-width: none;
  overflow-y: hidden;
  position: absolute;
}
</style>
