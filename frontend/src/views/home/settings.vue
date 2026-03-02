<template>
  <div class="settings-page">
    <div class="settings-title">个人设置</div>

    <div class="settings-content">
      <!-- 头像设置 -->
      <div class="settings-item avatar-item">
        <div class="item-label">头像：</div>
        <div class="item-value">
          <div class="avatar-wrapper">
            <img :src="formData.avatar || defaultAvatar" alt="用户头像" class="user-avatar" />
            <input
              ref="fileInputRef"
              type="file"
              accept="image/*"
              style="display: none"
              @change="handleFileChange"
            />
            <t-button variant="text" theme="primary" @click="handleChangeAvatar">
              更换
            </t-button>
          </div>
        </div>
      </div>

      <!-- 用户名设置 -->
      <div class="settings-item">
        <div class="item-label">用户名：</div>
        <div class="item-value">
          <t-input
            v-model="formData.username"
            placeholder="请输入用户名"
            class="setting-input"
          />
        </div>
      </div>

      <!-- 邮箱设置 -->
      <div class="settings-item">
        <div class="item-label">邮箱：</div>
        <div class="item-value">
          <t-input
            v-model="formData.email"
            placeholder="请输入邮箱"
            class="setting-input"
            disabled
          />
        </div>
      </div>

      <!-- 兴趣标签 -->
      <div class="settings-item tags-item">
        <div class="item-label">兴趣标签：</div>
        <div class="item-value">
          <div class="tags-select-container">
            <div 
              v-for="tag in availableTags" 
              :key="tag.value"
              class="tag-item"
              :class="{ active: formData.tags?.includes(tag.value) }"
              @click="toggleTag(tag.value)"
            >
              {{ tag.label }}
            </div>
          </div>
        </div>
      </div>

      <!-- 用户权益 -->
      <div class="settings-item">
        <div class="item-label">用户权益：</div>
        <div class="item-value">
          <div class="user-level">
            <span class="level-badge">基础版（免费）</span>
            <t-button variant="text" theme="primary" @click="handleUpgrade">
              升级
            </t-button>
          </div>
        </div>
      </div>

      <!-- 保存按钮 -->
      <div class="settings-actions">
        <t-button theme="primary" :loading="loading" @click="handleSave">保存</t-button>
      </div>
    </div>

    <!-- 升级弹窗 -->
    <t-dialog
      v-model:visible="showUpgradeDialog"
      header="用户权益升级"
      width="720px"
      :footer="false"
      @close="handleCloseUpgrade"
    >
      <div class="upgrade-content">
        <div class="plan-list">
          <!-- 基础版 -->
          <div class="plan-item">
            <div class="plan-header">
              <div class="plan-title">基础版（免费）</div>
            </div>
            <div class="plan-features">
              <div class="feature-item">
                <span class="bullet">•</span>
                <span>上限100个表格 / 表单 / 可视化大屏</span>
              </div>
              <div class="feature-item">
                <span class="bullet">•</span>
                <span>表格最大行数2000行</span>
              </div>
              <div class="feature-item">
                <span class="bullet">•</span>
                <span>AI使用token量1w</span>
              </div>
              <div class="feature-item">
                <span class="bullet">•</span>
                <span>团队最大人数10人</span>
              </div>
            </div>
          </div>

          <!-- 专业版 -->
          <div class="plan-item">
            <div class="plan-header">
              <div class="plan-title">专业版（199元/年）</div>
            </div>
            <div class="plan-features">
              <div class="feature-item">
                <span class="bullet">•</span>
                <span>上限300个表格 / 表单 / 可视化大屏</span>
              </div>
              <div class="feature-item">
                <span class="bullet">•</span>
                <span>表格最大行数10000行</span>
              </div>
              <div class="feature-item">
                <span class="bullet">•</span>
                <span>AI使用token量10w</span>
              </div>
              <div class="feature-item">
                <span class="bullet">•</span>
                <span>团队最大人数100人</span>
              </div>
            </div>
          </div>

          <!-- 旗舰版 -->
          <div class="plan-item">
            <div class="plan-header">
              <div class="plan-title">旗舰版（999元/年）</div>
            </div>
            <div class="plan-features">
              <div class="feature-item">
                <span class="bullet">•</span>
                <span>上限1000个表格 / 表单 / 可视化大屏</span>
              </div>
              <div class="feature-item">
                <span class="bullet">•</span>
                <span>表格最大行数20000行</span>
              </div>
              <div class="feature-item">
                <span class="bullet">•</span>
                <span>AI使用token量100wl</span>
              </div>
              <div class="feature-item">
                <span class="bullet">•</span>
                <span>团队最大人数1000人</span>
              </div>
            </div>
          </div>
        </div>

        <div class="upgrade-footer">
          <t-button theme="primary" size="large" @click="handlePurchase">
            添加购买升级会员权益
          </t-button>
        </div>
      </div>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { useUserStore } from '@/stores/user'
