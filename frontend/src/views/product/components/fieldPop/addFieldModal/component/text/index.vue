<template>
  <div class="w-100">
    <div class="flx-ce-bet nav-bar">
      <div>默认文本</div>
    </div>
    <t-input v-model="formData.defaultValue" placeholder="请输入内容" />
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

defineOptions({
  name: 'TextInput',
})
</script>

<style lang="less" scoped>
.nav-bar {
  height: 32px;
  line-height: 32px;
}
</style>
