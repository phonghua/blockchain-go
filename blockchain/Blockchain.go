package blockchain 

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"log"
)

type Blockchain struct {
	Chain		[]Block
}

func NewBlockchain() *Blockchain {
	var blockChain = &Blockchain {
		Chain: []Block{},
	}

	newBlock, _ := blockChain.genesis()
	newChain := append(blockChain.Chain, newBlock)
	blockChain.ReplaceChain(newChain)

	log.Printf("chain %v\n %v\n %v\n", blockChain, newBlock, newChain)
	return blockChain
}


func (blockchain *Blockchain) CalculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BMP) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (blockchain *Blockchain) genesis() (Block, error){
	var newBlock Block

	t := time.Now()

	newBlock.Index = 0
	newBlock.Timestamp = t.String()
	newBlock.BMP = 0
	newBlock.PrevHash = "0000000000000000000000000000000000000000000000000000000000000000"
	newBlock.Hash = blockchain.CalculateHash(newBlock)
	return newBlock, nil
}

func (blockchain *Blockchain) GenerateBlock(BMP int) (Block, error){
	oldBlock := blockchain.Chain[len(blockchain.Chain ) - 1]
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BMP = BMP
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = blockchain.CalculateHash(newBlock)

	return newBlock, nil
}

func (blockchain *Blockchain) IsBlockValid(newBlock, oldBlock Block) bool{
	if oldBlock.Index + 1 != newBlock.Index{
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash{
		return false
	}
	if blockchain.CalculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

func (blockchain *Blockchain) ReplaceChain(newBlocks []Block){
	if(len(blockchain.Chain) < len(newBlocks)){
		blockchain.Chain = newBlocks
	}
}

