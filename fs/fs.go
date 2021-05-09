package fs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

func GetDatabaseDirPath(dataDir string) string {
	return filepath.Join(dataDir, "database")
}

func GetGenesisJsonFilePath(dataDir string) string {
	return filepath.Join(GetDatabaseDirPath(dataDir), "genesis.json")
}

func GetBlocksDbJsonFilePath(dataDir string) string {
	return filepath.Join(GetDatabaseDirPath(dataDir), "block.db")
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func writeGenesisToDisk(path string) error {
	initBalances := make(map[string]uint)
	initBalances["andrej"] = 1000000

	genesis := make(map[string]interface{})
	genesis["genesis_time"] = time.Now()
	genesis["chain_id"] = "the-blockchain-bar-ledger"
	genesis["balances"] = initBalances

	jsonGenesis, err := json.Marshal(genesis)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, jsonGenesis, 0644); err != nil {
		return err
	}
	return nil
}

func writeEmptyBlocksDbToDisk(path string) error {
	if _, err := os.Create(path); err != nil {
		return err
	}
	return nil
}

func InitDataDirIfNotExists(dataDir string) error {
	databasePath := GetDatabaseDirPath(dataDir)
	if !pathExists(databasePath) {
		err := os.Mkdir(databasePath, 0700)
		if err != nil {
			return err
		}
	}
	genesisPath := GetGenesisJsonFilePath(dataDir)
	if !pathExists(genesisPath) {
		if err := writeGenesisToDisk(genesisPath); err != nil {
			return err
		}
	}
	blockDbPath := GetBlocksDbJsonFilePath(dataDir)
	if !pathExists(blockDbPath) {
		writeEmptyBlocksDbToDisk(blockDbPath)
	}

	return nil
}
