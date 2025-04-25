package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"math"
	"os"
	"strings"
)

func hash1(data string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(data))
	return h.Sum32()
}

func hash2(data string) uint32 {
	h := fnv.New32()
	h.Write([]byte(data))
	return h.Sum32()
}

func getHashes(data string, hash_function_count int, array_size uint32) []uint32 { // using double hashing with the data and index
	hashes := make([]uint32, hash_function_count)
	h1 := hash1(data)
	h2 := hash2(data)

	for i := 0; i < hash_function_count; i++ {
		combined := (h1 + uint32(i)*h2) % array_size
		hashes[i] = combined
	}
	return hashes
}

type filter struct {
	bitfield            []bool
	hash_function_count int
	array_size          uint32
}

func (f *filter) Set(s string) {
	hs := getHashes(s, f.hash_function_count, f.array_size)
	for _, pos := range hs {
		f.bitfield[pos] = true
	}

}

func (f *filter) Get(s string) bool {
	hs := getHashes(s, f.hash_function_count, f.array_size)
	for _, pos := range hs {
		if !f.bitfield[pos] {
			return false
		}
	}
	return true
}

func NewFilter(itemCount float64, accuracy float64) *filter {
	// compute array size and has function required based on acceptable false positive and expected item count
	arraySize := uint32(-itemCount*math.Log(accuracy)/math.Pow(math.Log(2), 2)) + 1
	hashCount := int(float64(arraySize)/itemCount*math.Log(2)) + 1

	return &filter{
		bitfield:            make([]bool, arraySize),
		hash_function_count: hashCount,
		array_size:          arraySize,
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var item_count float64
	var accuracy float64
	fmt.Println("Enter estimated item count (integer) and minimum acceptable False positive  (0<x<1)")
	_, err := fmt.Scan(&item_count, &accuracy)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	} else if accuracy <= 0 || accuracy >= 1 {
		fmt.Println("Accuracy must be between 0 and 1.")
		return
	}

	f := NewFilter(item_count, accuracy)
	fmt.Printf("Require %d hash functions and array size of %d.\n", f.hash_function_count, f.array_size)

	for {
		fmt.Println("Enter command [s=set, g=get, x=exit]:")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "x" {
			fmt.Println("Exit...")
			break
		} else if input == "s" || input == "g" {
			fmt.Printf("Enter string to %s: ", map[string]string{"s": "set", "g": "get"}[input])
			str, _ := reader.ReadString('\n')
			str = strings.TrimSpace(str)

			if input == "s" {
				f.Set(str)
			} else {
				fmt.Println(f.Get(str))
			}
		} else {
			fmt.Println("Unknown command. Use s, g or x.")
		}
	}
}
