<template>
  <t-card :bordered="false" :bodyStyle="{ padding: '10px 8px' }">
    <div class="mt-sidebar">
      <div class="flx-ce-bet mt-sidebar-toolbar">
        <!-- 搜索状态显示搜索框 -->
        <div v-if="isSearchMode" class="search-container">
          <t-select-input
            :value="selectValue"
            :popup-visible="popupVisible"
            :popup-props="{ overlayInnerStyle: { padding: '6px' } }"
            placeholder="搜索"
            clearable
            allow-input
            @clear="onClear"
            @blur="onBlur"
            @input-change="onInputChange"
          >
            <template #panel>
              <t-tree
                :data="cascaderMenuItems"
                :keys="{ label: 'name', value: 'id', children: 'children' }"
                expand-all
                activable
                @active="handleSearchTreeChange"
              >
                <template #icon="{ node }">
                  <template v-if="node.data.type !== 'folder'">
                    <t-icon :name="node.data.type" />
                  </template>
                  <template v-else>
                    <t-icon :name="node.expanded ? 'caretDownSmall' : 'caretRightSmall'" />
                  </template>
                </template>
              </t-tree>
            </template>
          </t-select-input>
        </div>
        <!-- 正常状态显示标题和工具栏 -->
        <template v-else>
          <div class="flex toolbar-title">
            <span class="title">数据</span>
          </div>
          <div class="flex toolbar-tools">
            <t-tooltip :content="item.tooltip" v-for="item in toolbarTools" :key="item.type">
              <template v-if="item.type != 'add'">
                <t-button
                  shape="square"
                  size="small"
                  variant="text"
                  class="toolbar-item"
                  @click="handleToolBarToolsClick(item)"
                >
                  <template #icon>
                    <t-icon :name="item.icon" />
                  </template>
                </t-button>
              </template>
              <template v-else>
                <t-dropdown
                  trigger="click"
                  :min-column-width="200"
                  :options="toolbarToolsAddOptions"
                  @click="handleToolbarToolsAddClick"
                >
                  <t-button
                    shape="square"
                    size="small"
                    variant="text"
                    class="toolbar-item"
                    @click="handleToolBarToolsClick(item)"
                  >
                    <template #icon>
                      <t-icon :name="item.icon" />
                    </template>
                  </t-button>
                </t-dropdown>
              </template>
            </t-tooltip>
          </div>
        </template>
      </div>
      <div class="siderBar-tree">
        <t-tree
          :data="treeData"
          :keys="{ label: 'name', value: 'id', children: 'children' }"
          v-model:actived="actived"
          activable
          expand-all
          hover
          @active="handleTreeChange"
        >
          <template #icon="{ node }">
            <template v-if="node.data.type !== 'folder'">
              <t-icon :name="node.data.type" />
            </template>
            <template v-else>
              <t-icon :name="node.expanded ? 'caretDownSmall' : 'caretRightSmall'" />
            </template>
          </template>
          <template #label="{ node }">
            <div class="flx-ce-bet" style="padding: 5px">
              <span>{{ node.data.name }} </span>
              <t-space :size="2">
                <t-dropdown
                  :min-column-width="200"
                  trigger="click"
                  :options="treeAddOptions"
                  @click="handleTreeAddClick(node.data, $event)"
                >
                  <t-button
                    v-if="node.data.type == 'folder'"
                    shape="square"
                    size="small"
                    variant="text"
                    @click.stop
                  >
                    <template #icon>
                      <t-icon name="add" />
                    </template>
                  </t-button>
                </t-dropdown>

                <t-dropdown
                  :min-column-width="200"
                  trigger="click"
                  :options="siderBarOptions(node.data.type)"
                  @click="handleSiderBarOptions(node.data, $event)"
                >
                  <t-button shape="square" size="small" variant="text" @click.stop>
                    <template #icon>
                      <t-icon name="more" />
                    </template>
                  </t-button>
                </t-dropdown>
              </t-space>
            </div>
          </template>
        </t-tree>
      </div>
    </div>
  </t-card>
  <t-dialog
    v-model:visible="visible"
    :footer="false"
    :header="createTitle"
    :confirm-on-enter="true"
    :on-cancel="onCancel"
    dialogClassName="from-dialog"
  >
    <t-form ref="form" :data="formData" :rules="FORM_RULES" label-align="top" @submit="onSubmit">
      <t-form-item label="名称" name="name">
        <t-input v-model="formData.name" placeholder="请输入内容"></t-input>
      </t-form-item>
      <t-form-item>
        <div class="w-100 flx-ce-end">
          <t-space>
            <t-button theme="default" @click="onCancel">取消</t-button>
            <t-button theme="primary" type="submit">确认</t-button>
          </t-space>
        </div>
      </t-form-item>
    </t-form>
  </t-dialog>
  <!-- 重命名 -->
  <RenameDialog
    v-if="renameVisible"
    v-model:visible="renameVisible"
    :name="renameName"
    @confirm="handleRenameConfirm"
  />
  <!-- 编辑描述 -->
  <EditTbDescDialog
    v-if="descVisible"
    v-model:visible="descVisible"
    :name="desc"
    @confirm="handleEditConfirm"
  />
  <!-- 复制数据表 -->
  <CopyTableDialog
    v-if="copyVisible"
    v-model:visible="copyVisible"
    :objData="copyName"
    @confirm="handleCopyTable"
  />
