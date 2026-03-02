<template>
  <div class="file-upload-wrapper">
    <!-- 上传区域 -->
    <div
      v-if="mode === 'card' && !hideUploadArea"
      class="upload-area"
      :class="{
        'is-dragover': isDragOver,
        'is-disabled': disabled || (max > 0 && fileList.length >= max),
      }"
      @click="triggerUpload"
      @dragover.prevent="handleDragOver"
      @dragleave.prevent="handleDragLeave"
      @drop.prevent="handleDrop"
    >
      <input
        ref="fileInputRef"
        type="file"
        :accept="accept"
        :multiple="multiple"
        :disabled="disabled"
        class="file-input"
        @change="handleFileChange"
      />

      <div class="upload-content">
        <t-icon name="cloud-upload" size="48px" class="upload-icon" />
        <p class="upload-text">{{ uploadText }}</p>
        <p class="upload-hint">{{ uploadHint }}</p>
      </div>
    </div>

    <!-- 按钮模式 -->
    <div v-if="mode === 'button'" class="upload-button-wrapper">
      <input
        ref="fileInputRef"
        type="file"
        :accept="accept"
        :multiple="multiple"
        :disabled="disabled"
        class="file-input"
        @change="handleFileChange"
      />
      <t-button
        :theme="buttonTheme"
        :variant="buttonVariant"
        :disabled="disabled || (max > 0 && fileList.length >= max)"
        @click="triggerUpload"
      >
        <template #icon>
          <t-icon :name="buttonIcon" />
        </template>
        {{ buttonText }}
      </t-button>
    </div>

    <!-- 图片墙模式 -->
    <div v-if="mode === 'image-wall'" class="image-wall">
      <div v-for="(file, index) in fileList" :key="file.id" class="image-wall-item">
        <img :src="file.url" :alt="file.name" class="image-preview" @click="handlePreview(file)" />
        <div class="image-mask">
          <div class="image-actions">
            <div class="action-icon-wrapper" @click.stop="handlePreview(file)">
              <t-icon name="zoomIn" size="20px" class="action-icon" />
            </div>
            <div class="action-icon-wrapper delete-wrapper" @click.stop="handleRemove(index)">
              <t-icon name="delete" size="20px" class="action-icon delete-icon" />
            </div>
          </div>
        </div>
        <div v-if="file.uploading" class="upload-progress">
          <t-loading size="small" />
          <span class="upload-text">上传中...</span>
        </div>
      </div>

      <div v-if="!max || fileList.length < max" class="image-wall-upload" @click="triggerUpload">
        <input
          ref="fileInputRef"
          type="file"
          :accept="accept"
          :multiple="multiple"
          :disabled="disabled"
          class="file-input"
          @change="handleFileChange"
        />
        <t-icon name="add" size="24px" />
        <span class="upload-label">上传图片</span>
      </div>
    </div>

    <!-- 链接模式 -->
    <div v-if="mode === 'link'" class="link-mode">
      <!-- 上传按钮 -->
      <div v-if="!hideUploadArea" class="link-upload-wrapper">
        <input
          ref="fileInputRef"
          type="file"
          :accept="accept"
          :multiple="multiple"
          :disabled="disabled"
          class="file-input"
          @change="handleFileChange"
        />
        <t-button
          size="small"
          variant="text"
          :disabled="disabled || (max > 0 && fileList.length >= max)"
          @click="triggerUpload"
        >
          <template #icon>
            <t-icon name="add" />
          </template>
          添加附件
        </t-button>
      </div>

      <!-- 文件链接列表 -->
      <div v-if="fileList.length > 0" class="link-list">
        <div
          v-for="(file, index) in fileList"
          :key="file.id"
          class="link-item"
          :class="{ 'is-uploading': file.uploading, 'is-error': file.error }"
        >
          <t-icon :name="getFileIcon(file)" size="16px" class="link-icon" />
          <a
            v-if="!file.uploading && !file.error"
            :href="file.url"
            target="_blank"
            class="link-name"
            :title="file.name"
          >
            {{ file.name }}
          </a>
          <span v-else class="link-name" :title="file.name">
            {{ file.name }}
          </span>
          <span class="link-size">{{ formatFileSize(file.size) }}</span>
          <span v-if="file.uploading" class="link-status uploading">
            <t-loading size="12px" />
            上传中...
          </span>
          <span v-else-if="file.error" class="link-status error">失败</span>
          <t-icon
            v-if="!file.uploading"
            name="close"
            size="14px"
            class="link-remove"
            @click="handleRemove(index)"
          />
        </div>
      </div>
    </div>

    <!-- 文件列表 -->
    <div
      v-if="showFileList && fileList.length > 0 && mode !== 'image-wall' && mode !== 'link'"
      class="file-list"
    >
      <div
        v-for="(file, index) in fileList"
        :key="file.id"
        class="file-item"
        :class="{ 'is-uploading': file.uploading, 'is-error': file.error }"
      >
        <!-- 文件图标 -->
        <div class="file-icon" :class="{ 'is-image': isImage(file) }">
          <img
            v-if="isImage(file)"
            :src="file.url"
            :alt="file.name"
            class="file-thumbnail"
            @click="handlePreview(file)"
          />
          <t-icon v-else :name="getFileIcon(file)" size="32px" />
        </div>

        <!-- 文件信息 -->
        <div class="file-info">
          <div class="file-name" :title="file.name">{{ file.name }}</div>
          <div class="file-meta">
            <span class="file-size">{{ formatFileSize(file.size) }}</span>
            <span v-if="file.uploading" class="file-status">上传中...</span>
            <span v-else-if="file.error" class="file-status error">上传失败</span>
            <span v-else class="file-status success">上传成功</span>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="file-actions">
          <t-icon
            v-if="isImage(file) && !file.uploading"
            name="zoom-in"
            size="18px"
            class="action-btn preview"
            @click="handlePreview(file)"
          />
          <t-icon
            v-if="!file.uploading"
            name="download"
            size="18px"
            class="action-btn"
            @click="handleDownload(file)"
          />
          <t-icon
            v-if="!file.uploading"
            name="delete"
            size="18px"
            class="action-btn delete"
            @click="handleRemove(index)"
          />
        </div>
      </div>
    </div>

    <!-- 图片预览对话框 -->
    <t-dialog 
      v-model:visible="previewVisible" 
      header="图片预览" 
      width="80%"
      :footer="false" 
      attach="body"
      :z-index="9999"
    >
      <div class="image-preview-wrapper">
        <img :src="previewUrl" alt="预览图片" class="preview-image" />
        <div class="preview-actions">
          <t-button variant="outline" size="small" @click="handleDownloadPreview">
            <template #icon>
              <t-icon name="download" />
            </template>
            下载
          </t-button>
        </div>
      </div>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { MessagePlugin } from 'tdesign-vue-next'
