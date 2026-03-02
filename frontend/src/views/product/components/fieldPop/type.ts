import { ObjType } from '@/types'

interface FieldsResult {
  fixedFields: ObjType[] // 固定列数组
  nonFixedFields: ObjType[] // 非固定列数组
}

export { FieldsResult }
