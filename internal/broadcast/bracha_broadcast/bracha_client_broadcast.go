package brachabroadcast

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"math/rand"
)

func BroadcastToServers(node Node) {
	ticker := time.NewTicker(1 * time.Second) // 每秒触发一次
	defer ticker.Stop()                       // 确保在函数结束时停止计时器

	// 生成1M的随机消息内容
	const n = 1024                      // 例如，n = 1024
	messageContent := make([]byte, n*n) // 创建n*n B的字节切片
	for i := range messageContent {
		messageContent[i] = byte('A' + rand.Intn(26)) // 随机生成字符'A'到'Z'
	}

	count := 0
	for range ticker.C {
		if count >= 200 {
			break
		}
		for i := 0; i < 10; i++ { // 每秒发送100条消息

			currentTime := time.Now().Format("2006-01-02 15:04:05")
			initialMessage := InitialMessage{
				Type:    0,                                                                                                            // 设置消息类型
				Message: fmt.Sprintf("基于bracha定时发送的消息,当前时间: %s, uniqueIndex:%d, 内容:%s\n", currentTime, count, string(messageContent)), // 消息内容
				NodeID:  node.Id,                                                                                                      // 设置节点ID
			}
			count++

			message, err := json.Marshal(initialMessage) // 序列化为JSON
			if err != nil {
				log.Printf("[ERROR] 序列化消息失败: %v", err)
				continue
			}

			message = append(message, '\n')
			// log.Print("[INFO] 广播消息: ", string(message))
			for addr, conn := range node.Conn {
				if _, err := fmt.Fprint(conn, string(message)); err != nil {
					log.Printf("[ERROR] 发送消息到 %s 失败: %v", addr, err)
					continue
				}
			}
			log.Printf("[INFO] 广播条数: %v", count)
		}
	}
}
