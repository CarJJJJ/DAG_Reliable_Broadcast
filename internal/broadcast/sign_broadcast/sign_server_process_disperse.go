package signbroadcast

import "log"

func (node *NodeExtention) ProcessDisperse() {
	select {
	case msg := <-node.BCBDispersePool:
		log.Printf("[INFO] 收到分片消息, unique_index: %s, from: %d", msg.UniqueIndex, msg.NodeID)
	}
}