import { updateUserInfo, uploadFile } from '@/api'

// 用户 Store
const userStore = useUserStore()

// 默认头像
const defaultAvatar = 'https://tdesign.gtimg.com/site/avatar.jpg'

// 表单数据
const formData = ref({
  avatar: '',
  username: '',
  email: '',
  tags: [] as string[],
})

// 可选标签列表
const availableTags = [
  { label: '互联网/科技', value: 'internet' },
  { label: '企业服务/To B', value: 'enterprise' },
  { label: '金融/财经', value: 'finance' },
  { label: '教育/培训', value: 'education' },
  { label: '电商/零售', value: 'ecommerce' },
  { label: '新能源/汽车', value: 'newenergy' },
  { label: '项目管理', value: 'project' },
  { label: '客户/人员管理', value: 'customer' },
  { label: '运营分析', value: 'operation' },
  { label: '资产管理', value: 'asset' },
  { label: '个人/爱好', value: 'personal' },
]

// 文件上传引用
const fileInputRef = ref<HTMLInputElement>()

// 升级弹窗显示状态
const showUpgradeDialog = ref(false)

// 保存加载状态
const loading = ref(false)

// 初始化用户信息
onMounted(() => {
  const userInfo = userStore.userInfo
  if (userInfo) {
    formData.value = {
      avatar: userInfo.avatar || '',
      username: userInfo.username || '',
      email: userInfo.email || '',
      tags: userInfo.tags || [],
    }
  }
})

// 切换标签选择
const toggleTag = (value: string) => {
  if (!formData.value.tags) {
    formData.value.tags = []
  }
  const index = formData.value.tags.indexOf(value)
  if (index > -1) {
    formData.value.tags.splice(index, 1)
  } else {
    formData.value.tags.push(value)
  }
}

// 更换头像
const handleChangeAvatar = () => {
  fileInputRef.value?.click()
}

// 处理文件选择
const handleFileChange = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    // 验证文件类型
    if (!file.type.startsWith('image/')) {
      MessagePlugin.warning('请选择图片文件')
      return
    }
    // 验证文件大小（限制5MB）
    if (file.size > 5 * 1024 * 1024) {
      MessagePlugin.warning('图片大小不能超过5MB')
      return
    }
    
    // 先显示本地预览
    const reader = new FileReader()
    reader.onload = (e) => {
      formData.value.avatar = e.target?.result as string
    }
    reader.readAsDataURL(file)
    
    // 上传到服务器
    try {
      loading.value = true
      const response = await uploadFile(file)
      
      // 接口返回 { url: string }
      if (response?.url) {
        formData.value.avatar = response.url
        MessagePlugin.success('头像上传成功')
      }
    } catch (error: any) {
      console.error('头像上传失败:', error)
      MessagePlugin.error(error.message || '头像上传失败，请稍后重试')
      // 恢复原始头像
      const userInfo = userStore.userInfo
      if (userInfo) {
        formData.value.avatar = userInfo.avatar || ''
      }
    } finally {
      loading.value = false
      // 清空文件输入，允许重复选择同一文件
      if (fileInputRef.value) {
        fileInputRef.value.value = ''
      }
    }
  }
}

// 升级权益
const handleUpgrade = () => {
  showUpgradeDialog.value = true
}

// 关闭升级弹窗
const handleCloseUpgrade = () => {
  showUpgradeDialog.value = false
}

// 添加购买会员按钮点击事件
const handlePurchase = () => {
  MessagePlugin.info('跳转到支付页面')
  // TODO: 实现支付功能
}

