package dnadbcache

import "github.com/antonybholmes/go-dna"

var cache = dna.NewDNADbCache()

func Dir(dir string) {
	cache.Dir(dir)
}

func Db(assembly string, format string, repeatMask string) (*dna.DNADb, error) {
	return cache.Db(assembly, format, repeatMask)
}
