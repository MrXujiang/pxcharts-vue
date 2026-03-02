// 工具函数统一出口
// export * from './common'
// export * from './validate'
import { fieldTypes } from '@/modal/options'

const storage = {
  get(key: string) {
    const data = localStorage.getItem(key)
    if (data) {
      try {
        return JSON.parse(data)
        // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/no-unused-vars
      } catch (err: any) {
        return data
      }
    }

    return null
  },
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  set(key: string, value: any) {
    if (typeof value !== 'object') {
      localStorage.setItem(key, value)
      return
    }

    localStorage.setItem(key, JSON.stringify(value))
  },
  remove(key: string) {
    localStorage.removeItem(key)
  },
  clear() {
    localStorage.clear()
  },
}

// 获取图标名的工具函数
const getIconNameByType = (type) => {
  const iconMap: Record<string, string> = {}
  Object.values(fieldTypes).forEach((category) => {
    category.forEach((field) => {
      iconMap[field.type] = field.iconName
    })
  })
  return iconMap[type]
}

const htmlToText = (html: string) => {
  // 将所有 <br> 标签（不区分大小写、自闭合形式）替换为一个空格
  const htmlWithSpaces = html.replace(/<br\s*\/?>/gi, ' ')

  const parser = new DOMParser()
  const doc = parser.parseFromString(htmlWithSpaces, 'text/html')

  // 获取文本并压缩连续空白（可选）
  let text = doc.body.textContent || doc.body.innerText || ''

  // 可选：将多个连续空格压缩为一个（避免因多个<br>产生过多空格）
  text = text.replace(/\s+/g, ' ').trim()

  return text
}

interface FileItem {
  id: string
  name: string
  size: number
  type: string
  url: string
  uploading?: boolean
  error?: boolean
}
// 判断是否是图片
const isImage = (file: FileItem): boolean => {
  return file.type.startsWith('image/') || /\.(jpg|jpeg|png|gif|bmp|webp|svg)$/i.test(file.name)
}

// 获取文件图标
const getFileIcon = (file: FileItem): string => {
  const ext = file.name.split('.').pop()?.toLowerCase()

  const iconMap: Record<string, string> = {
    pdf: 'filePdf',
    doc: 'fileWord',
    docx: 'fileWord',
    xls: 'fileExcel',
    xlsx: 'fileExcel',
    ppt: 'filePowerpoint',
    pptx: 'filePowerpoint',
    zip: 'fileZip',
    rar: 'fileZip',
    txt: 'fileText',
    mp4: 'video',
    avi: 'video',
    mp3: 'music',
    wav: 'music',
  }

  return iconMap[ext || ''] || 'file'
}
/**
 * 设置URL参数（覆盖或新增）
 * @param newParams 要添加的参数
 */
const setUrlParams = (newParams: Record<string, string>) => {
  // 获取当前URL的基础路径（不带参数）
  const baseUrl = window.location.origin + window.location.pathname
  // 创建URLSearchParams实例
  const searchParams = new URLSearchParams(window.location.search)
  // 遍历新参数，添加/覆盖到searchParams
  Object.entries(newParams).forEach(([key, value]) => {
    searchParams.set(key, value)
  })
  // 拼接新的URL并替换地址栏（不刷新页面）
  const newUrl = `${baseUrl}?${searchParams.toString()}`
  window.history.replaceState({}, '', newUrl)
}

// 清空 reactive 对象的所有属性（保留响应性）
const clearReactiveObject = <T extends object>(obj: T) => {
  // 遍历对象所有自有属性并删除
  Object.keys(obj).forEach((key) => {
    delete obj[key as keyof T]
  })
}
export {
  getIconNameByType,
  storage,
  htmlToText,
  isImage,
  getFileIcon,
  FileItem,
  setUrlParams,
  clearReactiveObject,
}
