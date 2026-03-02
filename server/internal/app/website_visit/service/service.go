package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	projectModel "mvtable/internal/app/mv_project/model"
	teamModel "mvtable/internal/app/team/model"
	userModel "mvtable/internal/app/user/model"
	visitModel "mvtable/internal/app/website_visit/model"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"
)

// WebsiteVisitService 网站访问记录服务
type WebsiteVisitService struct{}

// NewWebsiteVisitService 创建网站访问记录服务
func NewWebsiteVisitService() *WebsiteVisitService {
	return &WebsiteVisitService{}
}

// RecordVisit 记录访问
func (s *WebsiteVisitService) RecordVisit(ip, userAgent, referer, page, method, sessionID, userID string) {
	// 异步记录，避免影响响应性能
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("记录访问信息时发生panic", zap.Any("error", r))
			}
		}()

		visit := &visitModel.WebsiteVisit{
			ID:        uuid.New().String(),
			IP:        ip,
			UserAgent: userAgent,
			Referer:   referer,
			Page:      page,
			Method:    method,
			SessionID: sessionID,
			UserID:    userID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(db.GetDB(), visit); err != nil {
			log.Error("记录访问信息失败", zap.Error(err), zap.String("visit_id", visit.ID))
		}
	}()
}

// GetStatistics 获取网站统计数据
func (s *WebsiteVisitService) GetStatistics() (*visitModel.GetStatisticsRes, error) {
	response := &visitModel.GetStatisticsRes{}

	// 获取访问量统计
	visitStats, err := s.getVisitStatistics()
	if err != nil {
		return nil, err
	}
	response.VisitStats = *visitStats

	// 获取用户量统计
	userStats, err := s.getUserStatistics()
	if err != nil {
		return nil, err
	}
	response.UserStats = *userStats

	// 获取团队数统计
	teamStats, err := s.getTeamStatistics()
	if err != nil {
		return nil, err
	}
	response.TeamStats = *teamStats

	// 获取项目数统计
	projectStats, err := s.getProjectStatistics()
	if err != nil {
		return nil, err
	}
	response.ProjectStats = *projectStats

	return response, nil
}

// getVisitStatistics 获取访问量统计
func (s *WebsiteVisitService) getVisitStatistics() (*visitModel.VisitStats, error) {
	database := db.GetDB()

	// 总访问量（根据SessionID去重，同一个SessionID只算一次有效访问）
	var totalVisits int64
	if err := database.Model(&visitModel.WebsiteVisit{}).
		Distinct("session_id").
		Count(&totalVisits).Error; err != nil {
		log.Error("获取总访问量失败", zap.Error(err))
		return nil, err
	}

	// 今日访问量（根据SessionID去重，同一个SessionID当日只算一次有效访问）
	today := time.Now().Format("2006-01-02")
	var todayVisits int64
	if err := database.Model(&visitModel.WebsiteVisit{}).
		Where("DATE(created_at) = ?", today).
		Distinct("session_id").
		Count(&todayVisits).Error; err != nil {
		log.Error("获取今日访问量失败", zap.Error(err))
		return nil, err
	}

	// 近半年访问数据（每日根据SessionID去重）
	sixMonthData, xAxis, data, err := s.getSixMonthVisitDataWithChart()
	if err != nil {
		return nil, err
	}

	return &visitModel.VisitStats{
		TotalVisits:  totalVisits,
		TodayVisits:  todayVisits,
		SixMonthData: sixMonthData,
		XAxis:        xAxis,
		Data:         data,
	}, nil
}

