-- 创建支付方式表并迁移数据
-- 执行日期: 2026-03-17

USE account_management;

-- 1. 创建支付方式表
CREATE TABLE IF NOT EXISTS `t_payment_method` (
  `payment_method_id` VARCHAR(32) NOT NULL PRIMARY KEY COMMENT '支付方式唯一标识',
  `name` VARCHAR(50) NOT NULL COMMENT '支付方式名称',
  `description` VARCHAR(200) DEFAULT NULL COMMENT '描述',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
  `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '删除标识：0-未删除，1-已删除',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付方式表';

-- 2. 插入初始数据
INSERT INTO t_payment_method (payment_method_id, name, sort_order) VALUES
('pm_001', '现金', 1),
('pm_002', '公司转账', 2),
('pm_003', '微信', 3),
('pm_004', '支付宝', 4),
('pm_005', '银行转账', 5),
('pm_006', '支票', 6),
('pm_007', '其他', 7);

-- 3. 添加新字段 payment_method_id
ALTER TABLE t_transaction ADD COLUMN payment_method_id VARCHAR(32) DEFAULT NULL COMMENT '支付方式ID' AFTER person_id;

-- 4. 自动迁移现有数据（根据名称匹配）
UPDATE t_transaction t
JOIN t_payment_method pm ON t.payment_method = pm.name
SET t.payment_method_id = pm.payment_method_id;

-- 5. 删除旧字段
ALTER TABLE t_transaction DROP COLUMN payment_method;

-- 6. 添加外键约束
ALTER TABLE t_transaction
ADD CONSTRAINT fk_payment_method
FOREIGN KEY (payment_method_id) REFERENCES t_payment_method(payment_method_id)
ON DELETE SET NULL ON UPDATE CASCADE;

-- 验证表结构
DESC t_payment_method;
DESC t_transaction;
