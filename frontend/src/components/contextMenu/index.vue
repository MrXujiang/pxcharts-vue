<template>
  <div
    v-if="visible"
    class="context-menu"
    :style="{ left: position.x + 'px', top: position.y + 'px' }"
  >
    <div
      v-for="item in menuItems"
      :key="item.key"
      class="context-menu-item"
      @click="handleItemClick(item)"
    >
      <t-space size="small" class="flx-ce-ce">
        <span class="flx-ce-ce">
          <t-icon :name="item.iconName" />
        </span>
        <span>{{ item.label }}</span>
      </t-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

defineOptions({
  name: 'ContextMenu',
})

export interface MenuItem {
  key: string | number
  label: string
  iconName: string
  disabled?: boolean
}

// 定义props
interface Props {
  menuItems: MenuItem[]
  visible?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
})

// 定义事件
interface Emits {
  (e: 'select', item: MenuItem): void
  (e: 'close'): void
}

const emit = defineEmits<Emits>()

// 菜单位置
const position = ref({ x: 0, y: 0 })

// 设置菜单位置
const setPosition = (x: number, y: number) => {
  position.value = { x, y }
}

// 处理菜单项点击
const handleItemClick = (item: MenuItem) => {
  if (!item.disabled) {
    emit('select', item)
    emit('close')
  }
}

// 监听点击事件，隐藏菜单
onMounted(() => {
  const handleClick = () => {
    if (props.visible) {
      emit('close')
    }
  }

  window.addEventListener('click', handleClick)

  // 组件销毁时移除事件监听
  onUnmounted(() => {
    window.removeEventListener('click', handleClick)
  })
})

// 暴露方法给父组件
defineExpose({
  setPosition,
})
</script>

<style lang="less" scoped>
.context-menu {
  position: fixed;
  background: #ffffff;
  border: 1px solid #e5e5e5;
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 9999;
  min-width: 160px;
  padding: 4px 8px;
  box-sizing: border-box;
}

.context-menu-item {
  padding: 8px 16px;
  cursor: pointer;
  font-size: 14px;
  color: #333;
  transition: background-color 0.2s;
  border-radius: 8px;
}

.context-menu-item:hover {
  background-color: #f5f5f5;
}

.context-menu-item:active {
  background-color: #ebebeb;
}

.context-menu-item[disabled] {
  color: #ccc;
  cursor: not-allowed;
  background-color: transparent;
}
</style>
