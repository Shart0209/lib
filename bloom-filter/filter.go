package bloom_filter

import (
	"crypto/rand"
	"errors"
	"github.com/cespare/xxhash/v2"
	"math"
	"sync"
)

const (
	blockElem   uint64 = 64
	defaultHash uint64 = 3
)

type BloomFilter struct {
	m         uint64 // кол-во элементов
	k         uint64 // кол-во хешей
	blockSize uint64 //размер массива
	salt      []byte
	bitSet    []bitmap
	mu        *sync.Mutex
	hash      *xxhash.Digest
	// options
	fpRate        float64
	enableOptimal bool
}

// New n - число элементов
func New(n uint64, opts ...Option) (*BloomFilter, error) {
	var m, k, size uint64

	if n == 0 {
		return nil, errors.New("n element must be greater than 0")
	}

	bf := &BloomFilter{
		hash: xxhash.New(),
		mu:   &sync.Mutex{},
	}

	for _, opt := range opts {
		opt(bf)
	}

	m, k = n, defaultHash
	if size = uint64(math.Ceil(float64(n) / float64(blockElem))); size <= 1 {
		k = 1
	}

	if bf.enableOptimal {
		m, k = getOptimalParams(n, bf.fpRate)

		if size = uint64(math.Ceil(float64(m) / float64(blockElem))); size <= 1 {
			k = 1
		}
	}

	bf.m = m
	bf.k = k
	bf.blockSize = size
	bf.salt = generateSalt(int(k))
	bf.bitSet = make([]bitmap, int(size))

	return bf, nil
}

func (b *BloomFilter) Add(value []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, salt := range b.salt {
		hash := b.hashData(value, salt)
		b.bitSet[int(hash%b.blockSize)].Set(hash % blockElem)
	}
}

func (b *BloomFilter) Check(value []byte) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, salt := range b.salt {
		hash := b.hashData(value, salt)
		if !b.bitSet[int(hash%b.blockSize)].Check(hash % blockElem) {
			return false
		}
	}

	return true
}

func (b *BloomFilter) hashData(value []byte, salt byte) uint64 {
	b.hash.ResetWithSeed(uint64(salt))
	b.hash.Write(value)

	return b.hash.Sum64()
}

func generateSalt(size int) []byte {
	data := make([]byte, size)
	rand.Read(data)

	return data
}

// getOptimalParams
// возвращает m - Количество битов в наборе битов; k - Количество используемых хеш-функций
func getOptimalParams(n uint64, p float64) (m, k uint64) {
	if m = uint64(math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2))); m == 0 {
		m = 1
	}

	if k = uint64(math.Ceil((float64(m) / float64(n)) * math.Log(2))); k == 0 {
		k = 1
	}

	return
}
