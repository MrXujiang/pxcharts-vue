// 币种
const currencyOptions = [
  { label: 'CNY - ￥', format: '人民币', value: 'RMB' },
  { label: 'HKD - HK$', format: '港元', value: 'HKD' },
  { label: 'TWD - NT$', format: '新台币', value: 'TWD' },
  { label: 'USD - $', format: '美元', value: 'USD' },
  { label: 'EUR - €', format: '欧元', value: 'EUR' },
  { label: 'GBP - £', format: '英镑', value: 'GBP' },
]

const currencyMap = {
  RMB: '￥',
  HKD: 'HK$',
  TWD: 'NT$',
  USD: '$',
  EUR: '€',
  GBP: '£',
}
// 小数点位数
const decimalOptions = (length: number) => {
  return Array.from({ length }, (_, i) => ({
    label: i == 0 ? '整数' : `保留${i}位小数`,
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
  { label: '万', value: 'w', format: '1,0000' },
  { label: '百万', value: 'million', format: '1,000,000' },
  { label: '亿', value: 'yi', format: '100,000,000' },
]
export {
  currencyOptions,
  currencyMap,
  decimalOptions,
  thousandSeparatorOptions,
  abbreviationOptions,
}
