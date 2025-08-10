package domain

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type BlockChain struct {
	Pool  []*Mandate `json:"pool"`
	Chain []*Block   `json:"chain"`
}

func NewBlockChain() *BlockChain {
	bc := LoadDatabase()
	if len(bc.Chain) == 0 {
		bc.CreateGenesis()
		bc.CreateBlock(fmt.Sprintf("%x", [32]byte{}))
	}
	return bc
}

func LoadDatabase() *BlockChain {
	f, err := os.OpenFile("database/blockchain.db", os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		os.Exit(1)
	}
	scanner := bufio.NewScanner(f)
	blockChain := BlockChain{}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			os.Exit(1)
		}

		var blockSerialized BlockSerialized
		err := json.Unmarshal(scanner.Bytes(), &blockSerialized)
		if err != nil {
			os.Exit(1)
		}

		if len(blockChain.Chain) > 0 && (blockChain.LatestBlock().Hash() != blockSerialized.Value.Header.PrevHash) {
			log.Fatal("Invalid blockchain database")
		}

		blockChain.Chain = append(blockChain.Chain, blockSerialized.Value)
	}

	return &blockChain
}

func (bc *BlockChain) CreateBlock(prevHash string) *Block {
	b := NewBlock(prevHash, bc.Pool)
	bc.Chain = append(bc.Chain, b)
	bc.Pool = []*Mandate{}
	return b
}

func (bc *BlockChain) LatestBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *BlockChain) GiveMandate(from, to string, value int8) bool {
	if bc.CalculateMandate(from) < int64(value) {
		return false
	}
	m := NewMandate(from, to, value)
	bc.Pool = append(bc.Pool, m)
	return true
}

func (bc *BlockChain) CreateGenesis() {
	m := NewMandate("GOVERNMENT", "KPU", 20)
	bc.Pool = append(bc.Pool, m)
}

func (bc *BlockChain) CalculateMandate(user string) int64 {
	var total int64
	for _, v := range bc.Chain {
		for _, v2 := range v.Mandates {
			if v2.To == user {
				total += int64(v2.Value)
			}

			if v2.From == user {
				total -= int64(v2.Value)
			}

		}
	}

	return total
}

func (bc *BlockChain) PlenaryRecap() {
	for range time.Tick(time.Second * 10) {
		log.Println("Plenary Recap:", time.Now().Format(time.TimeOnly))
		if len(bc.Pool) > 0 {
			bc.CreateBlock(bc.LatestBlock().Hash())
		}
	}
}
