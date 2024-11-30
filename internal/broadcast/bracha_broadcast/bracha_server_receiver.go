package brachabroadcast

import (
	"log"
	"net"

	"DAG_Reliable_Broadcast/internal/broadcast"
)

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
		node.Conn[remoteAddr] = conn

		go broadcast.HandleConnection(conn)
	}
}
