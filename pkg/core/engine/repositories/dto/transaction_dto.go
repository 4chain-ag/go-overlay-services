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

type InsertAppliedTransaction struct {
	TxID  string
	Topic string
}

type FindTransactionOutput struct {
	TxID        string
	IncludeBEEF bool
}
