<template>
  <div class="home-layout">
    <!-- 顶部导航栏 -->
    <header class="home-header">
      <div class="header-left">
        <div class="logo">
          <div class="logo-icon">AI</div>
          <span class="logo-text">迈维AI表格</span>
        </div>
      </div>
      <div class="header-right">
        <t-button theme="primary" @click="handleCreate">
          <template #icon>
            <t-icon name="add" />
          </template>
          新建
        </t-button>
        <t-popup
          v-model="showSpaceDropdown"
          trigger="click"
          placement="bottom-left"
          :overlay-style="{ padding: 0 }"
        >
          <div class="space-dropdown">
            <t-icon name="user-circle" />
            <span>{{ teamStore.currentTeam?.name || '选择团队' }}</span>
            <svg class="dropdown-arrow" width="14" height="14" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M4 6L8 10L12 6" stroke="#999999" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <template #content>
            <div class="space-dropdown-content">
              <div class="menu-title">企业 / 组织 / 团队</div>
              <div class="space-list">
                <div
                  v-for="team in teamStore.teamList"
                  :key="team.id"
                  class="space-item"
                  :class="{ active: teamStore.currentTeam?.id === team.id }"
                  @click="handleSpaceChange(team)"
                >
                  <div v-if="team.logo" class="space-logo">
                    <img :src="team.logo" alt="团队logo" />
                  </div>
                  <div v-else class="space-icon" :style="{ background: team.color }">
                    <t-icon name="usergroup" />
                  </div>
                  <span class="space-name">{{ team.name }}</span>
                </div>
              </div>
              <div class="create-space" @click="handleCreateSpace">
                <t-icon name="add" />
                <span>创建企业 / 组织 / 团队</span>
              </div>
            </div>
          </template>
        </t-popup>
        <t-dropdown trigger="click" :min-column-width="120">
          <t-avatar size="32px" style="cursor: pointer;">用</t-avatar>
          <t-dropdown-menu>
            <t-dropdown-item @click="handleUserMenuClick({ value: 'profile' })">
              个人信息
            </t-dropdown-item>
            <t-dropdown-item @click="handleUserMenuClick({ value: 'logout' })">
              退出登录
            </t-dropdown-item>
          </t-dropdown-menu>
        </t-dropdown>
      </div>
    </header>

    <!-- 主体内容 -->
    <div class="home-body">
      <!-- 左侧导航 -->
      <aside class="home-sidebar">
        <t-menu
          :value="currentMenu"
          theme="light"
          @change="handleMenuChange"
          width="174px"
        >
          <t-menu-item value="home">
            <template #icon>
              <t-icon name="home" />
            </template>
            主页
          </t-menu-item>
          <t-menu-item value="team">
            <template #icon>
              <t-icon name="usergroup" />
            </template>
            团队管理
          </t-menu-item>
          <t-menu-item value="template">
            <template #icon>
              <t-icon name="template" />
            </template>
            模版中心
          </t-menu-item>
          <t-menu-item value="settings">
            <template #icon>
              <t-icon name="setting" />
            </template>
            个人设置
          </t-menu-item>
          <t-menu-item value="feedback">
            <template #icon>
              <t-icon name="chat" />
            </template>
            问题反馈
          </t-menu-item>
          <t-menu-item value="editor">
            富文本编辑器
          </t-menu-item>
          <t-menu-item value="upload-demo">
            <template #icon>
              <t-icon name="upload" />
            </template>
            文件上传Demo
          </t-menu-item>
          <t-menu-item value="watermark-demo">
            <template #icon>
              <t-icon name="secured" />
            </template>
            水印编辑器Demo
          </t-menu-item>
        </t-menu>
      </aside>

      <!-- 右侧内容区 -->
      <main class="home-content">
        <router-view />
      </main>
    </div>

    <!-- 全局浮动按钮组 -->
    <div class="floating-buttons">
      <!-- 全局创建按钮 -->
      <div class="floating-create">
        <t-button theme="primary" shape="circle" @click="handleCreate">
          <template #icon>
            <t-icon name="add" />
          </template>
        </t-button>
      </div>
      
      <!-- 私有化部署按钮 -->
      <div class="floating-deploy">
        <t-button theme="default" shape="circle" @click="showDeployDialog = true">
          <span class="deploy-icon">私</span>
        </t-button>
      </div>
    </div>

    <!-- 创建团队弹窗 -->
    <t-dialog
      v-model:visible="showCreateDialog"
      header="创建企业 / 组织 / 团队"
      width="520px"
      :on-confirm="handleConfirmCreate"
      :on-cancel="() => showCreateDialog = false"
    >
      <div class="create-space-content">
        <div class="form-item">
          <div class="form-label">名称：</div>
          <t-input
            v-model="newSpace.name"
            placeholder="请输入名称"
            clearable
          />
        </div>
        <div class="form-item">
          <div class="form-label">描述：</div>
          <t-textarea
            v-model="newSpace.description"
            placeholder="请输入描述（可选）"
            :autosize="{ minRows: 3, maxRows: 5 }"
          />
        </div>
        <div class="form-item">
          <div class="form-label">Logo：</div>
          <div class="logo-upload">
            <input
              ref="spaceLogoInputRef"
              type="file"
              accept="image/*"
              style="display: none"
              @change="handleSpaceLogoChange"
            />
            <div
              class="upload-area"
              :class="{ 'has-logo': newSpace.logo }"
              @click="handleUploadSpaceLogo"
            >
              <img v-if="newSpace.logo" :src="newSpace.logo" alt="Logo" class="logo-preview" />
              <t-icon v-else name="add" size="32px" />
            </div>
          </div>
        </div>
      </div>
    </t-dialog>

    <!-- 新建表格弹窗 -->
    <t-dialog
      v-model:visible="showNewTableDialog"
      header="新建"
      width="900px"
      :footer="false"
    >
      <div class="create-dialog-content">
        <!-- 新建方式 -->
        <div class="create-options">
          <div class="option-card" @click="handleCreateBlank">
            <div class="option-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <rect x="4" y="4" width="16" height="16" rx="2" stroke="white" stroke-width="2"/>
                <path d="M8 8h8M8 12h8M8 16h5" stroke="white" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </div>
            <div class="option-text">
              <div class="option-title">新建空白表格</div>
              <div class="option-desc">从空白开始</div>
            </div>
          </div>
          <div class="option-card" @click="handleImportExcel">
            <div class="option-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6z" stroke="white" stroke-width="2"/>
                <path d="M14 2v6h6M10 13l2 2 2-2M12 15v-4" stroke="white" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </div>
            <div class="option-text">
              <div class="option-title">导入 Excel/在线表格</div>
              <div class="option-desc">导入已有数据或连接现有</div>
            </div>
          </div>
          <div class="option-card" @click="handleCreateList">
            <div class="option-icon" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <path d="M9 5l7 7-7 7" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <div class="option-text">
              <div class="option-title">新建表单</div>
              <div class="option-desc">收集数据到表格</div>
            </div>
          </div>
          <div class="option-card" @click="handleCreateDashboard">
            <div class="option-icon" style="background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <path d="M3 3v18h18" stroke="white" stroke-width="2" stroke-linecap="round"/>
                <path d="M7 16l4-4 3 3 5-5" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <div class="option-text">
              <div class="option-title">新建仪表盘</div>
              <div class="option-desc">将数据自动转换成图表</div>
            </div>
          </div>
        </div>

        <!-- 从模板新建 -->
        <div class="template-section">
          <h3 class="section-title">从模板新建</h3>
          
          <!-- 模板分类 -->
          <div class="template-tabs">
            <div 
              class="template-tab"
              :class="{ active: activeTemplateTab === tab.value }"
              v-for="tab in templateTabs"
              :key="tab.value"
              @click="activeTemplateTab = tab.value"
            >
              {{ tab.label }}
            </div>
          </div>

          <!-- 热门推荐模板 -->
          <div class="template-grid">
            <div class="template-card" @click="handleCreateEmpty">
              <div class="template-image empty">
                <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
                  <path d="M24 14v20M14 24h20" stroke="#999" stroke-width="3" stroke-linecap="round"/>
                </svg>
                <div class="empty-text">新建空白 AI 表格</div>
              </div>
            </div>
            <div 
              class="template-card"
              v-for="template in templateList"
              :key="template.id"
              @click="handleUseTemplate(template)"
            >
              <div class="template-image">
                <img :src="template.image" :alt="template.name" />
              </div>
              <div class="template-info">
                <div class="template-header">
                  <div class="template-name" :title="template.name">{{ template.name }}</div>
                </div>
                <div class="template-desc" :title="template.description">{{ template.description }}</div>
                <div class="template-footer">
                  <div class="creator-info">
                    <t-avatar size="16px">{{ template.creator.charAt(0) }}</t-avatar>
                    <span class="creator-name">{{ template.creator }}</span>
                  </div>
                  <div class="usage-info">
                    <t-icon name="browse" size="12px" />
                    <span>{{ template.usage }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </t-dialog>
    
    <!-- 私有化部署弹窗 -->
    <t-dialog
      v-model:visible="showDeployDialog"
      header="私有化部署"
      width="600px"
      :footer="false"
    >
      <div class="deploy-content">
        <div class="qrcode-section">
          <div class="qrcode-item">
            <img src="https://flowmix.turntip.cn/fm/static/my.8ee63da4.png" alt="作者微信" class="qrcode-img" />
            <div class="qrcode-label">作者微信</div>
          </div>
          <div class="qrcode-item">
            <img src="https://jitword.com/assets/media/wechat.png" alt="产品公众号" class="qrcode-img" />
            <div class="qrcode-label">产品公众号</div>
          </div>
        </div>
      </div>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import type { MenuValue } from 'tdesign-vue-next'
import { createProject, getTemplateList, getTemplateTagList, createTeam, uploadFile } from '@/api'
import { useTeamStore } from '@/stores'

const router = useRouter()
const route = useRoute()
const currentMenu = ref('home')

// 使用团队 store
const teamStore = useTeamStore()

// 空间相关状态
const showCreateDialog = ref(false)
const showSpaceDropdown = ref(false)
const showNewTableDialog = ref(false)
const showDeployDialog = ref(false)
const activeTemplateTab = ref('')
const loadingTemplates = ref(false)

// 模板分类
const templateTabs = ref<Array<{ label: string; value: string }>>([])

// 模板列表
const templateList = ref<any[]>([])

// 新建空间表单
const newSpace = ref({
  name: '',
  description: '',
  logo: '',
})

// Logo上传输入框引用
const spaceLogoInputRef = ref<HTMLInputElement>()

const userMenuOptions = [
  { content: '个人信息', value: 'profile' },
  { content: '退出登录', value: 'logout' },
]

watch(() => route.name, (newName) => {
  if (newName === 'home') {
    currentMenu.value = 'home'
  } else if (newName === 'team') {
    currentMenu.value = 'team'
  } else if (newName === 'template') {
    currentMenu.value = 'template'
  } else if (newName === 'settings') {
    currentMenu.value = 'settings'
  } else if (newName === 'feedback') {
    currentMenu.value = 'feedback'
  } else if (newName === 'upload-demo') {
    currentMenu.value = 'upload-demo'
  } else if (newName === 'watermark-demo') {
    currentMenu.value = 'watermark-demo'
  }
}, { immediate: true })

const handleMenuChange = (value: MenuValue) => {
  currentMenu.value = String(value)
  const routeMap: Record<string, string> = {
    home: '/home',
    team: '/home/team',
    template: '/home/template',
    learning: '/home/learning',
    settings: '/home/settings',
    feedback: '/home/feedback',
    editor: '/home/editor',
    'upload-demo': '/home/upload-demo',
    'watermark-demo': '/home/watermark-demo',
  }
  if (routeMap[String(value)]) {
    router.push(routeMap[String(value)])
  }
}

const handleCreate = () => {
  showNewTableDialog.value = true
}

// 新建空白表格
const handleCreateBlank = async () => {
  showNewTableDialog.value = false
  
  try {
    const response = await createProject({
      name: '未命名',
      description: '',
    })
    
    if (response && response.id) {
      MessagePlugin.success('创建成功')
      // 跳转到表格编辑页面
      router.push(`/product/${response.id}`)
    }
  } catch (error: any) {
    console.error('创建项目失败:', error)
    MessagePlugin.error(error.message || '创建失败，请稍后重试')
  }
}

// 导入Excel
const handleImportExcel = () => {
  showNewTableDialog.value = false
  MessagePlugin.info('导入Excel功能开发中')
}

// 新建表单
const handleCreateList = () => {
  showNewTableDialog.value = false
  MessagePlugin.info('新建表单功能开发中')
}

// 新建仪表盘
const handleCreateDashboard = () => {
  showNewTableDialog.value = false
  MessagePlugin.info('新建仪表盘功能开发中')
}

// 创建空白AI表格
const handleCreateEmpty = async () => {
  showNewTableDialog.value = false
  
  try {
    const response = await createProject({
      name: '未命名',
      description: '',
    })
    
    if (response && response.id) {
      MessagePlugin.success('创建成功')
      // 跳转到表格编辑页面
      router.push(`/product/${response.id}`)
    }
  } catch (error: any) {
    console.error('创建项目失败:', error)
    MessagePlugin.error(error.message || '创建失败，请稍后重试')
  }
}

// 加载模板分类
const loadTemplateTags = async () => {
  try {
    const response = await getTemplateTagList({ searchWord: '' }) as any
    if (response && response.list) {
      // 添加"全部"标签
      templateTabs.value = [
        { label: '全部', value: '' },
        ...response.list.map((tag: any) => ({
          label: tag.name,
          value: tag.name, // 使用name字段作为tag参数
        }))
      ]
      // 默认选中第一个
      if (templateTabs.value.length > 0) {
        activeTemplateTab.value = templateTabs.value[0].value
      }
    }
  } catch (error: any) {
    console.error('加载模板分类失败:', error)
    // 设置默认分类
    templateTabs.value = [
      { label: '全部', value: '' },
      { label: '热门推荐', value: 'hot' },
      { label: '最新上线', value: 'new' },
    ]
    activeTemplateTab.value = ''
  }
}

// 加载模板列表
const loadTemplateList = async () => {
  loadingTemplates.value = true
  try {
    const response = await getTemplateList({
      page: 1,
      size: 20,
      tag: activeTemplateTab.value || '', // 确保tag参数存在，没有则传空字符串
    }) as any
    
    if (response && response.list) {
      templateList.value = response.list.map((item: any) => ({
        id: item.id,
        name: item.name || '未命名模板',
        usage: item.usage || '0人使用',
        image: item.cover || item.image || 'https://via.placeholder.com/200x120/4a7ff7/ffffff?text=' + encodeURIComponent(item.name || '模板'),
        description: item.description || '暂无描述',
        creator: item.creator || '系统',
      }))
    }
  } catch (error: any) {
    console.error('加载模板列表失败:', error)
    MessagePlugin.error(error.message || '加载模板列表失败')
    // 设置空数组
    templateList.value = []
  } finally {
    loadingTemplates.value = false
  }
}

// 使用模板
const handleUseTemplate = async (template: any) => {
  showNewTableDialog.value = false
  
  try {
    // 基于模板创建项目
    const response = await createProject({
      name: template.name,
      description: template.description || '',
    })
    
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

// 监听分类切换
watch(activeTemplateTab, () => {
  loadTemplateList()
})

// 监听弹窗显示状态
watch(showNewTableDialog, (newVal) => {
  if (newVal) {
    // 打开弹窗时加载数据
    if (templateTabs.value.length === 0) {
      loadTemplateTags()
    }
    loadTemplateList()
  }
})

// 页面加载时初始化模板分类
onMounted(async () => {
  loadTemplateTags()
  // 加载团队列表
  try {
    await teamStore.loadTeamList()
  } catch (error: any) {
    console.error('加载团队列表失败:', error)
  }
})

// 切换空间（团队）
const handleSpaceChange = (team: any) => {
  teamStore.setCurrentTeam(team)
  showSpaceDropdown.value = false
  MessagePlugin.success(`已切换到：${team.name}`)
}

// 打开创建空间弹窗
const handleCreateSpace = () => {
  showSpaceDropdown.value = false
  showCreateDialog.value = true
  newSpace.value = {
    name: '',
    description: '',
    logo: '',
  }
}

// 上传空间Logo
const handleUploadSpaceLogo = () => {
  spaceLogoInputRef.value?.click()
}

// Logo文件变化处理
const handleSpaceLogoChange = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    // 验证文件类型
    if (!file.type.startsWith('image/')) {
      MessagePlugin.warning('请选择图片文件')
      return
    }
    // 验证文件大小（5MB）
    if (file.size > 5 * 1024 * 1024) {
      MessagePlugin.warning('图片大小不能超过5MB')
      return
    }
    
    try {
      // 上传文件
      const response = await uploadFile(file)
      if (response && response.url) {
        newSpace.value.logo = response.url
        MessagePlugin.success('Logo已上传')
      }
    } catch (error: any) {
      console.error('上传Logo失败:', error)
      MessagePlugin.error(error.message || '上传失败')
    }
  }
}

// 确认创建空间
const handleConfirmCreate = async () => {
  if (!newSpace.value.name.trim()) {
    MessagePlugin.warning('请输入名称')
    return
  }
  
  try {
    const response = await createTeam({
      name: newSpace.value.name,
      description: newSpace.value.description,
      logo: newSpace.value.logo,
    }) as any
    
    if (response) {
      // 添加到 store
      const newTeamData = {
        id: response.id || Date.now().toString(),
        name: newSpace.value.name,
        description: newSpace.value.description,
        logo: newSpace.value.logo,
      }
      teamStore.addTeam(newTeamData)
      teamStore.setCurrentTeam(newTeamData)
      
      MessagePlugin.success(`已创建：${newSpace.value.name}`)
      showCreateDialog.value = false
    }
  } catch (error: any) {
    console.error('创建团队失败:', error)
    MessagePlugin.error(error.message || '创建失败')
  }
}

const handleUserMenuClick = (data: any) => {
  if (data.value === 'logout') {
    router.push('/user/login')
  } else if (data.value === 'profile') {
    currentMenu.value = 'settings'
    router.push('/home/settings')
  }
}
</script>

<style scoped lang="less">
.home-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #f5f7fa;
}

