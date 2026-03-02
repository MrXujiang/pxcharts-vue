import { defineAsyncComponent } from 'vue'

// 格式化数字，仅截断小数部分，不四舍五入
function formatNumberByDecimal(num, decimal) {
  // 先将输入转为数字类型（处理字符串输入）
  const number = Number(num)

  // 校验数字有效性
  if (isNaN(number)) {
    console.warn('输入不是有效的数字')
    return ''
  }

  // 核心逻辑：仅截断，不四舍五入
  const multiplier = Math.pow(10, decimal) // 10的n次方（n为小数位数）
  const truncatedNum = Math.trunc(number * multiplier) / multiplier // 截断处理

  // 补零并转为字符串（确保小数位数符合要求）
  const formatted = truncatedNum.toFixed(decimal)

  // 整数优化：去除末尾的 .0（可选，可按需删除）
  return decimal === 0 ? formatted.replace('.0', '') : formatted
}

const handleRowChangeValue = (value) => {
  const map = {
    1: 32,
    2: 56,
    3: 88,
    4: 128,
  }
  return map[value]
}

// 自定义编号格式处理
const customNumberFormat = (value, formData) => {
  const currentNumStr = String(value)
  // 确保位数设置在合理范围内（1-9）
  const actualDigitCount = Math.min(Math.max(formData.settings.rules.number, 1), 9)
  // 只有当数字位数小于设定的位数时才补0，防止出现大量0的情况
  let paddedNumber = currentNumStr
  if (currentNumStr.length < actualDigitCount) {
    paddedNumber = currentNumStr.padStart(actualDigitCount, '0')
  }
  let val = ''
  // 如果有固定文本，则放在数字前面
  if (formData.settings.rules.text) {
    val = `${formData.settings.rules.text}${paddedNumber}`
  } else {
    val = paddedNumber
  }
  return val
}

// 格式化单元格值
const formatCellValue = (value, field) => {
  if (value === null || value === undefined) return ''

  switch (field.type) {
    case 'number':
      return typeof value === 'number' ? value.toLocaleString() : value
    case 'person':
      return field.settings.allowMultiple ? value : [value] // 多选返回原数组，单选则需要将字符串加工为数组，方便preview组件处理
      break
    case 'date':
    case 'link':
    case 'image':
    case 'attachment':
      return value
    default:
      return String(value)
  }
}
// 获取行唯一标识
const getRowKey = (item) => {
  //   if (item.isGroupHeader) {
  //     return `group_${item.groupKey}`
  //   }
  //   if (item.isAddRow) {
  //     return `add_${item.groupKey}`
  //   }
  //   return `row_${item.id || item.originalIndex}`
  return item.rowId
}
// 特殊字段不区分显示or编辑模式方法
const isSpecialField = (field) => {
  return [
    'checkbox',
    'image',
    'attachment',
    'rate',
    'creator',
    'createTime',
    'autoNumber',
    'richText',
    'findRef',
  ].includes(field?.type || '')
}

// 筛选特殊字段不区分显示or编辑模式方法
const isSpecialFieldForFilter = (field) => {
  return [
    'checkbox',
    'image',
    'attachment',
    'rate',
    'creator',
    'autoNumber',
    'richText',
    'findRef',
  ].includes(field?.type || '')
}

// 文本
const TextInput = defineAsyncComponent(
  () => import('@/components/multiTable/components/textInput/index.vue'),
)
// 单选/多选
const SingleMultipleChoiceComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/singleMultipleChoice/index.vue'),
)
// 数字
const NumberComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/number/index.vue'),
)
// 日期
const DateComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/date/index.vue'),
)
// 人员
const PersonComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/person/index.vue'),
)
// 币种
const CurrencyComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/currency/index.vue'),
)
// 进度条
const ProgressComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/progress/index.vue'),
)
// 链接
const LinkComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/link/index.vue'),
)
// 复选框
const CheckboxComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/checkbox/index.vue'),
)
// 图片
const ImageComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/imageCom/index.vue'),
)
// 附件
const AttachmentComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/attachment/index.vue'),
)
// 评分
const RateEditComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/rate/index.vue'),
)

// 邮箱
const EmailComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/email/index.vue'),
)
// 身份证
const IdCardComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/idCard/index.vue'),
)

// 富文本
const RichTextComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/richText/index.vue'),
)

// 查找引用
const FindRefComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/findRef/index.vue'),
)

// 创建人
const CreatorComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/creator/index.vue'),
)

// 创建时间
const CreateTimeComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/createTime/index.vue'),
)

// 自动编号
const AutoNumberComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/autoNumber/index.vue'),
)
// 单向、双向关联
const AssociationComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/association/index.vue'),
)

// 单双向类型 筛选组件
const AssociationFilterComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/association/filter.vue'),
)

