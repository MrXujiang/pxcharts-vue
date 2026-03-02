import {
  TableIcon,
  FormIcon,
  DashboardIcon,
  FileIcon,
  FolderIcon,
  Edit2Icon,
  InfoCircleIcon,
  CopyIcon,
  FolderExportIcon,
  DeleteIcon,
  LinkIcon,
  AlignTopIcon,
} from 'tdesign-icons-vue-next'
import { TreeProps, DropdownProps } from 'tdesign-vue-next'

// 获取随机颜色面板颜色的函数
const getRandomColor = () => {
  const flatColors = colorRows.flat()
  const randomIndex = Math.floor(Math.random() * flatColors.length)
  return flatColors[randomIndex]
}

// 1. 定义基础筛选片段（原子能力）
const Equality = [
  { label: '等于', value: 'equal' },
  { label: '不等于', value: 'notEqual' },
] as const

const TextLike = [
  { label: '包含', value: 'contains' },
  { label: '不包含', value: 'notContains' },
] as const

const NullCheck = [
  { label: '为空', value: 'null' },
  { label: '不为空', value: 'notNull' },
] as const

const NumericCompare = [
  { label: '大于', value: 'greaterThan' },
  { label: '小于', value: 'lessThan' },
  { label: '大于等于', value: 'greaterThanOrEqual' },
  { label: '小于等于', value: 'lessThanOrEqual' },
] as const

const DateTimeCompare = [
  { label: '等于', value: 'equal' },
  { label: '早于(<)', value: 'before' },
  { label: '晚于(>)', value: 'after' },
  { label: '早于或等于(<=)', value: 'beforeOrEqual' },
  { label: '晚于或等于(>=)', value: 'afterOrEqual' },
] as const

// 表格视图删选条件枚举 (富文本类型不参与筛选)
const tableFilterEnum = {
  // 文本、单选、多选、人员、链接、邮箱、身份证、单双向关联
  default: [...Equality, ...TextLike, ...NullCheck],

  // 创建人
  creator: [...Equality, ...TextLike],

  // 数字、货币、进度、评分、查找引用
  number: [...Equality, ...NumericCompare, ...NullCheck],

  // 自动编号
  autoNumber: [...Equality, ...NumericCompare],

  // 创建时间
  createTime: [...DateTimeCompare],

  // 日期
  date: [...DateTimeCompare, ...NullCheck],

  // 图片、附件
  imageOrAttachment: [...NullCheck],

  // 复选框
  checkbox: [{ label: '等于', value: 'equal' }] as const,
} as const

const FIELD_TYPE_TO_FILTER_KEY: Record<string, keyof typeof tableFilterEnum> = {
  // default 类型
  text: 'default',
  radio: 'default',
  selectMultiple: 'default',
  person: 'default',
  link: 'default',
  email: 'default',
  idCard: 'default',
  singleAssociation: 'default',
  doubleAssociation: 'default',

  // number 类型
  number: 'number',
  currency: 'number',
  progress: 'number',
  rate: 'number',
  findRef: 'number',

  // 特殊类型
  autoNumber: 'autoNumber',
  createTime: 'createTime',
  date: 'date',

  // 图片/附件
  image: 'imageOrAttachment',
  attachment: 'imageOrAttachment',

  // 复选框
  checkbox: 'checkbox',
}
function getFilterOptions(fieldType: string) {
  const filterKey = FIELD_TYPE_TO_FILTER_KEY[fieldType] ?? 'default'
  return tableFilterEnum[filterKey] ?? tableFilterEnum.default
}

