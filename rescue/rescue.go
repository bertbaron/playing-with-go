package main

import (
	"fmt"
	"os"
	"log"
	"bufio"
)

func main() {
	fmt.Println("Hello world")
	
	file, err := os.Open("/home/bbaron/codeeval/examples/rescue.example33")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
    	fmt.Println(scanner.Text())
	}
}