<template>
  <div>
    <t-popup>
      <t-tag size="small" shape="round">
        <template #content>
          <div class="flx-ce-sta gap-4">
            <img class="img" :src="item.creator?.avatar" alt="" />
            <div style="font-size: 10px">{{ item.creator?.nickname }}</div>
          </div>
        </template>
      </t-tag>
      <template #content>
        <div class="flex-col gap-4 popup-content">
          <t-avatar size="42px" :image="item.creator?.avatar" />
          <p style="font-size: 16px; font-weight: 500">{{ item.creator?.nickname }}</p>
          <div>
            <t-button variant="outline">@ 提醒他关注记录</t-button>
          </div>
        </div>
      </template>
    </t-popup>
  </div>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'

interface Props {
  field: ObjType
  item: ObjType
  editingValue: string | number
}

const props = defineProps<Props>()
defineOptions({ name: 'CreatorComponent' })
</script>

<style lang="less" scoped>
.img {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  object-fit: cover;
}
.popup-content {
  width: 320px;
  padding: 20px 12px;
  align-items: center;
}
</style>
