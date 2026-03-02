// 显示格式
const formatOptions = [
  { label: '整数', value: 'integer' },
  { label: '小数', value: 'decimal' },
  { label: '百分比', value: 'percentage' },
]
// 小数点位数
const decimalOptions = (length: number) => {
  return Array.from({ length }, (_, i) => ({
    label: `保留${i}位小数`,
    value: i,
  }))
}
// 千分位分隔符
const thousandSeparatorOptions = [
  { label: '逗号 ,', value: 'comma', format: '1,234' },
  { label: '句点 .', value: 'dot', format: '1.234' },
]
// 大数缩写
const abbreviationOptions = [
  { label: '不缩写', value: 'none', format: '' },
  { label: 'K', value: 'k', format: '1,000' },
  { label: 'M', value: 'm', format: '1,000,000' },
  { label: 'B', value: 'b', format: '1,000,000,000' },
  { label: '千', value: 'thousand', format: '1,000' },
  { label: '万', value: 'wan', format: '1,0000' },
  { label: '百万', value: 'million', format: '1,000,000' },
  { label: '亿', value: 'yi', format: '1,000,000,000' },
]
export { formatOptions, decimalOptions, thousandSeparatorOptions, abbreviationOptions }
