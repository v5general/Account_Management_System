-- ====================================
-- 账务管理系统 - 表结构创建脚本
-- 版本: v2.1
-- 说明: 仅创建表结构，不删除现有表
-- ====================================

USE `account_management`;
SET NAMES utf8mb4;

-- ====================================
-- 1. 部门表 (t_department)
-- ====================================
CREATE TABLE IF NOT EXISTS `t_department` (
  `department_id` VARCHAR(32) NOT NULL PRIMARY KEY COMMENT '部门唯一标识',
  `name` VARCHAR(50) NOT NULL UNIQUE COMMENT '部门名称',
  `description` VARCHAR(200) DEFAULT NULL COMMENT '部门描述',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
  `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '删除标识：0-未删除，1-已删除',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_is_deleted` (`is_deleted`),
  INDEX `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门表';

-- ====================================
-- 2. 用户表 (t_user)
-- ====================================
CREATE TABLE IF NOT EXISTS `t_user` (
  `user_id` VARCHAR(32) NOT NULL PRIMARY KEY COMMENT '用户唯一标识',
  `username` VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名（仅数字、字母、下划线）',
  `password` VARCHAR(255) NOT NULL COMMENT '加密后的密码（bcrypt）',
  `real_name` VARCHAR(50) NOT NULL COMMENT '真实姓名',
  `role` VARCHAR(20) NOT NULL DEFAULT 'EMPLOYEE' COMMENT '角色：ADMIN/FINANCE/EMPLOYEE',
  `department_id` VARCHAR(32) DEFAULT NULL COMMENT '所属部门ID',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
  `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '删除标识：0-未删除，1-已删除',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_department_id` (`department_id`),
  INDEX `idx_is_deleted` (`is_deleted`),
  INDEX `idx_role` (`role`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ====================================
-- 3. 项目表 (t_project)
-- ====================================
CREATE TABLE IF NOT EXISTS `t_project` (
  `project_id` VARCHAR(32) NOT NULL PRIMARY KEY COMMENT '项目唯一标识',
  `name` VARCHAR(100) NOT NULL COMMENT '项目名称',
  `description` VARCHAR(500) DEFAULT NULL COMMENT '项目描述',
  `department_id` VARCHAR(32) DEFAULT NULL COMMENT '关联部门ID',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-进行中，0-已结束',
  `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '删除标识：0-未删除，1-已删除',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_department_id` (`department_id`),
  INDEX `idx_is_deleted` (`is_deleted`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目表';

-- ====================================
-- 4. 费用分类表 (t_category)
-- ====================================
CREATE TABLE IF NOT EXISTS `t_category` (
  `category_id` VARCHAR(32) NOT NULL PRIMARY KEY COMMENT '分类唯一标识',
  `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
  `type` VARCHAR(20) NOT NULL COMMENT '分类类型：INCOME（收入）、EXPENSE（支出）',
  `description` VARCHAR(200) DEFAULT NULL COMMENT '分类描述',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
  `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '删除标识：0-未删除，1-已删除',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_type` (`type`),
  INDEX `idx_is_deleted` (`is_deleted`),
  INDEX `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='费用分类表';

-- ====================================
-- 5. 收支流水表 (t_transaction)
-- ====================================
CREATE TABLE IF NOT EXISTS `t_transaction` (
  `record_id` VARCHAR(32) NOT NULL PRIMARY KEY COMMENT '记录唯一标识',
  `amount` DECIMAL(15,2) NOT NULL COMMENT '金额（正数=收入，负数=支出）',
  `category_id` VARCHAR(32) DEFAULT NULL COMMENT '关联费用分类ID',
  `project_id` VARCHAR(32) DEFAULT NULL COMMENT '关联项目ID',
  `person_id` VARCHAR(32) DEFAULT NULL COMMENT '关联人员user_id',
  `transaction_time` DATETIME NOT NULL COMMENT '交易发生时间',
  `remark` VARCHAR(500) DEFAULT NULL COMMENT '交易备注',
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态：0-待审核，1-已审核，2-已驳回',
  `creator_id` VARCHAR(32) NOT NULL COMMENT '录入人user_id',
  `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '删除标识：0-未删除，1-已删除',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_category_id` (`category_id`),
  INDEX `idx_project_id` (`project_id`),
  INDEX `idx_person_id` (`person_id`),
  INDEX `idx_creator_id` (`creator_id`),
  INDEX `idx_transaction_time` (`transaction_time`),
  INDEX `idx_status` (`status`),
  INDEX `idx_is_deleted` (`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收支流水表';

-- ====================================
-- 6. 凭证附件表 (t_attachment)
-- ====================================
CREATE TABLE IF NOT EXISTS `t_attachment` (
  `attachment_id` VARCHAR(32) NOT NULL PRIMARY KEY COMMENT '附件唯一标识',
  `record_id` VARCHAR(32) NOT NULL COMMENT '关联收支记录ID',
  `file_name` VARCHAR(255) NOT NULL COMMENT '附件原始名称',
  `file_path` VARCHAR(500) NOT NULL COMMENT '云存储路径',
  `file_size` BIGINT NOT NULL COMMENT '文件大小（字节）',
  `file_type` VARCHAR(20) NOT NULL COMMENT '文件类型：image/pdf/other',
  `upload_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `uploader_id` VARCHAR(32) NOT NULL COMMENT '上传人user_id',
  `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '删除标识：0-未删除，1-已删除',
  INDEX `idx_record_id` (`record_id`),
  INDEX `idx_uploader_id` (`uploader_id`),
  INDEX `idx_is_deleted` (`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='凭证附件表';

-- ====================================
-- 7. 操作日志表 (t_operation_log)
-- ====================================
CREATE TABLE IF NOT EXISTS `t_operation_log` (
  `log_id` VARCHAR(32) NOT NULL PRIMARY KEY COMMENT '日志唯一标识',
  `user_id` VARCHAR(32) NOT NULL COMMENT '操作人user_id',
  `operation_type` VARCHAR(50) NOT NULL COMMENT '操作类型：LOGIN/CREATE/UPDATE/DELETE/APPROVE',
  `module` VARCHAR(50) DEFAULT NULL COMMENT '操作模块',
  `description` VARCHAR(500) DEFAULT NULL COMMENT '操作描述',
  `request_ip` VARCHAR(50) DEFAULT NULL COMMENT '请求IP',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '操作状态：1-成功，0-失败',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_operation_type` (`operation_type`),
  INDEX `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- ====================================
-- 完成
-- ====================================
SELECT '表结构创建完成！' AS message;
