package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/block"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// the madlad which started everything
func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions
	opts.Dir = dbPath      // keys and meta data is stored here
	opts.ValueDir = dbPath // values are stored here :)

	db, err := badger.Open(opts) //db is a pointer to the db
	Handle(err)

	// update function allows us to read and write to the db
	err = db.Update(func(txn *badger.Txn) error {
		//checking if there is another blockchain already :)
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			// lh == lasthash, if it shows an error so it doesn't exist
			fmt.Println("no blockchains found, creating a new one....")
			init := Init()
			fmt.Println("Blockchain Initialised....")
			err = txn.Set(init.Hash, init.Serialize())
			//hash is the key and serializing it to put in db uwu
			Handle(err)

			// setting init blocks hash as the last :)
			err = txn.Set([]byte("lh"), init.Hash)
			lastHash = init.Hash

			return err
		} else { // is blockchain already exists
			item, err := txn.Get([]byte("lh")) // this returns a pointer to the LastHash in the DB :3
			Handle(err)
			lastHash, err = item.Value() // got dem valuez
			return err
		}
	})
	Handle(err) // handling any error that maybe got passed by the anon func

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

// adding the block to the chain
func (chain *BlockChain) AddBlock(data string) {
	/*
		######
		SHIT LEGACY CODE^ XD
		######

		prevBlock := chain.Blocks[len(chain.Blocks)-1] // determines previous block in the chain
		new := CreateBlock(data, prevBlock.Hash)       // creates the new block using the old one for hash
		chain.Blocks = append(chain.Blocks, new)       // appends :)
	*/

	var lastHash []byte

	// this is a read only type of transaction
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)

		lastHash, err = item.Value()

		return err
	})

	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err = txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		// same ol thing we did with init blockchain func :3

		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err)

}

// this is where the fun begins :D
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

//as we start from the last hash, we will be going backwards in this
// this does what it seems it does lol, goes thru the chain
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)

		encodedBlock, err := item.Value()
		block = block.Deserialize((encodedBlock))
		// to make the data viewable :3

		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash
	// for going backwardsthru our database

	return block
}
