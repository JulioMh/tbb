package database

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Snapshot [32]byte

type State struct {
	Balances  map[Account]uint
	txMempool []Tx
	dbFile    *os.File
	snapshot  Snapshot
}

func (s *State) Close() {
	s.dbFile.Close()
}

func (s *State) Add(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}
	s.txMempool = append(s.txMempool, tx)

	return nil
}

func (s *State) doSnapshot() error {
	_, err := s.dbFile.Seek(0, 0)
	if err != nil {
		return err
	}

	txsData, err := ioutil.ReadAll(s.dbFile)
	if err != nil {
		return err
	}
	s.snapshot = Snapshot(sha256.Sum256(txsData))

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

func (s *State) Persist() (Snapshot, error) {
	mempool := make([]Tx, len(s.txMempool))
	copy(mempool, s.txMempool)

	for _, tx := range s.txMempool {
		txJson, err := json.Marshal(tx)
		if err != nil {
			return Snapshot{}, err
		}
		fmt.Printf("Persisting new TX to disk:\n")
		fmt.Printf("\t%s\n\n", txJson)
		if _, err := s.dbFile.Write(append(txJson, '\n')); err != nil {
			return Snapshot{}, err
		}
		err = s.doSnapshot()
		if err != nil {
			return Snapshot{}, err
		}

		s.txMempool = s.txMempool[1:]
	}
	return s.snapshot, nil
}
