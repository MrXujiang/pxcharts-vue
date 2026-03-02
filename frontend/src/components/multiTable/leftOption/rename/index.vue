<template>
  <t-dialog
    :visible="visible"
    header="重命名"
    width="400px"
    @confirm="handleConfirm"
    @cancel="handleCancel"
    :confirm-btn="{ loading: confirmLoading }"
    @close="handleClose"
  >
    <div class="w-100 h-100" style="height: 80px">
      <t-form
        ref="formRef"
        :data="formData"
        :rules="rules"
        label-align="top"
        @submit="handleConfirm"
      >
        <t-form-item label="名称" name="name">
          <t-input v-model="formData.name" placeholder="请输入" clearable />
        </t-form-item>
      </t-form>
    </div>
  </t-dialog>
</template>

<script setup lang="ts">
import type { FormInstanceFunctions, FormRules } from 'tdesign-vue-next'

defineOptions({
  name: 'RenameDialog',
})

// 定义 props
const props = defineProps<{
  visible: boolean
  name?: string
}>()

// 定义 emits
const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'confirm', name: string): void
}>()

// 表单引用
const formRef = ref<FormInstanceFunctions>()

// 表单数据
const formData = reactive({
  name: '',
})

// 监听 props.name 的变化
watch(
  () => props.name,
  (newName) => {
    if (newName !== undefined) {
      formData.name = newName
    }
  },
  { immediate: true },
)

// 表单验证规则
const rules: FormRules = {
  name: [
    { required: true, message: '请输入名称', trigger: 'blur' },
    { min: 1, max: 50, message: '名称长度不能超过50个字符', trigger: 'blur' },
  ],
}

// 确认按钮加载状态
const confirmLoading = ref(false)

// 处理确认事件
const handleConfirm = async () => {
  if (!formRef.value) return

  const result = await formRef.value.validate()
  if (result === true) {
    confirmLoading.value = true
    emit('confirm', formData.name)
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
