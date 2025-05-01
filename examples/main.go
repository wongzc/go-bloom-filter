package main

import (
	"bufio"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"github.com/wongzc/go-bloom-filter/bloomfilter"
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

// usage example
func main() {
	reader := bufio.NewReader(os.Stdin)
	var item_count float64
	var accuracy float64
	fmt.Println("Enter estimated item count (integer) and acceptable false positive rate (0<x<1)")
	_, err := fmt.Scan(&item_count, &accuracy)
	if err != nil || accuracy <= 0 || accuracy >= 1 {
		fmt.Println("Invalid input")
		return
	}

	f := bloomfilter.New(item_count, accuracy, hash1, hash2)
	defer f.Close()
	
	fmt.Printf("Using %d hash functions and %d bytes.\n", f.HashFunctionCount, f.ArraySize/8+1)

	for {
		fmt.Print("[s=set, g=get, x=exit]: ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		if cmd == "x" {
			break
		} else if cmd == "s" || cmd == "g" {
			fmt.Print("Enter string: ")
			str, _ := reader.ReadString('\n')
			str = strings.TrimSpace(str)

			if cmd == "s" {
				f.Set(str)
			} else if cmd == "g" {
				if f.Get(str) {
					fmt.Printf("Possibly in set with %.2f%% confidence\n", 100.0-f.CalFPR())
				} else {
					fmt.Println("Definitely not in set")
				}
			}
		}
	}
}
