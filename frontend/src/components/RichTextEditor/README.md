# RichTextEditor 富文本编辑器

基于 TipTap 的 Vue 3 富文本编辑器组件，提供简约美观的界面和丰富的编辑功能。

## 安装依赖

```bash
npm install @tiptap/vue-3 @tiptap/starter-kit @tiptap/extension-image @tiptap/extension-link @tiptap/extension-placeholder @tiptap/extension-text-align @tiptap/extension-underline
```

## 基本用法

```vue
<template>
  <RichTextEditor v-model="content" placeholder="请输入内容..." />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import RichTextEditor from '@/components/RichTextEditor/index.vue'

const content = ref('')
</script>
```

## Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| modelValue | 编辑器内容（HTML） | string | '' |
| placeholder | 占位符文本 | string | '请输入内容...' |
| editable | 是否可编辑 | boolean | true |

## Events

| 事件名 | 说明 | 回调参数 |
|--------|------|----------|
| update:modelValue | 内容更新时触发 | (value: string) |
| change | 内容变化时触发 | (value: string) |

## 功能特性

### 文本格式
- ✅ 加粗、斜体、下划线、删除线
- ✅ 标题（H1/H2/H3）
- ✅ 清除格式

### 段落格式
- ✅ 无序列表
- ✅ 有序列表
- ✅ 引用块
- ✅ 左对齐/居中/右对齐

### 插入内容
- ✅ 图片上传（支持 10MB 以内）
- ✅ 链接插入
- ✅ 文件上传（支持 50MB 以内）

### 编辑操作
- ✅ 撤销/重做
- ✅ 快捷键支持

## 使用示例

### 只读模式

```vue
<RichTextEditor 
  v-model="content" 
  :editable="false"
/>
```

### 监听内容变化

```vue
<template>
  <RichTextEditor 
    v-model="content"
    @change="handleChange"
  />
</template>

<script setup lang="ts">
const handleChange = (html: string) => {
  console.log('内容已更新:', html)
}
</script>
```

### 自定义占位符

```vue
<RichTextEditor 
  v-model="content"
  placeholder="请输入文章内容，支持图片和文件上传..."
/>
```

### 表单集成

```vue
<template>
  <div class="form">
    <div class="form-item">
      <label>文章标题</label>
      <t-input v-model="form.title" />
    </div>
    <div class="form-item">
      <label>文章内容</label>
      <RichTextEditor v-model="form.content" />
    </div>
    <t-button @click="submit">发布文章</t-button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const form = ref({
  title: '',
  content: ''
})

const submit = () => {
  console.log('提交内容:', form.value)
}
</script>
```

## 样式定制

组件使用 scoped 样式，可通过 CSS 变量或 :deep() 选择器自定义样式：

```vue
<style>
/* 自定义编辑器高度 */
.rich-text-editor .editor-content {
  min-height: 400px;
  max-height: 800px;
}

/* 自定义工具栏背景 */
.rich-text-editor .editor-toolbar {
  background: #f0f0f0;
}
</style>
```

## 注意事项

1. 图片会以 Base64 格式存储，建议实际使用时改为上传到服务器
2. 文件上传仅插入文件名，需要配合后端实现真实文件上传
3. 编辑器内容以 HTML 格式存储，确保后端支持 HTML 内容存储
4. 建议在生产环境中对用户输入的 HTML 进行安全过滤
