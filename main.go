package main

import (
	"bufio"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"math"
	"os"
	"strings"
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

func getHashes(data string, hashFunctionCount int, arraySize uint32) []uint32 { // using double hashing with the data and index
	hashes := make([]uint32, hashFunctionCount)
	h1 := hash1(data)
	h2 := hash2(data)

	for i := 0; i < hashFunctionCount; i++ {
		combined := (h1 + uint32(i)*h2) % arraySize
		hashes[i] = combined
	}
	return hashes
}

type Filter struct {
	bitField            []byte
	hashFunctionCount int
	arraySize          uint32
}

func (f *Filter) Set(s string) {
	hs := getHashes(s, f.hashFunctionCount, f.arraySize)
	for _, pos := range hs {
		setBit(f.bitField, pos)
	}

}

func (f *Filter) Get(s string) bool {
	hs := getHashes(s, f.hashFunctionCount, f.arraySize)
	for _, pos := range hs {
		if !getBit(f.bitField, pos) {
			return false
		}
	}
	return true
}

func NewFilter(itemCount float64, accuracy float64) *Filter {
	// compute array size and has function required based on acceptable false positive and expected item count
	arraySize := uint32(-itemCount*math.Log(accuracy)/math.Pow(math.Log(2), 2))/8 + 1
	hashCount := int(float64(arraySize)/itemCount*math.Log(2)) + 1

	return &Filter{
		bitField:            make([]byte, arraySize),
		hashFunctionCount: hashCount,
		arraySize:          arraySize,
	}
}

func setBit(bitField []byte, pos uint32) {
	byteIndex := pos / 8
	bitOffset := pos % 8
	bitField[byteIndex] |= 1 << bitOffset
}

func getBit(bitField []byte, pos uint32) bool {
	byteIndex := pos / 8
	bitOffset := pos % 8
	return (bitField[byteIndex] & (1 << bitOffset)) != 0
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
	fmt.Printf("Require %d hash functions and byte array size of %d.\n", f.hashFunctionCount, f.arraySize)

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
