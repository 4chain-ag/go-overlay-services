CREATE TABLE transactions(
    txid BYTEA PRIMARY KEY,
    beef BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE outputs(
    topic TEXT NOT NULL,
    txid TEXT NOT NULL,
    vout INTEGER NOT NULL,
    height INTEGER,
    idx BIGINT NOT NULL DEFAULT 0,
    satoshis BIGINT NOT NULL,
    script BYTEA NOT NULL,
    consumes BYTEA[],
    consumed_by BYTEA[],
    spent BOOL NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(topic, txid, vout)
);
CREATE INDEX idx_outputs_txid_vout ON outputs(txid, vout);
CREATE INDEX idx_outputs_topic_height_idx ON outputs(topic, height, idx);

CREATE TABLE applied_transactions(
    topic TEXT NOT NULL,
    txid BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(topic, txid)
);