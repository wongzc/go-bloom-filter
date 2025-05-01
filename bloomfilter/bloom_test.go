package bloomfilter

import (
	"github.com/cespare/xxhash/v2"
	"testing"
	"time"
	"fmt"
	"github.com/wongzc/go-bloom-filter/randomstring"
	"math"
	"strconv"
)

func hash1(data string) uint32 {
	return uint32(xxhash.Sum64String(data))
}

func hash2(data string) uint32 {
	h := xxhash.New()
	h.Write([]byte(data + "salt")) 

	val := h.Sum64()
	if val == 0 {
		val = 1
	}
	return uint32(val)
}
func TestBloomFilterBasic(t *testing.T) {
	filter := New(1000, 0.01, hash1, hash2)

	filter.Set("apple")
	filter.Set("banana")
	filter.Set("cherry")

	time.Sleep(150*time.Millisecond) // give background setter some time

	if !filter.Get("apple") {
		t.Error("Expected 'apple' to be in set")
	}
	if !filter.Get("banana") {
		t.Error("Expected 'banana' to be in set")
	}

	if filter.Get("durian") {
		t.Error("Did not expect 'durian' to be in set")
	}
}

func TestFPRHighLoad(t *testing.T) {
	n := 1_000_000
	filter := New(float64(n), 0.001, hash1, hash2)

	// Insert n unique items
	for i := 0; i < n; i++ {
		filter.Set("data_" + strconv.Itoa(i) + randomstring.RandString())
	}

	fmt.Printf("\nHash Functions Count: %d\n", filter.HashFunctionCount)
	fmt.Printf("Array Size: %d\n", filter.ArraySize)
	fmt.Printf("Bit Saturation Rate: %.4f%%\n", filter.BitSaturation())
	fmt.Printf("Bit Distribution: %f\n", filter.BitDistribution())

	// Check for false positives using unseen keys
	falsePositives := 0
	trials := 1_000_000
	for i := 0; i < trials; i++ {
		unseen := "unseen_" + strconv.Itoa(i) + randomstring.RandString()
		if filter.Get(unseen) {
			falsePositives++
		}
	}

	actualFPR := float64(falsePositives) / float64(trials) * 100
	fmt.Printf("False Positive Rate: %.4f%% (expected ~%.2f%%)\n", actualFPR, filter.CalFPR())

	if math.Abs(actualFPR-filter.CalFPR()) > 0.03 {
		t.Errorf("False positive rate too high: %.2f%%", actualFPR)
	}

	filter.PrintRandomBitHeatmap(1000, 100)
}

func TestStressInsert(t *testing.T) {
	start := time.Now()

	n := 10_000_000
	filter := New(float64(n), 0.01, hash1, hash2)

	for i := 0; i < n; i++ {
		filter.Set("stress_" + strconv.Itoa(i))
	}

	duration := time.Since(start)
	fmt.Printf("Inserted %d elements in %s\n", n, duration)
}


