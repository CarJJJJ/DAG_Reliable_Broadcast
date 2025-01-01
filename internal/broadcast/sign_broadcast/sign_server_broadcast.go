package signbroadcast

import (
	"encoding/json"
	"fmt"
	"log"
)

// SendBCBRepToServer 发送BCBRepMessage给指定节点
func (instance *NodeExtention) SendBCBRepToServer(bcbRepMessage BCBRepMessage, nodeID int) {
	message, err := json.Marshal(bcbRepMessage) // 序列化为JSON
	if err != nil {
		log.Printf("[ERROR] 序列化消息失败: %v", err)
	}

	message = append(message, '\n')

	// Conn的构成是 key为serverAddr := fmt.Sprintf("%s:%s", server.Host, server.Port),value为连接
	// 根据NodeId去/config/host_config.json中找到对应的serverAddr
	serverAddr := instance.Config.Servers[nodeID]
	Host := serverAddr.Host
	Port := serverAddr.Port
	serverAddrStr := fmt.Sprintf("%s:%s", Host, Port)
	conn, ok := instance.Node.Conn[serverAddrStr]
	if !ok {
		log.Printf("[ERROR] 连接到 %s 失败", serverAddrStr)
		return
	}
	conn.Write(message)
}

func (instance *NodeExtention) SendDisperseMessage(disperseMessage DisperseMessage, nodeID int) {
	message, err := json.Marshal(disperseMessage) // 序列化为JSON
	if err != nil {
		log.Printf("[ERROR] 序列化消息失败: %v", err)
	}

	message = append(message, '\n')

	serverAddr := instance.Config.Servers[nodeID]
	Host := serverAddr.Host
	Port := serverAddr.Port
	serverAddrStr := fmt.Sprintf("%s:%s", Host, Port)
	conn, ok := instance.Node.Conn[serverAddrStr]
	if !ok {
		log.Printf("[ERROR] 连接到 %s 失败", serverAddrStr)
		return
	}

	_, err = conn.Write(message)

	if err != nil {
		log.Printf("[ERROR] 发送消息到 %s 失败: %v", serverAddrStr, err)
	}
}

// BroadcastBCBRepToServers 实例广播BCBSend
func (instance *NodeExtention) BroadcastBCBRepToServers(bcbRepMessage BCBRepMessage) {
	message, err := json.Marshal(bcbRepMessage) // 序列化为JSON
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

// BroadcastBCBFinalMessage 广播BCBFinalMessage给所有节点
func (instance *NodeExtention) BroadcastBCBFinalMessage(bcbFinalMessage BCBFinalMessage) {
	message, err := json.Marshal(bcbFinalMessage) // 序列化为JSON
	if err != nil {
		log.Printf("[ERROR] 序列化消息失败: %v", err)
	}

	message = append(message, '\n')
	for addr, conn := range instance.Node.Conn {
		if _, err := fmt.Fprint(conn, string(message)); err != nil {
			log.Printf("[ERROR] 发送消息到 %s 失败: %v", addr, err)
			continue
		}
	}
}

// BroadcastBCBReconstructMessage 广播BCBReconstructMessage给所有节点
func (instance *NodeExtention) BroadcastBCBReconstructMessage(bcbReconstructMessage ReconstructMessage) {
	message, err := json.Marshal(bcbReconstructMessage) // 序列化为JSON
	if err != nil {
		log.Printf("[ERROR] 序列化消息失败: %v", err)
	}

	message = append(message, '\n')
	for addr, conn := range instance.Node.Conn {
		if _, err := fmt.Fprint(conn, string(message)); err != nil {
			log.Printf("[ERROR] 发送消息到 %s 失败: %v", addr, err)
			continue
		}
	}
}

// BroadcastReadyMessage 广播ReadyMessage给所有节点
func (instance *NodeExtention) BroadcastReadyMessage(readyMessage ReadyMessage) {
	message, err := json.Marshal(readyMessage) // 序列化为JSON
	if err != nil {
		log.Printf("[ERROR] 序列化消息失败: %v", err)
	}

	message = append(message, '\n')
	for addr, conn := range instance.Node.Conn {
		if _, err := fmt.Fprint(conn, string(message)); err != nil {
			log.Printf("[ERROR] 发送消息到 %s 失败: %v", addr, err)
			continue
		}
	}
}
