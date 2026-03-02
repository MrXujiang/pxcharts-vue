<template>
  <div class="dashboard-designer">
    <!-- 顶部工具栏 -->
    <div class="dashboard-header">
      <div class="header-left">
        <t-input v-model="dashboardConfig.title" placeholder="仪表盘标题" class="dashboard-title-input" />
      </div>
      <div class="header-right">
        <!-- 预览模式切换 -->
        <t-radio-group v-model="previewMode" variant="default-filled" size="small" class="preview-switch">
          <t-radio-button value="pc">
            <t-icon name="desktop" />
            PC
          </t-radio-button>
          <t-radio-button value="mobile">
            <t-icon name="mobile" />
            移动
          </t-radio-button>
        </t-radio-group>
        
        <t-divider layout="vertical" />
        
        <t-button variant="text" @click="showChartSelector = true">
          <template #icon>
            <t-icon name="add" />
          </template>
          添加图表
        </t-button>
        <t-button variant="text" @click="toggleEditMode">
          <template #icon>
            <t-icon :name="isEditMode ? 'check' : 'edit'" />
          </template>
          {{ isEditMode ? '完成编辑' : '编辑布局' }}
        </t-button>
        <t-button variant="text" @click="showThemeSelector = true">
          <template #icon>
            <t-icon name="palette" />
          </template>
          主题
        </t-button>
        <t-button variant="text" @click="handleFullscreen">
          <template #icon>
            <t-icon name="fullscreen" />
          </template>
          全屏
        </t-button>
        <t-button theme="primary" @click="handleSave">
          <template #icon>
            <t-icon name="save" />
          </template>
          保存
        </t-button>
      </div>
    </div>

    <!-- 仪表盘画布 -->
    <div 
      class="dashboard-canvas" 
      :class="{ 'edit-mode': isEditMode, [`preview-${previewMode}`]: true, [`theme-${dashboardConfig.theme}`]: true }"
      :style="canvasBackgroundStyle"
    >
      <div 
        class="canvas-container"
        :style="previewMode === 'mobile' ? { height: mobileContainerHeight } : {}"
      >
        <GridLayout
          v-model:layout="currentLayout"
          :col-num="previewMode === 'mobile' ? 1 : 12"
          :row-height="previewMode === 'mobile' ? 120 : 60"
          :is-draggable="isEditMode && previewMode === 'pc'"
          :is-resizable="isEditMode && previewMode === 'pc'"
          :vertical-compact="true"
          :use-css-transforms="true"
          :margin="previewMode === 'mobile' ? [12, 12] : [16, 16]"
          :responsive="false"
        >
          <GridItem
            v-for="item in currentLayout"
            :key="item.i"
            :x="item.x"
            :y="item.y"
            :w="item.w"
            :h="item.h"
            :i="item.i"
            :static="previewMode === 'mobile'"
            :is-draggable="isEditMode && previewMode === 'pc'"
            :is-resizable="isEditMode && previewMode === 'pc'"
          >
            <div class="chart-wrapper" :class="{ 'editing': isEditMode && previewMode === 'pc' }">
              <!-- 图表操作按钮 -->
              <div v-if="isEditMode && previewMode === 'pc'" class="chart-actions">
                <t-button size="small" variant="text" @click="handleEditChartData(item)">
                  <t-icon name="chart-bar" />
                </t-button>
                <t-button size="small" variant="text" @click="handleConfigChart(item)">
                  <t-icon name="edit" />
                </t-button>
                <t-button size="small" variant="text" @click="handleDeleteChart(item)">
                  <t-icon name="delete" />
                </t-button>
              </div>

              <!-- 图表内容 -->
              <component
                :is="getChartComponent(item.type)"
                :config="item.config"
                :data="item.data"
                :chart-id="item.i"
                :theme-colors="currentChartTheme.chartColors"
              />
            </div>
          </GridItem>
        </GridLayout>
      </div>

      <!-- 空状态 -->
      <div v-if="layout.length === 0" class="empty-dashboard">
        <div class="empty-content">
          <t-icon name="chart-bar" size="64px" />
          <p class="empty-title">开始构建你的仪表盘</p>
          <p class="empty-desc">点击"添加图表"开始添加可视化组件</p>
          <t-button theme="primary" @click="showChartSelector = true">
            <template #icon><t-icon name="add" /></template>
            添加图表
          </t-button>
        </div>
      </div>
    </div>

    <!-- 图表选择器弹窗 -->
    <t-dialog
      v-model:visible="showChartSelector"
      header="选择图表类型"
      width="800px"
      :footer="false"
    >
      <div class="chart-selector">
        <div
          v-for="chart in chartTypes"
          :key="chart.type"
          class="chart-type-card"
          @click="handleAddChart(chart.type)"
        >
          <div class="chart-icon">
            <component :is="chart.icon" />
          </div>
          <div class="chart-info">
            <div class="chart-name">{{ chart.name }}</div>
            <div class="chart-desc">{{ chart.description }}</div>
          </div>
        </div>
      </div>
    </t-dialog>

    <!-- 图表配置弹窗 -->
    <t-dialog
      v-model:visible="showChartConfig"
      header="图表配置"
      width="600px"
      @confirm="handleSaveChartConfig"
    >
      <div v-if="currentEditChart" class="chart-config-form">
        <t-form :data="currentEditChart.config" label-align="top">
          <t-form-item label="图表标题">
            <t-input v-model="currentEditChart.config.title" placeholder="请输入图表标题" />
          </t-form-item>
          <t-form-item label="数据源">
            <t-select v-model="currentEditChart.config.dataSource" placeholder="选择数据源">
              <t-option value="table1" label="当前表格数据" />
              <t-option value="custom" label="自定义数据" />
            </t-select>
          </t-form-item>
          <!-- 根据不同图表类型显示不同配置项 -->
        </t-form>
      </div>
    </t-dialog>
    
    <!-- 删除确认对话框 -->
    <t-dialog
      v-model:visible="showDeleteConfirm"
      header="删除确认"
      width="480px"
      @confirm="confirmDelete"
    >
      <div class="delete-confirm-content">
        <p>确定要删除图表「<strong>{{ pendingDeleteChart?.config?.title || '未命名' }}</strong>」吗？</p>
        <p class="warning-text">删除后将无法恢复，请谨慎操作。</p>
      </div>
    </t-dialog>
    
    <!-- 数据编辑器弹窗 -->
    <t-dialog
      v-model:visible="showDataEditor"
      header="编辑图表数据"
      width="700px"
      @confirm="handleSaveChartData"
    >
      <div v-if="currentEditChart" class="data-editor-form">
        <t-alert theme="info" style="margin-bottom: 16px;">
          正在编辑图表：{{ currentEditChart.config.title || '未命名' }}
        </t-alert>
        
        <!-- 数字卡片数据编辑 -->
        <template v-if="currentEditChart.type === 'number-card'">
          <t-form label-align="top">
            <t-form-item label="数值">
              <t-input-number v-model="currentEditChart.data.value" placeholder="请输入数值" />
            </t-form-item>
            <t-form-item label="趋势">
              <t-radio-group v-model="currentEditChart.data.trend">
                <t-radio value="up">上升</t-radio>
                <t-radio value="down">下降</t-radio>
              </t-radio-group>
            </t-form-item>
            <t-form-item label="对比值">
              <t-input v-model="currentEditChart.data.compare" placeholder="例：12.5%" />
            </t-form-item>
          </t-form>
        </template>
        
        <!-- 饼图/环形图数据编辑 -->
        <template v-else-if="currentEditChart.type === 'pie' || currentEditChart.type === 'doughnut'">
          <div class="data-list">
            <div v-for="(dataItem, index) in currentEditChart.data" :key="index" class="data-item">
              <t-input v-model="dataItem.name" placeholder="名称" style="flex: 1;" />
              <t-input-number v-model="dataItem.value" placeholder="数值" style="flex: 1; margin-left: 8px;" />
              <t-button variant="text" @click="removeDataItem(index)" style="margin-left: 8px;">
                <t-icon name="delete" />
              </t-button>
            </div>
            <t-button theme="default" @click="addDataItem" style="margin-top: 8px; width: 100%;">
              <t-icon name="add" />
              添加数据
            </t-button>
          </div>
        </template>
        
        <!-- 柱状图数据编辑 -->
        <template v-else-if="currentEditChart.type === 'bar' || currentEditChart.type === 'line'">
          <t-form label-align="top">
            <t-form-item label="数据值（逗号分隔）">
              <t-input 
                :value="Array.isArray(currentEditChart.data) ? currentEditChart.data.join(',') : ''"
                @change="handleBarDataChange"
                placeholder="例：10,52,200,334,390" 
              />
            </t-form-item>
          </t-form>
        </template>
      </div>
    </t-dialog>
    
    <!-- 主题选择器弹窗 -->
    <t-dialog
      v-model:visible="showThemeSelector"
      header="主题设置"
      width="1000px"
      :footer="false"
    >
      <div class="theme-selector">
        <!-- 图表主题 -->
        <div class="theme-section">
          <h4 class="section-title">图表主题</h4>
          <div class="theme-grid chart-themes">
            <div
              v-for="theme in chartThemes"
              :key="theme.id"
              class="theme-card"
              :class="{ active: dashboardConfig.theme === theme.id }"
              @click="handleSelectChartTheme(theme.id)"
            >
              <div class="theme-preview chart-preview">
                <div class="preview-colors">
                  <span 
                    v-for="(color, i) in theme.chartColors.slice(0, 6)"
                    :key="i"
                    class="color-bar"
                    :style="{ backgroundColor: color }"
                  />
                </div>
              </div>
              <div class="theme-name">
                {{ theme.name }}
                <t-icon v-if="dashboardConfig.theme === theme.id" name="check-circle-filled" class="check-icon" />
              </div>
            </div>
          </div>
        </div>
        
        <!-- 背景主题 -->
        <div class="theme-section">
          <h4 class="section-title">背景主题</h4>
          <div class="theme-grid">
            <div
              v-for="bg in backgroundThemes"
              :key="bg.id"
              class="theme-card"
              :class="{ active: dashboardConfig.background === bg.id && !dashboardConfig.customBackground }"
              @click="handleSelectBackground(bg.id)"
            >
              <div 
                class="theme-preview"
                :style="{
                  background: bg.backgroundImage || bg.backgroundColor,
                }"
              >
                <div 
                  class="preview-card"
                  :style="{
                    background: bg.cardBackground,
                    border: `1px solid ${bg.cardBorder}`,
                    color: bg.textColor,
                  }"
                >
                  <div class="preview-title" :style="{ color: bg.textColor }">示例卡片</div>
                  <div class="preview-text" :style="{ color: bg.textSecondary }">文本内容</div>
                </div>
              </div>
              <div class="theme-name">
                {{ bg.name }}
                <t-icon v-if="dashboardConfig.background === bg.id && !dashboardConfig.customBackground" name="check-circle-filled" class="check-icon" />
              </div>
            </div>
          </div>
        </div>
        
        <!-- 自定义背景 -->
        <div class="theme-section">
          <h4 class="section-title">自定义背景</h4>
          <div class="custom-background">
            <t-radio-group v-model="customBgType" variant="default-filled" style="margin-bottom: 16px;">
              <t-radio-button value="color">纯色背景</t-radio-button>
              <t-radio-button value="image">图片背景</t-radio-button>
            </t-radio-group>
            
            <div class="bg-input-wrapper">
              <template v-if="customBgType === 'color'">
                <t-input 
                  v-model="customBgValue" 
                  placeholder="输入颜色值，如 #f5f7fa"
                  style="flex: 1;"
                />
                <input 
                  type="color" 
                  v-model="customBgValue"
                  class="color-picker"
                />
              </template>
              <template v-else>
                <input
                  ref="bgFileInputRef"
                  type="file"
                  accept="image/*"
                  style="display: none;"
                  @change="handleBgImageUpload"
                />
                <t-button theme="default" @click="triggerImageUpload" style="margin-right: 8px;">
                  <template #icon><t-icon name="upload" /></template>
                  上传图片
                </t-button>
                <t-input 
                  v-model="customBgValue" 
                  placeholder="或直接输入图片URL"
                  style="flex: 1;"
                />
              </template>
              <t-button theme="primary" @click="applyCustomBackground">应用</t-button>
            </div>
            
            <div v-if="dashboardConfig.customBackground" class="custom-bg-preview">
              <div 
                class="preview-box"
                :style="getCustomBackgroundStyle()"
              >
                <span>当前自定义背景</span>
              </div>
              <t-button variant="text" theme="danger" @click="clearCustomBackground">
                清除自定义背景
              </t-button>
            </div>
          </div>
        </div>
      </div>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, defineAsyncComponent } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { GridLayout, GridItem } from 'vue3-grid-layout-next'
