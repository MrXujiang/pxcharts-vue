<!--
 * @Author: 小磊
 * @Description: 单选框、多选框组件
-->
<template>
  <div class="w-100">
    <t-space class="w-100" size="small" direction="vertical">
      <div class="flx-ce-bet nav-bar">
        <div>选项内容</div>
        <div class="flx-ce-sta" style="gap: 8px">
          <t-tooltip content="批量设置颜色" placement="top">
            <span class="flex pointer" @click="batchSetColor">
              <svg width="16" height="16" fill="none" viewBox="0 0 16 16" class="_colorsslect16">
                <defs>
                  <clipPath id="b8e39f_0"><rect width="16" height="16" rx="0"></rect></clipPath>
                </defs>
                <g clip-path="url(#b8e39f_0)">
                  <path
                    fill="#FF5938"
                    d="M12.442 5.335q0 .109-.005.218-.006.109-.016.217-.01.109-.027.217-.016.107-.037.214-.021.107-.048.213-.026.106-.058.21-.032.105-.069.207-.036.103-.078.204-.042.1-.088.199-.047.099-.098.195-.052.096-.108.19-.056.093-.116.184-.061.09-.126.178-.065.087-.134.172-.07.084-.143.165-.073.08-.15.158-.077.077-.158.15-.08.073-.165.143-.084.069-.172.134-.087.065-.178.125-.09.06-.184.117-.094.056-.19.107-.096.052-.195.098-.098.047-.2.089-.1.041-.202.078-.103.037-.208.069-.104.031-.21.058-.105.026-.212.047-.107.022-.215.038t-.217.026q-.108.011-.217.017-.109.005-.218.005t-.218-.005q-.109-.006-.217-.017-.109-.01-.217-.026t-.215-.038q-.107-.02-.212-.047-.106-.027-.21-.058-.105-.032-.208-.069-.102-.037-.203-.078-.1-.042-.2-.089-.098-.046-.194-.098-.096-.051-.19-.107-.093-.056-.184-.117-.09-.06-.178-.125t-.172-.134q-.084-.07-.165-.143-.081-.073-.158-.15t-.15-.158q-.074-.08-.143-.165-.069-.085-.134-.172-.065-.088-.126-.178-.06-.091-.116-.185-.056-.093-.108-.19-.051-.095-.098-.194-.046-.099-.088-.2-.042-.1-.078-.203-.037-.102-.069-.207-.032-.104-.058-.21-.027-.106-.048-.213t-.037-.214q-.016-.108-.027-.217-.01-.108-.016-.217-.005-.11-.005-.218 0-.11.005-.218.006-.11.016-.218t.027-.216q.016-.108.037-.215.021-.107.048-.213.026-.105.058-.21.032-.104.069-.207.036-.102.078-.203.042-.101.088-.2.047-.098.098-.194.052-.096.108-.19t.116-.184q.061-.09.126-.178t.134-.172q.07-.085.143-.165.073-.081.15-.158t.158-.15q.08-.074.165-.143.084-.07.172-.134.087-.065.178-.126.09-.06.184-.116.094-.056.19-.108.096-.051.195-.098.098-.047.2-.088.1-.042.202-.079.103-.036.208-.068.104-.032.21-.058.105-.027.212-.048T7.348.94q.108-.016.217-.027.108-.01.217-.016Q7.891.893 8 .893q.11 0 .218.005.109.005.217.016.109.01.217.027.108.016.215.037.107.021.212.048.106.026.21.058.105.032.208.068.102.037.203.079.1.041.2.088.098.047.194.098.096.052.19.108.093.056.184.116.09.06.178.126.088.065.172.134.084.07.165.142.081.074.158.15.077.078.15.159.074.08.143.165.069.084.134.172.065.087.126.178.06.09.116.184.056.094.108.19.051.096.098.194.046.099.088.2.042.1.078.203.037.103.069.207.032.105.058.21.027.106.048.213t.037.215q.016.108.027.216.01.109.016.218.005.109.005.218Z"
                    style="opacity: 0.6"
                  ></path>
                  <path
                    fill="#00B23B"
                    d="M9.777 10.665q0 .11-.005.218-.006.11-.017.218-.01.108-.026.216t-.038.215q-.02.107-.047.213-.027.105-.058.21-.032.104-.069.207-.037.102-.078.203-.042.101-.089.2-.046.098-.098.194-.051.096-.107.19t-.117.184q-.06.09-.125.178t-.134.172q-.07.085-.143.165-.073.081-.15.158t-.158.15q-.08.074-.165.143-.085.07-.172.134-.088.065-.178.126-.091.06-.185.116-.093.056-.19.108-.095.051-.194.098-.099.047-.2.088-.1.042-.203.079-.102.036-.207.068-.104.032-.21.058-.106.027-.213.048t-.214.037q-.108.016-.217.027-.108.01-.217.016-.11.005-.218.005-.11 0-.218-.005-.11-.005-.218-.016-.108-.01-.216-.027-.108-.016-.215-.037-.107-.021-.213-.048-.105-.026-.21-.058-.104-.032-.207-.068-.102-.037-.203-.079-.101-.041-.2-.088-.098-.047-.194-.098-.096-.052-.19-.108t-.184-.116q-.09-.06-.178-.126-.088-.065-.172-.134-.085-.07-.165-.142-.081-.074-.158-.15-.077-.078-.15-.159-.074-.08-.143-.165-.07-.084-.134-.172-.065-.087-.126-.178-.06-.09-.116-.184-.056-.094-.108-.19-.051-.096-.098-.194-.047-.099-.088-.2-.042-.1-.079-.203-.036-.103-.068-.207-.032-.105-.058-.21-.027-.106-.048-.213t-.037-.215q-.016-.108-.027-.216-.01-.109-.016-.218-.005-.109-.005-.218t.005-.218q.005-.109.016-.217.01-.109.027-.217.016-.107.037-.214.021-.107.048-.213.026-.106.058-.21.032-.105.068-.207.037-.103.079-.204.041-.1.088-.199.047-.099.098-.195.052-.096.108-.19.056-.093.116-.184.06-.09.126-.178.065-.087.134-.172.07-.084.142-.165.074-.08.15-.158.078-.077.159-.15.08-.073.165-.143.084-.069.172-.134.087-.065.178-.125.09-.06.184-.117.094-.056.19-.107.096-.052.194-.098.099-.047.2-.089.1-.041.203-.078.103-.037.207-.069.105-.031.21-.058.106-.026.213-.048.107-.02.215-.037.108-.016.216-.026.109-.011.218-.017.109-.005.218-.005t.218.005q.109.006.217.017.109.01.217.026.107.016.214.037.107.022.213.048.106.027.21.058.105.032.207.069.103.037.204.078.1.042.199.089.099.046.195.098.096.051.19.107.093.056.184.117.09.06.178.125.087.065.172.134.084.07.165.143.08.073.158.15.077.077.15.158.073.08.143.165.069.085.134.172.065.088.125.178.06.091.117.185.056.093.107.19.052.095.098.194.047.099.089.2.041.1.078.203.037.102.069.207.031.104.058.21.026.106.047.213.022.107.038.214.016.108.026.217.011.108.017.217.005.11.005.218Z"
                    style="opacity: 0.6"
                  ></path>
                  <path
                    fill="#FFE747"
                    d="M15.107 10.665q0 .11-.005.218-.005.11-.016.218-.01.108-.027.216-.016.108-.037.215-.021.107-.048.213-.026.105-.058.21-.032.104-.068.207-.037.102-.079.203-.041.101-.088.2-.047.098-.098.194-.052.096-.108.19t-.116.184q-.06.09-.126.178-.065.088-.134.172-.07.085-.142.165-.074.081-.15.158-.078.077-.159.15-.08.074-.165.143-.084.07-.172.134-.087.065-.178.126-.09.06-.184.116-.094.056-.19.108-.096.051-.194.098-.099.047-.2.088-.1.042-.203.079-.103.036-.207.068-.105.032-.21.058-.106.027-.213.048t-.215.037q-.108.016-.216.027-.109.01-.218.016-.109.005-.218.005t-.218-.005q-.109-.005-.217-.016-.109-.01-.217-.027-.107-.016-.214-.037-.107-.021-.213-.048-.106-.026-.21-.058-.105-.032-.207-.068-.103-.037-.204-.079-.1-.041-.199-.088-.099-.047-.195-.098-.096-.052-.19-.108-.093-.056-.184-.116-.09-.06-.178-.126-.087-.065-.172-.134-.084-.07-.165-.142-.08-.074-.158-.15-.077-.078-.15-.159-.073-.08-.143-.165-.069-.084-.134-.172-.065-.087-.125-.178-.06-.09-.117-.184-.056-.094-.107-.19-.052-.096-.098-.194-.047-.099-.089-.2-.041-.1-.078-.203-.037-.103-.069-.207-.031-.105-.058-.21-.026-.106-.048-.213-.02-.107-.037-.215-.016-.108-.026-.216-.011-.109-.017-.218-.005-.109-.005-.218t.005-.218q.006-.109.017-.217.01-.109.026-.217.016-.107.037-.214.022-.107.048-.213.027-.106.058-.21.032-.105.069-.207.037-.103.078-.204.042-.1.089-.199.046-.099.098-.195.051-.096.107-.19.056-.093.117-.184.06-.09.125-.178.065-.087.134-.172.07-.084.143-.165.073-.08.15-.158.077-.077.158-.15.08-.073.165-.143.085-.069.172-.134.088-.065.178-.125.091-.06.185-.117.093-.056.19-.107.095-.052.194-.098.099-.047.2-.089.1-.041.203-.078.102-.037.207-.069.104-.031.21-.058.106-.026.213-.048.107-.02.214-.037.108-.016.217-.026.108-.011.217-.017.11-.005.218-.005.11 0 .218.005.11.006.218.017.108.01.216.026t.215.037q.107.022.213.048.105.027.21.058.104.032.207.069.102.037.203.078.101.042.2.089.098.046.194.098.096.051.19.107t.184.117q.09.06.178.125t.172.134q.085.07.165.143.081.073.158.15t.15.158q.074.08.143.165.07.085.134.172.065.088.126.178.06.091.116.185.056.093.108.19.051.095.098.194.047.099.088.2.042.1.079.203.036.102.068.207.032.104.058.21.027.106.048.213t.037.214q.016.108.027.217.01.108.016.217.005.11.005.218Z"
                    style="opacity: 0.6"
                    data-spm-anchor-id="a2q1e.b45810666.0.i204.7251Y8SgY8SgsV"
                  ></path>
                </g>
              </svg>
            </span>
          </t-tooltip>

          <t-divider style="margin: 0 4px" layout="vertical" />
          <t-checkbox v-model="formData.settings.isReferenced">引用选项</t-checkbox>
          <t-tooltip content="可引用当前数据表，或其他数据表的单选、多选字段的选项">
            <t-icon class="icon pointer" name="helpCircle" />
          </t-tooltip>
        </div>
      </div>
      <!-- 引用选项 -->
      <CitationOptions v-if="formData.settings.isReferenced" />
      <!-- 单选框选项 -->
      <t-card v-else size="small" class="w-100">
        <div ref="sortableContainer" class="sortable-container">
          <div
            v-for="(item, index) in formData.options"
            :key="item.id"
            class="flx-ce-sta sortable-item"
          >
            <!-- 拖拽手柄 -->
            <div class="move-handle flx-ce-ce">
              <t-icon name="gripVertical" />
            </div>
            <!-- 颜色选择器 -->
            <t-popup placement="bottom-left" trigger="click">
              <div class="circle pointer" :style="{ backgroundColor: item.color }"></div>
              <template #content>
                <ColorPanel
                  v-model:col-val="item.color"
                  @change="(color) => updateColor(index, color)"
                />
              </template>
            </t-popup>
            <!-- 输入框 -->
            <div class="flex1 content">
              <t-input v-model="item.label" placeholder="请输入内容" />
            </div>
            <!-- 关闭按钮 -->
            <div class="closeIcon pointer" @click="removeOption(index)">
              <t-icon name="close" />
            </div>
          </div>
        </div>

        <div class="flx-ce-ce w-100" style="gap: 20px">
          <!-- 添加选项按钮 -->
          <div class="flx-ce-ce pointer add-option" @click="addOption">
            <t-icon name="add" />
            <span>添加选项</span>
          </div>
          <!-- 批量编辑 -->
          <div class="flx-ce-ce pointer add-option" @click="openBatchEdit">
            <t-icon name="edit" />
            <span>批量编辑</span>
          </div>
        </div>
      </t-card>
      <div class="flex-col gap-8">
        <div>默认选项</div>
        <t-select
          v-model="formData.defaultValue"
          placeholder="请选择"
          :popup-props="{ overlayClassName: 'auto-number' }"
          :multiple="formData.type === 'radio' ? false : true"
        >
          <template #valueDisplay="{ value, onClose, label }">
            <template v-if="formData.type != 'radio'">
              <t-tag
                v-for="(item, index) in value"
                :key="index"
                shape="round"
                :closable="true"
                :color="formData.options.find((option) => option.id === item.value)?.color"
                :on-close="
                  ({ e }) => {
                    e.stopPropagation()
                    onClose(index)
                  }
                "
              >
                {{ item.label }}
              </t-tag>
            </template>
            <template v-else>
              <t-tag
                closable
                shape="round"
                :on-close="
                  ({ e }) => {
                    e.stopPropagation()
                    formData.defaultValue = ''
                  }
                "
                :color="formData.options.find((option) => option.id === value)?.color"
                >{{ label }}</t-tag
              >
            </template>
          </template>
          <t-option
            v-for="item in formData.options"
            :key="item.id"
            :value="item.id"
            :label="item.label"
          >
            <t-tag size="small" :color="item.color" shape="round">{{ item.label }}</t-tag>
          </t-option>
        </t-select>
      </div>
    </t-space>
    <t-dialog
      width="600px"
      :showOverlay="false"
      v-model:visible="batchEditVisible"
      @cancel="handleCancel"
      @confirm="handleBatchEditConfirm"
    >
      <template #header>
        <div>
          批量编辑
          <div style="font-weight: normal; font-size: 14px; color: #6e7286">
            粘贴或输入选项，一行文字生成一个选项
          </div>
        </div>
      </template>
      <t-card size="small">
        <t-textarea
          class="textarea"
          v-model="batchEditValue"
          name="description"
          :autosize="{ minRows: 10, maxRows: 12 }"
        />
        <span class="text-gray">此字段最多支持1000个选项，已有{{ formData.options.length }}个</span>
      </t-card>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import Sortable from 'sortablejs'
