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
		log.Fatal("ADD gen: Invalid register")
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
		log.Fatal("OUT gen: Invalid register")
	}
	print(s)
}

type in struct {
	register register
}

func (x *in) gen() {
	//IN A: 0010
	//IN B: 0110
	var s uint8
	if x.register == a {
		s = 0b0010000
	} else if x.register == b {
		s = 0b0110000
	} else {
		log.Fatal("IN gen: Invalid register")
	}
	print(s)
}

type mov struct {
	register1 register
	register2 register
	im        string
}

func (x *mov) gen() {
	//MOV A, Im: 0011
	//MOV B, Im: 0111
	//MOV A, B:  0001
	//MOV B, A:  0100
	var s uint8
	if x.register1 == a {
		if x.im != "" {
			s = 0b0011
			s = appendIm(s, x.im)
		} else if x.register2 == b {
			s = 0b00010000
		} else {
			log.Fatal("MOV gen: Invalid 2nd register")
		}
	} else if x.register1 == b {
		if x.im != "" {
			s = 0b0111
			s = appendIm(s, x.im)
		} else if x.register2 == a {
			s = 0b0100000
		} else {
			log.Fatal("MOV gen: Invalid 2nd register")
		}
	} else {
		log.Fatal("MOV gen: Invalid 1st register")
	}
	print(s)
}

type jmp struct {
	im string
}

func (x *jmp) gen() {
	//JMP: 1111
	var s uint8 = 0b1111
	s = appendIm(s, x.im)
	print(s)
}

type jnc struct {
	im string
}

func (x *jnc) gen() {
	//JNC: 1110
	var s uint8 = 0b1110
	s = appendIm(s, x.im)
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
		writer.WriteString(fmt.Sprintf(" %08b", i))
	} else {
		writer.WriteString(fmt.Sprintf(" %02x", i))
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
			log.Fatal("ADD consume: Invalid register")
		}
		im := consumeBinary(&l)
		if im == "" {
			log.Fatal("ADD consume: Invalid Immediate data")
		}
		addi.im = im
		i = &addi
	} else if consume(&l, "OUT ") {
		var outi out
		if consume(&l, "B") {
			outi = out{regiser: b}
		} else if im := consumeBinary(&l); im != "" {
			outi = out{im: im}
		} else {
			log.Fatal("OUT consume: Invalid Operand")
		}
		i = &outi
	} else if consume(&l, "IN ") {
		var ini in
		if consume(&l, "A") {
			ini = in{register: a}
		} else if consume(&l, "B") {
			ini = in{register: b}
		} else {
			log.Fatal("IN consume: Invalid Operand")
		}
		i = &ini
	} else if consume(&l, "MOV ") {
		var movi mov
		if consume(&l, "A, ") {
			if consume(&l, "B") {
				movi = mov{register1: a, register2: b}
			} else if im := consumeBinary(&l); im != "" {
				movi = mov{register1: a, im: im}
			} else {
				log.Fatal("MOV consume: Invalid Operand")
			}
		} else if consume(&l, "B, ") {
			if consume(&l, "A") {
				movi = mov{register1: b, register2: a}
			} else if im := consumeBinary(&l); im != "" {
				movi = mov{register1: b, im: im}
			} else {
				log.Fatal()
			}
		} else {
			log.Fatal("MOV consume: Invalid Operand")
		}
		i = &movi
	} else if consume(&l, "JMP ") {
		im := consumeBinary(&l)
		if im == "" {
			log.Fatal("JMP consume: Invalid Immediate data")
		}
		i = &jmp{im: im}
	} else if consume(&l, "JNC ") {
		im := consumeBinary(&l)
		if im == "" {
			log.Fatal("JNC consume: Invalid Immediate data")
		}
		i = &jnc{im: im}
	} else {
		log.Fatal("consume: Invalid instruction")
	}
	i.gen()
}

var binaryOption *bool
var writer *bufio.Writer

func main() {
	binaryOption = flag.Bool("b", false, "Specify output as binary format")
	ff := flag.String("f", "", "Specify assembly")
	of := flag.String("o", "", "Specify output")
	flag.Parse()

	f, err := os.Open(*ff)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	o, err := os.OpenFile(*of, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer o.Close()
	writer = bufio.NewWriter(o)

	scanner := bufio.NewScanner(f)

	writer.WriteString("v3.0 hex words addressed\n")
	writer.WriteString("0:")
	c := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, ";") {
			continue
		}
		gen(line)
		c++
	}
	if c < 15 {
		writer.WriteString(strings.Repeat(" 00", 16-c) + "\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
