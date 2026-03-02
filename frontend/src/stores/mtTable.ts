import { defineStore } from 'pinia'
import { storage } from '@/utils'
import { fields as defaultFields } from '@/modal/mock'
import { updateFieldApi } from '@/api'

//  类型定义
type TableColumn = {
  id: string
  title: string
  type?: string
  width?: number
  fixed?: boolean
  isShow?: boolean
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  [key: string]: any
}

// 表格记录项类型
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type TableRecord = Record<string, any>

// 表格配置完整类型（当前表）
type CurrentTableConfig = {
  fields: TableColumn[] // 表头字段
  records: TableRecord[] // 表格数据
  id: string // 表格唯一id（当前表的id）
  name: string // 表格名称
  version: number // 版本
  ct: Date // 创建时间
  ut: Date // 更新时间
  createdBy: string // 创建人
  updatedBy: string // 更新人
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  config: Record<string, any> // 其他配置(水印以及权限相关)
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  settings: Record<string, any> // 数据表配置（筛选、分组、排序等）
}

// 持久化存储键名（存储当前表）
const STORAGE_KEY = 'mtTable:current-config'

const useMtTableStore = defineStore('mtTable', () => {
  // 存储当前表格的配置（使用reactive直接管理单个对象）
  const currentTable = reactive<CurrentTableConfig>({
    id: '',
    name: '',
    version: 1,
    ct: new Date(),
    ut: new Date(),
    createdBy: '',
    updatedBy: '',
    fields: defaultFields || [],
    records: [],
    config: {},
    settings: {
      tableConfig: {
        rowHeight: 32,
      },
    },
  })

  // 从localStorage加载当前表配置
  const loadFromStorage = () => {
    const saved = storage.get(STORAGE_KEY)
    if (saved) {
      // 转换日期字符串为Date对象
      saved.ct = new Date(saved.ct)
      saved.ut = new Date(saved.ut)
      // 合并到当前表配置（保留响应式）
      Object.assign(currentTable, saved)
    }
  }

  // 初始化时加载本地存储的当前表
  loadFromStorage()

  // 保存当前表配置到localStorage
  const saveToStorage = () => {
    try {
      // 深克隆避免直接引用reactive对象
      const clone = JSON.parse(JSON.stringify(currentTable))
      storage.set(STORAGE_KEY, clone)
    } catch (err) {
      console.warn('保存当前表格配置失败', err)
    }
  }
  // 同步records与新增字段（
  const syncRecordsWithNewFields = (oldFields: TableColumn[], newFields: TableColumn[]) => {
    // 1. 获取旧字段的id集合
    const oldFieldIds = new Set(oldFields.map((field) => field.id))

    // 2. 筛选出新增的字段（旧字段中没有的id）
    const newAddedFields = newFields.filter((field) => !oldFieldIds.has(field.id))

    // 3. 若没有新增字段，直接返回原数据
    if (newAddedFields.length === 0) return currentTable.records

    // 4. 为每条记录添加新增字段的defaultValue
    const updatedRecords = currentTable.records.map((record) => {
      const newRecord = { ...record } // 不修改原记录，浅拷贝

      // 遍历新增字段，添加到记录中
      newAddedFields.forEach((field) => {
        // 优先使用字段的defaultValue，无则设为undefined（可根据需求调整默认值）
        newRecord[field.id] = field.defaultValue ?? ''
      })

      return newRecord
    })

    return updatedRecords
  }
  // 切换表格：重新设置当前表配置（切换表时调用）
  const setCurrentTable = (table: Partial<CurrentTableConfig>) => {
    // 合并新表格数据，缺失字段用默认值补充
    const now = new Date()
    Object.assign(currentTable, {
      // 默认值
      version: 1,
      ct: table.ct || now,
      ut: table.ut || now,
      createdBy: table.createdBy || 'system',
      updatedBy: table.updatedBy || 'system',
      fields: table.fields || defaultFields.map((f) => ({ ...f })),
      records: table.records || [],
      config: table.config || {},
      settings: table.settings || {},
      // 覆盖传入的字段
      ...table,
    })
    saveToStorage()
  }

  // 获取当前表的字段配置
  const getFields = (): TableColumn[] => {
    return currentTable.fields
  }

  // 更新当前表的字段配置
  const updateFields = async (newFields: TableColumn[]) => {
    //  保存更新前的旧字段（用于对比新增字段）
    const oldFields = [...currentTable.fields]

    // 更新字段配置
    currentTable.fields = newFields.map((f) => ({ ...f }))
    currentTable.version += 1
    currentTable.ut = new Date()

    // 同步records：为每条数据添加新增字段的defaultValue
    const updatedRecords = syncRecordsWithNewFields(oldFields, newFields)
    currentTable.records = updatedRecords // 更新数据

    // 保存到本地存储
    saveToStorage()
  }
  //   删除字段
  const delField = (fieldId: string) => {
    // 删除字段
    currentTable.fields = currentTable.fields.filter((field) => field.id !== fieldId)
    // 删除记录中的该字段数据
    currentTable.records = currentTable.records.map((record) => {
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const { [fieldId]: _, ...rest } = record
      return rest
    })
    saveToStorage()
  }

  //   删除行数据 参数为数组
  const delRecord = (recordIds: string[]) => {
    currentTable.records = currentTable.records.filter(
      (record) => !recordIds.includes(record.rowId),
    )
    saveToStorage()
  }

  // 更新当前表的记录数据
  const updateRecords = (newRecords: TableRecord[]) => {
    currentTable.records = [...newRecords]
    currentTable.version += 1
    currentTable.ut = new Date()
    saveToStorage()
  }
  //   在当前表插入新记录 数据源为对象
  const insertRecord = (newRecord: TableRecord, position: 'start' | 'end' = 'start') => {
    if (position === 'start') {
      currentTable.records = [newRecord, ...currentTable.records]
    } else {
      currentTable.records = [...currentTable.records, newRecord]
    }
    currentTable.version += 1
    currentTable.ut = new Date()
    saveToStorage()
  }
  // 重置当前表的配置（保留id和基础信息）
  const resetCurrentTable = () => {
    currentTable.fields = defaultFields.map((f) => ({ ...f }))
    currentTable.records = []
    currentTable.config = {}
    currentTable.settings = {}
    currentTable.version += 1
    currentTable.ut = new Date()
    saveToStorage()
  }

  // 切换当前表中某个字段的显示/隐藏
  const toggleFieldVisibility = async (fieldId: string) => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const field: any = currentTable.fields.find((f) => f.id === fieldId)
    if (field) {
      field.isShow = !field.isShow
      await updateFieldApi({ ...field })
      currentTable.version += 1
      currentTable.ut = new Date()
      saveToStorage()
    }
  }

  // 按固定/滚动分组获取当前表的字段
  const getGroupedFields = () => {
    const fields = getFields()
    return {
      fixed: fields.filter((f) => f.fixed),
      scrollable: fields.filter((f) => !f.fixed),
    }
  }

  // 获取当前表显示的固定列
  const getFixedFields = () => {
    return getFields().filter((f) => f.fixed && f.isShow)
  }

  // 获取当前表显示的滚动列
  const getScrollableFields = () => {
    return getFields().filter((f) => !f.fixed && f.isShow)
  }

  // 合并本地存储的字段与默认字段（用于初始化时）
  const mergeWithDefaults = (savedFields: TableColumn[], defaults: TableColumn[]) => {
    const savedMap = new Map(savedFields.map((f) => [f.id, f]))
    return defaults.map((df) => {
      const saved = savedMap.get(df.id)
      return saved ? { ...df, fixed: saved.fixed, isShow: saved.isShow } : { ...df }
    })
  }

  // 更新当前表的基本信息（名称、创建人等）
  const updateTableInfo = (
    info: Partial<Omit<CurrentTableConfig, 'id' | 'ct' | 'fields' | 'records' | 'config'>>,
  ) => {
    Object.assign(currentTable, {
      ...info,
      version: currentTable.version + 1,
      ut: new Date(),
    })
    saveToStorage()
  }

  // 更新当前表的其他配置
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const updateTableConfig = (newConfig: Record<string, any>) => {
    currentTable.config = { ...currentTable.config, ...newConfig }
    currentTable.version += 1
    currentTable.ut = new Date()
    saveToStorage()
  }

  // 生成测试数据并更新当前表
  const generateAndUpdateTestData = (count: number = 10) => {
    const departments = [
      '技术部',
      '产品部',
      '设计部',
      '运营部',
      '销售部',
      '人事部',
      '财务部',
      '法务部',
      '市场部',
      '客服部',
    ]
    const positions = [
      '工程师',
      '产品经理',
      '设计师',
      '运营专员',
      '销售经理',
      '人事专员',
      '财务专员',
      '法务顾问',
      '市场专员',
      '客服代表',
    ]
    const statuses = ['在职', '试用期', '实习生']
    const cities = ['北京', '上海', '深圳', '杭州', '广州', '成都', '南京', '苏州']
    const surnames = [
      '王',
      '李',
      '张',
      '刘',
      '陈',
      '杨',
      '黄',
      '赵',
      '吴',
      '周',
      '徐',
      '孙',
      '马',
      '朱',
      '胡',
      '郭',
      '何',
      '高',
      '林',
      '罗',
    ]
    const names = [
      '伟',
      '芳',
      '娜',
      '秀英',
      '敏',
      '非',
      '蕊',
      '静',
      '秀兰',
      '丹',
      '胜利',
      '欣',
      '婷',
      '晨',
      '佳',
      '秀梅',
      '晴',
      '凤英',
      '鑫',
      '丹丹',
    ]
    const data: TableRecord[] = []

    if (!storage.get(STORAGE_KEY)) {
      for (let i = 1; i <= count; i++) {
        const surname = surnames[Math.floor(Math.random() * surnames.length)]
        const name = names[Math.floor(Math.random() * names.length)]
        const department = departments[Math.floor(Math.random() * departments.length)]

        data.push({
          rowId: i + Math.floor(new Date().getTime() / 1000),
          idText: i,
          name: `${surname}${name}`,
          age: 22 + Math.floor(Math.random() * 38),
          department,
          position: positions[Math.floor(Math.random() * positions.length)],
          salary: (Math.floor(Math.random() * 50) + 8) * 1000,
          joinDate: new Date(
            2020 + Math.floor(Math.random() * 4),
            Math.floor(Math.random() * 12),
            Math.floor(Math.random() * 28) + 1,
          )
            .toISOString()
            .split('T')[0],
          email: `user${i}@company.com`,
          phone: `138${String(Math.floor(Math.random() * 100000000)).padStart(8, '0')}`,
          status: statuses[Math.floor(Math.random() * statuses.length)],
          city: cities[Math.floor(Math.random() * cities.length)],
        })
      }

      data.sort((a, b) => a.department.localeCompare(b.department))
    }
    updateRecords(storage.get(STORAGE_KEY)?.records || data)
    return storage.get(STORAGE_KEY)?.records || data
  }
  // 获取当前表的记录数据
  const getRecords = (): TableRecord[] => {
    return currentTable.records
  }

  // 更新当前表的 settings 配置（支持只更新部分字段）
  const updateTableSettings = (newSettings: Partial<CurrentTableConfig['settings']>) => {
    currentTable.settings = {
      ...currentTable.settings,
      ...newSettings,
      tableConfig: {
        ...currentTable.settings.tableConfig,
        ...(newSettings.tableConfig || {}),
        filterConfig:
          newSettings.tableConfig?.filterConfig ?? currentTable.settings.tableConfig.filterConfig,
        groupConfig:
          newSettings.tableConfig?.groupConfig ?? currentTable.settings.tableConfig.groupConfig,
        sortConfig:
          newSettings.tableConfig?.sortConfig ?? currentTable.settings.tableConfig.sortConfig,
        colorConfig:
          newSettings.tableConfig?.colorConfig ?? currentTable.settings.tableConfig.colorConfig,
      },
    }
    currentTable.version += 1
    currentTable.ut = new Date()
    saveToStorage()
  }

  // 更新当前表中指定rowId的记录
  const updateRecord = (rowId: string, updatedData: Partial<TableRecord>) => {
    const recordIndex = currentTable.records.findIndex((record) => record.rowId === rowId)
    if (recordIndex !== -1) {
      // 浅拷贝记录以保持响应性
      currentTable.records[recordIndex] = {
        ...currentTable.records[recordIndex],
        ...updatedData,
      }
      currentTable.version += 1
      currentTable.ut = new Date()
      saveToStorage()
    }
  }
  return {
    // 状态：当前表的完整配置
    currentTable,
    // 方法
    loadFromStorage, // 从本地存储加载当前表
    saveToStorage, // 保存当前表到本地存储
    setCurrentTable, // 切换表时调用，设置新的当前表
    getFields, // 获取当前表的字段配置
    updateFields, // 更新当前表的字段配置
    updateRecords, // 更新当前表的记录数据
    resetCurrentTable, // 重置当前表的配置（保留id和基础信息）
    toggleFieldVisibility, // 切换当前表中某个字段的显示/隐藏
    getGroupedFields, // 按固定/滚动分组获取当前表的字段
    getFixedFields, // 获取当前表显示的固定列
    getScrollableFields, // 获取当前表显示的滚动列字段
    updateTableInfo, // 更新当前表的基本信息（名称、创建人等）
    updateTableConfig, // 更新当前表的其他配置
    generateAndUpdateTestData, // 生成测试数据并更新当前表
    mergeWithDefaults, // 合并本地存储的字段与默认字段（用于初始化时）
    getRecords, // 获取当前表的记录数据
    updateTableSettings, // 更新当前表的 settings 配置
    insertRecord, // 在当前表插入新记录 数据源为对象
    delField, // 删除字段
    delRecord, // 删除记录
    updateRecord, // 更新行记录
  }
})

export { useMtTableStore, type TableColumn, type CurrentTableConfig, type TableRecord }
