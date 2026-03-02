<template>
  <PieChart 
    :config="{ ...config, isDoughnut: true }" 
    :data="data" 
    :chart-id="chartId"
    :theme-colors="themeColors"
  />
</template>

<script setup lang="ts">
import PieChart from './PieChart.vue'

interface Props {
  config: any
  data: Array<{ name: string; value: number }>
  chartId: string
  themeColors?: string[]
}

defineProps<Props>()
</script>