// 保存设置
const handleSave = async () => {
  loading.value = true
  try {
    // 调用更新用户信息接口
    await updateUserInfo({
      username: formData.value.username,
      avatar: formData.value.avatar,
      tags: formData.value.tags,
    })
    
    // 更新本地 store 中的用户信息
    userStore.updateUserInfo({
      username: formData.value.username,
      avatar: formData.value.avatar,
      tags: formData.value.tags,
    })
    
    MessagePlugin.success('保存成功')
  } catch (error: any) {
    console.error('保存用户信息失败:', error)
    MessagePlugin.error(error.message || '保存失败，请稍后重试')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="less">
.settings-page {
  background: #ffffff;
  border-radius: 8px;
  padding: 27px;
  min-height: 100%;

  .settings-title {
    font-size: 18px;
    font-weight: 500;
    color: #000000;
    margin-bottom: 30px;
  }

  .settings-content {
    max-width: 600px;

    .settings-item {
      display: flex;
      align-items: center;
      margin-bottom: 24px;

      &.avatar-item {
        align-items: flex-start;
      }

      &.tags-item {
        align-items: flex-start;
      }

      .item-label {
        width: 80px;
        flex-shrink: 0;
        font-size: 14px;
        color: #000000;
        text-align: right;
        padding-right: 16px;
        line-height: 32px;
        white-space: nowrap;
      }

      .item-value {
        flex: 1;

        .avatar-wrapper {
          display: flex;
          align-items: center;
          gap: 12px;

          .user-avatar {
            width: 60px;
            height: 60px;
            border-radius: 50%;
            object-fit: cover;
          }
        }

        .tags-select-container {
          display: flex;
          flex-wrap: wrap;
          gap: 12px;
          max-width: 600px;

          .tag-item {
            padding: 6px 16px;
            background: #f5f5f5;
            border-radius: 4px;
            font-size: 13px;
            color: #666;
            cursor: pointer;
            transition: all 0.3s;
            user-select: none;

            &:hover {
              background: #e8f0ff;
              color: #4a7ff7;
            }

            &.active {
              background: #4a7ff7;
              color: #ffffff;
              font-weight: 500;
            }
          }
        }

        .setting-input {
          width: 100%;
          max-width: 400px;
          background: #f5f5f5;
          border: none;

          &.disabled-input {
            color: #999;
            cursor: not-allowed;
          }

          :deep(.t-input) {
            background: #f5f5f5;
            border: none;
          }

          :deep(.t-input__inner) {
            background: #f5f5f5;
            border: none;
          }

          &:disabled {
            :deep(.t-input) {
              background: #f5f5f5;
              color: #999;
              cursor: not-allowed;
            }

            :deep(.t-input__inner) {
              background: #f5f5f5;
              color: #999;
              cursor: not-allowed;
            }
          }
        }

        .user-level {
          display: flex;
          align-items: center;
          gap: 12px;
          flex-wrap: nowrap;

          .level-badge {
            display: inline-block;
            padding: 4px 16px;
            background: #f5f5f5;
            border-radius: 4px;
            font-size: 14px;
            color: #000000;
            white-space: nowrap;
          }
        }
      }
    }

    .settings-actions {
      margin-top: 40px;
      padding-left: 96px;

      :deep(.t-button) {
        min-width: 80px;
      }
    }
  }
}

// 升级弹窗样式
.upgrade-content {
  .plan-list {
    display: flex;
    gap: 20px;
    margin-bottom: 30px;

    .plan-item {
      flex: 1;
      border: 1px solid #e8e8e8;
      border-radius: 8px;
      padding: 20px;
      background: #fafafa;

      .plan-header {
        margin-bottom: 16px;
        padding-bottom: 12px;
        border-bottom: 1px solid #e8e8e8;

        .plan-title {
          font-size: 16px;
          font-weight: 600;
          color: #000000;
        }
      }

      .plan-features {
        .feature-item {
          display: flex;
          align-items: flex-start;
          margin-bottom: 12px;
          font-size: 14px;
          color: #333;
          line-height: 1.6;

          &:last-child {
            margin-bottom: 0;
          }

          .bullet {
            color: #4a7ff7;
            margin-right: 8px;
            flex-shrink: 0;
          }
        }
      }
    }
  }

  .upgrade-footer {
    text-align: center;
    padding-top: 20px;
    border-top: 1px solid #e8e8e8;

    :deep(.t-button) {
      min-width: 200px;
      height: 44px;
      font-size: 16px;
    }
  }
}
</style>