// getSixMonthVisitData 获取近半年访问数据
func (s *WebsiteVisitService) getSixMonthVisitData() ([]visitModel.VisitDataPoint, error) {
	database := db.GetDB()

	// 计算6个月前的日期
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)

	var results []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	query := `
		SELECT
			DATE(created_at) as date,
			COUNT(*) as count
		FROM website_visit
		WHERE created_at >= ?
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at)
	`

	if err := database.Raw(query, sixMonthsAgo).Scan(&results).Error; err != nil {
		log.Error("获取近半年访问数据失败", zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	data := make([]visitModel.VisitDataPoint, len(results))
	for i, result := range results {
		data[i] = visitModel.VisitDataPoint{
			Date:  result.Date,
			Count: result.Count,
		}
	}

	return data, nil
}

// getSixMonthVisitDataWithChart 获取近半年访问数据（包含图表数据，根据SessionID去重）
func (s *WebsiteVisitService) getSixMonthVisitDataWithChart() ([]visitModel.VisitDataPoint, []string, []int64, error) {
	database := db.GetDB()

	// 计算6个月前的日期
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)

	var results []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	// 根据SessionID去重统计每日访问量
	query := `
		SELECT
			DATE(created_at) as date,
			COUNT(DISTINCT session_id) as count
		FROM website_visit
		WHERE created_at >= ?
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at)
	`

	if err := database.Raw(query, sixMonthsAgo).Scan(&results).Error; err != nil {
		log.Error("获取近半年访问数据失败", zap.Error(err))
		return nil, nil, nil, err
	}

	// 转换为响应格式
	sixMonthData := make([]visitModel.VisitDataPoint, len(results))
	xAxis := make([]string, len(results))
	data := make([]int64, len(results))

	for i, result := range results {
		sixMonthData[i] = visitModel.VisitDataPoint{
			Date:  result.Date,
			Count: result.Count,
		}
		xAxis[i] = result.Date
		data[i] = result.Count
	}

	return sixMonthData, xAxis, data, nil
}

// getUserStatistics 获取用户量统计
func (s *WebsiteVisitService) getUserStatistics() (*visitModel.UserStats, error) {
	database := db.GetDB()

	// 总用户量
	var totalUsers int64
	if err := database.Model(&userModel.User{}).Count(&totalUsers).Error; err != nil {
		log.Error("获取总用户量失败", zap.Error(err))
		return nil, err
	}

	// 计算周同比增长
	weeklyGrowthRate, err := s.calculateWeeklyUserGrowth()
	if err != nil {
		return nil, err
	}

	// 计算日环比增长
	dailyGrowthRate, dailyGrowthCount, err := s.calculateDailyUserGrowth()
	if err != nil {
		return nil, err
	}

	return &visitModel.UserStats{
		TotalUsers:       totalUsers,
		WeeklyGrowthRate: weeklyGrowthRate,
		DailyGrowthRate:  dailyGrowthRate,
		DailyGrowthCount: dailyGrowthCount,
	}, nil
}

// calculateWeeklyUserGrowth 计算周同比增长百分比
func (s *WebsiteVisitService) calculateWeeklyUserGrowth() (float64, error) {
	database := db.GetDB()

	now := time.Now()
	thisWeekStart := now.AddDate(0, 0, -int(now.Weekday()))
	lastWeekStart := thisWeekStart.AddDate(0, 0, -7)
	lastWeekEnd := thisWeekStart.AddDate(0, 0, -1)

	var thisWeekCount, lastWeekCount int64

	// 本周用户数
	if err := database.Model(&userModel.User{}).
		Where("created_at >= ? AND created_at <= ?", thisWeekStart, now).
		Count(&thisWeekCount).Error; err != nil {
		return 0, err
	}

	// 上周用户数
	if err := database.Model(&userModel.User{}).
		Where("created_at >= ? AND created_at <= ?", lastWeekStart, lastWeekEnd).
		Count(&lastWeekCount).Error; err != nil {
		return 0, err
	}

	if lastWeekCount == 0 {
		return 0, nil
	}

	growthRate := float64(thisWeekCount-lastWeekCount) / float64(lastWeekCount) * 100
	return growthRate, nil
}

