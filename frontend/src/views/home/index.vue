<template>
  <div class="home-page">
    <!-- Tab 切换 -->
    <t-tabs :value="activeTab" @change="handleTabChange">
      <t-tab-panel value="recent" label="最近打开" />
      <t-tab-panel value="mine" label="我创建的" />
      <t-tab-panel value="shared" label="共享给我" />
      <t-tab-panel value="collect" label="收藏" />
      
      <template #action>
        <div class="view-switcher">
          <div 
            class="view-btn"
            :class="{ active: viewMode === 'grid' }"
            @click="viewMode = 'grid'"
          >
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect x="3" y="3" width="7" height="7" rx="1" />
              <rect x="14" y="3" width="7" height="7" rx="1" />
              <rect x="3" y="14" width="7" height="7" rx="1" />
              <rect x="14" y="14" width="7" height="7" rx="1" />
            </svg>
          </div>
          <div 
            class="view-btn"
            :class="{ active: viewMode === 'list' }"
            @click="viewMode = 'list'"
          >
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect x="3" y="5" width="18" height="2" rx="1" />
              <rect x="3" y="11" width="18" height="2" rx="1" />
              <rect x="3" y="17" width="18" height="2" rx="1" />
            </svg>
          </div>
        </div>
      </template>
    </t-tabs>

    <!-- 表格列表 - 卡片视图 -->
    <div v-if="viewMode === 'grid'">
      <!-- 空状态 -->
      <div v-if="!loading && tableList.length === 0" class="empty-state">
        <div class="empty-content">
          <svg class="empty-icon" viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect x="8" y="8" width="48" height="48" rx="4" stroke="currentColor" stroke-width="2"/>
            <line x1="16" y1="20" x2="48" y2="20" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <line x1="16" y1="32" x2="40" y2="32" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <line x1="16" y1="44" x2="44" y2="44" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <p class="empty-text">{{ getEmptyText() }}</p>
          <p class="empty-desc">{{ getEmptyDesc() }}</p>
        </div>
      </div>
      
      <!-- 数据列表 -->
      <div v-else class="table-grid">
      <div 
        v-for="table in tableList" 
        :key="table.id"
        class="table-card"
        @click="handleOpenTable(table)"
      >
        <div class="card-icon" :style="{ background: table.color }">
          {{ table.icon }}
        </div>
        <div class="card-content">
          <h3 class="card-title">{{ table.name }}</h3>
          <p class="card-meta">{{ table.updateTime }} {{ table.creator }}</p>
        </div>
        <div class="card-actions" @click.stop>
          <t-popup
            v-model="showDropdown[table.id]"
            trigger="click"
            placement="bottom-right"
            :overlay-style="{ padding: 0 }"
          >
            <t-button variant="text" shape="square">
              <t-icon name="more" />
            </t-button>
            <template #content>
              <div class="card-dropdown-menu">
                <div class="menu-item" @click="handleShare(table)">
                  <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                    <path d="M13 11a2 2 0 1 0 0 4 2 2 0 0 0 0-4zM3 6a2 2 0 1 0 0 4 2 2 0 0 0 0-4zM13 1a2 2 0 1 0 0 4 2 2 0 0 0 0-4z" stroke="currentColor" stroke-width="1.5"/>
                    <path d="M5.5 7.5l5-2M5.5 8.5l5 2" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                  </svg>
                  <span>分享</span>
                </div>
                <div class="menu-item" @click="handleCopyLink(table)">
                  <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                    <path d="M6.5 9.5l3-3M8 6l1.5-1.5a2.5 2.5 0 0 1 3.5 3.5L11.5 9.5M4.5 6.5L3 8a2.5 2.5 0 0 0 3.5 3.5L8 10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                  </svg>
                  <span>复制链接</span>
                </div>
                <div class="menu-divider"></div>
                <div class="menu-item" @click="handleRename(table)">
                  <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                    <path d="M7 2.5h6.5v11H7M2 8h8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  <span>重命名</span>
                </div>
                <div class="menu-item" @click="handleToggleFavorite(table)">
                  <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                    <path d="M8 2l2 4.5 4.5 0.5-3.5 3 1 4.5L8 12l-4 2.5 1-4.5-3.5-3L6 6.5z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  <span>{{ table.favorite ? '取消收藏' : '添加收藏' }}</span>
                </div>
                <div class="menu-divider"></div>
                <div class="menu-item danger" @click="handleDeleteConfirm(table)">
                  <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                    <path d="M3 4h10M6.5 4V3a1 1 0 0 1 1-1h1a1 1 0 0 1 1 1v1M5 4v8a1 1 0 0 0 1 1h4a1 1 0 0 0 1-1V4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  <span>删除</span>
                </div>
              </div>
            </template>
          </t-popup>
        </div>
      </div>
      </div>
    </div>

    <!-- 表格列表 - 列表视图 -->
    <div v-else>
      <!-- 空状态 -->
      <div v-if="!loading && tableList.length === 0" class="empty-state">
        <div class="empty-content">
          <svg class="empty-icon" viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect x="8" y="8" width="48" height="48" rx="4" stroke="currentColor" stroke-width="2"/>
            <line x1="16" y1="20" x2="48" y2="20" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <line x1="16" y1="32" x2="40" y2="32" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <line x1="16" y1="44" x2="44" y2="44" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <p class="empty-text">{{ getEmptyText() }}</p>
          <p class="empty-desc">{{ getEmptyDesc() }}</p>
        </div>
      </div>
      
      <!-- 数据表格 -->
      <div v-else class="table-list">
      <t-table
        :data="tableList"
        :columns="tableColumns"
        row-key="id"
        hover
        @row-click="({ row }) => handleOpenTable(row as ProjectItem)"
      >
        <template #name="{ row }">
          <div class="table-name-cell">
            <div class="cell-icon" :style="{ background: row.color }">
              {{ row.icon }}
            </div>
            <span>{{ row.name }}</span>
          </div>
        </template>
        <template #action="{ row }">
          <div @click.stop>
            <t-popup
              v-model="showDropdown[row.id]"
              trigger="click"
              placement="bottom-right"
              :overlay-style="{ padding: 0 }"
            >
              <t-button variant="text" shape="square">
                <t-icon name="more" />
              </t-button>
              <template #content>
                <div class="card-dropdown-menu">
                  <div class="menu-item" @click="handleShare(row)">
                    <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                      <path d="M13 11a2 2 0 1 0 0 4 2 2 0 0 0 0-4zM3 6a2 2 0 1 0 0 4 2 2 0 0 0 0-4zM13 1a2 2 0 1 0 0 4 2 2 0 0 0 0-4z" stroke="currentColor" stroke-width="1.5"/>
                      <path d="M5.5 7.5l5-2M5.5 8.5l5 2" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                    </svg>
                    <span>分享</span>
                  </div>
                  <div class="menu-item" @click="handleCopyLink(row)">
                    <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                      <path d="M6.5 9.5l3-3M8 6l1.5-1.5a2.5 2.5 0 0 1 3.5 3.5L11.5 9.5M4.5 6.5L3 8a2.5 2.5 0 0 0 3.5 3.5L8 10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                    </svg>
                    <span>复制链接</span>
                  </div>
                  <div class="menu-divider"></div>
                  <div class="menu-item" @click="handleRename(row)">
                    <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                      <path d="M7 2.5h6.5v11H7M2 8h8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                    <span>重命名</span>
                  </div>
                  <div class="menu-item" @click="handleToggleFavorite(row)">
                    <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                      <path d="M8 2l2 4.5 4.5 0.5-3.5 3 1 4.5L8 12l-4 2.5 1-4.5-3.5-3L6 6.5z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                    <span>{{ row.favorite ? '取消收藏' : '添加收藏' }}</span>
                  </div>
                  <div class="menu-divider"></div>
                  <div class="menu-item danger" @click="handleDeleteConfirm(row)">
                    <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                      <path d="M3 4h10M6.5 4V3a1 1 0 0 1 1-1h1a1 1 0 0 1 1 1v1M5 4v8a1 1 0 0 0 1 1h4a1 1 0 0 0 1-1V4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                    <span>删除</span>
                  </div>
                </div>
              </template>
            </t-popup>
          </div>
        </template>
      </t-table>
      </div>
    </div>

    <!-- 加载更多 -->
    <div v-if="hasMore" class="load-more">
      <t-button theme="default" @click="loadMore">下一页</t-button>
    </div>

    <!-- 删除确认对话框 -->
    <t-dialog
      v-model:visible="showDeleteDialog"
      header="删除确认"
      width="480px"
      @confirm="confirmDelete"
    >
      <div class="delete-confirm-content">
        <p>确定要删除表格「<strong>{{ pendingDeleteTable?.name }}</strong>」吗？</p>
        <p class="warning-text">删除后将无法恢复，请谨慎操作。</p>
      </div>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import type { TabValue } from 'tdesign-vue-next'
