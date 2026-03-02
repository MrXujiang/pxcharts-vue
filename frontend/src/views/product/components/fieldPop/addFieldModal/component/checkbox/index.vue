<template>
  <div class="w-100">
    <t-space class="flex-col" size="small">
      <div class="flx-ce-bet">
        <div>默认状态</div>
      </div>
      <!-- eslint-disable vue/valid-v-model -->
      <t-radio-group v-model="formData.defaultValue">
        <t-radio :value="0">默认不勾选</t-radio>
        <t-radio :value="1">默认勾选</t-radio>
      </t-radio-group>
    </t-space>
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  formData: {
    type: Object,
    required: true,
  },
})
const { formData } = toRefs(props)

defineOptions({ name: 'CheckboxComponent' })
</script>

<style scoped></style>
