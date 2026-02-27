package controllers

import (
	"account-management-system/backend/database"
	"account-management-system/backend/models"
	"account-management-system/backend/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// StatisticsRequest 统计请求
type StatisticsRequest struct {
	Dimension string `form:"dimension" binding:"required,oneof=project person category"`
	Cycle     string `form:"cycle" binding:"required,oneof=month quarter year"`
	StartTime string `form:"start_time"` // 可选
	EndTime   string `form:"end_time"`   // 可选
}

// StatisticsResponse 统计响应
type StatisticsResponse struct {
	Dimension string              `json:"dimension"`
	Cycle     string              `json:"cycle"`
	StartTime string              `json:"start_time"`
	EndTime   string              `json:"end_time"`
	Summary   StatisticsSummary   `json:"summary"`
	Details   []StatisticsDetail  `json:"details"`
}

// StatisticsSummary 统计摘要
type StatisticsSummary struct {
	TotalIncome float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	NetAmount   float64 `json:"net_amount"`
	RecordCount int64   `json:"record_count"`
}

// StatisticsDetail 统计详情
type StatisticsDetail struct {
	Key         string  `json:"key"`
	Income      float64 `json:"income"`
	Expense     float64 `json:"expense"`
	NetAmount   float64 `json:"net_amount"`
	RecordCount int64   `json:"record_count"`
	Percentage  float64 `json:"percentage"`
}

// GetStatistics 获取统计数据
func GetStatistics(c *gin.Context) {
	var req StatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	// 解析时间范围（可选）
	var startTime, endTime time.Time
	var timeFilterEnabled bool

	if req.StartTime != "" && req.EndTime != "" {
		var err error
		startTime, err = time.Parse("2006-01-02", req.StartTime)
		if err != nil {
			c.JSON(200, utils.ErrorResponse(2001, "开始时间格式错误"))
			return
		}

		endTime, err = time.Parse("2006-01-02", req.EndTime)
		if err != nil {
			c.JSON(200, utils.ErrorResponse(2001, "结束时间格式错误"))
			return
		}

		// 设置结束时间为当天的23:59:59
		endTime = endTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		timeFilterEnabled = true
	}

	var totalIncome, totalExpense float64
	var recordCount int64

	// 获取总收入（只统计审核通过的记录）
	incomeQuery := database.DB.Model(&models.Transaction{}).Where("is_deleted = ?", 0).Where("status = ?", 1).Where("amount > 0")
	if timeFilterEnabled {
		incomeQuery = incomeQuery.Where("transaction_time >= ?", startTime).Where("transaction_time <= ?", endTime)
	}
	incomeQuery.Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome)

	// 获取总支出（只统计审核通过的记录）
	expenseQuery := database.DB.Model(&models.Transaction{}).Where("is_deleted = ?", 0).Where("status = ?", 1).Where("amount < 0")
	if timeFilterEnabled {
		expenseQuery = expenseQuery.Where("transaction_time >= ?", startTime).Where("transaction_time <= ?", endTime)
	}
	expenseQuery.Select("COALESCE(SUM(amount), 0)").Scan(&totalExpense)

	// 获取记录总数（只统计审核通过的记录）
	countQuery := database.DB.Model(&models.Transaction{}).Where("is_deleted = ?", 0).Where("status = ?", 1)
	if timeFilterEnabled {
		countQuery = countQuery.Where("transaction_time >= ?", startTime).Where("transaction_time <= ?", endTime)
	}
	countQuery.Count(&recordCount)

	netAmount := totalIncome + totalExpense // 总支出是负数

	// 获取明细数据
	details := getStatisticsDetails(req.Dimension, startTime, endTime, netAmount, timeFilterEnabled)

	response := StatisticsResponse{
		Dimension: req.Dimension,
		Cycle:     req.Cycle,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Summary: StatisticsSummary{
			TotalIncome: totalIncome,
			TotalExpense: totalExpense,
			NetAmount:   netAmount,
			RecordCount: recordCount,
		},
		Details: details,
	}

	// 确保 details 不为 nil
	if details == nil {
		details = []StatisticsDetail{}
	}

	c.JSON(200, utils.SuccessResponse(response))
}

