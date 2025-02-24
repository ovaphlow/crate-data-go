# Crate Data API

A flexible RESTful data API service written in Go that supports multiple database backends (MySQL, PostgreSQL, and SQLite).

一个使用 Go 语言编写的灵活的 RESTful 数据 API 服务，支持多种数据库后端（MySQL、PostgreSQL 和 SQLite）。

> Documentation written by GitHub Copilot in agent mode.
> 
> 文档由 GitHub Copilot agent 模式自动撰写。

> **⚠️ Security Notice | 安全提示**
>
> Due to potential data leakage risks, this API service is recommended for:
> - Server-Side Rendering (SSR) frontend applications
> - Backend service-to-service communication
>
> This ensures sensitive data can be properly processed on the server side before being transmitted to clients.
>
> 由于存在数据泄露风险，此 API 服务建议用于：
> - 服务器端渲染（SSR）前端应用
> - 后端服务之间的通信
>
> 这样可以确保敏感数据在传输到客户端之前在服务器端得到适当处理。

## Features | 特性

- Multiple database support (MySQL, PostgreSQL, SQLite)
  - 支持多种数据库（MySQL、PostgreSQL、SQLite）
- RESTful API endpoints
  - RESTful API 接口
- Middleware support for API versioning, CORS, and security headers
  - 支持 API 版本控制、CORS 和安全头的中间件
- Structured logging
  - 结构化日志记录
- Environment-based configuration
  - 基于环境变量的配置
- RFC9457-compliant HTTP responses
  - 符合 RFC9457 标准的 HTTP 响应

## Setup | 设置

1. Clone the repository | 克隆仓库
2. Create a `.env` file with the following configuration options | 创建包含以下配置选项的 `.env` 文件：

```env
PORT=8421  # Default port if not specified | 默认端口（如未指定）

# PostgreSQL Configuration | PostgreSQL 配置
POSTGRES_ENABLED=true  # or false | 启用或禁用
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DATABASE=your_database

# MySQL Configuration | MySQL 配置
MYSQL_ENABLED=true  # or false | 启用或禁用
MYSQL_USER=your_user
MYSQL_PASSWORD=your_password
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DATABASE=your_database

# SQLite Configuration | SQLite 配置
SQLITE_ENABLED=true  # or false | 启用或禁用
SQLITE_DATABASE=./data.db  # SQLite database file path | SQLite 数据库文件路径
```

3. Build and run the application | 构建并运行应用：
```bash
go build -o crate-api-data cmd/main.go
./crate-api-data
```

## Database Schema | 数据库表结构

Each table in the database must have the following required fields:

每个数据库表必须包含以下必要字段：

### Required Fields | 必要字段

| Field Name | Description | PostgreSQL | MySQL | SQLite | 说明 |
|------------|-------------|------------|-------|---------|------|
| id | Unique identifier | VARCHAR(27) | VARCHAR(27) | TEXT | 唯一标识符 |
| event_time | Event timestamp with timezone | TIMESTAMPTZ | TIMESTAMP | TEXT | 带时区的事件时间戳 |
| data_state | Record state metadata | JSONB | JSON | TEXT | 记录状态元数据 |

### Data State Structure | 数据状态结构

The `data_state` field is specifically designed to store record metadata and state information, NOT business data. It typically includes:

`data_state` 字段专门用于存储记录的元数据和状态信息，而不是业务数据。通常包含：

```json
{
    "created_at": "2024-03-20T10:00:00Z",    // Creation timestamp | 创建时间
    "updated_at": "2024-03-21T15:30:00Z",    // Last update timestamp | 最后更新时间
    "removed_at": "2024-03-22T08:00:00Z",  // Deprecation timestamp (if applicable) | 废弃时间（如果适用）
    "status": "active",                       // Record status (active/deprecated) | 记录状态（活动/废弃）
}
```

Important Notes | 重要说明：
- Business data should be stored in dedicated table columns, not in `data_state`
  - 业务数据应存储在专用的表字段中，而不是 `data_state` 中
- `data_state` is for tracking record lifecycle and metadata only
  - `data_state` 仅用于跟踪记录生命周期和元数据
- The structure is automatically managed by the API
  - 该结构由 API 自动管理

### Example Table Creation | 建表示例

#### PostgreSQL
```sql
CREATE TABLE your_table (
    id VARCHAR(27) PRIMARY KEY,
    event_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    data_state JSONB NOT NULL DEFAULT '{
        "created_at": CURRENT_TIMESTAMP,
        "status": "active",
    }'::jsonb,
    -- Add your business data columns here | 在此添加业务数据字段
    name VARCHAR(100),
    email VARCHAR(255),
    -- etc...
);
```

#### MySQL
```sql
CREATE TABLE your_table (
    id VARCHAR(27) PRIMARY KEY,
    event_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    data_state JSON NOT NULL DEFAULT (JSON_OBJECT(
        'created_at', CURRENT_TIMESTAMP,
        'status', 'active',
    )),
    -- Add your business data columns here | 在此添加业务数据字段
    name VARCHAR(100),
    email VARCHAR(255),
    -- etc...
);
```

