import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getTeamList } from '@/api'

export interface Team {
  id: string
  name: string
  description?: string
  logo?: string
  color?: string
  creator?: string
  memberCount?: number
  displayAvatarList?: string[]
  createdAt?: string
  updatedAt?: string
  isExternal?: boolean
}

export const useTeamStore = defineStore('team', () => {
  // 团队列表
  const teamList = ref<Team[]>([])
  
  // 当前选中的团队
  const currentTeam = ref<Team | null>(null)
  
  // 加载状态
  const loading = ref(false)

  // 获取随机颜色
  const getRandomColor = () => {
    const colors = ['#5B8FF9', '#61DDAA', '#F6BD16', '#6DC8EC', '#945FB9', '#FF9D4D', '#5AD8A6', '#5D7092']
    return colors[Math.floor(Math.random() * colors.length)]
  }

  // 加载团队列表
  const loadTeamList = async () => {
    try {
      loading.value = true
      const response = await getTeamList() as any
      
      // 处理响应数据
      let teams: Team[] = []
      if (response && Array.isArray(response)) {
        teams = response
      } else if (response && response.list && Array.isArray(response.list)) {
        teams = response.list
      }
      
      // 为每个团队分配颜色（如果没有）
      teamList.value = teams.map(team => ({
        ...team,
        color: team.color || getRandomColor()
      }))
      
      // 如果有团队且当前没有选中团队，选中第一个
      if (teamList.value.length > 0 && !currentTeam.value) {
        currentTeam.value = teamList.value[0]
      }
      
      return teamList.value
    } catch (error: any) {
      console.error('加载团队列表失败：', error)
      teamList.value = []
      throw error
    } finally {
      loading.value = false
    }
  }

  // 设置当前团队
  const setCurrentTeam = (team: Team) => {
    currentTeam.value = team
  }

  // 添加团队
  const addTeam = (team: Team) => {
    teamList.value.push({
      ...team,
      color: team.color || getRandomColor()
    })
  }

  // 更新团队
  const updateTeam = (teamId: string, updatedTeam: Partial<Team>) => {
    const index = teamList.value.findIndex(t => t.id === teamId)
    if (index !== -1) {
      teamList.value[index] = { ...teamList.value[index], ...updatedTeam }
      // 如果更新的是当前团队，也更新当前团队
      if (currentTeam.value?.id === teamId) {
        currentTeam.value = teamList.value[index]
      }
    }
  }

  // 删除团队
  const removeTeam = (teamId: string) => {
    const index = teamList.value.findIndex(t => t.id === teamId)
    if (index !== -1) {
      teamList.value.splice(index, 1)
      // 如果删除的是当前团队，切换到第一个团队
      if (currentTeam.value?.id === teamId) {
        currentTeam.value = teamList.value.length > 0 ? teamList.value[0] : null
      }
    }
  }

  return {
    teamList,
    currentTeam,
    loading,
    loadTeamList,
    setCurrentTeam,
    addTeam,
    updateTeam,
    removeTeam,
  }
})