import { uploadFile } from '@/api'
import { isImage, getFileIcon, FileItem } from '@/utils'

interface Props {
  // 上传模式: card-卡片区域, button-按钮, image-wall-图片墙, link-链接列表
  mode?: 'card' | 'button' | 'image-wall' | 'link'
  // 接受的文件类型
  accept?: string
  // 是否多选
  multiple?: boolean
  // 最大文件数量
  max?: number
  // 单个文件最大大小(MB)
  maxSize?: number
  // 是否禁用
  disabled?: boolean
  // 是否显示文件列表
  showFileList?: boolean
  // 上传提示文本
  uploadText?: string
  // 上传提示说明
  uploadHint?: string
  // 按钮主题
  buttonTheme?: 'default' | 'primary' | 'success' | 'warning' | 'danger'
  // 按钮样式
  buttonVariant?: 'base' | 'outline' | 'text'
  // 按钮文本
  buttonText?: string
  // 按钮图标
  buttonIcon?: string
  // 是否自动上传
  autoUpload?: boolean
  // 隐藏上传区域(用于只显示文件列表)
  hideUploadArea?: boolean
  // 默认文件列表
  defaultFileList?: FileItem[]
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'card',
  accept: '*',
  multiple: false,
  max: 0,
  maxSize: 10,
  disabled: false,
  showFileList: true,
  uploadText: '点击或拖拽文件到此区域上传',
  uploadHint: '支持单个或批量上传',
  buttonTheme: 'primary',
  buttonVariant: 'base',
  buttonText: '上传文件',
  buttonIcon: 'upload',
  autoUpload: true,
  hideUploadArea: false,
  defaultFileList: () => [],
})

interface Emits {
  (e: 'update:fileList', files: FileItem[]): void
  (e: 'success', file: FileItem): void
  (e: 'error', file: FileItem, error: any): void
  (e: 'remove', file: FileItem, index: number): void
  (e: 'preview', file: FileItem): void
}

const emit = defineEmits<Emits>()

