<template>
  <div class="w-100">
    <t-space size="small" direction="vertical" class="w-100">
      <div class="flx-ce-bet">
        <div>查找引用配置</div>
      </div>
      <t-card size="small">
        <t-space size="small" direction="vertical" class="w-100">
          <t-space size="small" direction="vertical" class="w-100">
            <div>要引用的字段</div>
            <div class="find-ref-select-container">
              <t-select
                v-model="formData.settings.sourceRef.tableSchemaId"
                class="find-ref-select"
                placeholder="请选择表"
                @change="handleRefTableChange"
              >
                <template #prefixIcon>
                  <span class="flx-ce-ce">
                    <t-icon style="width: 16px; height: 16px" name="database" />
                  </span>
                </template>
                <t-option
                  v-for="table in refTableList"
                  :key="table.id"
                  :value="table.id"
                  :label="table.name"
                >
                  <div class="flx-ce-sta gap-4 gray-col">
                    <span class="flx-ce-ce">
                      <t-icon style="width: 16px; height: 16px" name="database" />
                    </span>
                    <div>{{ table.name }}</div>
                  </div>
                </t-option>
              </t-select>
              <t-select
                :disabled="!formData.settings.sourceRef.tableSchemaId"
                v-model="formData.settings.sourceRef.fieldId"
                class="find-ref-select"
                placeholder="请选择字段"
              >
                <template v-if="formData.settings.sourceRef.fieldId" #prefixIcon>
                  <t-icon :name="getIconName(formData.settings.sourceRef.fieldId)" />
                </template>
                <template #panelTopContent>
                  <div style="padding: 6px 6px 0 6px">
                    <t-input v-model="fieldSearch" placeholder="输入搜索内容" @change="onSearch" />
                  </div>
                  <t-divider style="margin: 4px 0" />
                </template>
                <t-option
                  v-for="field in refFieldOptions"
                  :key="field.id"
                  :value="field.id"
                  :label="field.title"
                >
                  <div class="flx-ce-sta gap-4">
                    <t-icon :name="field.iconName" />
                    <div>{{ field.title }}</div>
                  </div>
                </t-option>
              </t-select>
            </div>
          </t-space>

          <!-- 查找条件 -->
          <t-space size="small" direction="vertical" class="w-100">
            <div class="find-ref-condition-header">
              <span>查找条件</span>
              <t-select
                borderless
                v-if="formData.settings.conditions.length > 1"
                v-model="formData.settings.matchCondition"
                class="find-ref-match-select"
                :options="matchOptions"
              />
            </div>

            <div
              v-for="(condition, index) in formData.settings.conditions"
              :key="index"
              class="find-ref-condition-item flex"
            >
              <!-- 引用表字段 -->
              <t-select
                :disabled="!formData.settings.sourceRef.tableSchemaId"
                v-model="condition.sourceTableFieldId"
                class="find-ref-condition-field"
                placeholder="请选择字段"
              >
                <template v-if="condition.sourceTableFieldId" #prefixIcon>
                  <t-icon :name="getIconName(condition.sourceTableFieldId)" />
                </template>
                <template #panelTopContent>
                  <div style="padding: 6px 6px 0 6px">
                    <t-input v-model="fieldSearch" placeholder="输入搜索内容" @change="onSearch" />
                  </div>
                  <t-divider style="margin: 4px 0" />
                </template>
                <t-option
                  v-for="field in refFieldOptions"
                  :key="field.id"
                  :value="field.id"
                  :label="field.title"
                >
                  <div class="flx-ce-sta gap-4">
                    <t-icon :name="field.iconName" />
                    <div>{{ field.title }}</div>
                  </div>
                </t-option>
              </t-select>
              <!-- 条件操作符 -->
              <t-select
                :disabled="!condition.sourceTableFieldId"
                v-model="condition.operator"
                class="find-ref-condition-operator"
                :options="[...getFilterOptions(getFieldType(condition.sourceTableFieldId))]"
              />
              <!-- 当前表字段 -->
              <t-select
                :disabled="!condition.sourceTableFieldId"
                v-model="condition.targetTableFieldId"
                class="find-ref-condition-value flex1"
                placeholder="请选择字段"
              >
                <template v-if="condition.targetTableFieldId" #prefixIcon>
                  <t-icon :name="getIconName(condition.targetTableFieldId)" />
                </template>
                <template #panelTopContent>
                  <div style="padding: 6px 6px 0 6px">
                    <t-input v-model="fieldSearch" placeholder="输入搜索内容" @change="onSearch" />
                  </div>
                  <t-divider style="margin: 4px 0" />
                </template>
                <t-option
                  v-for="field in currentTableFields"
                  :key="field.id"
                  :value="field.id"
                  :label="field.title"
                >
                  <div class="flx-ce-sta gap-4">
                    <t-icon :name="field.iconName" />
                    <div>{{ field.title }}</div>
                  </div>
                </t-option>
              </t-select>
              <!-- <t-button
                variant="text"
                @click="removeCondition(index)"
                v-show="formData.settings.conditions.length > 1"
              >
                <template #icon>
                  <t-icon name="delete" />
                </template>
              </t-button> -->
            </div>

            <!-- <t-button v-if="formData.settings.conditions.length < 9" @click="addCondition">添加筛选条件</t-button> -->
          </t-space>
        </t-space>
      </t-card>
    </t-space>
  </div>
