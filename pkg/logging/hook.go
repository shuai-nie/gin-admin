package logging

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type HookExecuter interface {
	Exec(extra map[string]string, b []byte) error
	Close() error
}

type hookOptions struct {
	maxJobs    int
	maxWorkers int
	extra      map[string]string
}

func SetHookMaxJobs(maxJobs int) HookOption {
	return func(o *hookOptions) {
		o.maxJobs = maxJobs
	}
}

func SetHookMaxWorkers(maxWorkers int) HookOption {
	return func(o *hookOptions) {
		o.maxWorkers = maxWorkers
	}
}

func SetHookExtra(extra map[string]string) HookOption {
	return func(o *hookOptions) {
		o.extra = extra
	}
}

type HookOption func(*hookOptions)

func NewHook(exec HookExecuter, opt ...HookOption) *Hook {
	opts := &hookOptions{
		maxJobs:    100,
		maxWorkers: 10,
	}

	for _, o := range opt {
		o(opts)
	}

	wg := new(sync.WaitGroup)
	wg.Add(opts.maxWorkers)

	h := &Hook{
		opts: opts,
		q:    make(chan []byte, opts.maxJobs),
		wg:   wg,
		e:    exec,
	}
	h.dispatch()
	return h
}

type Hook struct {
	opts   *hookOptions
	q      chan []byte
	wg     *sync.WaitGroup
	e      HookExecuter
	closed int32
}

func (h *Hook) dispatch() {
	for i := 0; i < h.opts.maxWorkers; i++ {
		go func() {
			defer func() {
				h.wg.Done()
				if r := recover(); r != nil {
					fmt.Println("panic:", r)
				}
			}()

			for data := range h.q {
				err := h.e.Exec(h.opts.extra, data)
				if err != nil {
					fmt.Println("exec error:", err)
				}
			}
		}()
	}
}

func (h *Hook) Write(p []byte) (int, error) {
	if atomic.LoadInt32(&h.closed) == 1 {
		return len(p), nil
	}

	if len(h.q) == h.opts.maxJobs {
		fmt.Println("queue full")
		return len(p), nil
	}

	data := make([]byte, len(p))
	copy(data, p)
	h.q <- data
	return len(p), nil
}

func (h *Hook) Flush() {
	if atomic.LoadInt32(&h.closed) == 1 {
		return
	}
	atomic.StoreInt32(&h.closed, 1)
	close(h.q)
	h.wg.Wait()
	err := h.e.Close()
	if err != nil {
		fmt.Println("close error:", err)
	}
}
