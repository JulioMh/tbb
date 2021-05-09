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

func NewStateFromDisk() (*State, Hash, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, Hash{}, err
	}
	genFilePath := filepath.Join(cwd, "database", "disk", "genesis.json")
	gen, err := loadGenesis(genFilePath)
	if err != nil {
		return nil, Hash{}, err
	}

	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	txDbFilePath := filepath.Join(cwd, "database", "disk", "block.db")
	f, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, Hash{}, err
	}

	scanner := bufio.NewScanner(f)
	state := &State{balances, make([]Tx, 0), f, Hash{}}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, Hash{}, err
		}

		var blockFs BlockFs
		json.Unmarshal(scanner.Bytes(), &blockFs)

		for _, tx := range blockFs.Block.TXs {
			if err := state.apply(tx); err != nil {
				return nil, Hash{}, err
			}
		}

		state.prevHash = blockFs.Hash
	}

	return state, state.prevHash, nil
}
