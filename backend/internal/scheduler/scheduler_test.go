package scheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/scheduler"
)

func vnLocation(t *testing.T) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	require.NoError(t, err)
	return loc
}

func TestDefaultCronSpecIsValid(t *testing.T) {
	s, err := scheduler.New("0 6,14,22 * * *", vnLocation(t),
		func(context.Context) error { return nil }, nil)
	require.NoError(t, err)

	s.Start()
	defer s.Stop(context.Background())

	next := s.NextRun()
	require.NotEmpty(t, next)

	parsed, err := time.Parse(time.RFC3339, next)
	require.NoError(t, err)

	// Lịch phải rơi đúng 06:00 / 14:00 / 22:00 GIỜ VIỆT NAM (ADR-004).
	inVN := parsed.In(vnLocation(t))
	assert.Contains(t, []int{6, 14, 22}, inVN.Hour(),
		"giờ chạy kế tiếp %s không thuộc 6/14/22 giờ VN", inVN)
	assert.Equal(t, 0, inVN.Minute())
	assert.True(t, parsed.After(time.Now()))
}

func TestInvalidCronSpecReturnsError(t *testing.T) {
	_, err := scheduler.New("không-phải-cron", vnLocation(t),
		func(context.Context) error { return nil }, nil)
	assert.Error(t, err)
}

func TestNilJobIsRejected(t *testing.T) {
	_, err := scheduler.New("0 6 * * *", vnLocation(t), nil, nil)
	assert.Error(t, err)
}

func TestSchedulerRunsJob(t *testing.T) {
	done := make(chan struct{}, 1)
	s, err := scheduler.New("@every 100ms", time.UTC, func(context.Context) error {
		select {
		case done <- struct{}{}:
		default:
		}
		return nil
	}, nil)
	require.NoError(t, err)

	s.Start()
	defer s.Stop(context.Background())

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("job không được kích hoạt trong 3 giây")
	}
}
