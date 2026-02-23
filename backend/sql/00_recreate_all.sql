-- ====================================
-- 账务管理系统 - 完全重建脚本
-- 版本: v2.1
-- 说明: 删除所有表并重新创建，用于完全重置数据库
-- ====================================

USE `account_management`;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ====================================
-- 删除所有表（按依赖关系逆序）
-- ====================================
DROP TABLE IF EXISTS `t_attachment`;
DROP TABLE IF EXISTS `t_transaction`;
DROP TABLE IF EXISTS `t_category`;
DROP TABLE IF EXISTS `t_project`;
DROP TABLE IF EXISTS `t_user`;
DROP TABLE IF EXISTS `t_department`;
DROP TABLE IF EXISTS `t_operation_log`;

SET FOREIGN_KEY_CHECKS = 1;

-- ====================================
-- 执行创建表脚本
-- ====================================
SOURCE 02_create_tables.sql;

-- ====================================
-- 执行初始化数据脚本
-- ====================================
SOURCE 03_init_data.sql;

-- ====================================
-- 完成
-- ====================================
SELECT '====================================' AS message;
SELECT '数据库完全重建完成！' AS message;
SELECT '====================================' AS message;
