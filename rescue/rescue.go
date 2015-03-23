package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	P       = 1 << iota
	U       = 1 << iota
	E       = 1 << iota
	visited = 1 << iota
)

// height, width and number of people, unpassable objects and exits
var n, m, p, u, e int

// slice of n*m holding the board
var matrix []byte

// set of people positions
var ps map[int]bool

func pack(x, y int) int {
	return x*m + y
}

func unpack(i int) (x, y int) {
	x = i / m
	y = i % m
	return
}

func parseInt(s string) int {
	x, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return x
}

func parseTuple(line string) (x, y int) {
	fields := strings.Fields(line)
	x = parseInt(fields[0])
	y = parseInt(fields[1])
	return
}
func parsePos(line string) int {
	return pack(parseTuple(line))
}

func nextLine(scanner *bufio.Scanner) string {
	if scanner.Scan() {
		return scanner.Text()
	}
	panic("No more lines")
}

func main() {
	file, err := os.Open("/home/bert/git/codeeval/examples/rescue.huge")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	n, m = parseTuple(nextLine(scanner))
	fmt.Printf("n=%d, m=%d\n", n, m)
	matrix = make([]byte, n*m)

	// parse people
	p = parseInt(nextLine(scanner))
	ps = make(map[int]bool)
	for i:=0 ; i<p; i++ {
		pos := parsePos(nextLine(scanner))
		ps[pos] = true
		matrix[pos] &= P
	}

	// parse unpassable objects
	u = parseInt(nextLine(scanner))
	for i:=0 ; i<u; i++ {
		pos := parsePos(nextLine(scanner))
		matrix[pos] &= U
	}

	// parse exits
	e = parseInt(nextLine(scanner))
	for i:=0 ; i<e; i++ {
		pos := parsePos(nextLine(scanner))
		matrix[pos] &= E
	}
	
	fmt.Println("Done...")
}
