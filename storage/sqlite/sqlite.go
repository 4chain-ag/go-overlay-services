package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/4chain-ag/go-overlay-services/storage"
	"github.com/4chain-ag/go-overlay-services/types"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

type SQLiteStorage struct {
	DB *sql.DB
}

func NewSQLiteStorage(ctx context.Context, conn string) (*SQLiteStorage, error) {
	if db, err := sql.Open("sqlite3", conn); err != nil {
		return nil, err
	} else {
		return &SQLiteStorage{DB: db}, nil
	}
}

func (s *SQLiteStorage) InsertOutput(ctx context.Context, utxo *types.Output) (err error) {
	consumed := []byte("[]")
	if len(utxo.OutputsConsumed) > 0 {
		if consumed, err = json.Marshal(utxo.OutputsConsumed); err != nil {
			return
		}
	}
	if _, err = s.DB.ExecContext(ctx, `
        INSERT INTO outputs(topic, txid, vout, height, idx, satoshis, script, spent, consumes)
        VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(topic, txid, vout) DO NOTHING`,
		utxo.Topic,
		utxo.Outpoint.Txid.String(),
		utxo.Outpoint.OutputIndex,
		utxo.BlockHeight,
		utxo.BlockIdx,
		utxo.Satoshis,
		utxo.Script,
		utxo.Spent,
		consumed,
	); err != nil {
		return err
	} else if len(utxo.Beef) > 0 {
		if _, err = s.DB.ExecContext(ctx, `
			INSERT INTO transactions(txid, beef)
			VALUES(?1, ?2)
			ON CONFLICT(txid) DO NOTHING`,
			utxo.Outpoint.Txid.String(),
			utxo.Beef,
		); err != nil {
			return err
		}
	}
	return
}

func (s *SQLiteStorage) FindOutput(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (output *types.Output, err error) {
	output = &types.Output{
		Outpoint: outpoint,
	}
	query := `SELECT topic, height, idx, satoshis, script, spent, consumes, consumed_by
        FROM outputs
        WHERE txid = ? AND vout = ? `
	args := []interface{}{outpoint.Txid.String(), outpoint.OutputIndex}
	if topic != nil {
		query += "AND topic = ? "
		args = append(args, *topic)
	}
	if spent != nil {
		query += "AND spent = ? "
		args = append(args, *spent)
	}
	var consumes []byte
	var consumedBy []byte
	if err := s.DB.QueryRowContext(ctx, query, args).Scan(
		&output.Topic,
		&output.BlockHeight,
		&output.BlockIdx,
		&output.Satoshis,
		&output.Script,
		&output.Spent,
		&consumes,
		&consumedBy,
	); err == sql.ErrNoRows {
		return nil, storage.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(consumes, &output.OutputsConsumed); err != nil {
		return nil, err
	} else if err := json.Unmarshal(consumedBy, &output.ConsumedBy); err != nil {
		return nil, err
	} else if includeBEEF {
		if err := s.DB.QueryRowContext(ctx, `
			SELECT beef
			FROM transactions
			WHERE txid = ?`,
			outpoint.Txid.String(),
		).Scan(&output.Beef); err != nil {
			return nil, err
		}
	}
	return
}

func (s *SQLiteStorage) FindOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) ([]*types.Output, error) {
	var outputs []*types.Output
	if len(outpoints) == 0 {
		return nil, nil
	}

	args := []interface{}{}
	outpointRows := make([]string, len(outpoints))
	for i, outpoint := range outpoints {
		outpointRows[i] = "(?,?)"
		args = append(args, outpoint.Txid.String(), outpoint.OutputIndex)
	}
	var query strings.Builder
	query.WriteString(`SELECT topic, txid, vout, height, idx, satoshis, script, spent, consumes, consumed_by
        FROM outputs
        WHERE (txid, vout) IN (` +
		strings.Join(outpointRows, ",") +
		") ")
	if topic != nil {
		query.WriteString("AND topic = ? ")
		args = append(args, *topic)
	}
	if spent != nil {
		query.WriteString("AND spent = ? ")
		args = append(args, *spent)
	}
	query.WriteString("LIMIT 1")
	rows, err := s.DB.QueryContext(ctx, query.String(), topic)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		output := &types.Output{
			Outpoint: &overlay.Outpoint{},
		}
		var txid string
		var consumes []byte
		var consumedBy []byte
		if err := rows.Scan(&txid,
			&output.Topic,
			&output.Outpoint.OutputIndex,
			&output.BlockHeight,
			&output.BlockIdx,
			&output.Satoshis,
			&output.Script,
			&output.Spent,
			&consumes,
			&consumedBy,
		); err != nil {
			return nil, err
		} else if output.Outpoint.Txid, err = chainhash.NewHashFromHex(txid); err != nil {
			return nil, err
		} else if err := json.Unmarshal(consumes, &output.OutputsConsumed); err != nil {
			return nil, err
		} else if err := json.Unmarshal(consumedBy, &output.ConsumedBy); err != nil {
			return nil, err
		}
		if includeBEEF {
			if err := s.DB.QueryRowContext(ctx, `
				SELECT beef
				FROM transactions
				WHERE txid = ?`,
				txid,
			).Scan(&output.Beef); err != nil {
				return nil, err
			}
		}
		outputs = append(outputs, output)

	}
	return outputs, nil
}

