<template>
  <div class="link-container flx-ce-sta">
    <t-popup ref="popupRef" placement="bottom-left" @visible-change="handleVisibleChange">
      <div class="link">
        <a :href="displayData.linkUrl" target="_blank">{{
          displayData.linkTitle ? displayData.linkTitle : displayData.linkUrl
        }}</a>
      </div>
      <template #content>
        <t-card size="small" :bordered="false" style="width: 360px">
          <t-space class="w-100" size="small" direction="vertical">
            <div class="flex-col gap-8">
              <div>文本</div>
              <t-input v-model="formData.linkTitle" placeholder="请输入链接名称" />
              <div>链接</div>
              <t-input v-model="formData.linkUrl" placeholder="请输入链接地址" />
            </div>
            <!-- 按钮区域 -->
            <div class="w-100 flx-ce-end gap-8">
              <t-button variant="outline" @click="cancelEdit">取消</t-button>
              <t-button @click="saveEdit">保存</t-button>
            </div>
          </t-space>
        </t-card>
      </template>
    </t-popup>
  </div>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'

interface Props {
  editingValue: ObjType
}

const props = defineProps<Props>()
const emit = defineEmits(['handleCellEdit', 'handleCancelCellEdit', 'update:editingValue'])

// 添加 popup 引用
const popupRef = ref()

// 使用单独的数据存储显示内容，避免直接绑定到 props
const displayData = computed(() => ({
  linkUrl: props.editingValue.linkUrl || '',
  linkTitle: props.editingValue.linkTitle || '',
}))

// 表单数据用于编辑
const formData = ref({
  linkUrl: props.editingValue.linkUrl || '',
  linkTitle: props.editingValue.linkTitle || '',
})

// 取消编辑
const cancelEdit = () => {
  emit('handleCancelCellEdit')
  // 重置表单数据为显示数据
  formData.value = {
    ...displayData.value,
  }
  // 关闭 popup
  if (popupRef.value) {
    popupRef.value.close()
  }
}

// 保存编辑
const saveEdit = () => {
  // 发送更新事件
  emit('update:editingValue', {
    ...formData.value,
  })
  emit('handleCellEdit')

  // 关闭 popup
  if (popupRef.value) {
    popupRef.value.close()
  }
}
const handleVisibleChange = (visible: boolean) => {
  if (!visible) {
    emit('handleCancelCellEdit')
    // 重置表单数据为显示数据
    formData.value = {
      ...displayData.value,
    }
  }
}
defineOptions({ name: 'LinkComponent' })
</script>

<style lang="less" scoped>
.text {
  font-size: 14px;
}
.link-container {
  border: 1px solid #dfe1e6;
  border-radius: 4px;
  height: 32px;
  padding: 0 12px;
  align-items: center;
}
.link {
  width: 100%;
  min-height: 20px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
  color: #2c3e50;
}
</style>