#### SQLite
```sql
CREATE TABLE your_table (
    id TEXT PRIMARY KEY,
    event_time TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    data_state TEXT NOT NULL DEFAULT json_object(
        'created_at', datetime('now'),
        'status', 'active',
    ),
    -- Add your business data columns here | 在此添加业务数据字段
    name TEXT,
    email TEXT,
    -- etc...
);
```

Note | 注意：
- The `id` field is automatically generated using KSUID (K-Sortable Unique Identifier)
  - `id` 字段使用 KSUID（K-Sortable Unique Identifier）自动生成
- `event_time` defaults to the current timestamp when not provided
  - `event_time` 在未提供时默认使用当前时间戳
- `data_state` stores the JSON data of your record
  - `data_state` 存储记录的 JSON 数据

## API Documentation | API 文档

Base URL | 基础 URL: `/crate-api-data`

### Endpoints | 接口端点

对于每种数据库类型（mysql/postgres/sqlite），都提供以下端点：

#### Create Record | 创建记录
- **POST** `/{db_type}/{table}`
- **Body | 请求体**: JSON object with record data | JSON 格式的记录数据
- **Response | 响应**: 201 Created on success | 成功时返回 201

#### Retrieve Records | 获取记录列表
- **GET** `/{db_type}/{table}`
- **Query Parameters | 查询参数**:
  - `l`: SQL suffix for query customization (e.g., ORDER BY, LIMIT, GROUP BY) | SQL 查询后缀（如排序、限制、分组等）
    - Examples | 示例:
      - `l=ORDER BY event_time DESC` - Sort by time descending | 按时间降序排序
      - `l=LIMIT 10` - Limit results to 10 records | 限制返回10条记录
      - `l=GROUP BY city` - Group results by city | 按城市分组
  - `f`: Filter criteria | 过滤条件
  - `c`: Column selection or SQL functions | 列选择或 SQL 函数
    - Examples | 示例:
      - `c=name,age` - Select specific columns | 选择特定列
      - `c=COUNT(*) as qty` - Count records | 统计记录数
      - `c=SUM(amount) as total` - Sum values | 求和

### Query Parameters Format | 查询参数格式

#### Filter Parameter (`f`) Format | 过滤参数格式

The filter parameter uses a special format: `operator,count,field1,value1,...`

过滤参数使用特殊格式：`操作符,参数数量,字段1,值1,...`

Available operators | 可用操作符：
- `eq` / `equal`: Equal comparison | 等于比较
- `ne` / `not-equal`: Not equal comparison | 不等于比较
- `in`: In array of values | 在值数组中
- `lk` / `like`: Like comparison | 模糊匹配
- `ge` / `greater-equal`: Greater than or equal | 大于等于
- `le` / `less-equal`: Less than or equal | 小于等于
- `gt` / `greater`: Greater than | 大于
- `lt` / `less`: Less than | 小于
- `act` / `array-contain`: JSON array contains | JSON数组包含
- `oct` / `object-contain`: JSON object contains | JSON对象包含

Examples | 示例：

```bash
# Simple equal comparison | 简单等于比较
# Format: eq,<count>,<field1>,<value1>,...
f=eq,2,id,1123
f=eq,4,name,john,age,25

# Multiple conditions with JSON object contains | 带JSON对象包含的多个条件
# Format: oct,<count>,<json_field>,<key>,<value>,...
f=oct,3,data_state,status,active

# Combining IN condition with other filters | 组合IN条件和其他过滤器
# Format: in,<count>,<value1>,<value2>,...
f=in,4,id,1123,2234,3345

# Complex example combining multiple conditions | 组合多个条件的复杂示例
f=eq,2,status,active,oct,3,data_state,type,user,in,3,region,east,west

# Like comparison for search | 用于搜索的模糊匹配
f=lk,2,name,john%

# Numeric comparisons | 数值比较
f=gt,2,age,18
f=le,2,price,100.50
```

Filter Format Rules | 过滤格式规则：
1. Each condition starts with an operator followed by count | 每个条件以操作符开始，后跟参数数量
2. Count indicates how many parameters follow | 参数数量表示后面有多少个参数
3. Multiple conditions can be combined with commas | 多个条件可以用逗号组合
4. For `eq`, `ne`, `like`, `ge`, `le`, `gt`, `lt`: count must be even | 对于`eq`,`ne`,`like`,`ge`,`le`,`gt`,`lt`：参数数量必须是偶数
5. For `in`: count indicates the number of values | 对于`in`：参数数量表示值的个数
6. For `oct`: count must be divisible by 3 | 对于`oct`：参数数量必须是3的倍数

### Advanced Query Examples | 高级查询示例

