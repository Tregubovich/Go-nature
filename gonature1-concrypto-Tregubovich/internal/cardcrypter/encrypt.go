package cardcrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"runtime"
	"sync"
	"unsafe"
)

type CardNumber [16]byte

type Card struct {
	ID     string
	Number CardNumber
}

var ErrNegativeWorkers = errors.New("negative workers")

type Crypter interface {
	Encrypt(cards []Card, key []byte) ([]string, error)
}

type crypterImpl struct {
	workers int
}

func New(opts ...CrypterOption) *crypterImpl {
	c := &crypterImpl{
		workers: runtime.GOMAXPROCS(0),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type CrypterOption func(*crypterImpl)

func WithWorkers(workers int) CrypterOption {
	return func(c *crypterImpl) {
		c.workers = workers
	}
}

func (c *crypterImpl) Encrypt(cards []Card, key []byte) ([]string, error) {
	if c.workers <= 0 {
		return nil, ErrNegativeWorkers
	}

	wg := sync.WaitGroup{}

	if len(cards) == 0 {
		return nil, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()

	out := make([]string, len(cards))

	workers := min(c.workers, len(cards))
	delta := len(cards)/workers + 1

	for w := range workers {
		nonce := make([]byte, nonceSize)
		buf := make([]byte, nonceSize, nonceSize+16+gcm.Overhead())

		wg.Go(func() {
			for i := w * delta; i < min((w+1)*delta, len(cards)); i++ {
				_, err := rand.Read(nonce)
				if err != nil {
					panic("unreachable")
				}

				copy(buf, nonce)

				result := gcm.Seal(
					buf,
					nonce,
					cards[i].Number[:],
					stringToSlice(cards[i].ID),
				)

				dst := make([]byte, len(result)*2)
				hex.Encode(dst, result)
				out[i] = sliceToString(dst)
			}
		})
	}
	wg.Wait()

	return out, nil
}

func stringToSlice(src string) []byte {
	return unsafe.Slice(unsafe.StringData(src), len(src))
}

func sliceToString(src []byte) string {
	return unsafe.String(unsafe.SliceData(src), len(src))
}
