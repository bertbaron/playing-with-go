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

type point int32
type points []point

func (p point) String() string {
	x, y := unpack(p)
	return fmt.Sprintf("(%d,%d)", x, y)
}

type graph struct {
	//Map from P node to a map of reachable P nodes with the path between them
	node2nodes map[point]map[point]points

	//Map from node to path to nearest exit (last element is the exit)
	node2exit map[point]points
}

func newGraph() graph {
	return graph{make(map[point]map[point]points), make(map[point]points)}
}

func (g graph) pprint() {
	fmt.Println("graph nodes:")
	for from, paths := range g.node2nodes {
		fmt.Println("  ", from)
		for to, path := range paths {
			fmt.Printf("    %s : %d\n", to, len(path))
		}
	}
	fmt.Println("nearest exits:")
	for node, path := range g.node2exit {
		exit := node
		if len(path) > 0 {
			exit = path[len(path)-1]
		}
		fmt.Printf("  %s -> %s: %d\n", node, exit, len(path))
	}
}

// height, width and number of people, unpassable objects and exits
var n, m, p, u, e int32

// slice of n*m holding the board
var matrix []byte

// set of people positions
var ps map[point]bool

// statistics
var ops = 0

// the graph, build up during the calculation phase.
var root = newGraph()

// sets of connected nodes, mapping to true if they have an exit or false if not
var nodesets []map[point]bool

// logs the time since the given start
func timed(msg string, start time.Time) {
	fmt.Println(msg, time.Since(start))
}

func pack(x, y int32) point {
	return point(x*m + y)
}

func unpack(i point) (x, y int32) {
	return int32(i) / m, int32(i) % m
}

func isValid(x, y int32) bool {
	return 0 <= x && x < n && 0 <= y && y < m
}

// could not find a library function for this (:
func reverse(slice points) points {
	count := len(slice)
	reversed := make(points, count, count)
	for i, e := range slice {
		reversed[count-1-i] = e
	}
	return reversed
}

func addPathToMap(from, to point, path points) {
	if _, ok := root.node2nodes[from]; !ok {
		root.node2nodes[from] = make(map[point]points)
	}
	nodeMap := root.node2nodes[from]
	nodeMap[to] = path
}

func extractPaths(from, to point, parents points) (forward, backward points) {
	backward = make(points, 0)
	for parents[to] >= 0 {
		backward = append(backward, to)
		to = parents[to]
	}
	backward = append(backward, to)
	forward = reverse(backward)
	return forward[1:], backward[1:]
}

func addPath(from, to point, parents points) {
	forward, backward := extractPaths(from, to, parents)
	addPathToMap(from, to, forward)
	addPathToMap(to, from, backward)
}

func addExit(from, to point, parents points) {
	forward, _ := extractPaths(from, to, parents)
	root.node2exit[from] = forward
	//	exit2node[to] = backward
}

// indicates that all necesary paths from the specified point have been calculated
func completed(pos point) {
	for _, set := range nodesets {
		if set[pos] { return }
	}
	set := make(map[point]bool)
	for node, _ := range root.node2nodes[pos] {
		set[node] = true
	}
}
// Performs a breadth-first search from the given position, finding all reachable
// nodes and the nearest exit
func calculatePaths(pos point) {
	x, y := unpack(pos)
	fmt.Printf("Calculating paths from (%d,%d)\n", x, y)
	start := time.Now()
	defer timed("Traversed in", start)

	queue := make(points, m*n)
	parents := make(points, m*n)
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

		x, y := unpack(p)
		for nx := x - 1; nx <= x+1; nx++ {
			for ny := y - 1; ny <= y+1; ny++ {
				if (nx != x || ny != y) && isValid(nx, ny) {
					nb := pack(nx, ny)
					if matrix[nb]&U == 0 && !visited[nb] {
						visited[nb], parents[nb] = true, p
						queue[tail] = nb
						tail += 1
					}
				}
			}
		}
	}
	
	completed(pos)
}

// Constructs the relevant graph with nodes and nearest exits.
func constructGraph(root graph) graph {
	for k := range ps {
		calculatePaths(k)
	}
	root.pprint()

	var subgraphs []graph

	size := 0 // size of biggest reachable subgraph

	// split the graph into subgraphs, forget about subgraphs with no exit
	for len(root.node2nodes) > 0 {
		// get any map of target nodes with their pahts from the root graph
		var nodes map[point]points
		for _, nodes = range root.node2nodes {
			break
		}

		g := newGraph()

		for to, _ := range nodes {
			g.node2nodes[to] = root.node2nodes[to]
			if exit, ok := root.node2exit[to]; ok {
				g.node2exit[to] = exit
			}
			delete(root.node2nodes, to)
		}
		if len(g.node2exit) > 0 {
			// only remember subgraphs with exit
			subgraphs = append(subgraphs, g)
			if len(nodes) > size {
				size = len(nodes)
			}
		}
	}
	biggest := newGraph()
	for _, g := range subgraphs {
		if len(g.node2nodes) == size {
			for k, v := range g.node2nodes {
				biggest.node2nodes[k] = v
			}
			for k, v := range g.node2exit {
				biggest.node2exit[k] = v
			}
		}
	}
	biggest.pprint()
	return biggest
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
func parsePos(line string) point {
	return pack(parseTuple(line))
}

func nextLine(scanner *bufio.Scanner) string {
	if scanner.Scan() {
		return scanner.Text()
	}
	panic("No more lines")
}

func parseInput() {
//	file, err := os.Open("/home/bert/git/codeeval/examples/rescue.large")
	file, err := os.Open("/home/bbaron/codeeval/examples/rescue.large")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	start := time.Now()
	defer timed("Parsed in", start)

	scanner := bufio.NewScanner(file)
	n, m = parseTuple(nextLine(scanner))
	fmt.Printf("n=%d, m=%d\n", n, m)

	matrix = make([]byte, n*m)

	// parse people
	p = parseInt(nextLine(scanner))
	fmt.Printf("p=%d", p)
	ps = make(map[point]bool)
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
}

func main() {
	start := time.Now()
	defer timed("Total time:", start)
	
	parseInput()
	constructGraph(root)

	fmt.Printf("%d operations\n", ops)
}
