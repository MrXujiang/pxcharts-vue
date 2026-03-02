<template>
  <div class="rich-text-editor">
    <!-- 工具栏 -->
    <div v-if="editor" class="editor-toolbar">
      <!-- 文本格式 -->
      <div class="toolbar-group">
        <button 
          class="toolbar-btn"
          :class="{ 'is-active': editor.isActive('bold') }"
          @click="editor.chain().focus().toggleBold().run()"
          title="加粗"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M4 2h5a3 3 0 0 1 0 6H4V2z" stroke="currentColor" stroke-width="1.5"/>
            <path d="M4 8h6a3 3 0 0 1 0 6H4V8z" stroke="currentColor" stroke-width="1.5"/>
          </svg>
        </button>
        <button 
          class="toolbar-btn"
          :class="{ 'is-active': editor.isActive('italic') }"
          @click="editor.chain().focus().toggleItalic().run()"
          title="斜体"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M9 2h5M2 14h5M11 2L5 14" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
        <button 
          class="toolbar-btn"
          :class="{ 'is-active': editor.isActive('underline') }"
          @click="editor.chain().focus().toggleUnderline().run()"
          title="下划线"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M4 2v5a4 4 0 0 0 8 0V2M2 14h12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
        <button 
          class="toolbar-btn"
          :class="{ 'is-active': editor.isActive('strike') }"
          @click="editor.chain().focus().toggleStrike().run()"
          title="删除线"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M2 8h12M5 3l6 10M11 3L5 13" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
      </div>

      <div class="toolbar-divider"></div>

      <!-- 标题 -->
      <div class="toolbar-group">
        <button 
          class="toolbar-btn"
          :class="{ 'is-active': editor.isActive('heading', { level: 1 }) }"
          @click="editor.chain().focus().toggleHeading({ level: 1 }).run()"
          title="标题 1"
        >
          H1
        </button>
        <button 
          class="toolbar-btn"
          :class="{ 'is-active': editor.isActive('heading', { level: 2 }) }"
          @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
          title="标题 2"
        >
          H2
        </button>
        <button 
          class="toolbar-btn"
          :class="{ 'is-active': editor.isActive('heading', { level: 3 }) }"
          @click="editor.chain().focus().toggleHeading({ level: 3 }).run()"
          title="标题 3"
        >
          H3
        </button>
      </div>

      <div class="toolbar-divider"></div>

      <!-- 列表 -->
      <div class="toolbar-group">
        <button 
          class="toolbar-btn"
          :class="{ 'is-active': editor.isActive('bulletList') }"
          @click="editor.chain().focus().toggleBulletList().run()"
          title="无序列表"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <circle cx="3" cy="3" r="1.5" fill="currentColor"/>
            <circle cx="3" cy="8" r="1.5" fill="currentColor"/>
            <circle cx="3" cy="13" r="1.5" fill="currentColor"/>
            <path d="M7 3h7M7 8h7M7 13h7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
        <button 
          class="toolbar-btn"
          :class="{ 'is-active': editor.isActive('orderedList') }"
          @click="editor.chain().focus().toggleOrderedList().run()"
          title="有序列表"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M2 3h1M2 8h1M2 13h1M7 3h7M7 8h7M7 13h7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
        <button 
          class="toolbar-btn"
          :class="{ 'is-active': editor.isActive('blockquote') }"
          @click="editor.chain().focus().toggleBlockquote().run()"
          title="引用"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M3 4h10v8H3V4z" stroke="currentColor" stroke-width="1.5"/>
            <path d="M1 6v4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </button>
      </div>

      <div class="toolbar-divider"></div>

      <!-- 对齐 -->
      <div class="toolbar-group">
        <button 
          class="toolbar-btn"
          :class="{ 'is-active': editor.isActive({ textAlign: 'left' }) }"
          @click="editor.chain().focus().setTextAlign('left').run()"
          title="左对齐"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M2 3h12M2 6h8M2 9h12M2 12h8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
        <button 
          class="toolbar-btn"
          :class="{ 'is-active': editor.isActive({ textAlign: 'center' }) }"
          @click="editor.chain().focus().setTextAlign('center').run()"
          title="居中"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M2 3h12M4 6h8M2 9h12M4 12h8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
        <button 
          class="toolbar-btn"
          :class="{ 'is-active': editor.isActive({ textAlign: 'right' }) }"
          @click="editor.chain().focus().setTextAlign('right').run()"
          title="右对齐"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M2 3h12M6 6h8M2 9h12M6 12h8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
      </div>

      <div class="toolbar-divider"></div>

      <!-- 插入 -->
      <div class="toolbar-group">
        <button 
          class="toolbar-btn"
          @click="handleImageUpload"
          title="插入图片"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <rect x="2" y="2" width="12" height="12" rx="2" stroke="currentColor" stroke-width="1.5"/>
            <circle cx="5.5" cy="5.5" r="1.5" fill="currentColor"/>
            <path d="M2 11l3-3 2 2 4-4 3 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
        <button 
          class="toolbar-btn"
          @click="handleLinkAdd"
          title="插入链接"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M6.5 9.5l3-3M8 6l1.5-1.5a2.5 2.5 0 0 1 3.5 3.5L11.5 9.5M4.5 6.5L3 8a2.5 2.5 0 0 0 3.5 3.5L8 10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
        <button 
          class="toolbar-btn"
          @click="handleFileUpload"
          title="插入文件"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M9 2H4a1 1 0 0 0-1 1v10a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V6L9 2z" stroke="currentColor" stroke-width="1.5"/>
            <path d="M9 2v4h4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
      </div>

      <div class="toolbar-divider"></div>

      <!-- 其他操作 -->
      <div class="toolbar-group">
        <button 
          class="toolbar-btn"
          @click="editor.chain().focus().undo().run()"
          :disabled="!editor.can().undo()"
          title="撤销"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M4 8h8a3 3 0 0 1 0 6H8M4 8l3-3M4 8l3 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
        <button 
          class="toolbar-btn"
          @click="editor.chain().focus().redo().run()"
          :disabled="!editor.can().redo()"
          title="重做"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M12 8H4a3 3 0 0 0 0 6h4M12 8l-3-3M12 8l-3 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
        <button 
          class="toolbar-btn"
          @click="editor.chain().focus().clearNodes().unsetAllMarks().run()"
          title="清除格式"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M3 3l10 10M13 3L3 13" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- 编辑器内容区 -->
    <editor-content :editor="editor" class="editor-content" />

    <!-- 隐藏的文件上传输入框 -->
    <input 
      ref="imageInputRef"
      type="file" 
      accept="image/*" 
      style="display: none"
      @change="handleImageChange"
    />
    <input 
      ref="fileInputRef"
      type="file" 
      style="display: none"
      @change="handleFileChange"
    />

    <!-- 链接弹窗 -->
    <t-dialog
      v-model:visible="showLinkDialog"
      header="插入链接"
      width="480px"
      @confirm="confirmAddLink"
    >
      <div class="link-form">
        <div class="form-item">
          <label>链接文本：</label>
          <t-input v-model="linkForm.text" placeholder="请输入链接文本" />
        </div>
        <div class="form-item">
          <label>链接地址：</label>
          <t-input v-model="linkForm.url" placeholder="https://" />
        </div>
      </div>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onBeforeUnmount } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import TextAlign from '@tiptap/extension-text-align'