const fileInputRef = ref<HTMLInputElement>()
const fileList = ref<FileItem[]>([...props.defaultFileList])
const isDragOver = ref(false)
const previewVisible = ref(false)
const previewUrl = ref('')

// 监听 defaultFileList 变化
watch(
  () => props.defaultFileList,
  (newVal) => {
    fileList.value = [...newVal]
  },
  { deep: true },
)

// 触发文件选择
const triggerUpload = () => {
  if (props.disabled || (props.max > 0 && fileList.value.length >= props.max)) {
    return
  }
  fileInputRef.value?.click()
}

// 处理文件选择
const handleFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  const files = target.files
  if (files && files.length > 0) {
    handleFiles(Array.from(files))
  }
  // 清空 input value，允许重复选择同一文件
  target.value = ''
}

// 处理拖拽上传
const handleDragOver = (event: DragEvent) => {
  if (!props.disabled) {
    isDragOver.value = true
  }
}

const handleDragLeave = () => {
  isDragOver.value = false
}

const handleDrop = (event: DragEvent) => {
  isDragOver.value = false
  if (props.disabled) return

  const files = event.dataTransfer?.files
  if (files && files.length > 0) {
    handleFiles(Array.from(files))
  }
}

// 处理文件
const handleFiles = async (files: File[]) => {
  // 检查文件数量限制
  if (props.max > 0 && fileList.value.length + files.length > props.max) {
    MessagePlugin.warning(`最多只能上传 ${props.max} 个文件`)
    return
  }

  for (const file of files) {
    // 验证文件类型
    if (props.accept !== '*' && !validateFileType(file)) {
      MessagePlugin.warning(`文件 ${file.name} 类型不符合要求`)
      continue
    }

    // 验证文件大小
    if (file.size > props.maxSize * 1024 * 1024) {
      MessagePlugin.warning(`文件 ${file.name} 大小超过 ${props.maxSize}MB`)
      continue
    }

    const fileItem: FileItem = {
      id: `${Date.now()}-${Math.random()}`,
      name: file.name,
      size: file.size,
      type: file.type,
      url: '',
      uploading: true,
      error: false,
    }

    // 如果是图片，先生成本地预览
    if (file.type.startsWith('image/')) {
      fileItem.url = URL.createObjectURL(file)
    }

    // 先添加到列表
    fileList.value.push(fileItem)
    emit('update:fileList', fileList.value)

    // 自动上传
    if (props.autoUpload) {
      // 使用 nextTick 确保响应式更新完成后再上传
      await uploadFileToServer(file, fileItem)
    }
  }
}

// 上传文件到服务器
const uploadFileToServer = async (file: File, fileItem: FileItem) => {
  try {
    const response = (await uploadFile(file)) as any

    console.log('文件上传响应:', response)

    // 兼容多种返回格式
    let fileUrl = ''

    // 格式1: 直接返回 { url: '...' }
    if (response && typeof response === 'object' && response.url) {
      fileUrl = response.url
    }
    // 格式2: 返回 { data: { url: '...' } }
    else if (response && response.data && response.data.url) {
      fileUrl = response.data.url
    }
    // 格式3: 返回 { state: 200, data: { url: '...' } }
    else if (response && response.state === 200 && response.data && response.data.url) {
      fileUrl = response.data.url
    }
    // 格式4: 直接返回 url 字符串
    else if (typeof response === 'string') {
      fileUrl = response
    }

    if (fileUrl) {
      // 关键修复: 从数组中查找对象引用并更新
      const targetFile = fileList.value.find((f) => f.id === fileItem.id)
      if (targetFile) {
        targetFile.url = fileUrl
        targetFile.uploading = false
        targetFile.error = false
        console.log('文件上传成功，URL:', fileUrl)
        emit('success', targetFile)
        // 移除组件内部的 Message 提示，由外部监听 success 事件处理
      }
    } else {
      // 如果没有找到 url 字段
      console.error('上传响应格式异常:', response)
      const targetFile = fileList.value.find((f) => f.id === fileItem.id)
      if (targetFile) {
        targetFile.uploading = false
        targetFile.error = true
      }
      MessagePlugin.error(`${file.name} 上传失败：返回数据格式错误`)
    }
  } catch (error: any) {
    console.error('文件上传失败:', error)
    const targetFile = fileList.value.find((f) => f.id === fileItem.id)
    if (targetFile) {
      targetFile.uploading = false
      targetFile.error = true
    }
    emit('error', fileItem, error)
    MessagePlugin.error(error.message || `${file.name} 上传失败`)
  }
}

