package main

import "fmt"

func main() {
	l := lexer{
		input: "ab cd",
		pos:   0,
		start: 0,
		width: len("ab cd"),
		line:  0,
	}

	for l.peek() != eof {
		fmt.Println(l.next())
	}
}