.home-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
  padding: 0 24px;
  background: #ffffff;
  border-bottom: 1px solid #e8e8e8;

  .header-left {
    .logo {
      display: flex;
      align-items: center;
      gap: 8px;

      .logo-icon {
        width: 32px;
        height: 32px;
        background: #4a7ff7;
        border-radius: 6px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: bold;
        color: #fff;
        font-size: 18px;
      }

      .logo-text {
        font-size: 16px;
        font-weight: 500;
        color: #333;
      }
    }
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 16px;

    .space-dropdown {
      display: flex;
      align-items: center;
      gap: 4px;
      padding: 4px 12px;
      cursor: pointer;
      border-radius: 4px;
      transition: background 0.3s;

      &:hover {
        background: #f5f5f5;
      }
      
      .dropdown-arrow {
        margin-left: 4px;
        transition: transform 0.3s;
      }
    }
  }
}

.home-body {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.home-sidebar {
  width: 200px;
  background: #ffffff;
  border-right: 1px solid #e8e8e8;
  overflow-y: auto;
  overflow-x: hidden;
  flex-shrink: 0;
  padding: 12px;

  :deep(.t-menu) {
    border: none;
    background: transparent;
  }

  :deep(.t-menu__item) {
    margin: 4px 0;
    border-radius: 6px;
    padding: 0 12px;
    height: 36px;
    line-height: 36px;
    font-size: 14px;
    color: #333;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;

    .t-icon {
      font-size: 16px;
      margin-right: 8px;
    }

    &:hover {
      background: #f5f7fa;
    }

    &.t-is-active {
      background: #e8f0ff;
      color: #4a7ff7;

      .t-icon {
        color: #4a7ff7;
      }
    }
  }
}

.home-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

// 全局浮动按钮组
.floating-buttons {
  position: fixed;
  bottom: 40px;
  right: 40px;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  gap: 16px;
  align-items: center;
}

.floating-create {
  :deep(.t-button) {
    width: 48px !important;
    height: 48px !important;
    min-width: 48px !important;
    padding: 0 !important;
    border-radius: 50% !important;
    box-shadow: 0 4px 12px rgba(74, 127, 247, 0.4);
    
    &:hover {
      box-shadow: 0 6px 16px rgba(74, 127, 247, 0.5);
      transform: scale(1.05);
    }
    
    &:active {
      transform: scale(0.95);
    }
    
    .t-icon {
      font-size: 20px;
    }
  }
}

.floating-deploy {
  :deep(.t-button) {
    width: 48px !important;
    height: 48px !important;
    min-width: 48px !important;
    padding: 0 !important;
    border-radius: 50% !important;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    background: #fff !important;
    border: 1px solid #ddd !important;
    display: flex;
    align-items: center;
    justify-content: center;
    
    &:hover {
      box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
      transform: scale(1.05);
      background: #f5f5f5 !important;
    }
    
    &:active {
      transform: scale(0.95);
    }
    
    .t-icon {
      font-size: 20px;
      color: #666 !important;
    }
  }
  
  .deploy-icon {
    font-size: 18px;
    font-weight: 600;
    color: #666;
  }
}

// 创建空间弹窗样式
.create-space-content {
  .form-item {
    margin-bottom: 20px;
    
    &:last-child {
      margin-bottom: 0;
    }
    
    .form-label {
      font-size: 14px;
      color: #333333;
      margin-bottom: 8px;
    }
    
    .logo-upload {
      .upload-area {
        width: 120px;
        height: 120px;
        background: #f5f5f5;
        border: 1px dashed #d9d9d9;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        transition: all 0.3s;
        position: relative;
        overflow: hidden;
        
        &:hover {
          border-color: #4a7ff7;
          background: #f0f5ff;
        }
        
        &.has-logo {
          border-style: solid;
          background: transparent;
          
          &:hover {
            border-color: #4a7ff7;
            
            &::after {
              content: '点击更换';
              position: absolute;
              top: 0;
              left: 0;
              right: 0;
              bottom: 0;
              background: rgba(0, 0, 0, 0.5);
              color: #ffffff;
              display: flex;
              align-items: center;
              justify-content: center;
              font-size: 14px;
            }
          }
        }
        
        .logo-preview {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }
        
        :deep(.t-icon) {
          color: #999999;
        }
      }
    }
  }
}
</style>

<style lang="less">
// 空间下拉菜单样式（全局样式）
.space-dropdown-content {
  padding: 8px;
  min-width: 280px;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  
  .menu-title {
    padding: 8px 12px;
    font-size: 13px;
    color: #999999;
    font-weight: 500;
  }
  
  .space-list {
    margin-bottom: 8px;
    
    .space-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 10px 12px;
      border-radius: 6px;
      cursor: pointer;
      transition: all 0.2s;
      
      &:hover {
        background: #f5f5f5;
      }
      
      &.active {
        background: #e8f0ff;
        
        .space-name {
          color: #4a7ff7;
          font-weight: 500;
        }
      }
      
      .space-icon {
        width: 32px;
        height: 32px;
        border-radius: 6px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #ffffff;
        flex-shrink: 0;
        
        .t-icon {
          font-size: 18px;
        }
      }
      
      .space-logo {
        width: 32px;
        height: 32px;
        border-radius: 6px;
        overflow: hidden;
        flex-shrink: 0;
        
        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }
      }
      
      .space-name {
        flex: 1;
        font-size: 14px;
        color: #333333;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
  
  .create-space {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 12px;
    border-top: 1px solid #eeeeee;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 14px;
    color: #4a7ff7;
    
    .t-icon {
      font-size: 16px;
    }
    
    &:hover {
      background: #f0f5ff;
    }
  }
}

// 新建表格弹窗样式
.create-dialog-content {
  padding: 0;
  
  // 新建选项
  .create-options {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 12px;
    padding: 20px;
    border-bottom: 1px solid #e8e8e8;
    
    .option-card {
      padding: 16px;
      border: 1px solid #e8e8e8;
      border-radius: 8px;
      cursor: pointer;
      transition: all 0.3s;
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 12px;
      
      &:hover {
        border-color: #4a7ff7;
        background: #f5f9ff;
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(74, 127, 247, 0.15);
      }
      
      .option-icon {
        width: 48px;
        height: 48px;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
      }
      
      .option-text {
        text-align: center;
        
        .option-title {
          font-size: 14px;
          font-weight: 500;
          color: #333;
          margin-bottom: 4px;
        }
        
        .option-desc {
          font-size: 12px;
          color: #999;
        }
      }
    }
  }
  
  // 模板区域
  .template-section {
    padding: 20px;
    
    .section-title {
      font-size: 16px;
      font-weight: 500;
      color: #333;
      margin: 0 0 16px;
    }
    
    // 模板分类
    .template-tabs {
      display: flex;
      gap: 8px;
      margin-bottom: 16px;
      overflow-x: auto;
      
      &::-webkit-scrollbar {
        height: 4px;
      }
      
      &::-webkit-scrollbar-thumb {
        background: #ddd;
        border-radius: 2px;
      }
      
      .template-tab {
        padding: 6px 16px;
        font-size: 13px;
        color: #666;
        background: #f5f5f5;
        border-radius: 16px;
        cursor: pointer;
        white-space: nowrap;
        transition: all 0.3s;
        
        &:hover {
          background: #e8e8e8;
          color: #333;
        }
        
        &.active {
          background: #e6f0ff;
          color: #4a7ff7;
          font-weight: 500;
        }
      }
    }
    
    // 模板列表
    .template-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 16px;
      
      .template-card {
        border: 1px solid #e8e8e8;
        border-radius: 8px;
        overflow: hidden;
        cursor: pointer;
        transition: all 0.3s;
        background: #fff;
        
        &:hover {
          border-color: #4a7ff7;
          box-shadow: 0 4px 12px rgba(74, 127, 247, 0.15);
          transform: translateY(-2px);
          
          .template-image img {
            transform: scale(1.05);
          }
        }
        
        .template-image {
          width: 100%;
          height: 120px;
          background: #f5f5f5;
          display: flex;
          align-items: center;
          justify-content: center;
          overflow: hidden;
          
          &.empty {
            flex-direction: column;
            gap: 8px;
            
            .empty-text {
              font-size: 12px;
              color: #999;
            }
          }
          
          img {
            width: 100%;
            height: 100%;
            object-fit: cover;
            transition: transform 0.3s ease;
          }
        }
        
        .template-info {
          padding: 10px 12px;
          
          .template-header {
            margin-bottom: 6px;
            
            .template-name {
              font-size: 13px;
              font-weight: 500;
              color: #1f2329;
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
              line-height: 1.4;
            }
          }
          
          .template-desc {
            font-size: 12px;
            color: #646a73;
            line-height: 1.4;
            margin-bottom: 8px;
            height: 32px;
            overflow: hidden;
            text-overflow: ellipsis;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
          }
          
          .template-footer {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding-top: 8px;
            border-top: 1px solid #f0f0f0;
            
            .creator-info {
              display: flex;
              align-items: center;
              gap: 4px;
              flex: 1;
              min-width: 0;
              
              .creator-name {
                font-size: 11px;
                color: #646a73;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
              }
            }
            
            .usage-info {
              display: flex;
              align-items: center;
              gap: 3px;
              font-size: 11px;
              color: #8f959e;
              flex-shrink: 0;
            }
          }
        }
      }
    }
  }
}

// 私有化部署弹窗样式
.deploy-content {
  .qrcode-section {
    display: flex;
    justify-content: center;
    gap: 40px;
    padding: 20px 0;
  }

  .qrcode-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
  }

  .qrcode-img {
    width: 200px;
    height: 200px;
    border-radius: 8px;
    border: 1px solid #f0f0f0;
    object-fit: contain;
  }

  .qrcode-label {
    font-size: 14px;
    font-weight: 500;
    color: #333;
    text-align: center;
  }
}
</style>
