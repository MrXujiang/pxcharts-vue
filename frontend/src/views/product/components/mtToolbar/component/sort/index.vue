<template>
  <div class="sort-popup">
    <div class="title bold">按以下条件排序</div>
    <div class="sort-content">
      <t-space size="small" direction="vertical">
        <div class="w-100 flx-ce-bet" v-for="sort in sortList" :key="sort.id">
          <div class="w-100 flx-ce-bet gap-8">
            <!-- 字段选择器 -->
            <t-select
              v-model="sort.fieldId"
              placeholder="请输入内容"
              clearable
              class="flex1"
              :on-change="(value, context) => handleFieldChangeForFilter(value, context, sort)"
            >
              <template #panelTopContent>
                <div style="padding: 6px; border-bottom: 1px solid #dcdfe6">
                  <t-input v-model="search" placeholder="请输入关键词搜索" @change="onSearch" />
                </div>
              </template>
              <t-option
                v-for="field in fields"
                :key="field.id"
                :value="field.id"
                :label="field.title"
              >
                <div class="w-100 flx-ce-sta gap-8">
                  <t-icon class="field-icon" :name="getIconNameByType(field.type)" />
                  <span class="field-label">{{ field.title }}</span>
                </div>
              </t-option>
            </t-select>
            <t-radio-group
              v-if="sort.fieldId"
              v-model="sort.order"
              variant="default-filled"
              @change="handleSortChange"
            >
              <t-radio-button
                v-for="option in sort.radioGroup"
                :key="option.value"
                :value="option.value"
                >{{ option.content }}</t-radio-button
              >
            </t-radio-group>
          </div>
          <!-- <t-button size="small" variant="text">
            <t-icon name="delete" />
          </t-button> -->
        </div>
      </t-space>
    </div>
    <!-- <div class="flex" style="padding: 8px 0">
      <t-button variant="text" @click="handleAddSort">
        <template #icon><t-icon name="add" /></template>
        添加排序条件
      </t-button>
    </div> -->
  </div>
</template>

<script setup lang="ts">
import { useMtTableStore } from '@/stores/mtTable'
import { getIconNameByType } from '@/utils'
import { SelectValue, SelectOption } from 'tdesign-vue-next'
import { ObjType } from '@/types'
import { sortConfigMap } from '@/modal/options'
import emitter from '@/utils/mitt'

interface Props {
  sortConfig: ObjType[]
}
const props = defineProps<Props>()

const { getFields } = useMtTableStore()
const sortList = ref<ObjType[]>([{ fieldId: '', order: '', radioGroup: [] }])
const search = ref<string>('')
const fields = ref<ObjType[]>(getFields() || []) // 显示的字段列表
const originalFields = ref<ObjType[]>(getFields() || []) // 保存原始字段列表

// 处理单选组列表
const handleRadioGroupList = (type: string) => {
  // 查找对应配置项
  const matchedConfig = sortConfigMap.find((item) => item.typeList.includes(type))

  // 匹配成功返回排序选项，失败返回空数组并告警（边界处理）
  if (matchedConfig) {
    return matchedConfig.sortOptions
  }
  return []
}

// const handleAddSort = () => {
//   sortList.value.push({ field: '' })
// }

// 处理字段选择器改变事件
const handleFieldChangeForFilter = (
  value: SelectValue,
  context: { option?: SelectOption; selectedOptions: SelectOption[]; trigger: string; e?: Event },
  filter: ObjType,
) => {
  // 根据选择的字段ID更新筛选项的字段类型
  const selectedField = originalFields.value.find((field) => field.id === value)
  console.log('selectedField', selectedField)
  if (selectedField) {
    filter.radioGroup = handleRadioGroupList(selectedField.type)
    filter.fieldId = selectedField.id // 保存字段对象
    filter.order = handleRadioGroupList(selectedField.type)[0].value
  }
  console.log('sortList.value', sortList.value)
  emitter.emit('handleSort', sortList.value)
}

// 处理排序改变事件
const handleSortChange = () => {
  emitter.emit('handleSort', sortList.value)
}

// 搜索
const onSearch = () => {
  if (search.value) {
    console.log('1', 1)
    // 如果有搜索关键词，过滤原始字段列表
    fields.value = originalFields.value.filter((field) => field.title.includes(search.value))
  } else {
    // 如果搜索关键词为空，显示所有原始字段
    fields.value = originalFields.value
  }
}

// 初始化 需要根据props.sortConfig来设置sortList中的radioGroup
const initSortList = () => {
  const selectedField = originalFields.value.find(
    (field) => field.id == props.sortConfig[0]?.fieldId,
  )
  if (selectedField) {
    sortList.value = props.sortConfig.map((item) => {
      return {
        fieldId: item.fieldId,
        order: item.order,
        radioGroup: handleRadioGroupList(selectedField.type),
      }
    })
  }
}
initSortList()

defineOptions({
  name: 'SortPopup',
})
</script>

<style lang="less" scoped>
.sort-popup {
  //   max-width: 460px;
  width: 100%;
  .title {
    font-size: 14px;
    line-height: 32px;
    border-bottom: 1px solid #e4e7ed;
  }
  .sort-content {
    padding: 8px;
    .search {
      padding: 6px 0;
      border-bottom: 1px solid #dcdfe6;
      margin-bottom: 6px;
    }
  }
}
.field-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}
.field-label {
  font-size: 12px;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex-shrink: 1; // 允许收缩
  flex-grow: 1; // 占据可用空间
}
</style>
