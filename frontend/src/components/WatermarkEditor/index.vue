<template>
  <div class="watermark-editor">
    <!-- 触发按钮 -->
    <slot name="trigger" :open="openDrawer">
      <t-button theme="primary" @click="openDrawer">
        <template #icon>
          <t-icon name="edit" />
        </template>
        编辑水印
      </t-button>
    </slot>

    <!-- 右侧抽屉 -->
    <t-drawer
      v-model:visible="drawerVisible"
      header="水印配置"
      size="480px"
      :footer="true"
      @confirm="handleConfirm"
      @cancel="handleCancel"
    >
      <div class="watermark-config">
        <!-- 配置表单 -->
        <div class="config-section">
          <!-- 水印内容 -->
          <div class="config-group">
            <h3 class="section-title">水印内容</h3>
            <t-form-item label="水印文本">
              <t-input
                v-model="config.content"
                placeholder="请输入水印文本"
                :maxlength="50"
              />
            </t-form-item>
            
            <t-form-item label="多行水印">
              <t-switch v-model="config.isMultiLine" />
              <span class="form-tip">开启后可输入多行文本</span>
            </t-form-item>

            <t-form-item v-if="config.isMultiLine" label="第二行文本">
              <t-input
                v-model="config.content2"
                placeholder="请输入第二行文本"
                :maxlength="50"
              />
            </t-form-item>
          </div>

          <!-- 样式配置 -->
          <div class="config-group">
            <h3 class="section-title">样式配置</h3>
            
            <t-form-item label="字体大小">
              <t-slider
                v-model="config.fontSize"
                :min="12"
                :max="48"
                :step="1"
                :marks="{ 12: '12px', 24: '24px', 36: '36px', 48: '48px' }"
              />
              <span class="value-label">{{ config.fontSize }}px</span>
            </t-form-item>

            <t-form-item label="字体颜色">
              <t-color-picker v-model="config.fontColor" />
            </t-form-item>

            <t-form-item label="透明度">
              <t-slider
                v-model="config.alpha"
                :min="0"
                :max="1"
                :step="0.1"
                :marks="{ 0: '0%', 0.5: '50%', 1: '100%' }"
              />
              <span class="value-label">{{ Math.round(config.alpha * 100) }}%</span>
            </t-form-item>

            <t-form-item label="旋转角度">
              <t-slider
                v-model="config.rotate"
                :min="-90"
                :max="90"
                :step="15"
                :marks="{ '-90': '-90°', '-45': '-45°', 0: '0°', 45: '45°', 90: '90°' }"
              />
              <span class="value-label">{{ config.rotate }}°</span>
            </t-form-item>
          </div>

          <!-- 布局配置 -->
          <div class="config-group">
            <h3 class="section-title">布局配置</h3>
            
            <t-form-item label="水平间距">
              <t-slider
                v-model="config.gapX"
                :min="50"
                :max="300"
                :step="10"
              />
              <span class="value-label">{{ config.gapX }}px</span>
            </t-form-item>

            <t-form-item label="垂直间距">
              <t-slider
                v-model="config.gapY"
                :min="50"
                :max="300"
                :step="10"
              />
              <span class="value-label">{{ config.gapY }}px</span>
            </t-form-item>

            <t-form-item label="水印宽度">
              <t-slider
                v-model="config.width"
                :min="100"
                :max="400"
                :step="10"
              />
              <span class="value-label">{{ config.width }}px</span>
            </t-form-item>

            <t-form-item label="水印高度">
              <t-slider
                v-model="config.height"
                :min="50"
                :max="200"
                :step="10"
              />
              <span class="value-label">{{ config.height }}px</span>
            </t-form-item>
          </div>

          <!-- 快捷预设 -->
          <div class="config-group">
            <h3 class="section-title">快捷预设</h3>
            <div class="preset-list">
              <t-button
                v-for="preset in presets"
                :key="preset.name"
                variant="outline"
                size="small"
                @click="applyPreset(preset)"
              >
                {{ preset.name }}
              </t-button>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="drawer-footer">
          <t-button variant="outline" @click="handleReset">重置</t-button>
          <t-button theme="primary" @click="handleConfirm">确定</t-button>
        </div>
      </template>
    </t-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'

interface WatermarkConfig {
  content: string
  content2: string
  isMultiLine: boolean
  fontSize: number
  fontColor: string
  alpha: number
  rotate: number
  gapX: number
  gapY: number
  width: number
  height: number
}

interface WatermarkPreset {
  name: string
  config: Partial<WatermarkConfig>
}

interface Props {
  defaultConfig?: Partial<WatermarkConfig>
}

interface Emits {
  (e: 'confirm', config: WatermarkConfig): void
  (e: 'cancel'): void
  (e: 'update:config', config: WatermarkConfig): void
}

const props = withDefaults(defineProps<Props>(), {
  defaultConfig: () => ({})
})

const emit = defineEmits<Emits>()

// 默认配置
const defaultWatermarkConfig: WatermarkConfig = {
  content: '机密文档',
  content2: '请勿外传',
  isMultiLine: false,
  fontSize: 16,
  fontColor: 'rgba(0,0,0,0.15)',
  alpha: 0.15,
  rotate: -22,
  gapX: 100,
  gapY: 100,
  width: 120,
  height: 64,
}

