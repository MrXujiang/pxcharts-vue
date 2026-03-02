<!-- 编辑表描述 -->
<template>
  <t-dialog
    :visible="visible"
    header="数据表描述"
    width="400px"
    @confirm="handleConfirm"
    @cancel="handleCancel"
    @close="handleClose"
  >
    <div class="w-100 h-100" style="min-height: 100px">
      <t-form
        ref="formRef"
        :data="formData"
        :rules="rules"
        label-align="top"
        @submit="handleConfirm"
      >
        <t-form-item label="描述" name="desc">
          <t-textarea v-model="formData.desc" placeholder="请输入" clearable />
        </t-form-item>
      </t-form>
    </div>
  </t-dialog>
</template>

<script setup lang="ts">
import type { FormInstanceFunctions, FormRules } from 'tdesign-vue-next'

defineOptions({
  name: 'EditTbDescDialog',
})

// 定义 props
const props = defineProps<{
  visible: boolean
  desc?: string
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
  desc: '',
})

// 监听 props.name 的变化
watch(
  () => props.desc,
  (newName) => {
    if (newName !== undefined) {
      formData.desc = newName
    }
  },
  { immediate: true },
)

// 表单验证规则
const rules: FormRules = {
  desc: [{ required: true, message: '请输入', trigger: 'blur' }],
}

// 确认按钮加载状态
const confirmLoading = ref(false)

// 处理确认事件
const handleConfirm = async () => {
  if (!formRef.value) return

  const result = await formRef.value.validate()
  if (result === true) {
    confirmLoading.value = true
    emit('confirm', formData.desc)
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
