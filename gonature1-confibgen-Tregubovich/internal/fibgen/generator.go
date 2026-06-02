package fibonacci

import (
	"errors"
	"runtime"
	"sync/atomic"
)

var ErrOverflow = errors.New("overflow")

const (
	unlocked = iota
	locked
)

type SpinLock struct {
	state atomic.Int64
}

func (s *SpinLock) Lock() {
	for {
		if s.state.CompareAndSwap(unlocked, locked) {
			return
		}

		runtime.Gosched()
	}
}

func (s *SpinLock) Unlock() {
	s.state.Store(unlocked)
}

type Generator interface {
	Next() uint64
}

var _ Generator = (*generatorImpl)(nil)

type generatorImpl struct {
	mx SpinLock
	a  uint64
	b  uint64
}

func NewGenerator() *generatorImpl {
	return &generatorImpl{
		mx: SpinLock{},
		a:  0,
		b:  1,
	}
}

func (g *generatorImpl) Next() uint64 {
	g.mx.Lock()
	defer g.mx.Unlock()

	g.isOverflow()

	a, b := g.a, g.b

	c := a + b
	if c < b {
		c = 0
	}

	g.a = b
	g.b = c
	return a
}

func (g *generatorImpl) isOverflow() {
	if g.a == 0 && g.b != 1 {
		panic(ErrOverflow)
	}
}