// 验证文件类型
const validateFileType = (file: File): boolean => {
  const acceptTypes = props.accept.split(',').map((type) => type.trim())

  for (const acceptType of acceptTypes) {
    if (acceptType.startsWith('.')) {
      // 扩展名匹配
      if (file.name.toLowerCase().endsWith(acceptType.toLowerCase())) {
        return true
      }
    } else if (acceptType.includes('*')) {
      // MIME类型通配符匹配
      const regex = new RegExp(acceptType.replace('*', '.*'))
      if (regex.test(file.type)) {
        return true
      }
    } else {
      // 精确MIME类型匹配
      if (file.type === acceptType) {
        return true
      }
    }
  }

  return false
}

// 移除文件
const handleRemove = (index: number) => {
  console.log('删除文件，索引:', index, '文件列表长度:', fileList.value.length)
  const file = fileList.value[index]
  if (!file) {
    console.error('文件不存在，索引:', index)
    return
  }
  console.log('删除文件:', file)
  fileList.value.splice(index, 1)
  emit('update:fileList', fileList.value)
  emit('remove', file, index)
  MessagePlugin.success('文件已移除')
}

// 下载文件
const handleDownload = (file: FileItem) => {
  const link = document.createElement('a')
  link.href = file.url
  link.download = file.name
  link.click()
}

// 预览图片
const handlePreview = (file: FileItem) => {
  console.log('预览图片:', file)
  if (!file.url || file.uploading) {
    MessagePlugin.warning('文件正在上传中，请稍候')
    return
  }
  previewUrl.value = file.url
  previewVisible.value = true
  emit('preview', file)
}

// 下载预览中的图片
const handleDownloadPreview = () => {
  if (previewUrl.value) {
    const link = document.createElement('a')
    link.href = previewUrl.value
    link.download = '图片' // 可以根据需要改进文件名
    link.click()
  }
}

// 格式化文件大小
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i]
}

// 暴露方法供外部调用
defineExpose({
  triggerUpload,
  clearFiles: () => {
    fileList.value = []
    emit('update:fileList', fileList.value)
  },
  getFileList: () => fileList.value,
})
</script>

