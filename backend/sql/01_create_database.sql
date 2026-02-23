-- ====================================
-- 账务管理系统 - 数据库创建脚本
-- 版本: v2.1
-- ====================================

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `account_management`
DEFAULT CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci
COMMENT '账务管理系统数据库';

-- 使用数据库
USE `account_management`;

-- 完成
SELECT '数据库创建完成！' AS message;