// 颜色面板值
const colorRows = [
  [
    '#FF775C',
    '#FFC670',
    '#FFEF85',
    '#AAC223',
    '#47B95C',
    '#00C7BC',
    '#A6E1FC',
    '#81A9FA',
    '#B592FD',
    '#E188CC',
    '#A6A6A6',
  ],
  [
    '#FFE9E5',
    '#FFF1DD',
    '#FFF7D1',
    '#F2F6D5',
    '#E0F1DF',
    '#D4F7F2',
    '#D9F2FF',
    '#E6EDFA',
    '#F2E9FC',
    '#FFE8F4',
    '#EEEEEE',
  ],
  [
    '#FFD0C7',
    '#FFE2B2',
    '#FFF5B8',
    '#DAE3AB',
    '#B8E0B8',
    '#B5E9E4',
    '#C7EBFF',
    '#CADCFF',
    '#DFD0FF',
    '#FFC6ED',
    '#D7D7D7',
  ],
  [
    '#EC3713',
    '#E48900',
    '#E5C600',
    '#7D8F00',
    '#009931',
    '#009990',
    '#00B2E5',
    '#2F75FD',
    '#905AFD',
    '#CC45B6',
    '#7D7D7D',
  ],
  [
    '#C7371A',
    '#B86E00',
    '#BDA400',
    '#657500',
    '#007B2B',
    '#008078',
    '#218FB5',
    '#2160D9',
    '#7E4ADC',
    '#B0369B',
    '#696969',
  ],
]

const cascaderMenuItems = ref([])
// 左侧侧边栏操作组
const toolbarTools = [
  { type: 'search', icon: 'search', tooltip: '搜索' },
  { type: 'add', icon: 'add', tooltip: '添加' },
  { type: 'chevronLeftDouble', icon: 'chevronLeftDouble', tooltip: '隐藏侧边栏' },
]
// 添加操作项
const toolbarToolsAddOptions: DropdownProps['options'] = [
  {
    content: '导入 Excel (.xlsx)',
    value: 'importExcel',
    prefixIcon: () => (
      <svg width="16" height="16" fill="none" viewBox="0 0 16 16" class="_microsoftexcel16">
        <defs>
          <clipPath id="7be1b1_0">
            <rect width="16" height="16" rx="0"></rect>
          </clipPath>
          <clipPath id="7be1b1_1">
            <rect width="14.545" height="14.545" x="0.1" y="0.747" rx="0"></rect>
          </clipPath>
        </defs>
        <g clip-path="url(#7be1b1_0)">
          <g clip-path="url(#7be1b1_1)">
            <path
              fill="#21A366"
              d="M9.797 1.352H4.342c-.363 0-.606.243-.606.606v3.03l6.06 3.03 2.425.91 2.424-.91v-3.03L9.797 1.353Z"
              style="mix-blend-mode: normal;"
            ></path>
            <path
              fill="#107C41"
              d="M3.736 4.989h6.06v3.03h-6.06v-3.03Z"
              style="mix-blend-mode: normal;"
            ></path>
            <path
              fill="#33C481"
              d="M14.645 1.958v3.03H9.797V1.353h4.242c.303 0 .606.303.606.606Z"
              style="mix-blend-mode: normal;"
            ></path>
            <path
              fill="#185C37"
              d="M9.796 8.019h-6.06v6.06c0 .364.242.607.606.607h9.697c.364 0 .606-.243.606-.606v-3.03L9.796 8.018Z"
              style="mix-blend-mode: normal;"
            ></path>
            <g style="opacity: 0.5;">
              <path
                fill="#000"
                d="M8.462 4.383H3.735v8.485h4.606c.424 0 .849-.425.849-.849V5.11a.7.7 0 0 0-.728-.727Z"
                style="mix-blend-mode: normal;"
              ></path>
            </g>
            <path
              fill="#107C41"
              d="M7.857 12.262H.827a.7.7 0 0 1-.727-.728v-7.03a.7.7 0 0 1 .727-.727h7.09c.364 0 .667.303.667.727v7.091c0 .364-.303.667-.727.667Z"
              style="mix-blend-mode: normal;"
            ></path>
            <path
              fill="#FFF"
              d="m2.16 10.444 1.576-2.425L2.28 5.595h1.152L4.22 7.11c.12.182.12.303.181.364l.182-.364.849-1.515h1.09L5.07 8.019l1.515 2.425H5.372l-.91-1.697c0-.061-.06-.122-.12-.243 0 .06-.06.121-.121.243l-.91 1.697H2.16Z"
              style="mix-blend-mode: normal;"
            ></path>
            <path
              fill="#107C41"
              d="M9.796 8.02h4.848v3.03H9.796V8.02Z"
              style="mix-blend-mode: normal;"
            ></path>
          </g>
        </g>
      </svg>
    ),
  },
  { content: '数据表', value: 'table', prefixIcon: () => <TableIcon /> },
  { content: '收集表', value: 'form', prefixIcon: () => <FormIcon /> },
  { content: '仪表盘', value: 'dashboard', prefixIcon: () => <DashboardIcon /> },
  { content: '文档', value: 'doc', prefixIcon: () => <FileIcon /> },
  { content: '文件夹', value: 'folder', prefixIcon: () => <FolderIcon /> },
]

