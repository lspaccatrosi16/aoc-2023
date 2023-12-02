package main

import (
	"fmt"

	"example.com/my/project/lib/challenge"
	"example.com/my/project/lib/ipt"
)

func main() {
	input := ipt.GetInput()
	output := challenge.Challenge(input)

	fmt.Println(output)
}
