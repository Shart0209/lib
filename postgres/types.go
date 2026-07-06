package postgres

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	User     string
	Password string
	Schema   string

	MaxConns        int32
	MinConns        int32
	SSLMode         string
	ConnectTimeout  time.Duration
	MaxConnLifeTime time.Duration
	MaxConnIdleTime time.Duration
}

type pgPool struct {
	master *pgConnect
	sync   *pgConnect
}

type pgConnect struct {
	conn *pgxpool.Pool
}
