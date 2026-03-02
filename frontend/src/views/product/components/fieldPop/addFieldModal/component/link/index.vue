<template>
  <div class="w-100">
    <t-space class="w-100" size="small" direction="vertical">
      <!-- 顶部预览栏 -->
      <div class="flx-ce-bet nav-bar">
        <div>默认链接</div>
        <div>
          <t-checkbox v-model="formData.settings.isLinkTitle" @change="handleIsLinkTitleChange"
            >链接标题</t-checkbox
          >
        </div>
      </div>
      <t-card size="small">
        <t-space class="w-100" size="small" direction="vertical">
          <div v-if="formData.settings.isLinkTitle" class="flex-col gap-8">
            <div>链接标题</div>
            <t-input v-model="formData.defaultValue.linkTitle" placeholder="请输入链接名称" />
          </div>
          <div class="flex-col gap-8">
            <div>链接地址</div>
            <t-input v-model="formData.defaultValue.linkUrl" placeholder="请输入链接地址" />
          </div>
        </t-space>
      </t-card>
    </t-space>
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  formData: {
    type: Object,
    required: true,
  },
})
const { formData } = toRefs(props)

const handleIsLinkTitleChange = () => {
  formData.value.defaultValue.linkTitle = ''
}

defineOptions({ name: 'LinkComponent' })
</script>

<style scoped></style>
