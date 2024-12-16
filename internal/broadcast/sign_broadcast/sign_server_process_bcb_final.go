package signbroadcast

import (
	"crypto/sha256"
	"log"

	"github.com/CarJJJJ/go-bls"
)

func (node *NodeExtention) ProcessBCBFinal() {
	select {
	case msg := <-node.BCBFinalPool:
		// log.Printf("[INFO] 收到BCBFinal消息, 唯一键: %s, from:%v", msg.UniqueIndex, msg.NodeID)

		// 验证签名
		sigmaCombine, err := node.System.SigFromBytes(msg.SigmaCombine)
		if err != nil {
			log.Printf("[ERROR] 签名转换失败: %v", err)
		}

		hash := sha256.Sum256(msg.Message)
		if !bls.Verify(sigmaCombine, hash, node.GroupKey) {
			log.Printf("[ERROR] 签名验证失败,份额来自:%v,uniqueIndex:%v", msg.NodeID, msg.UniqueIndex)
			return
		}

		// 验证成功，打印消息
		// log.Printf("[INFO] 门限签名验证成功,消息:%v,uniqueIndex:%v", msg.Message, msg.UniqueIndex)
		// 把签名加入到已验证的签名池
		node.Proof = append(node.Proof, sigmaCombine)
		// 构建Disperse消息

		// 纠删码对消息进行分片

		// 计算每个分片的大小
		shardSize := (len(msg.Message) + node.N - 1) / node.N

		// 创建数据分片
		data := make([][]byte, node.N)
		for i := range data {
			data[i] = make([]byte, shardSize)
		}

		// 将消息填充到数据分片中
		for i := 0; i < len(msg.Message); i++ {
			data[i/shardSize][i%shardSize] = msg.Message[i]
		}

		// log.Printf("[INFO] 分片前的数据: %v,纠删码分片数据: %v, uniqueIndex: %v", msg.Message, data, msg.UniqueIndex)

		err = node.ReedSolomonEncoder.Encode(data)
		if err != nil {
			log.Printf("[ERROR] 纠删码分片失败: %v", err)
			return
		}

		if _, ok := node.HadDisperseUniqueIndex.Get(msg.UniqueIndex); ok {
			// log.Printf("[INFO] 已经发送过分片消息, unique_index: %s", msg.UniqueIndex)
			return
		}

		// 把data的每一份分片发送给对应的人
		for i := 0; i < node.N; i++ {
			disperseMessage := DisperseMessage{
				NodeID:      node.Node.Id,
				Type:        BCBDisperseType,
				DataFrom:    data[i],
				UniqueIndex: msg.UniqueIndex,
			}
			node.SendDisperseMessage(disperseMessage, i)
		}
		node.HadDisperseUniqueIndex.Set(msg.UniqueIndex, 1)
	}
}
