package test_util

import (
	"crypto/sha256"
	"testing"
	"time"

	"math/rand"

	"github.com/enzoh/go-bls"
)

func TestBlsThreshold(test *testing.T) {
	message := "This is a message."

	// Generate key shares.
	params := bls.GenParamsTypeF(256)
	pairing := bls.GenPairing(params)
	system, err := bls.GenSystem(pairing)
	if err != nil {
		test.Fatal(err)
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 4
	t := 1
	groupKey, memberKeys, groupSecret, memberSecrets, err := bls.GenKeyShares(t, n, system)
	if err != nil {
		test.Fatal(err)
	}

	// Select group members.
	memberIds := rand.Perm(n)[:(n - t)]

	// Sign the message.
	hash := sha256.Sum256([]byte(message))
	shares := make([]bls.Signature, n-t)
	for i := 0; i < n-t; i++ {
		shares[i] = bls.Sign(hash, memberSecrets[memberIds[i]])
	}

	for i := 0; i < n-t; i++ {
		if !bls.Verify(shares[i], hash, memberKeys[memberIds[i]]) {
			test.Fatal("Failed to verify signature.")
		}
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

	// Clean up.
	signature.Free()
	groupKey.Free()
	groupSecret.Free()
	for i := 0; i < n-t; i++ {
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
