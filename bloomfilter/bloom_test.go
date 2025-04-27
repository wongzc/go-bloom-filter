package bloomfilter

import (
	"github.com/cespare/xxhash/v2"
	"testing"
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
func TestBloomFilterBasic(t *testing.T) {
	filter := New(1000, 0.01, hash1, hash2)

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
