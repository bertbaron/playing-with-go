package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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

var ops = 0

func pack(x, y int) int {
	return x*m + y
}

func unpack(i int) (x, y int) {
	x = i / m
	y = i % m
	return
}

func isValid(x, y int) bool {
	return 0 <= x && x < n && 0 <= y && y < m
}

func expand(pos int) []int {
	x, y := unpack(pos)
	expanded := make([]int, 0, 8)
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if (i != x || j != y) && isValid(i, j) {
				neighbour := pack(i, j)
				if matrix[neighbour] & U == 0 {
					expanded = append(expanded, neighbour)
				}
			}
		}
	}
	return expanded
}

func calculatePaths(pos int) {
	x, y := unpack(pos)
	fmt.Printf("Calculating paths from (%d,%d)\n", x, y)
	start := time.Now()

	queue := make([]int, m*n)
	paths := make([]int, m*n)
	head := 0
	tail := 1

	queue[head] = pos
	paths[head] = -1
	for head < tail {
		ops += 1
		p := queue[head]
		head += 1
		expanded := expand(p)
		for _, nb := range expanded {
			if matrix[nb] & visited == 0 {
				matrix[nb] |= visited
				queue[tail] = nb
				paths[tail] = p
				tail += 1
			}
		}
	}
	elapsed := time.Since(start)
    log.Printf("Traversed in %s\n", elapsed)

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

	start := time.Now()
	scanner := bufio.NewScanner(file)
	n, m = parseTuple(nextLine(scanner))
	fmt.Printf("n=%d, m=%d\n", n, m)
	
	matrix = make([]byte, n*m)
    
	// parse people
	p = parseInt(nextLine(scanner))
	ps = make(map[int]bool)
	for i := 0; i < p; i++ {
		pos := parsePos(nextLine(scanner))
		ps[pos] = true
		matrix[pos] |= P
	}

	// parse unpassable objects
	u = parseInt(nextLine(scanner))
	for i := 0; i < u; i++ {
		pos := parsePos(nextLine(scanner))
		matrix[pos] |= U
	}

	// parse exits
	e = parseInt(nextLine(scanner))
	for i := 0; i < e; i++ {
		pos := parsePos(nextLine(scanner))
		matrix[pos] |= E
	}

	fmt.Println("")

	elapsed := time.Since(start)
    log.Printf("Parsed in %s", elapsed)

	for k := range ps {
		calculatePaths(k)
	}
	
	fmt.Printf("%d operations\n", ops)
	fmt.Printf("Total time: %s", time.Since(start))
}