</template>

<script setup lang="tsx">
import TIcon from '@/components/TIcon/index.vue'
import {
  toolbarTools,
  toolbarToolsAddOptions,
  tableOrFormOptions,
  dashboardOptions,
  folderOptions,
  fileOptions,
} from '@/modal/options'
import { DropdownProps, FormProps, Button, TreeProps } from 'tdesign-vue-next'
import RenameDialog from '@/components/multiTable/leftOption/rename/index.vue' // 重命名
import EditTbDescDialog from '@/components/multiTable/leftOption/editTbDesc/index.vue' // 编辑描述
import CopyTableDialog from '@/components/multiTable/leftOption/copyTable/index.vue' // 复制数据表
import { ObjType } from '@/types'
import {
  getAllNodes,
  createProjectNode,
  searchFolder,
  deleteFolder,
  deleteTable,
  deleteDashboard,
  getFolderList,
  moveNode,
  renameProjectNode,
  copyTable,
  copyDashboard,
  updateTable,
} from '@/api'
import { emit } from 'process'

defineOptions({
  name: 'MtSidebar',
})

const emits = defineEmits(['handleHideSiderbar', 'handleRightChange'])
const route = useRoute()
const router = useRouter()

const productId = computed(() => route.params.id as string)
const renameVisible = ref<boolean>(false) // 重命名对话框是否可见
const renameName = ref<string>('') // 重命名名称
const createTitle = ref<string>('') // 创建标题
const createType = ref<string>('') // 创建类型

const actived = ref<TreeProps['actived']>([]) // 树选中值
const folderId = ref<string>('') // 文件夹ID
const cascaderMenuItems = ref<ObjType[]>([]) // 级联菜单项
const folderList = ref<ObjType[]>([]) // 文件夹列表

const treeData = ref<ObjType[]>([])

const descVisible = ref<boolean>(false) // 编辑描述对话框是否可见
const desc = ref<string>('') // 编辑描述内容
const copyVisible = ref<boolean>(false) // 复制数据表对话框是否可见
const copyName = reactive({ name: '', copyScope: '' }) // 复制数据表名称

const treeValue = reactive<ObjType>({}) // 树选中值
const currentTarget = reactive<ObjType>({}) // 当前目标值(id、type等)

const treeAddOptions = ref<DropdownProps['options']>(
  toolbarToolsAddOptions?.filter((item) => item.value !== 'importExcel'),
)
const visible = ref<boolean>(false)
const isSearchMode = ref(false)
const searchInput = ref<HTMLInputElement | null>(null)
const selectValue = ref('')
const popupVisible = ref(false)
const form = ref()
const formData: FormProps['data'] = reactive({
  name: '',
})
const FORM_RULES: FormProps['rules'] = {
  name: [
    {
      required: true,
      message: '请输入名称',
    },
  ],
}

// 处理树选中事件
const handleTreeChange: TreeProps['onChange'] = (checked, context) => {
  const { node } = context
  Object.assign(treeValue, node.data)
  if (node.data.type != 'folder') {
    router.replace({ params: { tableSchemaId: node.data.id } })
  }
  emits('handleRightChange', node.data)
}

