package main

import (
	"crypto/sha1"
	"fmt"
	"bufio"
	"os"
	"strings"
	"hash/fnv"
)

const (
	array_size uint32 = 1000
	hash_function = 4
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

func getHashes(data string) []uint32 { // using double hashing with the data and index
    hashes := make([]uint32, hash_function)
    h1 := hash1(data)
    h2 := hash2(data)

    for i := 0; i < hash_function; i++ {
        combined := (h1 + uint32(i)*h2) % array_size
        hashes[i] = combined
    }
    return hashes
}

type filter struct {
	bitfield [array_size]bool
}

func (f *filter)set(s string) {
	hs:= getHashes(s)
	for _,pos := range hs {
		f.bitfield[pos]=true
	}
	
}

func ( f *filter)get(s string) bool {
	hs:= getHashes(s)
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
	f:=filter{}
	reader := bufio.NewReader(os.Stdin)
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
