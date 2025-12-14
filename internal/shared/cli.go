package shared

import(
	"fmt"
	"strings"
)
func Prompt(message string) bool {
	fmt.Printf("%s [y/N]: ", message)
	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

func PromptChoice(message string, options []string) (string, error) {
	fmt.Println(message)
	for i, opt := range options {
		fmt.Printf("  %d) %s\n", i+1, opt)
	}
	fmt.Print("Choose [1]: ")

	var choice int
	if _, err := fmt.Scanln(&choice); err != nil {
		return options[0], nil // Default to first option
	}

	if choice < 1 || choice > len(options) {
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

	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}

	return input
}
