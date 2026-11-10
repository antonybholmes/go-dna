package dnadb

import (
	"sync"

	"github.com/antonybholmes/go-dna"
)

var (
	instance *dna.DnaDB
	once     sync.Once
)

func InitDnaDB(dir string) *dna.DnaDB {
	once.Do(func() {
		instance = dna.NewDnaDB(dir)
	})

	return instance
}

func Instance() *dna.DnaDB {
	return instance
}
