package shortlink_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/shortlink"
)

func TestLRUCacheGetPut(t *testing.T) {
	c := shortlink.NewLRUCache(3)

	_, ok := c.Get("a7Kx2p")
	assert.False(t, ok)

	c.Put("a7Kx2p", shortlink.Target{ArticleID: 1, URL: "https://vnexpress.net/x.html", Alive: true})
	got, ok := c.Get("a7Kx2p")
	require.True(t, ok)
	assert.Equal(t, "https://vnexpress.net/x.html", got.URL)
	assert.True(t, got.Alive)
}

func TestLRUCacheEvictsOldest(t *testing.T) {
	c := shortlink.NewLRUCache(2)
	c.Put("aaaaaa", shortlink.Target{URL: "https://cafef.vn/1"})
	c.Put("bbbbbb", shortlink.Target{URL: "https://cafef.vn/2"})

	// Chạm vào "aaaaaa" để nó trở thành mới nhất.
	_, _ = c.Get("aaaaaa")
	c.Put("cccccc", shortlink.Target{URL: "https://cafef.vn/3"})

	assert.Equal(t, 2, c.Len())
	_, ok := c.Get("bbbbbb")
	assert.False(t, ok, "phần tử cũ nhất phải bị đẩy ra")
	_, ok = c.Get("aaaaaa")
	assert.True(t, ok)
}

func TestLRUCacheConcurrentAccess(t *testing.T) {
	c := shortlink.NewLRUCache(100)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			code := fmt.Sprintf("code%02d", i)
			c.Put(code, shortlink.Target{ArticleID: int64(i), URL: "https://cafef.vn/x"})
			_, _ = c.Get(code)
		}(i)
	}
	wg.Wait()
	assert.Equal(t, 50, c.Len())
}

func TestClickCounterBatchesAndFlushes(t *testing.T) {
	var mu sync.Mutex
	flushed := map[string]int64{}

	counter := shortlink.NewClickCounter(func(_ context.Context, counts map[string]int64) error {
		mu.Lock()
		defer mu.Unlock()
		for k, v := range counts {
			flushed[k] += v
		}
		return nil
	}, 10*time.Millisecond, nil)

	for i := 0; i < 7; i++ {
		counter.Add("a7Kx2p")
	}
	counter.Add("b3Qw9m")
	assert.Equal(t, 2, counter.Pending())

	counter.Flush(context.Background())
	assert.Equal(t, 0, counter.Pending(), "flush xong phải rỗng bộ đệm")

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, int64(7), flushed["a7Kx2p"])
	assert.Equal(t, int64(1), flushed["b3Qw9m"])
}

func TestClickCounterSurvivesFlushError(t *testing.T) {
	counter := shortlink.NewClickCounter(func(context.Context, map[string]int64) error {
		return errors.New("DB tạm thời không khả dụng")
	}, time.Second, nil)

	counter.Add("a7Kx2p")
	// Lỗi flush không được panic và không được chặn luồng redirect.
	counter.Flush(context.Background())
	counter.Add("a7Kx2p")
	assert.Equal(t, 1, counter.Pending())
}

func TestClickCounterRunStopsWithContext(t *testing.T) {
	done := make(chan struct{})
	counter := shortlink.NewClickCounter(func(context.Context, map[string]int64) error {
		return nil
	}, 5*time.Millisecond, nil)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		counter.Run(ctx)
		close(done)
	}()

	counter.Add("a7Kx2p")
	time.Sleep(30 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("ClickCounter.Run không dừng khi context bị huỷ")
	}
}

func TestClickCounterIgnoresEmptyCode(t *testing.T) {
	counter := shortlink.NewClickCounter(nil, time.Second, nil)
	counter.Add("")
	assert.Equal(t, 0, counter.Pending())
}
