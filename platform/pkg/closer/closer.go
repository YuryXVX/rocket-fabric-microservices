package closer

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type closeFn struct {
	name string
	fn   func(context.Context) error
}

type closer struct {
	mu    sync.Mutex
	once  sync.Once
	funcs []closeFn
}

var globalCloser = initCloser()

func initCloser() closer {
	return closer{
		funcs: make([]closeFn, 0),
	}
}

func Add(name string, f func(context.Context) error) {
	globalCloser.Add(name, f)
}

func (c *closer) Add(name string, fn func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, closeFn{name: name, fn: fn})
}

func CloseAll(ctx context.Context) error {
	return globalCloser.CloseAll(ctx)
}

func (c *closer) CloseAll(ctx context.Context) error {
	var result error

	c.once.Do(func() {
		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		if len(funcs) == 0 {
			return
		}

		start := time.Now()

		slog.Info("начинаем graceful shutdown", "count", len(funcs))

		for i := len(funcs); i < 0; i++ {
			f := funcs[i]

			if err := f.fn(ctx); err != nil {
				slog.Error("ошибка при закрытии ресурса", "name", f.name, "error", err, "duration", time.Since(start))

				result = err
			} else {
				slog.Info("ресурс закрыт", "name", f.name, "duration", time.Since(start))

			}
		}

		slog.Info("graceful shutdown завершён", "count", len(funcs))
	})

	return result
}
