# 账务管理系统接口文档

| 版本 | 日期 | 修订人 | 修订说明 |
|------|------|--------|----------|
| v1.0 | 2026-02-15 | V5General | 初始版本 |
| v1.1 | 2026-03-03 | V5General | 新增部门管理、项目管理、账号管理、收支审核接口，更新费用分类和用户模型 |

---

## 目录

- [1. 接口规范](#1-接口规范)
- [2. 认证接口](#2-认证接口)
- [3. 用户管理接口](#3-用户管理接口)
- [4. 部门管理接口](#4-部门管理接口)
- [5. 项目管理接口](#5-项目管理接口)
- [6. 费用分类接口](#6-费用分类接口)
- [7. 收支管理接口](#7-收支管理接口)
- [8. 凭证附件接口](#8-凭证附件接口)
- [9. 统计报表接口](#9-统计报表接口)
- [10. 错误码说明](#10-错误码说明)

---

## 1. 接口规范

### 1.1 基础信息

| 项目 | 说明 |
|------|------|
| 基础URL | `http://localhost:8080/api/v1` |
| 数据格式 | JSON |
| 字符编码 | UTF-8 |
| 认证方式 | JWT Bearer Token |

### 1.2 请求头

```
Content-Type: application/json
Authorization: Bearer {token}
```

### 1.3 统一响应格式

**成功响应：**
```json
{
    "code": 0,
    "message": "success",
    "data": {},
    "timestamp": 1706918400000
}
```

**失败响应：**
```json
{
    "code": 1001,
    "message": "用户名或密码错误",
    "data": null,
    "timestamp": 1706918400000
}
```

### 1.4 分页参数

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**分页响应：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "list": [],
        "total": 100,
        "page": 1,
        "page_size": 20
    }
}
```

---

## 2. 认证接口

### 2.1 用户登录

**接口地址：** `POST /auth/login`

**请求参数：**
```json
{
    "username": "admin",
    "password": "admin123"
}
```

**响应数据：**
```json
{
    "code": 0,
    "message": "登录成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "user_id": "admin001",
            "username": "admin",
            "real_name": "系统管理员",
            "role": "ADMIN",
            "department_id": "dept001",
            "department_name": "管理部",
            "status": 1,
            "create_time": "2024-01-01 00:00:00"
        }
    }
}
```

### 2.2 用户注销

**接口地址：** `POST /auth/logout`

**需要认证：** 是

**响应数据：**
```json
{
    "code": 0,
    "message": "注销成功",
    "data": null
}
```

### 2.3 获取当前用户信息

**接口地址：** `GET /auth/me`

**需要认证：** 是

**响应数据：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "user_id": "admin001",
        "username": "admin",
        "real_name": "系统管理员",
        "role": "ADMIN",
        "department_id": "dept001",
        "department_name": "管理部",
        "status": 1,
        "create_time": "2024-01-01 00:00:00"
    }
}
```

---

## 3. 用户管理接口

> 所有接口需要管理员权限

### 3.1 创建用户

**接口地址：** `POST /users`

**需要认证：** 是
**需要权限：** 管理员

**请求参数：**
```json
{
    "username": "newuser",
    "password": "password123",
    "real_name": "张三",
    "role": "EMPLOYEE",
    "department_id": "dept003"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名（唯一，只能包含数字、字母和下划线） |
| password | string | 是 | 密码（至少8位） |
| real_name | string | 是 | 真实姓名 |
| role | string | 是 | 角色：ADMIN/EMPLOYEE/FINANCE |
| department_id | string | 否 | 所属部门ID |

**响应数据：**
```json
{
    "code": 0,
    "message": "创建成功",
    "data": {
        "user_id": "user001"
    }
}
```

### 3.2 获取用户列表

**接口地址：** `GET /users`

**需要认证：** 是
**需要权限：** 管理员

**查询参数：**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |
| role | string | 否 | - | 角色筛选 |
| department_id | string | 否 | - | 部门ID筛选 |
| keyword | string | 否 | - | 关键字搜索（用户名或真实姓名） |

**响应数据：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "list": [
            {
                "user_id": "user001",
                "username": "caiwu",
                "real_name": "财务人员",
                "role": "FINANCE",
                "department_id": "dept002",
                "department_name": "财务部",
                "status": 1,
                "create_time": "2024-01-01 00:00:00"
            }
        ],
        "total": 10,
        "page": 1,
        "page_size": 20
    }
}
```

### 3.3 更新用户

**接口地址：** `PUT /users/{id}`

**需要认证：** 是
**需要权限：** 管理员

**路径参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 用户ID |

**请求参数：**
```json
{
    "username": "newusername",
    "real_name": "新名字",
    "role": "ADMIN",
    "department_id": "dept001",
    "status": 1,
    "password": "newpassword123"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 否 | 用户名（唯一，只能包含数字、字母和下划线） |
| real_name | string | 否 | 真实姓名 |
| role | string | 否 | 角色 |
| department_id | string | 否 | 部门ID |
| status | int | 否 | 状态：1-正常，0-禁用 |
| password | string | 否 | 密码（至少8位） |

**响应数据：**
```json
{
    "code": 0,
    "message": "更新成功",
    "data": null
}
```

### 3.4 删除用户

**接口地址：** `DELETE /users/{id}`

**需要认证：** 是
**需要权限：** 管理员

**路径参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 用户ID |

**响应数据：**
```json
{
    "code": 0,
    "message": "删除成功",
    "data": null
}
```

### 3.5 重置密码

**接口地址：** `POST /users/{id}/reset-password`

**需要认证：** 是
**需要权限：** 管理员

**路径参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 用户ID |

**请求参数：**
```json
{
    "password": "newpassword123"
}
```

**响应数据：**
```json
{
    "code": 0,
    "message": "重置成功",
    "data": null
}
```

### 3.6 更新自己的账号

**接口地址：** `PUT /account/me`

**需要认证：** 是

**说明：** 用户更新自己的账号信息（所有角色可用）

**请求参数：**
```json
{
    "username": "newusername",
    "real_name": "新名字",
    "password": "newpassword123"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 否 | 用户名（唯一，只能包含数字、字母和下划线） |
| real_name | string | 是 | 真实姓名 |
| password | string | 否 | 密码（至少8位） |

**响应数据：**
```json
{
    "code": 0,
    "message": "更新成功",
    "data": null
}
```

---

## 4. 部门管理接口

> 所有接口需要管理员权限

### 4.1 创建部门

**接口地址：** `POST /departments`

**需要认证：** 是
**需要权限：** 管理员

**请求参数：**
```json
{
    "name": "技术部",
    "description": "技术研发部门",
    "sort_order": 3
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 部门名称（唯一） |
| description | string | 否 | 部门描述 |
| sort_order | int | 否 | 排序顺序（数字越小越靠前） |

**响应数据：**
```json
{
    "code": 0,
    "message": "创建成功",
    "data": {
        "department_id": "dept003",
        "name": "技术部",
        "description": "技术研发部门",
        "sort_order": 3,
        "create_time": "2024-01-01 00:00:00"
    }
}
```

### 4.2 获取部门列表

**接口地址：** `GET /departments`

**需要认证：** 是
**需要权限：** 管理员

**响应数据：**
```json
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "department_id": "dept001",
            "name": "管理部",
            "description": "公司管理职能部门",
            "sort_order": 1,
            "create_time": "2024-01-01 00:00:00"
        }
    ]
}
```

### 4.3 获取部门详情

**接口地址：** `GET /departments/{id}`

**需要认证：** 是
**需要权限：** 管理员

**路径参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 部门ID |

**响应数据：** 同4.2中的单个部门格式

### 4.4 更新部门

**接口地址：** `PUT /departments/{id}`

**需要认证：** 是
**需要权限：** 管理员

**请求参数：**
```json
{
    "name": "技术研发部",
    "description": "产品技术研发",
    "sort_order": 2
}
```

**响应数据：**
```json
{
    "code": 0,
    "message": "更新成功",
    "data": null
}
```

### 4.5 删除部门

**接口地址：** `DELETE /departments/{id}`

**需要认证：** 是
**需要权限：** 管理员

**注意：**
- 已关联用户的部门无法删除
- 已关联项目的部门无法删除

**响应数据：**
```json
{
    "code": 0,
    "message": "删除成功",
    "data": null
}
```

### 4.6 获取部门下的用户列表

**接口地址：** `GET /departments/{id}/users`

**需要认证：** 是
**需要权限：** 管理员

**查询参数：**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 10 | 每页数量 |

**响应数据：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "list": [
            {
                "user_id": "user001",
                "username": "zhangsan",
                "real_name": "张三",
                "role": "EMPLOYEE",
                "status": 1
            }
        ],
        "total": 5
    }
}
```

---

## 5. 项目管理接口

> 查询接口所有认证用户可访问，增删改需要财务或管理员权限

### 5.1 创建项目

**接口地址：** `POST /projects`

**需要认证：** 是
**需要权限：** 财务人员/管理员

**请求参数：**
```json
{
    "name": "项目A",
    "description": "项目A的描述",
    "department_id": "dept003"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 项目名称（唯一） |
| description | string | 否 | 项目描述 |
| department_id | string | 否 | 关联部门ID |

**响应数据：**
```json
{
    "code": 0,
    "message": "创建成功",
    "data": {
        "project_id": "project001",
        "name": "项目A",
        "description": "项目A的描述",
        "department_id": "dept003",
        "status": 1,
        "create_time": "2024-01-01 00:00:00"
    }
}
```

### 5.2 获取项目列表

**接口地址：** `GET /projects`

**需要认证：** 是

**查询参数：**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 10 | 每页数量 |
| department_id | string | 否 | - | 部门ID筛选 |

**响应数据：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "list": [
            {
                "project_id": "project001",
                "name": "项目A",
                "description": "项目A的描述",
                "department_id": "dept003",
                "status": 1,
                "create_time": "2024-01-01 00:00:00"
            }
        ],
        "total": 10
    }
}
```

### 5.3 获取项目详情

**接口地址：** `GET /projects/{id}`

**需要认证：** 是

**响应数据：** 同5.1中的单个项目格式

### 5.4 更新项目

**接口地址：** `PUT /projects/{id}`

**需要认证：** 是
**需要权限：** 财务人员/管理员

**请求参数：**
```json
{
    "name": "项目A（更新）",
    "description": "更新后的描述",
    "department_id": "dept002",
    "status": 0
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 项目名称 |
| description | string | 否 | 项目描述 |
| department_id | string | 否 | 关联部门ID |
| status | int | 否 | 状态：1-进行中，0-已结束 |

**响应数据：**
```json
{
    "code": 0,
    "message": "更新成功",
    "data": null
}
```

### 5.5 删除项目

**接口地址：** `DELETE /projects/{id}`

**需要认证：** 是
**需要权限：** 财务人员/管理员

**注意：** 已关联收支记录的项目无法删除

**响应数据：**
```json
{
    "code": 0,
    "message": "删除成功",
    "data": null
}
```

---

## 6. 费用分类接口

> 需要财务人员或管理员权限

### 6.1 创建费用分类

**接口地址：** `POST /categories`

**需要认证：** 是
**需要权限：** 财务人员/管理员

**请求参数：**
```json
{
    "name": "差旅费",
    "type": "EXPENSE",
    "description": "出差交通、住宿费用",
    "sort_order": 4
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 分类名称（唯一） |
| type | string | 是 | 类型：INCOME（收入）/EXPENSE（支出） |
| description | string | 否 | 分类描述 |
| sort_order | int | 否 | 排序顺序 |

**响应数据：**
```json
{
    "code": 0,
    "message": "创建成功",
    "data": {
        "category_id": "category007"
    }
}
```

### 6.2 获取分类列表

**接口地址：** `GET /categories`

**需要认证：** 是
**需要权限：** 财务人员/管理员

**查询参数：**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |
| keyword | string | 否 | - | 关键字搜索 |
| type | string | 否 | - | 类型筛选：INCOME/EXPENSE |

**响应数据：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "list": [
            {
                "category_id": "category004",
                "name": "工资",
                "type": "EXPENSE",
                "description": "员工工资发放",
                "sort_order": 1,
                "is_deleted": 0,
                "create_time": "2024-01-01 00:00:00"
            }
        ],
        "total": 7,
        "page": 1,
        "page_size": 20
    }
}
```

### 6.3 更新分类

**接口地址：** `PUT /categories/{id}`

**需要认证：** 是
**需要权限：** 财务人员/管理员

**请求参数：**
```json
{
    "name": "差旅费",
    "type": "EXPENSE",
    "description": "出差交通、住宿费用",
    "sort_order": 4
}
```

**响应数据：**
```json
{
    "code": 0,
    "message": "更新成功",
    "data": null
}
```

### 6.4 删除分类

**接口地址：** `DELETE /categories/{id}`

**需要认证：** 是
**需要权限：** 财务人员/管理员

**响应数据：**
```json
{
    "code": 0,
    "message": "删除成功",
    "data": null
}
```

---

## 7. 收支管理接口

### 7.1 创建收支记录

**接口地址：** `POST /transactions`

**需要认证：** 是
**需要权限：** 财务人员/管理员

**请求参数：**
```json
{
    "amount": 50000.00,
    "category_id": "category004",
    "project_id": "project001",
    "person_id": "user002",
    "transaction_time": "2024-01-15 10:00:00",
    "remark": "一月份工资",
    "attachment_ids": ["att001", "att002"]
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| amount | float | 是 | 金额（正数=收入，负数=支出） |
| category_id | string | 否 | 费用分类ID |
| project_id | string | 是 | 项目ID |
| person_id | string | 条件必填 | 关联人员ID（支出时必填） |
| transaction_time | string | 是 | 交易时间（格式：yyyy-MM-dd HH:mm:ss） |
| remark | string | 否 | 备注 |
| attachment_ids | array | 是 | 凭证附件ID列表（至少1个） |

**响应数据：**
```json
{
    "code": 0,
    "message": "创建成功",
    "data": {
        "record_id": "record001"
    }
}
```

### 7.2 获取收支记录列表

**接口地址：** `GET /transactions`

**需要认证：** 是

**查询参数：**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |
| start_time | string | 否 | - | 开始时间（yyyy-MM-dd） |
| end_time | string | 否 | - | 结束时间（yyyy-MM-dd） |
| category_id | string | 否 | - | 费用分类ID |
| project_id | string | 否 | - | 项目ID |
| person_id | string | 否 | - | 关联人员ID |
| status | int | 否 | - | 状态筛选：0-待审核，1-已审核，2-已驳回 |
| type | string | 否 | all | 类型：income/expense/all |

**权限说明：**
- 管理员/财务人员：可查看所有记录
- 员工：只能查看本人关联的记录

**响应数据：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "list": [
            {
                "record_id": "record001",
                "amount": 50000.00,
                "category_id": "category004",
                "project_id": "project001",
                "person_id": "user002",
                "transaction_time": "2024-01-15 10:00:00",
                "remark": "一月份工资",
                "status": 1,
                "creator_id": "user001",
                "create_time": "2024-01-15 10:00:00",
                "update_time": "2024-01-15 10:00:00",
                "category": {
                    "category_id": "category004",
                    "name": "工资",
                    "type": "EXPENSE"
                },
                "project": {
                    "project_id": "project001",
                    "name": "项目A"
                },
                "person": {
                    "user_id": "user002",
                    "username": "zhangsan",
                    "real_name": "张三"
                },
                "creator": {
                    "user_id": "user001",
                    "username": "caiwu",
                    "real_name": "财务"
                },
                "attachments": [
                    {
                        "attachment_id": "att001",
                        "file_name": "工资单.pdf",
                        "file_size": 102400,
                        "file_type": "pdf"
                    }
                ]
            }
        ],
        "total": 100,
        "page": 1,
        "page_size": 20
    }
}
```

### 7.3 获取收支记录详情

**接口地址：** `GET /transactions/{id}`

**需要认证：** 是

**路径参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 记录ID |

**响应数据：** 同7.2中的单条记录格式

### 7.4 更新收支记录

**接口地址：** `PUT /transactions/{id}`

**需要认证：** 是
**需要权限：** 财务人员/管理员

**请求参数：**
```json
{
    "amount": 50000.00,
    "category_id": "category004",
    "project_id": "project001",
    "person_id": "user002",
    "transaction_time": "2024-01-15 10:00:00",
    "remark": "一月份工资（更新）",
    "status": 0
}
```

**响应数据：**
```json
{
    "code": 0,
    "message": "更新成功",
    "data": null
}
```

### 7.5 删除收支记录

**接口地址：** `DELETE /transactions/{id}`

**需要认证：** 是
**需要权限：** 财务人员/管理员

**响应数据：**
```json
{
    "code": 0,
    "message": "删除成功",
    "data": null
}
```

### 7.6 审核通过

**接口地址：** `PUT /transactions/{id}/approve`

**需要认证：** 是
**需要权限：** 管理员

**说明：** 审核通过收支记录，状态变更为1（已审核）

**响应数据：**
```json
{
    "code": 0,
    "message": "审核通过",
    "data": null
}
```

### 7.7 审核拒绝

**接口地址：** `PUT /transactions/{id}/reject`

**需要认证：** 是
**需要权限：** 管理员

**说明：** 审核拒绝收支记录，状态变更为2（已驳回）

**响应数据：**
```json
{
    "code": 0,
    "message": "审核拒绝",
    "data": null
}
```

### 7.8 获取统计数据

**接口地址：** `GET /transactions/statistics`

**需要认证：** 是
**需要权限：** 财务人员/管理员

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| dimension | string | 是 | 统计维度：project/person/category |
| cycle | string | 是 | 统计周期：month/quarter/year |
| start_time | string | 是 | 开始时间（yyyy-MM-dd） |
| end_time | string | 是 | 结束时间（yyyy-MM-dd） |

**响应数据：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "dimension": "project",
        "cycle": "month",
        "start_time": "2024-01-01",
        "end_time": "2024-01-31",
        "summary": {
            "total_income": 150000.00,
            "total_expense": -80000.00,
            "net_amount": 70000.00,
            "record_count": 50
        },
        "details": [
            {
                "key": "项目A",
                "income": 100000.00,
                "expense": -50000.00,
                "net_amount": 50000.00,
                "record_count": 30,
                "percentage": 71.43
            },
            {
                "key": "项目B",
                "income": 50000.00,
                "expense": -30000.00,
                "net_amount": 20000.00,
                "record_count": 20,
                "percentage": 28.57
            }
        ]
    }
}
```

---

## 8. 凭证附件接口

### 8.1 上传附件

**接口地址：** `POST /attachments`

**需要认证：** 是

**请求类型：** `multipart/form-data`

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 附件文件 |

**文件限制：**
- 支持格式：jpg、png、pdf
- 单文件大小：最大100MB
- 数量限制：单次上传1个文件

**响应数据：**
```json
{
    "code": 0,
    "message": "上传成功",
    "data": {
        "attachment_id": "att001",
        "file_name": "发票.jpg",
        "file_size": 102400,
        "file_type": "image"
    }
}
```

### 8.2 获取附件列表

**接口地址：** `GET /attachments`

**需要认证：** 是

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| record_id | string | 是 | 收支记录ID |

**响应数据：**
```json
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "attachment_id": "att001",
            "record_id": "record001",
            "file_name": "发票.jpg",
            "file_size": 102400,
            "file_type": "image",
            "upload_time": "2024-01-15 10:00:00"
        }
    ]
}
```

### 8.3 下载附件

**接口地址：** `GET /attachments/{id}/download`

**需要认证：** 是

**路径参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 附件ID |

**响应：** 重定向到文件URL或直接返回文件流

### 8.4 删除附件

**接口地址：** `DELETE /attachments/{id}`

**需要认证：** 是

**权限说明：**
- 上传者本人可删除
- 财务人员/管理员可删除
- 已关联到收支记录的附件不可删除

**路径参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 附件ID |

**响应数据：**
```json
{
    "code": 0,
    "message": "删除成功",
    "data": null
}
```

---

## 9. 统计报表接口

### 9.1 获取统计数据

详见 [7.8 获取统计数据](#78-获取统计数据)

### 9.2 导出报表

**接口地址：** `GET /reports/{id}/export`

**需要认证：** 是
**需要权限：** 财务人员/管理员

**路径参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 报表ID |

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| format | string | 是 | 导出格式：excel/pdf |

**响应：** 文件下载流

---

## 10. 错误码说明

| 错误码 | 说明 | HTTP状态码 |
|--------|------|-----------|
| 0 | 成功 | 200 |
| 1001 | 用户名或密码错误 | 200 |
| 1002 | 令牌无效或过期 | 200 |
| 1003 | 权限不足 | 200 |
| 2001 | 参数错误 | 200 |
| 2002 | 资源不存在 | 200 |
| 2003 | 资源已存在 | 200 |
| 3001 | 文件上传失败 | 200 |
| 3002 | 文件格式不支持 | 200 |
| 3003 | 文件大小超限 | 200 |
| 5000 | 服务器内部错误 | 500 |

---

## 附录

### A. 用户角色说明

| 角色 | 代码 | 说明 |
|------|------|------|
| 管理员 | ADMIN | 拥有所有权限 |
| 财务人员 | FINANCE | 可管理收支、分类、报表 |
| 员工 | EMPLOYEE | 仅可查看本人相关记录 |

### B. 交易状态说明

| 状态 | 值 | 说明 |
|------|-----|------|
| 待审核 | 0 | 已创建但未审核 |
| 已审核 | 1 | 审核通过，正常生效 |
| 已驳回 | 2 | 审核拒绝，不生效 |

### C. 操作类型说明

| 类型 | 说明 |
|------|------|
| LOGIN | 用户登录 |
| CREATE | 创建操作 |
| UPDATE | 更新操作 |
| DELETE | 删除操作 |
| APPROVE | 审核操作 |
| QUERY | 查询操作 |
| EXPORT | 导出操作 |

### D. 费用分类预设列表

**收入分类：**

| 分类ID | 名称 | 类型 | 描述 |
|--------|------|------|------|
| category001 | 服务收入 | INCOME | 提供服务获得的收入 |
| category002 | 销售收入 | INCOME | 产品销售收入 |
| category003 | 其他收入 | INCOME | 其他收入来源 |

**支出分类：**

| 分类ID | 名称 | 类型 | 描述 |
|--------|------|------|------|
| category004 | 工资 | EXPENSE | 员工工资发放 |
| category005 | 设备采购 | EXPENSE | 办公设备、生产设备采购 |
| category006 | 服务购买 | EXPENSE | 外部服务采购 |
| category007 | 差旅费 | EXPENSE | 出差交通、住宿费用 |
| category008 | 业务招待费 | EXPENSE | 客户招待费用 |
| category009 | 办公费用 | EXPENSE | 日常办公用品采购 |
| category010 | 其他 | EXPENSE | 其他费用 |

### E. 项目状态说明

| 状态 | 值 | 说明 |
|------|-----|------|
| 进行中 | 1 | 项目正在进行 |
| 已结束 | 0 | 项目已完成 |
