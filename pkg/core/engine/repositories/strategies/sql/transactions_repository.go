package sql

import (
	"context"
	"fmt"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine/repositories/dto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionsRepository struct {
	db *gorm.DB
}

func (t *TransactionsRepository) Close() error {
	db, err := t.db.DB()
	if err != nil {
		return fmt.Errorf("failed to return db object: %w", err)
	}
	return db.Close()
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
		return fmt.Errorf("insert op failed: %w", err)
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
		return fmt.Errorf("update op failed: %w", err)
	}
	return nil
}

func (t *TransactionsRepository) InsertAppliedTransaction(ctx context.Context, command dto.InsertAppliedTransaction) error {
	conflicts := clause.OnConflict{
		Columns:   []clause.Column{{Name: "txid"}, {Name: "topic"}},
		DoNothing: true,
	}

	err := t.db.
		WithContext(ctx).
		Create(CreateAppliedTransaction(command)).
		Clauses(conflicts).
		Error
	if err != nil {
		return fmt.Errorf("insert op failed: %w", err)
	}
	return nil
}

func (t *TransactionsRepository) FindOutputsForTransaction(ctx context.Context, data dto.FindTransactionOutput) ([]*dto.Output, error) {
	return nil, nil
}

func (t *TransactionsRepository) DoesAppliedTransactionExist(ctx context.Context, txID string) (bool, error) {
	var exists bool
	err := t.db.
		WithContext(ctx).
		Model(&AppliedTransaction{}).
		Select("exists (select 1 from applied_transactions where txid = ?)", txID).
		Scan(&exists).Error
	if err != nil {
		return false, fmt.Errorf("query op failed: %w", err)
	}
	return exists, err
}

func NewTransactionsPostgresRepository() *TransactionsRepository {
	dsn := ""
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(Output{}, Transaction{}, AppliedTransaction{})
	if err != nil {
		panic(err)
	}

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

func CreateAppliedTransaction(data dto.InsertAppliedTransaction) *AppliedTransaction {
	return &AppliedTransaction{
		TxID:  data.TxID,
		Topic: data.Topic,
	}
}
