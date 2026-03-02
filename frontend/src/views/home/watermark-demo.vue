<template>
  <div class="watermark-demo-page">
    <h1 class="page-title">水印编辑器演示</h1>

    <!-- 实时预览 -->
    <section class="demo-section">
      <h2 class="section-title">水印实时预览</h2>
      <div class="demo-content full-width">
        <div class="demo-item">
          <h3>当前水印配置</h3>
          <p class="demo-desc">配置水印后实时查看效果</p>
          
          <div class="watermark-controls">
            <WatermarkEditor 
              ref="editorRef"
              :default-config="currentConfig"
              @confirm="handleWatermarkChange"
            />
            <t-button 
              theme="default" 
              variant="outline"
              style="margin-left: 12px"
              @click="clearWatermark"
            >
              清除水印
            </t-button>
          </div>

          <div class="watermark-preview">
            <t-watermark v-if="showWatermark" v-bind="watermarkProps">
              <div class="preview-area">
                <div class="preview-card">
                  <h4>文档标题</h4>
                  <p>这是一个示例文档内容区域</p>
                  <p>水印会覆盖在整个区域上</p>
                  <p>可以通过右侧编辑器调整水印样式</p>
                </div>
              </div>
            </t-watermark>
            <div v-else class="preview-area empty">
              <div class="empty-text">暂无水印，请点击编辑水印按钮配置</div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import WatermarkEditor from '@/components/WatermarkEditor/index.vue'

// 编辑器引用
const editorRef = ref()

// 当前配置
const currentConfig = ref<any>({
  content: '示例水印',
  fontSize: 16,
  alpha: 0.15,
  rotate: -22,
})

// 是否显示水印
const showWatermark = ref(false)

// 水印内容
const watermarkContent = computed(() => {
  if (currentConfig.value.isMultiLine && currentConfig.value.content2) {
    return `${currentConfig.value.content}\n${currentConfig.value.content2}`
  }
  return currentConfig.value.content
})

// 水印属性
const watermarkProps = computed(() => {
  console.log('水印配置：', currentConfig.value)
  
  let fontColor = currentConfig.value.fontColor || 'rgba(0,0,0,0.15)'
  
  // 颜色转换逻辑
  if (fontColor.startsWith('#')) {
    const r = parseInt(fontColor.slice(1, 3), 16)
    const g = parseInt(fontColor.slice(3, 5), 16)
    const b = parseInt(fontColor.slice(5, 7), 16)
    fontColor = `rgba(${r},${g},${b},${currentConfig.value.alpha || 0.15})`
  }

  const props = {
    y: currentConfig.value.gapY || 100,
    x: currentConfig.value.gapX || 100,
    rotate: currentConfig.value.rotate || -22,
    width: currentConfig.value.width || 120,
    height: currentConfig.value.height || 64,
    watermarkContent: {
      text: watermarkContent.value,
      fontSize: currentConfig.value.fontSize || 16,
      fontColor: fontColor,
    },
  }
  
  console.log('水印属性：', props)
  return props
})

// 处理水印变更
const handleWatermarkChange = (newConfig: any) => {
  console.log('水印配置变更：', newConfig)
  currentConfig.value = newConfig
  showWatermark.value = true
}

// 清除水印
const clearWatermark = () => {
  showWatermark.value = false
  MessagePlugin.info('水印已清除')
}
</script>

<style scoped lang="less">
.watermark-demo-page {
  min-height: 100%;
  background: #ffffff;
  border-radius: 8px;
  padding: 27px;
  
  .page-title {
    font-size: 24px;
    font-weight: 600;
    color: #333;
    margin: 0 0 32px;
    padding-bottom: 16px;
    border-bottom: 2px solid #e8e8e8;
  }
  
  .demo-section {
    margin-bottom: 48px;
    
    &:last-child {
      margin-bottom: 0;
    }
    
    .section-title {
      font-size: 18px;
      font-weight: 600;
      color: #333;
      margin: 0 0 24px;
      padding-left: 12px;
      border-left: 4px solid #4a7ff7;
    }
    
    .demo-content {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
      gap: 24px;
      
      &.full-width {
        grid-template-columns: 1fr;
      }
      
      .demo-item {
        padding: 20px;
        background: #fafafa;
        border-radius: 8px;
        border: 1px solid #e8e8e8;
        
        h3 {
          font-size: 14px;
          font-weight: 500;
          color: #666;
          margin: 0 0 8px;
        }
        
        .demo-desc {
          font-size: 12px;
          color: #999;
          margin: 0 0 16px;
        }
        
        .watermark-controls {
          display: flex;
          align-items: center;
          margin-bottom: 20px;
        }
        
        .watermark-preview {
          margin-top: 20px;
          
          .preview-area {
            min-height: 300px;
            border: 1px solid #e8e8e8;
            border-radius: 4px;
            background: #fff;
            
            &.empty {
              display: flex;
              align-items: center;
              justify-content: center;
              
              .empty-text {
                color: #999;
                font-size: 14px;
              }
            }
            
            .preview-card {
              padding: 32px;
              
              h4 {
                font-size: 20px;
                color: #333;
                margin: 0 0 16px;
              }
              
              p {
                font-size: 14px;
                color: #666;
                line-height: 1.6;
                margin: 8px 0;
              }
            }
          }
        }
        
        .scene-preview {
          margin-top: 16px;
          border: 1px solid #e8e8e8;
          border-radius: 4px;
          overflow: hidden;
          
          .document-content {
            padding: 24px;
            background: #fff;
            
            h4 {
              font-size: 16px;
              color: #333;
              margin: 0 0 12px;
            }
            
            p {
              font-size: 14px;
              color: #666;
              margin: 6px 0;
            }
          }
          
          .image-content {
            img {
              width: 100%;
              height: auto;
              display: block;
            }
          }
        }
      }
    }
  }
}
</style>
