package bloom_filter

type Option func(bf *BloomFilter)

// EnableOptimal опция установки оптимального кол-ва элементов и хеш функций
// fpRate - допустимая вероятность ложного срабатывания
func EnableOptimal(fpRate float64) Option {
	return func(bf *BloomFilter) {
		bf.fpRate = fpRate
		bf.enableOptimal = true
	}
}