import { getProjectList, setProjectFavorite, deleteProject } from '@/api'
import type { ProjectItem } from '@/api'

const router = useRouter()
const activeTab = ref('recent')
const viewMode = ref<'grid' | 'list'>('grid')
const hasMore = ref(false)
const loading = ref(false)

// 分页参数
const pagination = reactive({
  page: 1,
  size: 20,
  total: 0,
})

// 下拉菜单显示状态
const showDropdown = reactive<Record<string, boolean>>({})

// 删除确认对话框
const showDeleteDialog = ref(false)
const pendingDeleteTable = ref<ProjectItem | null>(null)

// 项目列表
const tableList = ref<ProjectItem[]>([])

const tableColumns = [
  {
    colKey: 'name',
    title: '名称',
    width: 400,
  },
  {
    colKey: 'creator',
    title: '创建者',
    width: 150,
  },
  {
    colKey: 'updateTime',
    title: '更新时间',
    width: 200,
  },
  {
    colKey: 'action',
    title: '操作',
    width: 80,
    align: 'center' as const,
  },
]

// Tab 与 type 的映射关系
const tabTypeMap: Record<string, number> = {
  recent: 1,  // 最近打开
  mine: 2,    // 我创建的
  shared: 3,  // 共享给我
  collect: 4, // 收藏
}

