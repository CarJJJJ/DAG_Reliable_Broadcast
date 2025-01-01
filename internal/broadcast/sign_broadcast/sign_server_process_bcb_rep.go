package signbroadcast

import (
	"crypto/sha256"
	"log"
	"strconv"

	"github.com/CarJJJJ/go-bls"
)

func (node *NodeExtention) ProcessBCBRep() {
	msg := <-node.BCBRepPool
	uniqueIndex := msg.UniqueIndex
	// log.Printf("[INFO] 收到BCBRep消息, 唯一键: %s, from:%v", uniqueIndex, msg.NodeID)
	// shareVerify
	sigmaFrom, err := node.System.SigFromBytes(msg.SigmaFrom)
	if err != nil {
		log.Printf("[ERROR] 签名转换失败: %v", err)
	}
	hash := sha256.Sum256(msg.Message)
	if !bls.Verify(sigmaFrom, hash, node.MemberKeys[msg.NodeID]) {
		log.Printf("[ERROR] 签名份额验证失败,份额来自:%v,uniqueIndex:%v", msg.NodeID, uniqueIndex)
		return
	}
	// log.Printf("[INFO] 签名份额验证成功,份额来自:%v,uniqueIndex:%v", msg.NodeID, uniqueIndex)

	// 将sigmaFrom加入到node.Pset中

	node.PsetMu.Lock()
	uniqueIndexInt, err := strconv.Atoi(uniqueIndex)
	if err != nil {
		log.Printf("[ERROR] 唯一键转换失败: %v", err)
	}
	if _, ok := node.Pset[uniqueIndexInt]; !ok {
		node.Pset[uniqueIndexInt] = make(map[int]bls.Signature)
	}
	node.Pset[uniqueIndexInt][msg.NodeID] = sigmaFrom
	node.PsetMu.Unlock()

	if len(node.Pset[uniqueIndexInt]) >= node.N-node.T {
		// 如果Pset的长度等于N-T，则需要恢复门限签名
		// 如果唯一键已经存在，则跳过
		if _, ok := node.HadFinalUniqueIndex.Get(uniqueIndex); ok {
			// log.Printf("[INFO] 已经广播过BCBFinal消息, 唯一键: %s", uniqueIndex)
			return
		}

		// 标记发送过Final
		node.HadFinalUniqueIndex.Set(uniqueIndex, 1)
		// log.Printf("[INFO] 第一次广播BCBFinal消息, 唯一键: %s", uniqueIndex)

		memberIds := []int{}
		shares := []bls.Signature{}
		for key := range node.Pset[uniqueIndexInt] {
			// 根据Pset的Key先拼出MemberIds

			memberIds = append(memberIds, key)
			// 根据Pset的Value拼出shares
			shares = append(shares, node.Pset[uniqueIndexInt][key])
		}
		// 恢复门限签名
		signature, err := bls.Threshold(shares, memberIds, node.System)
		if err != nil {
			log.Printf("[ERROR] 门限签名恢复失败: %v", err)
		}

		// 将恢复的门限签名转换为字节数组
		sigmaCombine := node.System.SigToBytes(signature)

		// 广播BCBFinalMessage给所有节点
		BCBFinalMessage := BCBFinalMessage{
			NodeID:       node.Node.Id,
			Type:         BCBFinalType,
			UniqueIndex:  uniqueIndex,
			SigmaCombine: sigmaCombine,
			Message:      msg.Message,
		}

		node.BroadcastBCBFinalMessage(BCBFinalMessage)
	}

}
