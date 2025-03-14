package postgres

import (
	"context"
	"strings"

	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shruggr/goverlay/storage"
	// "github.com/shruggr/goverlay/types"
)

type PostgresStorage struct {
	DB *pgxpool.Pool
}

func NewPostgresStorage(ctx context.Context, conn string) (*PostgresStorage, error) {
	if db, err := pgxpool.New(ctx, conn); err != nil {
		return nil, err
	} else {
		return &PostgresStorage{DB: db}, nil
	}
}

func (s *PostgresStorage) InsertOutput(ctx context.Context, utxo *overlay.Output) error {
	_, err := s.DB.Exec(ctx, `
		INSERT INTO outputs(outpoint, topic, height, idx, satoshis, script, spent)
		VALUES($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT(outpoint, topic) DO NOTHING`,
		utxo.Outpoint,
		utxo.Topic,
		utxo.BlockHeight,
		utxo.BlockIdx,
		utxo.Satoshis,
		utxo.Script,
		utxo.Spent,
	)
	return err
}

func (s *PostgresStorage) FindOutput(ctx context.Context, outpoint *overlay.Outpoint, topic string, spent bool, includeBEEF bool) (*overlay.Output, error) {
	var sql strings.Builder
	if includeBEEF {
		sql.WriteString("SELECT o.outpoint, o.topic, o.height, o.idx, o.satoshis, o.script, o.spent, t.beef ")
		sql.WriteString("FROM outputs o ")
		sql.WriteString("JOIN transactionst t ON o.txid = t.txid ")
	} else {
		sql.WriteString("SELECT outpoint, topic, height, idx, satoshis, script, spent, '\\x' AS beef ")
		sql.WriteString("FROM outputs ")
	}
	sql.WriteString("WHERE outpoint = $1 AND topic = $2 ")
	if spent {
		sql.WriteString("AND spent != '\\x' ")
	}

	var utxo overlay.Output
	row := s.DB.QueryRow(ctx, sql.String(), outpoint, topic)
	err := row.Scan(&utxo.Outpoint, &utxo.Topic, &utxo.BlockHeight, &utxo.BlockIdx, &utxo.Satoshis, &utxo.Script, &utxo.Spent)
	if err == pgx.ErrNoRows {
		err = storage.ErrNotFound
	}

	return &utxo, err
}

func (s *PostgresStorage) FindOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic string, spent bool, includeBEEF bool) ([]*overlay.Output, error) {
	var sql strings.Builder
	if includeBEEF {
		sql.WriteString("SELECT o.outpoint, o.topic, o.height, o.idx, o.satoshis, o.script, o.spent, t.beef ")
		sql.WriteString("FROM outputs o ")
		sql.WriteString("JOIN transactionst t ON o.txid = t.txid ")
	} else {
		sql.WriteString("SELECT outpoint, topic, height, idx, satoshis, script, spent, '\\x' AS beef ")
		sql.WriteString("FROM outputs ")
	}
	sql.WriteString("WHERE outpoint = $1 AND topic = $2 ")
	if spent {
		sql.WriteString("AND spent != '\\x' ")
	}
	sql.WriteString("ORDER BY height ASC, idx ASC ")

	if rows, err := s.DB.Query(ctx, sql.String(), outpoints, topic); err != nil {
		return nil, err
	} else {
		defer rows.Close()

		var utxosByOutpoint = make(map[string]*overlay.Output)
		// var utxos []*overlay.Output
		for rows.Next() {
			var utxo overlay.Output
			if err := rows.Scan(&utxo.Outpoint, &utxo.Topic, &utxo.BlockHeight, &utxo.BlockIdx, &utxo.Satoshis, &utxo.Script, &utxo.Spent); err != nil {
				return nil, err
			}
			utxosByOutpoint[utxo.Outpoint.String()] = &utxo
		}

		utxos := make([]*overlay.Output, len(outpoints))
		for i, outpoint := range outpoints {
			utxos[i] = utxosByOutpoint[outpoint.String()]
		}
		return utxos, nil
	}
}

func (s *PostgresStorage) FindOutputsForTransaction(ctx context.Context, txid *chainhash.Hash, includeBEEF bool) ([]*overlay.Output, error) {
	var sql strings.Builder
	if includeBEEF {
		sql.WriteString("SELECT o.outpoint, o.topic, o.height, o.idx, o.satoshis, o.script, o.spent, t.beef ")
		sql.WriteString("FROM outputs o ")
		sql.WriteString("JOIN transactionst t ON o.txid = t.txid ")
	} else {
		sql.WriteString("SELECT outpoint, topic, height, idx, satoshis, script, spent, '\\x' AS beef ")
		sql.WriteString("FROM outputs ")
	}
	sql.WriteString("WHERE txid = $1 ")
	sql.WriteString("ORDER BY vout ASC ")

	if rows, err := s.DB.Query(ctx, sql.String(), txid); err != nil {
		return nil, err
	} else {
		defer rows.Close()

		var utxos []*overlay.Output
		for rows.Next() {
			var utxo overlay.Output
			if err := rows.Scan(&utxo.Outpoint, &utxo.Topic, &utxo.BlockHeight, &utxo.BlockIdx, &utxo.Satoshis, &utxo.Script, &utxo.Spent); err != nil {
				return nil, err
			}
			utxos = append(utxos, &utxo)
		}
		return utxos, nil
	}
}

