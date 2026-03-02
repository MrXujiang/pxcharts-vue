<template>
  <div class="team-page">
    <div class="page-header">
      <h2 class="page-title">团队管理</h2>
      <t-button theme="primary" @click="handleCreateTeam">
        <template #icon>
          <t-icon name="add" />
        </template>
        创建团队
      </t-button>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <t-loading text="加载中..." />
    </div>

    <!-- 空状态 -->
    <div v-else-if="teamList.length === 0" class="empty-container">
      <div class="empty-content">
        <t-icon name="usergroup" size="64px" class="empty-icon" />
        <p class="empty-text">暂无团队</p>
        <p class="empty-desc">您还没有创建或加入任何团队</p>
        <t-button theme="primary" @click="handleCreateTeam">
          <template #icon>
            <t-icon name="add" />
          </template>
          创建团队
        </t-button>
      </div>
    </div>

    <!-- 团队列表 -->
    <div v-else class="team-list">
      <t-table
        :data="teamList"
        :columns="teamColumns"
        row-key="id"
        hover
      >
        <template #name="{ row }">
          <div class="team-name-cell">
            <!-- 优先显示logo图片，如果没有则显示图标 -->
            <div v-if="row.logo" class="team-logo">
              <img :src="row.logo" alt="团队logo" />
            </div>
            <div v-else class="team-icon" :style="{ background: row.color || getRandomColor() }">
              {{ getTeamIcon(row.name) }}
            </div>
            <span>{{ row.name }}</span>
            <t-tag v-if="row.isExternal" theme="warning" variant="light" size="small">外部</t-tag>
          </div>
        </template>
        <template #members="{ row }">
          <div class="members-cell">
            <t-avatar-group v-if="row.displayAvatarList && row.displayAvatarList.length > 0" :max="3" size="small">
              <t-avatar v-for="(avatar, index) in row.displayAvatarList" :key="index" :image="avatar" />
            </t-avatar-group>
            <span class="member-count">共{{ row.memberCount || 0 }}人</span>
          </div>
        </template>
        <template #creator="{ row }">
          <span>{{ row.creator || '暂无' }}</span>
        </template>
        <template #createdAt="{ row }">
          <span>{{ formatTime(row.createdAt) || '暂无' }}</span>
        </template>
        <template #updatedAt="{ row }">
          <span>{{ formatTime(row.updatedAt) || '暂无' }}</span>
        </template>
        <template #action="{ row }">
          <t-dropdown trigger="click" :min-column-width="160">
            <t-button variant="text" shape="square">
              <t-icon name="more" />
            </t-button>
            <t-dropdown-menu>
              <t-dropdown-item @click="handleTeamSettings(row)">
                <div class="dropdown-item-content">
                  <t-icon name="usergroup" />
                  <span>团队设置</span>
                </div>
              </t-dropdown-item>
              <t-dropdown-item @click="handleMemberManage(row)">
                <div class="dropdown-item-content">
                  <t-icon name="usergroup-add" />
                  <span>成员管理</span>
                </div>
              </t-dropdown-item>
              <t-dropdown-item @click="handleLeaveTeam(row)">
                <div class="dropdown-item-content">
                  <t-icon name="logout" />
                  <span>离开团队</span>
                </div>
              </t-dropdown-item>
              <t-dropdown-item class="danger-item" @click="handleDeleteTeam(row)">
                <div class="dropdown-item-content danger">
                  <t-icon name="delete" />
                  <span>删除团队</span>
                </div>
              </t-dropdown-item>
            </t-dropdown-menu>
          </t-dropdown>
        </template>
      </t-table>
    </div>

    <!-- 创建团队弹窗 -->
    <t-dialog
      v-model:visible="showCreateTeam"
      header="创建团队"
      width="600px"
      :footer="true"
      :on-confirm="handleConfirmCreateTeam"
      :on-cancel="() => showCreateTeam = false"
    >
      <div class="team-settings-content">
        <div class="form-item">
          <div class="form-label">团队名称：</div>
          <t-input v-model="newTeam.name" placeholder="请输入团队名称" />
        </div>
        <div class="form-item">
          <div class="form-label">团队描述：</div>
          <t-textarea
            v-model="newTeam.description"
            placeholder="请输入团队描述"
            :autosize="{ minRows: 3, maxRows: 5 }"
          />
        </div>
        <div class="form-item">
          <div class="form-label">团队Logo：</div>
          <div class="logo-upload">
            <input
              ref="newLogoInputRef"
              type="file"
              accept="image/*"
              style="display: none"
              @change="handleNewLogoChange"
            />
            <div
              class="upload-area"
              :class="{ 'has-logo': newTeam.logo }"
              @click="handleUploadNewLogo"
            >
              <img v-if="newTeam.logo" :src="newTeam.logo" alt="团队Logo" class="logo-preview" />
              <div v-else class="upload-placeholder">
                <t-icon name="add" size="32px" />
                <span>上传Logo</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </t-dialog>

    <!-- 团队设置弹窗 -->
    <t-dialog
      v-model:visible="showTeamSettings"
      header="团队设置"
      width="600px"
      :footer="true"
      :on-confirm="handleSaveTeamSettings"
      :on-cancel="() => showTeamSettings = false"
    >
      <div class="team-settings-content">
        <div class="form-item">
          <div class="form-label">团队名称：</div>
          <t-input v-model="currentTeam.name" placeholder="请输入团队名称" />
        </div>
        <div class="form-item">
          <div class="form-label">团队描述：</div>
          <t-textarea
            v-model="currentTeam.description"
            placeholder="请输入团队描述"
            :autosize="{ minRows: 3, maxRows: 5 }"
          />
        </div>
        <div class="form-item">
          <div class="form-label">团队Logo：</div>
          <div class="logo-upload">
            <input
              ref="logoInputRef"
              type="file"
              accept="image/*"
              style="display: none"
              @change="handleLogoChange"
            />
            <div
              class="upload-area"
              :class="{ 'has-logo': currentTeam.logo }"
              @click="handleUploadLogo"
            >
              <img v-if="currentTeam.logo" :src="currentTeam.logo" alt="团队Logo" class="logo-preview" />
              <t-icon v-else name="add" size="32px" />
            </div>
          </div>
        </div>
      </div>
    </t-dialog>

    <!-- 成员管理弹窗 -->
    <t-dialog
      v-model:visible="showMemberManage"
      header="成员管理"
      width="600px"
      :footer="false"
    >
      <div class="member-manage-content">
        <div class="manage-header">
          <t-input
            v-model="memberSearchKeyword"
            placeholder="搜索成员"
            clearable
          >
            <template #prefix-icon>
              <t-icon name="search" />
            </template>
          </t-input>
          <t-button theme="primary" @click="showInviteDialog = true">
            邀请成员
          </t-button>
        </div>
        <div class="member-list">
          <!-- 加载状态 -->
          <div v-if="memberList.length === 0" class="empty-members">
            <t-icon name="user" size="48px" class="empty-icon" />
            <p class="empty-text">暂无成员</p>
          </div>
          <div
            v-for="member in filteredMembers"
            :key="member.id"
            class="member-item"
          >
            <div class="member-info">
              <t-avatar v-if="member.avatar" :image="member.avatar" size="40px" />
              <t-avatar v-else size="40px">{{ member.name?.charAt(0) || 'U' }}</t-avatar>
              <span class="member-name">{{ member.name || '未知用户' }}</span>
            </div>
            <t-button
              theme="danger"
              variant="outline"
              size="small"
              @click="handleRemoveMember(member)"
            >
              移除
            </t-button>
          </div>
        </div>
      </div>
    </t-dialog>

    <!-- 邀请成员弹窗 -->
    <t-dialog
      v-model:visible="showInviteDialog"
      header="邀请成员加入"
      width="600px"
      :footer="false"
    >
      <div class="invite-content">
        <div class="invite-search">
          <t-input
            v-model="inviteSearchKeyword"
            placeholder="搜索成员"
            clearable
          >
            <template #prefix-icon>
              <t-icon name="search" />
            </template>
          </t-input>
          <t-button theme="primary" @click="handleInviteMember">
            邀请
          </t-button>
        </div>
        <div class="invite-list">
          <div
            v-for="user in searchedUsers"
            :key="user.id"
            class="invite-item"
          >
            <t-avatar size="40px">{{ user.name.charAt(0) }}</t-avatar>
            <span class="user-name">{{ user.name }}</span>
          </div>
        </div>
        <div class="invite-actions">
          <div class="action-item" @click="handleCopyInviteLink">
            <t-icon name="link" />
            <span>复制链接</span>
          </div>
          <div class="action-item" @click="handleShowQRCode">
            <t-icon name="qrcode" />
            <span>QR Code</span>
          </div>
        </div>
      </div>
    </t-dialog>

    <!-- 删除团队确认弹窗 -->
    <t-dialog
      v-model:visible="showDeleteConfirm"
      header="删除团队"
      width="480px"
      :on-confirm="confirmDeleteTeam"
      :on-cancel="() => showDeleteConfirm = false"
    >
      <div class="confirm-content">
        <p>确定要删除团队「<strong>{{ pendingTeam?.name }}</strong>」吗？</p>
        <p class="warning-text">删除后将无法恢复，请谨慎操作。</p>
      </div>
    </t-dialog>

    <!-- 离开团队确认弹窗 -->
    <t-dialog
      v-model:visible="showLeaveConfirm"
      header="离开团队"
      width="480px"
      :on-confirm="confirmLeaveTeam"
      :on-cancel="() => showLeaveConfirm = false"
    >
      <div class="confirm-content">
        <p>确定要离开团队「<strong>{{ pendingTeam?.name }}</strong>」吗？</p>
        <p class="warning-text">离开后将不再显示该团队的相关信息。</p>
      </div>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { getTeamList, createTeam, deleteTeam, updateTeam, addTeamMember, removeTeamMember, getTeamMemberList } from '@/api'
