/////////////////////////////////////////////////////////////////////////////////////
// 这里是一个简单的 MQTT 客户端示例，用于订阅和发布消息。


// package main

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"strings"

// 	mqtt "github.com/eclipse/paho.mqtt.golang"
// 	_ "github.com/lib/pq"
// )

// // RoomData 定义房间信息的结构体
// type RoomData struct {
// 	Number      string `json:"number"`
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	Status      string `json:"status"`
// }

// // DeviceData 设备基本信息
// type DeviceData struct {
// 	UUID         string `json:"uuid"`
// 	Name         string `json:"name"`
// 	Type         string `json:"type"`
// 	Model        string `json:"model"`
// 	Manufacturer string `json:"manufacturer"`
// 	RoomID       int    `json:"room_id"`
// }

// // DeviceStatusData 设备状态信息
// type DeviceStatusData struct {
// 	DeviceID   int    `json:"device_id"`
// 	Status     string `json:"status"`
// 	UpdatedAt  string `json:"updated_at"`
// 	LastReport string `json:"last_reported_at"`
// }

// // 全局数据库连接
// var db *sql.DB

// // 连接 PostgreSQL 数据库
// func connectDB() {
// 	connStr := "user=youruser password=yourpassword dbname=yourdb sslmode=disable"
// 	var err error
// 	db, err = sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatal("Failed to open database:", err)
// 	}
// 	if err = db.Ping(); err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}
// 	fmt.Println("Connected to PostgreSQL database")
// }

// // 处理 rooms 数据
// func insertOrUpdateRoom(data RoomData) {
// 	query := `
//         INSERT INTO rooms (number, name, description, status)
//         VALUES ($1, $2, $3, $4)
//         ON CONFLICT (number) DO UPDATE
//         SET name = $2, description = $3, status = $4, updated_at = CURRENT_TIMESTAMP`
// 	_, err := db.Exec(query, data.Number, data.Name, data.Description, data.Status)
// 	if err != nil {
// 		log.Println("Error inserting/updating room data:", err)
// 	} else {
// 		fmt.Println("Room data inserted/updated successfully:", data)
// 	}
// }

// // 处理 devices 数据
// func insertOrUpdateDevice(data DeviceData) {
// 	query := `
//         INSERT INTO devices (uuid, name, type, model, manufacturer, room_id)
//         VALUES ($1, $2, $3, $4, $5, $6)
//         ON CONFLICT (uuid) DO UPDATE
//         SET name = $2, type = $3, model = $4, manufacturer = $5, room_id = $6, updated_at = CURRENT_TIMESTAMP`
// 	_, err := db.Exec(query, data.UUID, data.Name, data.Type, data.Model, data.Manufacturer, data.RoomID)
// 	if err != nil {
// 		log.Println("Error inserting/updating device data:", err)
// 	} else {
// 		fmt.Println("Device data inserted/updated successfully:", data)
// 	}
// }

// // 处理 device_status 数据
// // insertOrUpdateDeviceStatus 函数用于插入或更新设备状态
// func insertOrUpdateDeviceStatus(data DeviceStatusData) {
// 	// 定义插入或更新设备的SQL语句
// 	query := `
//         INSERT INTO device_status (device_id, status, updated_at, last_reported_at)
//         VALUES ($1, $2, $3, $4)
//         ON CONFLICT (device_id) DO UPDATE
//         SET status = $2, updated_at = $3, last_reported_at = $4`
// 	// 执行SQL语句
// 	_, err := db.Exec(query, data.DeviceID, data.Status, data.UpdatedAt, data.LastReport)
// 	// 检查是否有错误
// 	if err != nil {
// 		log.Println("Error inserting/updating device status:", err)
// 	} else {
// 		fmt.Println("Device status data inserted/updated successfully:", data)
// 	}
// }

// // 解析 MQTT 消息并插入/更新数据库
// func parseAndInsert(topic string, payload []byte) {
// 	switch {
// 	case strings.Contains(topic, "rooms/data"):
// 		var data RoomData
// 		err := json.Unmarshal(payload, &data)
// 		if err != nil {
// 			log.Println("Error parsing room JSON:", err)
// 			return
// 		}
// 		if data.Number == "" || data.Name == "" {
// 			log.Println("Room Number and Name are required fields")
// 			return
// 		}
// 		insertOrUpdateRoom(data)

