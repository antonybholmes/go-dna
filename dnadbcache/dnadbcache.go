package dnadbcache

import "github.com/antonybholmes/go-dna"

var cache *dna.DNADBCache

func InitCache(dir string) {
	cache = dna.NewDNADBCache(dir)
}

func Dir() string {
	return cache.Dir
}

func Db(assembly string, format string, repeatMask string) (*dna.DNADB, error) {
	return cache.DB(assembly, format, repeatMask)
}
