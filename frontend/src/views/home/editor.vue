<template>
  <div class="editor-demo-page">
    <div class="demo-header">
      <h1>富文本编辑器演示</h1>
      <p>基于 TipTap 的 Vue 3 富文本编辑器组件</p>
    </div>

    <div class="demo-sections">
      <!-- 基础使用 -->
      <div class="demo-section">
        <h2>基础使用</h2>
        <div class="demo-content">
          <RichTextEditor 
            v-model="basicContent"
            placeholder="请输入内容，支持图片、链接、文件上传等功能..."
            @change="handleBasicChange"
          />
        </div>
        <div class="demo-info">
          <div class="info-label">字符数：{{ basicContent.length }}</div>
          <t-button size="small" @click="showBasicHTML = !showBasicHTML">
            {{ showBasicHTML ? '隐藏' : '查看' }} HTML
          </t-button>
        </div>
        <div v-if="showBasicHTML" class="demo-code">
          <pre>{{ basicContent || '暂无内容' }}</pre>
        </div>
      </div>

      <!-- 预设内容 -->
      <div class="demo-section">
        <h2>预设内容示例</h2>
        <div class="demo-content">
          <RichTextEditor 
            v-model="presetContent"
            placeholder="这里展示预设内容的效果"
          />
        </div>
        <div class="demo-info">
          <t-button size="small" @click="resetPresetContent">
            重置预设内容
          </t-button>
        </div>
      </div>

      <!-- 只读模式 -->
      <div class="demo-section">
        <h2>只读模式</h2>
        <div class="demo-content">
          <RichTextEditor 
            v-model="readonlyContent"
            :editable="false"
          />
        </div>
        <div class="demo-info">
          <span class="readonly-tag">只读模式</span>
        </div>
      </div>

      <!-- 表单集成 -->
      <div class="demo-section">
        <h2>表单集成示例</h2>
        <div class="demo-form">
          <div class="form-item">
            <label class="form-label">文章标题 <span class="required">*</span></label>
            <t-input 
              v-model="articleForm.title"
              placeholder="请输入文章标题"
              clearable
            />
          </div>
          <div class="form-item">
            <label class="form-label">文章分类</label>
            <t-select 
              v-model="articleForm.category"
              placeholder="请选择分类"
              clearable
            >
              <t-option value="tech" label="技术分享" />
              <t-option value="product" label="产品设计" />
              <t-option value="business" label="商业思考" />
              <t-option value="other" label="其他" />
            </t-select>
          </div>
          <div class="form-item">
            <label class="form-label">文章内容 <span class="required">*</span></label>
            <RichTextEditor 
              v-model="articleForm.content"
              placeholder="请输入文章内容，支持富文本格式..."
            />
          </div>
          <div class="form-item">
            <label class="form-label">文章摘要</label>
            <t-textarea 
              v-model="articleForm.summary"
              placeholder="请输入文章摘要（可选）"
              :autosize="{ minRows: 3, maxRows: 5 }"
            />
          </div>
          <div class="form-actions">
            <t-button theme="default" @click="resetForm">重置</t-button>
            <t-button theme="primary" @click="submitForm">发布文章</t-button>
          </div>
        </div>
      </div>

      <!-- 功能说明 -->
      <div class="demo-section">
        <h2>功能特性</h2>
        <div class="features-grid">
          <div class="feature-card">
            <div class="feature-icon">📝</div>
            <h3>文本格式</h3>
            <ul>
              <li>加粗、斜体、下划线</li>
              <li>删除线、清除格式</li>
              <li>H1/H2/H3 标题</li>
            </ul>
          </div>
          <div class="feature-card">
            <div class="feature-icon">📋</div>
            <h3>段落格式</h3>
            <ul>
              <li>无序列表</li>
              <li>有序列表</li>
              <li>引用块</li>
              <li>左/中/右对齐</li>
            </ul>
          </div>
          <div class="feature-card">
            <div class="feature-icon">🖼️</div>
            <h3>插入内容</h3>
            <ul>
              <li>图片上传（≤10MB）</li>
              <li>链接插入</li>
              <li>文件上传（≤50MB）</li>
            </ul>
          </div>
          <div class="feature-card">
            <div class="feature-icon">⚡</div>
            <h3>编辑操作</h3>
            <ul>
              <li>撤销/重做</li>
              <li>快捷键支持</li>
              <li>实时预览</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import RichTextEditor from '@/components/RichTextEditor/index.vue'

// 基础使用
const basicContent = ref('')
const showBasicHTML = ref(false)

const handleBasicChange = (html: string) => {
  console.log('内容更新:', html)
}

