package bloomfilter

import (
	"github.com/cespare/xxhash/v2"
	"math"
	"fmt"
)

func getHashes(data string, k int, m uint32) []uint32 {
	hashes := make([]uint32, k)
	for i := 0; i < k; i++ {
		salted := fmt.Sprintf("%s|%d", data, i)
		hashes[i] = uint32(xxhash.Sum64String(salted) % uint64(m))
	}
	return hashes
}

type Filter struct {
	BitField          []byte
	HashFunctionCount int
	ArraySize         uint32
	ElementCount      int
}

func (f *Filter) Set(s string) {
	hs := getHashes(s, f.HashFunctionCount, f.ArraySize)
	for _, pos := range hs {
		setBit(f.BitField, pos)
	}
	f.ElementCount++

}

func (f *Filter) Get(s string) bool {
	hs := getHashes(s, f.HashFunctionCount, f.ArraySize)
	for _, pos := range hs {
		if !getBit(f.BitField, pos) {
			return false
		}
	}
	return true
}

func (f *Filter) CalFPR() float64 {
	m := float64(f.ArraySize * 8) // total bits
	k := float64(f.HashFunctionCount)
	n := float64(f.ElementCount)

	FPR := math.Pow(1-math.Exp(-k*n/m), k)
	return FPR * 100
}
func NewFilter(itemCount float64, accuracy float64) *Filter {
	// compute array size and has function required based on acceptable false positive and expected item count
	ArraySize := uint32(-itemCount*math.Log(accuracy)/math.Pow(math.Log(2), 2)) + 1
	hashCount := int(float64(ArraySize)/itemCount*math.Log(2)) + 1
	byteSize := ArraySize/8 + 1 // convert to byte here

	return &Filter{
		BitField:          make([]byte, byteSize),
		HashFunctionCount: hashCount,
		ArraySize:         byteSize,
		ElementCount:      0,
	}
}

func setBit(BitField []byte, pos uint32) {
	byteIndex := pos / 8
	bitOffset := pos % 8
	BitField[byteIndex] |= 1 << bitOffset
}

func getBit(BitField []byte, pos uint32) bool {
	byteIndex := pos / 8
	bitOffset := pos % 8
	return (BitField[byteIndex] & (1 << bitOffset)) != 0
}
