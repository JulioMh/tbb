package database

type Account string

type Balances map[Account]uint

func NewAccount(account string) Account {
	return Account(account)
}

func NewBalances(account Account, value uint) Balances {
	b := make(map[Account]uint)
	b[account] = value
	return b
}
