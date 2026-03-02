import dayjs from 'dayjs'

const now = dayjs()

const formatOptions = [
  { label: now.format('YYYY-MM-DD'), value: 'YYYY-MM-DD' },
  { label: now.format('YYYY-MM-DD HH:mm'), value: 'YYYY-MM-DD HH:mm' },
  { label: now.format('YYYY-MM-DD HH:mm:ss'), value: 'YYYY-MM-DD HH:mm:ss' },
  { label: now.format('YYYY/MM/DD'), value: 'YYYY/MM/DD' },
  { label: now.format('YYYY/MM/DD HH:mm'), value: 'YYYY/MM/DD HH:mm' },
  { label: now.format('YYYY/MM/DD HH:mm:ss'), value: 'YYYY/MM/DD HH:mm:ss' },
  { label: now.format('HH:mm'), value: 'HH:mm' },
  { label: now.format('YYYY年MM月DD日'), value: 'YYYY年MM月DD日' },
  { label: now.format('YYYY年MM月'), value: 'YYYY年MM月' },
  { label: now.format('MM月DD日'), value: 'MM月DD日' },
]

export { formatOptions }
