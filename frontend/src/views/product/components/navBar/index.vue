<template>
  <div class="w-100 flex tabs-container">
    <div style="max-width: 95%" class="flex">
      <t-button v-if="isShowIcon" variant="text" style="height: 36px" @click="handleLeftShowChange">
        <t-icon class="gray-col" name="chevronRightDoubles" size="20"></t-icon>
      </t-button>
      <t-tabs
        drag-sort
        v-model="activeTab"
        theme="card"
        scroll-position="auto"
        @drag-sort="handleTabsDragend"
        @add="handleAddTab"
      >
        <t-tab-panel v-for="(item, index) in tabList" :key="item.id" :value="item.id">
          <template #label>
            <div
              class="flx-ce-ce"
              style="gap: 4px; font-size: 14px"
              @dblclick.stop="handleTabDbClick(item, index)"
            >
              <t-icon :name="item.type" />
              <span v-if="!editingTab || editingTab.id !== item.id">{{ item.name }}</span>
              <t-input
                v-else
                v-model="tempLabel"
                size="small"
                style="margin-left: 4px; width: 120px"
                @blur="handleInputConfirm(item)"
                @enter="handleInputConfirm(item)"
                :ref="(el) => setEditInputRef(el, index)"
              />
              <t-dropdown
                :min-column-width="200"
                trigger="click"
                :options="dynamicTabsMoreOptions"
                @click="handleMoreDrowChange(item, index, $event)"
              >
                <t-icon class="hide" name="more" />
              </t-dropdown>
            </div>
          </template>
        </t-tab-panel>
      </t-tabs>
      <t-dropdown trigger="click" :options="navbarViewOptions" @click="handleTabsAddOptionsChange">
        <t-button shape="square" variant="text" style="height: 36px; width: 36px">
          <t-icon name="add" />
        </t-button>
      </t-dropdown>
    </div>
  </div>
</template>

<script setup lang="tsx">
import { DropdownProps, TabsProps, Button } from 'tdesign-vue-next'
import type { ComponentPublicInstance } from 'vue'
import { navbarViewOptions, tabsMoreOptions } from '@/modal/options'
import { useMtTableStore } from '@/stores/mtTable'
import emitter from '@/utils/mitt'
import {
  getViewList,
  createView,
  updateView,
  deleteView,
  switchActiveView,
  getTableData,
} from '@/api'
import { ObjType } from '@/types'

const emit = defineEmits(['viewChange', 'handleLeftShowChange'])

const route = useRoute()

const mtTableStore = useMtTableStore()
const activeTab = ref<TabsProps['value']>() // tab初始值
const editingTab = ref<{ id: number | string } | null>(null)
const editInputRefs = ref<HTMLInputElement[]>([]) // editInputRef为数组形式，用于存储多个输入框引用
const tempLabel = ref('') // 临时存储输入框值
const isShowIcon = ref<boolean>(false) // 是否显示展开左侧面板图标
const tableSchemaId = computed(() => route.params.tableSchemaId as string) // 表格数据ID

emitter.on('handleRightPanelShowChange', (show) => {
  isShowIcon.value = show as boolean
})

// 监听updataTableData事件，触发数据更新
emitter.on('updataTableData', async () => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const res: any = await getTableData({ viewId: activeTab.value as string })
  mtTableStore.setCurrentTable(res) // 设置当前表格数据
})

const tabList = ref<ObjType[]>([]) // 视图列表

// 动态设置tabs更多操作选项的禁用状态
const dynamicTabsMoreOptions = computed(() => {
  if (!tabsMoreOptions) return []
  return tabsMoreOptions.map((option) => {
    if (option.value === 'delete') {
      return {
        ...option,
        disabled: tabList.value && tabList.value.length <= 1,
      }
    }
    return option
  })
})

// 监听activeTab变化，触发视图切换
watch(activeTab, async (newVal) => {
  if (newVal) {
    const currentTab = tabList.value.find((tab) => tab.id === newVal)

    await switchActiveView({
      tableSchemaId: route.params.tableSchemaId as string,
      viewId: currentTab?.id,
    })
    emit('viewChange', currentTab)
  }
})

// 更多操作 data-当前tab元数据 index-当前操作tab索引 event-当前操作tab元素
const handleMoreDrowChange = (data: ObjType, index: number, event: ObjType) => {
  console.log('data', data)
  console.log('event', event)
  switch (event.value) {
    case 'rename':
      handleTabDbClick(data, index)
      break
    case 'delete':
      const confirmDia = DialogPlugin({
        header: '删除视图',
        body: `确定要删除视图 「${data.name}」 吗？ 删除视图不会删除数据，不会影响其他视图。`,
        confirmBtn: () => (
          <Button
            theme="danger"
            onClick={async () => {
              await deleteView({ id: data.id })
              MessagePlugin.success('删除成功')
              await handleGetViewList()
              confirmDia.hide()
            }}
          >
            删除
          </Button>
        ),
        cancelBtn: '取消',
        onClose: ({ e, trigger }) => {
          console.log('e: ', e)
          console.log('trigger: ', trigger)
          confirmDia.hide()
        },
      })
      break
  }
}

// 获取视图列表
const handleGetViewList = async () => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const res: any = await getViewList({
    projectId: route.params.id as string,
    tableSchemaId: tableSchemaId.value,
  })
  tabList.value = res.list || []
  //   isDefault 字段为true的tab为默认tab
  activeTab.value = tabList.value.find((tab) => tab.isDefault)?.id || tabList.value[0]?.id || ''
}

// 左侧菜单显示事件
const handleLeftShowChange = () => {
  emitter.emit('handleLeftShowChange')
}

// tabs 添加项事件
const handleTabsAddOptionsChange: DropdownProps['onClick'] = async (data) => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const activeId: any = await createView({
    tableSchemaId: tableSchemaId.value,
    type: data.value as string,
  })
  await handleGetViewList()
  activeTab.value = activeId
}

// 添加tabs
const handleAddTab: TabsProps['onAdd'] = () => {
  console.log('handleAddTab')
}

// Tabs拖拽排序
const handleTabsDragend: TabsProps['onDragSort'] = ({ currentIndex, targetIndex }) => {
  ;[tabList.value[currentIndex], tabList.value[targetIndex]] = [
    tabList.value[targetIndex],
    tabList.value[currentIndex],
  ]
}

// tab双击
const handleTabDbClick = (item: ObjType, index: number) => {
  console.log('item', item)
  //保存当前编辑项和原始标签
  editingTab.value = {
    id: item.id,
  }
  console.log('editingTab.value', editingTab.value)
  tempLabel.value = item.name // 初始化临时输入值
  nextTick(() => {
    if (editInputRefs.value[index]) {
      editInputRefs.value[index].focus()
    }
  })
}

// 输入框确认 更新视图信息
const handleInputConfirm = async (item: ObjType) => {
  await updateView({
    id: item.id,
    name: tempLabel.value,
    description: '',
  })
  await handleGetViewList()
  editingTab.value = null
}

const setEditInputRef = (el: Element | ComponentPublicInstance | null, index: number) => {
  if (el) {
    editInputRefs.value[index] = el as HTMLInputElement
  }
}

watch(
  () => tableSchemaId.value,
  async (newVal) => {
    if (newVal) await handleGetViewList()
  },
  { immediate: true },
)
onMounted(async () => {
  await nextTick()
  if (tableSchemaId.value) {
    await handleGetViewList()
  }
})

defineOptions({
  name: 'NavBar',
})
</script>
<style lang="less" scoped>
@import './index.less';
</style>
