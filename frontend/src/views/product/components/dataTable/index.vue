<template>
  <div style="height: 100vh">
    <NavBar v-show="isShowNav" @view-change="handleViewChange" />
    <MtToolbar v-if="currentView.type === 'table'" :viewConfig="currentView" />

    <!-- 表格视图 -->
    <MultiTable
      :mvTableData="tableData"
      :fixedFieldsData="fixedFields"
      :scrollableFieldsData="scrollableFields"
      v-if="currentView.type === 'table'"
    />

    <!-- 表单视图 -->
    <FormDesigner
      v-else-if="currentView.type === 'form'"
      :table-id="tableId"
      :fields="fields"
      @save="handleSaveForm"
    />

    <!-- 看板视图 -->
    <KanbanView v-else-if="currentView.type === 'board'" />

    <!-- 仪表盘视图 -->
    <Dashboard
      v-else-if="currentView.type === 'dashboard'"
      :table-id="tableId"
      @save="handleSaveDashboard"
    />
  </div>
</template>

<script setup lang="tsx">
import { useRoute } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import NavBar from '@/views/product/components/navBar/index.vue'
import MtToolbar from '@/views/product/components/mtToolbar/index.vue'
import MultiTable from '@/components/multiTable/index.vue'
import FormDesigner from '@/components/FormDesigner/index.vue'
import Dashboard from '@/components/Dashboard/index.vue'
import KanbanView from '@/views/product/components/kanbanView/index.vue'
import { useMtTableStore } from '@/stores/mtTable'
import { getTableData, getViewDetail } from '@/api'
import { ObjType } from '@/types'

defineOptions({
  name: 'DataTable',
})

const route = useRoute()
const mtTableStore = useMtTableStore()
const isShowNav = ref(true)

// 当前视图参数
const currentView = reactive<ObjType>({})

// 表格ID
const tableId = computed(() => route.params.id as string)

// 获取字段配置
const fields = computed(() => mtTableStore.getFields())

// 表格数据
const tableData = computed(() =>
  mtTableStore.getRecords().length > 0 ? mtTableStore.getRecords() : [],
)

const fixedFields = computed(() => mtTableStore.getFixedFields()) // 固定列字段
const scrollableFields = computed(() => mtTableStore.getScrollableFields()) // 滚动列字段

// 获取表格数据
const handleGetTableData = async (viewId: string) => {
  try {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const res: any = await getTableData({ viewId })
    mtTableStore.setCurrentTable(res) // 设置当前表格数据
    await handleGetViewDetail(viewId) // 获取视图详情配置
  } catch (error) {
    console.error('获取表格数据失败:', error)
  }
}

// 获取视图详情(获取配置项)
const handleGetViewDetail = async (viewId: string) => {
  try {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const res: any = await getViewDetail({ id: viewId })
    mtTableStore.updateTableSettings(res) // 更新表格设置
  } catch (error) {
    console.error('获取视图详情失败:', error)
  }
}

// 视图切换
const handleViewChange = async (view: ObjType) => {
  Object.assign(currentView, view)

  switch (view.type) {
    case 'table':
      isShowNav.value = true
      handleGetTableData(view.id)
      break
    case 'form':
      break
    case 'board':
      break
    case 'dashboard':
      isShowNav.value = false
      break
  }
}

// 数据更新
// const handleDataUpdate = (newData: any[]) => {
//   console.log('数据已更新:', newData)
//   tableData.value = newData
// }

// 保存表单配置
const handleSaveForm = (formConfig: any) => {
  console.log('保存表单配置:', formConfig)
  // TODO: 调用API保存表单配置
  MessagePlugin.success('表单配置已保存')
}

// 保存仪表盘配置
const handleSaveDashboard = (dashboardConfig: any) => {
  console.log('保存仪表盘配置:', dashboardConfig)
  // TODO: 调用API保存仪表盘配置
  MessagePlugin.success('仪表盘配置已保存')
}

defineExpose({
  handleViewChange,
})
</script>
<style lang="less" scoped>
@import './index.less';
</style>
