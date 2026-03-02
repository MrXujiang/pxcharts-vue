<template>
  <t-drawer
    v-model:visible="visible"
    size="large"
    header="添加记录"
    :close-btn="true"
    @close="handleCancel"
    @confirm="handleConfirm"
    width="100px"
  >
    <div class="w-100 flex-col">
      <t-space direction="vertical">
        <div class="flx-ce-bet" v-for="item in fields" :key="item.id">
          <div class="flx-ce-sta gap-4" style="width: 30%">
            <t-icon class="field-icon" :name="item.iconName" />{{ item.title }}
          </div>
          <div class="flex1">
            <template v-if="isSpecialField(item)">
              <component
                :is="isActiveComponent(item?.type || '')"
                :field="item"
                :item="itemObj"
                :isEditor="true"
                @handle-special-change="handleSpecialChange"
              />
            </template>
            <template v-else>
              <component
                :is="isActiveComponent(item?.type || '')"
                v-model:editingValue="item.defaultValue"
                :field="item"
                :item="itemObj"
              />
            </template>
          </div>
        </div>
      </t-space>
    </div>
  </t-drawer>
</template>

<script setup lang="ts">
import { useMtTableStore } from '@/stores/mtTable'
import { ObjType } from '@/types'
import { isActiveComponent, isSpecialField } from '@/utils/mtable'
import { cloneDeep } from 'lodash'
import { updateRowRecordApi } from '@/api'

interface Props {
  visible: boolean
}
const props = defineProps<Props>()
const emit = defineEmits(['update:visible'])
const mtTableStore = useMtTableStore()

// 获取字段配置
const fields = ref(cloneDeep(mtTableStore.getFields()))
console.log('fields', fields)
const visible = computed(() => props.visible)

// itemObj 为 fields id值为key，defaultValue 字段值为value
const itemObj = computed(() => {
  return fields.value.reduce((acc, cur) => {
    acc[cur.id] = cur.defaultValue
    // 还需要将creator 也添加进去 value 为creator 对应的值
    acc['creator'] = mtTableStore.getRecords()[0]?.creator
    return acc
  }, {})
})

// 复选框类型字段、图片类型字段、附件类型字段、评分类型字段改变回调事件
const handleSpecialChange = (obj: ObjType) => {
  itemObj.value[obj.field.id] = obj.value
}

// 提交回调事件
const handleConfirm = async () => {
  console.log('itemObj.value', itemObj.value)

  await updateRowRecordApi({
    rowId: mtTableStore.getRecords()[0]?.rowId || '',
    ...Object.keys(itemObj.value).reduce((acc, key) => {
      if (key !== 'creator') {
        acc[key] = itemObj.value[key]
      }
      return acc
    }, {}),
  })
  emit('update:visible', false)
  mtTableStore.updateRecord(mtTableStore.getRecords()[0]?.rowId, itemObj.value)
}

// 抽屉关闭回调事件
const handleCancel = () => {
  emit('update:visible', false)
}

// const handleConfirm = () => {
//   console.log('fields.value', fields.value)
//   emit('update:visible', false)
//   mtTableStore.insertRecord({ ...itemObj.value, rowId: Math.floor(new Date().getTime() / 1000) })
// }

defineOptions({ name: 'AddRecord' })
</script>

<style lang="less" scoped>
.field-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}
</style>
