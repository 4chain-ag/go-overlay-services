package sql

import (
	"context"
	"fmt"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine/repositories/dto"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionsRepository struct {
	db *gorm.DB
}

func (t *TransactionsRepository) InsertTransaction(ctx context.Context, command dto.InsertTransaction) error {
	conflicts := clause.OnConflict{
		Columns:   []clause.Column{{Name: "beef"}},
		DoNothing: true,
	}

	err := t.db.
		WithContext(ctx).
		Create(CreateTransactionEntity(command)).
		Clauses(conflicts).
		Error
	if err != nil {
		return fmt.Errorf("insert transaction entity op failed: %w", err)
	}
	return nil
}

func (t *TransactionsRepository) UpdateTransactionBEEF(ctx context.Context, command dto.UpdateTransactionBEEF) error {
	err := t.db.
		WithContext(ctx).
		Model(&Transaction{}).
		Where("txid = ?").
		Update("beef", command.BEEF).
		Error
	if err != nil {
		return fmt.Errorf("update transaction entity op failed: %w", err)
	}

	return nil
}

func (*TransactionsRepository) InsertAppliedTransaction(ctx context.Context, tx *overlay.AppliedTransaction) error {
	return nil
}

func (*TransactionsRepository) FindOutputsForTransaction(ctx context.Context, txid *chainhash.Hash, includeBEEF bool) ([]*dto.Output, error) {
	return nil, nil
}

func (*TransactionsRepository) DoesAppliedTransactionExist(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
	return false, nil
}

func NewTransactionsRepository(db *gorm.DB) *TransactionsRepository {
	return &TransactionsRepository{db: db}
}

func CreateTransactionEntity(data dto.InsertTransaction) *Transaction {
	return &Transaction{
		ID:        data.TxID,
		BEEF:      data.BEEF,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}
