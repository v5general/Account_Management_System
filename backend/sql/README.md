# 账务管理系统 - SQL 脚本说明

## 文件说明

| 文件名 | 说明 | 用途 |
|--------|------|------|
| `00_recreate_all.sql` | 完全重建脚本 | 删除所有表并重新创建，用于完全重置数据库 |
| `01_create_database.sql` | 数据库创建脚本 | 创建数据库（首次安装时使用） |
| `02_create_tables.sql` | 表结构创建脚本 | 创建所有表结构（不删除现有数据） |
| `03_init_data.sql` | 初始数据脚本 | 插入系统预设的初始数据 |

## 使用方法

### 首次安装

```bash
# 1. 创建数据库
mysql -u root -p < 01_create_database.sql

# 2. 创建表结构
mysql -u root -p account_management < 02_create_tables.sql

# 3. 插入初始数据
mysql -u root -p account_management < 03_init_data.sql
```

### 完全重建（清空所有数据）

```bash
# 使用重建脚本（包含删除表、创建表、插入数据）
mysql -u root -p account_management < 00_recreate_all.sql
```

### Windows MySQL 客户端

```cmd
mysql -u root -p account_management < 00_recreate_all.sql
```

### 使用 SOURCE 命令（MySQL 命令行）

```sql
mysql -u root -p
USE account_management;
SOURCE 00_recreate_all.sql;
```

## 数据库表结构

### 核心表

| 表名 | 说明 | 主要字段 |
|------|------|----------|
| t_department | 部门表 | department_id, name, description, sort_order, is_deleted |
| t_user | 用户表 | user_id, username, password, real_name, role, department_id, is_deleted |
| t_project | 项目表 | project_id, name, description, department_id, status, is_deleted |
| t_category | 费用分类表 | category_id, name, **type**, description, sort_order, is_deleted |
| t_transaction | 收支流水表 | record_id, amount, category_id, project_id, person_id, status, is_deleted |
| t_attachment | 凭证附件表 | attachment_id, record_id, file_name, file_path, is_deleted |
| t_operation_log | 操作日志表 | log_id, user_id, operation_type, module, description |

## 版本 v2.1 主要变更

1. **费用分类表 (t_category)** 新增 `type` 字段
   - INCOME: 收入分类
   - EXPENSE: 支出分类

2. **所有表** 统一使用 `is_deleted` 字段实现软删除
   - 0: 未删除
   - 1: 已删除

3. **用户表 (t_user)** 新增 `real_name` 字段存储真实姓名

4. **操作日志表 (t_operation_log)** 移除软删除功能

## 默认登录信息

```
用户名: admin
密码: admin123
角色: 管理员
```

## 初始数据

### 预设部门（5个）
- 管理部
- 财务部
- 技术部
- 市场部
- 人事部

### 预设费用分类（10个）

**收入分类（3个）：**
- 服务收入
- 销售收入
- 其他收入

**支出分类（7个）：**
- 工资
- 设备采购
- 服务购买
- 差旅费
- 业务招待费
- 办公费用
- 其他

## 注意事项

1. 执行完全重建脚本会删除所有现有数据，请谨慎使用
2. 默认密码为 `admin123`，生产环境请及时修改
3. 建议定期备份数据库
4. 字符集为 utf8mb4，支持中文和特殊字符
