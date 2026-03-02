<!--
 * @Author: 小磊
 * @Description: 单向关联、双向关联
-->
<template>
  <t-space class="w-100 flx-ce-sta gap-4" direction="vertical">
    <div class="flx-ce-sta gap-4">
      <div>选择要关联的表</div>
      <t-tooltip content="选择关联一个数据表，就能在单元格里添加关联数据表的记录">
        <t-icon class="gray-col" name="helpCircle" />
      </t-tooltip>
    </div>
    <div class="w-100">
      <t-select
        v-model="formData.settings.sourceRef.tableSchemaId"
        placeholder="请选择关联的表"
        :options="tableData"
      />
    </div>
    <div>关联配置</div>
    <t-card size="small">
      <t-space size="small" direction="vertical">
        <t-checkbox v-model="formData.settings.conditions.allowMultiSelect"
          >允许多选记录</t-checkbox
        >
        <!-- <t-checkbox v-model="formData.settings.conditions.fromViewFilter">从视图筛选记录</t-checkbox> -->
      </t-space>
    </t-card>
  </t-space>
</template>

<script setup lang="ts">
import { getProjectTables } from '@/api'

const props = defineProps({
  formData: {
    type: Object,
    required: true,
  },
})

const route = useRoute()
const { formData } = toRefs(props)
defineOptions({ name: 'AssociationComponent' })

const tableData = ref([]) // 表数据

// 获取项目表数据
const getTables = async () => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const { list }: any = await getProjectTables({ projectId: route.params.id as string })
  tableData.value = list.map((item) => ({ label: item.name, value: item.id }))
}

getTables()
</script>

<style scoped></style>
