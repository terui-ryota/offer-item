package txhelper

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

// トランザクションを管理する
func WithTransaction(ctx context.Context, db *sql.DB, execute func(*sql.Tx) error) (err error) {
	boil.SetDB(db)
	tx, err := boil.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = execute(tx)
	return
}
