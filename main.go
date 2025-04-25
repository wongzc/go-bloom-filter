package main

import (
	"crypto/sha1"
	"fmt"
	"bufio"
	"os"
	"strings"
	"hash/fnv"
	"math"
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
	bitfield []bool
	hash_function_count int
	array_size uint32
}

func (f *filter)set(s string) {
	hs:= getHashes(s, f.hash_function_count, f.array_size)
	for _,pos := range hs {
		f.bitfield[pos]=true
	}
	
}

func ( f *filter)get(s string) bool {
	hs:= getHashes(s, f.hash_function_count, f.array_size)
	for _,pos := range hs {
		if !f.bitfield[pos] {
			return false
		}
	}
	return true
}

var (
	hasher=sha1.New()
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var item_count float64
	var accuracy float64
	fmt.Println("Enter estimated item count (integer) and minimum acceptable False positive  (0<x<1)")
	_, err:=fmt.Scan(&item_count, &accuracy)
	if err !=nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// compute array size and has function required based on acceptable false positive and expected item count
	array_s:=uint32(-item_count*math.Log(accuracy)/math.Pow(math.Log(2),2))+1 
	hash_c:=int(float64(array_s)/item_count*math.Log(2))+1
	fmt.Printf("require %d hash functions and array size of %d.", hash_c, array_s)
	
	f:=filter{
		bitfield: make([]bool, array_s),
	}

	f.array_size=array_s
	f.hash_function_count=hash_c
	
	for {
		fmt.Println("Enter:")
		input,_ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input =="x" {
			fmt.Println("Exit...")
			break
		} else if input =="s" {
			fmt.Println("Setting...")
			input,_=reader.ReadString('\n')
			input = strings.TrimSpace(input)
			f.set(input)
		} else if input =="g" {
			fmt.Println("Getting...")
			input,_=reader.ReadString('\n')
			input = strings.TrimSpace(input)
			fmt.Println(f.get(input)) 
		}
	}
}
