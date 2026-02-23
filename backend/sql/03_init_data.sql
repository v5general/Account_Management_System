-- ====================================
-- 账务管理系统 - 初始数据脚本
-- 版本: v2.1
-- 说明: 插入系统预设的初始数据
-- ====================================

USE `account_management`;
SET NAMES utf8mb4;

-- ====================================
-- 1. 初始化部门数据
-- ====================================
INSERT INTO `t_department` (`department_id`, `name`, `description`, `sort_order`) VALUES
('dept001', '管理部', '公司管理职能部门', 1),
('dept002', '财务部', '财务管理职能部门', 2),
('dept003', '技术部', '技术研发部门', 3),
('dept004', '市场部', '市场营销部门', 4),
('dept005', '人事部', '人力资源部门', 5)
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `description` = VALUES(`description`),
  `sort_order` = VALUES(`sort_order`);

-- ====================================
-- 2. 初始化管理员用户
-- 密码: admin123 (bcrypt加密后)
-- ====================================
INSERT INTO `t_user` (`user_id`, `username`, `password`, `real_name`, `role`, `department_id`, `status`) VALUES
('admin001', 'admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '系统管理员', 'ADMIN', 'dept001', 1)
ON DUPLICATE KEY UPDATE
  `password` = VALUES(`password`),
  `real_name` = VALUES(`real_name`),
  `role` = VALUES(`role`),
  `department_id` = VALUES(`department_id`),
  `status` = VALUES(`status`);

-- ====================================
-- 3. 初始化费用分类数据
-- ====================================
INSERT INTO `t_category` (`category_id`, `name`, `type`, `description`, `sort_order`) VALUES
-- 收入分类
('category001', '服务收入', 'INCOME', '提供服务获得的收入', 1),
('category002', '销售收入', 'INCOME', '产品销售收入', 2),
('category003', '其他收入', 'INCOME', '其他收入来源', 99),
-- 支出分类
('category004', '工资', 'EXPENSE', '员工工资发放', 1),
('category005', '设备采购', 'EXPENSE', '办公设备、生产设备采购', 2),
('category006', '服务购买', 'EXPENSE', '外部服务采购', 3),
('category007', '差旅费', 'EXPENSE', '出差交通、住宿费用', 4),
('category008', '业务招待费', 'EXPENSE', '客户招待费用', 5),
('category009', '办公费用', 'EXPENSE', '日常办公用品采购', 6),
('category010', '其他', 'EXPENSE', '其他费用', 99)
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `type` = VALUES(`type`),
  `description` = VALUES(`description`),
  `sort_order` = VALUES(`sort_order`);

-- ====================================
-- 初始化完成 - 统计信息
-- ====================================
SELECT '初始数据插入完成！' AS message;
SELECT '部门数量:' AS info, COUNT(*) AS count FROM t_department WHERE is_deleted = 0
UNION ALL
SELECT '用户数量:', COUNT(*) FROM t_user WHERE is_deleted = 0
UNION ALL
SELECT '分类数量:', COUNT(*) FROM t_category WHERE is_deleted = 0;

-- ====================================
-- 默认登录信息
-- ====================================
SELECT '====================================' AS info;
SELECT '默认登录账号:' AS info;
SELECT '用户名: admin' AS info;
SELECT '密码: admin123' AS info;
SELECT '====================================' AS info;
