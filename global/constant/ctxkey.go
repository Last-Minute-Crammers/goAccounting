package constant

type Key string

const (
	// about JWT
	Claims Key = "claims"
	UserId Key = "userId"
	User   Key = "user"
	Db     Key = "Db"
	// Tx means Transaction
	Tx       Key = "Tx"
	TxCommit Key = "TxCommit"
)
