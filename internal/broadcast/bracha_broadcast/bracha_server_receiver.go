package brachabroadcast

import (
	"log"
	"net"
	"sync"
)

var mu sync.Mutex // 添加局部互斥锁

func StartListener(node *Node, host, port string) {
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

		mu.Lock() // 加锁
		node.Conn[remoteAddr] = conn
		mu.Unlock() // 解锁

		// 新建实例
		instance = NewNodeExtentions(*node)

		// 消息存通道，开启通道处理器
		go HandleConnection(conn)

	}
}
