package store

import (
	"context"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

// createMigrationsTable là bước 0, luôn chạy trước mọi file migration.
const createMigrationsTable = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version    text PRIMARY KEY,
    applied_at timestamptz NOT NULL DEFAULT now()
)`

// Migrate chạy toàn bộ file .sql trong fsys theo thứ tự tên file và ghi nhận
// vào bảng schema_migrations. Idempotent — chạy lại không lỗi, không chạy lại
// file đã áp dụng.
func (s *Store) Migrate(ctx context.Context, fsys fs.FS) error {
	if _, err := s.pool.Exec(ctx, createMigrationsTable); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	applied, err := s.appliedVersions(ctx)
	if err != nil {
		return err
	}

	files, err := listMigrationFiles(fsys)
	if err != nil {
		return err
	}

	for _, name := range files {
		version := strings.TrimSuffix(name, ".sql")
		if applied[version] {
			continue
		}
		body, err := fs.ReadFile(fsys, name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		if err := s.applyMigration(ctx, version, string(body)); err != nil {
			return err
		}
		s.log.Info("migration applied", "version", version)
	}
	return nil
}

func (s *Store) appliedVersions(ctx context.Context) (map[string]bool, error) {
	rows, err := s.pool.Query(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, fmt.Errorf("read schema_migrations: %w", err)
	}
	defer rows.Close()

	out := map[string]bool{}
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("scan schema_migrations: %w", err)
		}
		out[v] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate schema_migrations: %w", err)
	}
	return out, nil
}

// applyMigration chạy nội dung file bằng simple query protocol — PostgreSQL bọc
// cả batch trong một transaction ngầm nên migration là nguyên tử.
func (s *Store) applyMigration(ctx context.Context, version, body string) error {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire connection for migration %s: %w", version, err)
	}
	defer conn.Release()

	if _, err := conn.Conn().PgConn().Exec(ctx, body).ReadAll(); err != nil {
		return fmt.Errorf("run migration %s: %w", version, err)
	}
	if _, err := conn.Exec(ctx,
		"INSERT INTO schema_migrations (version) VALUES ($1) ON CONFLICT DO NOTHING", version); err != nil {
		return fmt.Errorf("record migration %s: %w", version, err)
	}
	return nil
}

func listMigrationFiles(fsys fs.FS) ([]string, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("read migrations directory: %w", err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)
	if len(names) == 0 {
		return nil, fmt.Errorf("no migration files found")
	}
	return names, nil
}
