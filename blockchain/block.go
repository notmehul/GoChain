package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	// using byte slices cos UTF-8 support OP
	Nonce int
}

/*
* calculating the hash no big deal :)
* RIP this function got bangalored :D
*
* func (b *Block) DeriveHash() {
*     info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
*     hash := sha256.Sum256(info)
*     b.Hash = hash[:]
* }
 */

// creating the block and calculating it's hash :D
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	// block.DeriveHash() rekt xD

	pow := NewProof(block)

	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// adds the first block to the chain :3
func Init() *Block {
	return CreateBlock("INIT", []byte{})
}

// we need this function to transfer things to the badgerDB
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	//

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

// this is for getting things bacc from the db
func (b *Block) Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
