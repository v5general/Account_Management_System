package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 数据库连接信息
	dsn := "root:060928@tcp(localhost:3306)/account_management?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatal("数据库 ping 失败:", err)
	}

	fmt.Println("开始重建数据库表...")

	// 设置外键检查为0
	if _, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		log.Fatal("设置外键检查失败:", err)
	}

	// 删除所有表
	tables := []string{
		"t_attachment",
		"t_transaction",
		"t_category",
		"t_project",
		"user",  // 删除不带前缀的（如果存在）
		"t_user",
		"t_department",
		"t_operation_log",
	}

	for _, table := range tables {
		if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)); err != nil {
			log.Printf("删除表 %s 失败: %v", table, err)
		} else {
			fmt.Printf("删除表 %s 成功\n", table)
		}
	}

	// 设置外键检查恢复
	if _, err := db.Exec("SET FOREIGN_KEY_CHECKS = 1"); err != nil {
		log.Fatal("恢复外键检查失败:", err)
	}

	fmt.Println("\n旧表已删除，请重启后端程序以自动创建新表...")
	fmt.Println("或者手动执行 SQL 脚本完成初始化数据")
}
