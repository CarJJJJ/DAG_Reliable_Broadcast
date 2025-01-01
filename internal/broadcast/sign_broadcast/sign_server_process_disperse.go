package signbroadcast

func (node *NodeExtention) ProcessDisperse() {

	msg := <-node.BCBDispersePool
	// log.Printf("[INFO] 收到分片消息, unique_index: %s, from: %d", msg.UniqueIndex, msg.NodeID)

	if _, ok := node.RecvDisperseMessageForUniqueIndexNumber.Get(msg.UniqueIndex); !ok {
		node.RecvDisperseMessageForUniqueIndexNumber.Set(msg.UniqueIndex, 1)
	} else {
		number, _ := node.RecvDisperseMessageForUniqueIndexNumber.Get(msg.UniqueIndex)
		node.RecvDisperseMessageForUniqueIndexNumber.Set(msg.UniqueIndex, number+1)
	}

	// 如果收到T+1个分片消息，开始进行重构
	if number, ok := node.RecvDisperseMessageForUniqueIndexNumber.Get(msg.UniqueIndex); ok && number >= node.T+1 {
		// 如果还没有进行重构，则进行重构
		if _, ok := node.HadReconstructUniqueIndex.Get(msg.UniqueIndex); !ok {
			reconstructMessage := ReconstructMessage{
				NodeID:      node.Node.Id,
				Type:        BCBReconstructType,
				DataFrom:    msg.DataFrom,
				UniqueIndex: msg.UniqueIndex,
			}
			node.HadReconstructUniqueIndex.Set(msg.UniqueIndex, 1)
			node.BroadcastBCBReconstructMessage(reconstructMessage)
		}
	}

}
