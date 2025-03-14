CREATE TABLE transactions(
    txid BYTEA PRIMARY KEY,
    beef BYTEA NOT NULL
);

CREATE TABLE outputs(
    topic TEXT NOT NULL,
    outpoint BYTEA,
    txid BYTEA GENERATED ALWAYS AS (substring(outpoint from 1 for 32)) STORED,
    vout INTEGER GENERATED ALWAYS AS ((get_byte(outpoint, 32) << 24) |
       (get_byte(outpoint, 33) << 16) |
	   (get_byte(outpoint, 34) << 8) |
       (get_byte(outpoint, 35))) STORED,
    height INTEGER NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    idx BIGINT NOT NULL DEFAULT 0,
    satoshis BIGINT NOT NULL,
    script BYTEA NOT NULL,
    consumes BYTEA[],
    consumed_by BYTEA[],
    spent BYTEA NOT NULL DEFAULT '\x',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(outpoint, topic)
);
CREATE INDEX idx_outputs_txid_vout ON outputs(txid, vout);
CREATE INDEX idx_outputs_topic_height_idx ON outputs(topic, height, idx);

CREATE TABLE applied_transactions(
    txid BYTEA,
    topic TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(txid, topic)
);