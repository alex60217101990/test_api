package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	"go.uber.org/zap"

	"github.com/alex60217101990/test_api/internal/configs"
	"github.com/alex60217101990/test_api/internal/logger"
	"github.com/alex60217101990/test_api/internal/repository"
	sqldblogger "github.com/simukti/sqldb-logger"
	zaplogger "github.com/simukti/sqldb-logger/logadapter/zapadapter"

	"github.com/fatih/color"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// Repository ...
type Repository struct {
	pool *pgxpool.Pool
	conn *sql.DB
	// can use for metrics and monitoring...
	connCounter int32
}

// NewPostgresRepository ...
func NewPostgresRepository() repository.Repository {
	return &Repository{}
}

func (r *Repository) Connect(ctx context.Context) error {
	config, err := pgxpool.ParseConfig(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			configs.Conf.DB.UserName,
			configs.Conf.DB.Password,
			configs.Conf.DB.Host,
			configs.Conf.DB.Port,
			configs.Conf.DB.DbName,
		),
	)
	if err != nil {
		logger.AppLogger.Error(err)
		return err
	}

	config.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		atomic.AddInt32(&r.connCounter, 1)
		if configs.Conf.IsDebug {
			logger.AppLogger.Info(fmt.Sprintf("repo connections count was change: [%d]", r.connCounter))
		}
		return true
	}

	config.AfterRelease = func(conn *pgx.Conn) bool {
		delta := atomic.AddInt32(&r.connCounter, -1)
		if delta < 0 {
			delta = atomic.AddInt32(&r.connCounter, 1)
		}
		if configs.Conf.IsDebug {
			logger.AppLogger.Info(fmt.Sprintf("repo connections count was change: [%d]", r.connCounter))
		}
		return true
	}

	config.MaxConns = 10
	config.HealthCheckPeriod = 5
	if configs.Conf.LoggerType == configs.Zap {
		config.ConnConfig.Logger = zapadapter.NewLogger(logger.AppLogger.GetNativeLogger().(*zap.Logger))
	}

	if r.pool, err = pgxpool.ConnectConfig(ctx, config); err != nil {
		logger.AppLogger.Error(err)
		return err
	}

	r.conn, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		configs.Conf.DB.Host,
		configs.Conf.DB.Port,
		configs.Conf.DB.UserName,
		configs.Conf.DB.Password,
		configs.Conf.DB.DbName,
	))
	if err != nil {
		logger.AppLogger.Error(err)
		return err
	}

	if configs.Conf.LoggerType == configs.Zap {
		r.conn = sqldblogger.OpenDriver(
			fmt.Sprintf("host=%s port=%d user=%s "+
				"password=%s dbname=%s sslmode=disable",
				configs.Conf.DB.Host,
				configs.Conf.DB.Port,
				configs.Conf.DB.UserName,
				configs.Conf.DB.Password,
				configs.Conf.DB.DbName,
			),
			r.conn.Driver(),
			zaplogger.New(logger.AppLogger.GetNativeLogger().(*zap.Logger)),
		)
	}

	r.conn.SetConnMaxLifetime(0)
	r.conn.SetMaxIdleConns(100)
	r.conn.SetMaxOpenConns(10)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err = r.Ping(ctx); err != nil {
		logger.AppLogger.Error(err)
		return err
	}

	color.New(color.FgMagenta, color.Bold).Println("🐘 Postgres repository connect success.")

	return nil
}

func (r *Repository) Close() error {
	color.New(color.FgMagenta, color.Bold).Println("🐘 Postgres repository connect close.")
	r.pool.Close()
	return r.conn.Close()
}

func (r *Repository) GetDB() (interface{}, error) {
	return r.pool.Acquire(context.Background())
}

func (r *Repository) closeTx(ctx context.Context, tx *sql.Tx, err error) {
	if r := recover(); r != nil {
		tx.Rollback()
		logger.AppLogger.Fatal(err)
	}
	if err != nil {
		err = errors.WithMessage(err, "tx close posrgres repo")
		tx.Rollback()
		return
	}
	tx.Commit()
}

func (r *Repository) initTx(ctx context.Context) (tx *sql.Tx, err error) {
	tx, err = r.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "tx init posrgres repo")
		}
	}()

	return tx, err
}

func (r *Repository) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"acquired_conns":         r.pool.Stat().AcquiredConns(),
		"max_conns":              r.pool.Stat().MaxConns(),
		"acquire_count":          r.pool.Stat().AcquireCount(),
		"acquire_duration":       r.pool.Stat().AcquireDuration(),
		"canceled_acquire_count": r.pool.Stat().CanceledAcquireCount(),
		"constructing_conns":     r.pool.Stat().ConstructingConns(),
		"empty_acquire_count":    r.pool.Stat().EmptyAcquireCount(),
		"idle_conns":             r.pool.Stat().IdleConns(),
		"total_conns":            r.pool.Stat().TotalConns(),
	}
}

func (r *Repository) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	return r.pool.Acquire(ctx)
}

func (r *Repository) Ping(ctx context.Context) (err error) {
	var conn *pgxpool.Conn
	conn, err = r.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	eg, _ := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return conn.Conn().Ping(ctx)
	})

	eg.Go(func() error {
		return r.conn.PingContext(ctx)
	})

	return eg.Wait()

	// if err = conn.Conn().Ping(ctx); err != nil {
	// 	return err
	// }

	// if err = ; err != nil {
	// 	return err
	// }

	//return nil
}