// 处理搜索树选中事件
const handleSearchTreeChange: TreeProps['onChange'] = (checked, context) => {
  const { node } = context
  if (node.data.type != 'folder') {
    router.replace({ params: { tableSchemaId: node.data.id } })
    actived.value = [node.data.id]
  }
  popupVisible.value = false // 隐藏弹窗
  isSearchMode.value = false // 退出搜索模式
  selectValue.value = '' // 清空选中值
}

// 处理树添加事件
const handleTreeAddClick = (data, e) => {
  folderId.value = data.id
  handleToolbarToolsAddClick(e)
}

const handleLeftSidebarDel = async (data: ObjType) => {
  switch (data.type) {
    case 'folder':
      await deleteFolder({ id: data.id })
      break
    case 'table':
      await deleteTable({ id: data.id })
      break
    case 'dashboard':
      await deleteDashboard({ id: data.id })
      break
    default:
      break
  }
  MessagePlugin.success('删除成功')
}

// 处理重命名确认
const handleRenameConfirm = async (val: string) => {
  await renameProjectNode({
    projectId: productId.value,
    targetId: currentTarget.id,
    name: val,
    type: currentTarget.type,
  })
  await getProjectTreeData()
}

// 处理复制数据表
const handleCopyTable = async (val: ObjType) => {
  console.log('val', val)
  console.log('currentTarget', currentTarget)
  await copyTable({
    projectId: productId.value,
    name: val.dataTableName,
    range: val.copyScope,
    sourceId: currentTarget.id,
  })
  await getProjectTreeData()
  MessagePlugin.success('复制成功')
}

// 处理编辑描述确认
const handleEditConfirm = async (val: string) => {
  await updateTable({
    id: currentTarget.id,
    description: val,
  })
}