// 默认图标和颜色池
const defaultIcons = ['📊', '📋', '📁', '📝', '💼', '🎯', '⚡', '🔥', '💡', '🌟', '🚀', '📈']
const defaultColors = [
  '#52c41a', '#4a7ff7', '#597ef7', '#9254de', '#722ed1',
  '#d3adf7', '#b37feb', '#ff7a45', '#f5222d', '#fa541c',
  '#fa8c16', '#faad14', '#fadb14', '#a0d911', '#52c41a',
  '#13c2c2', '#1890ff', '#2f54eb', '#eb2f96', '#f759ab'
]

// 生成随机默认图标和颜色
const getRandomIconAndColor = () => {
  const randomIcon = defaultIcons[Math.floor(Math.random() * defaultIcons.length)]
  const randomColor = defaultColors[Math.floor(Math.random() * defaultColors.length)]
  return { icon: randomIcon, color: randomColor }
}

// 为项目添加默认图标和颜色
const processProjectList = (projects: ProjectItem[]) => {
  return projects.map(project => {
    // 如果没有图标或颜色,则生成随机的
    if (!project.icon || !project.color) {
      const { icon, color } = getRandomIconAndColor()
      return {
        ...project,
        icon: project.icon || icon,
        color: project.color || color,
      }
    }
    return project
  })
}

