# 账务管理系统

轻量化、易操作、权限清晰的BS架构账务管理系统，实现收支流水的规范化记录、凭证的统一管理及多维度数据统计分析。

## 技术栈

### 后端
- Go 1.18+
- Gin Web框架
- GORM (MySQL)
- Redis (可选)
- bcrypt 密码加密
- JWT 认证

### 前端
- Vue 3.4+
- TypeScript
- Element Plus UI
- ECharts 图表
- Pinia 状态管理
- Vue Router

### 数据库与存储
- MySQL 8.0+ (结构化数据)
- Redis 6.0+ (缓存，可选)
- 阿里云OSS / 腾讯云COS (凭证附件存储)

## 项目结构

```
account-management-system/
├── backend/                    # 后端代码
│   ├── main.go                 # 应用入口
│   ├── go.mod                  # Go 模块配置
│   ├── config/                 # 配置文件
│   ├── database/               # 数据库配置
│   ├── models/                 # 数据模型
│   ├── controllers/            # 控制器
│   ├── middlewares/            # 中间件
│   ├── routes/                 # 路由配置
│   └── utils/                  # 工具函数
├── frontend/                   # 前端代码
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts
│   └── src/
│       ├── main.ts
│       ├── api/                # API 封装
│       ├── router/             # 路由配置
│       ├── store/              # 状态管理
│       ├── components/         # 公共组件
│       └── views/              # 页面组件
├── database/                   # 数据库脚本
│   ├── schema.sql              # 表结构创建
│   └── init_data.sql           # 初始数据
└── docs/                       # 项目文档
    └── 需求规格说明书.md
```

## 功能模块

### 1. 用户认证
- 用户登录/注销
- JWT 令牌认证
- 密码错误 5 次锁定 30 分钟
- 令牌有效期 24 小时

### 2. 收入登记
- 费用类别（可选）
- 来源项目
- 金额（正数标识）
- 交易时间、备注
- 凭证附件上传（至少 1 个）

### 3. 支出登记
- 费用类别
- 关联项目
- 关联人员
- 金额（负数标识）
- 交易时间、备注
- 凭证附件上传（至少 1 个）

### 4. 收支记录查询
- 全局查询（管理者/财务人员）
- 按时间范围、费用类别、项目、人员筛选
- 个人查询（员工，仅本人记录）
- 凭证查看/下载

### 5. 费用分类管理
- 类别新增、编辑、删除
- 类别查询（支持模糊搜索）

### 6. 统计分析
- 多维度统计（项目/人员/类别）
- 多周期统计（月度/季度/年度）
- 可视化报表（表格、柱状图、折线图、饼图）
- 报表导出（Excel/PDF）

### 7. 用户管理（管理员）
- 用户注册
- 权限分配
- 密码重置
- 用户列表查询

### 8. 系统设置（管理员）
- 云存储配置
- 操作日志查看
- 数据备份

## 安装和运行

### 前置要求
1. Go 1.18 或更高版本
2. Node.js 18 或更高版本
3. MySQL 8.0 或更高版本
4. Redis 6.0 或更高版本（可选）

### 数据库初始化
```bash
mysql -u root -p < database/schema.sql
mysql -u root -p < database/init_data.sql
```

### 后端依赖安装
```bash
cd backend
go mod download
```

### 前端依赖安装
```bash
cd frontend
npm install
```

### 配置文件
编辑 `backend/config/config.yaml`：
```yaml
server:
  port: "8080"
  mode: debug

database:
  host: localhost
  port: 3306
  database: account_management
  username: root
  password: your_password

jwt:
  secret: your_jwt_secret_key
  expire: 86400
```

### 开发模式运行

**启动后端：**
```bash
cd backend
go run main.go
```

**启动前端：**
```bash
cd frontend
npm run dev
```

### 生产环境构建

**构建后端：**
```bash
cd backend
go build -o account-management main.go
```

**构建前端：**
```bash
cd frontend
npm run build
```

## 默认账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | admin123 | 管理员 |

## 用户角色

| 角色 | 权限 |
|------|------|
| **管理员 (ADMIN)** | 全部权限 |
| **财务人员 (FINANCE)** | 收支登记、费用分类、统计报表 |
| **员工 (EMPLOYEE)** | 查看个人收支记录 |

## API接口

### 认证相关
- POST `/api/v1/auth/login` - 用户登录
- POST `/api/v1/auth/logout` - 用户注销
- GET `/api/v1/auth/me` - 获取当前用户信息

### 收支管理
- POST `/api/v1/transactions` - 创建收支记录
- GET `/api/v1/transactions` - 查询收支记录列表
- GET `/api/v1/transactions/:id` - 获取单条记录详情
- PUT `/api/v1/transactions/:id` - 修改收支记录
- DELETE `/api/v1/transactions/:id` - 删除收支记录

### 费用分类
- POST `/api/v1/categories` - 创建费用分类
- GET `/api/v1/categories` - 获取分类列表
- PUT `/api/v1/categories/:id` - 修改分类信息
- DELETE `/api/v1/categories/:id` - 删除分类

## 核心业务规则

1. **记账模式**：采用简化流水账模式，正数代表收入，负数代表支出
2. **凭证强制**：每笔收支记录必须关联至少 1 个凭证附件
3. **角色隔离**：员工只能查看本人关联的记录，管理者和财务人员可查看所有记录
4. **金额规范**：所有金额字段保留 2 位小数
5. **时区统一**：所有时间字段使用 UTC+8 时区
6. **审计追踪**：关键操作记录日志，支持审计追溯

## 文档

- [需求规格说明书](./docs/需求规格说明书.md)
- [接口文档](./docs/接口文档.md)
- [部署文档](./docs/部署文档.md)
- [使用手册](./docs/使用手册.md)

## 许可证

MIT
