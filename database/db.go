package database

import (
	"account-management-system/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config 数据库配置
type Config struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

// DBService 数据库服务
type DBService struct {
	*gorm.DB
}

// InitDB 初始化数据库连接
func InitDB(cfg Config) (*DBService, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return &DBService{DB: db}, nil
}

// AutoMigrate 自动迁移表结构
func (db *DBService) AutoMigrate(models ...interface{}) error {
	return db.AutoMigrate(models...)
}

// Find 查询数据
func (db *DBService) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return db.Find(dest, conds...)
}

// Create 创建数据
func (db *DBService) Create(value interface{}) *gorm.DB {
	return db.Create(value)
}

// Save 保存数据
func (db *DBService) Save(value interface{}) *gorm.DB {
	return db.Save(value)
}

// Delete 删除数据
func (db *DBService) Delete(value interface{}) *gorm.DB {
	return db.Delete(value)
}

// Where 条件查询
func (db *DBService) Where(query interface{}, args ...interface{}) *gorm.DB {
	return db.Where(query, args...)
}

// First 查询单条数据
func (db *DBService) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return db.First(dest, conds...)
}
