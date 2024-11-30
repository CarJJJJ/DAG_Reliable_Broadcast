package networknode

import (
	"fmt"
	"log"
	"net"
	"time"

	"DAG_Reliable_Broadcast/internal/config"
	"DAG_Reliable_Broadcast/internal/network"
)

func StartClient(configPath string) {
	log.Println("[INFO] 启动客户端...")

	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("[ERROR] 加载配置失败: %v", err)
		return
	}

	node := NewNode("client", fmt.Sprintf("client-%d", 0))

	// 连接所有服务器
	for _, server := range config.Servers {
		address := fmt.Sprintf("%s:%s", server.Host, server.Port)
		if err := connectWithRetry(node, address); err != nil {
			log.Printf("[ERROR] 在多次尝试后仍无法连接到服务器 %s: %v", address, err)
			continue
		}
	}

	// 处理用户输入并广播消息
	handleUserInput(node)
}

func connectWithRetry(node *Node, address string) error {
	maxRetries := 3
	retryInterval := time.Second * 2

	for i := 0; i < maxRetries; i++ {
		conn, err := net.Dial("tcp", address) // 使用 net 进行连接
		if err == nil {
			node.conn[address] = conn
			log.Printf("[INFO] 成功连接到服务器: %s", address)
			go network.HandleConnection(conn)
			return nil
		}

		log.Printf("[ERROR] 连接服务器 %s 失败 (尝试 %d/%d): %v", address, i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(retryInterval)
		}
	}

	return fmt.Errorf("[ERROR] 达到最大重试次数")
}
func handleUserInput(node *Node) {
	ticker := time.NewTicker(1 * time.Second) // 每秒触发一次
	defer ticker.Stop()                       // 确保在函数结束时停止计时器

	for range ticker.C {
		message := "定时发送的消息\n" // 这里可以替换为你想要发送的文案
		log.Print("[INFO] 广播消息: ", message)
		for addr, conn := range node.conn {
			if _, err := fmt.Fprint(conn, message); err != nil {
				log.Printf("[ERROR] 发送消息到 %s 失败: %v", addr, err)
				continue
			}
		}
	}
}
