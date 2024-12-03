package brachabroadcast

import (
	"encoding/json"
	"fmt"
	"log"
)

// BroadcastEchoToServers 实例广播Echo
func (instance *NodeExtention) BroadcastEchoToServers(echoMessage EchoMessage) {
	message, err := json.Marshal(echoMessage) // 序列化为JSON
	if err != nil {
		log.Printf("[ERROR] 序列化消息失败: %v", err)
	}

	message = append(message, '\n')
	// log.Print("[INFO] 广播消息: ", string(message))
	for addr, conn := range instance.Node.Conn {
		if _, err := fmt.Fprint(conn, string(message)); err != nil {
			log.Printf("[ERROR] 发送消息到 %s 失败: %v", addr, err)
			continue
		}
	}
}

// BroadcastReadyToServers 实例广播Ready
func (instance *NodeExtention) BroadcastReadyToServers(readyMessage ReadyMessage) {
	message, err := json.Marshal(readyMessage) // 序列化为JSON
	if err != nil {
		log.Printf("[ERROR] 序列化消息失败: %v", err)
	}

	message = append(message, '\n')
	// log.Print("[INFO] 广播消息: ", string(message))
	for addr, conn := range instance.Node.Conn {
		if _, err := fmt.Fprint(conn, string(message)); err != nil {
			log.Printf("[ERROR] 发送消息到 %s 失败: %v", addr, err)
			continue
		}
	}
}