import { uploadFile } from '@/api'
import { useTeamStore } from '@/stores'

// 使用团队 store
const teamStore = useTeamStore()

// 加载状态（使用 store 的 loading）
const loading = computed(() => teamStore.loading)

// 弹窗显示状态
const showCreateTeam = ref(false)
const showTeamSettings = ref(false)
const showMemberManage = ref(false)
const showInviteDialog = ref(false)
const showDeleteConfirm = ref(false)
const showLeaveConfirm = ref(false)

// 新团队数据
const newTeam = ref({
  name: '',
  description: '',
  logo: '',
})

// 当前操作的团队
const currentTeam = ref<any>({
  id: '',
  name: '',
  description: '',
  logo: '',
})

// 待处理的团队（用于确认操作）
const pendingTeam = ref<any>(null)

// Logo上传输入框引用
const logoInputRef = ref<HTMLInputElement>()
const newLogoInputRef = ref<HTMLInputElement>()

// 搜索关键词
const memberSearchKeyword = ref('')
const inviteSearchKeyword = ref('')

// 成员列表
const memberList = ref<any[]>([])

// 模拟可邀请用户数据
const searchedUsers = ref([
  { id: '11', name: '小梦', avatar: '' },
])

// 过滤后的成员列表
const filteredMembers = computed(() => {
  if (!memberSearchKeyword.value) {
    return memberList.value
  }
  return memberList.value.filter(member =>
    member.name.toLowerCase().includes(memberSearchKeyword.value.toLowerCase())
  )
})

