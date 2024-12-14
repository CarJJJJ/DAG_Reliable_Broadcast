package signbroadcast

import "log"

func (node *NodeExtention) ProcessBCBFinal() {
	select {
	case msg := <-node.BCBFinalPool:
		log.Printf("[INFO] 收到BCBFinal消息, 唯一键: %s, from:%v", msg.UniqueIndex, msg.NodeID)
	}
}
