package hmi

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"../crypto"
)

func init() {
	register(newShowCommand())
}

type ShowCommand struct {
	pattern *regexp.Regexp
}

func (command *ShowCommand) Execute(cryptoStorage crypto.CryptoStorage, scanner *bufio.Scanner, input string, password []byte) bool {
	trimmedPassword := strings.TrimSpace(input)
	patternGroups := command.pattern.FindStringSubmatch(trimmedPassword)
	if patternGroups == nil {
		return false
	}

	content, err := cryptoStorage.GetContent(password, patternGroups[1])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(content))
		fmt.Printf("Show %s... done!\n", patternGroups[1])
	}

	return true
}

func (command *ShowCommand) GetPriority() int {
	return 0
}

func newShowCommand() *ShowCommand {
	pattern, _ := regexp.Compile("(?i)^show (\\S+)$")
	return &ShowCommand{pattern}
}
