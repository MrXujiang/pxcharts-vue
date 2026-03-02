<template>
  <div class="w-100">
    <t-space size="small" class="flex w-100" direction="vertical">
      <div class="flx-ce-bet">
        <div>评分选项</div>
      </div>

      <t-card size="small">
        <t-space class="flex w-100" size="small" direction="vertical">
          <div class="flex-col gap-4">
            <div>图形</div>
            <div class="rate-icons-container">
              <div
                v-for="icon in iconData"
                :key="icon.type"
                class="rate-icon"
                :class="{ selected: selectedIcon === icon.type }"
                @click="selectIcon(icon.type)"
                @mouseenter="handleMouseEnter(icon)"
                @mouseleave="handleMouseLeave"
              >
                <div
                  class="icon-inner"
                  :class="{
                    hover:
                      hoveredIcon && hoveredIcon.type === icon.type && selectedIcon !== icon.type,
                  }"
                >
                  <t-icon
                    v-if="['star', 'heart'].includes(icon.type)"
                    :name="icon.iconName"
                    :color="icon.color"
                  />
                  <span class="flx-ce-ce" v-else>
                    <svg
                      width="20"
                      height="20"
                      fill="none"
                      viewBox="0 0 20 20"
                      class="_rateshapenpscolor16"
                    >
                      <path
                        fill="#2F75FD"
                        fill-rule="evenodd"
                        d="M0 3.75v12.5q0 .184.018.368.018.183.054.364.036.18.09.357.053.176.123.346.071.17.158.333.087.162.189.315.102.154.22.296.116.142.246.273.13.13.273.247.142.117.296.219.153.102.315.19.163.086.333.157.17.07.346.123.177.054.357.09.181.036.364.054.184.018.368.018h12.5q.184 0 .368-.018.183-.018.364-.054.18-.036.357-.09.176-.053.346-.123.17-.071.333-.158.162-.087.315-.189.154-.102.296-.22.142-.116.273-.246.13-.13.247-.273.117-.142.219-.296.102-.153.19-.315.086-.163.157-.333.07-.17.123-.346.054-.177.09-.357.036-.181.054-.364.018-.184.018-.368V3.75q0-.184-.018-.368-.018-.183-.054-.364-.036-.18-.09-.357-.053-.176-.123-.346-.071-.17-.158-.333-.087-.162-.189-.315-.102-.154-.22-.296-.116-.142-.246-.273-.13-.13-.273-.247-.142-.117-.296-.219-.153-.102-.315-.19-.163-.086-.333-.157-.17-.07-.346-.124-.177-.053-.357-.089-.181-.036-.364-.054Q16.434 0 16.25 0H3.75q-.184 0-.368.018-.183.018-.364.054-.18.036-.357.09-.176.053-.346.123-.17.071-.333.158-.162.087-.315.189-.154.102-.296.22-.142.116-.273.246-.13.13-.247.273-.117.142-.219.296-.102.153-.19.315-.086.163-.157.333-.07.17-.124.346-.053.177-.089.357-.036.181-.054.364Q0 3.566 0 3.75Zm6.492 11.036v-.984L9.638 9.3q.567-.84.567-1.546v-.202q0-.587-.337-.963-.347-.345-.904-.345-.563 0-.89.345-.323.341-.323 1.016v.388H6.49v-.441q0-1.11.676-1.784.688-.689 1.824-.689.718 0 1.31.33.585.337.876.877.314.582.314 1.239v.229q0 .6-.212 1.128-.203.528-.623 1.095l-2.53 3.644h3.365v1.165h-5Zm6.49-.496q-.706-.642-.706-1.734v-.36h1.273v.32q0 .606.342.914.362.326.979.326.66 0 1.009-.383.366-.402.366-1.086v-.35q0-.701-.39-1.11-.394-.413-1.053-.413h-.629V9.236h.616q.602 0 .956-.37.352-.39.352-1.018v-.27q0-.6-.338-.977-.34-.357-.916-.357-.571 0-.905.333-.322.34-.322.907v.36h-1.232v-.386q0-1.09.663-1.732.668-.647 1.823-.647 1.114 0 1.81.675.675.675.675 1.77v.297q0 .73-.345 1.248-.257.402-.646.65.482.262.782.757.371.611.371 1.434v.323q0 1.217-.742 1.96-.729.728-1.919.728-1.2 0-1.874-.63ZM4.084 6.78l-1.6 1.197V6.512l1.682-1.298h1.177v9.572h-1.26V6.78Z"
                      ></path>
                    </svg>
                  </span>
                </div>
              </div>
            </div>
          </div>
          <!-- 分值设置 -->
          <div class="flx-ce-bet">
            <div>分值</div>
            <div class="flx-ce-ce gap-4">
              <t-select v-model="formData.settings.scoreStart" style="width: 100px">
                <t-option
                  v-for="option in scoreStartOptions"
                  :key="option.value"
                  :label="option.label"
                  :value="option.value"
                />
              </t-select>
              <div>至</div>
              <t-select v-model="formData.settings.scoreEnd" style="width: 100px">
                <t-option
                  v-for="option in scoreEndOptions"
                  :key="option.value"
                  :label="option.label"
                  :value="option.value"
                />
              </t-select>
            </div>
          </div>

          <!-- 两极文案设置 -->
          <!-- <div class="flx-ce-bet">
            <div>两极文案</div>
            <t-select v-model="formData.settings.text" style="width: 200px">
              <t-option
                v-for="option in textOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              />
            </t-select>
          </div> -->
        </t-space>
      </t-card>
      <div class="w-100 flex-col gap-8">
        <div>默认评分</div>
        <div class="flx-ce-sta" style="border: 1px solid #dfe1e6; padding: 8px; border-radius: 4px">
          <IconRate
            v-model="formData.defaultValue"
            :icon-type="selectedIcon"
            :icon-name="iconData.filter((icon) => icon.type === selectedIcon)[0]?.iconName || ''"
            :count="rateCount"
            :min="formData.settings.scoreStart"
            :max="formData.settings.scoreEnd"
          />
        </div>
      </div>
    </t-space>
  </div>
