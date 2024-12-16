package signbroadcast

import (
	"crypto/sha256"

	"github.com/CarJJJJ/go-bls"
)

func (node *NodeExtention) ProcessBCBSend() {
	select {
	case msg := <-node.BCBSendPool:
		uniqueIndex := msg.UniqueIndex
		receivedMessage := msg.Message

		// log.Printf("[INFO] 收到BCBSend消息: uniqueIndex: %s, from:%v", uniqueIndex, msg.NodeID)

		// 把msg.Message放到BCBSendMessage之后，发送给所有节点
		hash := sha256.Sum256(receivedMessage)
		sigmaFrom := bls.Sign(hash, node.MemberSecrets[node.Node.Id])
		defer sigmaFrom.Free()

		sigmaFromToBytes := node.System.SigToBytes(sigmaFrom)

		// 标记发送过Rep
		node.HadRepUniqueIndex.Set(uniqueIndex, 1)
		BCBSendMessage := BCBRepMessage{
			Type:        BCBRepType,
			Message:     msg.Message,
			SigmaFrom:   sigmaFromToBytes,
			NodeID:      node.Node.Id,
			UniqueIndex: uniqueIndex,
		}
		// 发送回去给Msg.NodeID
		node.SendBCBRepToServer(BCBSendMessage, msg.NodeID)
	}
}