const navbarViewOptions = [
  { content: '表格视图', value: 'table', prefixIcon: () => <TableIcon /> },
  { content: '表单视图', value: 'form', prefixIcon: () => <FormIcon /> },
  { content: '看板视图', value: 'board', prefixIcon: () => <DashboardIcon /> },
]

const items = ref<TreeProps['data']>([
  {
    label: '数据表',
    value: 'table',
    type: 'table',
    children: [],
  },
  {
    label: '收集表',
    value: 'form',
    type: 'form',
    children: [],
  },
  {
    label: '仪表盘',
    value: 'dashboard',
    type: 'dashboard',
    children: [],
  },
  {
    label: '文档',
    value: 'file',
    type: 'file',
    children: [],
  },
  {
    label: '数据可视化',
    value: 'folder',
    type: 'folder',
    children: true,
  },
])

// 数据表跟收集表操作项
const tableOrFormOptions: DropdownProps['options'] = [
  { content: '重命名', value: 'rename', prefixIcon: () => <Edit2Icon /> },
  { content: '编辑数据表描述', value: 'editDescription', prefixIcon: () => <InfoCircleIcon /> },
  { content: '复制数据表', value: 'copyTable', prefixIcon: () => <CopyIcon /> },
  { content: '移动至', value: 'move', prefixIcon: () => <FolderExportIcon />, divider: true },
  { content: '删除数据表', value: 'deleteTable', prefixIcon: () => <DeleteIcon /> },
]
// 仪表盘操作项
const dashboardOptions: DropdownProps['options'] = [
  { content: '重命名', value: 'rename', prefixIcon: () => <Edit2Icon /> },
  { content: '复制仪表盘', value: 'copyDashboard', prefixIcon: () => <CopyIcon /> },
  { content: '移动至', value: 'move', prefixIcon: () => <FolderExportIcon />, divider: true },
  { content: '删除仪表盘', value: 'deleteDashboard', prefixIcon: () => <DeleteIcon /> },
]
// 文件夹操作项
const folderOptions: DropdownProps['options'] = [
  { content: '重命名', value: 'rename', prefixIcon: () => <Edit2Icon />, divider: true },
  { content: '删除文件夹', value: 'deleteFolder', prefixIcon: () => <DeleteIcon /> },
]
const fileOptions: DropdownProps['options'] = [
  { content: '重命名', value: 'rename', prefixIcon: () => <Edit2Icon /> },
  { content: '移动至', value: 'move', prefixIcon: () => <FolderExportIcon />, divider: true },
  { content: '移除文档', value: 'deleteFile', prefixIcon: () => <DeleteIcon /> },
]

// tabs选项中more 操作项
const tabsMoreOptions: DropdownProps['options'] = [
  { content: '重命名', value: 'rename', prefixIcon: () => <Edit2Icon /> },
  {
    content: '复制视图链接',
    value: 'copyViewLink',
    prefixIcon: () => <LinkIcon />,
    divider: true,
  },
  { content: '删除视图', value: 'delete', prefixIcon: () => <DeleteIcon /> },
]

