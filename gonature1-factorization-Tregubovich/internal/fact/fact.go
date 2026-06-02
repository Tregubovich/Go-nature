package fact

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	ErrFactorizationCancelled = errors.New("cancelled")
	ErrWriterInteraction      = errors.New("writer interaction")
)

type Factorizer interface {
	Factorize(ctx context.Context, numbers []int, writer io.Writer) error
}

type factorizerImpl struct {
	factWorkers  int
	writeWorkers int
}

func New(opts ...FactorizeOption) (*factorizerImpl, error) {
	c := &factorizerImpl{
		factWorkers:  runtime.GOMAXPROCS(0),
		writeWorkers: runtime.GOMAXPROCS(0),
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.factWorkers <= 0 {
		return nil, errors.New("expected positive number of factorization workers, but got " + strconv.Itoa(c.factWorkers))
	}
	if c.writeWorkers <= 0 {
		return nil, errors.New("expected positive number of write workers, but got " + strconv.Itoa(c.writeWorkers))
	}

	return c, nil
}

type FactorizeOption func(*factorizerImpl)

func WithFactorizationWorkers(workers int) FactorizeOption {
	return func(c *factorizerImpl) {
		c.factWorkers = workers
	}
}

func WithWriteWorkers(workers int) FactorizeOption {
	return func(c *factorizerImpl) {
		c.writeWorkers = workers
	}
}

func (f *factorizerImpl) Factorize(
	parent context.Context,
	numbers []int,
	writer io.Writer,
) error {
	ctx, cancel := context.WithCancelCause(parent)
	defer cancel(nil)

	resCh := make(chan string)

	wgFact := sync.WaitGroup{}
	factWorkers := min(f.factWorkers, len(numbers))
	genCh := generator(ctx, numbers)

	workerPool(ctx, factWorkers, &wgFact, genCh, func(v int) {
		select {
		case <-ctx.Done():
			return
		case resCh <- factorizeNext(ctx, v):
		}
	})

	wgWrite := sync.WaitGroup{}
	writeWorkers := min(f.writeWorkers, len(numbers))

	workerPool(ctx, writeWorkers, &wgWrite, resCh, func(v string) {
		if err := writeNext(writer, v); err != nil {
			cancel(errors.Join(ErrWriterInteraction, err))
		}
	})

	wgFact.Wait()
	close(resCh)

	wgWrite.Wait()

	if parent.Err() != nil {
		return errors.Join(ErrFactorizationCancelled, context.Cause(parent))
	}

	return context.Cause(ctx)
}

func workerPool[T any](ctx context.Context, workers int, wg *sync.WaitGroup, input <-chan T, action func(T)) {
	for range workers {
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-input:
					if !ok {
						return
					}
					action(v)
				}
			}
		})
	}
}

func generator(ctx context.Context, data []int) <-chan int {
	result := make(chan int)

	go func() {
		defer close(result)

		for _, v := range data {
			select {
			case result <- v:
			case <-ctx.Done():
				return
			}
		}
	}()

	return result
}

func factorizeNext(ctx context.Context, num int) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d = ", num))
	flag := true

	if num < 0 {
		sb.WriteString("-1")
		num = -num
		flag = false
	}

	if num == math.MinInt {
		maxDiv(ctx, &num, 2, &flag, &sb)
		return sb.String()
	}

	for _, d := range []int{2, 3} {
		maxDiv(ctx, &num, d, &flag, &sb)
	}

	for d := 5; d <= num/d+1; {
		maxDiv(ctx, &num, d, &flag, &sb)
		if (d+1)%6 == 0 {
			d += 2
		} else {
			d += 4
		}
	}
	if num != 1 || flag {
		nextDivisor(flag, &sb, num)
	}
	return sb.String()
}

func maxDiv(ctx context.Context, num *int, d int, flag *bool, sb *strings.Builder) {
	n := *num
	for n != 0 && n%d == 0 {
		if ctx.Err() != nil {
			return
		}
		n /= d
		nextDivisor(*flag, sb, d)
		*flag = false
	}
	*num = n
}

func nextDivisor(flag bool, sb *strings.Builder, d int) {
	if flag {
		fmt.Fprintf(sb, "%d", d)
	} else {
		fmt.Fprintf(sb, " * %d", d)
	}
}

func writeNext(writer io.Writer, s string) error {
	if _, err := fmt.Fprintf(writer, "%s", s+"\n"); err != nil {
		return errors.Join(err, ErrWriterInteraction)
	}
	return nil
}
