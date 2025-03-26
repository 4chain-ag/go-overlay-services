package sql

import (
	"gorm.io/datatypes"
)

type Output struct {
	TxID        string         `gorm:"column:txid; primaryKey;not null"`
	Vout        uint32         `gorm:"column:vout;primarykey;not null"`
	Topic       string         `gorm:"column:topic;primaryKey;not null"`
	BlockHeight uint32         `gorm:"column:height"`
	BlockIdx    uint64         `gorm:"column:idx;not null;default:0"`
	Satoshis    uint64         `gorm:"column:satoshis;not null"`
	Script      []byte         `gorm:"column:script;not null"`
	Consumes    datatypes.JSON `gorm:"column:consumes;not null; default: '[]'"`
	ConsumedBy  datatypes.JSON `gorm:"column_consumed_by;not null;default:'[]'"`
	Spent       bool           `gorm:"column:spent;not null;default:false"`
	CreatedAt   string         `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   string         `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
}

func (Output) TableName() string { return "outputs" }

type Transaction struct {
	ID        string `gorm:"column:txid;primaryKey;not null"`
	BEEF      []byte `gorm:"column:beef;not null"`
	CreatedAt string `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt string `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
}

func (Transaction) TableName() string { return "transactions" }

type AppliedTransaction struct {
	TxID      string `gorm:"column:txid;primaryKey;not null"`
	Topic     string `gorm:"column:topic;primaryKey;not null"`
	CreatedAt string `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt string `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
}

func (AppliedTransaction) TableName() string { return "applied_transactions" }
