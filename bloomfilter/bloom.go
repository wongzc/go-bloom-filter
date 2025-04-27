package bloomfilter

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
)

const (
	resetColor = "\033[0m"
	redColor   = "\033[31m"
	greenColor = "\033[32m"
)

type Filter struct {
	BitField          []byte
	HashFunctionCount int
	ArraySize         uint32
	ElementCount      int
	HashFunc1         func(string) uint32
	HashFunc2         func(string) uint32
	mu                sync.RWMutex
}

func (f *Filter) getHashes(data string) []uint32 { // using double hashing with the data and index
	hashes := make([]uint32, f.HashFunctionCount)
	h1 := f.HashFunc1(data)
	h2 := f.HashFunc2(data)

	for i := 0; i < f.HashFunctionCount; i++ {
		combined := (h1 + uint32(i)*h2) % f.ArraySize
		hashes[i] = combined
	}
	return hashes
}

func (f *Filter) Set(s string) {
	hs := f.getHashes(s)

	f.mu.Lock()         //lock when setting bit
	defer f.mu.Unlock() // ensure unlock it even panic!!

	for _, pos := range hs {
		setBit(f.BitField, pos)
	}
	f.ElementCount++
}

func (f *Filter) Get(s string) bool {
	hs := f.getHashes(s)

	f.mu.RLock() // read lock it
	defer f.mu.RUnlock()

	for _, pos := range hs {
		if !getBit(f.BitField, pos) {
			return false
		}
	}
	return true
}

func (f *Filter) SetBatch(strings []string) {
	// batch setting method to save overhead from lock unlock
	bitPositions := []uint32{}

	for _, s := range strings {
		hs := f.getHashes(s)
		bitPositions = append(bitPositions, hs...)
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	for _, pos := range bitPositions {
		setBit(f.BitField, pos)
	}
	f.ElementCount += len(strings)
}

func (f *Filter) CalFPR() float64 {
	f.mu.RLock()
	defer f.mu.RUnlock()

	m := float64(f.ArraySize)
	k := float64(f.HashFunctionCount)
	n := float64(f.ElementCount)

	FPR := math.Pow(1-math.Exp(-k*n/m), k)
	return FPR * 100
}

func (f *Filter) BitSaturation() float64 {
	f.mu.RLock()
	defer f.mu.RUnlock()

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
	f.mu.RLock()
	defer f.mu.RUnlock()

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

func (f *Filter) PrintRandomBitHeatmap(sampleSize, columns int) {
	totalBits := len(f.BitField) * 8
	if sampleSize > totalBits {
		sampleSize = totalBits
	}

	startBit := rand.Intn(totalBits - sampleSize + 1)

	fmt.Printf("Colored Heatmap (bits %d to %d):\n", startBit, startBit+sampleSize-1)

	f.mu.RLock()
	defer f.mu.RUnlock()

	for i := startBit; i < startBit+sampleSize; i++ {
		if getBit(f.BitField, uint32(i)) {
			fmt.Print(redColor + "█" + resetColor) // Set bit → red
		} else {
			fmt.Print(greenColor + "·" + resetColor) // Unset bit → green
		}

		if (i-startBit+1)%columns == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func NewFilter(itemCount, accuracy float64, hashFunc1, hashFunc2 func(string) uint32) *Filter {
	// compute array size and hash function required based on acceptable false positive and expected item count
	ArraySize := uint32(-itemCount*math.Log(accuracy)/math.Pow(math.Log(2), 2)) + 1
	hashCount := int(float64(ArraySize)/itemCount*math.Log(2)) + 1
	byteArraySize := ArraySize/8 + 1 // convert to byte here

	return &Filter{
		BitField:          make([]byte, byteArraySize),
		HashFunctionCount: hashCount,
		ArraySize:         ArraySize,
		ElementCount:      0,
		HashFunc1:         hashFunc1,
		HashFunc2:         hashFunc2,
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
