<template>
  <t-card size="small" class="w-100">
    <div class="flex-col gap-8">
      <div class="flex gap-8">
        <t-select
          class="flex1"
          v-model="tableVal"
          :options="options1"
          clearable
          placeholder="选择数据表"
        >
          <template #panelTopContent>
            <div style="padding: 6px 6px 0 6px">
              <t-input v-model="search" placeholder="输入搜索内容" @change="onSearch" />
            </div>
            <t-divider style="margin: 4px 0" />
          </template>
        </t-select>
        <t-select
          class="flex1"
          v-model="value1"
          clearable
          placeholder="需要引用的字段"
          :disabled="!tableVal"
        >
          <template #panelTopContent>
            <div class="flex-col gap-4" style="padding: 6px 6px 0 6px">
              <div class="flx-ce-sta gap-2 col">
                <t-icon name="infoCircle" /> 仅支持选择单选或多选字段
              </div>
              <t-input v-model="search" placeholder="输入搜索内容" @change="onSearch" />
            </div>
            <t-divider style="margin: 4px 0" />
          </template>
          <template #valueDisplay="{ label, value }">
            <div class="flx-ce-sta gap-4">
              <t-icon :name="getIconName(value)" />
              {{ label }}
            </div>
          </template>
          <t-option
            v-for="item in options2"
            :key="item.value"
            :value="item.value"
            :label="item.label"
          >
            <div class="flx-ce-sta gap-4">
              <t-icon :name="item.iconName" />
              <div>{{ item.label }}</div>
            </div>
          </t-option>
        </t-select>
      </div>
      <!-- 新增筛选条件 -->
      <!-- <div style="margin: 0 auto" class="flx-ce-ce fit-content">
        <t-button block theme="primary" variant="text">
          <template #icon>
            <t-icon name="add" />
          </template>
          新增筛选条件
        </t-button>
      </div> -->
    </div>
  </t-card>
</template>

<script setup lang="ts">
import { InputProps } from 'tdesign-vue-next'

const value1 = ref('')
const tableVal = ref('')
const OPTIONS = [
  { label: '架构云', value: '1', iconName: 'formatVerticalAlignLeft' },
  { label: '大数据', value: '2', iconName: 'chevronDownCircle' },
]

const options1 = ref(OPTIONS)
const options2 = ref(OPTIONS)
const search = ref('')

const getIconName = (value: string) => {
  return OPTIONS.find((item) => item.value === value)?.iconName || 'infoCircle'
}
const onSearch: InputProps['onChange'] = () => {
  options1.value = OPTIONS.filter((item) => item.label.indexOf(search.value) !== -1)
}
defineOptions({ name: 'CitationOptions' })
</script>

<style lang="less" scoped>
@import './index.less';
</style>
