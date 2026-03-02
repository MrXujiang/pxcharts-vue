export interface NumberFormatOptions {
  currency?: string
  displayFormat?: 'decimal' | 'integer' | 'percentage'
  decimalPlaces: number
  useThousandSeparator: boolean
  thousandSeparator: 'comma' | 'dot'
  largeNumberAbbreviation: string
  disallowNegative?: boolean // 主要用于解析时校验
}

/**
 * 格式化数字为显示字符串
 */
export function formatNumber(num: number, options: NumberFormatOptions): string {
  const {
    displayFormat,
    decimalPlaces,
    useThousandSeparator,
    thousandSeparator,
    largeNumberAbbreviation,
  } = options

  if (num === null || isNaN(num)) return ''

  let value = num

  // 百分比：显示 ×100
  if (displayFormat === 'percentage') {
    value = num * 100
  }

  // 应用小数位或整数
  let str: string
  if (displayFormat === 'integer') {
    str = Math.floor(value).toString()
  } else {
    str = value.toFixed(decimalPlaces)
  }

  // 大数缩写（缩写后不加千分位）
  if (largeNumberAbbreviation !== 'none') {
    switch (largeNumberAbbreviation) {
      case 'k':
        str = (value / 1000).toFixed(decimalPlaces) + 'K'
        break
      case 'm':
        str = (value / 1000000).toFixed(decimalPlaces) + 'M'
        break
      case 'b':
        str = (value / 1000000000).toFixed(decimalPlaces) + 'B'
        break
      case 'thousand':
        str = (value / 1000).toFixed(decimalPlaces) + '千'
        break
      case 'w':
        str = (value / 10000).toFixed(decimalPlaces) + '万'
        break
      case 'million':
        str = (value / 1000000).toFixed(decimalPlaces) + '百万'
        break
      case 'yi':
        str = (value / 100000000).toFixed(decimalPlaces) + '亿'
        break
    }

    if (displayFormat === 'percentage') {
      str += '%'
    }
    return str
  }

  // 分离整数和小数部分（基于标准 .）
  const [intPart, decPart] = str.split('.')

  // 千分位符号映射
  const thousandCharMap: Record<string, string> = {
    comma: ',',
    dot: '.',
    // space: ' ',
    // apostrophe: "'",
  }
  const thousandChar = thousandCharMap[thousandSeparator] || ','

  // 格式化整数部分（仅当启用且非百分比）
  let formattedInt = intPart
  if (useThousandSeparator && displayFormat !== 'percentage') {
    formattedInt = intPart.replace(/\B(?=(\d{3})+(?!\d))/g, thousandChar)
  }

  // 组合
  let result = decPart ? `${formattedInt}.${decPart}` : formattedInt

  // 百分比后缀
  if (displayFormat === 'percentage') {
    result += '%'
  }

  return result
}

/**
 * 从显示字符串解析出原始数字（智能识别德式/美式格式）
 */
export function parseDisplayToRaw(display: string, options: NumberFormatOptions): number | null {
  const { displayFormat: format, disallowNegative = false } = options
  if (!display?.trim()) return null

  let text = display.trim()
  let scale = 1 // 缩写对应的放大倍数

  // 1. 处理百分比符号（先移除，最后还原）
  const isPercentage = format === 'percentage'
  if (isPercentage && text.endsWith('%')) {
    text = text.slice(0, -1).trim()
  }

  // 2. 处理大数缩写（使用正确的映射表，支持大小写输入）
  const abbreviationMap: Record<string, number> = {
    K: 1000,
    M: 1000000,
    B: 1000000000,
    千: 1000,
    万: 10000,
    百万: 1000000,
    亿: 100000000, // 正确值：1亿=100,000,000
  }

  // 正则匹配末尾的缩写符号（不区分大小写，支持数值与符号间的空格）
  // 匹配规则：数值部分（含-、.、,） + 可选空格 + 缩写符号（严格匹配映射表中的键，不忽略大小写）
  const abbrKeys = Object.keys(abbreviationMap).join('|')
  const abbreviationRegex = new RegExp(`^(-?[\\d.,]+)\\s*(${abbrKeys})$`)
  const match = text.match(abbreviationRegex)

  if (match) {
    const [_, numberPart, abbr] = match
    // 验证匹配到的缩写是否在支持列表中（严格匹配键名，如 K 匹配，k 不匹配，因映射表键为大写）
    if (Object.prototype.hasOwnProperty.call(abbreviationMap, abbr)) {
      text = numberPart.trim() // 保留数值部分
      scale = abbreviationMap[abbr] // 设置放大倍数
    } else {
      return null // 不支持的缩写符号
    }
  }

  // 3. 清理多余空格（保留数值部分的格式）
  text = text.replace(/\s+/g, '')

  // 4. 验证数值部分仅包含数字、-、.、,（移除缩写后必须是纯数值格式）
  if (!/^-?[\d.,]+$/.test(text)) {
    return null
  }

  // 5. 智能判断小数点位置（原逻辑保留，处理千分位/小数点区分）
  const lastDotIndex = text.lastIndexOf('.')
  const lastCommaIndex = text.lastIndexOf(',')

  let normalized = text
  let decimalChar = ''

  if (lastDotIndex > -1 && lastCommaIndex > -1) {
    // 两者都存在：最后一个特殊字符是小数点
    if (lastDotIndex > lastCommaIndex) {
      decimalChar = '.'
      normalized = text.replace(/,/g, '') // 移除逗号（千分位）
    } else {
      decimalChar = ','
      normalized = text.replace(/\./g, '') // 移除句点（千分位）
    }
  } else if (lastDotIndex > -1) {
    decimalChar = '.'
  } else if (lastCommaIndex > -1) {
    decimalChar = ','
  }

  // 统一小数点为 .
  if (decimalChar === ',') {
    normalized = normalized.replace(',', '.')
  }

  // 6. 验证是否为合法数字字符串
  if (!/^-?\d*\.?\d*$/.test(normalized)) {
    return null
  }

  // 7. 解析数值并应用放大倍数（缩写还原）
  let num = parseFloat(normalized)
  if (isNaN(num)) return null
  num *= scale // 还原缩写（如 1.2K → 1.2 * 1000 = 1200）

  // 8. 百分比还原（如果是百分比格式）
  if (isPercentage) {
    num = num / 100
  }

  // 9. 负数限制校验
  if (disallowNegative && num < 0) {
    return null
  }

  return num
}
