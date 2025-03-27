package dto

type InsertTransaction struct {
	TxID      string
	BEEF      []byte
	CreatedAt string
	UpdatedAt string
}

type UpdateTransactionBEEF struct {
	TxID string
	BEEF []byte
}