func (s *SQLiteStorage) FindOutputsForTransaction(ctx context.Context, txid *chainhash.Hash, includeBEEF bool) ([]*types.Output, error) {
	var outputs []*types.Output
	query := `SELECT topic, vout, height, idx, satoshis, script, spent, consumes, consumed_by
        FROM outputs
        WHERE txid = ?
        ORDER BY vout ASC`
	rows, err := s.DB.QueryContext(ctx, query, txid.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		output := &types.Output{
			Outpoint: &overlay.Outpoint{Txid: txid},
		}
		var consumes []byte
		var consumedBy []byte
		if err := rows.Scan(
			&output.Topic,
			&output.Outpoint.OutputIndex,
			&output.BlockHeight,
			&output.BlockIdx,
			&output.Satoshis,
			&output.Script,
			&output.Spent,
			&consumes,
			&consumedBy,
		); err != nil {
			return nil, err
		} else if err := json.Unmarshal(consumes, &output.OutputsConsumed); err != nil {
			return nil, err
		} else if err := json.Unmarshal(consumedBy, &output.ConsumedBy); err != nil {
			return nil, err
		} else if includeBEEF {
			if err := s.DB.QueryRowContext(ctx, `
				SELECT beef
				FROM transactions
				WHERE txid = ?`,
				txid.String(),
			).Scan(&output.Beef); err != nil {
				return nil, err
			}
		}
		outputs = append(outputs, output)
	}
	return outputs, nil
}

func (s *SQLiteStorage) FindUTXOsForTopic(ctx context.Context, topic string, since float64, includeBEEF bool) ([]*types.Output, error) {
	var outputs []*types.Output
	query := `
        SELECT txid, vout, height, idx, satoshis, script, spent, consumes, consumed_by
        FROM outputs
        WHERE topic = ? AND height >= ?
        ORDER BY height ASC, idx ASC`
	rows, err := s.DB.QueryContext(ctx, query, topic, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		output := &types.Output{
			Outpoint: &overlay.Outpoint{},
			Topic:    topic,
		}
		var txid string
		var consumes []byte
		var consumedBy []byte
		if err := rows.Scan(
			&txid,
			&output.Outpoint.OutputIndex,
			&output.BlockHeight,
			&output.BlockIdx,
			&output.Satoshis,
			&output.Script,
			&output.Spent,
			&consumes,
			&consumedBy,
		); err != nil {
			return nil, err
		} else if output.Outpoint.Txid, err = chainhash.NewHashFromHex(txid); err != nil {
			return nil, err
		} else if err := json.Unmarshal(consumes, &output.OutputsConsumed); err != nil {
			return nil, err
		} else if err := json.Unmarshal(consumedBy, &output.ConsumedBy); err != nil {
			return nil, err
		}
		if includeBEEF {
			if err := s.DB.QueryRowContext(ctx, `
				SELECT beef
				FROM transactions
				WHERE txid = ?`,
				txid,
			).Scan(&output.Beef); err != nil {
				return nil, err
			}
		}
		outputs = append(outputs, output)
	}
	return outputs, nil
}