</template>

<script setup lang="ts">
import { useMtTableStore } from '@/stores/mtTable'
import { getProjectTables, getFieldList } from '@/api'
import { ObjType } from '@/types'
import { InputProps, SelectProps } from 'tdesign-vue-next'
import { getFilterOptions } from '@/modal/options'

const { getFields } = useMtTableStore()

const props = defineProps({
  formData: {
    type: Object,
    required: true,
  },
})
const { formData } = toRefs(props)
const route = useRoute()
console.log('route', route)

const currentTableFields = getFields() // 当前表字段选项
const fieldSearch = ref('') // 要引用的字段-字段搜索
const refTableList = ref<ObjType[]>([]) // 引用表列表
const refFieldOptions = ref<ObjType[]>([]) // 引用表字段选项

// 匹配条件选项
const matchOptions = ref([
  { value: 'all', label: '满足所有条件' },
  { value: 'any', label: '满足任一条件' },
])

// 根据sourceTableFieldId 在引用表字段选项中找到对应的字段并且返回type
const getFieldType = (sourceTableFieldId: string) => {
  const field = refFieldOptions.value.find((item) => item.id === sourceTableFieldId)
  return field?.type || ''
}

// 获取当前项目下全部数据表
const handleGetAllTables = async () => {
  const { list }: ObjType = await getProjectTables({ projectId: route.params.id as string })
  refTableList.value = list || []
}

// 获取引用表下字段集
const handleGetRefTableFields = async (tableId: string) => {
  const { fields }: ObjType = await getFieldList({ tableSchemaId: tableId })
  refFieldOptions.value = fields || []
}

// 引用表字段选择变化
const handleRefTableChange: SelectProps['onChange'] = (value) => {
  console.log('value', value)
  handleGetRefTableFields(value as string)
  //   重置引用字段
  formData.value.settings.sourceRef.fieldId = ''
  //   重置条件
  formData.value.settings.conditions.forEach((condition) => {
    condition.sourceTableFieldId = ''
    condition.operator = 'equal'
    condition.targetTableFieldId = ''
  })
}
// 获取字段图标名称
const getIconName = (value: string) => {
  return refFieldOptions.value.find((item) => item.id === value)?.iconName || ''
}

// 引用字段模糊查询
const onSearch: InputProps['onChange'] = (val) => {
  console.log('val--', val)
}
// 添加条件
const addCondition = () => {
  formData.value.settings.conditions.push({
    sourceTableFieldId: '',
    operator: '=',
    targetTableFieldId: '',
  })
}

// 删除条件
const removeCondition = (index: number) => {
  formData.value.settings.conditions.splice(index, 1)
}

onMounted(async () => {
  await handleGetAllTables()
})

defineOptions({ name: 'FindRefComponent' })
</script>

<style scoped>
.find-ref-select-container {
  display: flex;
  gap: 16px;
}

.find-ref-select {
  flex: 1;
}

.find-ref-condition-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.find-ref-match-select {
  width: 120px;
}

.find-ref-condition-item {
  display: flex;
  gap: 8px;
  align-items: center;
}

.find-ref-condition-field,
.find-ref-condition-operator,
.find-ref-condition-value {
  width: 120px;
}

.find-ref-condition-operator {
  flex: 0.5;
}
</style>
