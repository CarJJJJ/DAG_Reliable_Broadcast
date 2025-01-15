package signbroadcast

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func BroadcastToServers(node Node) {
	ticker := time.NewTicker(1 * time.Second) // 每秒触发一次
	defer ticker.Stop()                       // 确保在函数结束时停止计时器
	// 定义n的大小
	const n = 16 // 例如，n = 16
	// 生成n*n字节的随机消息内容
	messageContent := make([]byte, n*n) // 创建n*n B的字节切片
	for i := range messageContent {
		messageContent[i] = byte('A' + i%26)
	}

	// 确保messageContent的长度是n*n
	if len(messageContent) < n*n {
		messageContent = append(messageContent, make([]byte, n*n-len(messageContent))...)
	} else if len(messageContent) > n*n {
		messageContent = messageContent[:n*n]
	}
	time.Sleep(10 * time.Second) // 等待10s连接建立完再发消息

	count := 0
	for range ticker.C {
		if count >= 200 {
			break
		}
		for i := 0; i < 10; i++ { // 每秒发送10条消息
			bcbSendMessage := BCBSendMessage{
				Type:        0,                        // 设置消息类型
				Message:     messageContent,           // 消息
				NodeID:      node.Id,                  // 设置节点ID
				UniqueIndex: fmt.Sprintf("%d", count), // 设置唯一索引
			}
			count++

			message, err := json.Marshal(bcbSendMessage) // 序列化为JSON
			if err != nil {
				log.Printf("[ERROR] 序列化消息失败: %v", err)
				continue
			}

			message = append(message, '\n')
			for addr, conn := range node.Conn {
				if _, err := fmt.Fprint(conn, string(message)); err != nil {
					log.Printf("[ERROR] 发送消息到 %s 失败: %v", addr, err)
					continue
				}
			}
			// log.Print("[INFO] 广播消息, uniqueIndex: ", count)

		}
	}
}