import ColorPanel from '@/components/colorPanel/index.vue'
import CitationOptions from '../citationOptions/index.vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { colorRows } from '@/modal/options'

const props = defineProps({
  formData: {
    type: Object,
    required: true,
  },
})
const { formData } = toRefs(props)
console.log('props-----', props)
defineOptions({ name: 'SingleMultipleChoiceComponent' })

const batchEditVisible = ref<boolean>(false) // 批量编辑弹窗可见性
const batchEditValue = ref<string>('') // 批量编辑的值

// Sortable 实例引用
const sortableContainer = ref<HTMLElement | null>(null)
let sortable: Sortable | null = null

// 取消批量编辑
const handleCancel = () => {
  batchEditVisible.value = false
  batchEditValue.value = ''
}

// 确认批量编辑
const handleBatchEditConfirm = () => {
  if (batchEditValue.value.trim()) {
    // 将输入的内容按行分割
    const lines = batchEditValue.value.split('\n').filter((line) => line.trim() !== '')

    // 限制最多1000个选项
    if (lines.length > 1000) {
      // 添加提示信息
      MessagePlugin.warning('选项数量不能超过1000个')
      return
    }

    // 保存现有的选项
    const existingOptions = [...formData.value.options]

    // 将二维数组扁平化为一维数组
    const flatColors = colorRows.flat()

    // 生成新的选项数组
    const newOptions = lines.map((line, index) => {
      // 如果索引在现有选项范围内，说明是修改现有选项
      if (index < existingOptions.length) {
        // 修改现有选项，保留原有颜色和ID，只更新标签
        return {
          id: existingOptions[index].id, // 保留原有ID
          label: line.trim(), // 更新标签
          color: existingOptions[index].color, // 保留原有颜色
        }
      } else {
        // 如果是新增选项，分配随机颜色
        const randomIndex = Math.floor(Math.random() * flatColors.length)
        return {
          id: Date.now() + index,
          label: line.trim(),
          color: flatColors[randomIndex],
        }
      }
    })

    // 更新选项
    formData.value.options = newOptions
  }

  // 关闭弹窗并清空值
  batchEditVisible.value = false
  batchEditValue.value = ''
}