// getStatisticsDetails 获取统计明细
func getStatisticsDetails(dimension string, startTime, endTime time.Time, netAmount float64, timeFilterEnabled bool) []StatisticsDetail {
	var results []StatisticsDetail

	switch dimension {
	case "project":
		// 按项目统计（只统计审核通过的记录）
		query := database.DB.Table("t_transaction").
			Select("COALESCE(t_project.name, '未分类') as `key`, "+
				"COALESCE(SUM(CASE WHEN t_transaction.amount > 0 THEN t_transaction.amount ELSE 0 END), 0) as income, "+
				"COALESCE(SUM(CASE WHEN t_transaction.amount < 0 THEN t_transaction.amount ELSE 0 END), 0) as expense, "+
				"COALESCE(SUM(t_transaction.amount), 0) as net_amount, "+
				"COUNT(*) as record_count").
			Joins("LEFT JOIN t_project ON t_project.project_id = t_transaction.project_id AND t_project.is_deleted = 0").
			Where("t_transaction.is_deleted = 0").
			Where("t_transaction.status = ?", 1)
		if timeFilterEnabled {
			query = query.Where("t_transaction.transaction_time >= ?", startTime).
				Where("t_transaction.transaction_time <= ?", endTime)
		}
		rows, err := query.Group("t_transaction.project_id, t_project.name").Rows()
		if err != nil {
			return results
		}
		defer rows.Close()

		for rows.Next() {
			var key string
			var income, expense, netAmt float64
			var count int64
			rows.Scan(&key, &income, &expense, &netAmt, &count)
			percentage := 0.0
			if netAmount != 0 {
				percentage = (netAmt / netAmount) * 100
			}
			results = append(results, StatisticsDetail{
				Key:         key,
				Income:      income,
				Expense:     expense,
				NetAmount:   netAmt,
				RecordCount: count,
				Percentage:  percentage,
			})
		}

	case "person":
		// 按人员统计（只统计审核通过的记录）
		query := database.DB.Table("t_transaction").
			Select("COALESCE(t_user.real_name, '未知') as `key`, "+
				"COALESCE(SUM(CASE WHEN t_transaction.amount > 0 THEN t_transaction.amount ELSE 0 END), 0) as income, "+
				"COALESCE(SUM(CASE WHEN t_transaction.amount < 0 THEN t_transaction.amount ELSE 0 END), 0) as expense, "+
				"COALESCE(SUM(t_transaction.amount), 0) as net_amount, "+
				"COUNT(*) as record_count").
			Joins("LEFT JOIN t_user ON t_user.user_id = t_transaction.person_id AND t_user.is_deleted = 0").
			Where("t_transaction.is_deleted = 0").
			Where("t_transaction.status = ?", 1)
		if timeFilterEnabled {
			query = query.Where("t_transaction.transaction_time >= ?", startTime).
				Where("t_transaction.transaction_time <= ?", endTime)
		}
		rows, err := query.Group("t_transaction.person_id, t_user.real_name").Rows()
		if err != nil {
			return results
		}
		defer rows.Close()

		for rows.Next() {
			var key string
			var income, expense, netAmt float64
			var count int64
			rows.Scan(&key, &income, &expense, &netAmt, &count)
			percentage := 0.0
			if netAmount != 0 {
				percentage = (netAmt / netAmount) * 100
			}
			results = append(results, StatisticsDetail{
				Key:         key,
				Income:      income,
				Expense:     expense,
				NetAmount:   netAmt,
				RecordCount: count,
				Percentage:  percentage,
			})
		}

	case "category":
		// 按分类统计（只统计审核通过的记录）
		query := database.DB.Table("t_transaction").
			Select("COALESCE(t_category.name, '未分类') as `key`, "+
				"COALESCE(SUM(CASE WHEN t_transaction.amount > 0 THEN t_transaction.amount ELSE 0 END), 0) as income, "+
				"COALESCE(SUM(CASE WHEN t_transaction.amount < 0 THEN t_transaction.amount ELSE 0 END), 0) as expense, "+
				"COALESCE(SUM(t_transaction.amount), 0) as net_amount, "+
				"COUNT(*) as record_count").
			Joins("LEFT JOIN t_category ON t_category.category_id = t_transaction.category_id AND t_category.is_deleted = 0").
			Where("t_transaction.is_deleted = 0").
			Where("t_transaction.status = ?", 1)
		if timeFilterEnabled {
			query = query.Where("t_transaction.transaction_time >= ?", startTime).
				Where("t_transaction.transaction_time <= ?", endTime)
		}
		rows, err := query.Group("t_transaction.category_id, t_category.name").Rows()
		if err != nil {
			return results
		}
		defer rows.Close()

		for rows.Next() {
			var key string
			var income, expense, netAmt float64
			var count int64
			rows.Scan(&key, &income, &expense, &netAmt, &count)
			percentage := 0.0
			if netAmount != 0 {
				percentage = (netAmt / netAmount) * 100
			}
			results = append(results, StatisticsDetail{
				Key:         key,
				Income:      income,
				Expense:     expense,
				NetAmount:   netAmt,
				RecordCount: count,
				Percentage:  percentage,
			})
		}
	}

	return results
}

// ExportReport 导出报表
func ExportReport(c *gin.Context) {
	// 报表导出功能可以后续实现
	c.JSON(200, utils.ErrorResponse(2001, "功能开发中"))
}
