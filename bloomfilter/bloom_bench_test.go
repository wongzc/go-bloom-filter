package bloomfilter

import (
	"strconv"
	"testing"
	"fmt"
)

func BenchmarkInsert(b *testing.B) {
	filter := New(float64(b.N), 0.01, hash1, hash2)
	for i := 0; i < b.N; i++ {
		filter.Set("item_" + strconv.Itoa(i))
	}
}

func BenchmarkLookup(b *testing.B) {
	filter := New(float64(b.N), 0.01, hash1, hash2)
	for i := 0; i < b.N; i++ {
		filter.Set("item_" + strconv.Itoa(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Get("item_" + strconv.Itoa(i))
	}
}

func BenchmarkMemoryUsage(b *testing.B) {
	filter := New(float64(b.N), 0.01, hash1, hash2)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		filter.Set("item_" + strconv.Itoa(i))
	}
}

func BenchmarkFalsePositiveRate(b *testing.B) {
	n := 1_000_000
	filter := New(float64(n), 0.01, hash1, hash2)

	for i:=0;i<n;i++ {
		filter.Set("exists_"+ strconv.Itoa(i))
	}

	b.ResetTimer()
	falsePositiveCount:=0
	for i:=0;i<n;i++ {
		if filter.Get("nonexists_"+ strconv.Itoa(i)) {
			falsePositiveCount++
		}
	}
	fmt.Printf("False Positives: %d out of %d\n", falsePositiveCount, b.N)
}
