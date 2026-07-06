package postgres

import "context"

type OptionFunc func(ctx context.Context, pool *pgPool) error

func WithNewSyncConnect(cfg *Config) OptionFunc {
	return func(ctx context.Context, pool *pgPool) error {
		db, err := pool.createPool(ctx, cfg)
		if err != nil {
			return err
		}

		pool.sync = db

		return nil
	}
}
