package database

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"tbb_blockchain/fs"

	"time"
)

type Genesis struct {
	Genesis_time time.Time `json:"genesis_time"`
	Chain_id     string    `json:"chain_id"`
	Balances     Balances  `json:"balances"`
}

func loadGenesis(genFilePath string) (*State, error) {
	// Read genesis file
	f, err := ioutil.ReadFile(genFilePath)
	if err != nil {
		return nil, err
	}
	// Parse to map[string]interface{}
	var genesis Genesis
	json.Unmarshal(f, &genesis)
	// Create balances from genesis
	genesisBalances := make(map[Account]uint)
	for account, balance := range genesis.Balances {
		genesisBalances[Account(account)] = balance
	}
	// Create first state
	var genesisState State
	genesisState.Balances = genesisBalances
	return &genesisState, nil
}

func NewStateFromDisk(dataDir string) (*State, Hash, error) {
	if err := fs.InitDataDirIfNotExists(dataDir); err != nil {
		return nil, Hash{}, err
	}

	genFilePath := fs.GetGenesisJsonFilePath(dataDir)
	gen, err := loadGenesis(genFilePath)
	if err != nil {
		return nil, Hash{}, err
	}

	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	txDbFilePath := fs.GetBlocksDbJsonFilePath(dataDir)
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
