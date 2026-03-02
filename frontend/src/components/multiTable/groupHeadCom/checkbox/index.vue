<template>
  <div>
    <t-checkbox disabled v-model="checkboxValue"> </t-checkbox>
  </div>
</template>

<script setup lang="ts">
interface Props {
  value: string
}

const props = defineProps<Props>()

const checkboxValue = computed(() => (props.value == '未分组' ? 0 : 1))
defineOptions({ name: 'GroupHeadCheckbox' })
</script>

<style scoped></style>
