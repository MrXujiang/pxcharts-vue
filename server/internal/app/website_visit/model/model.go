package model

import (
	"time"
)

// WebsiteVisit 网站访问记录模型
type WebsiteVisit struct {
	ID        string     `gorm:"column:id;type:varchar(36);primaryKey;not null" json:"id"`
	IP        string     `gorm:"column:ip;type:text;default:''" json:"ip"`
	UserAgent string     `gorm:"column:user_agent;type:text;default:''" json:"userAgent"`
	Referer   string     `gorm:"column:referer;type:text;default:''" json:"referer"`
	Page      string     `gorm:"column:page;type:text;default:''" json:"page"`
	Method    string     `gorm:"column:method;type:text;default:''" json:"method"`
	SessionID string     `gorm:"column:session_id;type:text;default:''" json:"sessionId"`
	UserID    string     `gorm:"column:user_id;type:text;default:''" json:"userId"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamptz" json:"-"`
}

func (*WebsiteVisit) TableName() string {
	return "website_visit"
}

// GetStatisticsRes 统计数据响应
type GetStatisticsRes struct {
	// 访问量统计
	VisitStats VisitStats `json:"visitStats"`
	// 用户量统计
	UserStats UserStats `json:"userStats"`
	// 团队数统计
	TeamStats TeamStats `json:"teamStats"`
	// 项目数统计
	ProjectStats ProjectStats `json:"projectStats"`
}

// VisitStats 访问量统计
type VisitStats struct {
	TotalVisits  int64            `json:"totalVisits"`  // 总访问量
	TodayVisits  int64            `json:"todayVisits"`  // 今日访问量
	SixMonthData []VisitDataPoint `json:"sixMonthData"` // 近半年访问数据
	XAxis        []string         `json:"xAxis"`        // 时间轴标签
	Data         []int64          `json:"data"`         // 对应时间轴的数据
}

// UserStats 用户量统计
type UserStats struct {
	TotalUsers       int64   `json:"totalUsers"`       // 总用户量
	WeeklyGrowthRate float64 `json:"weeklyGrowthRate"` // 周同比增长百分比
	DailyGrowthRate  float64 `json:"dailyGrowthRate"`  // 日环比增长百分比
	DailyGrowthCount int64   `json:"dailyGrowthCount"` // 日增长量
}

// TeamStats 团队数统计
type TeamStats struct {
	TotalTeams int64 `json:"totalTeams"` // 总团队数
}

// ProjectStats 项目数统计
type ProjectStats struct {
	TotalProjects int64 `json:"totalProjects"` // 总项目数
}

// VisitDataPoint 访问数据点（用于echarts折线图）
type VisitDataPoint struct {
	Date  string `json:"date"`  // 日期，格式：YYYY-MM-DD
	Count int64  `json:"count"` // 访问量
}

// GetNewUsersReq 获取新增用户请求
type GetNewUsersReq struct {
	TimeRange string `json:"timeRange" form:"timeRange" binding:"required,oneof=day week month year"` // 时间范围：day, week, month, year
}

// GetNewUsersRes 获取新增用户响应
type GetNewUsersRes struct {
	XAxis []string `json:"xAxis"` // X轴标签
	Data  []int64  `json:"data"`  // 数据值
}

// TimeRangeXAxis 不同时间范围的X轴标签
var TimeRangeXAxis = map[string][]string{
	"day": {
		"02h", "04h", "06h", "08h", "10h", "12h",
		"14h", "16h", "18h", "20h", "22h", "24h",
	},
	"week": {
		"周一", "周二", "周三", "周四", "周五", "周六", "周日",
	},
	"month": {
		"1-3", "4-6", "7-9", "10-12", "13-15",
		"16-18", "19-21", "22-24", "25-27", "28-30",
	},
	"year": {
		"1月", "2月", "3月", "4月", "5月", "6月",
		"7月", "8月", "9月", "10月", "11月", "12月",
	},
}
