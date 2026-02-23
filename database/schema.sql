-- 账务管理系统数据库表结构
-- 创建时间: 2026-02-15

-- 创建数据库
CREATE DATABASE IF NOT EXISTS account_management DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE account_management;

-- 用户表
DROP TABLE IF EXISTS t_user;
CREATE TABLE t_user (
    user_id VARCHAR(32) PRIMARY KEY COMMENT '用户唯一标识',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '加密后的密码',
    role ENUM('ADMIN', 'EMPLOYEE', 'FINANCE') NOT NULL COMMENT '角色：ADMIN-管理员,EMPLOYEE-员工,FINANCE-财务人员',
    department VARCHAR(100) DEFAULT NULL COMMENT '所属部门',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_username (username),
    INDEX idx_role (role),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 费用分类表
DROP TABLE IF EXISTS t_category;
CREATE TABLE t_category (
    category_id VARCHAR(32) PRIMARY KEY COMMENT '分类唯一标识',
    name VARCHAR(50) NOT NULL UNIQUE COMMENT '分类名称',
    description VARCHAR(200) DEFAULT NULL COMMENT '分类描述',
    sort_order INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
    is_deleted TINYINT NOT NULL DEFAULT 0 COMMENT '删除标记：0-正常，1-已删除',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_name (name),
    INDEX idx_sort_order (sort_order),
    INDEX idx_is_deleted (is_deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='费用分类表';

-- 收支流水表
DROP TABLE IF EXISTS t_transaction;
CREATE TABLE t_transaction (
    record_id VARCHAR(32) PRIMARY KEY COMMENT '记录唯一标识',
    amount DECIMAL(15,2) NOT NULL COMMENT '金额（正数=收入，负数=支出）',
    category_id VARCHAR(32) DEFAULT NULL COMMENT '关联费用分类ID',
    project_name VARCHAR(100) NOT NULL COMMENT '关联项目名称',
    person_id VARCHAR(32) DEFAULT NULL COMMENT '关联人员user_id',
    transaction_time DATETIME NOT NULL COMMENT '交易发生时间',
    remark VARCHAR(500) DEFAULT NULL COMMENT '交易备注',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-已生效，0-待审核，2-已作废',
    creator_id VARCHAR(32) NOT NULL COMMENT '录入人user_id',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_category_id (category_id),
    INDEX idx_person_id (person_id),
    INDEX idx_transaction_time (transaction_time),
    INDEX idx_project_name (project_name),
    INDEX idx_status (status),
    INDEX idx_creator_id (creator_id),
    FOREIGN KEY (category_id) REFERENCES t_category(category_id) ON DELETE SET NULL,
    FOREIGN KEY (person_id) REFERENCES t_user(user_id) ON DELETE SET NULL,
    FOREIGN KEY (creator_id) REFERENCES t_user(user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收支流水表';

-- 凭证附件表
DROP TABLE IF EXISTS t_attachment;
CREATE TABLE t_attachment (
    attachment_id VARCHAR(32) PRIMARY KEY COMMENT '附件唯一标识',
    record_id VARCHAR(32) NOT NULL COMMENT '关联收支记录ID',
    file_name VARCHAR(255) NOT NULL COMMENT '附件原始名称',
    file_path VARCHAR(500) NOT NULL COMMENT '云存储路径',
    file_size BIGINT NOT NULL COMMENT '文件大小（字节）',
    file_type ENUM('image', 'pdf') NOT NULL COMMENT '文件类型',
    upload_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
    uploader_id VARCHAR(32) NOT NULL COMMENT '上传人user_id',
    INDEX idx_record_id (record_id),
    INDEX idx_uploader_id (uploader_id),
    FOREIGN KEY (record_id) REFERENCES t_transaction(record_id) ON DELETE CASCADE,
    FOREIGN KEY (uploader_id) REFERENCES t_user(user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='凭证附件表';

-- 操作日志表
DROP TABLE IF EXISTS t_operation_log;
CREATE TABLE t_operation_log (
    log_id VARCHAR(32) PRIMARY KEY COMMENT '日志唯一标识',
    user_id VARCHAR(32) NOT NULL COMMENT '操作人user_id',
    operation_type ENUM('LOGIN', 'CREATE', 'UPDATE', 'DELETE', 'QUERY', 'EXPORT') NOT NULL COMMENT '操作类型',
    module VARCHAR(50) NOT NULL COMMENT '操作模块',
    description VARCHAR(500) DEFAULT NULL COMMENT '操作描述',
    request_ip VARCHAR(50) DEFAULT NULL COMMENT '请求IP',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '操作状态：1-成功，0-失败',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
    INDEX idx_user_id (user_id),
    INDEX idx_operation_type (operation_type),
    INDEX idx_create_time (create_time),
    FOREIGN KEY (user_id) REFERENCES t_user(user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';
