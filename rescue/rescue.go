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
	type node struct {
		p      int
		parent *node
	}
	queue := make([]node, m*n)
	head := 0
	tail := 1
	queue[head] = node{pos, nil}
	for head < tail {
		ops += 1
		path := queue[head]
		head += 1
		expanded := expand(path.p)
		for _, nb := range expanded {
			if matrix[nb] & visited == 0 {
				matrix[nb] |= visited
				queue[tail] = node{nb, &path}
				tail += 1
			}
		}
	}
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

	expanded := expand(pack(1, 0))
	fmt.Printf("Expand: ")
	for _, p := range expanded {
		x, y := unpack(p)
		fmt.Printf(" (%d,%d)", x, y)
	}
	fmt.Println("")
	calculatePaths(0)
	fmt.Printf("%d operations", ops)
	fmt.Println("Done...")
}
