package dagbroadcast

import "log"

func (node *NodeExtention) ProcessSend() {
	select {
	case msg := <-node.SendPool:
		uniqueIndex := msg.UniqueIndex
		log.Printf("[INFO] 收到Send消息: uniqueIndex: %s, from:%v", uniqueIndex, msg.NodeID)
	}
}