// 筛选日期、创建时间组件
const FilterDateComponent = defineAsyncComponent(
  () => import('@/components/multiTable/components/filterDate/index.vue'),
)
// 激活的组件
const isActiveComponent = (type: string) => {
  const componentMap = {
    findRef: FindRefComponent, // 查找引用
    text: TextInput, // 文本
    radio: SingleMultipleChoiceComponent, // 单选
    number: NumberComponent, // 数字
    selectMultiple: SingleMultipleChoiceComponent, // 多选
    date: DateComponent, // 日期
    person: PersonComponent, // 人员
    currency: CurrencyComponent, // 币种
    progress: ProgressComponent, // 进度条
    link: LinkComponent, // 链接
    checkbox: CheckboxComponent, // 复选框
    image: ImageComponent, // 图片
    attachment: AttachmentComponent, // 附件
    rate: RateEditComponent, // 评分
    email: EmailComponent, // 邮箱
    idCard: IdCardComponent, // 身份证
    richText: RichTextComponent, // 富文本
    creator: CreatorComponent, // 创建人
    createTime: CreateTimeComponent, // 创建时间
    autoNumber: AutoNumberComponent, // 自动编号
    singleAssociation: AssociationComponent, // 单向关联
    doubleAssociation: AssociationComponent, // 双向关联
  }
  return componentMap[type] || TextInput
}

// 筛选组件
const filterComponent = (type: string) => {
  const componentMap = {
    findRef: NumberComponent, // 查找引用
    text: TextInput, // 文本
    radio: SingleMultipleChoiceComponent, // 单选
    number: NumberComponent, // 数字
    selectMultiple: SingleMultipleChoiceComponent, // 多选
    date: FilterDateComponent, // 日期
    person: PersonComponent, // 人员
    currency: CurrencyComponent, // 币种
    progress: ProgressComponent, // 进度条
    link: TextInput, // 链接
    checkbox: CheckboxComponent, // 复选框
    image: ImageComponent, // 图片
    attachment: AttachmentComponent, // 附件
    rate: RateEditComponent, // 评分
    email: EmailComponent, // 邮箱
    idCard: IdCardComponent, // 身份证
    richText: RichTextComponent, // 富文本
    creator: CreatorComponent, // 创建人
    createTime: FilterDateComponent, // 创建时间
    autoNumber: TextInput, // 自动编号
    singleAssociation: AssociationFilterComponent, // 单向关联
    doubleAssociation: AssociationFilterComponent, // 双向关联
  }
  return componentMap[type] || TextInput
}

// 分组组件
// 单选/多选
const GroupHeadSingleMultipleChoice = defineAsyncComponent(
  () => import('@/components/multiTable/groupHeadCom/singleMultipleChoice/index.vue'),
)
// 日期
const GroupHeadDate = defineAsyncComponent(
  () => import('@/components/multiTable/groupHeadCom/date/index.vue'),
)
// 货币
const GroupHeadCurrency = defineAsyncComponent(
  () => import('@/components/multiTable/groupHeadCom/currency/index.vue'),
)
// 进度条
const GroupHeadProgress = defineAsyncComponent(
  () => import('@/components/multiTable/groupHeadCom/progress/index.vue'),
)
// 链接
const GroupHeadLink = defineAsyncComponent(
  () => import('@/components/multiTable/groupHeadCom/link/index.vue'),
)

// 复选框
const GroupHeadCheckbox = defineAsyncComponent(
  () => import('@/components/multiTable/groupHeadCom/checkbox/index.vue'),
)

// 评分
const GroupHeadRate = defineAsyncComponent(
  () => import('@/components/multiTable/groupHeadCom/rate/index.vue'),
)

// 创建人
const GroupHeadCreator = defineAsyncComponent(
  () => import('@/components/multiTable/groupHeadCom/creator/index.vue'),
)

// 创建时间
const GroupHeadCreateTime = defineAsyncComponent(
  () => import('@/components/multiTable/groupHeadCom/createTime/index.vue'),
)

// 单双向关联
const GroupHeadAssociation = defineAsyncComponent(
  () => import('@/components/multiTable/groupHeadCom/association/index.vue'),
)

const isActiveGroupComponent = (type: string) => {
  const componentMap = {
    radio: GroupHeadSingleMultipleChoice,
    selectMultiple: GroupHeadSingleMultipleChoice,
    date: GroupHeadDate,
    person: GroupHeadSingleMultipleChoice,
    currency: GroupHeadCurrency,
    progress: GroupHeadProgress,
    link: GroupHeadLink,
    checkbox: GroupHeadCheckbox,
    rate: GroupHeadRate,
    creator: GroupHeadCreator,
    createTime: GroupHeadCreateTime,
    singleAssociation: GroupHeadAssociation,
    doubleAssociation: GroupHeadAssociation,
  }
  return componentMap[type]
}

export {
  formatCellValue,
  getRowKey,
  isActiveComponent,
  filterComponent,
  isSpecialField,
  customNumberFormat,
  isSpecialFieldForFilter,
  handleRowChangeValue,
  isActiveGroupComponent,
  formatNumberByDecimal,
}