// 字段类型(常用、业务、高级)
const fieldTypes = {
  // 常用
  common: [
    {
      width: 120,
      label: '文本',
      title: '',
      desc: '',
      id: '',
      type: 'text',
      iconName: 'formatVerticalAlignLeft',
      defaultValue: '',
    },
    {
      width: 120,
      label: '单选',
      title: '',
      desc: '',
      id: '',
      type: 'radio',
      iconName: 'chevronDownCircle',
      options: [
        {
          id: '1',
          label: '选项一',
          color: getRandomColor(),
        },
        {
          id: '2',
          label: '选项二',
          color: getRandomColor(),
        },
      ],
      defaultValue: '',
      settings: {
        isReferenced: false, // 是否引用
      },
      ref: {}, // 引用选项，空对象时不引用，如果有值，如下
      //   ref: {
      //     table_id: '',
      //     field_id: ''
      //   }
    },
    {
      width: 120,
      label: '数字',
      title: '',
      desc: '',
      id: '',
      type: 'number',
      iconName: 'artboard',
      settings: {
        // 数字相关的设置配置
        useThousandSeparator: false, // 是否使用千位分隔符（开关状态，默认关闭）
        displayFormat: 'decimal', // 显示格式（下拉选项，如“小数”对应“decimal”）
        decimalPlaces: 1, // 小数点位数（下拉选项，如“保留2位小数”）
        thousandSeparator: 'comma', // 千位分隔符
        largeNumberAbbreviation: 'none', // 大数缩写方式（下拉选项，如“不缩写”）
        disallowNegative: false, // 是否不允许输入负数（复选框状态，默认未勾选）
      },
      defaultValue: '1234.0', // 默认数字（对应“请输入数字”输入框）
    },
    {
      width: 120,
      label: '多选',
      title: '',
      desc: '',
      id: '',
      type: 'selectMultiple',
      iconName: 'componentCheckbox',
      options: [
        {
          id: '1',
          label: '选项一',
          color: getRandomColor(),
        },
        {
          id: '2',
          label: '选项二',
          color: getRandomColor(),
        },
      ],
      defaultValue: [],
      settings: {
        isReferenced: false, // 是否引用
      },
      ref: {}, // 引用选项，空对象时不引用，如果有值，如下
      //   ref: {
      //     table_id: '',
      //     field_id: ''
      //   }
    },
    {
      width: 120,
      label: '日期',
      title: '',
      desc: '',
      id: '',
      type: 'date',
      iconName: 'time',
      format: 'YYYY-MM-DD', // 显示格式
      settings: {
        isDefaultCreateTime: false, // 是否默认创建时间
      },
      defaultValue: '', // 默认日期
    },
    {
      label: '人员',
      title: '',
      desc: '',
      id: '',
      type: 'person',
      iconName: 'user1',
      settings: {
        showNickname: true, // 是否展示昵称
        showExternalTag: true, // 是否展示组织外部标签
        allowMultiple: true, // 是否允许多选
      },
      defaultValue: [],
    },
    {
      label: '货币',
      title: '',
      desc: '',
      id: '',
      type: 'currency',
      iconName: 'currencyExchange',
      settings: {
        // 数字相关的设置配置
        useThousandSeparator: false, // 使用千位分隔符：是否启用千位分隔符（如1,234.56）
        currency: 'RMB', // 币种
        decimalPlaces: 1, // 小数点位数：保留的小数位数（示例中2位对应¥1.00）
        thousandSeparator: 'comma', //  千分位分隔符
        largeNumberAbbreviation: 'none', // 大数缩写：是否对大额数字进行缩写（如1k代表1000）
        disallowNegative: false, // 不允许输入负数：是否禁止输入负数金额
      },
      defaultValue: '1234.0', // 默认金额：默认的初始金额值（用户可输入数字）
    },
    {
      width: 150,
      label: '图片',
      title: '',
      desc: '',
      id: '',
      type: 'image',
      iconName: 'image',
      defaultValue: [],
    },
    {
      width: 150,
      label: '附件',
      title: '',
      desc: '',
      id: '',
      type: 'attachment',
      iconName: 'fileAttachment',
    },
    {
      width: 150,
      label: '进度',
      title: '',
      desc: '',
      id: '',
      type: 'progress',
      iconName: 'componentInput',
      settings: {
        numberFormat: 'numVal', // 数值格式
        decimalPlaces: 0, // 小数点位数
        color: '1', // 颜色
        enableCustomProgress: false, // 自定义进度条值开关
        // startValue: 0, // 起始值
        // targetValue: 100, // 目标值
      },
      defaultValue: null,
    },
    {
      width: 200,
      label: '链接',
      title: '',
      desc: '',
      id: '',
      type: 'link',
      iconName: 'link',
      linkTitle: '',
      settings: {
        isLinkTitle: false, // 是否显示链接标题
      },
      defaultValue: {
        linkTitle: '', // 链接标题
        linkUrl: '', // 链接URL
      },
    },
    {
      width: 120,
      label: '公式',
      title: '',
      desc: '',
      id: '',
      type: 'formula',
      iconName: 'functions1',
    },
  ],
  // 业务
  business: [
    {
      width: 120,
      label: '复选框',
      title: '',
      desc: '',
      id: '',
      type: 'checkbox',
      iconName: 'checkRectangle',
      defaultValue: 0, // 默认不勾选
    },
    {
      width: 120,
      label: '人员',
      title: '',
      desc: '',
      id: '',
      type: 'person',
      iconName: 'user1',
      settings: {
        showNickname: true, // 是否展示昵称
        showExternalTag: true, // 是否展示组织外部标签
        allowMultiple: true, // 是否允许多选
      },
      defaultValue: [],
    },
    {
      width: 150,
      label: '评分',
      title: '',
      desc: '',
      id: '',
      type: 'rate',
      iconName: 'star',
      settings: {
        icon: 'star',
        scoreStart: 1, // 分值起始
        scoreEnd: 5, // 分值结束
        // text: '', // 两极文案
      },
      defaultValue: null,
    },
    {
      width: 120,
      label: '邮箱',
      title: '',
      desc: '',
      id: '',
      type: 'email',
      iconName: 'email',
      defaultValue: '',
    },
    // {
    //   label: '地址',
    //   type: 'address',
    //   iconName: 'address',
    // },
    {
      width: 180,
      label: '身份证',
      title: '',
      desc: '',
      id: '',
      type: 'idCard',
      iconName: 'idCard',
      defaultValue: '',
    },
    {
      width: 150,
      label: '富文本',
      title: '',
      desc: '',
      id: '',
      type: 'richText',
      iconName: 'textInitial',
      defaultValue: '',
    },
  ],
  // 高级
  advanced: [
    {
      width: 180,
      label: '查找引用',
      title: '',
      desc: '',
      id: '',
      type: 'findRef',
      iconName: 'dataSearch',
      defaultValue: '',
      settings: {
        /* 引用的表和字段 */
        sourceRef: {
          tableSchemaId: '',
          fieldId: '',
        },
        matchCondition: 'all', // 匹配条件选项 (满足所有条件、满足任一条件)
        /* 引用条件 */
        conditions: [
          {
            sourceTableFieldId: '',
            operator: 'equal', //  大于（gt），等于（equal），大于等于（ge），小于等于（le），不为空（not null），为空（null）
            targetTableFieldId: '',
          },
        ],
      },
    },
    {
      width: 120,
      label: '创建人',
      title: '',
      desc: '',
      id: '',
      type: 'creator',
      iconName: 'userAdd',
      settings: {
        showNickname: true, // 是否展示昵称
        showExternalTag: true, // 是否展示组织外部标签
      },
      defaultValue: null,
    },
    {
      width: 150,
      label: '单向关联',
      title: '',
      desc: '',
      id: '',
      type: 'singleAssociation',
      iconName: 'swapRight',
      settings: {
        sourceRef: {
          tableSchemaId: '',
        },
        conditions: {
          allowMultiSelect: true,
        },
      },
      defaultValue: null,
    },
    {
      width: 150,
      label: '双向关联',
      title: '',
      desc: '',
      id: '',
      type: 'doubleAssociation',
      iconName: 'swap',
      settings: {
        sourceRef: {
          tableSchemaId: '',
        },
        conditions: {
          allowMultiSelect: true,
        },
      },
      defaultValue: null,
    },
    {
      width: 200,
      label: '创建时间',
      title: '',
      desc: '',
      id: '',
      type: 'createTime',
      iconName: 'calendarPlus',
      format: 'YYYY-MM-DD', // 显示格式
      defaultValue: '',
    },
    {
      width: 120,
      label: '自动编号',
      title: '',
      desc: '',
      id: '',
      type: 'autoNumber',
      iconName: 'arrowUp01',
      settings: {
        idType: 'autoIncrement', // 编号类型 autoIncrement | customNumber
        rules: {
          number: 3,
          text: '',
        },
      },
      defaultValue: '',
    },
  ],
}