</template>

<script setup lang="ts">
import { iconData, scoreStartOptions, scoreEndOptions, textOptions, IconData } from './options'
import IconRate from '@/components/fieldRate/index.vue'
defineOptions({ name: 'RateComponent' })

const props = defineProps({
  formData: {
    type: Object,
    required: true,
  },
})
const { formData } = toRefs(props)

const selectedIcon = ref<string>(formData.value.settings.icon)
const hoveredIcon = ref<IconData | null>(null)
const rate = ref<number | null>(null)

// 计算rate的count值，即分值范围内的选项数量
const rateCount = computed(() => {
  // 确保结束值大于等于起始值
  if (formData.value.settings.scoreEnd >= formData.value.settings.scoreStart) {
    // 选项数量 = 结束值 - 起始值 + 1
    return formData.value.settings.scoreEnd - formData.value.settings.scoreStart + 1
  }
  return 5 // 默认返回5
})

// 选择图标
const selectIcon = (type: string) => {
  selectedIcon.value = type
  formData.value.settings.icon = type
}

// 鼠标进入事件
const handleMouseEnter = (icon: IconData) => {
  hoveredIcon.value = icon
}

// 鼠标离开事件
const handleMouseLeave = () => {
  hoveredIcon.value = null
}

watch(selectedIcon, () => {
  rate.value = null
})
watch(formData.value.settings, () => {
  formData.value.defaultValue = null
})
watch(rate, () => {
  console.log('rate.value', rate.value)
})
</script>

<style scoped>
.rate-icons-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.rate-icon {
  width: 28px;
  height: 28px;
  border: 1px solid transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}
.rate-icon .icon-inner {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}
.rate-icon:hover .icon-inner.hover {
  border: 1.5px solid #61a4e8;
  border-radius: 4px;
}
.rate-icon.selected .icon-inner {
  border: 1.5px solid #0052d9;
  border-radius: 4px;
}
</style>
