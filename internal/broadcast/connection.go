package broadcast

import (
	"bufio"
	"log"
	"net"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("[ERROR] 读取消息错误: %v", err)
			return
		}
		log.Printf("[INFO] 收到消息: %s", message)
	}
}