func (s *PostgresStorage) FindUTXOsForTopic(ctx context.Context, topic string, since float64, includeBEEF bool) ([]*overlay.Output, error) {
	var sql strings.Builder
	if includeBEEF {
		sql.WriteString("SELECT o.outpoint, o.topic, o.height, o.idx, o.satoshis, o.script, o.spent, t.beef ")
		sql.WriteString("FROM outputs o ")
		sql.WriteString("JOIN transactionst t ON o.txid = t.txid ")
	} else {
		sql.WriteString("SELECT outpoint, topic, height, idx, satoshis, script, spent, '\\x' AS beef ")
		sql.WriteString("FROM outputs ")
	}
	sql.WriteString("WHERE topic = $1 AND height >= $2 ")
	sql.WriteString("ORDER BY height ASC, idx ASC ")

	if rows, err := s.DB.Query(ctx, sql.String(), topic, since); err != nil {
		return nil, err
	} else {
		defer rows.Close()

		var utxos []*overlay.Output
		for rows.Next() {
			var utxo overlay.Output
			if err := rows.Scan(&utxo.Outpoint, &utxo.Topic, &utxo.BlockHeight, &utxo.BlockIdx, &utxo.Satoshis, &utxo.Script, &utxo.Spent, &utxo.Beef); err != nil {
				return nil, err
			}
			utxos = append(utxos, &utxo)
		}
		return utxos, nil
	}
}

func (s *PostgresStorage) DeleteOutput(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	_, err := s.DB.Exec(ctx, `
		DELETE FROM outputs
		WHERE outpoint = $1 AND topic = $2`,
		outpoint,
		topic,
	)
	return err
}

func (s *PostgresStorage) DeleteOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic string) error {
	_, err := s.DB.Exec(ctx, `
		DELETE FROM outputs
		WHERE outpoint = ANY($1) AND topic = $2`,
		outpoints,
		topic,
	)
	return err
}

func (s *PostgresStorage) MarkUTXOAsSpent(ctx context.Context, outpoint *overlay.Outpoint, topic string, spendTxid *chainhash.Hash) error {
	_, err := s.DB.Exec(ctx, `
		UPDATE outputs
		SET spent = $1
		WHERE outpoint = $2 AND topic = $3`,
		spendTxid,
		outpoint,
		topic,
	)
	return err
}

func (s *PostgresStorage) MarkUTXOsAsSpent(ctx context.Context, outpoints []*overlay.Outpoint, topic string, spendTxid *chainhash.Hash) error {
	_, err := s.DB.Exec(ctx, `
		UPDATE outputs
		SET spent = $1
		WHERE outpoint = ANY($2) AND topic = $3`,
		spendTxid,
		outpoints,
		topic,
	)
	return err
}

func (s *PostgresStorage) UpdateConsumedBy(ctx context.Context, outpoint *overlay.Outpoint, topic string, consumedBy []*overlay.Outpoint) error {
	_, err := s.DB.Exec(ctx, `
		UPDATE outputs
		SET consumed_by = $1
		WHERE outpoint = $2 AND topic = $3`,
		consumedBy,
		outpoint,
		topic,
	)
	return err
}

func (s *PostgresStorage) UpdateTransactionBEEF(ctx context.Context, txid *chainhash.Hash, beef []byte) error {
	_, err := s.DB.Exec(ctx, `
		UPDATE transaction
		SET beef = $1
		WHERE txid = $2`,
		beef,
		txid,
	)
	return err
}

func (s *PostgresStorage) UpdateOutputBlockHeight(ctx context.Context, outpoint *overlay.Outpoint, topic string, blockHeight uint32, blockIndex uint64) error {
	_, err := s.DB.Exec(ctx, `
		UPDATE outputs
		SET height = $1, idx = $2
		WHERE outpoint = $3 AND topic = $4`,
		blockHeight,
		blockIndex,
		outpoint,
		topic,
	)
	return err
}

func (s *PostgresStorage) InsertAppliedTransaction(ctx context.Context, tx *overlay.AppliedTransaction) error {
	_, err := s.DB.Exec(ctx, `
		INSERT INTO applied_transactions(txid, topic)
		VALUES($1, $2)
		ON CONFLICT(txid) DO NOTHING`,
		tx.Txid,
		tx.Topic,
	)
	return err
}

func (s *PostgresStorage) DoesAppliedTransactionExist(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
	var exists bool
	err := s.DB.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM applied_transactions WHERE txid = $1)`,
		tx.Txid,
	).Scan(&exists)
	return exists, err
}
