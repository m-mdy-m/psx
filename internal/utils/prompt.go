package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var stdin = bufio.NewReader(os.Stdin)

func Prompt(message string) bool {
	fmt.Printf("%s [y/N]: ", message)
	line, _ := stdin.ReadString('\n')
	resp := strings.ToLower(strings.TrimSpace(line))
	return resp == "y" || resp == "yes"
}
func PromptChoice(message string, options []string) (string, error) {
	fmt.Println(message)
	for i, opt := range options {
		fmt.Printf("  %d) %s\n", i+1, opt)
	}
	fmt.Print("Choose [1]: ")

	line, _ := stdin.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return options[0], nil
	}

	choice, err := strconv.Atoi(line)
	if err != nil || choice < 1 || choice > len(options) {
		return options[0], nil
	}
	return options[choice-1], nil
}
func PromptInput(message string, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", message, defaultValue)
	} else {
		fmt.Printf("%s: ", message)
	}

	line, _ := stdin.ReadString('\n')
	input := strings.TrimSpace(line)
	if input == "" {
		return defaultValue
	}
	return input
}
