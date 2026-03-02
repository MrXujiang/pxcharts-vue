<template>
  <div class="pie-chart-container">
    <div v-if="config.title" class="chart-title">{{ config.title }}</div>
    <div ref="chartRef" class="chart" :style="{ height: chartHeight }"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import * as echarts from 'echarts'

interface Props {
  config: {
    title?: string
    isDoughnut?: boolean
    showLegend?: boolean
  }
  data: Array<{ name: string; value: number }>
  chartId: string
  themeColors?: string[]
}

const props = withDefaults(defineProps<Props>(), {
  config: () => ({}),
  data: () => [],
  themeColors: () => ['#3b82f6', '#8b5cf6', '#22c55e', '#eab308', '#f97316', '#ec4899'],
})

const chartRef = ref<HTMLElement>()
let chartInstance: echarts.ECharts | null = null
let resizeObserver: ResizeObserver | null = null

const chartHeight = ref('100%')

const initChart = () => {
  if (!chartRef.value) return

  chartInstance = echarts.init(chartRef.value)
  
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)',
    },
    legend: props.config.showLegend !== false ? {
      orient: 'vertical',
      right: 10,
      top: 'center',
      textStyle: {
        fontSize: 12,
      },
    } : undefined,
    series: [
      {
        name: props.config.title || '数据分布',
        type: 'pie',
        radius: props.config.isDoughnut ? ['40%', '70%'] : '60%',
        center: props.config.showLegend !== false ? ['40%', '50%'] : ['50%', '50%'],
        data: props.data.length > 0 ? props.data : [
          { value: 6, name: '进行中' },
          { value: 4, name: '空值' },
          { value: 6, name: '未开始' },
          { value: 10, name: '已完成' },
        ],
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)',
          },
        },
        label: {
          formatter: '{b}\n{d}%',
          fontSize: 12,
        },
        color: props.themeColors,
      },
    ],
  }

  chartInstance.setOption(option)
}

const resizeChart = () => {
  chartInstance?.resize()
}

watch(() => props.data, () => {
  if (chartInstance) {
    initChart()
  }
}, { deep: true })

watch(() => props.themeColors, () => {
  if (chartInstance) {
    initChart()
  }
}, { deep: true })

onMounted(() => {
  nextTick(() => {
    initChart()
    
    // 使用ResizeObserver监听容器大小变化
    if (chartRef.value) {
      resizeObserver = new ResizeObserver(() => {
        resizeChart()
      })
      resizeObserver.observe(chartRef.value)
    }
  })
  
  window.addEventListener('resize', resizeChart)
})

onUnmounted(() => {
  window.removeEventListener('resize', resizeChart)
  resizeObserver?.disconnect()
  chartInstance?.dispose()
})
</script>

<style scoped lang="less">
.pie-chart-container {
  width: 100%;
  height: 100%;
  padding: 16px;
  display: flex;
  flex-direction: column;

  .chart-title {
    font-size: 16px;
    font-weight: 500;
    color: #333;
    margin-bottom: 12px;
  }

  .chart {
    flex: 1;
    min-height: 200px;
  }
}
</style>
