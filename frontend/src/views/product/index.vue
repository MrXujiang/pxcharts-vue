<template>
  <div class="flex-col mt-container">
    <MtHeader :projectDetail="projectDetail" />
    <div class="flex1 overflow-hidden">
      <div class="container">
        <!-- 左侧导航 -->
        <div ref="sidebarRef" class="left" :style="{ width: `${sidebarWidth}px` }">
          <MtSidebar
            @handle-hide-siderbar="handleHideSiderbar"
            @handle-right-change="handleRightChange"
          />
        </div>
        <!-- 左侧和中间的拖拽分割线 -->
        <div ref="leftResizeHandleRef" class="segmentation" @mousedown="startLeftResize"></div>
        <div class="w-100">
          <DataTable ref="dataTableRef" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import DataTable from './components/dataTable/index.vue'
import MtHeader from '@/components/mtHeader/index.vue'
import MtSidebar from '@/views/product/components/mtSidebar/index.vue'
import { useResize } from '@/utils/useResize'
import emitter from '@/utils/mitt'
import { getProjectDetail } from '@/api'
import { ObjType } from '@/types'

defineOptions({
  name: 'ProductEdit',
})

const route = useRoute()
const productId = computed(() => route.params.id as string)
const projectDetail = reactive<ObjType>({})
const dataTableRef = ref()

// 可以在这里根据 productId 加载对应的表格数据
watchEffect(() => {
  console.log('当前编辑的表格ID:', productId.value)
  // TODO: 根据 productId 加载表格数据
})

const { sidebarWidth, startLeftResize, stopResize } = useResize()
const sidebarRef = ref<HTMLElement | null>(null)
const leftResizeHandleRef = ref<HTMLElement | null>(null)

//  左侧及分割线显示
emitter.on('handleLeftShowChange', () => {
  handleHideSiderbar()
})
const handleHideSiderbar = () => {
  emitter.emit('handleRightPanelShowChange', !sidebarRef.value?.classList.contains('hide'))
  sidebarRef.value?.classList.toggle('hide')
  leftResizeHandleRef.value?.classList.toggle('hide')
}

// 处理右侧面板对应内容显示
const handleRightChange = (data: any) => {
  dataTableRef.value?.handleViewChange(data)
}

onMounted(async () => {
  const response = await getProjectDetail({
    projectId: productId.value,
  })
  Object.assign(projectDetail, response)
})
onUnmounted(() => {
  document.removeEventListener('mouseup', stopResize)
  emitter.off('handleLeftShowChange') // 移除左侧面板显示事件监听
})
</script>

<style scoped lang="less">
@import './index.less';
</style>