// 团队列表（使用 store）
const teamList = computed(() => teamStore.teamList)

const teamColumns = [
  {
    colKey: 'name',
    title: '团队名称',
    width: 400,
  },
  {
    colKey: 'creator',
    title: '创建人',
    width: 150,
  },
  {
    colKey: 'members',
    title: '成员',
    width: 200,
  },
  {
    colKey: 'createdAt',
    title: '创建时间',
    width: 200,
  },
  {
    colKey: 'updatedAt',
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

// 格式化时间
const formatTime = (time: string) => {
  if (!time) return ''
  
  const date = new Date(time)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

// 获取随机颜色
const getRandomColor = () => {
  const colors = ['#d3adf7', '#b37feb', '#597ef7', '#52c41a', '#fa8c16', '#4a7ff7', '#f759ab', '#ffc53d']
  return colors[Math.floor(Math.random() * colors.length)]
}

// 获取团队图标（取名称前两个字）
const getTeamIcon = (name: string) => {
  if (!name) return 'T'
  return name.substring(0, 2)
}

// 加载团队列表
const loadTeamList = async () => {
  try {
    await teamStore.loadTeamList()
  } catch (error: any) {
    console.error('加载团队列表失败：', error)
    MessagePlugin.error(error.message || '加载团队列表失败')
  }
}

// 创建团队
const handleCreateTeam = () => {
  newTeam.value = {
    name: '',
    description: '',
    logo: '',
  }
  showCreateTeam.value = true
}

// 上传新团队Logo
const handleUploadNewLogo = () => {
  newLogoInputRef.value?.click()
}

// 新Logo文件变化处理
const handleNewLogoChange = async (event: Event) => {
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
      const response = await uploadFile(file) as any
      if (response && response.url) {
        newTeam.value.logo = response.url
        MessagePlugin.success('Logo上传成功')
      }
    } catch (error: any) {
      MessagePlugin.error(error.message || 'Logo上传失败')
    }
  }
}

// 确认创建团队
const handleConfirmCreateTeam = async () => {
  if (!newTeam.value.name) {
    MessagePlugin.warning('请输入团队名称')
    return
  }
  
  try {
    const response = await createTeam(newTeam.value) as any
    console.log('创建团队响应：', response)
    MessagePlugin.success('团队创建成功')
    showCreateTeam.value = false
    // 重新加载团队列表
    loadTeamList()
  } catch (error: any) {
    console.error('创建团队失败：', error)
    MessagePlugin.error(error.message || '创建团队失败')
  }
}

// 团队设置
const handleTeamSettings = (team: any) => {
  currentTeam.value = { 
    id: team.id,
    name: team.name,
    description: team.description || '',
    logo: team.logo || '',
  }
  showTeamSettings.value = true
}

// 保存团队设置
const handleSaveTeamSettings = async () => {
  if (!currentTeam.value.name) {
    MessagePlugin.warning('请输入团队名称')
    return
  }
  
  try {
    await updateTeam({
      id: currentTeam.value.id,
      name: currentTeam.value.name,
      description: currentTeam.value.description,
      logo: currentTeam.value.logo,
    })
    MessagePlugin.success('团队设置已保存')
    showTeamSettings.value = false
    // 重新加载团队列表
    loadTeamList()
  } catch (error: any) {
    console.error('保存团队设置失败：', error)
    MessagePlugin.error(error.message || '保存失败')
  }
}

// 上传Logo
const handleUploadLogo = () => {
  logoInputRef.value?.click()
}

// Logo文件变化处理
const handleLogoChange = async (event: Event) => {
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
      const response = await uploadFile(file) as any
      if (response && response.url) {
        currentTeam.value.logo = response.url
        MessagePlugin.success('Logo已更新')
      }
    } catch (error: any) {
      MessagePlugin.error(error.message || 'Logo上传失败')
    }
  }
}

