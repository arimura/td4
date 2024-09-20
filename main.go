package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type register int

const (
	a register = iota
	b
)

//Instruction format
//    MSB <----------------------> LSB
//bit| 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
//   |Op Cod e       |Immediate data |

type add struct {
	// Register name
	register register
	// Immediate binary data, e.g., "0010"
	im string
}

func (x *add) bin() {
	//   MSB <-> LSB
	//ADD A: 0000
	//ADD B: 0101

}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No file name")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
