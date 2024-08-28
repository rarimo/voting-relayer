-- +migrate Up

CREATE TABLE state (
    id BIGSERIAL PRIMARY KEY,
    operation_id CHAR(64) NOT NULL ,
    tx_hash CHAR(64) NOT NULL,
    root CHAR(64) NOT NULL,
    proof TEXT NOT NULL,
    chain_id INTEGER,
    event TEXT,
    block_height INT
);

-- +migrate Down
DROP TABLE state;
