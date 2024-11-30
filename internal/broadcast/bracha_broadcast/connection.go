package brachabroadcast

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("[ERROR] 读取消息错误: %v", err)
			return
		}
		log.Printf("[INFO] 收到消息: %s", message)

		// 处理接收到的消息
		var msgType struct {
			Type int `json:"type"`
		}
		if err := json.Unmarshal([]byte(message), &msgType); err != nil {
			log.Printf("[ERROR] 反序列化消息失败: %v", err)
			continue
		}

		switch msgType.Type {
		case 0: // InitialMessage
			var initialMsg InitialMessage
			if err := json.Unmarshal([]byte(message), &initialMsg); err == nil {
				instance.InitialPool <- initialMsg
				log.Printf("[INFO] 存储 InitialMessage: %+v", initialMsg) // 添加日志
				// 每份消息都开启一个协程去处理
				go instance.ProcessInitial()
			}
		case 1: // EchoMessage
			var echoMsg EchoMessage
			if err := json.Unmarshal([]byte(message), &echoMsg); err == nil {
				instance.EchoPool <- echoMsg
				log.Printf("[INFO] 存储 EchoMessage: %+v", echoMsg) // 添加日志
			}
		case 2: // ReadyMessage
			var readyMsg ReadyMessage
			if err := json.Unmarshal([]byte(message), &readyMsg); err == nil {
				instance.ReadyPool <- readyMsg
				log.Printf("[INFO] 存储 ReadyMessage: %+v", readyMsg) // 添加日志
			}
		default:
			log.Printf("[ERROR] 未知消息类型: %d", msgType.Type)
		}
	}
}
