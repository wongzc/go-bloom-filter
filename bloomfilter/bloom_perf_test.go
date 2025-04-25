package bloomfilter

import (
	"fmt"
	"strconv"
	"testing"
	"math"
	"time"
)

func BenchmarkInsert(b *testing.B) {
	filter := NewFilter(float64(b.N), 0.01)
	for i := 0; i < b.N; i++ {
		filter.Set("item_" + strconv.Itoa(i))
	}
}

func BenchmarkLookup(b *testing.B) {
	filter := NewFilter(float64(b.N), 0.01)
	for i := 0; i < b.N; i++ {
		filter.Set("item_" + strconv.Itoa(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Get("item_" + strconv.Itoa(i))
	}
}

func TestFPRHighLoad(t *testing.T) {
	n := 1_000_000
	filter := NewFilter(float64(n), 0.001)

	// Insert n unique items
	for i := 0; i < n; i++ {
		filter.Set("data_" + strconv.Itoa(i))
	}

	fmt.Printf("Bit saturation: %.4f%%\n", filter.BitSaturation())

	// Check for false positives using unseen keys
	falsePositives := 0
	trials := 700_000
	for i := 0; i < trials; i++ {
		unseen := "unseen_" + strconv.Itoa(i)
		if filter.Get(unseen) {
			falsePositives++
		}
	}

	actualFPR := float64(falsePositives) / float64(trials) * 100
	fmt.Printf("High Load FPR: %.4f%% (expected ~%.2f%%)\n", actualFPR, filter.CalFPR())

	if math.Abs(actualFPR-filter.CalFPR()) > 0.03 { 
		t.Errorf("False positive rate too high: %.2f%%", actualFPR)
	}
}

func TestStressInsert(t *testing.T) {
	start := time.Now()

	n := 10_000_000
	filter := NewFilter(float64(n), 0.01)

	for i := 0; i < n; i++ {
		filter.Set("stress_" + strconv.Itoa(i))
	}

	duration := time.Since(start)
	fmt.Printf("Inserted %d elements in %s\n", n, duration)
}

// go test -v ./bloomfilter   
// go test -v -run ^TestFPRHighLoad$ ./bloomfilter
// go test -v -bench=. ./bloomfilter
// go test -v -run ^TestStressInsert$ ./bloomfilter