import { uploadFile } from '@/api'

// 异步加载图表组件
const NumberCard = defineAsyncComponent(() => import('./charts/NumberCard.vue'))
const PieChart = defineAsyncComponent(() => import('./charts/PieChart.vue'))
const BarChart = defineAsyncComponent(() => import('./charts/BarChart.vue'))
const DoughnutChart = defineAsyncComponent(() => import('./charts/DoughnutChart.vue'))
const TableChart = defineAsyncComponent(() => import('./charts/TableChart.vue'))

interface ChartItem {
  i: string
  x: number
  y: number
  w: number
  h: number
  type: string
  config: any
  data: any
}

interface ThemeConfig {
  id: string
  name: string
  chartColors: string[]
}

interface BackgroundConfig {
  id: string
  name: string
  backgroundColor: string
  backgroundImage?: string
  cardBackground: string
  cardBorder: string
  textColor: string
  textSecondary: string
}

interface DashboardConfig {
  title: string
  theme: string // 图表主题ID
  background: string // 背景ID
  customBackground?: {
    type: 'color' | 'image'
    value: string
  }
}

defineOptions({
  name: 'Dashboard',
})

const props = defineProps<{
  tableId?: string
  initialLayout?: ChartItem[]
}>()

const emit = defineEmits(['save'])