// calculateDailyUserGrowth 计算日环比增长百分比和增长量
func (s *WebsiteVisitService) calculateDailyUserGrowth() (float64, int64, error) {
	database := db.GetDB()

	yesterday := time.Now().AddDate(0, 0, -1)
	todayStart := time.Now().Truncate(24 * time.Hour)
	yesterdayStart := yesterday.Truncate(24 * time.Hour)
	yesterdayEnd := todayStart.Add(-time.Second)

	var todayCount, yesterdayCount int64

	// 今日用户数
	if err := database.Model(&userModel.User{}).
		Where("created_at >= ?", todayStart).
		Count(&todayCount).Error; err != nil {
		return 0, 0, err
	}

	// 昨日用户数
	if err := database.Model(&userModel.User{}).
		Where("created_at >= ? AND created_at <= ?", yesterdayStart, yesterdayEnd).
		Count(&yesterdayCount).Error; err != nil {
		return 0, 0, err
	}

	var growthRate float64
	var growthCount int64

	if yesterdayCount == 0 {
		growthRate = 0
		growthCount = todayCount
	} else {
		growthRate = float64(todayCount-yesterdayCount) / float64(yesterdayCount) * 100
		growthCount = todayCount - yesterdayCount
	}

	return growthRate, growthCount, nil
}

// getTeamStatistics 获取团队数统计
func (s *WebsiteVisitService) getTeamStatistics() (*visitModel.TeamStats, error) {
	database := db.GetDB()

	var totalTeams int64
	if err := database.Model(&teamModel.Team{}).Count(&totalTeams).Error; err != nil {
		log.Error("获取总团队数失败", zap.Error(err))
		return nil, err
	}

	return &visitModel.TeamStats{
		TotalTeams: totalTeams,
	}, nil
}

// getProjectStatistics 获取项目数统计
func (s *WebsiteVisitService) getProjectStatistics() (*visitModel.ProjectStats, error) {
	database := db.GetDB()

	var totalProjects int64
	if err := database.Model(&projectModel.MvProject{}).Count(&totalProjects).Error; err != nil {
		log.Error("获取总项目数失败", zap.Error(err))
		return nil, err
	}

	return &visitModel.ProjectStats{
		TotalProjects: totalProjects,
	}, nil
}

// GetNewUsers 获取新增用户数据
func (s *WebsiteVisitService) GetNewUsers(timeRange string) (*visitModel.GetNewUsersRes, error) {
	switch timeRange {
	case "day":
		return s.getNewUsersByDay()
	case "week":
		return s.getNewUsersByWeek()
	case "month":
		return s.getNewUsersByMonth()
	case "year":
		return s.getNewUsersByYear()
	default:
		return nil, fmt.Errorf("不支持的时间范围: %s", timeRange)
	}
}

// getNewUsersByDay 获取今日各时段新增用户数
func (s *WebsiteVisitService) getNewUsersByDay() (*visitModel.GetNewUsersRes, error) {
	database := db.GetDB()
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	xAxis := visitModel.TimeRangeXAxis["day"]
	data := make([]int64, len(xAxis))

	// 查询每个2小时时段的新增用户数
	for i, hourStr := range xAxis {
		var hour int
		if hourStr == "24h" {
			hour = 24
		} else {
			// 移除"h"后缀并转换为int
			hourStr = strings.TrimSuffix(hourStr, "h")
			hour, _ = strconv.Atoi(hourStr)
		}

		var startHour, endHour int
		if hour == 24 {
			startHour = 22
			endHour = 23
		} else {
			startHour = hour - 2
			endHour = hour - 1
		}

		startTime := todayStart.Add(time.Duration(startHour) * time.Hour)
		endTime := todayStart.Add(time.Duration(endHour+1)*time.Hour - time.Second)

		var count int64
		if err := database.Model(&userModel.User{}).
			Where("created_at >= ? AND created_at <= ?", startTime, endTime).
			Count(&count).Error; err != nil {
			log.Error("获取每日新增用户数失败", zap.Error(err), zap.String("hourRange", hourStr))
			return nil, err
		}

		data[i] = count
	}

	return &visitModel.GetNewUsersRes{
		XAxis: xAxis,
		Data:  data,
	}, nil
}

