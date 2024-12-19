package dagbroadcast

import (
	"log"
	"net"
)

func StartListener(node *Node, host, port string) {
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Printf("[ERROR] 服务器启动失败: %v", err)
		return
	}
	defer listener.Close()

	log.Printf("[INFO] 服务器正在监听端口 %s...", port)

	// 新建实例
	Instance = NewNodeExtentions(*node)
	go Instance.CountTPS()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[ERROR] 接受连接失败: %v", err)
			continue
		}

		remoteAddr := conn.RemoteAddr().String()
		ip, port, err := net.SplitHostPort(remoteAddr)
		if err != nil {
			log.Printf("[ERROR] 获取 IP 地址失败: %v", err)
			continue
		}
		log.Printf("[INFO] 新的连接: %s:%s", ip, port)

		// 消息存通道，开启通道处理器
		go HandleConnection(conn)
	}
}
