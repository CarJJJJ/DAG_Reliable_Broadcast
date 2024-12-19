package dagbroadcast

import "log"

func (node *NodeExtention) ProcessResponse() {
	select {
	case msg := <-node.ResponsePool:
		log.Printf("[INFO] 收到Response消息: uniqueIndex: %s, from:%v", msg.UniqueIndex, msg.NodeID)
	}
}
