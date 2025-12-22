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

func GetInstance() *dna.DnaDB {
	return instance
}

func Dir() string {
	return instance.Dir
}

func Db(assembly string) (*dna.AssemblyDB, error) {
	return instance.DB(assembly)
}