// getNewUsersByWeek 获取本周每日新增用户数
func (s *WebsiteVisitService) getNewUsersByWeek() (*visitModel.GetNewUsersRes, error) {
	database := db.GetDB()
	now := time.Now()

	// 计算本周开始时间（周一）
	weekday := int(now.Weekday())
	if weekday == 0 { // 周日
		weekday = 7
	}
	weekStart := now.AddDate(0, 0, -(weekday - 1))
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())

	xAxis := visitModel.TimeRangeXAxis["week"]
	data := make([]int64, len(xAxis))

	// 查询每天的新增用户数
	for i := 0; i < 7; i++ {
		startTime := weekStart.AddDate(0, 0, i)
		endTime := startTime.AddDate(0, 0, 1).Add(-time.Second)

		var count int64
		if err := database.Model(&userModel.User{}).
			Where("created_at >= ? AND created_at <= ?", startTime, endTime).
			Count(&count).Error; err != nil {
			log.Error("获取每周新增用户数失败", zap.Error(err), zap.Int("dayIndex", i))
			return nil, err
		}

		data[i] = count
	}

	return &visitModel.GetNewUsersRes{
		XAxis: xAxis,
		Data:  data,
	}, nil
}

// getNewUsersByMonth 获取本月各时间段新增用户数
func (s *WebsiteVisitService) getNewUsersByMonth() (*visitModel.GetNewUsersRes, error) {
	database := db.GetDB()
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	xAxis := visitModel.TimeRangeXAxis["month"]
	data := make([]int64, len(xAxis))

	// 查询每个3天时间段的新增用户数
	for i, periodStr := range xAxis {
		// 解析时间段，如 "1-3" 表示1-3号
		var startDay, endDay int
		parts := strings.Split(periodStr, "-")
		if len(parts) == 2 {
			startDay, _ = strconv.Atoi(parts[0])
			endDay, _ = strconv.Atoi(parts[1])
		}

		startTime := monthStart.AddDate(0, 0, startDay-1)
		endTime := monthStart.AddDate(0, 0, endDay).Add(-time.Second)

		// 确保不超过当月最后一天
		if endTime.After(now) {
			endTime = now
		}

		var count int64
		if err := database.Model(&userModel.User{}).
			Where("created_at >= ? AND created_at <= ?", startTime, endTime).
			Count(&count).Error; err != nil {
			log.Error("获取每月新增用户数失败", zap.Error(err), zap.String("period", periodStr))
			return nil, err
		}

		data[i] = count
	}

	return &visitModel.GetNewUsersRes{
		XAxis: xAxis,
		Data:  data,
	}, nil
}

// getNewUsersByYear 获取今年各月新增用户数
func (s *WebsiteVisitService) getNewUsersByYear() (*visitModel.GetNewUsersRes, error) {
	database := db.GetDB()
	now := time.Now()
	yearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	xAxis := visitModel.TimeRangeXAxis["year"]
	data := make([]int64, len(xAxis))

	// 查询每月的新增用户数
	for i := 0; i < 12; i++ {
		startTime := yearStart.AddDate(0, i, 0)
		endTime := startTime.AddDate(0, 1, 0).Add(-time.Second)

		// 如果超过当前时间，则截止到当前时间
		if endTime.After(now) {
			endTime = now
		}

		var count int64
		if err := database.Model(&userModel.User{}).
			Where("created_at >= ? AND created_at <= ?", startTime, endTime).
			Count(&count).Error; err != nil {
			log.Error("获取每年新增用户数失败", zap.Error(err), zap.Int("monthIndex", i))
			return nil, err
		}

		data[i] = count
	}

	return &visitModel.GetNewUsersRes{
		XAxis: xAxis,
		Data:  data,
	}, nil
}
