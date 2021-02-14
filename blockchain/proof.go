/*
* proof of work algorithm will be implemented here :)
* this is used to secure the blockchain by making the network do computational work to add a block to the chain
* work here is computational work, calculations and shiz :)
* like how it works on bitcoin
* proof of work is mainly used to make the chain secure, like jisne kaam kiya uska select hoga etc etc
*
* the work should also be difficult to do and easy to prove :D
 */

package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// take data from block
// create counter which starts from 0
// create a hash of the data plus the counter
// check the hash if it meets the requirements
// Requirements:
// First few bytes must have 0s

const Difficulty = 13

// we are keeping this constant but on a real blockchain this is incremented overtime to account for increasing computational power of miners

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	// lsh is left shift lol
	// 256 is the number of bytes in our hash
	// subtractiung the difficulty and using it to left shift :)

	pow := &ProofOfWork{b, target}

	return pow
	// what we did here is similar to the derive hash function in the other file, so that is kil
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

// this is made just to include the difficulty and the nonce in the [][]bytes thingy
// I don't understand this fully either xD
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	// big endian has the significant information(bytes) is stored first
	// the write func will take our number and decode it into bytes

	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
	// here we just return the bytes portion of our buffer
}

// basically the entire proof of work algorithm
func (pow *ProofOfWork) Run() (int, []byte) {

	/*
	* In this function we will:
	* -> prepare out data
	* -> hash it into sha256
	* -> convert it into a big integer
	* -> compare it with our target big int(inside ProofOfWork)
	 */

	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce) // prepare data
		hash = sha256.Sum256(data)  // hash to sha256

		fmt.Printf("\r%x", hash)  // just to see the process :)
		intHash.SetBytes(hash[:]) // convert to big int

		if intHash.Cmp(pow.Target) == -1 {
			break // this means that our hash is less than the target
			// which is like we have already signed the block
		} else {
			nonce++
		}

	}
	fmt.Println() // just to have a lil gap

	return nonce, hash[:]
}

// validation method to verify blocks added
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
