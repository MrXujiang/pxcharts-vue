<template>
  <div class="mt-header">
    <!-- 左侧内容 -->
    <div class="header-left">
      <!-- 文本标题（双击可编辑） -->
      <div class="flx-ce-sta gap-8">
        <t-icon style="font-size: 24px" class="pointer" name="home" @click="handleHomeClick" />
        <div class="title-section" @dblclick="enableEditMode">
          <div v-if="!isEditing" class="title-text">{{ title }}</div>
          <t-input
            v-else
            ref="titleInput"
            v-model="editableTitle"
            class="title-input"
            @blur="saveTitle"
            @enter="saveTitle"
          />
        </div>
      </div>

      <!-- 收藏按钮 isFavorite -->
      <t-icon
        class="favorite-icon"
        style="cursor: pointer"
        @click="toggleFavorite"
        :name="isFavorite ? 'heartFilled' : 'heart'"
        :style="{ color: isFavorite ? 'orange' : '#333' }"
      />
    </div>

    <!-- 右侧内容 -->
    <div class="header-right">
      <!-- 保存状态 -->
      <div class="save-status">
        <t-icon name="objectStorage"></t-icon>
        <span v-if="isSaving">保存中...</span>
        <span v-else>已保存</span>
      </div>

      <!-- 分享按钮 -->
      <t-button theme="primary" size="medium" shape="round">
        <template #icon><t-icon name="share" /></template>
        分享
      </t-button>

      <!-- 高级权限按钮 -->
      <t-button shape="round" variant="outline">
        <template #icon><t-icon name="shieldError" /></template>
        高级权限
      </t-button>
      <!-- 更多 -->
      <t-dropdown
        trigger="click"
        :min-column-width="200"
        :options="headMoreOptions"
        @click="handleToolMoreClick"
      >
        <t-button shape="square" variant="text">
          <template #icon><t-icon name="more" /></template>
        </t-button>
      </t-dropdown>

      <!-- 用户头像 -->
      <t-dropdown trigger="click" :min-column-width="200" :options="headMoreOptions">
        <template #panelTopContent>
          <div class="mt-sidebar-head-user flx-ce-sta">
            <t-comment
              :avatar="userStore.userInfo?.avatar"
              author="徐小夕"
              content="重庆橙讯智科网络科技有限公司"
            />
          </div>
        </template>
        <t-avatar class="pointer" :image="userStore.userInfo?.avatar" />
      </t-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
import { DropdownProps } from 'tdesign-vue-next'
import TIcon from '../TIcon/index.vue'
import { headMoreOptions } from './data'
import { ObjType } from '@/types'
import { useUserStore } from '@/stores/user'
import { setProjectFavorite, updateProject } from '@/api'

const userStore = useUserStore()

interface Props {
  projectDetail: ObjType
}

const props = withDefaults(defineProps<Props>(), {
  projectDetail: () => ({}),
})
// 为组件设置多词名称，避免 ESLint 警告
defineOptions({
  name: 'MtHeader',
})

const route = useRoute()
const router = useRouter()

// 响应式数据
const title = ref('')
const isEditing = ref(false)
const editableTitle = ref('')
const titleInput = ref<HTMLInputElement | null>(null)
const isFavorite = ref(false)
const isSaving = ref(false)

// moer dropdown click
const handleToolMoreClick: DropdownProps['onClick'] = (data) => {
  switch (data.value) {
    case 'importExcel':
      break
    default:
      break
  }
}

const handleHomeClick = () => {
  router.push('/home')
}

const enableEditMode = () => {
  isEditing.value = true
  editableTitle.value = title.value
  nextTick(() => {
    titleInput.value?.focus()
  })
}

const saveTitle = async () => {
  if (editableTitle.value.trim()) {
    title.value = editableTitle.value.trim()
  }
  await updateProject({ id: route.params.id as string, name: title.value })
  isEditing.value = false
}

const toggleFavorite = async () => {
  const response = await setProjectFavorite({
    projectId: route.params.id as string,
    isFavorite: !isFavorite.value,
  })
  isFavorite.value = !isFavorite.value
  MessagePlugin.success(isFavorite.value ? '收藏成功' : '取消收藏成功')
  console.log('response', response)
}

watch(
  () => props.projectDetail,
  (newVal) => {
    title.value = newVal.name
    isFavorite.value = newVal.isFavorite || false
  },
  { immediate: true, deep: true },
)
</script>

<style scoped lang="less">
@import './index.less';
</style>
