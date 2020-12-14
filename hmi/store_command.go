package hmi

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"../crypto"
)

func init() {
	register(newStoreCommand())
}

type StoreCommand struct {
	pattern *regexp.Regexp
}

func (command *StoreCommand) Execute(cryptoStorage crypto.CryptoStorage, scanner *bufio.Scanner, input string, password []byte) bool {
	trimmedPassword := strings.TrimSpace(input)
	patternGroups := command.pattern.FindStringSubmatch(trimmedPassword)
	if patternGroups == nil {
		return false
	}

	content, err := command.readContentWithConfirmation(scanner)
	if err != nil {
		fmt.Println(err)
		return true
	}

	err = cryptoStorage.SaveContent(password, patternGroups[1], content)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Store %s... done!\n", patternGroups[1])
	}

	return true
}

func (command *StoreCommand) GetPriority() int {
	return 0
}

func (command *StoreCommand) readContentWithConfirmation(scanner *bufio.Scanner) ([]byte, error) {
	fmt.Println("And now content please")
	setSecureTerminalMode(true)
	scanner.Scan()
	content := scanner.Text()
	setSecureTerminalMode(false)
	fmt.Println("Once again")
	setSecureTerminalMode(true)
	scanner.Scan()
	confirmedContent := scanner.Text()
	setSecureTerminalMode(false)

	if content == confirmedContent {
		return []byte(content), nil
	}
	return nil, fmt.Errorf("Miss")
}

func newStoreCommand() *StoreCommand {
	pattern, _ := regexp.Compile("(?i)^store (\\S+)$")
	return &StoreCommand{pattern}
}
