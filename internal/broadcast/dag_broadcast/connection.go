package dagbroadcast

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
		// log.Printf("[INFO] 收到消息: %s", message)

		// 处理接收到的消息
		var msgType struct {
			Type int `json:"type"`
		}
		if err := json.Unmarshal([]byte(message), &msgType); err != nil {
			log.Printf("[ERROR] 反序列化消息失败: %v", err)
			continue
		}

		switch msgType.Type {
		case 0: // SendMessage
			var sendMsg SendMessage
			if err := json.Unmarshal([]byte(message), &sendMsg); err == nil {
				Instance.SendPool <- sendMsg
				// 每份消息都开启一个协程去处理
				go Instance.ProcessSend()
			}
		case 1: // ResponseMessage
			var responseMsg ResponseMessage
			if err := json.Unmarshal([]byte(message), &responseMsg); err == nil {
				Instance.ResponsePool <- responseMsg
				// 每份消息都开启一个协程去处理
				go Instance.ProcessResponse()
			}
		}
	}
}
