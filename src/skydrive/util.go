package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

func sha1hash(str string) string {
	// calculate sha-1 hash
	hasher := sha1.New()
	hasher.Write([]byte(str))

	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func distance(a, b []byte, bits int) *big.Int {
	var ring big.Int
	ring.Exp(big.NewInt(2), big.NewInt(int64(bits)), nil)

	var a_int, b_int big.Int
	(&a_int).SetBytes(a)
	(&b_int).SetBytes(b)

	var dist big.Int
	(&dist).Sub(&b_int, &a_int)

	(&dist).Mod(&dist, &ring)
	return &dist
}

//True if key is in [id1, id2)
func between(id1, id2, key []byte) bool {
	// 0 if a==b, -1 if a < b, and +1 if a > b

	if bytes.Compare(key, id1) == 0 { // key == id1
		return true
	}

	if bytes.Compare(id2, id1) == 1 { // id2 > id1
		if bytes.Compare(key, id2) == -1 && bytes.Compare(key, id1) == 1 { // key < id2 && key > id1
			return true
		} else {
			return false
		}
	} else { // id1 > id2
		if bytes.Compare(key, id1) == 1 || bytes.Compare(key, id2) == -1 { // key > id1 || key < id2
			return true
		} else {
			return false
		}
	}
}

// True if key is in (id1, id2]
func between2(id1, id2, key []byte) bool {
	// 0 if a==b, -1 if a < b, and +1 if a > b

	if bytes.Compare(key, id1) == 0 { // key == id1
		return false
	}

	if bytes.Compare(key, id2) == 0 { // key == id2
		return true
	}

	if bytes.Compare(id2, id1) == 1 { // id2 > id1
		if bytes.Compare(key, id2) == -1 && bytes.Compare(key, id1) == 1 { // key < id2 && key > id1
			return true
		} else {
			return false
		}
	} else { // id1 > id2
		if bytes.Compare(key, id1) == 1 || bytes.Compare(key, id2) == -1 { // key > id1 || key < id2
			return true
		} else {
			return false
		}
	}
}

// False if key == id1 || key == id2
//True if key is in (id1, id2)
func inside(id1, id2, key []byte) bool {

	if bytes.Compare(key, id1) == 0 { // key == id1
		return false
	}

	if bytes.Compare(key, id2) == 0 { // key == id2
		return false
	}
	// 0 if a==b, -1 if a < b, and +1 if a > b
	if bytes.Compare(id2, id1) == 1 { // id2 > id1
		if bytes.Compare(key, id2) == -1 && bytes.Compare(key, id1) == 1 { // key < id2 && key > id1
			return true
		} else {
			return false
		}
	} else { // id1 > id2
		if bytes.Compare(key, id1) == 1 || bytes.Compare(key, id2) == -1 { // key > id1 || key < id2
			return true
		} else {
			return false
		}
	}
}

// (n + 2^(k-1)) mod (2^m)
func calcFinger(n []byte, k int, m int) (string, []byte) {

	// convert the n to a bigint
	nBigInt := big.Int{}
	nBigInt.SetBytes(n)

	// get the right addend, i.e. 2^(k-1)
	two := big.NewInt(2)
	addend := big.Int{}
	addend.Exp(two, big.NewInt(int64(k-1)), nil)

	// calculate sum
	sum := big.Int{}
	sum.Add(&nBigInt, &addend)

	// calculate 2^m
	ceil := big.Int{}
	ceil.Exp(two, big.NewInt(int64(m)), nil)

	// apply the mod
	result := big.Int{}
	result.Mod(&sum, &ceil)

	resultBytes := result.Bytes()
	resultHex := fmt.Sprintf("%x", resultBytes)

	return resultHex, resultBytes
}

func calcFingerSha(n []byte, k int) (string, []byte) {
	m := 160
	return calcFinger(n, k, m)
}

func add(id string, added int64) (string, []byte) {

	idByte, err := hex.DecodeString(id)
	if err != nil {
		log.Fatal("decodeString error :", err)
	}

	idBigInt := big.Int{}
	idBigInt.SetBytes(idByte)

	// calculate sum
	sum := big.Int{}
	addInt := big.NewInt(added)
	sum.Add(&idBigInt, addInt)

	resultBytes := sum.Bytes()
	resultHex := fmt.Sprintf("%x", resultBytes)
	return resultHex, resultBytes
}

// (n - 2^(k-1)) mod 2^m
func calcLastFinger(n []byte, k int) (string, []byte) {

	// convert the n to a bigint
	nBigInt := big.Int{}
	nBigInt.SetBytes(n)

	// get the right addend, i.e. 2^(k-1)
	two := big.NewInt(2)
	addend := big.Int{}
	addend.Exp(two, big.NewInt(int64(k-1)), nil)

	addend.Mul(&addend, big.NewInt(-1))
	//Soustraction
	neg := big.Int{}
	neg.Add(&addend, &nBigInt)

	// calculate 2^m
	m := 160
	ceil := big.Int{}
	ceil.Exp(two, big.NewInt(int64(m)), nil)

	// apply the mod
	result := big.Int{}
	result.Mod(&neg, &ceil)

	resultBytes := result.Bytes()
	resultHex := fmt.Sprintf("%x", resultBytes)

	return resultHex, resultBytes
}

func DHTnodeToNode(dhtNode DHTnode) Node {
	node := new(Node)
	node.Id = dhtNode.Id
	node.IdByte = dhtNode.IdByte
	node.Ip = dhtNode.Ip
	node.Port = dhtNode.Port
	return *node
}