// 成员管理
const handleMemberManage = async (team: any) => {
  currentTeam.value = team
  showMemberManage.value = true
  // 加载团队成员列表
  await loadTeamMembers(team.id)
}

// 加载团队成员列表
const loadTeamMembers = async (teamId: string) => {
  try {
    const response = await getTeamMemberList({ teamId }) as any
    console.log('团队成员列表响应：', response)
    
    // 处理响应数据
    if (response && Array.isArray(response)) {
      memberList.value = response
    } else if (response && response.list && Array.isArray(response.list)) {
      memberList.value = response.list
    } else {
      memberList.value = []
    }
  } catch (error: any) {
    console.error('加载团队成员失败：', error)
    MessagePlugin.error(error.message || '加载成员列表失败')
    memberList.value = []
  }
}

// 移除成员
const handleRemoveMember = async (member: any) => {
  try {
    await removeTeamMember({
      teamId: currentTeam.value.id,
      userId: member.id,
    })
    MessagePlugin.success(`已移除成员：${member.name}`)
    const index = memberList.value.findIndex(m => m.id === member.id)
    if (index > -1) {
      memberList.value.splice(index, 1)
    }
  } catch (error: any) {
    console.error('移除成员失败：', error)
    MessagePlugin.error(error.message || '移除成员失败')
  }
}

// 邀请成员
const handleInviteMember = async () => {
  // TODO: 实现邀请逻辑
  MessagePlugin.success('邀请已发送')
  showInviteDialog.value = false
}

// 复制邀请链接
const handleCopyInviteLink = () => {
  MessagePlugin.success('链接已复制到剪贴板')
}

// 显示二维码
const handleShowQRCode = () => {
  MessagePlugin.info('二维码功能开发中')
}

// 离开团队
const handleLeaveTeam = (team: any) => {
  pendingTeam.value = team
  showLeaveConfirm.value = true
}

// 确认离开团队
const confirmLeaveTeam = () => {
  if (pendingTeam.value) {
    const index = teamList.value.findIndex(t => t.id === pendingTeam.value.id)
    if (index > -1) {
      teamList.value.splice(index, 1)
      MessagePlugin.success(`已离开团队：${pendingTeam.value.name}`)
    }
  }
  showLeaveConfirm.value = false
  pendingTeam.value = null
}

// 删除团队
const handleDeleteTeam = (team: any) => {
  pendingTeam.value = team
  showDeleteConfirm.value = true
}

// 确认删除团队
const confirmDeleteTeam = async () => {
  if (pendingTeam.value) {
    try {
      await deleteTeam({ id: pendingTeam.value.id })
      const index = teamList.value.findIndex(t => t.id === pendingTeam.value.id)
      if (index > -1) {
        teamList.value.splice(index, 1)
      }
      MessagePlugin.success(`已删除团队：${pendingTeam.value.name}`)
    } catch (error: any) {
      console.error('删除团队失败：', error)
      MessagePlugin.error(error.message || '删除团队失败')
    }
  }
  showDeleteConfirm.value = false
  pendingTeam.value = null
}

