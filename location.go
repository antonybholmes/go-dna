package dna

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/antonybholmes/go-basemath"
	"github.com/antonybholmes/go-sys"
)

type (
	PromoterRegion struct {
		upstream   int
		downstream int
	}

	Location struct {
		chr    string
		strand string
		start  int
		end    int
	}

	jsonLocation struct {
		Chr    string `json:"chr"`
		Strand string `json:"strand,omitempty"`
		Start  int    `json:"start"`
		End    int    `json:"end"`
	}
)

const (
	PosStrand      = "+"
	NegStrand      = "-"
	StrandNotGiven = "."
)

var (
	defaultPromoterRegion = NewPromoterRegion(2000, 1000)
	//chrRegex              = regexp.MustCompile(`(?i)(?:chr)?([0-9]+|[a-z_]+)`)
	locRegex = regexp.MustCompile(`(?i)(?:chr)?([0-9]+|[a-z_]+):([0-9,]+)-([0-9,]+)`)
)

func DefaultPromoterRegion() *PromoterRegion {
	// once.Do(func() {
	// 	defaultPromoterRegion = NewPromoterRegion(2000, 1000)
	// })
	return defaultPromoterRegion

	//return NewPromoterRegion(2000, 1000)
}

func NewPromoterRegion(upstream int, downstream int) *PromoterRegion {
	return &PromoterRegion{basemath.AbsInt(upstream), basemath.AbsInt(downstream)}
}

func (promoterRegion *PromoterRegion) Upstream() int {
	return promoterRegion.upstream
}

func (promoterRegion *PromoterRegion) Downstream() int {
	return promoterRegion.downstream
}

func NewLocation(chr string, start int, end int) (*Location, error) {
	return NewStrandedLocation(chr, start, end, StrandNotGiven)
}

func NewStrandedLocation(chr string, start int, end int, strand string) (*Location, error) {
	// standardize chromosome names so that letters other than chr are capitalized
	// e.g. chr1, chr2, ..., chrX, chrY, chrM
	// This is to ensure that the chromosome names are consistent
	// and can be easily parsed and compared.
	chr, err := ParseChr(chr)

	if err != nil {
		return nil, fmt.Errorf("invalid chromosome name: %s", chr)
	}

	// limit strand values
	strand = ParseStrand(strand)

	loc := ParseStartEnd(start, end)

	return &Location{chr: chr, start: loc[0], end: loc[1], strand: strand}, nil
}

// Returns the base chromosome string without the "chr" prefix.
func (location *Location) BaseChr() string {
	return strings.TrimPrefix(location.chr, "chr")
}

// Returns the chromosome string with the "chr" prefix.
func (location *Location) Chr() string {
	return location.chr
}

func (location *Location) Start() int {
	return location.start
}

func (location *Location) End() int {
	return location.end
}

func (location *Location) Strand() string {
	return location.strand
}

// Returns the string representation of the location in the format "chrX:start-end".
func (location *Location) String() string {
	return location.chr + ":" + strconv.Itoa(location.start) + "-" + strconv.Itoa(location.end)

	//fmt.Sprintf("%s:%d-%d", location.Chr, location.Start, location.End)
}

// func (location *Location) MarshalJSON() ([]byte, error) {
// 	// Customize the JSON output here
// 	return []byte(location.String()), nil
// }

// func (location *Location) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(location.String())
// }

// Returns the midpoint of the location rounded down to the nearest integer.
func (location *Location) Mid() int {
	return (location.start + location.end) / 2
}

// Returns the length of the location.
func (location *Location) Len() int {
	return location.end - location.start + 1
}

// Create a JSON representation of the coordinate
func (location *Location) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonLocation{
		Chr:    location.chr,
		Strand: location.strand,
		Start:  location.start,
		End:    location.end,
	})
}

func (location *Location) UnmarshalJSON(data []byte) error {
	var jl jsonLocation

	if err := json.Unmarshal(data, &jl); err != nil {
		return err
	}

	// Still parse it so we are not loading complete rubbish

	chr, err := ParseChr(jl.Chr)

	if err != nil {
		return err
	}

	loc := ParseStartEnd(jl.Start, jl.End)

	location.chr = chr
	location.strand = ParseStrand(jl.Strand)
	location.start = loc[0]
	location.end = loc[1]

	return nil
}

