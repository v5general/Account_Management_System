USE account_management;

-- 插入默认管理员账号
-- 密码: admin123 (使用bcrypt加密)
INSERT INTO t_user (user_id, username, password, role, department, status) VALUES
('admin001', 'admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'ADMIN', '管理部', 1);

-- 插入测试用户
INSERT INTO t_user (user_id, username, password, role, department, status) VALUES
('user001', 'caiwu', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'FINANCE', '财务部', 1),
('user002', 'yuangong1', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'EMPLOYEE', '技术部', 1),
('user003', 'yuangong2', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'EMPLOYEE', '市场部', 1);

-- 插入预设费用分类
INSERT INTO t_category (category_id, name, description, sort_order, is_deleted) VALUES
('category001', '工资', '员工工资发放', 1, 0),
('category002', '设备采购', '办公设备、生产设备采购', 2, 0),
('category003', '服务购买', '外部服务采购', 3, 0),
('category004', '差旅费', '出差交通、住宿费用', 4, 0),
('category005', '业务招待费', '客户招待费用', 5, 0),
('category006', '办公费用', '日常办公用品采购', 6, 0),
('category007', '其他', '其他费用', 99, 0);
