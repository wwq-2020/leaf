package generator

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/wwq1988/leaf/pkg/conf"
)

const (
	defaultStep uint64 = 1000
)

type segment struct {
	db *sql.DB
}

func NewSegment() (Generator, error) {
	dsn := conf.GetDsn()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.Wrapf(err, "dsn", dsn)
	}
	return &segment{db: db}, nil
}

func (s *segment) Generate(key string) ([]uint64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer tx.Rollback()
	var maxId uint64
	var step uint64
	if err := tx.QueryRow("select max_id, step  from leaf_alloc where biz_tag = ? for update", key).Scan(&maxId, &step); err != nil {
		return nil, errors.WithStack(err)
	}

	if step == 0 {
		step = defaultStep
	}

	ids := make([]uint64, 0, step)
	for i := uint64(1); i < step; i++ {
		ids = append(ids, maxId+i)
	}
	if _, err := tx.Exec("update leaf_alloc set max_id = ? where biz_tag = ?", maxId+step, key); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.WithStack(err)
	}

	return ids, nil
}

func (s *segment) Stop() {
	s.db.Close()
}
