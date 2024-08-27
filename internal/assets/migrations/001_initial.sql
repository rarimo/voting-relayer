-- +migrate Up

CREATE TABLE state (
    id BIGSERIAL PRIMARY KEY,
    operation_id CHAR(64) NOT NULL ,
    tx_hash CHAR(64) NOT NULL,
    proof CHAR(64) NOT NULL,
    root CHAR(64) NOT NULL,
    chain_id INTEGER,
    event TEXT
);

-- +migrate Down
DROP TABLE state;
