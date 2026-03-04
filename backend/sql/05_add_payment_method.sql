-- 添加支付方式字段到t_transaction表
-- 执行日期: 2026-03-03

USE account_management;

-- 添加payment_method字段
ALTER TABLE t_transaction
ADD COLUMN payment_method VARCHAR(50) DEFAULT '' COMMENT '支付方式' AFTER person_id;

-- 验证字段已添加
DESC t_transaction;
