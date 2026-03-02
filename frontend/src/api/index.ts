import { ObjType } from '@/types'
import req from '@/utils/req'

// 邀请码注册
interface RegisterByInvitationCode {
  verifyCode?: string
  password: string
  inviteCode: string
  email: string
}
const registerByInvitationCode = (data: RegisterByInvitationCode) =>
  req.post('/auth/invite-code/register', data)

// 邮箱验证码注册
interface RegisterByEmailCode {
  verifyCode: string
  password: string
  email: string
}
const registerByEmailCode = (data: RegisterByEmailCode) => req.post('/auth/register', data)

// 登录
interface Login {
  email: string
  password: string
}
const login = (data: Login) => req.post('/auth/login', data)

// 发送邮箱验证码
interface SendEmailCode {
  email: string
}
const sendEmailCode = (data: SendEmailCode) => req.post('/auth/register/send-code', data)

// 重置密码
interface ResetPassword {
  email: string
  verifyCode: string
  password: string
}
const resetPassword = (data: ResetPassword) => req.post('/user/password/reset', data)

// 更新用户信息
interface UpdateUserInfo {
  nickname?: string
  username?: string
  avatar?: string
  tags?: string[]
}
const updateUserInfo = (data: UpdateUserInfo) => req.post('/user/update', data)

// 获取用户信息
const getUserInfo = () => req.get('/user/info')

/**
 * 项目相关
 */
// 获取项目下的所有数据表
const getProjectTables = (data: { projectId: string }) =>
  req.get('/project/tables', { params: data })

interface RenameProjectNode {
  projectId: string
  targetId: string
  name: string
  type: string
}
// 重命名项目节点
const renameProjectNode = (data: RenameProjectNode) => req.post('/project/node/rename', data)

// 更新项目 project/update
const updateProject = (data: {
  id: string
  name?: string
  description?: string
  enableAdvancedPerm?: boolean
}) => req.post('/project/update', data)

// 获取项目列表
interface ProjectListParams {
  page?: number
  size?: number
  type?: number // 1=最近打开 2=我创建的 3=共享给我 4=收藏
}

interface ProjectItem {
  id: string
  name: string
  icon?: string
  color?: string
  creator?: string
  updateTime?: string
  favorite?: boolean
  [key: string]: any
}

interface ProjectListResponse {
  list: ProjectItem[]
  total: number
  page: number
  size: number
}

const getProjectList = (params?: ProjectListParams): Promise<ProjectListResponse> =>
  req.get('/project/query', { params })

// 设置项目为收藏或取消收藏
interface SetProjectFavorite {
  projectId: string
  isFavorite: boolean
}
// 设置项目为收藏或取消收藏
const setProjectFavorite = (data: SetProjectFavorite) => req.post('/project/favorite/set', data)

// 创建项目
interface CreateProject {
  name: string
  description: string
}

interface CreateProjectResponse {
  id: string
  name: string
  description: string
  [key: string]: any
}

const createProject = (data: CreateProject): Promise<CreateProjectResponse> =>
  req.post('/project/create', data)

// 删除项目
const deleteProject = (data: { id: string }) => req.post('/project/delete', data)

// 获取项目详情
const getProjectDetail = (data: { projectId: string }) => req.get('/project/get', { params: data })

/** 文件相关 */

