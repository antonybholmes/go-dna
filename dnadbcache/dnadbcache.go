package dnadbcache

import (
	"sync"

	"github.com/antonybholmes/go-dna"
)

var (
	instance *dna.DNADBCache
	once     sync.Once
)

func InitCache(dir string) *dna.DNADBCache {
	once.Do(func() {
		instance = dna.NewDNADBCache(dir)
	})

	return instance
}

func GetInstance() *dna.DNADBCache {
	return instance
}

func Dir() string {
	return instance.Dir
}

func Db(assembly string) (*dna.DNADB, error) {
	return instance.DB(assembly)
}
