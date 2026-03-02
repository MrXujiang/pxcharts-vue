<template>
  <div class="flex gap-2">
    <template v-if="selectedOption.length === 0">
      <span class="no-content">无内容</span>
    </template>
    <template v-else v-for="(option, index) in selectedOption" :key="option.id">
      <t-tag
        v-if="index === 0"
        size="small"
        shape="round"
        :color="option.color || undefined"
        theme="primary"
      >
        {{ option.label }}
      </t-tag>
      <template v-else-if="index === 1">
        <t-tooltip :content="getTooltipContent" placement="top">
          <t-tag size="small" color="#f0f0f0" theme="primary">
            +{{ selectedOption.length - 1 }}
          </t-tag>
        </t-tooltip>
      </template>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useMtTableStore } from '@/stores/mtTable'

interface Option {
  id: string
  label: string
  color?: string
  // 可能还有其他属性，根据实际需要扩展
}

interface Props {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  value: any
}

const props = defineProps<Props>()

const groupByField = computed(() =>
  mtTableStore.currentTable.settings.tableConfig.groupConfig?.length > 0
    ? mtTableStore.currentTable.settings.tableConfig.groupConfig[0]
    : '',
)
const mtTableStore = useMtTableStore()

const field = computed(() => mtTableStore.getFields())

// 根据groupByField值在field中找到id相等的 options 值，再根据value值找到对应的选项，最后返回显示文本数组
const selectedOption = computed(() => {
  // 查找当前分组字段
  const fieldItem = field.value.find((item) => item.id === groupByField.value)
  if (!fieldItem || !fieldItem.options) {
    return []
  }

  // 获取字段的选项列表
  const options = fieldItem.options as Option[]

  // 处理 value，如果是字符串（逗号分隔）则转换为数组
  let valueArray: string[] = []
  if (typeof props.value === 'string') {
    // 如果是逗号分隔的字符串，分割成数组
    valueArray = props.value.split(',').filter((v) => v.trim() !== '')
  } else if (Array.isArray(props.value)) {
    // 如果已经是数组，直接使用
    valueArray = props.value
  } else {
    // 其他情况转为空数组
    valueArray = []
  }

  // 根据 valueArray 中的值查找对应的选项
  const matchedOptions: Option[] = []
  valueArray.forEach((value) => {
    const option = options.find((opt) => opt.id === value)
    if (option) {
      matchedOptions.push(option)
    }
  })

  return matchedOptions
})
console.log('selectedOption', selectedOption.value)

// 计算额外选项的工具提示内容
const getTooltipContent = computed(() => {
  if (selectedOption.value.length <= 1) {
    return ''
  }

  const additionalOptions = selectedOption.value.slice(1)
  return additionalOptions.map((option) => option.label).join(', ')
})

defineOptions({ name: 'GroupHeadSingleMultipleChoice' })
</script>

<style scoped></style>