func ParseChr(location string) (string, error) {
	location = strings.ToUpper(strings.TrimSpace(location))

	// remove any prefix
	location = strings.TrimPrefix(location, "CHR")

	// should test if remaining is either all digits or a known letter
	// loop over string to check

	// digitCount := 0
	// letterCount := 0

	// for _, c := range location {
	// 	if unicode.IsDigit(c) {
	// 		digitCount++
	// 	}

	// 	if unicode.IsLetter(c) {
	// 		letterCount++
	// 	}
	// }

	// // we must either be all letter
	// if letterCount > 0 && digitCount > 0 {
	// 	return "", fmt.Errorf("%s does not seem like a valid chr", location)
	// }

	// matches := chrRegex.FindStringSubmatch(location)

	// if len(matches) < 1 {
	// 	return "", fmt.Errorf("%s does not seem like a valid chr", location)
	// }

	// chr := matches[1]

	return "chr" + location, nil
}

func ParseStrand(strand string) string {
	if strand != "+" && strand != "-" {
		return "."
	}

	return strand
}

func ParseStartEnd(start int, end int) []int {
	start = basemath.AbsInt(start)
	end = basemath.AbsInt(end)

	// s is 1 based and the smaller of the two
	// this is to ensure coordinates are consistent
	// on the forward strand
	s := basemath.Max(1, basemath.Min(start, end))
	e := basemath.Max(s, end)

	return []int{s, e}
}

func ParseLocation(location string) (*Location, error) {

	matches := locRegex.FindStringSubmatch(location)

	if len(matches) < 3 {
		return nil, fmt.Errorf("%s does not seem like a valid location", location)
	}

	chr := matches[1]

	start, err := sys.Atoi(matches[2])

	if err != nil {
		return nil, fmt.Errorf("%s does not seem like a valid start", matches[2])
	}

	if start < 1 {
		return nil, fmt.Errorf("start position %d is less than 1", start)
	}

	end, err := sys.Atoi(matches[3])

	if err != nil {
		return nil, fmt.Errorf("%s does not seem like a valid end", matches[3])
	}

	if end < 1 {
		return nil, fmt.Errorf("end position %d is less than 1", end)
	}

	if end < start {
		return nil, fmt.Errorf("end position %d is less than start position %d", end, start)
	}

	return NewLocation(chr, start, end)
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
// as 23 in human, but we do not presume the species so we use
// 1023 to allow for lots of chromosomes. chrN is converted to N where
// N is an integer. chrX, chrY, chrM are converted to large numbers to
// ensure they sort after numbered chromosomes. Unknown chromosomes
// are converted to 9999 to ensure they sort last.
func ChromToInt(chr string) int {
	chr = strings.TrimPrefix(strings.ToLower(chr), "chr")

	switch chr {
	case "x":
		return 1023 //23
	case "y":
		return 1024 //24
	case "m", "mt":
		return 1025 //25
	default:
		n, err := sys.Atoi(chr)
		if err != nil {
			return 9999 // // Put unknown chromosomes last
		}
		return n
	}
}

// SortLocations sorts locations in place by chr, start, end
func SortLocations(locations []*Location) {
	slices.SortFunc(locations, func(a, b *Location) int {
		ci := ChromToInt(a.Chr())
		cj := ChromToInt(b.Chr())

		// on different chrs so sort by chr
		if ci != cj {
			return int(ci) - int(cj)
		}

		// same chr so sort by start
		diff := a.Start() - b.Start()

		if diff != 0 {
			return diff
		}

		// same start so sort by end
		return a.End() - b.End()
	})
}

// SortLocationsFunc is a comparison function for sorting locations by chr, start, end
// suitable for use with slices.SortFunc
func SortLocationsFunc(a, b *Location) bool {

	ci := ChromToInt(a.Chr())
	cj := ChromToInt(b.Chr())

	// on different chrs so sort by chr
	if ci != cj {
		return ci < cj
	}

	if a.Start()-b.Start() != 0 {
		return a.Start() < b.Start()
	}

	// same start so sort by end
	return a.End() < b.End()
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
