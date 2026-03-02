# MVTable Server

## 功能特性

### 核心功能
- **多视图数据管理** - 支持表格视图、看板视图、表单视图等多种数据展示方式
- **项目管理** - 创建和管理多个数据项目，支持项目收藏和最近访问
- **字段系统** - 多种字段类型：文本、数字、日期、选择器、关联等
- **数据记录** - 完整的 CRUD 操作，支持批量处理
- **权限管理** - 精细化的权限控制，支持高级权限配置

### 协作功能
- **团队协作** - 创建团队邀请成员，共同管理项目
- **实时协作** - WebSocket 支持实时数据更新
- **文档管理** - 内置富文本文档功能

### 系统功能
- **用户系统** - 邮箱注册登录，支持多级权益体系
- **仪表盘** - 数据可视化仪表盘，支持多种图表类型
- **邀请码** - 灵活的邀请注册机制
- **网站配置** - 可定制的网站配置

### 基础设施
- **对象存储** - 支持阿里云 OSS 和 MinIO
- **邮件服务** - SMTP 邮件发送
- **缓存** - Redis 缓存支持
- **日志** - 结构化日志记录

## 技术栈

| 类别 | 技术 |
|------|------|
| 语言 | Go 1.24.9 |
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | PostgreSQL |
| 缓存 | Redis |
| 认证 | JWT |
| 依赖注入 | Google Wire |
| CLI | Cobra |
| 日志 | Zap |
| 对象存储 | 阿里云 OSS / MinIO |
| API 文档 | Swagger (OpenAPI 3.0) |

## 快速开始

### 环境要求

- Go 1.24+
- PostgreSQL 12+
- Redis 6+
- Make

### 安装步骤

1. **安装依赖**
```bash
make install
```

2. **配置环境**

项目供配置文件根据需要修改：
- `configs/config.dev.yaml` - 开发环境配置
- `configs/config.prod.yaml` - 生产环境配置

3. **初始化数据库**
```bash
# 创建数据库
createdb mvtable

# 运行迁移
psql -h localhost -U postgres -d mvtable -f scripts/create.sql
```

4. **启动服务**
```bash
# 开发模式 (热重载)
make dev

# 生产模式
make build
make run ENV=prod
```

服务启动后访问: http://localhost:8901

## 项目结构

```
mvtable-server/
├── cmd/                   # 命令行入口
│   ├── root.go            # 根命令
│   ├── server.go          # 服务启动命令
│   └── version.go         # 版本信息
├── configs/               # 配置文件
│   ├── config.dev.yaml    # 开发环境配置
│   └── config.prod.yaml   # 生产环境配置
├── docs/                  # API 文档
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/             # 内部包
│   ├── app/              # 业务模块
│   │   ├── user/         # 用户模块
│   │   ├── team/         # 团队模块
│   │   ├── mv_project/   # 项目模块
│   │   ├── mv_view/      # 视图模块
│   │   ├── mv_field/     # 字段模块
│   │   ├── mv_record/    # 记录模块
│   │   ├── file/         # 文件模块
│   │   ├── collaboration/# 协作模块
│   │   └── ...
│   ├── di/               # 依赖注入
│   ├── middleware/       # 中间件
│   ├── pkg/              # 公共包
│   ├── server/           # HTTP 服务器
│   └── storage/          # 存储层
├── pkg/                  # 公共库
│   ├── jwt/             # JWT 认证
│   ├── lexorank/        # 排名算法
│   ├── log/             # 日志
│   ├── mail/            # 邮件
│   └── oss/             # 对象存储
├── scripts/             # 脚本
│   └── create.sql       # 数据库建表脚本
├── main.go              # 程序入口
└── Makefile             # 构建脚本
```

## 配置说明

### 配置文件结构

```yaml
server:
  host: "0.0.0.0"          # 监听地址
  port: 8901               # 监听端口
  mode: "debug"            # 运行模式: debug, release, test
  read_timeout: "30s"      # 读取超时
  write_timeout: "30s"     # 写入超时
  idle_timeout: "120s"     # 空闲连接超时

database:
  driver: "postgres"       # 数据库驱动: postgres, mysql
  host: "localhost"        # 数据库主机地址
  port: 5432               # 数据库端口
  database: "mvtable"      # 数据库名称
  username: "postgres"     # 数据库用户名
  password: ""             # 数据库密码
  charset: "utf8mb4"       # 字符编码
  max_idle_conns: 10       # 最大空闲连接数
  max_open_conns: 100      # 最大打开连接数
  conn_max_lifetime: 3600  # 连接最大生命周期(秒)
  log_level: "info"        # 日志级别: debug, info, warn, error

redis:
  host: "localhost"        # Redis 主机地址
  port: 6379               # Redis 端口
  password: ""             # Redis 密码
  database: 0              # Redis 数据库编号
  pool_size: 10            # 连接池大小

jwt:
  secret: ""               # JWT 密钥
  access_time: "24h"       # 访问令牌有效期
  refresh_time: "168h"     # 刷新令牌有效期
  issuer: "mvtable"        # JWT 发行者

log:
  level: "debug"           # 日志级别: debug, info, warn, error
  format: ""               # 日志格式: 空字符串为 console, json 为 JSON 格式
  output: "stdout"         # 输出方式: stdout, file
  filename: "logs/server.log"  # 日志文件路径
  max_size: 100            # 单个日志文件最大大小(MB)
  max_age: 30              # 日志文件保留天数
  max_backups: 10          # 保留的日志文件数量
  compress: true           # 是否压缩日志文件

mail:
  host: ""                # SMTP 服务器地址
  port: 465               # SMTP 端口
  from_name: ""           # 发件人显示名称
  from_mail: ""           # 发件人邮箱
  user_name: ""           # SMTP 用户名
  auth_code: ""           # SMTP 授权码

oss:
  provider: "aliyun"       # 对象存储提供商: aliyun, minio
  minio_conf:
    endpoint: ""           # MinIO 端点地址
    access_url: ""         # 访问 URL
    access_key_id: ""      # Access Key ID
    secret_access_key: ""  # Secret Access Key
    bucket_name: ""        # 存储桶名称
    use_ssl: false         # 是否使用 SSL
  aliyun_conf:
    access_url: ""         # 阿里云 OSS 访问地址
    region: ""             # 阿里云区域
    access_key_id: ""      # Access Key ID
    access_key_secret: ""  # Access Key Secret
    bucket_name: ""        # 存储桶名称
```

## API 文档

启动服务后访问:
- Swagger UI: http://localhost:8901/swagger/index.html
- OpenAPI JSON: http://localhost:8901/swagger/doc.json
- OpenAPI YAML: http://localhost:8901/swagger/doc.yaml

## Makefile 命令

| 命令 | 说明 |
|------|------|
| `make help` | 显示帮助信息 |
| `make build` | 编译项目 |
| `make run` | 运行服务 |
| `make dev` | 开发模式(热重载) |
| `make clean` | 清理构建产物 |
| `make deps` | 下载依赖 |
| `make install` | 安装项目依赖和工具 |
| `make init` | 初始化项目 |
| `make generate` | 生成代码 |
| `make swag` | 生成 Swagger 文档 |
| `make wire` | 生成依赖注入代码 |
| `make version` | 显示版本信息 |
| `make fmt` | 格式化代码 |
| `make lint` | 运行代码检查 |

## 使用方法

### 命令行选项

```bash
# 查看版本
./mvtable version

# 启动服务（开发环境）
./mvtable server -e dev

# 启动服务（生产环境）
./mvtable server -e prod
```
