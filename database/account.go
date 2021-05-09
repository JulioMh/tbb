package database

type Account string

func NewAccount(account string) Account {
	return Account(account)
}
