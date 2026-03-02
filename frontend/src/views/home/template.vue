<template>
  <div class="template-page">
    <!-- 分类导航 -->
    <div class="category-tabs">
      <t-tabs :value="activeCategory" @change="handleCategoryChange">
        <t-tab-panel 
          v-for="category in categoryList" 
          :key="category.value"
          :value="category.value" 
          :label="category.label" 
        />
      </t-tabs>
    </div>

    <!-- 模板列表 -->
    <div v-if="loading && templateList.length === 0" class="loading-state">
      <t-loading size="large" text="加载中..." />
    </div>
    
    <div v-else-if="templateList.length === 0" class="empty-state">
      <div class="empty-content">
        <svg class="empty-icon" viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
          <rect x="8" y="8" width="48" height="48" rx="4" stroke="currentColor" stroke-width="2"/>
          <line x1="16" y1="20" x2="48" y2="20" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <line x1="16" y1="32" x2="40" y2="32" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <line x1="16" y1="44" x2="44" y2="44" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
        <p class="empty-text">暂无模板</p>
        <p class="empty-desc">当前分类下暂无模板数据</p>
      </div>
    </div>
    
    <div v-else class="template-list">
      <div 
        v-for="template in templateList" 
        :key="template.id"
        class="template-card"
      >
        <div class="card-preview" @click="handlePreview(template)">
          <img :src="template.cover" :alt="template.name" />
          <div class="preview-overlay">
            <t-icon name="browse" size="32px" />
          </div>
        </div>
        <div class="card-info">
          <div class="card-header">
            <h3 class="card-title" :title="template.name">{{ template.name }}</h3>
            <t-tag v-if="template.tag" size="small" theme="primary" variant="light">
              {{ template.tag }}
            </t-tag>
          </div>
          <p class="card-desc" :title="template.description">{{ template.description }}</p>
          <div class="card-footer">
            <div class="creator-info">
              <t-avatar size="20px">{{ template.creator.charAt(0) }}</t-avatar>
              <span class="creator-name">{{ template.creator }}</span>
            </div>
            <div class="usage-info">
              <t-icon name="browse" size="14px" />
              <span>{{ template.usage }}</span>
            </div>
          </div>
        </div>
        <div class="card-actions">
          <t-button 
            class="preview-btn"
            variant="outline" 
            size="small"
            @click="handlePreview(template)"
          >
            <t-icon name="browse" />
            预览
          </t-button>
          <t-button 
            class="use-btn"
            theme="primary" 
            size="small"
            @click="handleUseTemplate(template)"
          >
            <t-icon name="check" />
            使用模版
          </t-button>
        </div>
      </div>
    </div>

    <!-- 加载更多 -->
    <div v-if="hasMore" class="load-more">
      <t-button 
        theme="default" 
        :loading="loading"
        @click="loadMore"
      >
        {{ loading ? '加载中...' : '加载更多' }}
      </t-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import type { TabValue } from 'tdesign-vue-next'
import { getTemplateList, getTemplateTagList, createProject } from '@/api'

const router = useRouter()
const activeCategory = ref('热门推荐') // 默认为"热门推荐"
const hasMore = ref(false)
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const totalCount = ref(0)

// 默认预览图
const defaultPreview = 'https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=600&h=400&fit=crop'

// 模板数据
const templateList = ref<any[]>([])

// 分类列表
const categoryList = ref<Array<{ label: string; value: string }>>([])

// 加载模板分类
const loadCategories = async () => {
  try {
    const response = await getTemplateTagList({ searchWord: '' }) as any
    if (response && response.list) {
      categoryList.value = response.list.map((tag: any) => ({
        label: tag.name,
        value: tag.name, // 使用name字段作为tag参数
      }))
      // 查找是否有"热门推荐"分类
      const hotCategory = categoryList.value.find(cat => cat.value === '热门推荐')
      if (hotCategory) {
        activeCategory.value = '热门推荐'
      } else if (categoryList.value.length > 0) {
        // 如果没有"热门推荐"，选择第一个分类
        activeCategory.value = categoryList.value[0].value
      }
    }
  } catch (error: any) {
    console.error('加载模板分类失败:', error)
    // 设置默认分类
    categoryList.value = [
      { label: '热门推荐', value: '热门推荐' },
      { label: '营销主题', value: '营销主题' },
      { label: '电子商务', value: '电子商务' },
      { label: '销售', value: '销售' },
      { label: '运营', value: '运营' },
    ]
    activeCategory.value = '热门推荐'
  }
}

// 加载模板列表
const loadTemplateList = async (isLoadMore = false) => {
  if (loading.value) return
  
  loading.value = true
  try {
    const response = await getTemplateList({
      page: currentPage.value,
      size: pageSize.value,
      tag: activeCategory.value || '', // 确保tag参数存在，没有则传空字符串
    }) as any
    
    if (response && response.list) {
      const newTemplates = response.list.map((item: any) => ({
        id: item.id,
        name: item.name || '未命名模板',
        usage: item.usage || '0人使用',
        preview: item.cover || item.image || defaultPreview,
        tag: item.tag || (item.name?.includes('AI') ? 'AI' : ''),
        description: item.description || '暂无描述',
        creator: item.creator || '系统',
        cover: item.cover || defaultPreview,
      }))
      
      // 加载更多时追加，否则替换
      if (isLoadMore) {
        templateList.value = [...templateList.value, ...newTemplates]
      } else {
        templateList.value = newTemplates
      }
      
      // 更新分页信息
      totalCount.value = response.total || 0
      hasMore.value = templateList.value.length < totalCount.value
    }
  } catch (error: any) {
    console.error('加载模板列表失败:', error)
    MessagePlugin.error(error.message || '加载模板列表失败')
    templateList.value = []
    hasMore.value = false
  } finally {
    loading.value = false
  }
}

