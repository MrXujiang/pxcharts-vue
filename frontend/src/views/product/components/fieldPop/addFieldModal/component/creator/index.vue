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
            <t-icon class="gray-col pointer" name="helpCircle" />
          </t-tooltip>
        </div>
        <div class="flx-ce-sta gap-4">
          <t-checkbox v-model="formData.settings.showExternalTag">展示组织外部标签</t-checkbox>
          <t-tooltip content="开启时将展示外部组织人员标签"> </t-tooltip>
          <t-icon class="gray-col pointer" name="helpCircle" />
        </div>
      </t-space>
    </t-card>
  </t-space>
</template>

<script setup lang="ts">
const props = defineProps({
  formData: {
    type: Object,
    required: true,
  },
})
const { formData } = toRefs(props)
defineOptions({ name: 'CreatorComponent' })
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
