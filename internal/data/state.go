package data

type StateQ interface {
	New() StateQ

	Get() (*State, error)
	Insert(data State) (int64, error)
	FilterByOperationId(ids ...string) StateQ
}

type State struct {
	ID          int64    `db:"id" structs:"-"`
	OperationId [32]byte `db:"operation_id" structs:"operation_id"`
	TxHash      [32]byte `db:"tx_hash" structs:"tx_hash"`
	Proof       [32]byte `db:"proof" structs:"proof"`
	Root        [32]byte `db:"root" structs:"root"`
	ChainId     int64    `db:"chain_id" structs:"chain_id"`

	Event string `db:"event" structs:"event"`
}
