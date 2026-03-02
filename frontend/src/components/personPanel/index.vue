<template>
  <t-select
    v-model="defaultPerson"
    :multiple="allowMultiple"
    style="width: 100%"
    @change="handlePersonChange"
    @blur="handleBlur"
    @enter="handleEnter"
    @keydown.esc="handleEsc"
    @popup-visible-change="handlePopupVisibleChange"
  >
    <template #valueDisplay="{ value, onClose, label }">
      <template v-if="Array.isArray(value)">
        <t-tag
          v-for="(item, index) in value"
          :key="item?.value || index"
          shape="round"
          :closable="true"
          :on-close="
            ({ e }) => {
              e.stopPropagation()
              onClose(index)
            }
          "
        >
          <template #content>
            <div class="flx-ce-sta gap-4">
              <span
                style="width: 16px; height: 16px"
                class="flex"
                v-if="!findPersonById(item?.value)?.avatar"
              >
                <t-icon name="user1" />
              </span>
              <t-image
                v-else
                :src="findPersonById(item?.value)?.avatar"
                :style="{ width: '16px', height: '16px' }"
                shape="circle"
                fit="cover"
              />
              <span>{{ item?.label }}</span>
              <span class="external" v-if="findPersonById(item?.value)?.isExternal">外部</span>
            </div>
          </template>
        </t-tag>
      </template>
      <template v-else>
        <t-tag shape="round" :closable="true" :on-close="handleClose">
          <template #content>
            <div class="flx-ce-sta gap-4">
              <span
                style="width: 16px; height: 16px"
                class="flex"
                v-if="!findPersonById(value)?.avatar"
              >
                <t-icon name="user1" />
              </span>
              <t-image
                v-else
                :src="findPersonById(value)?.avatar"
                :style="{ width: '16px', height: '16px' }"
                shape="circle"
                fit="cover"
              />
              <span>{{ label }}</span>
              <span class="external" v-if="findPersonById(value)?.isExternal">外部</span>
            </div>
          </template>
        </t-tag>
      </template>
    </template>
    <template #panelTopContent>
      <div style="padding: 6px 6px 0 6px">
        <t-input v-model="personSearch" placeholder="输入搜索内容" clearable @change="onSearch" />
      </div>
      <t-divider style="margin: 4px 0" />
    </template>
    <t-option
      style="height: 32px; line-height: 32px"
      v-for="item in personDefaultData"
      :key="item.id"
      :value="item.id"
      :label="item.nickname"
    >
      <div>
        <t-avatar size="small" v-if="!item.avatar">
          <template #icon>
            <t-icon name="user1" />
          </template>
        </t-avatar>
        <t-avatar v-else size="small" :image="item.avatar" />
        {{ item.nickname }}
      </div>
    </t-option>
  </t-select>
</template>

<script setup lang="ts">
// import { personData } from '@/modal/mock'
import { ObjType } from '@/types'
import type { SelectValue, SelectOption } from 'tdesign-vue-next'
import { searchUser, batchGetUsers } from '@/api'

interface props {
  modelValue: string | string[]
  allowMultiple: boolean
}
const props = defineProps<props>()

const emit = defineEmits(['update:modelValue', 'handleBlur', 'handleEnter', 'handleEsc'])
const defaultPerson = ref<string | string[]>(props.modelValue)
// 缓存已选中的人员完整对象
const selectedCache = ref<Map<string | number, ObjType>>(new Map())

// 存储当前API搜索回来的原始结果
const lastSearchResults = ref<ObjType[]>([])

const personSearch = ref<string>('')
const personDefaultData = ref<ObjType[]>([])

// 监听外部 modelValue 变化
watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal !== undefined && newVal !== null) {
      defaultPerson.value = newVal as string | string[]
      // 如果外部传进来的值为空，也进行重置操作
      checkAndResetState()
      refreshData()
    } else {
      defaultPerson.value = props.allowMultiple ? [] : ''
      checkAndResetState() // 空值重置
      refreshData()
    }
  },
)

// 监听搜索框输入，如果用户手动删完了搜索词，也要清空结果
watch(personSearch, (newVal) => {
  if (!newVal) {
    lastSearchResults.value = []
    refreshData()
  }
})

// 初始化获取下拉人员数据
const initPersonData = async (userIds: string[] = []) => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const res: any = await batchGetUsers({ userIds })
  personDefaultData.value = res.list || []
}

// 监听下拉弹窗显示隐藏
const handlePopupVisibleChange = (visible: boolean) => {
  if (!visible) {
    personSearch.value = ''
  }
}

// 只要 defaultPerson 为空，就强制把搜索框和搜索结果清空，回归初始态
const checkAndResetState = () => {
  const val = defaultPerson.value
  const isEmpty = Array.isArray(val) ? val.length === 0 : !val

  if (isEmpty) {
    personSearch.value = '' // 清空搜索框文字
    lastSearchResults.value = [] // 清空搜索结果
  }
}

const updateSelectedCache = () => {
  const val = defaultPerson.value
  const currentIds = Array.isArray(val) ? val : val ? [val] : []

  // 1. 减法
  for (const id of selectedCache.value.keys()) {
    if (!currentIds.includes(id as string)) {
      selectedCache.value.delete(id)
    }
  }

  // 2. 加法
  personDefaultData.value.forEach((item) => {
    if (currentIds.includes(item.id)) {
      selectedCache.value.set(item.id, item)
    }
  })
}

const mergeAndDisplay = () => {
  const mergedMap = new Map()

  // 1. 放入搜索结果 (仅当搜索框有值时)
  if (personSearch.value) {
    lastSearchResults.value.forEach((item) => {
      mergedMap.set(item.id, item)
    })
  }

  // 2. 放入已选缓存
  selectedCache.value.forEach((item, id) => {
    mergedMap.set(id, item)
  })

  personDefaultData.value = Array.from(mergedMap.values())
}

const refreshData = () => {
  updateSelectedCache()
  mergeAndDisplay()
}

const handleEsc = () => emit('handleEsc')
const handleBlur = () => emit('handleBlur')
const handleEnter = () => emit('handleEnter')

// 处理搜索
const onSearch = async (value: string | number) => {
  const searchStr = value as string
  updateSelectedCache()

  if (!searchStr) {
    lastSearchResults.value = []
  } else {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const res: any = await searchUser({ searchWord: searchStr })
    lastSearchResults.value = res.list || []
  }
  mergeAndDisplay()
}

// 处理关闭标签 (Tag删除)
const handleClose = () => {
  defaultPerson.value = ''
  emit('update:modelValue', '')

  checkAndResetState() // 触发重置：清空搜索
  refreshData()
}

// 处理人员选择变化 (含点击清除图标、勾选、取消勾选)
const handlePersonChange = (value: SelectValue<SelectOption>) => {
  const stringValue = value as string | string[]
  defaultPerson.value = stringValue

  // 关键：值变化后，检查是否变成了空值。如果是，重置搜索状态。
  checkAndResetState()

  refreshData()
  emit('update:modelValue', stringValue)
}

const findPersonById = (id: number | string) => {
  if (selectedCache.value.has(id)) {
    return selectedCache.value.get(id)
  }
  return personDefaultData.value.find((p) => p.id === id)
}

onMounted(() => {})
if (props.modelValue) {
  const ids = Array.isArray(props.modelValue) ? props.modelValue : [props.modelValue]
  if (ids.length > 0) {
    initPersonData(ids as string[])
  }
}
defineOptions({ name: 'PersonPanel' })
</script>

<style lang="less" scoped>
.external {
  color: orange;
  font-size: 10px;
}
</style>