const handleCategoryChange = (value: TabValue) => {
  activeCategory.value = String(value)
  currentPage.value = 1
  loadTemplateList()
}

const handlePreview = (template: any) => {
  console.log('预览模板:', template)
  MessagePlugin.info(`预览功能开发中：${template.name}`)
}

const handleUseTemplate = async (template: any) => {
  try {
    // 基于模板创建项目
    const response = await createProject({
      name: template.name,
      description: template.description || '',
    }) as any
    
    if (response && response.id) {
      MessagePlugin.success(`已基于模板「${template.name}」创建项目`)
      // 跳转到表格编辑页面
      router.push(`/product/${response.id}`)
    }
  } catch (error: any) {
    console.error('使用模板创建项目失败:', error)
    MessagePlugin.error(error.message || '创建失败，请稍后重试')
  }
}

const loadMore = () => {
  if (hasMore.value && !loading.value) {
    currentPage.value += 1
    loadTemplateList(true)
  }
}

// 监听分类变化
watch(activeCategory, (newVal, oldVal) => {
  // 只有当值真正改变时才重新加载（避免初始化时触发）
  if (oldVal !== undefined && newVal !== oldVal) {
    currentPage.value = 1
    loadTemplateList()
  }
})

// 页面加载时初始化
onMounted(async () => {
  await loadCategories()
  // 分类加载完成后，加载模板列表
  loadTemplateList()
})
</script>

<style scoped lang="less">
.template-page {
  min-height: 100%;
  background: #ffffff;
  border-radius: 8px;
  padding: 24px;

  .category-tabs {
    margin-bottom: 24px;

    :deep(.t-tabs__nav) {
      border-bottom: 1px solid #e8e8e8;
    }

    :deep(.t-tab-panel) {
      padding: 0;
    }
  }
  
  // 加载状态
  .loading-state {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 400px;
    padding: 60px 20px;
  }
  
  // 空状态
  .empty-state {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 400px;
    padding: 60px 20px;
    
    .empty-content {
      text-align: center;
      max-width: 400px;
      
      .empty-icon {
        width: 120px;
        height: 120px;
        color: #dcdcdc;
        margin: 0 auto 24px;
      }
      
      .empty-text {
        font-size: 16px;
        font-weight: 500;
        color: #333;
        margin: 0 0 8px;
      }
      
      .empty-desc {
        font-size: 14px;
        color: #999;
        margin: 0;
        line-height: 1.6;
      }
    }
  }

  .template-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;

    .template-card {
      position: relative;
      background: #ffffff;
      border-radius: 12px;
      border: 1px solid #e8e8e8;
      overflow: hidden;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      cursor: pointer;

      &:hover {
        transform: translateY(-6px);
        border-color: #4a7ff7;
        box-shadow: 0 12px 32px rgba(74, 127, 247, 0.15);

        .card-preview {
          .preview-overlay {
            opacity: 1;
          }
          
          img {
            transform: scale(1.08);
          }
        }

        .card-actions {
          opacity: 1;
          transform: translateY(0);
        }
      }

      .card-preview {
        position: relative;
        width: 100%;
        height: 180px;
        background: #f5f7fa;
        overflow: hidden;
        
        .preview-overlay {
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          background: rgba(0, 0, 0, 0.4);
          display: flex;
          align-items: center;
          justify-content: center;
          color: white;
          opacity: 0;
          transition: opacity 0.3s ease;
        }

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
          transition: transform 0.4s cubic-bezier(0.4, 0, 0.2, 1);
        }
      }

      .card-info {
        padding: 14px 16px;

        .card-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          gap: 8px;
          margin-bottom: 10px;

          .card-title {
            flex: 1;
            font-size: 15px;
            font-weight: 500;
            color: #1f2329;
            margin: 0;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            line-height: 1.4;
          }
        }

        .card-desc {
          font-size: 13px;
          color: #646a73;
          line-height: 1.5;
          margin: 0 0 12px 0;
          height: 40px;
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
        }
        
        .card-footer {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding-top: 10px;
          border-top: 1px solid #f0f0f0;
          
          .creator-info {
            display: flex;
            align-items: center;
            gap: 6px;
            flex: 1;
            min-width: 0;
            
            .creator-name {
              font-size: 12px;
              color: #646a73;
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
            }
          }
          
          .usage-info {
            display: flex;
            align-items: center;
            gap: 4px;
            font-size: 12px;
            color: #8f959e;
            flex-shrink: 0;
          }
        }
      }

      .card-actions {
        display: flex;
        gap: 8px;
        padding: 0 16px 14px;
        opacity: 0;
        transform: translateY(-8px);
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

        .preview-btn {
          flex: 1;
          border-color: #4a7ff7;
          color: #4a7ff7;
          
          :deep(.t-icon) {
            margin-right: 4px;
          }

          &:hover {
            background: #f0f5ff;
            border-color: #4a7ff7;
          }
        }

        .use-btn {
          flex: 1.2;
          
          :deep(.t-icon) {
            margin-right: 4px;
          }
        }
      }
    }
  }

  .load-more {
    display: flex;
    justify-content: center;
    margin-top: 32px;
    padding-bottom: 20px;
  }
}
</style>