// 仪表盘配置
const dashboardConfig = ref<DashboardConfig>({
  title: '数据仪表盘',
  theme: 'default',
  background: 'light',
})

// 布局数据
const layout = ref<ChartItem[]>([])
const pcLayout = ref<ChartItem[]>([]) // PC端布局
const mobileLayout = ref<ChartItem[]>([]) // 移动端布局

// 当前使用的布局
const currentLayout = computed({
  get: () => previewMode.value === 'mobile' ? mobileLayout.value : pcLayout.value,
  set: (val) => {
    if (previewMode.value === 'mobile') {
      mobileLayout.value = val
    } else {
      pcLayout.value = val
    }
  }
})

// 编辑模式和预览模式
const isEditMode = ref(true)
const previewMode = ref<'pc' | 'mobile'>('pc')

// 计算移动端容器高度
const mobileContainerHeight = computed(() => {
  if (previewMode.value !== 'mobile' || mobileLayout.value.length === 0) return 'auto'
  // 计算所有图表的总高度
  const maxY = Math.max(...mobileLayout.value.map(item => item.y + item.h))
  const rowHeight = 120 // 移动端行高
  const margin = 12 // margin
  return `${maxY * rowHeight + (maxY - 1) * margin + 24}px` // 24px为容器padding
})

