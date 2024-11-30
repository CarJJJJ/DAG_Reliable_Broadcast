package dagbroadcast

import (
	"fmt"
	"log"
	"time"
)

func BroadcastToServers(node Node) {
	ticker := time.NewTicker(1 * time.Second) // 每秒触发一次
	defer ticker.Stop()                       // 确保在函数结束时停止计时器

	for range ticker.C {
		message := "基于有向无环图的可靠广播定时发送的消息\n" // 这里可以替换为你想要发送的文案
		log.Print("[INFO] 广播消息: ", message)
		for addr, conn := range node.Conn {
			if _, err := fmt.Fprint(conn, message); err != nil {
				log.Printf("[ERROR] 发送消息到 %s 失败: %v", addr, err)
				continue
			}
		}
	}
}
