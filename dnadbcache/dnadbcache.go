package dnadbcache

import "github.com/antonybholmes/go-dna"

var cache = dna.NewDNADbCache()

func SetDir(dir string) {
	cache.SetDir(dir)
}

func Dir() string {
	return cache.Dir()
}

func Db(assembly string, format string, repeatMask string) (*dna.DNADb, error) {
	return cache.Db(assembly, format, repeatMask)
}