// 文件上传
interface UploadFileResponse {
  url: string
}
const uploadFile = (file: File): Promise<UploadFileResponse> => {
  const formData = new FormData()
  formData.append('file', file)
  return req.post('/file/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

interface CreateProjectNodeResponse {
  folderId?: string
  name: string
  projectId: string
  type: 'table' | 'form' | 'dashboard' | 'folder'
}
// 在项目下创建节点 /v1/project/node/create
const createProjectNode = (data: CreateProjectNodeResponse) =>
  req.post('/project/node/create', data)

// 查询项目下所有文件夹（扁平列表）
const getFolderList = (data: { projectId: string }) => req.get('/folder/list', { params: data })

interface MoveNode {
  projectId: string
  targetFolderId: string // 目标文件夹ID，为空字符串表示移动到根目录
  targetId: string // 要移动的节点ID
  type: string // 节点类型
}
// 移动节点 folder/move
const moveNode = (data: MoveNode) => req.post('/folder/move', data)

// 创建文件夹
const createFolder = (data: { parentId: string; projectId: string; name: string }) =>
  req.post('/folder/create', data)

interface SearchFolder {
  projectId: string
  keyword: string
}
// 左侧文件、文件夹搜索 folder/search
const searchFolder = (data: SearchFolder) => req.get('/folder/search', { params: data })

// 删除文件夹
const deleteFolder = (data: { id: string }) => req.post('/folder/delete', data)

// 更新文件夹
const updateFolder = (data: { id: string; name?: string }) => req.post('/folder/update', data)

// 获取子节点
const getSubFolders = (data: { parentId?: string; projectId?: string }) =>
  req.get('/folder/subquery', { params: data })

// 查询所有节点 folder/all-nodes
const getAllNodes = (data: { projectId: string }) => req.get('/folder/all-nodes', { params: data })

/** 表格字段相关 */

// 获取表格所有记录
const getTableAllRecords = (data: { tableSchemaId: string }) =>
  req.get('/record/list', { params: data })

// 更新整条记录
const updateRowRecordApi = (data: { rowId: string; [key: string]: any }) =>
  req.post('/record/row/update', data)

// 更新单元格数据
interface UpdateCellData {
  fieldId: string
  recordId: string
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  value?: any
}
const updateCellData = (data: UpdateCellData) => req.post('/record/cell/update', data)

// 获取表格数据
const getTableData = (data: { viewId: string }) => req.get('/table/get', { params: data })

// 复制表格
interface CopyTable {
  projectId: string
  name: string
  range: string
  sourceId: string
}
const copyTable = (data: CopyTable) => req.post('/table/copy', data)

// 复制仪表盘
const copyDashboard = (data: { id: string; name: string }) => req.post('/dashboard/copy', data)

// 更新后获取对应更新此条记录 tableSchemaId-表格ID recordId-记录ID(rowId)
const getRecord = (data: { tableSchemaId: string; recordId: string }) =>
  req.get('/record/get', { params: data })

// 更新表格(描述等)
interface UpdateTable {
  id: string
  description?: string
  name?: string
  rowName?: string
}
const updateTable = (data: UpdateTable) => req.post('/table/update', data)

// 创建字段
interface CreateField {
  title: string
  tableSchemaId: string
  type: string
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  [key: string]: any
}
const createField = (data: CreateField) => req.post('/field/create', data)

// 字段排序（列拖拽排序）
const sortField = (data: { fieldIds: string[]; tableSchemaId: string }) =>
  req.post('/field/sort', data)

// 获取字段列表
const getFieldList = (data: { tableSchemaId: string }) => req.get('/field/list', { params: data })

// 删除字段
const deleteField = (data: { id: string }) => req.post('/field/delete', data)

// 更新字段(id - 字段id title - 字段标题 type - 字段类型)
const updateFieldApi = (data: { id: string; title: string; type: string }) =>
  req.post('/field/update', data)

/** 团队相关 */

// 搜索用户
const searchUser = (data: { searchWord: string }) => req.get('/team/search-user', { params: data })

// 批量查询用户
const batchGetUsers = (data: { userIds: string[] }) => req.post('/team/batch-get-users', data)

// 获取团队列表
const getTeamList = () => req.get('/team/list')

// 创建团队
interface CreateTeam {
  name: string
  description: string
  logo: string
}
const createTeam = (data: CreateTeam) => req.post('/team/create', data)

// 删除团队
interface DeleteTeam {
  id: string
}
const deleteTeam = (data: DeleteTeam) => req.post('/team/delete', data)

// 更新团队
interface UpdateTeam {
  id: string
  name: string
  description: string
  logo: string
}
const updateTeam = (data: UpdateTeam) => req.post('/team/update', data)

// 添加团队成员
interface AddTeamMember {
  teamId: string
  userId: string
  identity: string
}
const addTeamMember = (data: AddTeamMember) => req.post('/team/add-member', data)

// 删除团队成员
interface RemoveTeamMember {
  teamId: string
  userId: string
}
const removeTeamMember = (data: RemoveTeamMember) => req.post('/team/delete-member', data)

// 更新团队成员
interface UpdateTeamMember {
  teamId: string
  userId: string
  identity: string
}
const updateTeamMember = (data: UpdateTeamMember) => req.post('/team/update-member', data)

// 获取团队成员列表
interface GetTeamMemberList {
  teamId: string
  page?: number
  size?: number
  searchWord?: string
}
const getTeamMemberList = (data: GetTeamMemberList) =>
  req.get('/team/member-list', { params: data })

export type { ProjectItem, ProjectListResponse, ProjectListParams, CreateProjectResponse }

// 查询模版列表
interface TplParams {
  page?: number // 默认为1
  size?: number // 每页数量
  tag: string // 模版分类，必填，没有则传空字符串
}
const getTemplateList = (params: TplParams) =>
  req.get('/template-project/query', {
    params: {
      ...params,
      tag: params.tag ?? '', // 确保tag字段存在，默认为空字符串
    },
  })

// 查询模版分类列表
interface TplTagParams {
  searchWord: string // 标签关键词
}
const getTemplateTagList = (params: TplTagParams) => req.get('/template-tag/query', { params })

// 获取视图详情(获取配置)
const getViewDetail = (data: { id: string }) => req.get('/view/get', { params: data })

// 更新视图表格(更新配置)
export interface UpdateViewTableConfig {
  colorConfig?: ObjType[]
  filterConfig?: ObjType[]
  groupConfig?: string[]
  viewId: string
  rowHeight?: number
  sortConfig?: ObjType[]
}
const updateViewTable = (data: UpdateViewTableConfig) => req.post('/view/table/update', data)

// 获取表格记录选项列表
const getRecordOptions = (data: { tableSchemaId: string }) =>
  req.get('/record/options', { params: data })

// 获取视图列表 view/query
interface ViewItem {
  projectId: string
  tableSchemaId: string
}
const getViewList = (data: ViewItem) => req.get('/view/query', { params: data })

// 创建视图 view/create
interface CreateView {
  tableSchemaId: string
  type: string
}
const createView = (data: CreateView) => req.post('/view/create', data)

// 更新视图 view/update
interface UpdateView {
  description?: string
  id: string
  name: string
}
const updateView = (data: UpdateView) => req.post('/view/update', data)

// 删除视图
const deleteView = (data: { id: string }) => req.post('/view/delete', data)

// 切换激活视图
const switchActiveView = (data: { tableSchemaId: string; viewId: string }) =>
  req.post('/view/switch-active', data)

// 删除表格
const deleteTable = (data: { id: string }) => req.post('/table/delete', data)

// 删除仪表盘
const deleteDashboard = (data: { id: string }) => req.post('/dashboard/delete', data)

/** 记录相关 - 开始 */
// 插入记录(添加记录) prevId
interface InsertRecord {
  prevId: string // 上一条记录的ID
  tableSchemaId: string // 表格的ID
  rowData?: ObjType // 分组时，当前分组的字段值
}
const insertRecordApi = (data: InsertRecord) => req.post('/record/insert', data)

// 删除记录
const deleteRecord = (data: { recordIds: string[] }) => req.post('/record/delete', data)

/** 记录相关 - 结束 */

// 获取富文本内容
const getRichTextContent = (data: { recordId: string; fieldId: string }) =>
  req.get('/rich-text/get', { params: data })

export {
  registerByInvitationCode,
  registerByEmailCode,
  login,
  sendEmailCode,
  resetPassword,
  updateUserInfo,
  getUserInfo,
  getProjectList,
  setProjectFavorite,
  createProject,
  deleteProject,
  getProjectDetail,
  uploadFile,
  getTeamList,
  createTeam,
  deleteTeam,
  updateTeam,
  addTeamMember,
  removeTeamMember,
  updateTeamMember,
  getTeamMemberList,
  createFolder,
  deleteFolder,
  updateFolder,
  getSubFolders,
  createField,
  getFieldList,
  deleteField,
  updateFieldApi,
  getTemplateList,
  getTemplateTagList,
  getViewList,
  createProjectNode,
  getAllNodes,
  searchFolder,
  getTableData,
  deleteTable,
  deleteDashboard,
  getFolderList,
  moveNode,
  renameProjectNode,
  copyTable,
  copyDashboard,
  updateTable,
  createView,
  updateView,
  deleteView,
  getViewDetail,
  updateViewTable,
  updateCellData,
  sortField,
  insertRecordApi,
  switchActiveView,
  deleteRecord,
  getProjectTables,
  searchUser,
  batchGetUsers,
  getRichTextContent,
  getTableAllRecords,
  getRecord,
  getRecordOptions,
  updateRowRecordApi,
  updateProject,
}
