package test_util

import (
	"crypto/sha256"
	"fmt"
	"os"
	"testing"
	"time"

	"math/rand"

	"github.com/CarJJJJ/go-bls"
)

func TestBlsThreshold(test *testing.T) {
	message := "This is a message."

	// Generate key shares.
	const paramsFile = "bls_params.txt"
	var params bls.Params
	if _, err := os.Stat(paramsFile); os.IsNotExist(err) {
		params = bls.GenParamsTypeF(256)
		paramsBytes, _ := params.ToBytes()
		if err := os.WriteFile(paramsFile, paramsBytes, 0600); err != nil {
			test.Fatalf("保存参数失败: %v", err)
		}
	} else {
		data, err := os.ReadFile(paramsFile)
		if err != nil {
			test.Fatalf("读取参数文件失败: %v", err)
		}
		params, err = bls.ParamsFromBytes(data)
		if err != nil {
			test.Fatalf("加载参数失败: %v", err)
		}
		fmt.Printf("Params content: %s\n", string(data))
	}

	pairing := bls.GenPairing(params)

	const systemFile = "bls_system.bin"
	var system bls.System

	if _, err := os.Stat(systemFile); os.IsNotExist(err) {
		system, err = bls.GenSystem(pairing)
		if err != nil {
			test.Fatal(err)
		}

		systemBytes := system.ToBytes()
		if err := os.WriteFile(systemFile, systemBytes, 0600); err != nil {
			test.Fatalf("保存系统失败: %v", err)
		}
	} else {
		systemData, err := os.ReadFile(systemFile)
		if err != nil {
			test.Fatalf("读取系统文件失败: %v", err)
		}
		system, err = bls.SystemFromBytes(pairing, systemData)
		if err != nil {
			test.Fatalf("加载系统失败: %v", err)
		}
	}

	// 打印当前的 System 的字节表示
	systemBytes := system.ToBytes()
	fmt.Printf("本次的 System 字节表示: %x\n", systemBytes)

	rand.Seed(int64(1))
	n := 4
	t := 3
	groupKey, memberKeys, groupSecret, memberSecrets, err := bls.GenKeyShares(t, n, system)
	if err != nil {
		test.Fatal(err)
	}

	startTime := time.Now()

	// Select group members.
	memberIds := []int{0, 1, 2}

	hash := sha256.Sum256([]byte(message))
	// test groupkey
	signGroupKey := bls.Sign(hash, groupSecret)
	test.Logf("hash:%x", hash)
	test.Logf("本次groupSecret的签名:%v", fmt.Sprintf("%x", system.SigToBytes(signGroupKey)))
	// Sign the message.
	shares := make([]bls.Signature, t)
	for i := 0; i < t; i++ {
		shares[i] = bls.Sign(hash, memberSecrets[i])
		test.Logf("份额%d:%v", i, fmt.Sprintf("%x", system.SigToBytes(shares[memberIds[i]])))
	}
	for i := 0; i < t; i++ {
		// 验证份额
		if !bls.Verify(shares[i], hash, memberKeys[i]) {
			test.Fatal("Failed to verify signature.")
		}
		// 验证成功，打印份额
		test.Logf("份额%d验证成功,份额:%v", i, fmt.Sprintf("%x", system.SigToBytes(shares[memberIds[i]])))
	}

	// Recover the threshold signature.
	signature, err := bls.Threshold(shares, memberIds, system)
	if err != nil {
		test.Fatal(err)
	}

	// Verify the threshold signature.
	if !bls.Verify(signature, hash, groupKey) {
		test.Fatal("Failed to verify signature.")
	}

	test.Logf("本次的签名:%v", fmt.Sprintf("%x", system.SigToBytes(signature)))
	// 计算总耗时
	duration := time.Since(startTime)
	test.Logf("门限签名（除去生成参数）总耗时:%v", duration)

	// Clean up.
	signature.Free()
	groupKey.Free()
	groupSecret.Free()
	for i := 0; i < t; i++ {
		shares[i].Free()
	}
	for i := 0; i < n; i++ {
		memberKeys[i].Free()
		memberSecrets[i].Free()
	}
	system.Free()
	pairing.Free()
	params.Free()
}