import Underline from '@tiptap/extension-underline'
import { MessagePlugin } from 'tdesign-vue-next'

interface Props {
  modelValue?: string
  placeholder?: string
  editable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: '请输入内容...',
  editable: true,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'change': [value: string]
}>()

// 编辑器实例
const editor = useEditor({
  content: props.modelValue,
  editable: props.editable,
  extensions: [
    StarterKit,
    Underline,
    Image.configure({
      HTMLAttributes: {
        class: 'editor-image',
      },
    }),
    Link.configure({
      openOnClick: false,
      HTMLAttributes: {
        class: 'editor-link',
      },
    }),
    Placeholder.configure({
      placeholder: props.placeholder,
    }),
    TextAlign.configure({
      types: ['heading', 'paragraph'],
    }),
  ],
  onUpdate: ({ editor }) => {
    const html = editor.getHTML()
    emit('update:modelValue', html)
    emit('change', html)
  },
})

// 文件上传引用
const imageInputRef = ref<HTMLInputElement>()
const fileInputRef = ref<HTMLInputElement>()

// 链接表单
const showLinkDialog = ref(false)
const linkForm = ref({
  text: '',
  url: '',
})

// 监听外部值变化
watch(() => props.modelValue, (value) => {
  if (editor.value && value !== editor.value.getHTML()) {
    editor.value.commands.setContent(value)
  }
})

