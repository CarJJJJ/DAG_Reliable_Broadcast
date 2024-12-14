package signbroadcast

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

		// 处理接收到的消息
		var msgType struct {
			Type int `json:"type"`
		}
		if err := json.Unmarshal([]byte(message), &msgType); err != nil {
			log.Printf("[ERROR] 反序列化消息失败: %v", err)
			continue
		}

		switch msgType.Type {
		case 0: // BCBSendMessage
			var bcbSendMsg BCBSendMessage
			if err := json.Unmarshal([]byte(message), &bcbSendMsg); err == nil {
				Instance.BCBSendPool <- bcbSendMsg
				// 每份消息都开启一个协程去处理
				go Instance.ProcessBCBSend()
			}
		case 1: // BCBRspMessage
			var bcbRspMsg BCBRepMessage
			if err := json.Unmarshal([]byte(message), &bcbRspMsg); err == nil {
				Instance.BCBRepPool <- bcbRspMsg
				// 每份消息都开启一个协程去处理
				go Instance.ProcessBCBRep()
			}
		case 2: // BCBFinalMessage
			var bcbFinalMsg BCBFinalMessage
			if err := json.Unmarshal([]byte(message), &bcbFinalMsg); err == nil {
				Instance.BCBFinalPool <- bcbFinalMsg
				// 每份消息都开启一个协程去处理
				go Instance.ProcessBCBFinal()
			}
		case 3: // DisperseMessage
			var disperseMsg DisperseMessage
			if err := json.Unmarshal([]byte(message), &disperseMsg); err == nil {
				Instance.BCBDispersePool <- disperseMsg
				// 每份消息都开启一个协程去处理
				go Instance.ProcessDisperse()
			}
		case 4: // ReconstructMessage
			var reconstructMsg ReconstructMessage
			if err := json.Unmarshal([]byte(message), &reconstructMsg); err == nil {
				Instance.BCBReconstructPool <- reconstructMsg
				// 每份消息都开启一个协程去处理
				go Instance.ProcessReconstruct()
			}

		default:
			log.Printf("[ERROR] 未知消息类型: %d", msgType.Type)
		}
	}
}
