# Account Management System

[![Version](https://img.shields.io/badge/version-v1.0.2-blue.svg)](https://github.com)
[![License](https://img.shields.io/badge/license-GPL--3.0-green.svg)](LICENSE)

**English** | [简体中文](./README.zh-CN.md)

A lightweight, easy-to-use, BS-architecture accounting management system with clearly defined permissions, enabling standardized recording of revenue and expense transactions, unified voucher management, and multi-dimensional data analysis.

## Features

- **User Authentication** - JWT token auth, lockout mechanism on repeated password failures
- **Transaction Management** - Registration, approval, and querying of income/expense records
- **Expense Categories** - Flexible expense category management
- **Statistical Reports** - Multi-dimensional data analysis and visualization
- **Access Control** - Three roles: Admin / Finance / Employee
- **Operation Logs** - Complete operation audit trail
- **Version History** - System version update information display

## Tech Stack

### Backend
- Go 1.18+
- Gin web framework
- GORM (MySQL)
- JWT authentication
- bcrypt password hashing

### Frontend
- Vue 3.4+
- TypeScript
- Element Plus UI
- ECharts
- Pinia state management
- Vue Router

### Database
- MySQL 8.0+

## Project Structure

```
account-management-system/
├── backend/                    # Backend code
│   ├── main.go                 # Application entry
│   ├── go.mod                  # Go module definition
│   ├── config/                 # Configuration
│   │   ├── config.go           # Config struct
│   │   └── config.yaml         # YAML configuration file
│   ├── cmd/                    # CLI tools
│   │   └── recreate_db/        # Database rebuild tool
│   ├── database/               # Database connection
│   ├── models/                 # Data models
│   │   ├── user.go             # User model
│   │   ├── transaction.go      # Transaction model
│   │   ├── category.go         # Category model
│   │   ├── project.go          # Project model
│   │   ├── department.go       # Department model
│   │   ├── operation_log.go    # Operation log model
│   │   └── attachment.go       # Attachment model
│   ├── controllers/            # Controllers
│   │   ├── auth.go             # Auth controller
│   │   ├── user.go             # User controller
│   │   ├── transaction.go      # Transaction controller
│   │   ├── category.go         # Category controller
│   │   ├── project.go          # Project controller
│   │   ├── department.go       # Department controller
│   │   ├── statistics.go       # Statistics controller
│   │   ├── log.go              # Log controller
│   │   └── attachment.go       # Attachment controller
│   ├── middlewares/            # Middlewares
│   │   ├── auth.go             # JWT auth middleware
│   │   ├── cors.go             # CORS middleware
│   │   └── logger.go           # Logger middleware
│   ├── routes/                 # Route configuration
│   ├── utils/                  # Utility functions
│   │   ├── common.go           # Common utilities
│   │   ├── crypto.go           # Crypto utilities
│   │   ├── jwt.go              # JWT utilities
│   │   └── oss.go              # OSS storage utilities
│   ├── sql/                    # SQL scripts
│   │   ├── 01_create_database.sql
│   │   ├── 02_create_tables.sql
│   │   ├── 03_init_data.sql
│   │   ├── 04_fix_department_unique_index.sql
│   │   └── 05_add_payment_method.sql
│   └── uploads/                # File upload directory
├── frontend/                   # Frontend code
│   ├── src/
│   │   ├── main.ts             # Application entry
│   │   ├── App.vue             # Root component
│   │   ├── api/                # API wrappers
│   │   ├── router/             # Routing configuration
│   │   ├── store/              # State management
│   │   ├── components/         # Shared components
│   │   ├── utils/              # Utilities
│   │   │   ├── request.ts      # HTTP request wrapper
│   │   │   └── format.ts       # Formatting utilities
│   │   └── views/              # Page components
│   │       ├── Dashboard.vue   # Home dashboard
│   │       ├── Login.vue       # Login page
│   │       ├── transaction/    # Transaction management
│   │       ├── category/       # Expense categories
│   │       ├── statistics/     # Statistical reports
│   │       └── settings/       # System settings
│   └── dist/                   # Build output
├── docs/                       # Project documentation
│   ├── Requirements_Specification.md
│   ├── API_Documentation.md
│   ├── User_Manual.md
│   ├── Deployment_Guide.md
│   └── Resubmission_Feature.md
└── README.md                   # Project README
```

## Modules

### Dashboard
- Stat cards for total income, total expense, net balance, and record count
- Recent transaction list
- Project-based pie chart

### Transaction Management
- **Income Registration** - Finance staff register income records
- **Expense Registration** - Finance staff register expense records
- **Transaction Approval** - Administrators approve transactions
- **Transaction List** - View and filter transactions

### Expense Categories
- CRUD operations on categories
- Supports fuzzy search

### Statistical Reports
- Multi-dimensional analysis (project/personnel/category)
- Visualized charts

### System Settings
- **User Management** - User CRUD and permission assignment (Admin)
- **Department Management** - Department information management (Admin)
- **Project Management** - Project information management (Admin)
- **Operation Logs** - User operation records (Admin/Finance)
- **Account Management** - Personal account maintenance
- **Version History** - System version update information

## User Roles

| Role | Permissions |
|------|-------------|
| **Admin (ADMIN)** | Full permissions, including user management, department management, project management, and transaction approval |
| **Finance (FINANCE)** | Transaction registration, expense categories, statistical reports, and operation log viewing |
| **Employee (EMPLOYEE)** | View personal transaction records, account management |

## Quick Start

### Prerequisites
- Go 1.18+
- Node.js 18+
- MySQL 8.0+

### Database Initialization
```bash
# Create database
mysql -u root -p < backend/sql/01_create_database.sql
# Create table structure
mysql -u root -p < backend/sql/02_create_tables.sql
# Initialize base data
mysql -u root -p < backend/sql/03_init_data.sql
```

### Install Dependencies

**Backend:**
```bash
cd backend
go mod download
```

**Frontend:**
```bash
cd frontend
npm install
```

### Configuration

Edit `backend/config/config.yaml`:
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

### Development Mode

**Start backend:**
```bash
cd backend
go run main.go
```

**Start frontend:**
```bash
cd frontend
npm run dev
```

### Production Build

**Build backend (Linux x86_64):**
```bash
cd backend
GOOS=linux GOARCH=amd64 go build -o account-management main.go
```

**Build frontend:**
```bash
cd frontend
npx vite build
```

## Default Account

| Username | Password | Role |
|----------|----------|------|
| admin | admin123 | Administrator |

## Version History

### v1.0.2 (2026-03-18)
- Added payment method field to support multiple payment options
- Expense registration page allows adding new payment methods directly
- Implemented resubmission of rejected records
- Optimized amount display with thousands separators
- Unified formatting utility functions
- Improved table column width display
- Refined version info UI, opened version history viewing permission

### v1.0.1 (2026-03-01)
- Optimized login page layout with a two-column design
- Fixed record count display issue on the dashboard
- Updated deployment documentation

### v1.0.0 (2026-03-01)
- Initial stable release
- Complete accounting management system features

## Documentation

- [Requirements Specification](./docs/Requirements_Specification.md) (Chinese)
- [API Documentation](./docs/API_Documentation.md) (Chinese)
- [Deployment Guide](./docs/Deployment_Guide.md) (Chinese)
- [User Manual](./docs/User_Manual.md) (Chinese)
- [Resubmission Feature Notes](./docs/Resubmission_Feature.md) (Chinese)

## License

This project is licensed under the [GPL-3.0 License](./LICENSE).
