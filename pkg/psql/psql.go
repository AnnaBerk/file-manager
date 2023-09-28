package psql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"time"
)

const (
	defaultConnAttempts = 5
	defaultConnTimeout  = 2 * time.Second
)

type PgxConn interface {
	Close(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type Postgres struct {
	connAttempts int
	connTimeout  time.Duration

	Conn PgxConn
}

func New(url string) (*Postgres, error) {
	pg := &Postgres{
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	var err error

	for pg.connAttempts > 0 {
		pg.Conn, err = pgx.Connect(context.Background(), url)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("pgdb - New - pgx.Connect: %w", err)
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Conn != nil {
		err := p.Conn.Close(context.Background())
		if err != nil {
			return
		}
	}
}
