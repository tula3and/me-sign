package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"sync"

	"github.com/tula3and/me-sign/db"
	"github.com/tula3and/me-sign/utils"
)

type Block struct {
	FileName string
	Email    string
	Hash     string
	PrevHash string
	Height   int
}

type blockchain struct {
	NewestHash string
	Height     int
}

var b *blockchain
var once sync.Once

func toBytes(i interface{}) []byte {
	var blockBuffer bytes.Buffer
	encoder := gob.NewEncoder(&blockBuffer)
	utils.HandleErr(encoder.Encode(b))
	return blockBuffer.Bytes()
}

func fromBytes(i interface{}, data []byte) {
	encoder := gob.NewDecoder(bytes.NewReader(data)).Decode(b)
	utils.HandleErr(encoder)
}

func (b *blockchain) AddBlock(fileName, email string) {
	block := Block{fileName, email, "", b.NewestHash, b.Height + 1}
	payload := block.FileName + block.Email + block.PrevHash + fmt.Sprintf("%x", block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	db.SaveBlock(block.Hash, toBytes(block))
	b.NewestHash = block.Hash
	b.Height = block.Height
	db.SaveBlockchain(toBytes(b))
}

var ErrNotFound = errors.New("block not found")

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	fromBytes(block, blockBytes)
	return block, nil
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis", "None")
			} else {
				fromBytes(b, checkpoint)
			}
		})
	}
	return b
}
