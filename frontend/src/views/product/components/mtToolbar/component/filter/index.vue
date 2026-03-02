<template>
  <div class="filter-popup">
    <div class="title bold">筛选符合以下条件的记录</div>
    <div class="filter-content" v-if="filterList.length > 0">
      <t-space size="small" direction="vertical">
        <div class="w-100 flx-ce-bet gap-4" v-for="(filter, index) in filterList" :key="index">
          <span style="width: 30px">{{ index == 0 ? '当' : '且' }}</span>
          <t-space size="small" class="w-100 flx-ce-bet">
            <!-- 字段选择器 -->
            <t-select
              v-model="filter.fieldId"
              placeholder="请输入内容"
              class="flex1"
              :on-change="
                (value, context) => handleFieldChangeForFilter(value, context, filter, index)
              "
              @popup-visible-change="handlePopupVisibleChange"
            >
              <template #prefixIcon>
                <span style="width: 16px; height: 16px" class="flx-ce-bet">
                  <t-icon :name="getIconNameByType(filter.type)" />
                </span>
              </template>
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
            <t-select
              v-model="filter.operator"
              theme="primary"
              :options="[...getFilterOptions(filter.type)]"
            />

            <div style="width: 200px" v-if="!['image', 'attachment', 'rate'].includes(filter.type)">
              <template v-if="isSpecialFieldForFilter(handleFilterFieldsItem(filter))">
                <component
                  :is="filterComponent(filter?.type || '')"
                  :field="handleFilterFieldsItem(filter)"
                  :item="itemObj"
                  :isEditor="true"
                  @handle-special-change="handleSpecialChange"
                />
              </template>
              <template v-else>
                <component
                  :is="filterComponent(handleFilterFieldsItem(filter)?.type || '')"
                  v-model:editingValue="filter.value"
                  :field="handleFilterFieldsItem(filter)"
                />
              </template>
            </div>
          </t-space>
          <t-button size="small" variant="text" @click="handleDeleteFilter(index)">
            <t-icon name="delete" />
          </t-button>
        </div>
      </t-space>
    </div>
    <div class="flx-ce-bet" style="padding: 8px 0">
      <t-button variant="text" @click="handleAddFilter">
        <template #icon><t-icon name="add" /></template>
        添加筛选条件
      </t-button>
      <t-button v-if="filterList.length > 0" theme="primary" @click="handleFilter"> 筛选 </t-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useMtTableStore } from '@/stores/mtTable'
import { isSpecialFieldForFilter, filterComponent } from '@/utils/mtable'
import { SelectValue, SelectOption } from 'tdesign-vue-next'
import { getIconNameByType } from '@/utils'
import { getFilterOptions } from '@/modal/options'
import { ObjType } from '@/types'
import emitter from '@/utils/mitt'

interface Props {
  filterConfig: ObjType[]
}

const props = defineProps<Props>()

const store = useMtTableStore()
const { getFields } = store

const originalFields = ref<ObjType[]>(getFields() || []) // 保存原始字段列表
// 过滤掉富文本字段
const fields = ref<ObjType[]>(getFields().filter((field) => field.type !== 'richText') || []) // 显示的字段列表
console.log('fields', fields.value)
const search = ref<string>('') // 搜索关键词

const filterList = ref<ObjType[]>(props.filterConfig || []) // 筛选列表
// itemObj 为 fields id值为key，defaultValue 字段值为value
const itemObj = computed(() => {
  return fields.value.reduce((acc, cur) => {
    acc[cur.id] = cur.defaultValue
    return acc
  }, {})
})

// 筛选事件
const handleFilter = () => {
  console.log('filterList.value', filterList.value)
  emitter.emit('handleFilter', filterList.value)
}

// 处理删除筛选条件
const handleDeleteFilter = (index: number) => {
  filterList.value = filterList.value.filter((f, i) => i !== index)
  if (filterList.value.length === 0) {
    handleFilter()
  }
}

// 处理弹出层可见性改变事件
const handlePopupVisibleChange = (visible: boolean) => {
  if (!visible) {
    search.value = ''
    nextTick(() => {
      onSearch()
    })
  }
}

// 处理字段选择器改变事件
const handleFieldChangeForFilter = (
  value: SelectValue,
  context: { option?: SelectOption; selectedOptions: SelectOption[]; trigger: string; e?: Event },
  filter: ObjType,
  index: number,
) => {
  // 将当前数据中的value置空
  filterList.value[index].value = ''

  // 根据选择的字段ID更新筛选项的字段类型
  const selectedField = originalFields.value.find((field) => field.id === value)
  console.log('selectedField', selectedField)
  if (selectedField) {
    filter.type = selectedField.type
    filter.fieldId = value as string // 保存字段值
    filter.value = selectedField.defaultValue
  }
  // 将当前数据中的operator添加默认值
  filterList.value[index].operator = getFilterOptions(filter.type)[0].value
  console.log('filterList.value', filterList.value)
}

const handleFilterFieldsItem = (filter: ObjType) => {
  // 如果筛选项中已有字段对象，则直接返回
  if (filter.field) {
    return filter.field
  }

  // 否则根据筛选项中的字段值查找对应的字段信息
  return originalFields.value.find((f) => f.id === filter.fieldId)
}
// 复选框类型字段、图片类型字段、附件类型字段、评分类型字段改变回调事件
const handleSpecialChange = (obj: ObjType) => {
  itemObj.value[obj.field.id] = obj.value
}

// 添加筛选条件
const handleAddFilter = () => {
  if (filterList.value.length >= 10) {
    return MessagePlugin.error('最多添加10个筛选条件')
  }
  filterList.value.push({
    fieldId: '', // 字段值
    operator: '', // 条件值
    value: '', // 值
    type: '', // 字段类型
  })
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

watch(
  () => props.filterConfig,
  (newVal) => {
    filterList.value = newVal || []
  },
)
defineOptions({
  name: 'FilterPopup',
})
</script>

<style lang="less" scoped>
.filter-popup {
  min-width: 520px;
  .title {
    font-size: 14px;
    line-height: 32px;
    border-bottom: 1px solid #e4e7ed;
  }
  .filter-content {
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
