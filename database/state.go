package database

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type State struct {
	Balances  map[Account]uint
	txMempool []Tx
	dbFile    *os.File
	prevHash  Hash
}

func (s *State) Close() {
	s.dbFile.Close()
}

func (s *State) AddTx(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}
	s.txMempool = append(s.txMempool, tx)

	return nil
}

func (s *State) AddBlock(block Block) error {
	for _, tx := range block.TXs {
		err := s.AddTx(tx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s State) apply(tx Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}
	if s.Balances[tx.From] < tx.Value {
		return fmt.Errorf("insufficient balance")
	}
	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value
	return nil
}

func (s *State) Persist() (Hash, error) {
	block := NewBlock(
		s.prevHash,
		uint64(time.Now().Unix()),
		s.txMempool,
	)
	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, err
	}

	blockFs := BlockFs{blockHash, block}

	blockFsJson, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, err
	}

	fmt.Printf("Persisting new Block to disk:\n")
	fmt.Printf("\t%s\n", blockFsJson)

	_, err = s.dbFile.Write(append(blockFsJson, '\n'))
	if err != nil {
		return Hash{}, err
	}
	s.prevHash = blockHash
	s.txMempool = []Tx{}

	return s.prevHash, nil
}