// 监听可编辑状态
watch(() => props.editable, (value) => {
  if (editor.value) {
    editor.value.setEditable(value)
  }
})

// 图片上传
const handleImageUpload = () => {
  imageInputRef.value?.click()
}

const handleImageChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  // 验证文件类型
  if (!file.type.startsWith('image/')) {
    MessagePlugin.warning('请选择图片文件')
    return
  }

  // 验证文件大小（10MB）
  if (file.size > 10 * 1024 * 1024) {
    MessagePlugin.warning('图片大小不能超过10MB')
    return
  }

  // 读取并插入图片
  const reader = new FileReader()
  reader.onload = (e) => {
    const url = e.target?.result as string
    editor.value?.chain().focus().setImage({ src: url }).run()
    MessagePlugin.success('图片已插入')
  }
  reader.readAsDataURL(file)

  // 清空input
  target.value = ''
}

// 文件上传
const handleFileUpload = () => {
  fileInputRef.value?.click()
}

const handleFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  // 验证文件大小（50MB）
  if (file.size > 50 * 1024 * 1024) {
    MessagePlugin.warning('文件大小不能超过50MB')
    return
  }

  // 插入文件链接
  const fileName = file.name
  editor.value
    ?.chain()
    .focus()
    .insertContent(`<p><a href="#" class="file-link" data-filename="${fileName}">📎 ${fileName}</a></p>`)
    .run()
  
  MessagePlugin.success('文件已插入')

  // 清空input
  target.value = ''
}

// 链接操作
const handleLinkAdd = () => {
  const previousUrl = editor.value?.getAttributes('link').href || ''
  linkForm.value.url = previousUrl
  linkForm.value.text = editor.value?.state.doc.textBetween(
    editor.value.state.selection.from,
    editor.value.state.selection.to
  ) || ''
  showLinkDialog.value = true
}

const confirmAddLink = () => {
  if (!linkForm.value.url) {
    MessagePlugin.warning('请输入链接地址')
    return
  }

  if (linkForm.value.text) {
    editor.value
      ?.chain()
      .focus()
      .insertContent(`<a href="${linkForm.value.url}">${linkForm.value.text}</a>`)
      .run()
  } else {
    editor.value
      ?.chain()
      .focus()
      .setLink({ href: linkForm.value.url })
      .run()
  }

  showLinkDialog.value = false
  linkForm.value = { text: '', url: '' }
}

// 清理
onBeforeUnmount(() => {
  editor.value?.destroy()
})
</script>

