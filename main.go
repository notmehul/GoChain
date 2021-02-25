package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/notmehul/blockchain-go/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - adds a block to the chain")
	fmt.Println(" print - prints the blocks in the chain")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
		// this initiates a shutdown of the goroutine and gives badgerDB time
		// to get the keys and values correct in order to prevent data loss
	}
}

func (cli *CommandLine) AddBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Block Added to the chain :D")
}

func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()

		// copied from prev version of main lol :P
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)                           // approving :)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate())) //printing the validation output
		fmt.Println()

		if len(block.PrevHash) == 0 {
			// if the block has no previous hash :3
			// ie the chain ended :/
			break
		}
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	// just making flags for the commandline to accept
	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	// this is inside the add thingy, if a user types add then says block
	// then they can enter in the data
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		// parsing everything that comes after the first argument on the CLI

		blockchain.Handle(err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" { // if it's an empty string
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.AddBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}

func main() {

	defer os.Exit(0)
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close() // to close the db before the main ends

	// chain.AddBlock("First Block")
	// chain.AddBlock("Second Block")
	// chain.AddBlock("Third Block")

	cli := CommandLine{chain}
	cli.run()
}
