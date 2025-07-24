package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	prompt := "Pokedex > "
	fmt.Print(prompt)
	for scanner.Scan() {
		input := cleanInput(scanner.Text())
		fmt.Println("Your command was:", input[0])
		fmt.Print(prompt)
	}
}

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	out := strings.Fields(text)
	if len(out) == 0 {
		return []string{""}
	}
	return out
}
