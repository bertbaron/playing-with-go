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
//	visited = 1 << iota
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

func expand(pos int, buffer *[8]int) int {
	x, y := unpack(pos)
	count := 0
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if (i != x || j != y) && isValid(i, j) {
				neighbour := pack(i, j)
				if matrix[neighbour] & U == 0 {
					buffer[count] = neighbour
					count++
				}
			}
		}
	}
	return count
}

//Map from P node to a map of reachable P nodes with the path between them
var paths map[int]map[int][]int = make(map[int]map[int][]int)

func addPathToMap(from, to int, path []int) {
	if _, ok := paths[from]; !ok {
		paths[from]=make(map[int][]int)
	}
	nodeMap := paths[from]
	nodeMap[to] = path
}

func addPath(from, to int, links []int) {
//	fmt.Printf("  from %d to %d\n", from, to)
	path := make([]int, 0)
	for links[to] >= 0 {
		path = append(path, to)
		to = links[to]
	}
	path = append(path, to)
	backward := path[1:]
	count := len(path)
	forward := make([]int, count, count)
	for i, point := range path {
		forward[count-1-i] = point
	}
	forward = forward[1:]
//	fmt.Println("  forward=%s", forward)
//	fmt.Println("  backward=%s", backward)
	addPathToMap(from, to, forward)
	addPathToMap(to, from, backward)
}

func calculatePaths(pos int) {
	x, y := unpack(pos)
	fmt.Printf("Calculating paths from (%d,%d)\n", x, y)
	start := time.Now()

	var buffer [8]int
	queue := make([]int, m*n)
	parents := make([]int, m*n)
	visited := make([]bool, m*n)
	head := 0
	tail := 1
	exitFound := false
	queue[head] = pos
	parents[pos] = -1
	visited[pos] = true
	for head < tail {
		ops += 1
		p := queue[head]
		head += 1
		
		if matrix[p] & P != 0 {
			x, y := unpack(p)
			fmt.Printf("  found (%d, %d)\n", x, y)
			addPath(pos, p, parents)
		}
		if !exitFound && matrix[p] & E != 0 {
			x, y := unpack(p)
			fmt.Printf("  found exit (%d, %d)\n", x, y)
			exitFound = true
		}
		
		expanded := expand(p, &buffer)
		for i := 0; i < expanded; i++ {
			nb := buffer[i]
			if !visited[nb] {
				visited[nb] = true
				if parents[nb] != 0 { panic("there is already a parent configured") }
				parents[nb] = p
				queue[tail] = nb
				tail += 1
			}
		}
	}
    fmt.Printf("Traversed in %s\n", time.Since(start))
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
//	file, err := os.Open("/home/bert/git/codeeval/examples/rescue.huge")
	file, err := os.Open("/home/bbaron/codeeval/examples/rescue.huge")
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
    fmt.Printf("Parsed in %s\n", elapsed)

	for k := range ps {
		calculatePaths(k)
	}

//	fmt.Printf("paths: %s", paths)
	
	fmt.Printf("%d operations\n", ops)
	fmt.Printf("Total time: %s", time.Since(start))
}
