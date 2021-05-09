package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"tbb_blockchain/database"
)

const httpPort int = 8080

type BalancesRes struct {
	Hash     database.Hash             `json:"hash"`
	Balances map[database.Account]uint `json:"balances"`
}

type TxAddReq struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint   `json:"value"`
	Data  string `json:"data"`
}

func writeRes(w http.ResponseWriter, balances BalancesRes) {
	json.NewEncoder(w).Encode(balances)
}

func writeErrRes(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

func listBalancesHanlder(w http.ResponseWriter, r *http.Request, state *database.State) {
	writeRes(w, BalancesRes{state.LatestBlockHash(), state.Balances})
}

func addTxHandler(w http.ResponseWriter, r *http.Request, state *database.State) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrRes(w, http.StatusBadRequest, err)
		return
	}

	var txAddReq TxAddReq
	json.Unmarshal(body, &txAddReq)

	tx := database.NewTx(
		database.Account(txAddReq.From),
		database.Account(txAddReq.To),
		txAddReq.Value,
		txAddReq.Data,
	)

	if err := state.AddTx(tx); err != nil {
		writeErrRes(w, http.StatusInternalServerError, err)
		return
	}

	hash, err := state.Persist()
	if err != nil {
		writeErrRes(w, http.StatusInternalServerError, err)
		return
	}

	writeRes(w, BalancesRes{hash, state.Balances})
}

func Run(dataDir string) error {
	state, _, err := database.NewStateFromDisk(dataDir)
	if err != nil {
		return err
	}
	defer state.Close()

	http.HandleFunc("/balances/list", func(rw http.ResponseWriter, r *http.Request) {
		listBalancesHanlder(rw, r, state)
	})
	http.HandleFunc("/tx/add", func(rw http.ResponseWriter, r *http.Request) {
		addTxHandler(rw, r, state)
	})
	http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)

	return nil
}
