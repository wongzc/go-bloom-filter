package bloomfilter

import (
	"testing"
)

func TestBloomFilterBasic(t *testing.T) {
	filter := NewFilter(1000, 0.01)

	// Insert some elements
	filter.Set("apple")
	filter.Set("banana")
	filter.Set("cherry")

	// Should return true (possibly in set)
	if !filter.Get("apple") {
		t.Error("Expected 'apple' to be in set")
	}
	if !filter.Get("banana") {
		t.Error("Expected 'banana' to be in set")
	}

	// Should return false (definitely not in set)
	if filter.Get("durian") {
		t.Error("Did not expect 'durian' to be in set")
	}
}

// run with go test ./bloomfilter