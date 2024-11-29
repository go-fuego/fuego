package main

import (
	"fmt"
	"os/exec"
)

func main() {
	output, _ := exec.Command("go run").Output()
	fmt.Println(output)
}