// 弹窗状态
const showChartSelector = ref(false)
const showChartConfig = ref(false)
const showDeleteConfirm = ref(false)
const showDataEditor = ref(false) // 数据编辑器
const showThemeSelector = ref(false) // 主题选择器
const currentEditChart = ref<ChartItem | null>(null)
const pendingDeleteChart = ref<ChartItem | null>(null)

// 自定义背景状态
const customBgType = ref<'color' | 'image'>('color')
const customBgValue = ref('#f5f7fa')
const bgFileInputRef = ref<HTMLInputElement>()

// 图表主题配置（仅包含图表颜色）
const chartThemes: ThemeConfig[] = [
  {
    id: 'default',
    name: '默认主题',
    chartColors: ['#3b82f6', '#8b5cf6', '#22c55e', '#eab308', '#f97316', '#ec4899'],
  },
  {
    id: 'warm',
    name: '暖色主题',
    chartColors: ['#f97316', '#fb923c', '#f59e0b', '#fbbf24', '#f87171', '#fb7185'],
  },
  {
    id: 'cool',
    name: '冷色主题',
    chartColors: ['#3b82f6', '#60a5fa', '#06b6d4', '#0891b2', '#14b8a6', '#10b981'],
  },
  {
    id: 'purple',
    name: '紫色主题',
    chartColors: ['#8b5cf6', '#a78bfa', '#c084fc', '#e879f9', '#f472b6', '#ec4899'],
  },
  {
    id: 'green',
    name: '绿色主题',
    chartColors: ['#22c55e', '#10b981', '#14b8a6', '#84cc16', '#a3e635', '#4ade80'],
  },
  {
    id: 'mono',
    name: '黑白主题',
    chartColors: ['#000000', '#404040', '#737373', '#a3a3a3', '#d4d4d4', '#e5e5e5'],
  },
]

// 背景主题配置
const backgroundThemes: BackgroundConfig[] = [
  {
    id: 'light',
    name: '浅色背景',
    backgroundColor: '#f5f7fa',
    cardBackground: '#ffffff',
    cardBorder: 'rgba(0, 0, 0, 0.08)',
    textColor: '#333333',
    textSecondary: '#666666',
  },
  {
    id: 'dark',
    name: '深色背景',
    backgroundColor: '#1a1a1a',
    cardBackground: '#2d2d2d',
    cardBorder: 'rgba(255, 255, 255, 0.1)',
    textColor: '#e5e5e5',
    textSecondary: '#a3a3a3',
  },
  {
    id: 'tech',
    name: '科技蓝',
    backgroundColor: '#0f172a',
    backgroundImage: 'linear-gradient(135deg, #0f172a 0%, #1e293b 100%)',
    cardBackground: 'rgba(30, 41, 59, 0.5)',
    cardBorder: 'rgba(59, 130, 246, 0.3)',
    textColor: '#e0f2fe',
    textSecondary: '#93c5fd',
  },
  {
    id: 'purple-gradient',
    name: '渐变紫',
    backgroundColor: '#1e1b4b',
    backgroundImage: 'linear-gradient(135deg, #1e1b4b 0%, #4c1d95 50%, #6d28d9 100%)',
    cardBackground: 'rgba(139, 92, 246, 0.1)',
    cardBorder: 'rgba(168, 85, 247, 0.3)',
    textColor: '#e9d5ff',
    textSecondary: '#c4b5fd',
  },
  {
    id: 'business',
    name: '商务灰',
    backgroundColor: '#262626',
    cardBackground: '#404040',
    cardBorder: 'rgba(115, 115, 115, 0.3)',
    textColor: '#fafafa',
    textSecondary: '#a3a3a3',
  },
  {
    id: 'green-gradient',
    name: '清新绿',
    backgroundColor: '#f0fdf4',
    backgroundImage: 'linear-gradient(135deg, #f0fdf4 0%, #dcfce7 100%)',
    cardBackground: '#ffffff',
    cardBorder: 'rgba(34, 197, 94, 0.2)',
    textColor: '#064e3b',
    textSecondary: '#16a34a',
  },
]

// 获取当前图表主题
const currentChartTheme = computed(() => {
  return chartThemes.find(t => t.id === dashboardConfig.value.theme) || chartThemes[0]
})

