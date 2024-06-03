package pgx

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Connection struct {
	Host         string
	Port         int
	User         string
	Password     string
	SchemaName   string
	DatabaseName string
	SslMode      string
	Logger       *zap.Logger
}

type Pgx struct {
	Db *pgxpool.Pool
}

// NewDBConnection membuat dan mengembalikan pool koneksi ke database PostgreSQL
func NewDBConnection(conf *Connection) (*Pgx, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%v sslmode=%s", conf.User, conf.DatabaseName, conf.Password, conf.Host, conf.Port, conf.SslMode)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 5 // Sesuaikan dengan kebutuhan Anda
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &Pgx{Db: pool}, nil
}

func (conn *Pgx) QueryRow(ctx context.Context, tx pgx.Tx, sqlStatement string, args ...interface{}) (row pgx.Row) {
	if tx != nil {
		row = tx.QueryRow(ctx, sqlStatement, args...)
	} else {
		row = conn.Db.QueryRow(ctx, sqlStatement, args...)
	}

	fmt.Println("sql statement :", sqlStatement, " args", args)
	return row
}

func (conn *Pgx) QueryRows(ctx context.Context, tx pgx.Tx, sqlStatement string, args ...interface{}) (rows pgx.Rows, err error) {
	if tx != nil {
		rows, err = tx.Query(ctx, sqlStatement, args...)
	} else {
		rows, err = conn.Db.Query(ctx, sqlStatement, args...)
	}

	fmt.Println("err :", err, "sql statement :", sqlStatement, " args", args)
	return rows, err
}

func (conn *Pgx) Exec(ctx context.Context, tx pgx.Tx, sqlStatement string, args ...any) (res pgconn.CommandTag, err error) {
	if tx != nil {
		res, err = tx.Exec(ctx, sqlStatement, args...)
	} else {

		res, err = conn.Db.Exec(ctx, sqlStatement, args...)
	}

	fmt.Println("err :", err, "sql statement :", sqlStatement, " args", args)
	return res, err

}
