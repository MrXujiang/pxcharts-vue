<template>
  <div style="width: 100%">
    <t-tooltip
      size="small"
      :content="
        formatNumberByDecimal(value, settings.decimalPlaces) +
        (settings.numberFormat == 'percentage' ? '%' : '')
      "
    >
      <t-progress
        class="pointer"
        :label="false"
        :percentage="value"
        :color="colorOptions.find((option) => option.value === settings.color)?.colors"
      >
        <template #label>
          <span
            >{{ formatNumberByDecimal(value, settings.decimalPlaces)
            }}{{ settings.numberFormat == 'percentage' ? '%' : '' }}</span
          >
        </template>
      </t-progress>
    </t-tooltip>
  </div>
</template>

<script setup lang="ts">
import { useMtTableStore } from '@/stores/mtTable'
import { formatNumberByDecimal } from '@/utils/mtable'
import { colorOptions } from '@/views/product/components/fieldPop/addFieldModal/component/progress/options'

interface Props {
  value: number
}

const mtTableStore = useMtTableStore()
const props = defineProps<Props>()

const field = computed(() => mtTableStore.getFields())
const groupByField = computed(() =>
  mtTableStore.currentTable.settings.tableConfig.groupConfig?.length > 0
    ? mtTableStore.currentTable.settings.tableConfig.groupConfig[0]
    : '',
)
const settings = computed(() => {
  // 根据groupByField 在field中找到id相等的 settings 值
  return field.value.find((item) => item.id === groupByField.value)?.settings
})
defineOptions({ name: 'GroupHeadProgress' })
</script>

<style scoped>
.date-text {
  color: #999;
  font-size: 12px;
}
</style>