func (s *SQLiteStorage) DeleteOutput(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	_, err := s.DB.ExecContext(ctx, `
        DELETE FROM outputs
        WHERE topic = ? AND txid = ? AND vout = ?`,
		topic,
		outpoint.Txid,
		outpoint.OutputIndex,
	)
	return err
}

func (s *SQLiteStorage) DeleteOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic string) error {
	query := `
        DELETE FROM outputs
        WHERE topic = ? AND (txid, vout) IN (`
	for i, outpoint := range outpoints {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("('%s', %d)", outpoint.Txid, outpoint.OutputIndex)
	}
	query += ")"
	_, err := s.DB.ExecContext(ctx, query, topic)
	return err
}

func (s *SQLiteStorage) MarkUTXOAsSpent(ctx context.Context, outpoint *overlay.Outpoint, topic string, spendTxid *chainhash.Hash) error {
	_, err := s.DB.ExecContext(ctx, `
        UPDATE outputs
        SET spent = ?
        WHERE topic = ? AND txid = ? AND vout = ?`,
		spendTxid.String(),
		topic,
		outpoint.Txid,
		outpoint.OutputIndex,
	)
	return err
}

func (s *SQLiteStorage) MarkUTXOsAsSpent(ctx context.Context, outpoints []*overlay.Outpoint, topic string, spendTxid *chainhash.Hash) error {
	query := `
        UPDATE outputs
        SET spent = ?
        WHERE topic = ? AND (txid, vout) IN (`
	for i, outpoint := range outpoints {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("('%s', %d)", outpoint.Txid, outpoint.OutputIndex)
	}
	query += ")"
	_, err := s.DB.ExecContext(ctx, query, spendTxid.String(), topic)
	return err
}

func (s *SQLiteStorage) UpdateConsumedBy(ctx context.Context, outpoint *overlay.Outpoint, topic string, consumedBy []*overlay.Outpoint) error {
	consumedByStr := make([]string, len(consumedBy))
	for i, op := range consumedBy {
		consumedByStr[i] = fmt.Sprintf("%s:%d", op.Txid, op.OutputIndex)
	}
	_, err := s.DB.ExecContext(ctx, `
        UPDATE outputs
        SET consumed_by = ?
        WHERE topic = ? AND txid = ? AND vout = ?`,
		consumedByStr,
		topic,
		outpoint.Txid,
		outpoint.OutputIndex,
	)
	return err
}

func (s *SQLiteStorage) UpdateTransactionBEEF(ctx context.Context, txid *chainhash.Hash, beef []byte) error {
	_, err := s.DB.ExecContext(ctx, `
        UPDATE transactions
        SET beef = ?
        WHERE txid = ?`,
		beef,
		txid.String(),
	)
	return err
}

func (s *SQLiteStorage) UpdateOutputBlockHeight(ctx context.Context, outpoint *overlay.Outpoint, topic string, blockHeight uint32, blockIndex uint64) error {
	_, err := s.DB.ExecContext(ctx, `
        UPDATE outputs
        SET height = ?, idx = ?
        WHERE topic = ? AND txid = ? AND vout = ?`,
		blockHeight,
		blockIndex,
		topic,
		outpoint.Txid,
		outpoint.OutputIndex,
	)
	return err
}

func (s *SQLiteStorage) InsertAppliedTransaction(ctx context.Context, tx *overlay.AppliedTransaction) error {
	_, err := s.DB.ExecContext(ctx, `
        INSERT INTO applied_transactions(topic, txid)
        VALUES(?, ?)
        ON CONFLICT(topic, txid) DO NOTHING`,
		tx.Topic,
		tx.Txid,
	)
	return err
}

func (s *SQLiteStorage) DoesAppliedTransactionExist(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
	var exists bool
	err := s.DB.QueryRowContext(ctx, `
        SELECT EXISTS(SELECT 1 FROM applied_transactions WHERE topic = ? AND txid = ?)`,
		tx.Topic,
		tx.Txid,
	).Scan(&exists)
	return exists, err
}
