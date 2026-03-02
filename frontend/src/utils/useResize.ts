import { ref } from 'vue'

/**
 * 自定义Hook：处理页面元素拖拽调整大小功能
 * 主要用于左侧侧边栏和右侧面板的宽度调整
 */
export function useResize() {
  // 左侧侧边栏宽度，默认280px
  const sidebarWidth = ref(280)
  // 中间内容区域分割比例，默认0.58
  const contentSplitRatio = ref(0.58)
  // 右侧面板是否可见，默认显示
  const rightPanelVisible = ref(true)

  // 拖拽状态标识
  let isLeftResizing = false // 左侧是否正在拖拽调整
  let isRightResizing = false // 右侧是否正在拖拽调整
  let startX = 0 // 鼠标按下时的起始X坐标
  let startSidebarWidth = 0 // 拖拽开始时的侧边栏宽度
  let startContentRatio = 0 // 拖拽开始时的内容分割比例

  // 存储事件处理函数引用，方便移除监听器
  let leftMoveHandler: (e: MouseEvent) => void
  let rightMoveHandler: (e: MouseEvent) => void

  /**
   * 开始左侧拖拽调整
   * @param e 鼠标事件对象
   */
  const startLeftResize = (e: MouseEvent) => {
    e.preventDefault()
    isLeftResizing = true
    startX = e.clientX
    startSidebarWidth = sidebarWidth.value

    // 鼠标移动事件处理函数
    leftMoveHandler = (e: MouseEvent) => {
      if (!isLeftResizing) return
      const deltaX = e.clientX - startX
      // 限制侧边栏宽度在180px到400px之间
      sidebarWidth.value = Math.max(180, Math.min(400, startSidebarWidth + deltaX))
    }

    // 添加全局鼠标事件监听器
    document.addEventListener('mousemove', leftMoveHandler)
    document.addEventListener('mouseup', stopResize)
    // 设置鼠标样式为列调整样式
    document.body.style.cursor = 'col-resize'
    document.body.style.userSelect = 'none'
  }

  /**
   * 开始右侧拖拽调整
   * @param e 鼠标事件对象
   * @param chatPanelRef 右侧面板DOM引用
   */
  const startRightResize = (e: MouseEvent, chatPanelRef: HTMLElement | null) => {
    console.log('rightPanelVisible.value', rightPanelVisible.value)
    // 如果右侧面板不可见或DOM引用不存在，则不执行调整
    if (!rightPanelVisible.value || !chatPanelRef) return
    e.preventDefault()
    isRightResizing = true
    startX = e.clientX
    startContentRatio = contentSplitRatio.value

    // 获取容器元素
    const container = chatPanelRef.parentElement
    if (!container) return

    // 鼠标移动事件处理函数
    rightMoveHandler = (e: MouseEvent) => {
      if (!isRightResizing) return
      const containerWidth = container.clientWidth
      const deltaX = e.clientX - startX
      // 计算新的分割比例，限制在0.3到0.7之间
      const newRatio = startContentRatio + deltaX / containerWidth
      contentSplitRatio.value = Math.max(0.3, Math.min(0.7, newRatio))
    }

    // 添加全局鼠标事件监听器
    document.addEventListener('mousemove', rightMoveHandler)
    document.addEventListener('mouseup', stopResize)
    // 设置鼠标样式为列调整样式
    document.body.style.cursor = 'col-resize'
    document.body.style.userSelect = 'none'
  }

  /**
   * 停止拖拽调整
   * 清理所有事件监听器并恢复默认样式
   */
  const stopResize = () => {
    isLeftResizing = false
    isRightResizing = false

    // 移除鼠标移动事件监听器
    if (leftMoveHandler) {
      document.removeEventListener('mousemove', leftMoveHandler)
    }
    if (rightMoveHandler) {
      document.removeEventListener('mousemove', rightMoveHandler)
    }

    // 移除鼠标释放事件监听器
    document.removeEventListener('mouseup', stopResize)
    // 恢复默认鼠标样式
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
  }

  // 返回所有需要的响应式数据和方法
  return {
    sidebarWidth,
    contentSplitRatio,
    rightPanelVisible,
    startLeftResize,
    startRightResize,
    stopResize,
  }
}
