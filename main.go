package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type BlockChain struct {
	blocks []*Block
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
	prevBlock := chain.blocks[len(chain.blocks)-1] // determines previous block in the chain
	new := CreateBlock(data, prevBlock.Hash)       // creates the new block using the old one for hash
	chain.blocks = append(chain.blocks, new)       // appends :)
}

// adds the first block to the chain :3
func Init() *Block {
	return CreateBlock("INIT", []byte{})
}

// the madlad which started everything
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Init()}}
}

func main() {
	chain := InitBlockChain()

	chain.AddBlock("First Block")
	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")

	for _, block := range chain.blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
	}
}
