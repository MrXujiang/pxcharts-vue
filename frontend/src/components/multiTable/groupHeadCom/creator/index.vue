<template>
  <div>
    <template v-if="creatorInfo?.id">
      <t-tag size="small" shape="round">
        <template #content>
          <div class="flx-ce-sta gap-4">
            <img class="img" :src="creatorInfo?.avatar" alt="" />
            <div style="font-size: 10px">{{ creatorInfo?.nickname }}</div>
          </div>
        </template>
      </t-tag>
    </template>
    <span class="noGroup" v-else>无创建人</span>
  </div>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'

interface Props {
  value: ObjType
}

const props = defineProps<Props>()

const creatorInfo = computed(() => props.value)

defineOptions({
  name: 'CreatorGroupHeadCom',
})
</script>

<style scoped lang="less">
.img {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  object-fit: cover;
}
.noGroup {
  font-size: 12px;
  color: #495057;
}
</style>
