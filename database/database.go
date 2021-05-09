package database

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func loadGenesis(genFilePath string) (*State, error) {
	// Read genesis file
	f, err := ioutil.ReadFile(genFilePath)
	if err != nil {
		return nil, err
	}
	// Parse to map[string]interface{}
	var genesis map[string]interface{}
	json.Unmarshal(f, &genesis)
	// Create balances from genesis
	genesisBalances := make(map[Account]uint)
	for account, balance := range genesis["balances"].(map[string]interface{}) {
		genesisBalances[Account(account)] = uint(balance.(float64))
	}
	// Create first state
	var genesisState State
	genesisState.Balances = genesisBalances
	return &genesisState, nil
}

func NewStateFromDisk() (*State, Snapshot, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, Snapshot{}, err
	}
	genFilePath := filepath.Join(cwd, "database", "disk", "genesis.json")
	gen, err := loadGenesis(genFilePath)
	if err != nil {
		return nil, Snapshot{}, err
	}

	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	txDbFilePath := filepath.Join(cwd, "database", "disk", "tx.db")
	f, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, Snapshot{}, err
	}

	scanner := bufio.NewScanner(f)
	state := &State{balances, make([]Tx, 0), f, Snapshot{}}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, Snapshot{}, err
		}

		var tx Tx
		json.Unmarshal(scanner.Bytes(), &tx)

		if err := state.apply(tx); err != nil {
			return nil, Snapshot{}, err
		}
	}

	err = state.doSnapshot()
	if err != nil {
		return nil, Snapshot{}, err
	}

	return state, state.snapshot, nil
}
