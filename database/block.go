package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Hash [32]byte

type BlockFs struct {
	Hash  Hash  `json:"hash"`
	Block Block `json:"block"`
}

type Block struct {
	Header BlockHeader `json:"header"`
	TXs    []Tx        `json:"payload"`
}

type BlockHeader struct {
	Parent Hash   `json:"parent"`
	Time   uint64 `json:"time"`
}

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(h.Hex()), nil
}

func (h *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

func NewBlock(prevHash Hash, time uint64, txs []Tx) Block {
	blockHeader := BlockHeader{prevHash, time}
	return Block{blockHeader, txs}
}

func (b Block) Hash() (Hash, error) {
	blockJson, err := json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}
	return sha256.Sum256(blockJson), nil
}
