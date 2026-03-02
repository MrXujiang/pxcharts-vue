<template>
  <t-dialog
    dialogClassName="default-dialog"
    v-model:visible="visible"
    top="100px"
    :closeBtn="false"
    :header="false"
    @cancel="handleCancel"
    @confirm="handleConfirm"
    width="auto"
  >
    <t-form ref="formRef" :data="formData" label-align="top">
      <t-form-item label="标题" name="title">
        <t-input
          style="min-width: 420px"
          v-model="formData.title"
          placeholder="请输入字段标题"
        ></t-input>
      </t-form-item>
      <t-form-item label="字段类型">
        <t-popup placement="right" trigger="click" v-model:visible="typePopVis">
          <div class="w-100 flx-ce-sta field-type pointer">
            <t-space :size="8" class="flx-ce-sta flex1">
              <span class="field-icon flx-ce-ce"> <t-icon :name="formData.iconName" /></span>
              <div class="flex1">{{ formData.label }}</div>
            </t-space>
            <t-icon name="chevronRight" />
          </div>
          <template #content>
            <t-tabs v-model="value">
              <t-tab-panel
                v-for="(types, category, index) in fieldTypes"
                :key="category"
                :value="index + 1"
                :label="getCategoryLabel(category)"
                :destroy-on-hide="false"
              >
                <div class="field-types-grid">
                  <t-space
                    :size="8"
                    v-for="item in types"
                    :key="item.type"
                    class="field-type-item"
                    @click="selectFieldType(item)"
                  >
                    <span class="field-icon flx-ce-ce">
                      <t-icon :name="item.iconName" />
                    </span>
                    <span class="field-label">{{ item.label }}</span>
                  </t-space>
                </div>
              </t-tab-panel>
            </t-tabs>
          </template>
        </t-popup>
      </t-form-item>
      <!-- 动态字段配置 -->
      <t-form-item
        v-if="
          formData.type && !['image', 'attachment', 'formula', 'idCard'].includes(formData.type)
        "
      >
        <component :is="getComponentByType(formData.type)" v-model:formData="formData" />
      </t-form-item>
      <t-form-item label="字段描述" name="desc">
        <t-textarea v-model="formData.desc" placeholder="请输入字段描述"></t-textarea>
      </t-form-item>
    </t-form>
  </t-dialog>
</template>

<script setup lang="ts">
import { DialogProps } from 'tdesign-vue-next'
import { fieldTypes } from '@/modal/options'
import { TableColumn, useMtTableStore } from '@/stores/mtTable'
import { cloneDeep } from 'lodash'
import { createField, updateFieldApi, getTableAllRecords } from '@/api'
import { ObjType } from '@/types'
import emitter from '@/utils/mitt'
import { customNumberFormat } from '@/utils/mtable'

const TextInput = defineAsyncComponent(() => import('./component/text/index.vue')) // 文本
const SingleMultipleChoiceComponent = defineAsyncComponent(
  () => import('./component/singleMultipleChoice/index.vue'),
) // 多选 or 单选
const NumberComponent = defineAsyncComponent(() => import('./component/number/index.vue')) // 数值
const DateComponent = defineAsyncComponent(() => import('./component/date/index.vue')) // 日期
const PersonComponent = defineAsyncComponent(() => import('./component/person/index.vue')) // 人员
const CurrencyComponent = defineAsyncComponent(() => import('./component/currency/index.vue')) // 货币
const ProgressComponent = defineAsyncComponent(() => import('./component/progress/index.vue')) // 进度
const LinkComponent = defineAsyncComponent(() => import('./component/link/index.vue')) // 链接
const CheckboxComponent = defineAsyncComponent(() => import('./component/checkbox/index.vue')) // 复选框
const RateComponent = defineAsyncComponent(() => import('./component/rate/index.vue')) // 评分
const EmailComponent = defineAsyncComponent(() => import('./component/email/index.vue')) // 邮件
const FindRefComponent = defineAsyncComponent(() => import('./component/findRef/index.vue')) // 引用
const CreatorComponent = defineAsyncComponent(() => import('./component/creator/index.vue')) // 创建人
const AssociationComponent = defineAsyncComponent(() => import('./component/association/index.vue')) // 关联(单向、双向)
const createTimeComponent = defineAsyncComponent(() => import('./component/createTime/index.vue')) // 创建时间
const AutoNumberComponent = defineAsyncComponent(() => import('./component/autoNumber/index.vue')) // 自动编号
const RichTextComponent = defineAsyncComponent(() => import('./component/richText/index.vue')) // 富文本

