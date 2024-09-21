package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type register int

const (
	none register = iota
	a
	b
)

//Instruction format
//    MSB <----------------------> LSB
//bit| 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
//   |Op Code        |Immediate data |

type instruction interface {
	gen()
}

type add struct {
	// Register name
	register register
	// Immediate binary data, e.g., "0010"
	im string
}

func (x *add) gen() {
	//   MSB <-> LSB
	//ADD A: 0000
	//ADD B: 0101
	var s uint8
	if x.register == a {
		s = 0b0000
	} else if x.register == b {
		s = 0b0101
	} else {
		log.Fatal("Invalid register")
	}
	s = appendIm(s, x.im)

	print(s)
}

type out struct {
	regiser register
	im      string
}

func (x *out) gen() {
	//OUT Im: 1011
	//OUT B: 1001
	var s uint8
	if x.regiser == none {
		s = 0b1011
		s = appendIm(s, x.im)
	} else if x.regiser == b {
		s = 0b10110000
	} else {
		log.Fatalf("Invalid register")
	}
	print(s)
}

func appendIm(bin uint8, im string) uint8 {
	for _, c := range im {
		bin = (bin << 1) | uint8(c-'0')
	}
	return bin
}

func print(i uint8) {
	if *binaryOption {
		fmt.Printf("%08b\n", i)
	} else {
		fmt.Printf("%X\n", i)
	}
}

func consume(str *string, keyword string) bool {
	if strings.HasPrefix(*str, keyword) {
		*str = strings.TrimPrefix(*str, keyword)
		return true
	}
	return false
}

func consumeBinary(str *string) string {
	var b []rune
	for i, c := range *str {
		if c == '1' || c == '0' {
			b = append(b, c)
		} else {
			*str = (*str)[i:]
		}
	}
	return string(b)
}

func gen(l string) {
	var i instruction
	if consume(&l, "ADD ") {
		var addi add
		if consume(&l, "A, ") {
			addi = add{register: a}
		} else if consume(&l, "B, ") {
			addi = add{register: b}
		} else {
			log.Fatal("Unsupported Operand")
		}
		im := consumeBinary(&l)
		if im == "" {
			log.Fatal("Unsupported Immediate data")
		}
		addi.im = im
		i = &addi
	} else if consume(&l, "OUT ") {
		var outi out
		if consume(&l, " B") {
			outi = out{regiser: b}
		} else if im := consumeBinary(&l); im != "" {
			outi = out{im: im}
		} else {
			log.Fatal("Unsupported Op code")
		}
		i = &outi
	} else {
		log.Fatal("Unsupported Op Code")
	}
	i.gen()
}

var binaryOption *bool

func main() {
	binaryOption = flag.Bool("b", false, "Specify output as binary format")
	ff := flag.String("f", "", "Specify assembly")
	flag.Parse()

	f, err := os.Open(*ff)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		gen(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
