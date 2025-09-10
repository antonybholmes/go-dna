package dna

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/antonybholmes/go-basemath"
)

const (
	STRAND_POS  = "+"
	STRAND_NEG  = "-"
	STRAND_NONE = "."
)

type PromoterRegion struct {
	offset5p uint
	offset3p uint
}

var DEFAULT_PROMOTER_REGION = PromoterRegion{2000, 1000}

func DefaultPromoterRegion() PromoterRegion {
	return DEFAULT_PROMOTER_REGION
}

func NewPromoterRegion(offset5p uint, offset3p uint) *PromoterRegion {
	return &PromoterRegion{offset5p, offset3p}
}

func (promoterRegion *PromoterRegion) Offset5P() uint {
	return promoterRegion.offset5p
}

func (promoterRegion *PromoterRegion) Offset3P() uint {
	return promoterRegion.offset3p
}

type Location struct {
	Chr    string `json:"chr"`
	Strand string `json:"strand,omitempty"`
	Start  uint   `json:"start"`
	End    uint   `json:"end"`
}

func NewLocation(chr string, start uint, end uint) *Location {
	return NewStrandedLocation(chr, start, end, STRAND_NONE)
}

func NewStrandedLocation(chr string, start uint, end uint, strand string) *Location {

	// standardize chromosome names so that letters other than chr are capitalized
	// e.g. chr1, chr2, ..., chrX, chrY, chrM
	// This is to ensure that the chromosome names are consistent
	// and can be easily parsed and compared.
	chr = strings.Replace(strings.ToUpper(chr), "CHR", "chr", 1)

	if !strings.HasPrefix(chr, "chr") {
		chr = fmt.Sprintf("chr%s", chr)
	}

	s := basemath.Max(1, basemath.Min(start, end))

	return &Location{Chr: chr, Start: s, End: basemath.Max(s, end), Strand: strand}
}

func (location *Location) String() string {
	return fmt.Sprintf("%s:%d-%d", location.Chr, location.Start, location.End)
}

// func (location *Location) MarshalJSON() ([]byte, error) {
// 	// Customize the JSON output here
// 	return []byte(location.String()), nil
// }

// func (location *Location) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(location.String())
// }

func (location *Location) Mid() uint {
	return uint((location.Start + location.End) / 2)
}

func (location *Location) Len() uint {
	return uint(location.End - location.Start + 1)
}

func ParseLocation(location string) (*Location, error) {
	matched, err := regexp.MatchString(`^chr([0-9]+|[xyXY]):\d+-\d+$`, location)

	if !matched || err != nil {
		return nil, fmt.Errorf("%s does not seem like a valid location", location)
	}

	tokens := strings.Split(location, ":")
	chr := tokens[0]
	tokens = strings.Split(tokens[1], "-")

	start, err := strconv.ParseInt(tokens[0], 10, 32)

	if err != nil {
		return nil, fmt.Errorf("%s does not seem like a valid start", tokens[0])
	}

	end, err := strconv.ParseInt(tokens[1], 10, 32)

	if err != nil {
		return nil, fmt.Errorf("%s does not seem like a valid end", tokens[1])
	}

	return NewLocation(chr, uint(start), uint(end)), nil
}

func ParseLocations(locations []string) ([]*Location, error) {
	ret := make([]*Location, 0, len(locations))

	for _, l := range locations {
		loc, err := ParseLocation(l)

		if err != nil {
			return nil, err
		}

		ret = append(ret, loc)
	}

	return ret, nil
}

// Convert a chromosome string to a number suitable for sorting
// These numbers are to ensure a sort order and do not necessarily
// correspond to conventions, for example chrX is often represented
// as 23, but to allow for more chromosomes we use 1000.
func ChromToInt(chr string) uint {
	chr = strings.TrimPrefix(strings.ToLower(chr), "chr")
	switch chr {
	case "x":
		return 1000 //23
	case "y":
		return 2000 //24
	case "m", "mt":
		return 3000 //25
	default:
		n, err := strconv.Atoi(chr)
		if err != nil {
			return 9999 // // Put unknown chromosomes last
		}
		return uint(n)
	}
}

func SortLocations(locations []*Location) {
	slices.SortFunc(locations, func(a, b *Location) int {
		ci := ChromToInt(a.Chr)
		cj := ChromToInt(b.Chr)

		// on different chrs so sort by chr
		if ci != cj {
			return int(ci) - int(cj)
		}

		return int(a.Start) - int(b.Start)
	})
}

// // -------- Position-based sorter --------
// type SortLocByPos []*Location

// func (locations SortLocByPos) Len() int      { return len(locations) }
// func (locations SortLocByPos) Swap(i, j int) { locations[i], locations[j] = locations[j], locations[i] }
// func (locations SortLocByPos) Less(i, j int) bool {
// 	ci := ChromToInt(locations[i].Chr)
// 	cj := ChromToInt(locations[j].Chr)

// 	// on different chrs so sort by chr
// 	if ci != cj {
// 		return ci < cj
// 	}

// 	return locations[i].Start < locations[j].Start
// }
