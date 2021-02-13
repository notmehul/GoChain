package blockchain

import (
	"bytes"
	"crypto/sha256"
)

type BlockChain struct {
	Blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	// using byte slices cos UTF-8 support OP
}

// calculating the hash no big dead :)
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

// creating the block and calculating it's hash :D
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

// adding the block to the chain
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1] // determines previous block in the chain
	new := CreateBlock(data, prevBlock.Hash)       // creates the new block using the old one for hash
	chain.Blocks = append(chain.Blocks, new)       // appends :)
}

// adds the first block to the chain :3
func Init() *Block {
	return CreateBlock("INIT", []byte{})
}

// the madlad which started everything
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Init()}}
}
