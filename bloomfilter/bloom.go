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

// func getHashes(data string, k int, m uint32) []uint32 {
// 	hashes := make([]uint32, k)
// 	for i := 0; i < k; i++ {
// 		salted := fmt.Sprintf("%s|%d", data, i)
// 		hashes[i] = uint32(xxhash.Sum64String(salted) % uint64(m))
// 	}
// 	return hashes
// }

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
	m := float64(f.ArraySize)
	k := float64(f.HashFunctionCount)
	n := float64(f.ElementCount)

	FPR := math.Pow(1-math.Exp(-k*n/m), k)
	return FPR * 100
}

func (f *Filter) BitSaturation() float64 {
	count := 0
	for _, b := range f.BitField {
		for i := 0; i < 8; i++ {
			if b&(1<<i) != 0 {
				count++
			}
		}
	}
	return float64(count) / float64(len(f.BitField)*8) * 100
}

func (f *Filter) BitDistribution() float64 {
	byteBitCounts := make([]int, len(f.BitField))
	for i, b := range f.BitField {
		for j := 0; j < 8; j++ {
			if b&(1<<j) != 0 {
				byteBitCounts[i]++
			}
		}
	}

	sum := 0
	for _, c := range byteBitCounts {
		sum += c
	}
	mean := float64(sum) / float64(len(byteBitCounts))

	var varianceSum float64
	for _, c := range byteBitCounts {
		diff := float64(c) - mean
		varianceSum += math.Pow(diff, 2)
	}
	variance := varianceSum / float64(len(f.BitField))

	return variance
}

func NewFilter(itemCount float64, accuracy float64) *Filter {
	// compute array size and has function required based on acceptable false positive and expected item count
	ArraySize := uint32(-itemCount*math.Log(accuracy)/math.Pow(math.Log(2), 2)) + 1
	hashCount := int(float64(ArraySize)/itemCount*math.Log(2)) + 1
	byteArraySize := ArraySize/8 + 1 // convert to byte here

	return &Filter{
		BitField:          make([]byte, byteArraySize),
		HashFunctionCount: hashCount,
		ArraySize:         ArraySize,
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
