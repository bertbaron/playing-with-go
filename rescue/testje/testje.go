package main

import (
	"fmt"
)


func main() {
	var m map[int]map[int]bool
	fmt.Printf("map=%v\n", m)
	x := m[1]
	fmt.Printf("map=%v, x=%s\n", m, x)
	x[3]=true
	fmt.Printf("map=%v, x=%v\n", m, x)
}