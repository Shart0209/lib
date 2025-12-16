package bloom_filter

type Option func(bf *Service)

// EnableOptimal опция установки оптимального кол-ва элементов и хеш функций
// fpRate - допустимая вероятность ложного срабатывания
func EnableOptimal(fpRate float64) Option {
	return func(bf *Service) {
		bf.fpRate = fpRate
	}
}
