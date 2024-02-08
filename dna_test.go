package dna

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestWithin(t *testing.T) {
	location, err := ParseLocation("chr1:100000-100100")

	if err != nil {
		t.Fatalf(`err %s`, err)
	}

	dna, err := GetDNA("/ifs/scratch/cancer/Lab_RDF/ngs/dna/hg19", location, false, false)

	if err != nil {
		t.Fatalf(`err %s`, err)
	}

	out, err := json.Marshal(dna)
	if err != nil {
		t.Fatalf(`err %s`, err)
	}

	fmt.Println(string(out))
}
