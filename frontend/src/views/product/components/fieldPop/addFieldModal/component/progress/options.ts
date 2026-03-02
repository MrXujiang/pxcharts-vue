// 定义颜色选项类型
interface ColorOption {
  label: string
  value: string
  colors?: string[] // 颜色数组，可以是单色、双色或更多颜色
}

// 格式选项
const formatOptions = [
  { label: '数值', value: 'numVal' },
  { label: '百分比', value: 'percentage' },
]

// 小数位数选项
const decimalOptions = [
  { label: '整数', value: 0 },
  { label: '保留1位小数', value: 1 },
  { label: '保留2位小数', value: 2 },
  { label: '保留3位小数', value: 3 },
]

// 颜色选项数组，根据colors字段渲染颜色
const colorOptions: ColorOption[] = [
  // 单色选项（没有colors字段，使用value作为颜色）
  { label: '海洋蓝', value: '1', colors: ['#2f75fd'] },
  { label: '葡萄紫', value: '2', colors: ['#905afd'] },
  { label: '翡翠绿', value: '3', colors: ['#00b23b'] },
  { label: '天空蓝', value: '4', colors: ['#48ccff'] },
  { label: '火焰红', value: '5', colors: ['#ff5938'] },
  {
    label: '红绿双色',
    value: '9',
    colors: ['#ff5938', '#00b23b'],
  },
  {
    label: '橙蓝双色',
    value: '10',
    colors: ['#ff9900', '#2f75fd'],
  },
  {
    label: '梯度蓝',
    value: '6',
    colors: ['#cadcff', '#a3c1fe', '#81a9fa', '#5a91ff', '#2f75fd'],
  },
  {
    label: '梯度紫',
    value: '7',
    colors: ['#dfd0ff', '#c8adff', '#b592fd', '#a378fc', '#905afd'],
  },
  {
    label: '梯度橙',
    value: '8',
    colors: ['#fff1dd', '#ffe2b2', '#ffc670', '#ffad33', '#ff9900'],
  },
  {
    label: '彩虹色',
    value: '11',
    colors: ['#ff5938', '#ff9900', '#ffdd00', '#00b23b', '#2f75fd'],
  },
]

export { formatOptions, decimalOptions, colorOptions }
