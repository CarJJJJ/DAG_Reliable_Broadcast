package networknode

import (
	"fmt"
	"log"
	"net"
	"time"

	bracha_broadcast "DAG_Reliable_Broadcast/internal/broadcast/bracha_broadcast"
	dag_broadcast "DAG_Reliable_Broadcast/internal/broadcast/dag_broadcast"
	"DAG_Reliable_Broadcast/internal/config"
)

func StartServer(host, port string, broadcastType, nodeId int) {
	log.Println("[INFO] 启动服务端...")

	config, err := config.LoadConfig("config/host_config.json")
	if err != nil {
		log.Printf("[ERROR] 加载配置失败: %v", err)
		return
	}

	node := NewNode("server", nodeId)

	// 连接到其他服务器
	connectToOtherServers(node, config)

	// 启动监听服务
	if broadcastType == dag_broadcastType {
		dag_broadcast.StartListener(&dag_broadcast.Node{
			NodeType: node.NodeType,
			Id:       node.Id,
			Conn:     node.Conn,
		}, host, port)
	} else if broadcastType == bracha_broadcastType {
		bracha_broadcast.StartListener(&bracha_broadcast.Node{
			NodeType: node.NodeType,
			Id:       node.Id,
			Conn:     node.Conn,
		}, host, port)
	}
}

func connectToOtherServers(node *Node, config *config.Config) {
	for _, server := range config.Servers {
		serverAddr := fmt.Sprintf("%s:%s", server.Host, server.Port)
		// 使用指数退避重试连接
		go func(addr string) {
			retryCount := 0
			maxRetries := 5 // 最大重试次数
			for {
				conn, err := net.Dial("tcp", addr)
				if err == nil {
					node.mu.Lock() // 在写入 map 之前加锁
					node.Conn[addr] = conn
					node.mu.Unlock() // 写入后解锁
					log.Printf("[INFO] 成功连接到服务器: %s", addr)
					return
				}

				retryCount++
				if retryCount > maxRetries {
					log.Printf("[ERROR] 连接服务器 %s 失败，已达到最大重试次数: %v", addr, err)
					return
				}

				// 计算等待时间：1秒、2秒、4秒、8秒、16秒
				waitTime := time.Second * time.Duration(1<<uint(retryCount-1))
				log.Printf("[ERROR] 连接服务器 %s 失败，%d秒后重试 (第%d次): %v",
					addr, waitTime/time.Second, retryCount, err)
				time.Sleep(waitTime)
			}
		}(serverAddr)
	}
}
