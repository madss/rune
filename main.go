package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	i := NewInterpreter()

	fmt.Println("Welcome do the ᚱune interpreter")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">>> ")
		if !scanner.Scan() {
			break
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}
		val, err := i.Eval(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			continue
		}
		fmt.Printf("--> %v\n", val)
	}
}
