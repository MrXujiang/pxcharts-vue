<template>
  <t-tooltip content="系统自动填充，不可修改">
    <div class="time-text">{{ defaultModelValue }}</div>
  </t-tooltip>
</template>

<script setup lang="ts">
interface Props {
  editingValue: string | number
}

const props = defineProps<Props>()
const defaultModelValue = computed(() => props.editingValue)

defineOptions({ name: 'FindRefComponent' })
</script>

<style lang="less" scoped>
.time-text {
  font-size: 14px;
}
</style>
