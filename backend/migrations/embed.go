// Package migrations nhúng toàn bộ file .sql vào binary để service tự chạy
// migration khi khởi động (không cần psql trên máy đích).
package migrations

import "embed"

// FS chứa mọi file migration, sắp xếp chạy theo thứ tự tên file.
//
//go:embed *.sql
var FS embed.FS
