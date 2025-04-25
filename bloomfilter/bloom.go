package bloomfilter

import (
	"github.com/cespare/xxhash/v2"
	"math"
)

func hash1(data string) uint32 {
	return uint32(xxhash.Sum64String(data))
}

func hash2(data string) uint32 {
	h := xxhash.New()
	h.Write([]byte(data + "salt")) // change a bit to get a different hash

	val := h.Sum64() // to avoid h2 = 0, which may cause all the value in hashes slice become the same
	if val == 0 {
		val = 1
	}
	return uint32(val)
}

func getHashes(data string, HashFunctionCount int, ArraySize uint32) []uint32 { // using double hashing with the data and index
	hashes := make([]uint32, HashFunctionCount)
	h1 := hash1(data)
	h2 := hash2(data)

	for i := 0; i < HashFunctionCount; i++ {
		combined := (h1 + uint32(i)*h2) % ArraySize
		hashes[i] = combined
	}
	return hashes
}

type Filter struct {
	BitField          	[]byte
	HashFunctionCount 	int
	ArraySize         	uint32
	ElementCount		int
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

func (f *Filter)CalFPR() float64 {
	FPR:=math.Pow(1-math.Exp(-float64(f.HashFunctionCount)*float64(f.ElementCount)/float64(f.ArraySize)/8.0),float64(f.HashFunctionCount))
	return (1-FPR)*100
}
func NewFilter(itemCount float64, accuracy float64) *Filter {
	// compute array size and has function required based on acceptable false positive and expected item count
	ArraySize := uint32(-itemCount*math.Log(accuracy)/math.Pow(math.Log(2), 2)) + 1
	hashCount := int(float64(ArraySize)/itemCount*math.Log(2)) + 1
	byteSize := ArraySize/8+1 // convert to byte here

	return &Filter{
		BitField:            make([]byte, byteSize),
		HashFunctionCount: 	hashCount,
		ArraySize:          byteSize,
		ElementCount:		0,
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