# HttpScheduleBE

HttpScheduleBE 是一个基于 Go 语言开发的 HTTP 定时任务调度后端服务。它支持通过 Web API 管理定时任务，并自动调度执行 HTTP 请求，记录执行结果。

## 主要功能

- 任务管理：创建、更新、删除、查询定时任务
- 任务调度：基于 cron 表达式自动调度任务
- HTTP 请求执行：支持自定义请求方法、Header、Body
- 执行记录：记录每次任务执行的状态、时间和错误日志
- 配置化：通过 `etc/config.yaml` 配置数据库和自动执行参数

## 快速开始

1. **准备数据库**  
   创建 MySQL 数据库，并在 `etc/config.yaml` 配置连接信息。

2. **安装依赖**  
   ```bash
   go mod tidy
   ```

3. **启动服务**  
   ```bash
   go run ./cmd/main.go
   ```
   服务默认监听 `:8080`。

## 配置说明

`etc/config.yaml` 示例：
```yaml
db_host: 127.0.0.1
db_port: 3306
db_user: root
db_password: yourpassword
db_name: yourdbname
execute_automatic: true
```

## 依赖

- Go 1.18+
- Gin
- GORM
- robfig/cron
- MySQL8