// 预设内容
const defaultPresetContent = `
<h2>欢迎使用富文本编辑器</h2>
<p>这是一个基于 <strong>TipTap</strong> 的 Vue 3 富文本编辑器组件。</p>
<h3>主要特性</h3>
<ul>
  <li>简洁美观的界面设计</li>
  <li>丰富的文本编辑功能</li>
  <li>支持图片和文件上传</li>
  <li>完善的快捷键支持</li>
</ul>
<blockquote>提示：你可以尝试使用工具栏的各种功能，或者上传图片和文件。</blockquote>
<p style="text-align: center;"><em>现在就开始编辑吧！</em></p>
`

const presetContent = ref(defaultPresetContent)

const resetPresetContent = () => {
  presetContent.value = defaultPresetContent
  MessagePlugin.success('已重置为预设内容')
}

// 只读模式
const readonlyContent = ref(`
<h3>这是只读模式的演示</h3>
<p>在只读模式下，编辑器的工具栏会被隐藏，内容无法编辑。</p>
<p>这种模式适合用于：</p>
<ul>
  <li>内容展示页面</li>
  <li>文章详情页</li>
  <li>评论或反馈的展示</li>
</ul>
<p>你可以通过设置 <code>:editable="false"</code> 来启用只读模式。</p>
`)

// 表单集成
const articleForm = ref({
  title: '',
  category: '',
  content: '',
  summary: '',
})

const resetForm = () => {
  articleForm.value = {
    title: '',
    category: '',
    content: '',
    summary: '',
  }
  MessagePlugin.success('表单已重置')
}

const submitForm = () => {
  if (!articleForm.value.title) {
    MessagePlugin.warning('请输入文章标题')
    return
  }
  if (!articleForm.value.content) {
    MessagePlugin.warning('请输入文章内容')
    return
  }
  
  console.log('提交表单:', articleForm.value)
  MessagePlugin.success('文章发布成功！')
}
</script>

<style scoped lang="less">
.editor-demo-page {
  min-height: 100%;
  background: #ffffff;
  border-radius: 8px;
  padding: 27px;

  .demo-header {
    margin-bottom: 32px;
    padding-bottom: 24px;
    border-bottom: 1px solid #e8e8e8;

    h1 {
      font-size: 28px;
      font-weight: 600;
      color: #1a1a1a;
      margin: 0 0 8px;
    }

    p {
      font-size: 15px;
      color: #666666;
      margin: 0;
    }
  }

  .demo-sections {
    display: flex;
    flex-direction: column;
    gap: 32px;
  }

  .demo-section {
    h2 {
      font-size: 20px;
      font-weight: 600;
      color: #1a1a1a;
      margin: 0 0 16px;
    }

    .demo-content {
      margin-bottom: 12px;
    }

    .demo-info {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 12px 16px;
      background: #f5f7fa;
      border-radius: 6px;
      font-size: 14px;
      color: #666666;

      .info-label {
        font-weight: 500;
      }

      .readonly-tag {
        display: inline-block;
        padding: 4px 12px;
        background: #fff3e0;
        color: #f57c00;
        border-radius: 4px;
        font-size: 13px;
        font-weight: 500;
      }
    }

    .demo-code {
      margin-top: 12px;
      padding: 16px;
      background: #1a1a1a;
      border-radius: 6px;
      overflow-x: auto;

      pre {
        margin: 0;
        color: #ffffff;
        font-family: 'Courier New', monospace;
        font-size: 13px;
        line-height: 1.6;
        white-space: pre-wrap;
        word-break: break-all;
      }
    }
  }

  .demo-form {
    padding: 24px;
    background: #fafafa;
    border-radius: 8px;

    .form-item {
      margin-bottom: 20px;

      &:last-of-type {
        margin-bottom: 24px;
      }

      .form-label {
        display: block;
        margin-bottom: 8px;
        font-size: 14px;
        font-weight: 500;
        color: #333333;

        .required {
          color: #f5222d;
          margin-left: 2px;
        }
      }
    }

    .form-actions {
      display: flex;
      gap: 12px;
      justify-content: flex-end;
      padding-top: 16px;
      border-top: 1px solid #e8e8e8;
    }
  }

  .features-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 16px;

    .feature-card {
      padding: 20px;
      background: #fafafa;
      border-radius: 8px;
      border: 1px solid #e8e8e8;
      transition: all 0.3s;

      &:hover {
        border-color: #4a7ff7;
        box-shadow: 0 4px 12px rgba(74, 127, 247, 0.15);
        transform: translateY(-2px);
      }

      .feature-icon {
        font-size: 32px;
        margin-bottom: 12px;
      }

      h3 {
        font-size: 16px;
        font-weight: 600;
        color: #1a1a1a;
        margin: 0 0 12px;
      }

      ul {
        margin: 0;
        padding-left: 20px;
        list-style-type: disc;

        li {
          font-size: 14px;
          color: #666666;
          line-height: 1.8;
          margin: 4px 0;
        }
      }
    }
  }
}
</style>
