package dna

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/antonybholmes/go-utils"
)

// var DNA_4BIT_DECODE_MAP = map[byte]byte{
// 	1:  65,
// 	2:  67,
// 	3:  71,
// 	4:  84,
// 	5:  97,
// 	6:  99,
// 	7:  103,
// 	8:  116,
// 	9:  78,
// 	10: 110,
// }

// use an array for speed since
// we only have 16 values and we
// know explicitly what each value
// maps to
var DNA_4BIT_DECODE_MAP = [16]byte{
	0,
	65,
	67,
	71,
	84,
	97,
	99,
	103,
	116,
	78,
	110,
	0,
	0,
	0,
	0,
	0,
}

// This is simple complementary lookup
// map for DNA bases represented as
// ASCII code bytes, e.g. 65 = 'A' maps
// to 84 = 'T'
var DNA_COMPLEMENT_MAP = map[byte]byte{
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

type TSSRegion struct {
	offset5p int
	offset3p int
}

func NewTSSRegion(offset5p int, offset3p int) *TSSRegion {
	return &TSSRegion{
		offset5p: utils.AbsInt(offset5p),
		offset3p: utils.AbsInt(offset3p),
	}
}

func (t *TSSRegion) Offset5P() int {
	return t.offset5p
}

func (t *TSSRegion) Offset3P() int {
	return t.offset3p
}

type Location struct {
	Chr   string `json:"chr"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

func NewLocation(chr string, start int, end int) (*Location, error) {
	chr = strings.ToLower(chr)

	if !strings.Contains(chr, "chr") {
		return nil, fmt.Errorf("chr %s is invalid", chr)
	}

	s := utils.IntMax(1, utils.IntMin(start, end))

	return &Location{
		Chr:   chr,
		Start: s,
		End:   utils.IntMax(s, end),
	}, nil
}

func (location *Location) Mid() int {
	return (location.Start + location.End) / 2
}

func (location *Location) String() string {
	return fmt.Sprintf("%s:%d-%d", location.Chr, location.Start, location.End)
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

func Rev(dna []byte) {
	l := len(dna)
	lastIndex := l - 1

	for i := 0; i < l/2; i++ {
		j := lastIndex - i
		dna[i], dna[j] = dna[j], dna[i]
	}
}

func Comp(dna []byte) {

	for i, v := range dna {
		dna[i] = DNA_COMPLEMENT_MAP[v]
	}
}

// Reverse complement a dna byte sequence in situ.
func RevComp(dna []byte) {
	Rev(dna)
	Comp(dna)

	// l := len(dna)
	// lastIndex := l - 1
	// n := l / 2

	// // reverse the byte order and complement each base
	// for i := 0; i < n; i++ {
	// 	i2 := lastIndex - i
	// 	b := DNA_COMPLEMENT_MAP[dna[i]]
	// 	dna[i] = DNA_COMPLEMENT_MAP[dna[i2]]
	// 	dna[i2] = b
	// }
}

func GetDNA(dir string, location *Location, rev bool, comp bool) (string, error) {
	s := location.Start - 1
	e := location.End - 1
	l := e - s + 1
	bs := s / 2
	be := e / 2
	bl := be - bs + 1

	d := make([]byte, bl)

	file := filepath.Join(dir, fmt.Sprintf("%s.dna.4bit", strings.ToLower(location.Chr)))

	f, err := os.Open(file)

	if err != nil {
		return "", err
	}

	f.Seek(int64(1+bs), 0)

	_, err = f.Read(d)

	f.Close()

	if err != nil {
		return "", err
	}

	dna := make([]byte, l)

	// which byte we are scanning (each byte contains 2 bases)
	byteIndex := 0
	var v byte

	for i := 0; i < l; i++ {
		// Which base we want in the byte
		// If the start position s is even, we want the first
		// 4 bits of the byte, else the lower 4 bits.
		baseIndex := s % 2

		v = d[byteIndex]

		if baseIndex == 0 {
			v = v >> 4
		} else {
			// if we are on the second base of the byte, on the
			// next loop we must proceed to the next byte to get
			// the base
			byteIndex++
		}

		// mask for lower 4 bits since these
		// contain the dna base code
		dna[i] = DNA_4BIT_DECODE_MAP[v&15]

		s++
	}

	if rev {
		Rev(dna)
	}

	if comp {
		Comp(dna)
	}

	return string(dna), nil
}
