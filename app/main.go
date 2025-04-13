package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	command := args[0]
	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
	}

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	} else if len(args) == 2 {
		err := runFile(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
	} else {
		runPrompt()
	}
}
func runFile(fileName string) error {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	run(string(bytes))
	return nil
}

func runPrompt() error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		line := scanner.Text()
		run(line)
		fmt.Print("> ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning")
		return err
	}
	return nil
}

func run(source string) error {
	scanner := NewScanner(source)
	scanner.ScanTokens()
	for _, token := range scanner.Tokens {
		fmt.Println(token)
	}
	return nil
}
