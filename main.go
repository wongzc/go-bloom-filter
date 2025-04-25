package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"go-bloom-filter/bloomfilter" 
)

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

	f := bloomfilter.NewFilter(item_count, accuracy)
	fmt.Printf("Using %d hash functions and %d bytes.\n", f.HashFunctionCount, f.ArraySize)

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
