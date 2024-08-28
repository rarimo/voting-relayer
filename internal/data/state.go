package data

type SortOrder string

type StateQ interface {
	New() StateQ

	Get() (*State, error)
	Insert(data State) (int64, error)
	FilterByRoot(root ...string) StateQ
	FilterByBlockHeight(blockHeight ...string) StateQ
	SortByBlockHeight(order SortOrder) StateQ
}

type State struct {
	ID                 int64    `db:"id" structs:"-"`
	OperationId        [32]byte `db:"operation_id" structs:"operation_id"`
	TxHash             [32]byte `db:"tx_hash" structs:"tx_hash"`
	Proof              [32]byte `db:"proof" structs:"proof"`
	Root               [32]byte `db:"root" structs:"root"`
	ChainId            int64    `db:"chain_id" structs:"chain_id"`
	BlockHeight        uint64   `db:"block_height" structs:"block_height"`
	DestinationAddress [20]byte `db:"destination_address" structs:"destination_address"`
	Event              string   `db:"event" structs:"event"`
}

const (
	ASC  SortOrder = "ASC"
	DESC SortOrder = "DESC"
)
