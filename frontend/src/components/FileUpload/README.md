# FileUpload 文件上传组件

通用文件上传组件，支持多种上传模式、文件类型验证、拖拽上传等功能。

## 基本用法

```vue
<template>
  <FileUpload 
    accept="image/*"
    :max-size="5"
    @success="handleSuccess"
  />
</template>

<script setup>
import FileUpload from '@/components/FileUpload/index.vue'

const handleSuccess = (file) => {
  console.log('上传成功:', file)
}
</script>
```

## Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| mode | 上传模式：`card` 卡片区域 / `button` 按钮 / `image-wall` 图片墙 | String | `'card'` |
| accept | 接受的文件类型，如 `image/*`, `.pdf,.doc` | String | `'*'` |
| multiple | 是否支持多选 | Boolean | `false` |
| max | 最大文件数量，0表示不限制 | Number | `0` |
| maxSize | 单个文件最大大小(MB) | Number | `10` |
| disabled | 是否禁用 | Boolean | `false` |
| showFileList | 是否显示文件列表 | Boolean | `true` |
| uploadText | 上传区域提示文本 | String | `'点击或拖拽文件到此区域上传'` |
| uploadHint | 上传区域说明文本 | String | `'支持单个或批量上传'` |
| buttonTheme | 按钮主题 | String | `'primary'` |
| buttonVariant | 按钮样式 | String | `'base'` |
| buttonText | 按钮文本 | String | `'上传文件'` |
| buttonIcon | 按钮图标 | String | `'upload'` |
| autoUpload | 是否自动上传 | Boolean | `true` |
| hideUploadArea | 隐藏上传区域（只显示文件列表） | Boolean | `false` |
| defaultFileList | 默认文件列表 | Array | `[]` |

## Events

| 事件名 | 说明 | 回调参数 |
|--------|------|----------|
| update:fileList | 文件列表更新 | `(files: FileItem[])` |
| success | 文件上传成功 | `(file: FileItem)` |
| error | 文件上传失败 | `(file: FileItem, error: any)` |
| remove | 文件移除 | `(file: FileItem, index: number)` |
| preview | 图片预览 | `(file: FileItem)` |

## Methods

| 方法名 | 说明 | 参数 |
|--------|------|------|
| triggerUpload | 手动触发文件选择 | - |
| clearFiles | 清空文件列表 | - |
| getFileList | 获取文件列表 | - |

## 使用场景

### 1. 卡片上传模式（默认）

```vue
<FileUpload 
  accept="image/*"
  multiple
  :max="5"
  :max-size="10"
  upload-text="点击或拖拽图片到此区域"
  upload-hint="支持 jpg、png、gif 格式，单个文件不超过 10MB"
/>
```

### 2. 按钮上传模式

```vue
<FileUpload 
  mode="button"
  accept=".pdf,.doc,.docx"
  button-text="上传文档"
  button-icon="file-add"
  button-theme="primary"
  :show-file-list="true"
/>
```

### 3. 图片墙模式

```vue
<FileUpload 
  mode="image-wall"
  accept="image/*"
  multiple
  :max="9"
  @success="handleImageSuccess"
/>
```

### 4. 附件上传

```vue
<FileUpload 
  mode="button"
  accept="*"
  multiple
  :max-size="20"
  button-text="添加附件"
  button-variant="outline"
/>
```

### 5. 头像上传

```vue
<FileUpload 
  mode="image-wall"
  accept="image/*"
  :max="1"
  :max-size="2"
  :show-file-list="false"
/>
```

### 6. 只显示文件列表

```vue
<FileUpload 
  :default-file-list="existingFiles"
  :hide-upload-area="true"
  :show-file-list="true"
/>
```

## 文件类型示例

- 图片：`image/*` 或 `image/png,image/jpeg,image/gif`
- 文档：`.pdf,.doc,.docx,.xls,.xlsx`
- 视频：`video/*` 或 `.mp4,.avi,.mov`
- 音频：`audio/*` 或 `.mp3,.wav`
- 压缩包：`.zip,.rar,.7z`
- 所有文件：`*`

## 注意事项

1. 组件会自动调用后端 `/file/upload` 接口上传文件
2. 上传的文件会先进行类型和大小验证
3. 图片文件会在上传前生成本地预览
4. 支持拖拽上传（仅 card 模式）
5. 图片墙模式专为图片上传优化