// 获取当前背景主题
const currentBackground = computed(() => {
  return backgroundThemes.find(t => t.id === dashboardConfig.value.background) || backgroundThemes[0]
})

// 图表类型
const chartTypes = [
  {
    type: 'number-card',
    name: '数字卡片',
    description: '展示单个数值指标',
    icon: 'NumberCardIcon',
  },
  {
    type: 'pie',
    name: '饼图',
    description: '展示数据占比关系',
    icon: 'PieChartIcon',
  },
  {
    type: 'bar',
    name: '柱状图',
    description: '对比不同类别数据',
    icon: 'BarChartIcon',
  },
  {
    type: 'line',
    name: '折线图',
    description: '展示数据趋势变化',
    icon: 'LineChartIcon',
  },
  {
    type: 'doughnut',
    name: '环形图',
    description: '展示数据占比（环形）',
    icon: 'DoughnutChartIcon',
  },
  {
    type: 'table',
    name: '数据表格',
    description: '展示详细数据列表',
    icon: 'TableIcon',
  },
]

// 获取图表组件
const getChartComponent = (type: string) => {
  const componentMap: Record<string, any> = {
    'number-card': NumberCard,
    'pie': PieChart,
    'bar': BarChart,
    'line': BarChart, // 暂时用柱状图代替
    'doughnut': DoughnutChart,
    'table': TableChart,
  }
  return componentMap[type] || NumberCard
}

// 切换编辑模式
const toggleEditMode = () => {
  isEditMode.value = !isEditMode.value
  if (!isEditMode.value) {
    MessagePlugin.success('已退出编辑模式')
  }
}

// 添加图表
const handleAddChart = (type: string) => {
  // 计算新图表的y位置（放在最下面）
  const maxY = pcLayout.value.length > 0 
    ? Math.max(...pcLayout.value.map(item => item.y + item.h))
    : 0
  
  const newChart: ChartItem = {
    i: `chart_${Date.now()}`,
    x: 0,
    y: maxY, // 放在所有现有图表下面
    w: type === 'number-card' ? 3 : type === 'table' ? 12 : 6,
    h: type === 'number-card' ? 2 : type === 'table' ? 6 : 4,
    type,
    config: {
      title: getDefaultTitle(type),
      dataSource: 'table1',
      showPagination: type === 'table' ? true : undefined,
      pageSize: type === 'table' ? 10 : undefined,
    },
    data: getDefaultData(type),
  }
  
  pcLayout.value = [...pcLayout.value, newChart]
  generateMobileLayout()
  showChartSelector.value = false
  MessagePlugin.success('图表已添加')
}

// 获取默认标题
const getDefaultTitle = (type: string): string => {
  const titleMap: Record<string, string> = {
    'number-card': '数字指标',
    'pie': '饼图',
    'bar': '柱状图',
    'line': '折线图',
    'doughnut': '环形图',
    'table': '数据表格',
  }
  return titleMap[type] || '图表'
}

// 获取默认数据
const getDefaultData = (type: string): any => {
  // 根据图表类型返回模拟数据
  if (type === 'number-card') {
    return { value: 26, trend: 'up', compare: '--' }
  }
  if (type === 'pie' || type === 'doughnut') {
    return [
      { value: 6, name: '进行中' },
      { value: 4, name: '空值' },
      { value: 6, name: '未开始' },
      { value: 10, name: '已完成' },
    ]
  }
  if (type === 'bar' || type === 'line') {
    return [10, 52, 200, 334, 390, 330, 220]
  }
  if (type === 'table') {
    return {
      columns: [
        { colKey: 'name', title: '姓名', width: 120, align: 'left' },
        { colKey: 'status', title: '状态', width: 100, align: 'center' },
        { colKey: 'value', title: '数值', width: 100, align: 'right' },
        { colKey: 'date', title: '日期', width: 120, align: 'center' },
      ],
      rows: [
        { id: 1, name: '张三', status: '进行中', value: 1250, date: '2024-01-15' },
        { id: 2, name: '李四', status: '已完成', value: 890, date: '2024-01-16' },
        { id: 3, name: '王五', status: '未开始', value: 0, date: '2024-01-17' },
        { id: 4, name: '赵六', status: '进行中', value: 2340, date: '2024-01-18' },
        { id: 5, name: '钱七', status: '已完成', value: 1560, date: '2024-01-19' },
      ],
    }
  }
  return []
}

// 配置图表
const handleConfigChart = (item: ChartItem) => {
  currentEditChart.value = { ...item }
  showChartConfig.value = true
}

// 编辑图表数据
const handleEditChartData = (item: ChartItem) => {
  currentEditChart.value = { ...item }
  showDataEditor.value = true
}

// 删除图表（显示确认弹窗）
const handleDeleteChart = (item: ChartItem) => {
  pendingDeleteChart.value = item
  showDeleteConfirm.value = true
}