// 组件挂载时加载数据
onMounted(() => {
  loadTeamList()
})
</script>

<style scoped lang="less">
.team-page {
  min-height: 100%;
  background: #ffffff;
  border-radius: 8px;
  padding: 24px;
  
  .page-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 20px;
  }
  
  .page-title {
    font-size: 20px;
    font-weight: 500;
    color: #333;
    margin: 0;
  }
  
  // 加载状态
  .loading-container {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 400px;
  }
  
  // 空状态
  .empty-container {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 400px;
    
    .empty-content {
      text-align: center;
      
      .empty-icon {
        color: #dcdcdc;
        margin-bottom: 16px;
      }
      
      .empty-text {
        font-size: 16px;
        color: #333;
        margin: 0 0 8px;
        font-weight: 500;
      }
      
      .empty-desc {
        font-size: 14px;
        color: #999;
        margin: 0 0 24px;
      }
    }
  }
  
  .team-list {
    background: #ffffff;
    border-radius: 8px;
    overflow: hidden;
    
    .team-name-cell {
      display: flex;
      align-items: center;
      gap: 12px;
      
      .team-logo {
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
      
      .team-icon {
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
    
    .members-cell {
      display: flex;
      align-items: center;
      gap: 12px;
      
      .member-count {
        font-size: 14px;
        color: #666;
      }
    }
  }
}

// 团队设置弹窗样式
.team-settings-content {
  .form-item {
    margin-bottom: 24px;
    
    &:last-child {
      margin-bottom: 0;
    }
    
    .form-label {
      font-size: 14px;
      color: #333333;
      margin-bottom: 12px;
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
        
        .upload-placeholder {
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: 8px;
          
          span {
            font-size: 12px;
            color: #999;
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

// 确认弹窗样式
.confirm-content {
  padding: 10px 0;
  
  p {
    margin: 0 0 12px;
    font-size: 14px;
    color: #333333;
    line-height: 1.6;
    
    &:last-child {
      margin-bottom: 0;
    }
    
    strong {
      color: #4a7ff7;
      font-weight: 500;
    }
  }
  
  .warning-text {
    color: #e34d59;
    font-size: 13px;
  }
}

// 成员管理弹窗样式
.member-manage-content {
  .manage-header {
    display: flex;
    gap: 12px;
    margin-bottom: 20px;
  }
  
  .member-list {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
    max-height: 400px;
    overflow-y: auto;
    
    .empty-members {
      grid-column: 1 / -1;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 40px 20px;
      
      .empty-icon {
        color: #dcdcdc;
        margin-bottom: 12px;
      }
      
      .empty-text {
        font-size: 14px;
        color: #999;
        margin: 0;
      }
    }
    
    .member-item {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 12px;
      border-radius: 8px;
      background: #f5f5f5;
      transition: all 0.2s;
      
      &:hover {
        background: #f0f0f0;
      }
      
      .member-info {
        display: flex;
        align-items: center;
        gap: 12px;
        
        .member-name {
          font-size: 14px;
          color: #333333;
        }
      }
    }
  }
}

// 邀请成员弹窗样式
.invite-content {
  .invite-search {
    display: flex;
    gap: 12px;
    margin-bottom: 20px;
  }
  
  .invite-list {
    margin-bottom: 24px;
    
    .invite-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 12px;
      border-radius: 8px;
      background: #f5f5f5;
      margin-bottom: 8px;
      
      .user-name {
        font-size: 14px;
        color: #333333;
      }
    }
  }
  
  .invite-actions {
    display: flex;
    gap: 12px;
    padding-top: 16px;
    border-top: 1px solid #eeeeee;
    
    .action-item {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 8px 16px;
      border-radius: 6px;
      background: #f5f5f5;
      cursor: pointer;
      transition: all 0.2s;
      font-size: 14px;
      color: #333333;
      
      :deep(.t-icon) {
        font-size: 16px;
        color: #666666;
      }
      
      &:hover {
        background: #e6e6e6;
      }
    }
  }
}
</style>

<style lang="less">
// 下拉菜单项样式（全局样式）
.dropdown-item-content {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 2px 0;
  font-size: 14px;
  color: #333333;
  
  .t-icon {
    font-size: 16px;
    color: #666666;
  }
  
  &.danger {
    color: #e34d59;
    
    .t-icon {
      color: #e34d59;
    }
  }
}

// 危险操作项hover样式
.t-dropdown__item.danger-item {
  &:hover {
    background: #fff1f0 !important;
  }
}
</style>
