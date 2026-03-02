// 定义图标数据接口
interface IconData {
  iconName: string
  type: string
  color: string
}

// 定义选项接口
interface SelectOption {
  label: string
  value: string | number
}

// 图标数据数组
const iconData: IconData[] = [
  { iconName: 'starFilled', type: 'star', color: 'orange' },
  { iconName: '', type: 'number', color: 'orange' },
  { iconName: 'heartFilled', type: 'heart', color: 'red' },
]

// 分值起始选项
const scoreStartOptions: SelectOption[] = [
  { label: '0', value: 0 },
  { label: '1', value: 1 },
]

// 分值结束选项
const scoreEndOptions: SelectOption[] = [
  { label: '2', value: 2 },
  { label: '3', value: 3 },
  { label: '4', value: 4 },
  { label: '5', value: 5 },
  { label: '6', value: 6 },
  { label: '7', value: 7 },
  { label: '8', value: 8 },
  { label: '9', value: 9 },
  { label: '10', value: 10 },
]

// 两极文案选项
const textOptions: SelectOption[] = [
  { label: '无', value: '' },
  { label: '满意度', value: 'satisfaction' },
  { label: '可能性', value: 'possibility' },
  { label: '认同度', value: 'agreement' },
  { label: '重要性', value: 'importance' },
  { label: '愿意度', value: 'willingness' },
  { label: '符合度', value: 'conformity' },
]

export { iconData, scoreStartOptions, scoreEndOptions, textOptions, IconData }
