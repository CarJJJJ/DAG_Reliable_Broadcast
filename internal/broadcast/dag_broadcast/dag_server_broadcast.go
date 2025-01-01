package dagbroadcast

import (
	"encoding/json"
	"fmt"
	"log"
)

func (node *NodeExtention) BroadcastResposneToServers(responseMessage *ResponseMessage) {
	message, err := json.Marshal(responseMessage) // 序列化为JSON
	if err != nil {
		log.Printf("[ERROR] 序列化消息失败: %v", err)
		return
	}

	message = append(message, '\n')
	// log.Print("[INFO] 广播消息: ", string(message))
	for addr, conn := range node.Node.Conn {
		if _, err := fmt.Fprint(conn, string(message)); err != nil {
			log.Printf("[ERROR] 发送消息到 %s 失败: %v", addr, err)
			continue
		}
	}
}
