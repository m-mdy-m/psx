package utils

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/m-mdy-m/psx/internal/logger"
)

var stdin = bufio.NewReader(os.Stdin)

func Prompt(message string) bool {
	logger.Infof("%s [y/N]:\n> ", message)
	line, _ := stdin.ReadString('\n')
	resp := strings.ToLower(strings.TrimSpace(line))
	return resp == "y" || resp == "yes"
}

func PromptChoice(message string, options []string) (string, error) {
	logger.Infof("%s\n", message)
	for i, opt := range options {
		logger.Infof("  %d) %s\n", i+1, opt)
	}
	logger.Info("Choose [1]:\n> ")

	line, _ := stdin.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return options[0], nil
	}

	c, err := strconv.Atoi(line)
	if err != nil || c < 1 || c > len(options) {
		return options[0], nil
	}
	return options[c-1], nil
}

func PromptInput(message string, defaultValue string) string {
	if defaultValue != "" {
		logger.Infof("%s [%s]:\n> ", message, defaultValue)
	} else {
		logger.Infof("%s:\n> ", message)
	}

	line, _ := stdin.ReadString('\n')
	input := strings.TrimSpace(line)
	if input == "" {
		return defaultValue
	}
	return input
}
