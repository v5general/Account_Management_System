# 账务管理系统

[![Version](https://img.shields.io/badge/version-v1.0.2-blue.svg)](https://github.com)
[![License](https://img.shields.io/badge/license-GPL--3.0-green.svg)](LICENSE)

轻量化、易操作、权限清晰的BS架构账务管理系统，实现收支流水的规范化记录、凭证的统一管理及多维度数据统计分析。

## 功能特性

- **用户认证** - JWT令牌认证，密码错误锁定机制
- **收支管理** - 收入/支出登记、审核、查询
- **费用分类** - 灵活的费用类别管理
- **统计报表** - 多维度数据统计与可视化
- **权限控制** - 管理员/财务/员工三种角色
- **操作日志** - 完整的操作审计追踪
- **版本记录** - 系统版本更新信息展示

## 技术栈

### 后端
- Go 1.18+
- Gin Web框架
- GORM (MySQL)
- JWT 认证
- bcrypt 密码加密

### 前端
- Vue 3.4+
- TypeScript
- Element Plus UI
- ECharts 图表
- Pinia 状态管理
- Vue Router

### 数据库
- MySQL 8.0+

## 项目结构

```
account-management-system/
├── backend/                    # 后端代码
│   ├── main.go                 # 应用入口
│   ├── config/                 # 配置文件
│   ├── database/               # 数据库配置
│   ├── models/                 # 数据模型
│   ├── controllers/            # 控制器
│   ├── middlewares/            # 中间件
│   ├── routes/                 # 路由配置
│   └── utils/                  # 工具函数
├── frontend/                   # 前端代码
│   ├── src/
│   │   ├── api/                # API 封装
│   │   ├── router/             # 路由配置
│   │   ├── store/              # 状态管理
│   │   ├── components/         # 公共组件
│   │   └── views/              # 页面组件
│   │       ├── dashboard/      # 首页仪表盘
│   │       ├── transaction/    # 收支管理
│   │       ├── category/       # 费用分类
│   │       ├── statistics/     # 统计报表
│   │       └── settings/       # 系统设置
│   └── dist/                   # 构建输出
├── database/                   # 数据库脚本
└── docs/                       # 项目文档
```

## 功能模块

### 首页仪表盘
- 总收入、总支出、净收支、记录数统计卡片
- 最近收支记录列表
- 按项目统计饼图

### 收支管理
- **收入登记** - 财务人员登记收入记录
- **支出登记** - 财务人员登记支出记录
- **收支审核** - 管理员审核收支记录
- **收支列表** - 查看和筛选收支记录

### 费用分类
- 类别增删改查
- 支持模糊搜索

### 统计报表
- 多维度统计（项目/人员/类别）
- 可视化图表展示

### 系统设置
- **用户管理** - 用户增删改查、权限分配（管理员）
- **部门管理** - 部门信息管理（管理员）
- **项目管理** - 项目信息管理（管理员）
- **操作日志** - 用户操作记录（管理员/财务）
- **账号管理** - 个人账号信息维护
- **版本记录** - 系统版本更新信息

## 用户角色

| 角色 | 权限说明 |
|------|----------|
| **管理员 (ADMIN)** | 全部权限，包括用户管理、部门管理、项目管理、收支审核 |
| **财务人员 (FINANCE)** | 收支登记、费用分类、统计报表、操作日志查看 |
| **员工 (EMPLOYEE)** | 查看个人收支记录、账号管理 |

## 快速开始

### 前置要求
- Go 1.18+
- Node.js 18+
- MySQL 8.0+

### 数据库初始化
```bash
mysql -u root -p < database/schema.sql
mysql -u root -p < database/init_data.sql
```

### 安装依赖

**后端：**
```bash
cd backend
go mod download
```

**前端：**
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

### 开发模式

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

### 生产构建

**构建后端（Linux x86_64）：**
```bash
cd backend
GOOS=linux GOARCH=amd64 go build -o account-management main.go
```

**构建前端：**
```bash
cd frontend
npx vite build
```

## 默认账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | admin123 | 管理员 |

## 版本历史

### v1.0.2 (2026-03-05)
- 新增支付方式字段，支持多种支付方式选择
- 实现被驳回记录的重新提交功能
- 优化金额显示，采用千位分隔符格式
- 统一使用格式化工具函数
- 优化表格列宽显示效果

### v1.0.1 (2026-03-01)
- 优化登录界面设计，采用左右分栏布局
- 修复首页记录数显示问题
- 更新部署文档

### v1.0.0 (2026-03-01)
- 首个正式版本发布
- 完整的财务管理系统功能

## 文档

- [需求规格说明书](./docs/需求规格说明书.md)
- [接口文档](./docs/接口文档.md)
- [部署文档](./docs/部署文档.md)
- [使用手册](./docs/使用手册.md)

## 许可证

GPL-3.0
