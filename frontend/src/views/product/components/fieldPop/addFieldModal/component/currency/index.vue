<template>
  <div class="w-100">
    <t-space size="small" class="flex w-100" direction="vertical">
      <!-- 顶部预览栏 -->
      <div class="flx-ce-bet nav-bar">
        <div>货币格式</div>
        <div>
          <span class="gray-col">显示预览:</span>
          <span class="bold">
            <span>{{ currencyMap[formData.settings.currency] }}</span>
            {{ previewText }}
          </span>
        </div>
      </div>

      <!-- 配置面板 -->
      <t-card size="small" class="w-100">
        <!-- 千位分隔符 -->
        <div class="flx-ce-bet option-row">
          <span>使用千位分隔符</span>
          <t-switch v-model="formData.settings.useThousandSeparator" />
        </div>

        <!-- 币种 -->
        <div class="flx-ce-bet option-row">
          <span>币种</span>
          <t-select v-model="formData.settings.currency" style="width: 220px">
            <template #valueDisplay="{ value, label }">
              <div class="flx-ce-bet gap-8">
                <span>{{ label }}</span>
                <span class="gray-col">{{ getCurrencyFormat(value) }}</span>
              </div>
            </template>
            <t-option
              v-for="option in currencyOptions"
              :key="option.value"
              :value="option.value"
              :label="option.label"
              class="customize-option"
            >
              <div class="flx-ce-bet gap-8 w-100">
                <div>{{ option.label }}</div>
                <div class="gray-col" v-if="option.format">{{ option.format }}</div>
              </div>
            </t-option>
          </t-select>
        </div>

        <!-- 小数点位数 -->
        <div class="flx-ce-bet option-row">
          <span>小数点位数</span>
          <t-select
            v-model="formData.settings.decimalPlaces"
            :options="decimalOptions(5)"
            style="width: 220px"
          />
        </div>

        <!-- 千分位分隔符 -->
        <div v-if="formData.settings.useThousandSeparator" class="flx-ce-bet option-row">
          <span>千分位分隔符</span>
          <t-select v-model="formData.settings.thousandSeparator" style="width: 220px">
            <template #valueDisplay="{ label }">
              <div>{{ label }}</div>
            </template>
            <t-option
              v-for="option in thousandSeparatorOptions"
              :key="option.value"
              :value="option.value"
              :label="option.label"
            >
              <div class="flx-ce-sta gap-8">
                <div>{{ option.label }}</div>
                <div class="gray-col" v-if="option.format">{{ option.format }}</div>
              </div>
            </t-option>
          </t-select>
        </div>

        <!-- 大数缩写 -->
        <div class="flx-ce-bet option-row">
          <span>大数缩写</span>
          <t-select
            v-model="formData.settings.largeNumberAbbreviation"
            style="width: 220px"
            @change="handleLargeNumberChange"
          >
            <t-option
              v-for="option in abbreviationOptions"
              :key="option.value"
              :value="option.value"
              :label="option.label"
            >
              <div class="flx-ce-sta gap-8">
                <div>{{ option.label }}</div>
                <div class="gray-col" v-if="option.format">{{ option.format }}</div>
              </div>
            </t-option>
          </t-select>
        </div>

        <!-- 不允许负数 -->
        <div class="flx-ce-sta option-row">
          <t-checkbox v-model="formData.settings.disallowNegative">不允许输入负数</t-checkbox>
        </div>
      </t-card>

      <!-- 输入框 -->
      <div class="w-100 flex-col gap-8">
        <div>默认金额</div>
        <t-input
          v-model="formData.defaultValue"
          @blur="handleBlur"
          @focus="handleFocus"
          @change="handleInputChange"
          placeholder="请输入数字"
          style="width: 100%"
        />
      </div>
    </t-space>
  </div>
</template>

<script setup lang="ts">
import { InputProps } from 'tdesign-vue-next'
import { formatNumber, parseDisplayToRaw, NumberFormatOptions } from '@/utils/numberFormatter'
import {
  currencyOptions,
  decimalOptions,
  thousandSeparatorOptions,
  abbreviationOptions,
  currencyMap,
} from './options'
const props = defineProps({
  formData: {
    type: Object,
    required: true,
  },
})
const { formData } = toRefs(props)

const rawValue = ref<number | null>(1234.0) // 内部真实数值（用于计算）

const getCurrencyFormat = (value: string) => {
  return currencyOptions.find((option) => option.value === value)?.format || ''
}

const getFormatOptions = (): NumberFormatOptions => ({
  currency: formData.value.settings.currency,
  decimalPlaces: formData.value.settings.decimalPlaces as number,
  useThousandSeparator: formData.value.settings.useThousandSeparator,
  thousandSeparator: formData.value.settings.thousandSeparator,
  largeNumberAbbreviation: formData.value.settings.largeNumberAbbreviation,
  disallowNegative: formData.value.settings.disallowNegative,
})
// 大数缩写改变时处理
const handleLargeNumberChange = () => {
  rawValue.value = null
  formData.value.defaultValue = ''
}

// ===== 输入处理 =====
const handleInputChange: InputProps['onChange'] = (val) => {
  if (!/^-?[\d.,%\s]*$/.test(val as string)) {
    MessagePlugin.warning('仅支持输入数字、小数点、逗号、负号和百分号')
    formData.value.defaultValue =
      rawValue.value !== null ? formatNumber(rawValue.value, getFormatOptions()) : ''
    return
  }

  const parsed = parseDisplayToRaw(val as string, getFormatOptions())
  console.log('parsed', parsed)
  if (parsed !== null) {
    rawValue.value = parsed
    formData.value.defaultValue = val as string
  } else {
    formData.value.defaultValue = val as string
  }
}
// ===== 聚焦时处理 =====
const handleFocus = () => {
  formData.value.defaultValue = rawValue.value
}

// ===== 失焦时格式化 =====
function handleBlur() {
  if (rawValue.value !== null) {
    formData.value.defaultValue = formatNumber(rawValue.value, getFormatOptions())
  } else {
    formData.value.defaultValue = ''
  }
}

const previewText = computed(() => {
  const base = 1234
  if (base === undefined) return '0'
  return formatNumber(base, getFormatOptions())
})

watch(
  [formData.value.settings],
  () => {
    if (rawValue.value !== null) {
      formData.value.defaultValue = formatNumber(rawValue.value, getFormatOptions())
    }
  },
  { immediate: true },
)
defineOptions({ name: 'CurrencyComponent' })
</script>

<style lang="less" scoped>
@import './index.less';

.option-row {
  height: 32px;
  line-height: 32px;
  margin-bottom: 8px;

  &:last-child {
    margin-bottom: 0;
  }
}
</style>