<style scoped lang="less">
.rich-text-editor {
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  background: #ffffff;
  overflow: hidden;

  .editor-toolbar {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 8px 12px;
    border-bottom: 1px solid #e8e8e8;
    background: #fafafa;
    flex-wrap: wrap;

    .toolbar-group {
      display: flex;
      align-items: center;
      gap: 4px;
    }

    .toolbar-divider {
      width: 1px;
      height: 20px;
      background: #e8e8e8;
      margin: 0 4px;
    }

    .toolbar-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 32px;
      height: 32px;
      border: none;
      background: transparent;
      border-radius: 4px;
      cursor: pointer;
      color: #666666;
      font-size: 13px;
      font-weight: 500;
      transition: all 0.2s;

      &:hover:not(:disabled) {
        background: #e8e8e8;
        color: #333333;
      }

      &.is-active {
        background: #e6f0ff;
        color: #4a7ff7;
      }

      &:disabled {
        opacity: 0.4;
        cursor: not-allowed;
      }

      svg {
        flex-shrink: 0;
      }
    }
  }

  .editor-content {
    min-height: 300px;
    max-height: 600px;
    overflow-y: auto;
    padding: 16px;

    :deep(.ProseMirror) {
      outline: none;
      min-height: 268px;

      > * + * {
        margin-top: 0.75em;
      }

      h1 {
        font-size: 28px;
        font-weight: 600;
        line-height: 1.4;
        color: #1a1a1a;
      }

      h2 {
        font-size: 24px;
        font-weight: 600;
        line-height: 1.4;
        color: #1a1a1a;
      }

      h3 {
        font-size: 20px;
        font-weight: 600;
        line-height: 1.4;
        color: #1a1a1a;
      }

      p {
        font-size: 15px;
        line-height: 1.8;
        color: #333333;
      }

      strong {
        font-weight: 600;
      }

      em {
        font-style: italic;
      }

      u {
        text-decoration: underline;
      }

      s {
        text-decoration: line-through;
      }

      a {
        color: #4a7ff7;
        text-decoration: underline;
        cursor: pointer;

        &:hover {
          color: #3a6fe6;
        }

        &.file-link {
          display: inline-flex;
          align-items: center;
          gap: 4px;
          padding: 4px 12px;
          background: #f5f5f5;
          border-radius: 4px;
          text-decoration: none;
          color: #333333;
          font-size: 14px;

          &:hover {
            background: #e8e8e8;
          }
        }
      }

      ul, ol {
        padding-left: 1.5em;

        li {
          margin: 0.25em 0;
          line-height: 1.8;
        }
      }

      ul {
        list-style-type: disc;
      }

      ol {
        list-style-type: decimal;
      }

      blockquote {
        padding-left: 1em;
        border-left: 3px solid #4a7ff7;
        color: #666666;
        font-style: italic;
        margin: 1em 0;
      }

      code {
        background: #f5f5f5;
        padding: 2px 6px;
        border-radius: 3px;
        font-family: 'Courier New', monospace;
        font-size: 0.9em;
      }

      pre {
        background: #1a1a1a;
        color: #ffffff;
        padding: 1em;
        border-radius: 6px;
        overflow-x: auto;

        code {
          background: transparent;
          padding: 0;
          color: inherit;
        }
      }

      img {
        max-width: 100%;
        height: auto;
        border-radius: 6px;
        display: block;
        margin: 1em 0;

        &.ProseMirror-selectednode {
          outline: 2px solid #4a7ff7;
        }
      }

      p.is-editor-empty:first-child::before {
        content: attr(data-placeholder);
        float: left;
        color: #999999;
        pointer-events: none;
        height: 0;
      }
    }

    &::-webkit-scrollbar {
      width: 6px;
    }

    &::-webkit-scrollbar-thumb {
      background: #d9d9d9;
      border-radius: 3px;

      &:hover {
        background: #bfbfbf;
      }
    }
  }
}

.link-form {
  .form-item {
    margin-bottom: 16px;

    &:last-child {
      margin-bottom: 0;
    }

    label {
      display: block;
      margin-bottom: 8px;
      font-size: 14px;
      color: #333333;
      font-weight: 500;
    }
  }
}
</style>