// 	case strings.Contains(topic, "devices/data"):
// 		var data DeviceData
// 		err := json.Unmarshal(payload, &data)
// 		if err != nil {
// 			log.Println("Error parsing device JSON:", err)
// 			return
// 		}
// 		if data.UUID == "" || data.Name == "" || data.Type == "" {
// 			log.Println("Device UUID, Name, and Type are required fields")
// 			return
// 		}
// 		insertOrUpdateDevice(data)

// 	case strings.Contains(topic, "device_status/data"):
// 		var data DeviceStatusData
// 		err := json.Unmarshal(payload, &data)
// 		if err != nil {
// 			log.Println("Error parsing device status JSON:", err)
// 			return
// 		}
// 		if data.DeviceID == 0 {
// 			log.Println("Device ID is required")
// 			return
// 		}
// 		insertOrUpdateDeviceStatus(data)

// 	default:
// 		log.Println("Unknown topic, ignoring message:", topic)
// 	}
// }

// // MQTT 消息处理函数
// var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
// 	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
// 	parseAndInsert(msg.Topic(), msg.Payload())
// }

// func main() {
// 	// 连接数据库
// 	connectDB()

// 	// 配置 MQTT 客户端
// 	opts := mqtt.NewClientOptions().
// 		// AddBroker("tcp://localhost:1883").
// 		// AddBroker("tcp://emqx@127.0.0.1:1883").
// 		AddBroker("tcp://broker.emqx.io:1883").
// 		SetClientID("go-mqtt-client").
// 		SetDefaultPublishHandler(messageHandler)

// 	// 创建并连接 MQTT 客户端
// 	client := mqtt.NewClient(opts)
// 	if token := client.Connect(); token.Wait() && token.Error() != nil {
// 		log.Fatal("Failed to connect to MQTT broker:", token.Error())
// 	}
// 	fmt.Println("Connected to MQTT broker")

// 	// 订阅不同的数据主题
// 	topics := []string{"rooms/data", "devices/data", "device_status/data"}
// 	for _, topic := range topics {
// 		if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
// 			log.Fatalf("Failed to subscribe to topic %s: %v", topic, token.Error())
// 		}
// 		fmt.Println("Subscribed to topic:", topic)
// 	}

// 	// 保持程序运行
// 	select {}
// }
////////////////////////////////////////////////////////////////////////////////////

// 这个 Go 语言程序用于接收来自 EMQX 的 Webhook 请求，并根据 topic（如 rooms/data、devices/data、
// device_status/data）将数据解析并存入 PostgreSQL 数据库。
// 程序提供 /mqtt/webhook 接口，接收 Webhook 请求后，根据 topic 解析 JSON 数据，并写入相应的数据表
// （rooms、devices、device_status）。它还包含数据库连接管理、数据去重与更新逻辑，确保数据一致性。
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

