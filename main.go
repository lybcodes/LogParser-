package main

import (
	"bufio"
	"fmt"
	"os"

	"logparser/parser"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	parser := parser.NewLogParser()

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parsed, err := parser.Parse(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "解析错误: %v\n", err)
			continue
		}

		fmt.Println(parsed.String())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "读取错误: %v\n", err)
		os.Exit(1)
	}
}
