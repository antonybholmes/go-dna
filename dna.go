package dna

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var DNA_4BIT_DECODE_DICT = map[byte]byte{
	1:  65,
	2:  67,
	3:  71,
	4:  84,
	5:  97,
	6:  99,
	7:  103,
	8:  116,
	9:  78,
	10: 110,
}

var DNA_4BIT_COMP_DICT = map[byte]byte{
	0:   0,
	65:  84,
	67:  71,
	71:  67,
	84:  65,
	97:  116,
	99:  103,
	103: 99,
	116: 97,
	78:  78,
	110: 10,
}

type Location struct {
	Chr   string `json:"chr"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

type DNA struct {
	Location *Location `json:"location"`
	DNA      string    `json:"dna"`
}

func ParseLocation(location string) (*Location, error) {
	matched, err := regexp.MatchString(`^chr([0-9]+|[xyXY]):\d+-\d+$`, location)

	if !matched || err != nil {
		return nil, fmt.Errorf("%s does not seem like a valid location", location)
	}

	tokens := strings.Split(location, ":")
	chr := tokens[0]
	tokens = strings.Split(tokens[1], "-")

	start, err := strconv.Atoi(tokens[0])

	if err != nil {
		return nil, fmt.Errorf("%s does not seem like a valid start", tokens[0])
	}

	end, err := strconv.Atoi(tokens[1])

	if err != nil {
		return nil, fmt.Errorf("%s does not seem like a valid end", tokens[1])
	}

	return &Location{Chr: chr, Start: start, End: end}, nil
}

func RevComp(dna []byte) {

	i2 := len(dna) - 1

	l := int(len(dna) / 2)

	for i := 0; i < l; i++ {
		b := DNA_4BIT_COMP_DICT[dna[i]]
		dna[i] = DNA_4BIT_COMP_DICT[dna[i2]]
		dna[i2] = b
		i2 -= 1
	}
}

func GetDNA(dir string, location *Location) (*DNA, error) {
	s := location.Start - 1
	e := location.End - 1
	l := e - s + 1
	bs := s / 2
	be := e / 2
	bl := be - bs + 1

	d := make([]byte, bl)

	file := filepath.Join(dir, fmt.Sprintf("%s.dna.4bit", strings.ToLower(location.Chr)))

	fmt.Printf("%s\n", file)

	f, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	fmt.Printf("seek %d\n", 1+bs)

	f.Seek(int64(1+bs), 0)

	_, err = f.Read(d)

	f.Close()

	if err != nil {
		return nil, err
	}

	ret := make([]byte, l)

	bi := 0
	var v byte

	for i := 0; i < l; i++ {
		block := s % 2

		v = d[bi]

		if block == 0 {
			v = v >> 4
		}

		ret[i] = DNA_4BIT_DECODE_DICT[v&15]

		if block == 1 {
			bi++
		}

		s++
	}

	return &DNA{Location: location, DNA: string(ret)}, nil
}
