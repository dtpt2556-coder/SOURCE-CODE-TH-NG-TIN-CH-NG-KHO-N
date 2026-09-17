// Package store là tầng truy cập PostgreSQL bằng pgx.
// Chỉ tầng này biết về SQL; domain không bao giờ import package này.
package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound được trả khi truy vấn không có kết quả.
var ErrNotFound = errors.New("record not found")

// PipelineLockKey là khoá advisory lock chống 2 pipeline chạy chồng
// (mục 11 architecture.md).
const PipelineLockKey int64 = 8_102_026

// Store giữ connection pool. Inject qua struct, không dùng biến global.
type Store struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

// New mở pool và kiểm tra kết nối.
func New(ctx context.Context, dsn string, log *slog.Logger) (*Store, error) {
	if log == nil {
		log = slog.Default()
	}
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}
	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 15 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL pool: %w", err)
	}
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
	}
	return &Store{pool: pool, log: log}, nil
}

// Close đóng pool.
func (s *Store) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

// Ping kiểm tra DB còn sống (dùng cho /readyz).
func (s *Store) Ping(ctx context.Context) error {
	if err := s.pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping PostgreSQL: %w", err)
	}
	return nil
}

// Pool trả về pool cho các thao tác nâng cao.
func (s *Store) Pool() *pgxpool.Pool { return s.pool }

// AdvisoryLock giữ pg_advisory_lock trên một connection riêng.
// Bắt buộc gọi Release() để trả khoá và connection về pool.
type AdvisoryLock struct {
	conn *pgxpool.Conn
	key  int64
}

// TryAdvisoryLock thử lấy khoá. acquired=false nghĩa là đang có pipeline khác chạy.
func (s *Store) TryAdvisoryLock(ctx context.Context, key int64) (*AdvisoryLock, bool, error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("acquire connection for advisory lock: %w", err)
	}
	var ok bool
	if err := conn.QueryRow(ctx, "SELECT pg_try_advisory_lock($1)", key).Scan(&ok); err != nil {
		conn.Release()
		return nil, false, fmt.Errorf("pg_try_advisory_lock(%d): %w", key, err)
	}
	if !ok {
		conn.Release()
		return nil, false, nil
	}
	return &AdvisoryLock{conn: conn, key: key}, true, nil
}

// Release trả khoá và connection.
func (l *AdvisoryLock) Release(ctx context.Context) error {
	if l == nil || l.conn == nil {
		return nil
	}
	defer func() {
		l.conn.Release()
		l.conn = nil
	}()
	if _, err := l.conn.Exec(ctx, "SELECT pg_advisory_unlock($1)", l.key); err != nil {
		return fmt.Errorf("pg_advisory_unlock(%d): %w", l.key, err)
	}
	return nil
}

// wrapNoRows đổi pgx.ErrNoRows thành ErrNotFound của tầng store.
func wrapNoRows(err error, what string) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%s: %w", what, ErrNotFound)
	}
	return fmt.Errorf("%s: %w", what, err)
}