```bash
# Get record count by city | 按城市统计记录数
curl "http://localhost:8421/crate-api-data/mysql/users?c=city,COUNT(*) as qty&l=GROUP BY city"

# Get latest 10 records sorted by time | 获取最新10条记录
curl "http://localhost:8421/crate-api-data/mysql/users?l=ORDER BY event_time DESC LIMIT 10"

# Get total amount by category | 按类别统计总金额
curl "http://localhost:8421/crate-api-data/mysql/orders?c=category,SUM(amount) as total&l=GROUP BY category"

# Combine filters with aggregation | 组合过滤条件和聚合
curl "http://localhost:8421/crate-api-data/mysql/orders?f=status:eq:completed&c=category,COUNT(*) as qty&l=GROUP BY category"

# Get active records only | 只获取活动状态的记录
curl "http://localhost:8421/crate-api-data/mysql/users?f=data_state->>'status':eq:active"

# Get records modified in last 24 hours | 获取最近24小时内修改的记录
curl "http://localhost:8421/crate-api-data/mysql/users?f=data_state->>'updated_at':gt:2024-03-20T00:00:00Z"

# Get version history count by status | 按状态统计版本历史数量
curl "http://localhost:8421/crate-api-data/mysql/users?c=data_state->>'status' as status,COUNT(*) as count&l=GROUP BY data_state->>'status'"

# Get deprecated records with their deprecation time | 获取已废弃记录及其废弃时间
curl "http://localhost:8421/crate-api-data/mysql/users?c=id,name,data_state->>'removed_at' as deprecated_time&f=data_state->>'status':eq:deprecated"

# Get records by version number | 按版本号获取记录
curl "http://localhost:8421/crate-api-data/mysql/users?f=data_state->>'version':eq:2"

# Filter by multiple equal conditions | 多个等于条件过滤
curl "http://localhost:8421/crate-api-data/mysql/users?f=eq,4,status,active,role,admin"

# Filter by JSON object contains with IN clause | JSON对象包含条件与IN子句组合
curl "http://localhost:8421/crate-api-data/mysql/users?f=oct,3,data_state,status,active,in,3,region,east,west"

# Complex filter with multiple conditions | 多条件复杂过滤
curl "http://localhost:8421/crate-api-data/mysql/orders?f=eq,2,status,pending,gt,2,amount,1000"

# Search with LIKE operator | 使用LIKE操作符搜索
curl "http://localhost:8421/crate-api-data/mysql/products?f=lk,2,name,phone%"
```

Note on JSON field access | JSON 字段访问说明：
- PostgreSQL uses `->>` operator for JSON text access | PostgreSQL 使用 ->> 运算符访问 JSON 文本
- MySQL uses ->>, -> or JSON_EXTRACT() | MySQL 使用 ->>, -> 或 JSON_EXTRACT()
- SQLite uses json_extract() | SQLite 使用 json_extract()

#### Retrieve Single Record | 获取单条记录
- **GET** `/{db_type}/{table}/{id}`
- **Response | 响应**: Single record object | 单条记录对象

#### Update Record | 更新记录
- **PUT** `/{db_type}/{table}/{id}`
- **Query Parameters | 查询参数**:
  - `d`: Set to "true" or "1" to mark as deprecated | 设置为 "true" 或 "1" 表示标记为废弃
- **Body | 请求体**: JSON object with updated data | JSON 格式的更新数据
- **Response | 响应**: 200 OK on success | 成功时返回 200

#### Delete Record | 删除记录
- **DELETE** `/{db_type}/{table}/{id}`
- **Response | 响应**: 200 OK on success | 成功时返回 200

## Error Handling | 错误处理

The API follows RFC9457 for HTTP response formatting. All error responses include:

API 遵循 RFC9457 标准进行 HTTP 响应格式化。所有错误响应都包含：

- A descriptive message | 描述性消息
- HTTP status code | HTTP 状态码
- Request details | 请求详情

## Security | 安全性

The API includes several security measures | API 包含多项安全措施：

- CORS middleware for cross-origin request handling | 用于处理跨域请求的 CORS 中间件
- Security headers middleware | 安全头中间件
- API version middleware for version control | 用于版本控制的 API 版本中间件

## Build Instructions | 构建说明

### Initialize Build Directory | 初始化构建目录

Create and initialize the `build` directory in the project root:

在项目根目录下创建并初始化 `build` 目录：

```shell
mkdir -p build
cd build
cmake ..
```

### Sync Dependencies | 同步依赖

```shell
cmake --build . --target tidy
```

### Build Targets | 构建目标

#### Build for Linux | 构建 Linux 版本

```shell
cmake --build . --target build-linux
```

### Clean Build Files | 清理构建文件

```shell
cmake --build . --target clean_target
```

### Using CMake on Windows | 在 Windows 系统上使用 CMake

To use CMake on Windows, you need to install a compatible tool such as Git Bash or Cygwin. After installation, you can run `cmake` and `cmake --build` commands in these terminals.

在 Windows 系统上使用 CMake，您需要安装一个兼容的工具，例如 Git Bash 或 Cygwin。安装完成后，您可以在这些终端中运行 `cmake` 和 `cmake --build` 命令。

## License | 许可证

See LICENSE file for details. | 详见 LICENSE 文件。