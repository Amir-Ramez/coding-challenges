package main

import (
	"fmt"
	"strings"
)

func main() {
	input := `{"name" : "amir" , "key" : "value"}`
	p := NewParser(strings.NewReader(input))
	err := p.Parse()
	fmt.Println(err)
}
