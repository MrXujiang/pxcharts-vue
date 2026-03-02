<template>
  <div class="w-100">
    <t-space size="small" class="flex w-100" direction="vertical">
      <div>进度选项</div>
      <t-card size="small" class="w-100">
        <!-- 数字格式和小数点位数 -->
        <div class="option-row flx-ce-bet">
          <div class="flex-col option-item gap-4">
            <div>数字格式</div>
            <t-select
              v-model="formData.settings.numberFormat"
              :options="formatOptions"
              style="width: 90%"
            />
          </div>
          <div class="flex-col option-item gap-4">
            <div>小数点位数</div>
            <t-select
              v-model="formData.settings.decimalPlaces"
              :options="decimalOptions"
              style="width: 90%"
            />
          </div>
        </div>

        <!-- 颜色选择 -->
        <div class="option-row">
          <div>颜色</div>
          <t-select v-model="formData.settings.color" style="width: 100%">
            <template #valueDisplay>
              <div class="flx-ce-bet gap-8">
                <div
                  class="color-preview"
                  :style="getColorPreviewStyle(formData.settings.color)"
                ></div>
                <span>{{ getColorName(formData.settings.color) }}</span>
              </div>
            </template>
            <t-option
              v-for="color in colorOptions"
              :key="color.value"
              :value="color.value"
              :label="color.label"
            >
              <div class="flx-ce-bet gap-8">
                <!-- 根据colors字段渲染颜色 -->
                <div class="multi-color-option">
                  <div
                    v-for="(c, index) in color.colors"
                    :key="index"
                    class="multi-color-segment"
                    :style="{ backgroundColor: c, width: `${100 / (color?.colors?.length || 1)}%` }"
                  ></div>
                </div>
                <span>{{ color.label }}</span>
              </div>
            </t-option>
          </t-select>
        </div>

        <!-- 自定义进度条值开关 -->
        <div class="option-row flx-ce-bet">
          <span>自定义进度条值</span>
          <t-switch v-model="formData.settings.enableCustomProgress" />
        </div>

        <!-- 起始值和默认值 -->
        <!-- <div v-if="formData.settings.enableCustomProgress" class="flex-col">
          <div class="option-row flx-ce-bet">
            <div class="flex-col option-item gap-4">
              <span>起始值</span>
              <t-input-number
                v-model="formData.settings.startValue"
                theme="normal"
                style="width: 90%"
              >
                <template v-if="formData.settings.numberFormat == 'percentage'" #suffix
                  ><span>%</span></template
                >
              </t-input-number>
            </div>
            <div class="flex-col option-item gap-4">
              <span>目标值</span>
              <t-input-number
                v-model="formData.settings.targetValue"
                theme="normal"
                style="width: 90%"
              >
                <template v-if="formData.settings.numberFormat == 'percentage'" #suffix
                  ><span>%</span></template
                >
              </t-input-number>
            </div>
          </div>
          <div
            style="color: red"
            v-if="formData.settings.targetValue <= formData.settings.startValue"
          >
            目标值应大于起始值
          </div>
        </div> -->
      </t-card>
      <div class="w-100 flex-col gap-8">
        <div>默认数字</div>
        <t-input-number
          v-model="formData.defaultValue"
          placeholder="请输入数字"
          theme="normal"
          style="width: 100%"
          :decimal-places="formData.settings.decimalPlaces"
          :key="formData.settings.decimalPlaces"
        >
          <template v-if="formData.settings.numberFormat == 'percentage'" #suffix
            ><span>%</span></template
          >
        </t-input-number>
      </div>
    </t-space>
  </div>
</template>

<script setup lang="ts">
import { formatOptions, decimalOptions, colorOptions } from './options'

const props = defineProps({
  formData: {
    type: Object,
    required: true,
  },
})
const { formData } = toRefs(props)

// 获取颜色预览样式
const getColorPreviewStyle = (color: string) => {
  // 查找颜色对象
  const colorObj = colorOptions.find((c) => c.value === color)
  if (!colorObj) {
    return { backgroundColor: color }
  }

  // 如果有colors字段
  if (colorObj.colors && colorObj.colors.length > 0) {
    // 如果只有一个颜色，使用纯色背景
    if (colorObj.colors.length === 1) {
      return { backgroundColor: colorObj.colors[0] }
    }
    // 如果有多个颜色，创建与选项中一致的多色背景（使用渐变实现）
    else {
      // 生成渐变色定义，每个颜色平均分配位置
      const colorStops = colorObj.colors
        .map((c, index) => {
          const start = (index / colorObj.colors!.length) * 100
          const end = ((index + 1) / colorObj.colors!.length) * 100
          return `${c} ${start}%, ${c} ${end}%`
        })
        .join(', ')

      return { backgroundImage: `linear-gradient(to right, ${colorStops})` }
    }
  }
  // 如果没有colors字段，使用value作为背景色
  else {
    return { backgroundColor: colorObj.value }
  }
}

// 获取颜色名称
const getColorName = (color: string) => {
  const colorObj = colorOptions.find((c) => c.value === color)
  return colorObj ? colorObj.label : ''
}

defineOptions({ name: 'ProgressComponent' })
</script>

<style lang="less" scoped>
@import './index.less';
</style>
