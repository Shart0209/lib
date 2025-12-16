package bloom_filter

import (
	"crypto/rand"
	"errors"
	"math"
	"strings"
	"sync"

	"github.com/shart0209/lib/bitmap"

	"github.com/cespare/xxhash/v2"
)

const (
	blockElem   uint64 = 64
	defaultHash uint64 = 3
)

type Service struct {
	m         uint64 // кол-во элементов
	k         uint64 // кол-во хешей
	blockSize uint64 //размер массива
	salt      []byte
	bitSet    []bitmap.Bitmap
	mu        *sync.Mutex
	hash      *xxhash.Digest
	fpRate    float64
}

func New(elemCount uint64, opts ...Option) (*Service, error) {
	if elemCount == 0 {
		return nil, errors.New("element must be greater than 0")
	}

	bf := &Service{
		hash: xxhash.New(),
		mu:   &sync.Mutex{},
	}

	for _, opt := range opts {
		opt(bf)
	}

	if bf.fpRate > 0 {
		bf.m, bf.k = getOptimalParams(elemCount, bf.fpRate)
	} else {
		bf.m, bf.k = elemCount, defaultHash
	}

	if bf.blockSize = uint64(math.Ceil(float64(bf.m) / float64(blockElem))); bf.blockSize <= 1 {
		bf.k = 1
	}

	bf.salt = generateSalt(int(bf.k))
	bf.bitSet = make([]bitmap.Bitmap, int(bf.blockSize))

	return bf, nil
}

func (s *Service) Add(value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, salt := range s.salt {
		hash := s.hashData(value, salt)
		s.bitSet[int(hash%s.blockSize)].Set(hash % blockElem)
	}
}

func (s *Service) Check(value []byte) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, salt := range s.salt {
		hash := s.hashData(value, salt)
		if !s.bitSet[int(hash%s.blockSize)].Check(hash % blockElem) {
			return false
		}
	}

	return true
}

func (s *Service) Clear(value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, salt := range s.salt {
		hash := s.hashData(value, salt)
		s.bitSet[int(hash%s.blockSize)].Clear(hash % blockElem)
	}
}

func (s *Service) ClearAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.bitSet = make([]bitmap.Bitmap, int(s.blockSize))
}

func (s *Service) String(value []byte) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := make([]string, 0)
	for _, salt := range s.salt {
		hash := s.hashData(value, salt)
		res = append(res, s.bitSet[int(hash%s.blockSize)].String())
	}

	return strings.Join(res, "; ")
}

func (s *Service) hashData(value []byte, salt byte) uint64 {
	s.hash.ResetWithSeed(uint64(salt))
	s.hash.Write(value)

	return s.hash.Sum64()
}

func generateSalt(size int) []byte {
	data := make([]byte, size)
	rand.Read(data)

	return data
}

// getOptimalParams
// возвращает m - Количество битов в наборе битов; k - Количество используемых хеш
func getOptimalParams(n uint64, p float64) (m, k uint64) {
	if m = uint64(math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2))); m == 0 {
		m = 1
	}

	if k = uint64(math.Ceil((float64(m) / float64(n)) * math.Log(2))); k == 0 {
		k = 1
	}

	return
}
