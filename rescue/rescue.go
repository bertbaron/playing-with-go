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

// bit masks for P, U and E in the byte matrix
const (
	P = 1 << iota
	U = 1 << iota
	E = 1 << iota
)

// height, width and number of people, unpassable objects and exits
var n, m, p, u, e int32

// slice of n*m holding the board
var matrix []byte

// set of people positions
var ps map[int32]bool

// statistics
var ops = 0

//Map from P node to a map of reachable P nodes with the path between them
var node2nodes map[int32]map[int32][]int32 = make(map[int32]map[int32][]int32)

//Map from node to path to nearest exit
var node2exit map[int32][]int32 = make(map[int32][]int32)

//Map from exit to path to nearest
//var exit2node map[int32][]int32 = make(map[int32][]int32)

func pack(x, y int32) int32 {
	return x*m + y
}

func unpack(i int32) (x, y int32) {
	return i / m, i % m
}

func isValid(x, y int32) bool {
	return 0 <= x && x < n && 0 <= y && y < m
}

// could not find a library function for this (:
func reverse(slice []int32) []int32 {
	count := len(slice)
	reversed := make([]int32, count, count)
	for i, e := range slice {
		reversed[count-1-i] = e
	}
	return reversed
}

func addPathToMap(from, to int32, path []int32) {
	if _, ok := node2nodes[from]; !ok {
		node2nodes[from] = make(map[int32][]int32)
	}
	nodeMap := node2nodes[from]
	nodeMap[to] = path
}

func extractPaths(from, to int32, parents []int32) (forward, backward []int32) {
	backward = make([]int32, 0)
	for parents[to] >= 0 {
		backward = append(backward, to)
		to = parents[to]
	}
	backward = append(backward, to)
	forward = reverse(backward)
	return forward[1:], backward[1:]
}

func addPath(from, to int32, parents []int32) {
	forward, backward := extractPaths(from, to, parents)
	addPathToMap(from, to, forward)
	addPathToMap(to, from, backward)
}

func addExit(from, to int32, parents []int32) {
	forward, _ := extractPaths(from, to, parents)
	node2exit[from] = forward
//	exit2node[to] = backward
}

func expand(pos int32, buffer *[8]int32) int {
	x, y := unpack(pos)
	count := 0
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if (i != x || j != y) && isValid(i, j) {
				neighbour := pack(i, j)
				if matrix[neighbour]&U == 0 {
					buffer[count] = neighbour
					count++
				}
			}
		}
	}
	return count
}

func calculatePaths(pos int32) {
	x, y := unpack(pos)
	fmt.Printf("Calculating paths from (%d,%d)\n", x, y)
	start := time.Now()

	var buffer [8]int32
	queue := make([]int32, m*n)
	parents := make([]int32, m*n)
	visited := make([]bool, m*n)
	head := 0
	tail := 1
	exitFound := false
	queue[head] = pos
	visited[pos], parents[pos] = true, -1
	for head < tail {
		ops += 1
		p := queue[head]
		head += 1

		if matrix[p]&P != 0 {
			x, y := unpack(p)
			fmt.Printf("  found (%d, %d)\n", x, y)
			addPath(pos, p, parents)
		}
		if !exitFound && matrix[p]&E != 0 {
			x, y := unpack(p)
			fmt.Printf("  found exit (%d, %d)\n", x, y)
			addExit(pos, p, parents)
			exitFound = true
		}

		expanded := expand(p, &buffer)
		for i := 0; i < expanded; i++ {
			nb := buffer[i]
			if !visited[nb] {
				visited[nb], parents[nb] = true, p
				queue[tail] = nb
				tail += 1
			}
		}
	}
	fmt.Printf("Traversed in %s\n", time.Since(start))
}

type graph struct {
    //Map from P node to a map of reachable P nodes with the path between them
	nodes map[int32]map[int32][]int32
	
	//Map from node to path to nearest exit
	node2exit map[int32][]int32
}

// Constructs the graph with nodes and nearest exits.
// Requires the global variables for the input to be initialized.
func constructGraph() []graph {
	for k := range ps {
		calculatePaths(k)
	}
	var subgraphs []graph
	
//	size := 0 // size of biggest reachable subgraph

	// split the graph into subgraphs
	for len(node2nodes) > 0 {
		// wow, its really not trivial to get any single element from a map
		var from int32
		var nodes map[int32][]int32
		for from, nodes = range node2nodes { break }
		
		delete(node2nodes, from)
		g := graph{make(map[int32]map[int32][]int32), make(map[int32][]int32)}
		for to, _ := range nodes {
			g.nodes[from] = nodes
			delete(node2nodes, to)
		}
		subgraphs = append(subgraphs, g)
	}
	return subgraphs
}

func parseInt(s string) int32 {
	x, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return int32(x)
}

func parseTuple(line string) (x, y int32) {
	fields := strings.Fields(line)
	x = parseInt(fields[0])
	y = parseInt(fields[1])
	return
}
func parsePos(line string) int32 {
	return pack(parseTuple(line))
}

func nextLine(scanner *bufio.Scanner) string {
	if scanner.Scan() {
		return scanner.Text()
	}
	panic("No more lines")
}

func parseInput() {
	file, err := os.Open("/home/bert/git/codeeval/examples/rescue.example")
	//	file, err := os.Open("/home/bbaron/codeeval/examples/rescue.huge")
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
	ps = make(map[int32]bool)
	for i := 0; i < int(p); i++ {
		pos := parsePos(nextLine(scanner))
		ps[pos] = true
		matrix[pos] |= P
	}

	// parse unpassable objects
	u = parseInt(nextLine(scanner))
	for i := 0; i < int(u); i++ {
		pos := parsePos(nextLine(scanner))
		matrix[pos] |= U
	}

	// parse exits
	e = parseInt(nextLine(scanner))
	for i := 0; i < int(e); i++ {
		pos := parsePos(nextLine(scanner))
		matrix[pos] |= E
	}

	fmt.Println("")

	elapsed := time.Since(start)
	fmt.Printf("Parsed in %s\n", elapsed)
}

func main() {
	start := time.Now()

	parseInput()
	constructGraph()

	fmt.Printf("%d operations\n", ops)
	fmt.Printf("Total time: %s", time.Since(start))
}
