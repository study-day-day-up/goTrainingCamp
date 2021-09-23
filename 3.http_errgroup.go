package main

// 因无法下载 errgroup 包，故将 errgroup 源码直接放入该文件
// 源码地址：https://github.com/golang/sync/blob/master/errgroup/errgroup.go

// Package errgroup provides synchronization, error propagation, and Context
// cancelation for groups of goroutines working on subtasks of a common task.
import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero Group is valid and does not cancel on error.
type Group struct {
	cancel func()

	wg sync.WaitGroup

	errOnce sync.Once
	err     error
}

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}

func main() {
	sigs := []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sigs...)

	eg, ctx := WithContext(context.Background())

	eg.Go(func() error {
		go HandelRequest(ctx)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("HttpServe have Done")
				return ctx.Err()
			case <-ch:
				fmt.Println("HttpServe shutdown by system signal")
				eg.cancel()
			}
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		fmt.Println(err)
	}
}

func HandelRequest(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hi", hiHandler)

	server := &http.Server{
		Addr:    ":8085",
		Handler: mux,
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("HttpServe cancel")
			return nil
		default:
			fmt.Println("HttpServe start up")
			return server.ListenAndServe()
		}
	}
}

func hiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hi, Go HandleFunc")
}
