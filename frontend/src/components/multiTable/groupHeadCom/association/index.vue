<template>
  <div>xx</div>
</template>

<script setup lang="ts">
interface Props {
  value: string
}

const props = defineProps<Props>()
console.log('propsasdasdasdasdas', props)
defineOptions({ name: 'GroupHeadAssociation' })
</script>

<style scoped></style>
