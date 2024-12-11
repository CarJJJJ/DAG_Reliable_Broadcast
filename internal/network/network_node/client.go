package networknode

import (
	"fmt"
	"log"
	"net"
	"time"

	brachabroadcast "DAG_Reliable_Broadcast/internal/broadcast/bracha_broadcast"
	dagbroadcast "DAG_Reliable_Broadcast/internal/broadcast/dag_broadcast"
	signbroadcast "DAG_Reliable_Broadcast/internal/broadcast/sign_broadcast"
	"DAG_Reliable_Broadcast/internal/config"
)

func StartClient(configPath string, broadcastType int, nodeId int) {
	log.Println("[INFO] 启动客户端...")

	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("[ERROR] 加载配置失败: %v", err)
		return
	}

	node := NewNode("client", nodeId)

	// 连接所有服务器
	for _, server := range config.Servers {
		address := fmt.Sprintf("%s:%s", server.Host, server.Port)
		if err := connectWithRetry(node, address); err != nil {
			log.Printf("[ERROR] 在多次尝试后仍无法连接到服务器 %s: %v", address, err)
			continue
		}
	}

	if broadcastType == dag_broadcastType {
		log.Printf("[DEBUG] 进入 DAG 广播逻辑") // 添加调试日志
		go dagbroadcast.BroadcastToServers(dagbroadcast.Node{
			NodeType: node.NodeType,
			Id:       node.Id,
			Conn:     node.Conn,
		})
	} else if broadcastType == bracha_broadcastType {
		log.Printf("[DEBUG] 进入 Bracha 广播逻辑") // 添加调试日志
		go brachabroadcast.BroadcastToServers(brachabroadcast.Node{
			NodeType: node.NodeType,
			Id:       node.Id,
			Conn:     node.Conn,
		})
	} else if broadcastType == sign_broadcastType {
		log.Printf("[DEBUG] 进入 Sign 广播逻辑") // 添加调试日志
		go signbroadcast.BroadcastToServers(signbroadcast.Node{
			NodeType: node.NodeType,
			Id:       node.Id,
			Conn:     node.Conn,
		})
	}

	// 阻塞主协程
	select {}
}

func connectWithRetry(node *Node, address string) error {
	maxRetries := 3
	retryInterval := time.Second * 2

	for i := 0; i < maxRetries; i++ {
		conn, err := net.Dial("tcp", address) // 使用 net 进行连接
		if err == nil {
			node.Conn[address] = conn
			log.Printf("[INFO] 成功连接到服务器: %s", address)
			return nil
		}

		log.Printf("[ERROR] 连接服务器 %s 失败 (尝试 %d/%d): %v", address, i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(retryInterval)
		}
	}

	return fmt.Errorf("[ERROR] 达到最大重试次数")
}
