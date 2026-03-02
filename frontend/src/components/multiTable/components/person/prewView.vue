<template>
  <div class="person-preview-container flx-ce-sta">
    <!-- 如果是多选，循环显示所有人员 -->
    <template v-if="isMultiple">
      <t-popup v-for="person in selectedPersons" :key="person.id">
        <t-tag size="small" shape="round">
          <template #content>
            <div class="flx-ce-sta gap-4">
              <img class="img" :src="person.avatar || defaultAvatar" alt="" />
              <div style="font-size: 10px">{{ person.nickname }}</div>
              <!-- <span class="external" v-if="person.isExternal">外部</span> -->
            </div>
          </template>
        </t-tag>
        <template #content>
          <div class="flex-col gap-4 popup-content">
            <t-avatar size="42px" :image="person.avatar || defaultAvatar" />
            <p style="font-size: 16px; font-weight: 500">{{ person.nickname }}</p>
            <div>
              <t-button variant="outline">@ 提醒他关注记录</t-button>
            </div>
          </div>
        </template>
      </t-popup>
    </template>
    <!-- 如果是单选，只显示一个人员 -->
    <template v-else>
      <t-popup v-if="selectedPerson">
        <t-tag size="small" shape="round">
          <template #content>
            <div class="flx-ce-sta gap-4">
              <img class="img" :src="selectedPerson.avatar || defaultAvatar" alt="" />
              <div style="font-size: 10px">{{ selectedPerson.nickname }}</div>
              <span class="external" v-if="selectedPerson.isExternal">外部</span>
            </div>
          </template>
        </t-tag>
        <template #content>
          <div class="flex-col gap-4 popup-content">
            <t-avatar size="42px" :image="selectedPerson.avatar || defaultAvatar" />
            <p style="font-size: 16px; font-weight: 500">{{ selectedPerson.nickname }}</p>
            <div>
              <t-button variant="outline">@ 提醒他关注记录</t-button>
            </div>
          </div>
        </template>
      </t-popup>
      <!-- 默认人员 -->
      <!-- <t-popup v-else>
        <t-tag size="small" shape="round">
          <template #content>
            <div class="flx-ce-sta gap-4">
              <img class="img" :src="defaultPerson.avatar || defaultAvatar" alt="" />
              <div style="font-size: 10px">{{ defaultPerson.nickname }}</div>
            </div>
          </template>
        </t-tag>
        <template #content>
          <div class="flex-col gap-4 popup-content">
            <t-avatar size="42px" :image="defaultPerson.avatar || defaultAvatar" />
            <p style="font-size: 16px; font-weight: 500">{{ defaultPerson.nickname }}</p>
            <div>
              <t-button variant="outline">@ 提醒他关注记录</t-button>
            </div>
          </div>
        </template>
      </t-popup> -->
    </template>
  </div>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'
import { batchGetUsers } from '@/api'

interface Props {
  editingValue: ObjType[]
  field: ObjType
}
const props = defineProps<Props>()
// 默认头像
const defaultAvatar = 'https://tdesign.gtimg.com/site/avatar.jpg'
const selectedPersons = ref<ObjType[]>([])

// 获取字段配置，判断是否允许多选
const fieldConfig = computed(() => props.field)

// 判断是否是多选
const isMultiple = computed(() => {
  return fieldConfig.value?.settings?.allowMultiple || false
})
// 获取默认人员（索引为0的人员）
// const defaultPerson = computed(() => {})

// 获取人员
const fetchPersons = async (ids: string[]) => {
  if (ids.length == 0) {
    selectedPersons.value = []
    return
  }

  try {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const res: any = await batchGetUsers({ userIds: ids })
    selectedPersons.value = res.list || []
  } catch (error) {
    console.error('Failed to fetch persons:', error)
    selectedPersons.value = []
  }
}

// 获取选中的单个人员
const selectedPerson = computed(() => {
  return selectedPersons.value[0] || null
})

// 监听 editingValue 变化，自动重新拉取
watch(
  () => props.editingValue,
  (newVal) => {
    if (Array.isArray(newVal)) {
      fetchPersons(newVal as unknown as string[])
    } else {
      selectedPersons.value = []
    }
  },
  { immediate: true },
)
defineOptions({ name: 'PersonPreview' })
</script>

<style lang="less" scoped>
.person-preview-container {
  width: 100%;
  overflow-x: auto;
  scrollbar-width: none;
  &::-webkit-scrollbar {
    display: none;
  }
  &::-webkit-scrollbar {
    display: none;
  }
  -ms-overflow-style: none;
}

.img {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  object-fit: cover;
}

.popup-content {
  width: 320px;
  padding: 20px 12px;
  align-items: center;
}

.external {
  color: orange;
  font-size: 10px;
}
</style>
