package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const packageName = "postgres"

func New(ctx context.Context, cfg *Config, ops ...OptionFunc) (PgPooler, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	pgpool := pgPool{}

	db, err := pgpool.createPool(ctx, cfg)
	if err != nil {
		return nil, err
	}

	pgpool.master = db

	for _, option := range ops {
		if err = option(ctx, &pgpool); err != nil {
			return nil, err
		}
	}

	if pgpool.sync == nil {
		pgpool.sync = db
	}

	return &pgpool, nil
}

func (p *pgPool) checkConnect(ctx context.Context, conn *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := conn.Ping(ctx); err != nil {
		conn.Close()
		return fmt.Errorf("failed to check connect database: %w", err)
	}

	return nil
}

func (p *pgPool) buildDSN(cfg *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Schema, cfg.SSLMode)
}

func (p *pgPool) createPool(ctx context.Context, cfg *Config) (*pgConnect, error) {
	poolConfig, err := pgxpool.ParseConfig(p.buildDSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifeTime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.ConnConfig.ConnectTimeout = cfg.ConnectTimeout

	conn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err = p.checkConnect(ctx, conn); err != nil {
		return nil, err
	}

	return &pgConnect{conn: conn}, nil
}

func (p *pgPool) Close(_ context.Context) error {
	p.master.conn.Close()
	p.sync.conn.Close()
	return nil
}

func (p *pgPool) Name() string {
	return packageName
}

func (p *pgPool) Master() Connector {
	return p.master
}

func (p *pgPool) Sync() Connector {
	return p.sync
}