<style scoped lang="less">
.file-upload-wrapper {
  width: 100%;

  .file-input {
    display: none;
  }

  // 卡片上传区域
  .upload-area {
    border: 2px dashed #dcdcdc;
    border-radius: 8px;
    padding: 40px 20px;
    text-align: center;
    cursor: pointer;
    transition: all 0.3s;
    background: #fafafa;

    &:hover {
      border-color: #4a7ff7;
      background: #f0f5ff;
    }

    &.is-dragover {
      border-color: #4a7ff7;
      background: #e6f0ff;
    }

    &.is-disabled {
      cursor: not-allowed;
      opacity: 0.6;

      &:hover {
        border-color: #dcdcdc;
        background: #fafafa;
      }
    }

    .upload-content {
      .upload-icon {
        color: #4a7ff7;
        margin-bottom: 12px;
      }

      .upload-text {
        font-size: 14px;
        color: #333;
        margin: 0 0 8px;
      }

      .upload-hint {
        font-size: 12px;
        color: #999;
        margin: 0;
      }
    }
  }

  // 按钮模式
  .upload-button-wrapper {
    display: inline-block;
  }

  // 图片墙
  .image-wall {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(104px, 1fr));
    gap: 8px;

    .image-wall-item {
      position: relative;
      width: 104px;
      height: 104px;
      border: 1px solid #dcdcdc;
      border-radius: 4px;
      overflow: hidden;
      cursor: pointer;
      transition: all 0.3s;

      &:hover {
        border-color: #4a7ff7;
        box-shadow: 0 2px 8px rgba(74, 127, 247, 0.2);
        transform: translateY(-2px);
      }

      .image-preview {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.3s;
      }

      &:hover .image-preview {
        transform: scale(1.05);
      }

      .image-mask {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.6);
        display: flex;
        align-items: center;
        justify-content: center;
        opacity: 0;
        visibility: hidden;
        transition: all 0.3s;
        z-index: 2;

        .image-actions {
          display: flex;
          gap: 12px;

          .action-icon-wrapper {
            cursor: pointer;
            transition: all 0.2s;
            padding: 6px;
            border-radius: 4px;
            background: rgba(255, 255, 255, 0.2);
            display: flex;
            align-items: center;
            justify-content: center;

            &:hover {
              transform: scale(1.1);
              background: rgba(255, 255, 255, 0.3);
            }

            &:active {
              transform: scale(1);
            }

            &.delete-wrapper:hover {
              background: rgba(255, 77, 79, 0.25);

              .action-icon {
                color: #ff4d4f;
              }
            }

            .action-icon {
              color: #fff;
              transition: color 0.2s;
            }
          }
        }
      }

      &:hover .image-mask {
        opacity: 1;
        visibility: visible;
      }

      .upload-progress {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(255, 255, 255, 0.9);
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 8px;

        .upload-text {
          font-size: 12px;
          color: #666;
        }
      }
    }

    .image-wall-upload {
      width: 104px;
      height: 104px;
      border: 1px dashed #dcdcdc;
      border-radius: 4px;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      transition: all 0.3s;
      background: #fafafa;

      &:hover {
        border-color: #4a7ff7;
        background: #f0f5ff;
        color: #4a7ff7;
      }

      .upload-label {
        font-size: 12px;
        margin-top: 4px;
      }
    }
  }

  // 文件列表
  .file-list {
    margin-top: 16px;

    .file-item {
      display: flex;
      align-items: center;
      padding: 12px;
      border: 1px solid #e8e8e8;
      border-radius: 4px;
      margin-bottom: 8px;
      transition: all 0.3s;

      &:hover {
        background: #fafafa;
      }

      &.is-uploading {
        opacity: 0.6;
      }

      &.is-error {
        border-color: #f5222d;
        background: #fff2f0;
      }

      .file-icon {
        width: 40px;
        height: 40px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 12px;
        flex-shrink: 0;

        &.is-image {
          cursor: pointer;
          transition: transform 0.2s;

          &:hover {
            transform: scale(1.05);
          }
        }

        .file-thumbnail {
          width: 40px;
          height: 40px;
          object-fit: cover;
          border-radius: 4px;
        }
      }

      .file-info {
        flex: 1;
        min-width: 0;

        .file-name {
          font-size: 14px;
          color: #333;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          margin-bottom: 4px;
        }

        .file-meta {
          font-size: 12px;
          color: #999;
          display: flex;
          gap: 12px;

          .file-status {
            &.success {
              color: #52c41a;
            }

            &.error {
              color: #f5222d;
            }
          }
        }
      }

      .file-actions {
        display: flex;
        gap: 8px;

        .action-btn {
          cursor: pointer;
          color: #666;
          transition: color 0.3s;

          &:hover {
            color: #4a7ff7;
          }

          &.preview:hover {
            color: #52c41a;
          }

          &.delete:hover {
            color: #f5222d;
          }
        }
      }
    }
  }

  // 链接模式
  .link-mode {
    .link-upload-wrapper {
      margin-bottom: 8px;
    }

    .link-list {
      .link-item {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 6px 8px;
        border-radius: 4px;
        transition: all 0.2s;

        &:hover {
          background: #f5f5f5;
        }

        &.is-uploading {
          opacity: 0.7;
        }

        &.is-error {
          .link-name {
            color: #f5222d;
            text-decoration: line-through;
          }
        }

        .link-icon {
          color: #666;
          flex-shrink: 0;
        }

        .link-name {
          flex: 1;
          font-size: 14px;
          color: #4a7ff7;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          text-decoration: none;

          &:hover {
            text-decoration: underline;
          }
        }

        .link-size {
          font-size: 12px;
          color: #999;
          flex-shrink: 0;
        }

        .link-status {
          font-size: 12px;
          flex-shrink: 0;
          display: flex;
          align-items: center;
          gap: 4px;

          &.uploading {
            color: #1890ff;
          }

          &.error {
            color: #f5222d;
          }
        }

        .link-remove {
          color: #999;
          cursor: pointer;
          flex-shrink: 0;
          transition: color 0.2s;

          &:hover {
            color: #f5222d;
          }
        }
      }
    }
  }
}
</style>

<style lang="less">
// 图片预览弹窗全局样式（非 scoped，用于 attach="body" 的 dialog）
.image-preview-wrapper {
  text-align: center;
  width: 100%;
  max-width: 100%;
  overflow: hidden;
  
  .preview-image {
    max-width: 100%;
    max-height: 70vh;
    width: 100%;
    height: auto;
    border-radius: 4px;
    object-fit: contain;
    display: block;
    margin: 0 auto;
  }
  
  .preview-actions {
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid #e8e8e8;
    display: flex;
    justify-content: center;
    gap: 12px;
  }
}
</style>
