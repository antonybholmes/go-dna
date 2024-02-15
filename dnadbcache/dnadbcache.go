package dnadbcache

import "github.com/antonybholmes/go-dna"

var Cache = dna.NewDNADbCache()

func Dir(dir string) *dna.DNADbCache {
	Cache.Dir(dir)
	return Cache
}
func Db(assembly string, format string, repeatMask string) (*dna.DNADb, error) {
	return Cache.Db(assembly, format, repeatMask)
}
