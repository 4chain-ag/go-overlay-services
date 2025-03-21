CREATE TABLE transactions(
    txid TEXT PRIMARY KEY,
    beef BLOB NOT NULL,
    created_at TEXT NOT NULL DEFAULT current_timestamp,
    updated_at TEXT NOT NULL DEFAULT current_timestamp
);

CREATE TABLE outputs(
    txid TEXT NOT NULL,
    vout INTEGER NOT NULL,
    topic TEXT NOT NULL,
    height INTEGER,
    idx BIGINT NOT NULL DEFAULT 0,
    satoshis BIGINT NOT NULL,
    script BLOB NOT NULL,
    consumes TEXT NOT NULL DEFAULT json('[]'),
    consumed_by TEXT NOT NULL DEFAULT json('[]'),
    spent BOOL NOT NULL DEFAULT false,
    created_at TEXT NOT NULL DEFAULT current_timestamp,
    updated_at TEXT NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY(txid, vout, topic)
);
CREATE INDEX idx_outputs_txid_vout ON outputs(txid, vout);
CREATE INDEX idx_outputs_topic_height_idx ON outputs(topic, height, idx);

CREATE TABLE applied_transactions(
    txid TEXT NOT NULL,
    topic TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT current_timestamp,
    updated_at TEXT NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY(txid, topic)
);