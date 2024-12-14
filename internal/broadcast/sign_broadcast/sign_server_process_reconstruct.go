package signbroadcast

import "log"

func (node *NodeExtention) ProcessReconstruct() {
	select {
	case msg := <-node.BCBReconstructPool:
		log.Printf("[INFO] 收到重构消息, unique_index: %s, from: %d", msg.UniqueIndex, msg.NodeID)
	}
}