// RoomData 定义房间信息的结构体
type RoomData struct {
	Number      string `json:"number"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// DeviceData 设备基本信息
type DeviceData struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Model        string `json:"model"`
	Manufacturer string `json:"manufacturer"`
	RoomID       int    `json:"room_id"`
}

// DeviceStatusData 设备状态信息
type DeviceStatusData struct {
	DeviceID   int    `json:"device_id"`
	Status     string `json:"status"`
	UpdatedAt  string `json:"updated_at"`
	LastReport string `json:"last_reported_at"`
}

// WebhookPayload EMQX Webhook负载结构
type WebhookPayload struct {
	Action  string          `json:"action"`
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
	// 其他可能的字段
	ClientID  string `json:"clientid"`
	Username  string `json:"username"`
	Timestamp int64  `json:"timestamp"`
}

// 全局数据库连接
var db *sql.DB

// 连接 PostgreSQL 数据库
func connectDB() {
	connStr := "user=youruser password=yourpassword dbname=yourdb sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Connected to PostgreSQL database")
}

// 处理 rooms 数据
func insertOrUpdateRoom(data RoomData) {
	query := `
        INSERT INTO rooms (number, name, description, status)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (number) DO UPDATE
        SET name = $2, description = $3, status = $4, updated_at = CURRENT_TIMESTAMP
    `
	_, err := db.Exec(query, data.Number, data.Name, data.Description, data.Status)
	if err != nil {
		log.Println("Error inserting/updating room data:", err)
	} else {
		fmt.Println("Room data inserted/updated successfully:", data)
	}
}

// 处理 devices 数据
func insertOrUpdateDevice(data DeviceData) {
	query := `
        INSERT INTO devices (uuid, name, type, model, manufacturer, room_id)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (uuid) DO UPDATE
        SET name = $2, type = $3, model = $4, manufacturer = $5, room_id = $6, updated_at = CURRENT_TIMESTAMP
    `
	_, err := db.Exec(query, data.UUID, data.Name, data.Type, data.Model, data.Manufacturer, data.RoomID)
	if err != nil {
		log.Println("Error inserting/updating device data:", err)
	} else {
		fmt.Println("Device data inserted/updated successfully:", data)
	}
}

// 处理 device_status 数据
func insertOrUpdateDeviceStatus(data DeviceStatusData) {
	query := `
        INSERT INTO device_status (device_id, status, updated_at, last_reported_at)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (device_id) DO UPDATE
        SET status = $2, updated_at = $3, last_reported_at = $4
    `
	_, err := db.Exec(query, data.DeviceID, data.Status, data.UpdatedAt, data.LastReport)
	if err != nil {
		log.Println("Error inserting/updating device status:", err)
	} else {
		fmt.Println("Device status data inserted/updated successfully:", data)
	}
}

// 解析Webhook负载并处理数据
func parseAndInsert(topic string, payload []byte) {
	switch {
    case strings.Contains(topic, "rooms/data"):
        var data RoomData
        // 尝试直接解析 payload
        err := json.Unmarshal([]byte(payload), &data)
        if err != nil {
            // 如果失败，可能是字符串，尝试二次解析
            var payloadStr string
            if json.Unmarshal([]byte(payload), &payloadStr) == nil {
                // 如果 payload 是一个字符串，尝试解析该字符串为 RoomData
                err = json.Unmarshal([]byte(payloadStr), &data)
                if err != nil {
                    log.Println("Error parsing room JSON from string payload:", err)
                    return
                }
            } else {
                log.Println("Error parsing room JSON:", err)
                return
            }
        }
        if data.Number == "" || data.Name == "" {
            log.Println("Room Number and Name are required fields")
            return
        }
        insertOrUpdateRoom(data)

	case strings.Contains(topic, "devices/data"):
		var data DeviceData
        // 尝试直接解析 payload
        err := json.Unmarshal([]byte(payload), &data)
        if err != nil {
            // 如果失败，可能是字符串，尝试二次解析
            var payloadStr string
            if json.Unmarshal([]byte(payload), &payloadStr) == nil {
                // 如果 payload 是一个字符串，尝试解析该字符串为 RoomData
                err = json.Unmarshal([]byte(payloadStr), &data)
                if err != nil {
                    log.Println("Error parsing room JSON from string payload:", err)
                    return
                }
            } else {
                log.Println("Error parsing room JSON:", err)
                return
            }
        }
		if data.UUID == "" || data.Name == "" || data.Type == "" {
			log.Println("Device UUID, Name, and Type are required fields")
			return
		}
		insertOrUpdateDevice(data)

	case strings.Contains(topic, "device_status/data"):
		var data DeviceStatusData
        // 尝试直接解析 payload
        err := json.Unmarshal([]byte(payload), &data)
        if err != nil {
            // 如果失败，可能是字符串，尝试二次解析
            var payloadStr string
            if json.Unmarshal([]byte(payload), &payloadStr) == nil {
                // 如果 payload 是一个字符串，尝试解析该字符串为 RoomData
                err = json.Unmarshal([]byte(payloadStr), &data)
                if err != nil {
                    log.Println("Error parsing room JSON from string payload:", err)
                    return
                }
            } else {
                log.Println("Error parsing room JSON:", err)
                return
            }
        }
		if data.DeviceID == 0 {
			log.Println("Device ID is required")
			return
		}
		insertOrUpdateDeviceStatus(data)

	default:
		log.Println("Unknown topic, ignoring message:", topic)
	}
}

// handleWebhook 处理来自EMQX的Webhook请求
func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	log.Printf("Received webhook payload: %s", string(body))

	var webhookPayload WebhookPayload
	if err := json.Unmarshal(body, &webhookPayload); err != nil {
		log.Println("Error parsing webhook payload:", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received webhook: Action=%s, Topic=%s", webhookPayload.Action, webhookPayload.Topic)

	if webhookPayload.Action == "" {
		parseAndInsert(webhookPayload.Topic, webhookPayload.Payload)
	} else {
		log.Printf("Unrecognized action: %s", webhookPayload.Action)
	}

	// 返回成功响应
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

func main() {
	// 连接数据库
	connectDB()

	// 配置HTTP服务器处理Webhook
	http.HandleFunc("/mqtt/webhook", handleWebhook)

	// 启动HTTP服务器
	port := 8080
	fmt.Printf("Webhook server starting on port %d...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