// 确认删除
const confirmDelete = () => {
  if (!pendingDeleteChart.value) return
  
  const index = pcLayout.value.findIndex(chart => chart.i === pendingDeleteChart.value!.i)
  if (index > -1) {
    pcLayout.value = pcLayout.value.filter(chart => chart.i !== pendingDeleteChart.value!.i)
    generateMobileLayout()
    MessagePlugin.success('图表已删除')
  }
  
  showDeleteConfirm.value = false
  pendingDeleteChart.value = null
}

// 保存图表配置
const handleSaveChartConfig = () => {
  if (!currentEditChart.value) return
  
  const index = pcLayout.value.findIndex(chart => chart.i === currentEditChart.value!.i)
  if (index > -1) {
    const newLayout = [...pcLayout.value]
    newLayout[index] = { ...currentEditChart.value }
    pcLayout.value = newLayout
  }
  
  showChartConfig.value = false
  MessagePlugin.success('配置已保存')
}

// 保存图表数据
const handleSaveChartData = () => {
  if (!currentEditChart.value) return
  
  const index = pcLayout.value.findIndex(chart => chart.i === currentEditChart.value!.i)
  if (index > -1) {
    const newLayout = [...pcLayout.value]
    newLayout[index] = { ...newLayout[index], data: currentEditChart.value.data }
    pcLayout.value = newLayout
  }
  
  showDataEditor.value = false
  MessagePlugin.success('数据已保存')
}

// 选择图表主题
const handleSelectChartTheme = (themeId: string) => {
  dashboardConfig.value.theme = themeId
  MessagePlugin.success('图表主题已应用')
}

// 选择背景主题
const handleSelectBackground = (bgId: string) => {
  dashboardConfig.value.background = bgId
  // 清除自定义背景
  dashboardConfig.value.customBackground = undefined
  MessagePlugin.success('背景已应用')
}

// 触发图片上传
const triggerImageUpload = () => {
  bgFileInputRef.value?.click()
}

// 处理图片上传
const handleBgImageUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  
  // 验证文件类型
  if (!file.type.startsWith('image/')) {
    MessagePlugin.warning('请选择图片文件')
    return
  }
  
  // 验证文件大小（5MB）
  if (file.size > 5 * 1024 * 1024) {
    MessagePlugin.warning('图片大小不能超过5MB')
    return
  }
  
  try {
    MessagePlugin.loading('上传中...')
    const response = await uploadFile(file)
    if (response?.url) {
      customBgValue.value = response.url
      MessagePlugin.success('图片上传成功')
    }
  } catch (error: any) {
    console.error('图片上传失败:', error)
    MessagePlugin.error(error.message || '图片上传失败')
  } finally {
    // 清空文件输入
    if (bgFileInputRef.value) {
      bgFileInputRef.value.value = ''
    }
  }
}

// 应用自定义背景
const applyCustomBackground = () => {
  if (!customBgValue.value) {
    MessagePlugin.warning('请输入背景值')
    return
  }
  
  dashboardConfig.value.customBackground = {
    type: customBgType.value,
    value: customBgValue.value,
  }
  MessagePlugin.success('自定义背景已应用')
}

// 清除自定义背景
const clearCustomBackground = () => {
  dashboardConfig.value.customBackground = undefined
  MessagePlugin.success('已清除自定义背景')
}

// 获取自定义背景样式
const getCustomBackgroundStyle = () => {
  if (!dashboardConfig.value.customBackground) return {}
  
  const { type, value } = dashboardConfig.value.customBackground
  if (type === 'color') {
    return { backgroundColor: value }
  } else {
    // 支持图片URL和渐变
    if (value.startsWith('linear-gradient') || value.startsWith('radial-gradient')) {
      return { background: value }
    }
    return { 
      backgroundImage: `url(${value})`,
      backgroundSize: 'cover',
      backgroundPosition: 'center',
    }
  }
}

// 计算画布背景样式
const canvasBackgroundStyle = computed(() => {
  // 自定义背景优先
  if (dashboardConfig.value.customBackground) {
    return getCustomBackgroundStyle()
  }
  
  // 使用背景主题
  const bg = currentBackground.value
  if (bg.backgroundImage) {
    return { background: bg.backgroundImage }
  }
  return { backgroundColor: bg.backgroundColor }
})

// 全屏
const handleFullscreen = () => {
  const elem = document.querySelector('.dashboard-canvas') as any
  if (elem.requestFullscreen) {
    elem.requestFullscreen()
  }
}

// 保存仪表盘
const handleSave = () => {
  emit('save', {
    config: dashboardConfig.value,
    layout: layout.value,
  })
  MessagePlugin.success('仪表盘已保存')
}

// 初始化
onMounted(() => {
  if (props.initialLayout) {
    pcLayout.value = props.initialLayout
    layout.value = props.initialLayout
    // 生成移动端布局
    generateMobileLayout()
  }
})

