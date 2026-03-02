<template>
  <div class="w-100 h-100 flx-ce-sta" @mouseenter="handleMouseEnter" @mouseleave="handleMouseLeave">
    <div class="w-100" v-if="listData.length > 0">
      <div class="attachment-list">
        <div class="attachment-item" v-for="(item, index) in listData" :key="index">
          <t-image-viewer v-if="isImage(item)" :images="[previewUrl]">
            <template #trigger="{ open }">
              <img @click="handlePre(item, open)" :src="item.url" :alt="item.name" />
            </template>
          </t-image-viewer>

          <div class="file-icon-wrapper" v-else>
            <t-tooltip :content="item.name" placement="top">
              <t-icon :name="getFileIcon(item)" size="14px" />
            </t-tooltip>
          </div>
        </div>
      </div>
    </div>
    <t-popup
      ref="popupRef"
      trigger="click"
      placement="bottom-left"
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
          <h3>多文件上传（最多5个）</h3>
          <FileUpload
            accept="*"
            multiple
            :max="5"
            :max-size="10"
            upload-text="点击或拖拽多张图片到此区域"
            upload-hint="支持批量上传，最多 5 张，单个文件不超过 10MB"
            @update:fileList="handleAttachmentsUpdate"
            :default-file-list="listData"
          />
        </div>
      </template>
    </t-popup>
  </div>
</template>

<script setup lang="ts">
import { ObjType } from '@/types'
import { isImage, getFileIcon, FileItem } from '@/utils'

interface Props {
  field: ObjType
  item: ObjType
  isEditor: boolean
}

const props = defineProps<Props>()
const emit = defineEmits(['handleSpecialChange'])
// 事件日志
const listData = ref<FileItem[]>([])
// 控制鼠标悬停状态
const isHovering = ref(false)
const previewUrl = ref('')
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
  console.log('上传列表更新:', files)
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
.attachment-list {
  display: flex;
  align-items: center;
  min-height: 22px;
  gap: 6px;
  overflow-x: auto;
  padding: 4px 6px;
  border-radius: 4px;
  scrollbar-width: none;
  &::-webkit-scrollbar {
    display: none;
  }
  &::-webkit-scrollbar {
    display: none;
  }
  -ms-overflow-style: none;
  .attachment-item {
    display: flex;
    align-items: center;
    justify-content: center;
    min-width: 22px;
    width: auto;
    height: 22px;
    border-radius: 4px;
    flex-shrink: 0;
    img {
      max-width: 100%;
      max-height: 100%;
      border-radius: 4px;
      object-fit: contain;
    }
  }
}
.file-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  color: #666;
}
.upload-text {
  text-align: center;
  color: #9a9dac;
  font-size: 12px;
  overflow: hidden;
}
</style>
