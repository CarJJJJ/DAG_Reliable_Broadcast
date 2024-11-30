package brachabroadcast

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func BroadcastToServers(node Node) {
	ticker := time.NewTicker(1 * time.Second) // 每秒触发一次
	defer ticker.Stop()                       // 确保在函数结束时停止计时器

	for range ticker.C {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		initialMessage := InitialMessage{
			Type:    0,                                                           // 设置消息类型
			Message: fmt.Sprintf("基于bracha的可靠广播定时发送的消息,当前时间: %s\n", currentTime), // 消息内容
			NodeID:  node.Id,                                                     // 设置节点ID
		}

		message, err := json.Marshal(initialMessage) // 序列化为JSON
		if err != nil {
			log.Printf("[ERROR] 序列化消息失败: %v", err)
			continue
		}

		message = append(message, '\n')
		log.Print("[INFO] 广播消息: ", string(message))
		for addr, conn := range node.Conn {
			if _, err := fmt.Fprint(conn, string(message)); err != nil {
				log.Printf("[ERROR] 发送消息到 %s 失败: %v", addr, err)
				continue
			}
		}
	}
}