// 监听预览模式切换
watch(previewMode, (newMode, oldMode) => {
  if (newMode === 'mobile' && oldMode === 'pc') {
    // PC切换到移动端，生成移动端布局
    generateMobileLayout()
    // 移动端自动退出编辑模式
    if (isEditMode.value) {
      isEditMode.value = false
      MessagePlugin.info('移动端仅支持预览模式')
    }
  }
})

// 监听PC布局变化，同步到layout（用于保存）
watch(pcLayout, (newVal) => {
  layout.value = newVal
}, { deep: true })

// 生成移动端布局（从上到下，宽度铺满）
const generateMobileLayout = () => {
  let currentY = 0
  mobileLayout.value = pcLayout.value.map((item) => {
    // 根据图表类型设置不同高度
    let mobileHeight = 3 // 默认高度
    if (item.type === 'number-card') {
      mobileHeight = 2
    } else if (item.type === 'table') {
      mobileHeight = 5 // 表格需要更高的高度
    }
    
    const mobileItem = {
      ...item,
      x: 0,
      y: currentY,
      w: 1, // 移动端宽度占满
      h: mobileHeight,
    }
    currentY += mobileHeight // 累加y位置
    return mobileItem
  })
}

// 添加数据项（饼图/环形图）
const addDataItem = () => {
  if (!currentEditChart.value) return
  if (!Array.isArray(currentEditChart.value.data)) {
    currentEditChart.value.data = []
  }
  currentEditChart.value.data.push({ name: '新数据', value: 0 })
}

// 删除数据项
const removeDataItem = (index: number) => {
  if (!currentEditChart.value || !Array.isArray(currentEditChart.value.data)) return
  currentEditChart.value.data.splice(index, 1)
}

// 处理柱状图数据变化
const handleBarDataChange = (val: string | number) => {
  if (!currentEditChart.value) return
  const strVal = String(val)
  currentEditChart.value.data = strVal.split(',').map(v => Number(v.trim())).filter(v => !isNaN(v))
}
</script>

