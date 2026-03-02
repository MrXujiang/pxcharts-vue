<template>
  <t-space :size="10" direction="vertical" class="w-100">
    <div class="flx-ce-bet">
      <div>人员字段选项</div>
    </div>
    <t-card size="small" class="w-100">
      <t-space :size="10" direction="vertical">
        <div class="flx-ce-sta gap-4">
          <t-checkbox v-model="formData.settings.showNickname">展示昵称</t-checkbox>
          <t-tooltip content="用户将显姓名-昵称">
            <t-icon class="gray-col" name="helpCircle" />
          </t-tooltip>
        </div>
        <div class="flx-ce-sta gap-4">
          <t-checkbox v-model="formData.settings.showExternalTag">展示组织外部标签</t-checkbox>
          <t-tooltip content="开启时将展示外部组织人员标签"> </t-tooltip>
          <t-icon class="gray-col" name="helpCircle" />
        </div>
        <div class="flx-ce-sta gap-4">
          <t-checkbox v-model="formData.settings.allowMultiple" @change="handleAllowMultipleChange"
            >允许选择多个人员</t-checkbox
          >
        </div>
      </t-space>
    </t-card>
    <!-- 使用提取的组件 -->
    <div class="flex-col gap-4">
      <div class="option-row">默认人员</div>
      <PersonPanel
        v-model="formData.defaultValue"
        :allow-multiple="formData.settings.allowMultiple"
      />
    </div>
  </t-space>
</template>

<script setup lang="ts">
import PersonPanel from '@/components/personPanel/index.vue'

const props = defineProps({
  formData: {
    type: Object,
    required: true,
  },
})
const { formData } = toRefs(props)

const handleAllowMultipleChange = () => {
  formData.value.defaultValue = formData.value.settings.allowMultiple ? [] : ''
}

defineOptions({ name: 'PersonComponent' })
</script>

<style lang="less" scoped>
.option-row {
  height: 32px !important;
  line-height: 32px !important;
}
.external {
  color: orange;
  font-size: 10px;
}
</style>
