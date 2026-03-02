import { DropdownProps } from 'tdesign-vue-next'
import {
  LinkIcon,
  FolderExportIcon,
  InternetIcon,
  FileExportIcon,
  CatalogIcon,
  HistoryIcon,
  Delete1Icon,
  LockOnIcon,
} from 'tdesign-icons-vue-next'

const headMoreOptions: DropdownProps['options'] = [
  {
    content: '权限管理',
    value: 'permission',
    prefixIcon: () => <LockOnIcon />,
  },
  { content: '复制链接', value: 'copyLink', prefixIcon: () => <LinkIcon /> },
  { content: '移动到', value: 'move', prefixIcon: () => <FolderExportIcon />, divider: true },
  {
    content: '保存为模板',
    value: 'savaTemplate',
    prefixIcon: () => (
      <svg width="16" height="16" fill="none" viewBox="0 0 24 24" class="_templatecolor24">
        <path
          fill="#EC3713"
          d="m16.867 10.567 5.566 9.543a1.058 1.058 0 0 1-.914 1.591H10.386a1.058 1.058 0 0 1-.914-1.591l5.567-9.543c.408-.7 1.42-.7 1.828 0Z"
        ></path>
        <path
          fill="#FFAD33"
          fill-rule="evenodd"
          d="M9 4.117v7.59c0 1.168.948 2.116 2.117 2.116h.17l2.37-4.062q.769-1.319 2.296-1.319 1.527 0 2.296 1.32l1.995 3.418a2.11 2.11 0 0 0 .597-1.474v-7.59c0-1.168-.948-2.116-2.117-2.116h-7.607C9.947 2 9 2.948 9 4.117Z"
        ></path>
        <path
          fill="#256CFB"
          fill-rule="evenodd"
          d="M7.4 11.706V6.001a6.5 6.5 0 1 0 .896 12.95l2.096-3.593q-1.08-.2-1.903-1.024Q7.4 13.246 7.4 11.706Z"
        ></path>
      </svg>
    ),
  },
  { content: '发布到模版中心', value: 'pubTemplate', prefixIcon: () => <InternetIcon /> },
  { content: '导出为', value: 'export', prefixIcon: () => <FileExportIcon />, divider: true },
  { content: '文档信息', value: 'docInfo', prefixIcon: () => <CatalogIcon /> },
  { content: '历史记录', value: 'history', prefixIcon: () => <HistoryIcon />, divider: true },
  { content: '删除', value: 'delete', prefixIcon: () => <Delete1Icon /> },
]
export { headMoreOptions }
