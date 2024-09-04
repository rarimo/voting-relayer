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
	ID          int64  `db:"id" structs:"-"`
	OperationId string `db:"operation_id" structs:"operation_id"`
	TxHash      string `db:"tx_hash" structs:"tx_hash"`
	Proof       string `db:"proof" structs:"proof"`
	Root        string `db:"root" structs:"root"`
	ChainId     int64  `db:"chain_id" structs:"chain_id"`
	BlockHeight uint64 `db:"block_height" structs:"block_height"`
	Event       string `db:"event" structs:"event"`
}

const (
	ASC  SortOrder = "ASC"
	DESC SortOrder = "DESC"
)
