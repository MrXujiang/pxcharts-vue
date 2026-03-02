<template>
  <div class="bar-chart-container">
    <div v-if="config.title" class="chart-title">{{ config.title }}</div>
    <div ref="chartRef" class="chart"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import * as echarts from 'echarts'

interface Props {
  config: {
    title?: string
    xAxisData?: string[]
    seriesName?: string
  }
  data: number[]
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

const initChart = () => {
  if (!chartRef.value) return

  chartInstance = echarts.init(chartRef.value)
  
  // 使用主题颜色，如果没有则使用默认蓝色
  const primaryColor = props.themeColors?.[0] || '#3b82f6'
  
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow',
      },
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '10%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: props.config.xAxisData || ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
      axisTick: {
        alignWithLabel: true,
      },
    },
    yAxis: {
      type: 'value',
    },
    series: [
      {
        name: props.config.seriesName || '数据',
        type: 'bar',
        barWidth: '60%',
        data: props.data.length > 0 ? props.data : [10, 52, 200, 334, 390, 330, 220],
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: primaryColor },
            { offset: 0.5, color: primaryColor },
            { offset: 1, color: primaryColor },
          ]),
        },
        emphasis: {
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: primaryColor },
              { offset: 0.7, color: primaryColor },
              { offset: 1, color: primaryColor },
            ]),
          },
        },
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
.bar-chart-container {
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
