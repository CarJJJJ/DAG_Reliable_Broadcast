package dagctbroadcast

import (
	"log"
	"time"
)

func (node *NodeExtention) StartInBufferListener() {
	// 每秒监听Inbuffer
	for {
		length := 0
		node.InBuffer.Range(func(key, value interface{}) bool {
			length++
			return true
		})

		if length >= node.Y {
			node.InBuffer.Range(func(key, value interface{}) bool {
				uniqueIndex := key.(string)
				sendMessage := value.(*SendMessage)
				// 1. 选择下两个要发送消息的节点
				selectedNodes := node.SelectNextNodes()

				// 2. 选择NodeID下未的UniqueIndex未被Echo的消息
				echoReferences := node.chooseEchoReferences(selectedNodes) // 至多2条Echo,至少0条

				// 附加：如果echoReferences为空，那么重选
				for len(echoReferences) == 0 {
					selectedNodes = node.SelectNextNodes()
					echoReferences = node.chooseEchoReferences(selectedNodes) // 至多2条Echo,至少0条
				}

				// 3. 选择Ready的UniqueIndex
				readyReferences := node.chooseReadyReferences() // 只有1条Ready

				// 4. 发送消息
				responseMessage := ResponseMessage{
					Type:            ResponseType,
					SendMessage:     *sendMessage,
					EchoReferences:  echoReferences,
					ReadyReferences: readyReferences,
					NodeID:          node.Node.Id,
					UniqueIndex:     uniqueIndex,
				}
				log.Printf("[INFO] Inbuffer发送消息: R%v_%v", responseMessage.UniqueIndex, responseMessage.NodeID)
				node.BroadcastResposneToServers(&responseMessage)
				node.InBuffer.Delete(key)
				return true
			})
		}
		time.Sleep(time.Second)
	}
}
