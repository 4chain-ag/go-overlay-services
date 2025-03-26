package sql

import (
	"context"

	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"gorm.io/gorm"
)

type TransactionsRepository struct {
	db *gorm.DB
}

func (*TransactionsRepository) UpdateTransactionBEEF(ctx context.Context, txid *chainhash.Hash, beef []byte) error {
	return nil
}

func (*TransactionsRepository) InsertAppliedTransaction(ctx context.Context, tx *overlay.AppliedTransaction) error {
	return nil
}

func (*TransactionsRepository) DoesAppliedTransactionExist(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
	return false, nil
}
