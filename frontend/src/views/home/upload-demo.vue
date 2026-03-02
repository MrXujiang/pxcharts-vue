<template>
  <div class="upload-demo-page">
    <h1 class="page-title">文件上传组件演示</h1>
    
    <!-- 卡片上传模式 -->
    <section class="demo-section">
      <h2 class="section-title">1. 卡片上传模式（支持拖拽）</h2>
      <div class="demo-content">
        <div class="demo-item">
          <h3>单文件上传</h3>
          <FileUpload 
            accept="image/*"
            :max-size="5"
            upload-text="点击或拖拽图片到此区域"
            upload-hint="支持 jpg、png、gif 格式，单个文件不超过 5MB"
            @success="handleSuccess"
            @error="handleError"
          />
        </div>
        
        <div class="demo-item">
          <h3>多文件上传（最多5个）</h3>
          <FileUpload 
            accept="image/*"
            multiple
            :max="5"
            :max-size="10"
            upload-text="点击或拖拽多张图片到此区域"
            upload-hint="支持批量上传，最多 5 张，单个文件不超过 10MB"
          />
        </div>
      </div>
    </section>

    <!-- 按钮上传模式 -->
    <section class="demo-section">
      <h2 class="section-title">2. 按钮上传模式</h2>
      <div class="demo-content">
        <div class="demo-item">
          <h3>文档上传</h3>
          <FileUpload 
            mode="button"
            accept=".pdf,.doc,.docx,.xls,.xlsx"
            multiple
            button-text="上传文档"
            button-icon="file-add"
            button-theme="primary"
            :max-size="20"
          />
        </div>
        
        <div class="demo-item">
          <h3>附件上传（Outline 样式）</h3>
          <FileUpload 
            mode="button"
            accept="*"
            multiple
            button-text="添加附件"
            button-icon="attach"
            button-theme="default"
            button-variant="outline"
            :max-size="50"
          />
        </div>

        <div class="demo-item">
          <h3>视频上传</h3>
          <FileUpload 
            mode="button"
            accept="video/*"
            button-text="上传视频"
            button-icon="video"
            button-theme="success"
            :max-size="100"
          />
        </div>
      </div>
    </section>

    <!-- 图片墙模式 -->
    <section class="demo-section">
      <h2 class="section-title">3. 图片墙模式</h2>
      <div class="demo-content">
        <div class="demo-item">
          <h3>多图上传（最多9张）</h3>
          <FileUpload 
            mode="image-wall"
            accept="image/*"
            multiple
            :max="9"
            :max-size="5"
            @preview="handlePreview"
          />
        </div>
        
        <div class="demo-item">
          <h3>头像上传（单张）</h3>
          <FileUpload 
            mode="image-wall"
            accept="image/*"
            :max="1"
            :max-size="2"
            :show-file-list="false"
          />
        </div>
      </div>
    </section>

    <!-- 高级用法 -->
    <section class="demo-section">
      <h2 class="section-title">4. 高级用法</h2>
      <div class="demo-content">
        <div class="demo-item">
          <h3>链接模式（轻量附件展示）</h3>
          <FileUpload 
            mode="link"
            accept="*"
            multiple
            :max="10"
            :max-size="20"
          />
        </div>
        
        <div class="demo-item">
          <h3>只显示文件列表（隐藏上传区域）</h3>
          <FileUpload 
            :default-file-list="defaultFiles"
            :hide-upload-area="true"
            :show-file-list="true"
          />
        </div>
        
        <div class="demo-item">
          <h3>手动控制上传</h3>
          <div class="manual-upload">
            <FileUpload 
              ref="uploadRef"
              mode="button"
              accept="*"
              multiple
              button-text="选择文件"
              button-variant="outline"
            />
            <t-button 
              theme="danger" 
              variant="outline"
              style="margin-left: 12px"
              @click="clearFiles"
            >
              清空文件
            </t-button>
          </div>
        </div>

        <div class="demo-item">
          <h3>压缩包上传</h3>
          <FileUpload 
            mode="button"
            accept=".zip,.rar,.7z"
            button-text="上传压缩包"
            button-icon="folder-zip"
            button-theme="warning"
            :max-size="100"
          />
        </div>
      </div>
    </section>

    <!-- 实际应用示例 -->
    <section class="demo-section">
      <h2 class="section-title">5. 实际应用示例</h2>
      <div class="demo-content">
        <div class="demo-item full-width">
          <h3>表单中的文件上传</h3>
          <t-form :data="formData" label-width="100px">
            <t-form-item label="项目名称">
              <t-input v-model="formData.name" placeholder="请输入项目名称" />
            </t-form-item>
            
            <t-form-item label="项目描述">
              <t-textarea 
                v-model="formData.description" 
                placeholder="请输入项目描述"
                :autosize="{ minRows: 3, maxRows: 6 }"
              />
            </t-form-item>
            
            <t-form-item label="项目封面">
              <FileUpload 
                mode="image-wall"
                accept="image/*"
                :max="1"
                :max-size="5"
                @success="handleCoverUpload"
              />
            </t-form-item>
            
            <t-form-item label="项目附件">
              <FileUpload 
                mode="button"
                accept="*"
                multiple
                :max="10"
                button-text="添加附件"
                button-variant="outline"
                @update:fileList="handleAttachmentsUpdate"
              />
            </t-form-item>
            
            <t-form-item>
              <t-button theme="primary" @click="submitForm">提交</t-button>
              <t-button theme="default" variant="outline" style="margin-left: 12px">取消</t-button>
            </t-form-item>
          </t-form>
        </div>
      </div>
    </section>

    <!-- 事件监听演示 -->
    <section class="demo-section">
      <h2 class="section-title">6. 事件监听</h2>
      <div class="demo-content">
        <div class="demo-item">
          <h3>监听所有事件</h3>
          <FileUpload 
            mode="button"
            accept="*"
            multiple
            button-text="上传文件（查看控制台）"
            @success="handleSuccess"
            @error="handleError"
            @remove="handleRemove"
            @update:fileList="handleFileListUpdate"
          />
          
          <div class="event-log">
            <h4>事件日志：</h4>
            <div class="log-list">
              <div v-for="(log, index) in eventLogs" :key="index" class="log-item">
                <span class="log-time">{{ log.time }}</span>
                <span class="log-event" :class="`log-${log.type}`">{{ log.event }}</span>
                <span class="log-message">{{ log.message }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import FileUpload from '@/components/FileUpload/index.vue'

// 上传组件引用
const uploadRef = ref()

// 默认文件列表
const defaultFiles = [
  {
    id: '1',
    name: '示例文档.pdf',
    size: 1024 * 1024 * 2,
    type: 'application/pdf',
    url: 'https://example.com/sample.pdf',
  },
  {
    id: '2',
    name: '示例图片.jpg',
    size: 1024 * 500,
    type: 'image/jpeg',
    url: 'https://via.placeholder.com/150',
  },
]

// 表单数据
const formData = reactive({
  name: '',
  description: '',
  cover: '',
  attachments: [] as any[],
})

// 事件日志
const eventLogs = ref<Array<{ time: string; event: string; message: string; type: string }>>([])

const addLog = (event: string, message: string, type: string = 'info') => {
  const now = new Date()
  const time = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}:${now.getSeconds().toString().padStart(2, '0')}`
  
  eventLogs.value.unshift({
    time,
    event,
    message,
    type,
  })
  
  // 只保留最近20条
  if (eventLogs.value.length > 20) {
    eventLogs.value = eventLogs.value.slice(0, 20)
  }
}

// 上传成功
const handleSuccess = (file: any) => {
  console.log('上传成功:', file)
  addLog('success', `${file.name} 上传成功`, 'success')
  // 添加成功提示
  MessagePlugin.success(`${file.name} 上传成功`)
}

// 上传失败
const handleError = (file: any, error: any) => {
  console.error('上传失败:', file, error)
  addLog('error', `${file.name} 上传失败: ${error.message}`, 'error')
}

// 移除文件
const handleRemove = (file: any, index: number) => {
  console.log('移除文件:', file, index)
  addLog('remove', `移除文件: ${file.name}`, 'warning')
}

// 文件列表更新
const handleFileListUpdate = (files: any[]) => {
  console.log('文件列表更新:', files)
  addLog('update', `文件列表更新，当前 ${files.length} 个文件`, 'info')
}

// 预览图片
const handlePreview = (file: any) => {
  console.log('预览图片:', file)
  addLog('preview', `预览图片: ${file.name}`, 'info')
}

// 清空文件
const clearFiles = () => {
  if (uploadRef.value) {
    uploadRef.value.clearFiles()
    MessagePlugin.success('文件已清空')
    addLog('clear', '清空所有文件', 'warning')
  }
}

// 封面上传成功
const handleCoverUpload = (file: any) => {
  formData.cover = file.url
  console.log('封面上传成功:', file.url)
}

// 附件列表更新
const handleAttachmentsUpdate = (files: any[]) => {
  formData.attachments = files
  console.log('附件列表更新:', files)
}

// 提交表单
const submitForm = () => {
  console.log('提交表单:', formData)
  MessagePlugin.success('表单提交成功')
  addLog('submit', `提交表单: ${formData.name}`, 'success')
}
</script>

<style scoped lang="less">
.upload-demo-page {
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
      
      .demo-item {
        padding: 20px;
        background: #fafafa;
        border-radius: 8px;
        border: 1px solid #e8e8e8;
        
        &.full-width {
          grid-column: 1 / -1;
        }
        
        h3 {
          font-size: 14px;
          font-weight: 500;
          color: #666;
          margin: 0 0 16px;
        }
        
        .manual-upload {
          display: flex;
          align-items: center;
        }
        
        .event-log {
          margin-top: 20px;
          padding: 16px;
          background: #fff;
          border-radius: 4px;
          border: 1px solid #e8e8e8;
          
          h4 {
            font-size: 14px;
            font-weight: 500;
            color: #333;
            margin: 0 0 12px;
          }
          
          .log-list {
            max-height: 300px;
            overflow-y: auto;
            
            .log-item {
              display: flex;
              align-items: center;
              padding: 8px 0;
              font-size: 12px;
              border-bottom: 1px solid #f5f5f5;
              
              &:last-child {
                border-bottom: none;
              }
              
              .log-time {
                color: #999;
                margin-right: 12px;
                font-family: monospace;
              }
              
              .log-event {
                padding: 2px 8px;
                border-radius: 4px;
                margin-right: 12px;
                font-weight: 500;
                
                &.log-success {
                  background: #f6ffed;
                  color: #52c41a;
                }
                
                &.log-error {
                  background: #fff2f0;
                  color: #f5222d;
                }
                
                &.log-warning {
                  background: #fff7e6;
                  color: #fa8c16;
                }
                
                &.log-info {
                  background: #e6f7ff;
                  color: #1890ff;
                }
              }
              
              .log-message {
                color: #666;
                flex: 1;
              }
            }
          }
        }
      }
    }
  }
}
</style>