interface Props {
  isAddVisible: boolean
  fieldConfig?: ObjType // 字段配置(编辑字段时)
}

const props = withDefaults(defineProps<Props>(), {
  fieldConfig: () => ({}),
})
const emit = defineEmits(['update:isAddVisible'])

const route = useRoute()
const formRef = ref()
const value = ref(1)
const typePopVis = ref<boolean>(false) // 字段类型弹pop窗显示状态
const formData = reactive<ObjType>({ ...fieldTypes.common[0], ...props.fieldConfig })
console.log('formData', formData)
const { updateFields, getFields, getRecords, updateRecords } = useMtTableStore()
const visible = computed(() => props.isAddVisible) // 控制模态框的显示状态

// Dialog 的 confirm 事件处理
const handleConfirm = async () => {
  // 验证单双向关联表是否选择
  if (
    ['singleAssociation', 'doubleAssociation'].includes(formData.type) &&
    !formData.settings.sourceRef.tableSchemaId
  ) {
    return MessagePlugin.error('请选择关联的表')
  }

  const currentFields = getFields() // 获取当前字段
  const currentRecords = getRecords() // 获取当前表数据

  const newField: ObjType = {
    ...formData,
    title: formData.title ? formData.title : formData.label,
    isShow: true,
  }
  console.log('newField', newField)

  if (formData.id) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    await updateFieldApi({ ...(newField as any) })
  } else {
    const params = {
      ...newField,
      tableSchemaId: route.params.tableSchemaId as string,
    }
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const id = await createField(params as any)
    newField.id = id as unknown as string
    if (newField.type === 'findRef') {
      // 触发更新表数据
      emitter.emit('updataTableData')
    }
  }

  // 添加新字段
  const updatedFields = [...currentFields, { ...newField }]

  updateFields(updatedFields as TableColumn[]) // 更新字段
  if (newField.type === 'autoNumber') {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const response: any = await getTableAllRecords({
      tableSchemaId: route.params.tableSchemaId as string,
    })
    console.log('response', response)
    if (response.records.length > 0) {
      const fieldId = newField.id

      for (let index = 0; index < response.records.length; index++) {
        const record = currentRecords[index]
        const responseItem = response.records[index]

        if (responseItem && typeof responseItem === 'object' && fieldId in responseItem) {
          record[fieldId] =
            newField.settings.idType == 'autoIncrement'
              ? responseItem[fieldId]
              : customNumberFormat(responseItem[fieldId], newField)
        }
      }
      console.log('currentRecords', currentRecords)
      updateRecords(currentRecords)
    }
  }
  emit('update:isAddVisible', false)
}

// Dialog 的 cancel 事件处理
const handleCancel: DialogProps['onCancel'] = () => {
  // 清空formData
  Object.keys(formData).forEach((key) => {
    formData[key] = ''
  })
  console.log('formData', formData)
  emit('update:isAddVisible', false)
}
// 根据type 对应组件
const getComponentByType = (type: string) => {
  console.log('fieldTypes', fieldTypes)
  const components = {
    text: TextInput,
    radio: SingleMultipleChoiceComponent,
    number: NumberComponent,
    selectMultiple: SingleMultipleChoiceComponent,
    date: DateComponent,
    person: PersonComponent,
    currency: CurrencyComponent,
    progress: ProgressComponent,
    link: LinkComponent,
    checkbox: CheckboxComponent,
    rate: RateComponent,
    email: EmailComponent,
    richText: RichTextComponent,
    findRef: FindRefComponent,
    creator: CreatorComponent,
    singleAssociation: AssociationComponent,
    doubleAssociation: AssociationComponent,
    createTime: createTimeComponent,
    autoNumber: AutoNumberComponent,
  }
  return components[type] || TextInput
}

// 获取分类标签
const getCategoryLabel = (category: string) => {
  const labels: Record<string, string> = {
    common: '常用',
    business: '业务',
    advanced: '高级',
  }
  return labels[category] || category
}

// 选择字段类型
const selectFieldType = (item: Record<string, unknown>) => {
  formData.type = item.type as string
  Object.assign(formData, cloneDeep(item))
  console.log('formData', formData)
  typePopVis.value = false
}

defineOptions({ name: 'AddFieldModal' })
</script>

<style lang="less" scoped>
.field-type {
  height: 32px;
  padding: 0 8px;
  border: 1px solid #ddd;
  border-radius: 3px;
}

.field-types-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  padding: 12px;
  width: 260px;
}
.field-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}
.field-type-item {
  display: flex;
  align-items: center;
  padding: 4px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: #f5f5f5;
  }

  .field-label {
    font-size: 14px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}
</style>
