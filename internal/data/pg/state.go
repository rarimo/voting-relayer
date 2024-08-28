package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/rarimo/voting-relayer/internal/data"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const stateTableName = "state"

func NewStateQ(db *pgdb.DB) data.StateQ {
	return &StateQ{
		db:  db.Clone(),
		sql: sq.Select("b.*").From(fmt.Sprintf("%s as b", stateTableName)),
	}
}

type StateQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

func (q *StateQ) New() data.StateQ {
	return NewStateQ(q.db)
}

func (q *StateQ) Get() (*data.State, error) {
	var result data.State
	err := q.db.Get(&result, q.sql)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &result, err
}

func (q *StateQ) Insert(value data.State) (int64, error) {
	clauses := structs.Map(value)
	var id int64

	stmt := sq.Insert(stateTableName).SetMap(clauses).Suffix("returning id")
	err := q.db.Get(&id, stmt)

	return id, err
}

func (q *StateQ) FilterByRoot(root ...string) data.StateQ {
	q.sql = q.sql.Where(sq.Eq{"b.root": root})
	return q
}

func (q *StateQ) FilterByBlockHeight(blockHeight ...string) data.StateQ {
	q.sql = q.sql.Where(sq.Eq{"b.block_height": blockHeight})
	return q
}

func (q *StateQ) SortByBlockHeight(order data.SortOrder) data.StateQ {
	q.sql = q.sql.OrderBy(fmt.Sprintf("b.block_height %s", order))
	return q
}
