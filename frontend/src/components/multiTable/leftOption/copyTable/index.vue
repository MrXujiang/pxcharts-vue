<template>
  <t-dialog
    :visible="visible"
    header="复制数据表"
    width="400px"
    @confirm="handleConfirm"
    @cancel="handleCancel"
    :confirm-btn="{ loading: confirmLoading }"
    @close="handleClose"
  >
    <div class="w-100 h-100" style="height: 160px">
      <t-form
        ref="formRef"
        :data="formData"
        :rules="rules"
        label-align="top"
        @submit="handleConfirm"
      >
        <t-form-item label="数据表名称" name="dataTableName">
          <t-input v-model="formData.dataTableName" placeholder="请输入" clearable />
        </t-form-item>
        <t-form-item label="复制范围" name="copyScope">
          <t-radio-group v-model="formData.copyScope">
            <t-radio :value="1">仅数据表结构</t-radio>
            <t-radio :value="2">数据表结构跟所有记录</t-radio>
          </t-radio-group>
        </t-form-item>
      </t-form>
    </div>
  </t-dialog>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'
import type { FormInstanceFunctions, FormRules } from 'tdesign-vue-next'

defineOptions({
  name: 'CopyTable',
})

// 定义 props
const props = defineProps<{
  visible: boolean
  objData: ObjType
}>()

// 定义 emits
const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'confirm', data: { dataTableName: string; copyScope: string }): void
}>()

// 表单引用
const formRef = ref<FormInstanceFunctions>()

// 表单数据
const formData = reactive({
  dataTableName: '',
  copyScope: 1, // 默认选择仅数据表结构
})

// 监听 props.name 的变化
watch(
  () => props.objData,
  (newName) => {
    if (newName !== undefined) {
      formData.dataTableName = newName.name
    }
  },
  { immediate: true },
)

// 表单验证规则
const rules: FormRules = {
  dataTableName: [
    { required: true, message: '请输入名称', trigger: 'blur' },
    { min: 1, max: 50, message: '名称长度不能超过50个字符', trigger: 'blur' },
  ],
  copyScope: [{ required: true, message: '请输入', trigger: 'blur' }],
}

// 确认按钮加载状态
const confirmLoading = ref(false)

// 处理确认事件
const handleConfirm = async () => {
  if (!formRef.value) return

  const result = await formRef.value.validate()
  if (result === true) {
    confirmLoading.value = true
    emit('confirm', { ...formData })
    confirmLoading.value = false
    emit('update:visible', false)
  }
}

// 处理取消事件
const handleCancel = () => {
  emit('update:visible', false)
}

// 处理关闭事件
const handleClose = () => {
  emit('update:visible', false)
}
</script>
