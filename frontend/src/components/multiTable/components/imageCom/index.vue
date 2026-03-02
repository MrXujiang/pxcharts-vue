<template>
  <div class="w-100 h-100 flx-ce-sta" @mouseenter="handleMouseEnter" @mouseleave="handleMouseLeave">
    <div class="w-100" v-if="listData.length > 0">
      <div class="image-list">
        <t-image-viewer :images="[previewUrl]">
          <template #trigger="{ open }">
            <img
              @click="handlePre(item, open)"
              v-for="(item, index) in listData"
              :key="index"
              :src="item.url"
              alt=""
            />
          </template>
        </t-image-viewer>
      </div>
    </div>
    <t-popup
      ref="popupRef"
      trigger="click"
      placement="bottom-left"
      attach="body"
      @visible-change="handleVisibleChange"
    >
      <template v-if="isEditor">
        <div v-if="!listData.length" class="w-100 upload-text pointer">点击上传</div>
        <t-button
          v-else-if="isHovering"
          size="small"
          theme="default"
          variant="base"
          class="upload-btn"
        >
          <template #icon>
            <t-icon size="16px" name="fullscreen1" />
          </template>
        </t-button>
      </template>

      <template #content>
        <div class="upload-content">
          <h3>多图上传（最多9张）</h3>
          <FileUpload
            mode="image-wall"
            accept="image/*"
            multiple
            :max="9"
            :max-size="5"
            @preview="handlePreview"
            :defaultFileList="listData"
            @update:fileList="handleAttachmentsUpdate"
          />
        </div>
      </template>
    </t-popup>
  </div>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'
import { FileItem } from '@/utils'

interface Props {
  field: ObjType
  item: ObjType
  isEditor: boolean
}

const props = defineProps<Props>()
const emit = defineEmits(['handleSpecialChange'])
const previewUrl = ref('')

// 事件日志
const eventLogs = ref<Array<{ time: string; event: string; message: string; type: string }>>([])
const listData = ref([])

// 控制鼠标悬停状态
const isHovering = ref(false)

watch(
  () => props.item,
  (val) => {
    if (val) {
      listData.value = val[props.field.id] || []
    }
  },
  { immediate: true },
)

const handlePre = async (file, open) => {
  previewUrl.value = file.url
  open()
}

const handleMouseEnter = () => {
  isHovering.value = true
}

const handleMouseLeave = () => {
  isHovering.value = false
}

const handleVisibleChange = (visible: boolean) => {
  if (!visible) {
    emit('handleSpecialChange', {
      field: props.field,
      item: props.item,
      value: listData.value,
    })
  }
}
// 图片列表数据获取
const handleAttachmentsUpdate = (files: FileItem[]) => {
  listData.value = files
}
const addLog = (event: string, message: string, type: string = 'info') => {
  const now = new Date()
  const time = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}:${now.getSeconds().toString().padStart(2, '0')}`

  eventLogs.value.unshift({
    time,
    event,
    message,
    type,
  })

  // 只保留最近20条
  if (eventLogs.value.length > 20) {
    eventLogs.value = eventLogs.value.slice(0, 20)
  }
}
// 预览图片
const handlePreview = (file: FileItem) => {
  console.log('预览图片:', file)
  addLog('preview', `预览图片: ${file.name}`, 'info')
}
defineOptions({
  name: 'ImageCom',
})
</script>

<style lang="less" scoped>
.upload-btn {
  position: absolute;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
}
.upload-content {
  width: 420px;
  padding: 8px;
}
.image-list {
  display: flex;
  align-items: center;
  height: 22px;
  gap: 6px;
  overflow-x: auto;
  padding: 0 4px;
  border-radius: 4px;
  scrollbar-width: none;
  &::-webkit-scrollbar {
    display: none;
  }
  &::-webkit-scrollbar {
    display: none;
  }
  -ms-overflow-style: none;
  img {
    height: 100%;
    border-radius: 4px;
    width: auto;
    min-width: 0;
    object-fit: contain;
    flex-shrink: 0;
  }
}
.upload-text {
  text-align: center;
  color: #9a9dac;
  font-size: 12px;
  overflow: hidden;
}
</style>
