package signbroadcast

import (
	"log"
)

func (node *NodeExtention) ProcessBCBRsp() {
	select {
	case msg := <-node.BCBRepPool:
		uniqueIndex := msg.UniqueIndex
		// 如果唯一键不存在，放入 map 中并处理消息
		node.HadRepUniqueIndex.Set(uniqueIndex, 1)
		log.Printf("[INFO] 收到BCBRep消息，唯一键: %s, from:%v", uniqueIndex, msg.NodeID)
	}
}
