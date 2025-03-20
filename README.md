# MQTT to PostgreSQL in Go

## 项目简介

本项目是一个使用 Go 语言开发的简单程序，它用于接收 MQTT 消息，解析后将数据写入 PostgreSQL 数据库。

## 环境要求

- Go 语言环境（可使用 `go version` 检查是否已安装）
- 运行中的 MQTT 服务器（如 [Mosquitto](https://mosquitto.org/)）
- PostgreSQL 数据库

## 依赖库

本项目依赖以下 Go 库：

- MQTT 客户端库：[github.com/eclipse/paho.mqtt.golang](https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang)
- PostgreSQL 驱动：[github.com/lib/pq](https://pkg.go.dev/github.com/lib/pq)

安装这些库的命令：

```sh
go get github.com/eclipse/paho.mqtt.golang
go get github.com/lib/pq
```

## 功能介绍

该程序实现了以下功能：

1. 连接到 MQTT 服务器并订阅指定主题。
2. 接收 MQTT 消息并解析 JSON 数据。
3. 连接到 PostgreSQL 数据库。
4. 将解析后的数据插入数据库。

## 数据库设置

在 PostgreSQL 中创建数据库和表：

```sql
CREATE DATABASE sensor_data_db;

\c sensor_data_db;

CREATE TABLE sensor_data (
    id SERIAL PRIMARY KEY,
    sensor_id VARCHAR(50),
    temperature FLOAT,
    humidity FLOAT,
    received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 配置文件

请在 `config.yaml`（或代码中直接修改）中设置数据库连接信息和 MQTT 服务器地址。例如：

```yaml
mqtt:
  broker: "tcp://localhost:1883"
  topic: "test/topic"
database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "yourpassword"
  dbname: "sensor_data_db"
```

## API 说明

### `POST /mqtt/webhook`

- **功能**: 接收 MQTT Webhook 请求。
- **请求头**: `Content-Type: application/json`
- **请求体示例**:

```json
{
    "action": "publish",
    "topic": "rooms/data",
    "payload": "{\"number\": \"101\", \"name\": \"Conference Room\", \"description\": \"A large meeting room\", \"status\": \"available\"}"
}
```

## 数据库表结构
### `rooms` 表

| 字段名      | 类型      | 说明             |
| ----------- | --------- | ---------------- |
| number      | TEXT      | 房间编号（主键） |
| name        | TEXT      | 房间名称         |
| description | TEXT      | 房间描述         |
| status      | TEXT      | 房间状态         |
| updated_at  | TIMESTAMP | 更新时间         |

### `devices` 表

| 字段名       | 类型      | 说明                 |
| ------------ | --------- | -------------------- |
| uuid         | TEXT      | 设备唯一标识（主键） |
| name         | TEXT      | 设备名称             |
| type         | TEXT      | 设备类型             |
| model        | TEXT      | 设备型号             |
| manufacturer | TEXT      | 设备制造商           |
| room_id      | INT       | 所属房间ID           |
| updated_at   | TIMESTAMP | 更新时间             |

### `device_status` 表

| 字段名           | 类型      | 说明           |
| ---------------- | --------- | -------------- |
| device_id        | INT       | 设备ID（主键） |
| status           | TEXT      | 设备状态       |
| updated_at       | TIMESTAMP | 更新时间       |
| last_reported_at | TIMESTAMP | 最后上报时间   |
## 运行步骤

### 1. 确保环境就绪

- MQTT 服务器（如 Mosquitto）运行在 `localhost:1883`
- PostgreSQL 数据库已创建，并包含 `sensor_data` 表
- 更新代码中的数据库连接字符串（或使用 `config.yaml`）

### 2. 编译和运行

在代码所在目录下运行：

```sh
go run main.go
```

或编译后运行：

```sh
go build -o mqtt_postgres
./mqtt_postgres
```

### 3. 测试

使用 MQTT 客户端（如 `mosquitto_pub`）发布测试消息：

```sh
mosquitto_pub -h localhost -t test/topic -m '{"sensor_id": "001", "temperature": 25.5, "humidity": 60.0}'
```

### 4. 验证数据

在 PostgreSQL 中查询数据：

```sql
SELECT * FROM sensor_data;
```

## 代码结构

```
.
├── main.go          # 主程序
├── config.yaml      # 配置文件
├── go.mod           # Go 模块管理文件
├── go.sum           # 依赖管理文件
└── README.md        # 项目文档
```

## 可能遇到的问题

1. **MQTT 连接失败**
   - 确保 MQTT 服务器正在运行，并且使用正确的 `broker` 地址。
   - 检查防火墙或网络配置。

2. **数据库连接失败**
   - 确保 PostgreSQL 服务正常运行，并允许远程连接（如 `pg_hba.conf` 配置）。
   - 确保数据库凭据正确。

3. **数据未正确插入**
   - 检查 `main.go` 中的 JSON 解析逻辑。
   - 确保 MQTT 消息格式正确。

## 未来改进

- 增加日志记录。
- 支持 TLS 认证的 MQTT 连接。
- 使用 Docker 容器化部署。
- 增加 Web API 接口以查询数据。

## 许可证

本项目基于 MIT 许可证发布，欢迎贡献和改进！