// 加载项目列表
const loadProjectList = async (isLoadMore = false) => {
  if (loading.value) return
  
  loading.value = true
  try {
    const response = await getProjectList({
      page: pagination.page,
      size: pagination.size,
      type: tabTypeMap[activeTab.value],
    })
    
    if (response) {
      // 处理项目列表，添加默认图标和颜色
      const processedList = processProjectList(response.list)
      
      // 如果是加载更多，则追加到列表，否则替换列表
      if (isLoadMore) {
        tableList.value = [...tableList.value, ...processedList]
      } else {
        tableList.value = processedList
      }
      
      pagination.total = response.total
      pagination.page = response.page
      pagination.size = response.size
      
      // 判断是否还有更多数据
      hasMore.value = tableList.value.length < pagination.total
    }
  } catch (error: any) {
    console.error('加载项目列表失败:', error)
    MessagePlugin.error(error.message || '加载失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 页面加载时获取数据
onMounted(() => {
  loadProjectList()
})

const handleTabChange = (value: TabValue) => {
  activeTab.value = String(value)
  // 切换 tab 时重置分页并重新加载数据
  pagination.page = 1
  loadProjectList()
}

const handleOpenTable = (table: ProjectItem) => {
  // 跳转到表格编辑页面
  router.push(`/product/${table.id}`)
}

// 分享
const handleShare = (table: ProjectItem) => {
  showDropdown[table.id] = false
  MessagePlugin.info(`分享功能开发中：${table.name}`)
}

// 复制链接
const handleCopyLink = (table: ProjectItem) => {
  showDropdown[table.id] = false
  MessagePlugin.success('链接已复制到剪贴板')
}

// 重命名
const handleRename = (table: ProjectItem) => {
  showDropdown[table.id] = false
  MessagePlugin.info(`重命名功能开发中：${table.name}`)
}

// 切换收藏
const handleToggleFavorite = async (table: ProjectItem) => {
  showDropdown[table.id] = false
  
  try {
    const newFavoriteState = !table.favorite
    await setProjectFavorite({
      projectId: table.id,
      isFavorite: newFavoriteState,
    })
    
    // 更新本地状态
    const item = tableList.value.find(t => t.id === table.id)
    if (item) {
      item.favorite = newFavoriteState
      MessagePlugin.success(item.favorite ? '已添加到收藏' : '已取消收藏')
    }
  } catch (error: any) {
    console.error('更新收藏状态失败:', error)
    MessagePlugin.error(error.message || '操作失败，请稍后重试')
  }
}

// 删除确认
const handleDeleteConfirm = (table: ProjectItem) => {
  showDropdown[table.id] = false
  pendingDeleteTable.value = table
  showDeleteDialog.value = true
}

// 确认删除
const confirmDelete = async () => {
  if (!pendingDeleteTable.value) return
  
  try {
    await deleteProject({ id: pendingDeleteTable.value.id })
    
    // 从列表中移除
    const index = tableList.value.findIndex(t => t.id === pendingDeleteTable.value!.id)
    if (index > -1) {
      tableList.value.splice(index, 1)
      MessagePlugin.success(`已删除：${pendingDeleteTable.value.name}`)
    }
    
    // 更新总数
    pagination.total = Math.max(0, pagination.total - 1)
  } catch (error: any) {
    console.error('删除项目失败:', error)
    MessagePlugin.error(error.message || '删除失败，请稍后重试')
  } finally {
    showDeleteDialog.value = false
    pendingDeleteTable.value = null
  }
}

const loadMore = () => {
  pagination.page += 1
  loadProjectList(true)
}

// 获取空状态文字
const getEmptyText = () => {
  const textMap: Record<string, string> = {
    recent: '暂无最近打开的项目',
    mine: '暂无创建的项目',
    shared: '暂无共享的项目',
    collect: '暂无收藏的项目',
  }
  return textMap[activeTab.value] || '暂无数据'
}

// 获取空状态描述
const getEmptyDesc = () => {
  const descMap: Record<string, string> = {
    recent: '您还没有打开过任何项目',
    mine: '点击“+新建”按钮创建您的第一个项目',
    shared: '目前没有人分享项目给您',
    collect: '您还没有收藏任何项目',
  }
  return descMap[activeTab.value] || '暂无数据'
}
</script>

<style scoped lang="less">
.home-page {
  position: relative;
  min-height: 100%;
  background: #ffffff;
  border-radius: 8px;
  padding: 24px;
  
  :deep(.t-tabs) {
    margin-bottom: 20px;
  }
  
  .view-switcher {
    display: flex;
    gap: 4px;
    padding: 3px;
    background: #f5f7fa;
    border-radius: 4px;
    
    .view-btn {
      width: 24px;
      height: 24px;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      border-radius: 3px;
      transition: all 0.2s;
      
      svg {
        width: 16px;
        height: 16px;
        stroke: #666;
        stroke-width: 1.5;
        fill: none;
      }
      
      &:hover {
        background: #e8e8e8;
        
        svg {
          stroke: #333;
        }
      }
      
      &.active {
        background: #ffffff;
        
        svg {
          stroke: #4a7ff7;
          fill: #4a7ff7;
        }
      }
    }
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
  
  // 卡片视图
  .table-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 16px;
    
    .table-card {
      display: flex;
      align-items: center;
      padding: 16px;
      background: #ffffff;
      border-radius: 8px;
      border: 1px solid #e8e8e8;
      cursor: pointer;
      transition: all 0.3s;
      
      &:hover {
        border-color: #4a7ff7;
        box-shadow: 0 2px 8px rgba(74, 127, 247, 0.15);
      }
      
      .card-icon {
        width: 48px;
        height: 48px;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #ffffff;
        font-size: 16px;
        font-weight: 500;
        flex-shrink: 0;
        margin-right: 12px;
      }
      
      .card-content {
        flex: 1;
        min-width: 0;
        
        .card-title {
          font-size: 14px;
          font-weight: 500;
          color: #333;
          margin: 0 0 4px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
        
        .card-meta {
          font-size: 12px;
          color: #999;
          margin: 0;
        }
      }
      
      .card-actions {
        flex-shrink: 0;
        opacity: 0;
        transition: opacity 0.3s;
      }
      
      &:hover .card-actions {
        opacity: 1;
      }
    }
  }
  
  // 列表视图
  .table-list {
    background: #ffffff;
    border-radius: 8px;
    
    .table-name-cell {
      display: flex;
      align-items: center;
      gap: 12px;
      
      .cell-icon {
        width: 32px;
        height: 32px;
        border-radius: 6px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #ffffff;
        font-size: 14px;
        font-weight: 500;
        flex-shrink: 0;
      }
    }
  }
  
  // 加载更多
  .load-more {
    display: flex;
    justify-content: center;
    margin-top: 32px;
    padding-bottom: 80px;
  }
}

// 卡片下拉菜单样式
.card-dropdown-menu {
  min-width: 160px;
  padding: 6px 0;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 3px 14px rgba(0, 0, 0, 0.12);

  .menu-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 16px;
    cursor: pointer;
    transition: background 0.2s;
    color: #333333;
    font-size: 14px;

    svg {
      flex-shrink: 0;
      color: #666666;
    }

    &:hover {
      background: #f5f7fa;
    }

    &.danger {
      color: #e34d59;

      svg {
        color: #e34d59;
      }

      &:hover {
        background: #fff1f0;
      }
    }
  }

  .menu-divider {
    height: 1px;
    background: #e8e8e8;
    margin: 6px 0;
  }
}

// 删除确认对话框样式
.delete-confirm-content {
  padding: 8px 0;

  p {
    margin: 0 0 12px;
    font-size: 14px;
    line-height: 1.6;
    color: #333333;

    &:last-child {
      margin-bottom: 0;
    }

    strong {
      font-weight: 600;
      color: #4a7ff7;
    }
  }

  .warning-text {
    color: #e34d59;
    font-size: 13px;
  }
}
</style>
