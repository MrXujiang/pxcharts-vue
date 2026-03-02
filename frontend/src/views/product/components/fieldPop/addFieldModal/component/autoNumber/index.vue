<template>
  <t-space size="small" class="w-100" direction="vertical">
    <div class="flex-col gap-8">
      <div>编号类型</div>
      <t-select
        v-model="formData.settings.idType"
        :tips="
          autoNumberTypeOptions.find((option) => option.value === formData.settings.idType)
            ?.description
        "
        :popup-props="{ overlayClassName: 'auto-number' }"
        @change="handleIdTypeChange"
      >
        <t-option
          v-for="item in autoNumberTypeOptions"
          :key="item.value"
          :value="item.value"
          :label="item.label"
        >
          <div>{{ item.label }}</div>
          <div class="description gray-col">{{ item.description }}</div>
        </t-option>
      </t-select>
    </div>
    <div
      v-if="formData.settings.idType === 'customNumber'"
      class="flex-col gap-8"
      style="margin-top: 20px"
    >
      <div class="flx-ce-bet">
        <div>编号格式</div>
        <div>编号预览：{{ previewText }}</div>
      </div>
      <t-card size="small">
        <!-- 第一行：自增数字位数设置 -->
        <div class="flx-ce-bet" style="margin-bottom: 12px">
          <div class="flx-ce-sta gap-4">
            <span>自增数字</span>
            <t-tooltip content="超过当前位数最大值后，数字将继续递增">
              <t-icon class="gray-col" name="helpCircle" size="16" />
            </t-tooltip>
          </div>
          <t-input-number
            v-model="formData.settings.rules.number"
            theme="column"
            :allowInputOverLimit="false"
            :min="1"
            :max="9"
            style="width: 50%"
          >
            <template #suffix><span>位</span></template>
          </t-input-number>
        </div>

        <!-- 第二行：固定文本设置 -->
        <div class="flx-ce-bet">
          <div>固定文本</div>
          <t-input
            v-model="formData.settings.rules.text"
            :maxlength="18"
            show-limit-number
            placeholder="请输入"
            style="width: 50%"
          />
        </div>
      </t-card>
    </div>
  </t-space>
</template>

<script setup lang="ts">
defineOptions({ name: 'AutoNumberComponent' })

const props = defineProps({
  formData: {
    type: Object,
    required: true,
  },
})
const { formData } = toRefs(props)
const currentNumber = ref<number>(1) // 当前自增数字

// 编号类型改变处理(编号类型改变，需还原自增数字跟固定位文本为默认值，并且重置当前编号为1)
const handleIdTypeChange = () => {
  formData.value.settings.rules.number = 3
  formData.value.settings.rules.text = ''
  currentNumber.value = 1
}
// 计算预览文本（数字前补0，文字在后）
const previewText = computed(() => {
  // 获取当前数字的字符串表示
  const currentNumStr = String(currentNumber.value)

  // 确保位数设置在合理范围内（1-9）
  const actualDigitCount = Math.min(Math.max(formData.value.settings.rules.number, 1), 9)

  // 只有当数字位数小于设定的位数时才补0，防止出现大量0的情况
  let paddedNumber = currentNumStr
  if (currentNumStr.length < actualDigitCount) {
    paddedNumber = currentNumStr.padStart(actualDigitCount, '0')
  }

  // 如果有固定文本，则放在数字前面
  if (formData.value.settings.rules.text) {
    return `${formData.value.settings.rules.text}${paddedNumber}`
  }
  return paddedNumber
})

const autoNumberTypeOptions = [
  {
    label: '自增编号',
    value: 'autoIncrement',
    description: '从 1 开始连续递增的数字编号（1, 2, 3...）',
  },
  { label: '自定义编号', value: 'customNumber', description: '自定义编号规则和格式' },
]
</script>

<style lang="less" scoped>
.description {
  font-size: 12px;
}
</style>
