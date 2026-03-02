<template>
  <div class="field-preview">
    <!-- 文本输入 -->
    <t-input
      v-if="field.type === 'text'"
      :placeholder="`请输入${field.config.description || '内容'}`"
      :disabled="readonly"
    />

    <!-- 富文本 -->
    <t-textarea
      v-else-if="field.type === 'richText'"
      :placeholder="`请输入${field.config.description || '内容'}`"
      :disabled="readonly"
      :autosize="{ minRows: 3, maxRows: 6 }"
    />

    <!-- 单选 -->
    <t-radio-group
      v-else-if="field.type === 'singleChoice'"
      :disabled="readonly"
    >
      <t-radio
        v-for="option in field.config.options"
        :key="option.value"
        :value="option.value"
      >
        {{ option.label }}
      </t-radio>
    </t-radio-group>

    <!-- 多选 -->
    <t-checkbox-group
      v-else-if="field.type === 'multipleChoice'"
      :disabled="readonly"
    >
      <t-checkbox
        v-for="option in field.config.options"
        :key="option.value"
        :value="option.value"
      >
        {{ option.label }}
      </t-checkbox>
    </t-checkbox-group>

    <!-- 日期 -->
    <t-date-picker
      v-else-if="field.type === 'date'"
      :placeholder="`选择${field.config.description || '日期'}`"
      :disabled="readonly"
    />

    <!-- 数字 -->
    <t-input-number
      v-else-if="field.type === 'number'"
      :placeholder="`请输入${field.config.description || '数字'}`"
      :disabled="readonly"
      :min="field.config.min"
      :max="field.config.max"
    />

    <!-- 人员 -->
    <t-select
      v-else-if="field.type === 'person'"
      :placeholder="`选择${field.config.description || '人员'}`"
      :disabled="readonly"
    />

    <!-- 部门 -->
    <t-select
      v-else-if="field.type === 'department'"
      :placeholder="`选择${field.config.description || '部门'}`"
      :disabled="readonly"
    />

    <!-- 链接 -->
    <t-input
      v-else-if="field.type === 'link'"
      :placeholder="`请输入${field.config.description || '链接'}`"
      :disabled="readonly"
    />

    <!-- 电话 -->
    <t-input
      v-else-if="field.type === 'phone'"
      :placeholder="`请输入${field.config.description || '电话'}`"
      :disabled="readonly"
    />

    <!-- 地理位置 -->
    <t-input
      v-else-if="field.type === 'location'"
      :placeholder="`请输入${field.config.description || '位置'}`"
      :disabled="readonly"
    />

    <!-- 行政区域 -->
    <t-cascader
      v-else-if="field.type === 'area'"
      :placeholder="`选择${field.config.description || '区域'}`"
      :disabled="readonly"
    />

    <!-- 评分 -->
    <t-rate
      v-else-if="field.type === 'rating'"
      :disabled="readonly"
    />

    <!-- 复选框 -->
    <t-checkbox
      v-else-if="field.type === 'checkbox'"
      :disabled="readonly"
    >
      {{ field.config.description || '同意' }}
    </t-checkbox>

    <!-- 图片和附件 -->
    <t-upload
      v-else-if="field.type === 'file'"
      theme="file"
      :disabled="readonly"
    >
      <t-button variant="outline">
        <template #icon>
          <t-icon name="upload" />
        </template>
        上传文件
      </t-button>
    </t-upload>

    <!-- 默认 -->
    <t-input
      v-else
      :placeholder="`请输入${field.config.description || '内容'}`"
      :disabled="readonly"
    />
  </div>
</template>

<script setup lang="ts">
interface Props {
  field: any
  readonly?: boolean
}

defineProps<Props>()
</script>

<style scoped lang="less">
.field-preview {
  width: 100%;
}
</style>
