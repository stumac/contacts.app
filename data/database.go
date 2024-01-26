package data

import (
	"database/sql"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

func NewDB(dbName string) (*sql.DB, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, dbName)
	return sqldb, err
}

func NewBun(sqldb *sql.DB) (*bun.DB, error) {
	// orm := bun.NewDB(db, sqlitedialect.New())
	db := bun.NewDB(sqldb, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))
	return db, nil
}
