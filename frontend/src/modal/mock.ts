const fields = [
  {
    id: 'name',
    title: '姓名',
    type: 'text',
    iconName: 'formatVerticalAlignLeft',
    width: 120,
    fixed: true,
    isShow: true,
  },
  {
    id: 'idText',
    title: 'ID',
    type: 'text',
    iconName: 'formatVerticalAlignLeft',
    width: 180,
    fixed: false,
    isShow: true,
  },

  {
    id: 'department',
    title: '部门',
    type: 'text',
    iconName: 'formatVerticalAlignLeft',
    width: 120,
    isShow: true,
  },
  {
    id: 'joinDate',
    title: '入职日期',
    type: 'date',
    iconName: 'time',
    desc: '',
    format: 'YYYY-MM-DD', // 显示格式
    settings: {
      isDefaultCreateTime: false, // 是否默认创建时间
    },
    defaultValue: '2023-01-01', // 默认日期
    width: 120,
    isShow: true,
  },
]

const personData = [
  {
    id: 0,
    nickname: '添加此记录的人',
    email: '',
    avatar: '',
    department: '',
    position: '',
    isExternal: false,
  },
]
export { fields, personData }
