<template>
  <component :is="iconComponent" v-bind="$attrs" />
</template>

<script setup lang="ts">
import { localIcons } from '@/assets/icons'
import { computed } from 'vue'

// 为组件设置多词名称，避免 ESLint 警告
defineOptions({
  name: 'TIcon',
})

const props = defineProps({
  name: {
    type: String,
    required: true,
  },
})

const iconComponent = computed(() => localIcons[props.name])
</script>