<style scoped lang="less">
.dashboard-designer {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;

  .dashboard-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 24px;
    background: #fff;
    border-bottom: 1px solid #e8e8e8;

    .header-left {
      .dashboard-title-input {
        width: 300px;
        
        :deep(.t-input) {
          border: none;
          font-size: 16px;
          font-weight: 600;
        }
      }
    }

    .header-right {
      display: flex;
      align-items: center;
      gap: 8px;
      
      .preview-switch {
        margin-right: 8px;
        
        :deep(.t-radio-button) {
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }
    }
  }

  .dashboard-canvas {
    flex: 1;
    overflow: auto;
    padding: 24px;
    position: relative;
    display: flex;
    justify-content: center;
    align-items: flex-start; // 顶部对齐
    transition: background 0.3s ease; // 背景切换动画
    
    .canvas-container {
      width: 100%;
      max-width: 1400px;
      transition: all 0.3s;
      min-height: 100%; // 确保容器至少填满高度
    }
    
    &.preview-mobile {
      padding: 24px 16px; // 移动端减小左右内边距
      
      .canvas-container {
        max-width: 375px;
        min-height: 0; // 移动端完全根据内容自适应
        height: auto;
        background: #fff;
        border-radius: 12px;
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
        padding: 12px;
      }
    }

    // vue-grid-layout 样式覆盖
    :deep(.vue-grid-layout) {
      min-height: 400px;
      position: relative;
    }
    
    // 移动端布局样式优化
    &.preview-mobile :deep(.vue-grid-layout) {
      min-height: 0;
      height: auto;
      padding-bottom: 12px; // 底部留白
    }
    
    :deep(.vue-grid-item) {
      background: transparent;
      touch-action: none;
      transition: all 200ms ease;
      transition-property: left, top;
      
      &.vue-grid-placeholder {
        background: rgba(0, 82, 217, 0.1);
        border: 2px dashed #0052d9;
        border-radius: 8px;
        transition: all 150ms ease;
        z-index: 2;
      }
      
      &.resizing {
        opacity: 0.9;
        z-index: 3;
      }
      
      &.vue-draggable-dragging {
        transition: none;
        z-index: 3;
      }
      
      // 拖动句柄
      &.cssTransforms {
        transition-property: transform;
      }
      
      // 移动端静态项
      &.static {
        background: none;
      }
    }

    .chart-wrapper {
      position: relative;
      width: 100%;
      height: 100%;
      background: #fff;
      border-radius: 8px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
      transition: all 0.3s;
      overflow: hidden;
      display: flex;
      flex-direction: column;
      
      // 确保图表组件自适应容器
      > * {
        flex: 1;
        min-height: 0;
        min-width: 0;
      }

      &.editing {
        border: 2px solid #e8e8e8;
        
        &:hover {
          border-color: #0052d9;
        }
      }

      .chart-actions {
        position: absolute;
        top: 8px;
        right: 8px;
        z-index: 10;
        display: flex;
        gap: 4px;
        opacity: 0;
        transition: opacity 0.3s;
      }

      &:hover {
        .chart-actions {
          opacity: 1;
        }
      }
    }
    
    // 移动端图表样式优化
    &.preview-mobile .chart-wrapper {
      box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
      
      // 移动端图表间距由margin控制
      margin-bottom: 0;
    }

    .empty-dashboard {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      text-align: center;

      .empty-content {
        color: #999;

        .empty-title {
          font-size: 18px;
          font-weight: 500;
          color: #333;
          margin: 16px 0 8px;
        }

        .empty-desc {
          font-size: 14px;
          margin-bottom: 24px;
        }
      }
    }
  }

  .chart-selector {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
    padding: 16px 0;

    .chart-type-card {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 16px;
      border: 1px solid #e8e8e8;
      border-radius: 8px;
      cursor: pointer;
      transition: all 0.3s;

      &:hover {
        border-color: #0052d9;
        box-shadow: 0 4px 12px rgba(0, 82, 217, 0.15);
      }

      .chart-icon {
        width: 48px;
        height: 48px;
        display: flex;
        align-items: center;
        justify-content: center;
        background: #f5f7fa;
        border-radius: 8px;
        font-size: 24px;
      }

      .chart-info {
        flex: 1;

        .chart-name {
          font-size: 14px;
          font-weight: 500;
          color: #333;
          margin-bottom: 4px;
        }

        .chart-desc {
          font-size: 12px;
          color: #999;
        }
      }
    }
  }

  .chart-config-form {
    padding: 16px 0;
  }
  
  .delete-confirm-content {
    padding: 8px 0;

    p {
      margin: 0 0 12px;
      font-size: 14px;
      line-height: 1.6;
      color: #333;

      &:last-child {
        margin-bottom: 0;
      }

      strong {
        font-weight: 600;
        color: #0052d9;
      }
    }

    .warning-text {
      color: #e34d59;
      font-size: 13px;
    }
  }
  
  .data-editor-form {
    padding: 16px 0;
    
    .data-list {
      .data-item {
        display: flex;
        align-items: center;
        margin-bottom: 8px;
      }
    }
  }
  
  // 主题选择器样式
  .theme-selector {
    .theme-section {
      margin-bottom: 32px;
      
      &:last-child {
        margin-bottom: 0;
      }
      
      .section-title {
        font-size: 16px;
        font-weight: 600;
        color: #333;
        margin: 0 0 16px;
      }
    }
    
    .theme-grid {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 16px;
      
      &.chart-themes {
        .theme-card {
          .chart-preview {
            background: linear-gradient(135deg, #f5f7fa 0%, #e5e7eb 100%);
            
            .preview-colors {
              display: flex;
              gap: 4px;
              width: 100%;
              
              .color-bar {
                flex: 1;
                height: 60px;
                border-radius: 4px;
              }
            }
          }
        }
      }
      
      .theme-card {
        cursor: pointer;
        border-radius: 8px;
        overflow: hidden;
        border: 2px solid transparent;
        transition: all 0.3s;
        
        &:hover {
          border-color: #0052d9;
          transform: translateY(-2px);
          box-shadow: 0 4px 12px rgba(0, 82, 217, 0.15);
        }
        
        &.active {
          border-color: #0052d9;
          box-shadow: 0 4px 12px rgba(0, 82, 217, 0.2);
        }
        
        .theme-preview {
          height: 120px;
          padding: 16px;
          display: flex;
          align-items: center;
          justify-content: center;
          position: relative;
          
          .preview-card {
            width: 100%;
            height: 80px;
            border-radius: 6px;
            padding: 12px;
            display: flex;
            flex-direction: column;
            justify-content: space-between;
            
            .preview-title {
              font-size: 12px;
              font-weight: 500;
            }
            
            .preview-text {
              font-size: 11px;
            }
            
            .preview-colors {
              display: flex;
              gap: 6px;
              
              .color-dot {
                width: 16px;
                height: 16px;
                border-radius: 50%;
              }
            }
          }
        }
        
        .theme-name {
          padding: 12px;
          font-size: 14px;
          font-weight: 500;
          color: #333;
          background: #fff;
          display: flex;
          align-items: center;
          justify-content: space-between;
          
          .check-icon {
            color: #0052d9;
            font-size: 18px;
          }
        }
      }
    }
    
    .custom-background {
      .bg-input-wrapper {
        display: flex;
        gap: 12px;
        margin-top: 12px;
        align-items: center;
        
        .color-picker {
          width: 40px;
          height: 32px;
          border: 1px solid #ddd;
          border-radius: 4px;
          cursor: pointer;
        }
      }
      
      .custom-bg-preview {
        margin-top: 16px;
        padding: 16px;
        border: 1px solid #e8e8e8;
        border-radius: 8px;
        background: #f5f7fa;
        
        .preview-box {
          height: 100px;
          border-radius: 6px;
          display: flex;
          align-items: center;
          justify-content: center;
          color: #fff;
          font-size: 14px;
          font-weight: 500;
          text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
          margin-bottom: 12px;
        }
      }
    }
  }
}
</style>