// 当前配置
const config = ref<WatermarkConfig>({
  ...defaultWatermarkConfig,
  ...props.defaultConfig,
})

// 备份配置（用于重置）
const backupConfig = ref<WatermarkConfig>({ ...config.value })

// 抽屉显示状态
const drawerVisible = ref(false)

// 预设模板
const presets: WatermarkPreset[] = [
  {
    name: '默认',
    config: defaultWatermarkConfig,
  },
  {
    name: '机密文档',
    config: {
      content: '机密文档',
      content2: '请勿外传',
      isMultiLine: true,
      fontSize: 16,
      alpha: 0.15,
      rotate: -22,
    },
  },
  {
    name: '内部资料',
    config: {
      content: '内部资料',
      content2: '仅供参考',
      isMultiLine: true,
      fontSize: 18,
      alpha: 0.2,
      rotate: -30,
    },
  },
  {
    name: '草稿版本',
    config: {
      content: 'DRAFT',
      isMultiLine: false,
      fontSize: 24,
      alpha: 0.1,
      rotate: -45,
      fontColor: 'rgba(255,0,0,0.1)',
    },
  },
  {
    name: '版权保护',
    config: {
      content: '© Copyright',
      isMultiLine: false,
      fontSize: 14,
      alpha: 0.1,
      rotate: 0,
    },
  },
]

// 计算预览内容
const previewContent = computed(() => {
  if (config.value.isMultiLine && config.value.content2) {
    return `${config.value.content}\n${config.value.content2}`
  }
  return config.value.content
})

// 计算水印组件属性
const watermarkProps = computed(() => {
  // 将 fontColor 转换为正确格式
  let fontColor = config.value.fontColor
  
  // 如果是 hex 颜色，转换为 rgba
  if (fontColor.startsWith('#')) {
    const r = parseInt(fontColor.slice(1, 3), 16)
    const g = parseInt(fontColor.slice(3, 5), 16)
    const b = parseInt(fontColor.slice(5, 7), 16)
    fontColor = `rgba(${r},${g},${b},${config.value.alpha})`
  } else if (fontColor.startsWith('rgb(')) {
    // 将 rgb 转换为 rgba
    fontColor = fontColor.replace('rgb(', 'rgba(').replace(')', `,${config.value.alpha})`)
  } else if (fontColor.startsWith('rgba(')) {
    // 替换现有的 alpha 值
    fontColor = fontColor.replace(/,\s*[\d.]+\s*\)/, `,${config.value.alpha})`)
  }

  console.log('水印组件配置：', {
    config: config.value,
    previewContent: previewContent.value,
  })

  return {
    y: config.value.gapY,
    x: config.value.gapX,
    rotate: config.value.rotate,
    width: config.value.width,
    height: config.value.height,
    watermarkContent: {
      text: previewContent.value,
      fontSize: config.value.fontSize,
      fontColor: fontColor,
    },
  }
})

// 打开抽屉
const openDrawer = () => {
  backupConfig.value = { ...config.value }
  drawerVisible.value = true
}

// 应用预设
const applyPreset = (preset: WatermarkPreset) => {
  config.value = {
    ...config.value,
    ...preset.config,
  }
  MessagePlugin.success(`已应用「${preset.name}」预设`)
}

// 重置配置
const handleReset = () => {
  config.value = { ...defaultWatermarkConfig }
  MessagePlugin.success('已重置为默认配置')
}

// 确认
const handleConfirm = () => {
  emit('confirm', config.value)
  emit('update:config', config.value)
  drawerVisible.value = false
  MessagePlugin.success('水印配置已保存')
}

// 取消
const handleCancel = () => {
  config.value = { ...backupConfig.value }
  emit('cancel')
  drawerVisible.value = false
}

// 监听配置变化，实时更新
watch(() => config.value, (newConfig) => {
  emit('update:config', newConfig)
}, { deep: true })

// 暴露方法
defineExpose({
  openDrawer,
  getConfig: () => config.value,
  setConfig: (newConfig: Partial<WatermarkConfig>) => {
    config.value = { ...config.value, ...newConfig }
  },
})
</script>

<style scoped lang="less">
.watermark-editor {
  display: inline-block;
}

.watermark-config {
  .config-section {
    .config-group {
      margin-bottom: 32px;
      padding-bottom: 24px;
      border-bottom: 1px solid #f0f0f0;
      
      &:last-child {
        border-bottom: none;
      }
      
      .section-title {
        font-size: 14px;
        font-weight: 500;
        color: #333;
        margin: 0 0 16px;
      }
      
      :deep(.t-form-item) {
        margin-bottom: 20px;
        
        .t-form__label {
          font-size: 13px;
          color: #666;
        }
        
        .t-form__controls {
          position: relative;
          
          .value-label {
            position: absolute;
            right: 0;
            top: 50%;
            transform: translateY(-50%);
            font-size: 12px;
            color: #999;
            background: #fff;
            padding-left: 8px;
          }
        }
      }
      
      .form-tip {
        font-size: 12px;
        color: #999;
        margin-left: 8px;
      }
      
      .preset-list {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
      }
    }
  }
}

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
