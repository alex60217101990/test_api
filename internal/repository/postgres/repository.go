package postgres

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/alex60217101990/test_api/internal/configs"
	"github.com/alex60217101990/test_api/internal/logger"
	"github.com/alex60217101990/test_api/internal/repository"
	"github.com/fatih/color"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
	// can use for metrics and monitoring...
	connCounter int32
}

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

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err = r.Ping(ctx); err != nil {
		logger.AppLogger.Error(err)
		return err
	}

	color.New(color.FgMagenta, color.Bold).Println("🐘 Postgres repository connect success.")

	return nil
}

func (r *Repository) Close() {
	logger.CmdServer.Println("🐘 Postgres repository connect close.")
	r.pool.Close()
}

func (r *Repository) GetDB() (interface{}, error) {
	return r.pool.Acquire(context.Background())
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

func (r *Repository) Ping(ctx context.Context) error {
	conn, err := r.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	if err := conn.Conn().Ping(ctx); err != nil {
		return err
	}

	return nil
}
