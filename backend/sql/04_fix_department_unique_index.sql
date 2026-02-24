-- ====================================
-- 修复部门表唯一索引问题
-- 问题：软删除后无法创建同名部门
-- 原因：name 字段有 UNIQUE 约束
-- 解决：删除 name 字段的唯一索引
-- ====================================

USE `account_management`;

-- 删除 t_department 表的 name 唯一索引
-- MySQL 自动为 UNIQUE 字段创建索引，索引名为 'idx_t_department_name'
ALTER TABLE `t_department` DROP INDEX `idx_t_department_name`;

-- 验证索引是否已删除
SELECT '索引删除完成！现在软删除后可以创建同名部门了。' AS message;
