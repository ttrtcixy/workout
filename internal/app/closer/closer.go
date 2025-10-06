package closer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var (
	closerInstance *closer
	once           sync.Once
)

type Logger interface {
	Info(format string, a ...any)
	Error(format string, a ...any)
	Fatal(format string, a ...any)
	Debug(format string, a ...any)
	Warning(format string, a ...any)
}

type Closer interface {
	Add(name string, close func(ctx context.Context) error)
	Tasks() tasks
	Close()
}

type task struct {
	name  string
	close func(ctx context.Context) error
}

type tasks []task

type closer struct {
	log           Logger
	mu            sync.Mutex
	tasks         tasks
	totalDuration time.Duration
	funcDuration  time.Duration
}

// Config for closer
type Config struct {
	TotalDuration time.Duration
	FuncDuration  time.Duration
	Logger        Logger
}

// New create or returns instance of closer (singleton)
//
// default config:
// - totalDuration: infinity
// - funcDuration: infinity
func New(config ...Config) Closer {

	once.Do(func() {
		cfg := Config{}
		if len(config) > 0 {
			cfg = config[0]
		}

		if cfg.TotalDuration < 0 {
			cfg.TotalDuration = 0
		}
		if cfg.FuncDuration < 0 {
			cfg.FuncDuration = 0
		}

		closerInstance = &closer{
			log:           cfg.Logger,
			totalDuration: cfg.TotalDuration,
			funcDuration:  cfg.FuncDuration,
		}
		closerInstance.start()
	})

	return closerInstance
}

func (c *closer) start() {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch
		signal.Stop(ch)
		c.Close()
	}()
}

// todo add custom close time

func (c *closer) Add(name string, close func(ctx context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.tasks = append(c.tasks, task{
		name:  name,
		close: close,
	})
}

func (c *closer) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.log.Info("[*] closer is starting to close tasks")

	if len(c.tasks) == 0 {
		c.log.Info("[*] no tasks to close")
		os.Exit(0)
	}

	timer := time.Now()

	var (
		ctx           = context.Background()
		cancel        context.CancelFunc
		wg            sync.WaitGroup
		done          = make(chan struct{})
		mu            = sync.Mutex{}
		closeErrTasks = make([]string, 0, len(c.tasks))
	)

	if c.totalDuration != 0 {
		ctx, cancel = context.WithTimeout(ctx, c.totalDuration)
		defer cancel()
	}

	for _, task := range c.tasks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := c.closeWithTimeout(ctx, task)
			if err != nil {
				mu.Lock()
				closeErrTasks = append(closeErrTasks, task.name)
				mu.Unlock()

				c.log.Error("[-] task: %s, close error: %s", task.name, err.Error())
				return
			}

			c.log.Info("[+] task: %s, complete", task.name)

		}()
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		if len(closeErrTasks) > 0 {
			c.log.Warning("[-] closer finished, with errors, total duration: " + time.Since(timer).String())
			c.log.Warning("tasks failed to close: [%s]", strings.Join(closeErrTasks, ", "))
			break
		}
		c.log.Info("[+] closer finished, all tasks closed, total duration: " + time.Since(timer).String())
	}
}

func (c *closer) closeWithTimeout(globalCtx context.Context, t task) error {
	var (
		fnCtx  context.Context
		cancel context.CancelFunc
		done   = make(chan error, 1)
	)

	if c.funcDuration != 0 {
		fnCtx, cancel = context.WithTimeout(context.Background(), c.funcDuration)
		defer cancel()
	} else {
		fnCtx = context.Background()
	}

	go func() {
		// close all tasks with timeout
		done <- t.close(fnCtx)
	}()

	select {
	case <-globalCtx.Done():
		return fmt.Errorf("timeout exceeded, max total closing task duration: %v", c.totalDuration)
	case <-fnCtx.Done():
		return fmt.Errorf("timeout exceeded, max closing task duration: %v", c.funcDuration)
	case err := <-done:
		return err
	}
}

func (c *closer) Tasks() tasks {
	return c.tasks
}
