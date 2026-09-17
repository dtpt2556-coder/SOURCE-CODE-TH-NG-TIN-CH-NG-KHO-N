package store

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/migrations"
)

// Migration được chạy theo thứ tự tên file — test offline, không cần DB.
func TestListMigrationFilesSortedByName(t *testing.T) {
	fsys := fstest.MapFS{
		"0010_late.sql": {Data: []byte("SELECT 1")},
		"0001_init.sql": {Data: []byte("SELECT 1")},
		"0002_seed.sql": {Data: []byte("SELECT 1")},
		"readme.md":     {Data: []byte("bỏ qua")},
		"0003_demo.sql": {Data: []byte("SELECT 1")},
		"sub/other.sql": {Data: []byte("SELECT 1")},
	}
	got, err := listMigrationFiles(fsys)
	require.NoError(t, err)
	assert.Equal(t, []string{"0001_init.sql", "0002_seed.sql", "0003_demo.sql", "0010_late.sql"}, got)
}

func TestListMigrationFilesRequiresAtLeastOne(t *testing.T) {
	_, err := listMigrationFiles(fstest.MapFS{"readme.md": {Data: []byte("x")}})
	assert.Error(t, err)
}

// Thư mục migrations thật phải nhúng được và đủ 4 file bắt buộc.
func TestEmbeddedMigrationsArePresent(t *testing.T) {
	got, err := listMigrationFiles(migrations.FS)
	require.NoError(t, err)
	assert.Equal(t, []string{
		"0001_init.sql",
		"0002_seed_sources.sql",
		"0003_seed_tickers.sql",
		"0004_seed_demo.sql",
		"0005_edge_cases.sql",
		"0006_qa_p0_fixes.sql",
		"0007_qa_round2_fixes.sql",
	}, got)
}