// 处理siderbar 操作项
const handleSiderBarOptions = async (data, e: ObjType) => {
  console.log('data', data)
  console.log('e', e)
  Object.assign(currentTarget, data)
  switch (e.value) {
    // 重命名
    case 'rename':
      renameName.value = data.name
      renameVisible.value = true
      break
    //   编辑描述
    case 'editDescription':
      descVisible.value = true
      break
    //   复制数据表
    case 'copyTable':
      copyVisible.value = true
      break
    case 'copyDashboard':
      await copyDashboard({
        id: data.id,
        name: '',
      })
      await getProjectTreeData()
      MessagePlugin.success('复制成功')
      break
    //   删除数据表、仪表盘
    case 'deleteTable':
    case 'deleteDashboard':
    case 'deleteFolder':
    case 'deleteFile':
      const message = {
        deleteTable: '确定删除数据表吗？',
        deleteDashboard: '确定删除仪表盘吗？',
        deleteFolder: '确定删除文件夹吗？',
        deleteFile: '确定删除文件吗？',
      }
      const confirmDia = DialogPlugin({
        header: message[e.value],
        body: data.name,
        confirmBtn: () => (
          <Button
            theme="danger"
            onClick={async () => {
              console.log('confirm button has been clicked!')
              console.log('data', data)
              await handleLeftSidebarDel(data)
              await getProjectTreeData()
              const tableNode = treeData.value.find((item) => item.type === 'table')
              if (tableNode) {
                actived.value = [tableNode?.id]

                router.replace({ params: { tableSchemaId: tableNode.id } })
              }
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
    default:
      await moveNode({
        projectId: productId.value,
        targetFolderId: e.id,
        targetId: data.id,
        type: data.type,
      })
      MessagePlugin.success('移动成功')
      await getProjectTreeData()
      break
  }
}

// 处理siderbar 操作项
const siderBarOptions = (type: string) => {
  switch (type) {
    case 'table':
    case 'form':
      //tableOrFormOptions  value为move类型 添加children
      return (
        tableOrFormOptions &&
        tableOrFormOptions.map((item) => {
          if (item.value === 'move') {
            item.children = folderList.value || []
          }
          return item
        })
      )
    case 'dashboard':
      return (
        dashboardOptions &&
        dashboardOptions.map((item) => {
          if (item.value === 'move') {
            item.children = folderList.value || []
          }
          return item
        })
      )
    case 'folder':
      return folderOptions
    default:
      return fileOptions
  }
}
// 左侧新增操作
const handleLeftSidebarAdd = async () => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const response: any = await createProjectNode({
    folderId: folderId.value || '',
    projectId: productId.value,
    name: formData.name,
    type: createType.value as 'table' | 'form' | 'dashboard' | 'folder',
  })
  actived.value = [response.id]
  emits('handleRightChange', response)
  router.replace({ params: { tableSchemaId: response.id } })
}

// 弹窗关闭
const onCancel = () => {
  folderId.value = ''
  visible.value = false
}
// 重置表单
const handleReset = () => {
  form.value?.reset?.()
}
// 表单提交
const onSubmit: FormProps['onSubmit'] = async ({ validateResult, firstError }) => {
  if (validateResult === true) {
    await handleLeftSidebarAdd()
    await getProjectTreeData()
    handleReset()
    visible.value = false
  } else {
    console.log('Validate Errors: ', firstError, validateResult)
    // MessagePlugin.warning(firstError as string)
  }
}
const handleToolbarToolsAddClick = (data) => {
  console.log('data', data)
  const map = {
    table: '新建数据表',
    form: '新建收集表',
    dashboard: '新建仪表盘',
    folder: '新建文件夹',
    file: '新建文件',
  }
  createTitle.value = map[data.value as keyof typeof map]
  createType.value = data.value as string
  visible.value = true
}

//  工具栏点击事件
const handleToolBarToolsClick = (item: { type: string; icon: string }) => {
  switch (item.type) {
    case 'search':
      enterSearchMode()
      break
    //   隐藏侧边栏
    case 'chevronLeftDouble':
      emits('handleHideSiderbar')
    default:
      break
  }
}

// 进入搜索模式
const enterSearchMode = () => {
  isSearchMode.value = true
  selectValue.value = ''
  nextTick(() => {
    searchInput.value?.focus()
  })
}

const onClear = () => {
  selectValue.value = ''
}

// 搜索框失焦
const onBlur = () => {
  selectValue.value = ''
  isSearchMode.value = false
}

const onInputChange = async (value: string) => {
  selectValue.value = value
  if (!value) {
    cascaderMenuItems.value = []
    popupVisible.value = false
    return
  }
  const response: ObjType = await searchFolder({ projectId: productId.value, keyword: value })

  cascaderMenuItems.value = response.list || []
  popupVisible.value = true
}

// 获取项目树数据(查询左侧树节点)
const getProjectTreeData = async () => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const response: any = await getAllNodes({
    projectId: productId.value as string,
  })
  treeData.value = response.list || []
  const res: ObjType = await getFolderList({ projectId: productId.value as string })
  folderList.value = res.list.map((item) => {
    return {
      ...item,
      value: 'move',
      content: item.name,
    }
  })
}

onMounted(async () => {
  await getProjectTreeData()

  // 获取路由参数中的ID
  const routeTableSchemaId = route.params.tableSchemaId

  // 检查路由参数对应的节点是否存在且类型有效
  let activeNodeId = null
  let activeNodeType = null

  if (routeTableSchemaId) {
    // 检查路由参数对应的节点是否存在于树中
    const routeNode = treeData.value.find((item) => item.id === routeTableSchemaId)
    if (routeNode) {
      // 如果路由参数对应的节点存在，使用它
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      activeNodeId = routeTableSchemaId as any
      activeNodeType = routeNode.type
    } else {
      // 如果路由参数对应的节点不存在，查找第一个table类型节点
      const tableNode = treeData.value.find((item) => item.type === 'table')
      if (tableNode) {
        activeNodeId = tableNode.id
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        activeNodeType = 'table' as any
      }
    }
  } else {
    // 如果没有路由参数，查找第一个table类型节点
    const tableNode = treeData.value.find((item) => item.type === 'table')
    if (tableNode) {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      activeNodeId = tableNode.id as any
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      activeNodeType = 'table' as any
    }
  }
  // 设置激活节点
  actived.value = activeNodeId ? [activeNodeId] : []
  if (activeNodeId) {
    if (activeNodeType === 'dashboard') {
      emits('handleRightChange', {
        type: 'dashboard',
      })
    } else {
      if (route.params.tableSchemaId !== activeNodeId) {
        router.replace({ params: { tableSchemaId: activeNodeId } })
      }
    }
  }
})
</script>

<style lang="less" scoped>
@import './index.less';
</style>
