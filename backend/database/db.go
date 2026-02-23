package database

import (
	"fmt"
	"time"

	"account-management-system/backend/config"
	"account-management-system/backend/models"
	"account-management-system/backend/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) error {
	dsn := cfg.Database.GetDSN()

	var logLevel logger.LogLevel
	if cfg.Server.Mode == "debug" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
	})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	// 自动迁移表结构
	if err := AutoMigrate(); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 初始化数据
	if err := InitData(); err != nil {
		return fmt.Errorf("初始化数据失败: %w", err)
	}

	return nil
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.Department{},
		&models.Project{},
		&models.User{},
		&models.Category{},
		&models.Transaction{},
		&models.Attachment{},
		&models.OperationLog{},
	)
}

// InitData 初始化数据
func InitData() error {
	// 检查是否已有部门
	var deptCount int64
	DB.Model(&models.Department{}).Count(&deptCount)
	if deptCount == 0 {
		// 创建预设部门
		departments := []models.Department{
			{DepartmentID: "dept001", Name: "管理部", Description: "公司管理职能部门", SortOrder: 1},
			{DepartmentID: "dept002", Name: "财务部", Description: "财务管理职能部门", SortOrder: 2},
			{DepartmentID: "dept003", Name: "技术部", Description: "技术研发部门", SortOrder: 3},
			{DepartmentID: "dept004", Name: "市场部", Description: "市场营销部门", SortOrder: 4},
			{DepartmentID: "dept005", Name: "人事部", Description: "人力资源部门", SortOrder: 5},
		}
		if err := DB.Create(&departments).Error; err != nil {
			return fmt.Errorf("创建部门失败: %w", err)
		}
	}

	// 检查是否已有用户
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		// 创建默认管理员
		hashedPassword, _ := utils.HashPassword("admin123")
		users := []models.User{
			{
				UserID:       "admin001",
				Username:     "admin",
				Password:     hashedPassword,
				RealName:     "系统管理员",
				Role:         "ADMIN",
				DepartmentID: "dept001",
				Status:       1,
			},
		}
		if err := DB.Create(&users).Error; err != nil {
			return fmt.Errorf("创建用户失败: %w", err)
		}
	}

	// 检查是否已有分类
	var catCount int64
	DB.Model(&models.Category{}).Count(&catCount)
	if catCount == 0 {
		// 创建预设费用分类
		categories := []models.Category{
			{CategoryID: "category001", Name: "工资", Description: "员工工资发放", SortOrder: 1},
			{CategoryID: "category002", Name: "设备采购", Description: "办公设备、生产设备采购", SortOrder: 2},
			{CategoryID: "category003", Name: "服务购买", Description: "外部服务采购", SortOrder: 3},
			{CategoryID: "category004", Name: "差旅费", Description: "出差交通、住宿费用", SortOrder: 4},
			{CategoryID: "category005", Name: "业务招待费", Description: "客户招待费用", SortOrder: 5},
			{CategoryID: "category006", Name: "办公费用", Description: "日常办公用品采购", SortOrder: 6},
			{CategoryID: "category007", Name: "其他", Description: "其他费用", SortOrder: 99},
		}
		if err := DB.Create(&categories).Error; err != nil {
			return fmt.Errorf("创建分类失败: %w", err)
		}
	}

	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
