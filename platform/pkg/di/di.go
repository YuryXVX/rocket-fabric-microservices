package di

import (
	"context"
	"sync"
)

type Value[T any] struct {
	once sync.Once
	val  T
}

// Get возвращает инициализированный объект.
// Если объект еще не был создан, вызывается переданная функция-конструктор initFn.
func (v *Value[T]) Get(ctx context.Context, initFn func(ctx context.Context) T) T {
	v.once.Do(func() {
		v.val = initFn(ctx)
	})

	return v.val
}
