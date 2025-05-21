# MC Dashboard

一个基于 Go + Vue3 开发的 Minecraft 服务器管理面板，提供服务器管理、配置修改、日志查看等功能。

## 功能特性

- 服务器管理
  - 创建、启动、停止、删除服务器
  - 服务器状态监控
  - 服务器配置管理
  - 控制台日志实时查看
  - 服务器命令执行
- 存档管理
  - 上传、下载、删除存档
  - 应用存档到服务器
- 白名单管理
  - 添加、删除白名单用户
  - 白名单状态控制
- 用户认证
  - JWT 认证
  - 用户信息管理

## 技术栈

### 后端
- Go 1.21+
- Gin Web 框架
- GORM ORM 框架
- MySQL 数据库
- Logrus 日志
- JWT 认证

### 前端
- Vue 3
- Element Plus UI 框架
- Vite 构建工具
- Pinia 状态管理
- Vue Router 路由管理

## 项目结构

```
.
├── cmd/                    # 主程序入口
│   └── server/
├── internal/              # 内部包
│   ├── app/              # 应用层
│   │   ├── handler/      # HTTP 处理器
│   │   ├── infra/        # 基础设施
│   │   │   └── minecraft/# Minecraft 相关实现
│   │   ├── service/      # 业务服务
│   │   └── usecase/      # 用例层
│   ├── domain/           # 领域模型
│   └── pkg/              # 公共包
│       ├── middleware/   # 中间件
│       └── utils/        # 工具函数
├── web/                  # 前端项目
│   ├── src/             # 源代码
│   │   ├── api/         # API 接口
│   │   ├── components/  # 组件
│   │   ├── composables/ # 组合式函数
│   │   ├── router/      # 路由配置
│   │   ├── store/       # 状态管理
│   │   └── views/       # 页面视图
│   └── public/          # 静态资源
├── configs/             # 配置文件
├── go.mod
└── README.md
```

## 快速开始

1. 克隆项目
```bash
git clone https://github.com/candbright/mc-dashboard.git
cd mc-dashboard
```

2. 安装后端依赖
```bash
go mod download
```

3. 安装前端依赖
```bash
cd web
npm install
```

4. 配置
- 创建 MySQL 数据库
- 修改配置文件 `configs/config.yaml`
- 配置环境变量（可选）

5. 运行项目
```bash
# 运行后端
go run cmd/server/main.go

# 运行前端（开发模式）
cd web
npm run dev
```

## 环境变量

- MC_DASHBOARD_DB_HOST: 数据库主机地址
- MC_DASHBOARD_DB_PORT: 数据库端口
- MC_DASHBOARD_DB_USER: 数据库用户名
- MC_DASHBOARD_DB_PASSWORD: 数据库密码
- MC_DASHBOARD_DB_NAME: 数据库名称
- MC_DASHBOARD_JWT_SECRET: JWT 密钥
- MC_DASHBOARD_SERVER_PORT: 服务器端口
- MC_DASHBOARD_ROOT_DIR: Minecraft 服务器根目录

## API 文档

### 服务器管理

#### 获取服务器列表
- 请求方式：GET
- 路径：/api/v1/servers
- 参数：
  - page: 页码
  - size: 每页数量
  - order: 排序方向（asc/desc）
  - order_by: 排序字段

#### 创建服务器
- 请求方式：POST
- 路径：/api/v1/servers
- 请求体：
```json
{
    "name": "服务器名称",
    "description": "服务器描述",
    "version": "服务器版本"
}
```

#### 启动服务器
- 请求方式：POST
- 路径：/api/v1/servers/{id}/start

#### 停止服务器
- 请求方式：POST
- 路径：/api/v1/servers/{id}/stop

#### 获取控制台日志
- 请求方式：GET
- 路径：/api/v1/servers/{id}/console_log
- 参数：
  - line: 获取最后几行日志

### 存档管理

#### 获取存档列表
- 请求方式：GET
- 路径：/api/v1/saves
- 参数：
  - page: 页码
  - size: 每页数量

#### 上传存档
- 请求方式：POST
- 路径：/api/v1/saves
- 请求体：multipart/form-data
  - file: 存档文件

#### 应用存档
- 请求方式：POST
- 路径：/api/v1/saves/apply
- 请求体：
```json
{
    "save_id": "存档ID",
    "server_id": "服务器ID"
}
```

### 白名单管理

#### 获取白名单列表
- 请求方式：GET
- 路径：/api/v1/servers/{id}/allowlist

#### 添加白名单用户
- 请求方式：POST
- 路径：/api/v1/servers/{id}/allowlist/{username}

#### 删除白名单用户
- 请求方式：DELETE
- 路径：/api/v1/servers/{id}/allowlist/{username}

## 许可证

MIT License