// 打开批量编辑弹窗时，将当前选项值带入
const openBatchEdit = () => {
  batchEditValue.value = formData.value.options.map((option) => option.label).join('\n')
  batchEditVisible.value = true
}

// 批量设置颜色 - 给全部数据重新设置颜色
const batchSetColor = () => {
  // 将二维数组扁平化为一维数组
  const flatColors = colorRows.flat()

  // 给所有选项重新分配随机颜色
  formData.value.options = formData.value.options.map((option) => {
    const randomIndex = Math.floor(Math.random() * flatColors.length)
    return {
      ...option,
      color: flatColors[randomIndex],
    }
  })
}

// 更新指定选项的颜色
const updateColor = (index: number, color: string) => {
  formData.value.options[index].color = color
}

// 添加新选项
const addOption = () => {
  // 将二维数组扁平化为一维数组
  const flatColors = colorRows.flat()
  // 随机选择一个颜色
  const randomIndex = Math.floor(Math.random() * flatColors.length)

  formData.value.options.push({
    id: Date.now(),
    label: '',
    color: flatColors[randomIndex],
  })
}

// 删除指定选项
const removeOption = (index: number) => {
  if (formData.value.options.length > 1) {
    formData.value.options.splice(index, 1)
  }
}

// 初始化 Sortable
const initSortable = () => {
  if (sortableContainer.value) {
    // 如果已存在 sortable 实例，先销毁
    if (sortable) {
      sortable.destroy()
    }

    // 创建新的 sortable 实例
    sortable = new Sortable(sortableContainer.value, {
      animation: 150,
      ghostClass: 'sortable-ghost',
      chosenClass: 'sortable-chosen',
      handle: '.move-handle',
      onEnd: handleDragEnd,
    })
  }
}

// 拖拽结束
const handleDragEnd = (evt: Sortable.SortableEvent) => {
  const { oldIndex, newIndex } = evt
  if (oldIndex === newIndex || oldIndex === undefined || newIndex === undefined) return

  const newOptions = [...formData.value.options]
  const movedItem = newOptions.splice(oldIndex, 1)[0]
  newOptions.splice(newIndex, 0, movedItem)

  formData.value.options = newOptions
}

// 组件挂载时初始化 Sortable
onMounted(() => {
  nextTick(() => {
    initSortable()
  })
})

// 组件卸载时销毁 Sortable
onUnmounted(() => {
  if (sortable) {
    sortable.destroy()
  }
})
</script>

<style lang="less" scoped>
@import './index.less';
</style>
