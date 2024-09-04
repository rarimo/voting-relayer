-- +migrate Up

CREATE TABLE state (
    id BIGSERIAL PRIMARY KEY,
    operation_id TEXT NOT NULL ,
    tx_hash TEXT NOT NULL,
    root TEXT NOT NULL,
    proof TEXT NOT NULL,
    chain_id INTEGER,
    event TEXT,
    block_height INT
);

-- +migrate Down
DROP TABLE state;