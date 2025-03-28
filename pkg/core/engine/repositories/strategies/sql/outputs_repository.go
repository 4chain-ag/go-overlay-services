package sql

import (
	"context"
	"fmt"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine/repositories/dto"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/script"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OutputsRepository struct {
	db *gorm.DB
}

func (o *OutputsRepository) Close() error {
	db, err := o.db.DB()
	if err != nil {
		return fmt.Errorf("failed to return db object: %w", err)
	}
	return db.Close()
}

func (o *OutputsRepository) InsertOutput(ctx context.Context, command dto.InsertOutput) error {
	conflicts := clause.OnConflict{
		Columns:   []clause.Column{{Name: "outpoint"}, {Name: "topic"}},
		DoNothing: true,
	}

	err := o.db.
		WithContext(ctx).
		Create(CreateOutputEntity(command)).
		Clauses(conflicts).
		Error
	if err != nil {
		return fmt.Errorf("insert op failed: %w", err)
	}
	return nil
}

func (o *OutputsRepository) FindOutput(ctx context.Context, query dto.FindOutput) (*dto.Output, error) {
	var res Output

	err := o.db.WithContext(ctx).
		Where("txid = ?", query.TxID).
		Where("vout = ?", query.OutputIndex).
		Where("topic = ?", query.Topic).
		Where("spent = ?", query.Spent).
		First(&res).
		Error
	if err != nil {
		return nil, fmt.Errorf("query op failed: %w", err)
	}

	return &dto.Output{
		Outpoint:    overlay.Outpoint{OutputIndex: res.Vout},
		Topic:       res.Topic,
		Script:      &script.Script{},
		Satoshis:    res.Satoshis,
		Spent:       res.Spent,
		BlockHeight: res.BlockHeight,
		BlockIdx:    res.BlockIdx,
	}, nil
}

func (*OutputsRepository) FindOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) ([]*dto.Output, error) {
	return nil, nil
}

func (*OutputsRepository) DeleteOutput(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	return nil
}

func (*OutputsRepository) DeleteOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic string) error {
	return nil
}

func (*OutputsRepository) FindUTXOsForTopic(ctx context.Context, topic string, since float64, includeBEEF bool) ([]*dto.Output, error) {
	return nil, nil
}

func (*OutputsRepository) MarkUTXOAsSpent(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	return nil
}

func (*OutputsRepository) MarkUTXOsAsSpent(ctx context.Context, outpoints []*overlay.Outpoint, topic string) error {
	return nil
}

func (*OutputsRepository) UpdateConsumedBy(ctx context.Context, outpoint *overlay.Outpoint, topic string, consumedBy []*overlay.Outpoint) error {
	return nil
}

func (*OutputsRepository) UpdateOutputBlockHeight(ctx context.Context, outpoint *overlay.Outpoint, topic string, blockHeight uint32, blockIndex uint64) error {
	return nil
}

func NewOutputsPostgresRepository() *OutputsRepository {
	dsn := ""
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(Output{}, Transaction{}, AppliedTransaction{})
	if err != nil {
		panic(err)
	}
	return &OutputsRepository{db: db}
}

func CreateOutputEntity(data dto.InsertOutput) *Output {
	return &Output{
		TxID:        data.TxID,
		Vout:        data.Vout,
		Topic:       data.Topic,
		BlockHeight: data.BlockHeight,
		BlockIdx:    data.BlockIdx,
		Satoshis:    data.Satoshis,
		Script:      data.Script,
		Spent:       data.Spent,
	}
}