// 字段管理more 操作项
const fieldManageMoreOptions: DropdownProps['options'] = [
  { content: '移到顶部', value: 'moveToTop', prefixIcon: () => <AlignTopIcon />, divider: true },
  { content: '编辑', value: 'edit', prefixIcon: () => <Edit2Icon />, divider: true },
  { content: '删除字段', theme: 'error', value: 'deleteField', prefixIcon: () => <DeleteIcon /> },
]

// 排序选项
const sortConfigMap = [
  {
    typeList: [
      'text',
      'person',
      'link',
      'email',
      'idCard',
      'creator',
      'singleAssociation',
      'doubleAssociation',
    ],
    sortOptions: [
      { content: 'A-Z', value: 'asc' },
      { content: 'Z-A', value: 'desc' },
    ],
    defaultValue: 'asc',
  },
  {
    typeList: ['radio', 'selectMultiple'],
    sortOptions: [
      { content: '选项顺序', value: 'asc' },
      { content: '选项逆序', value: 'desc' },
    ],
    defaultValue: 'asc',
  },
  {
    typeList: ['number', 'currency', 'progress', 'rate', 'findRef', 'autoNumber'],
    sortOptions: [
      { content: '1-9', value: 'asc' },
      { content: '9-1', value: 'desc' },
    ],
    defaultValue: 'asc',
  },
  {
    typeList: ['date', 'createTime'],
    sortOptions: [
      { content: '早-晚', value: 'asc' },
      { content: '晚-早', value: 'desc' },
    ],
    defaultValue: 'asc',
  },
  {
    typeList: ['checkbox'],
    sortOptions: [
      { content: '顺序', value: 'asc' },
      { content: '逆序', value: 'desc' },
    ],
    defaultValue: 'asc',
  },
]

const rowHeightOptions = [
  { content: '常规', value: 1 },
  { content: '中等', value: 2 },
  { content: '高', value: 3 },
  { content: '超高', value: 4 },
]

export {
  cascaderMenuItems,
  toolbarTools,
  toolbarToolsAddOptions,
  items,
  tableOrFormOptions,
  fileOptions,
  dashboardOptions,
  folderOptions,
  tabsMoreOptions,
  fieldTypes,
  fieldManageMoreOptions,
  colorRows,
  sortConfigMap,
  navbarViewOptions,
  rowHeightOptions,
  getRandomColor,
  getFilterOptions,
}
