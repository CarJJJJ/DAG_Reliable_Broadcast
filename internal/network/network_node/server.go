package networknode

import (
	"fmt"
	"log"
	"net"
	"time"

	"DAG_Reliable_Broadcast/internal/config"
	"DAG_Reliable_Broadcast/internal/network"
)

func StartServer(host, port string) {
	log.Println("[INFO] 启动服务端...")

	config, err := config.LoadConfig("config/host_config.json")
	if err != nil {
		log.Printf("[ERROR] 加载配置失败: %v", err)
		return
	}

	node := NewNode("server", fmt.Sprintf("%s:%s", host, port))

	// 连接到其他服务器
	connectToOtherServers(node, config)

	// 启动监听服务
	startListener(node, host, port)
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
					node.conn[addr] = conn
					go network.HandleConnection(conn)
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

func startListener(node *Node, host, port string) {
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Printf("[ERROR] 服务器启动失败: %v", err)
		return
	}
	defer listener.Close()

	log.Printf("[INFO] 服务器正在监听端口 %s...", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[ERROR] 接受连接失败: %v", err)
			continue
		}

		remoteAddr := conn.RemoteAddr().String()
		ip, _, err := net.SplitHostPort(remoteAddr)
		if err != nil {
			log.Printf("[ERROR] 获取 IP 地址失败: %v", err)
			continue
		}
		log.Printf("[INFO] 新的连接: %s", ip)
		node.conn[remoteAddr] = conn

		go network.HandleConnection(conn)
	}
}
